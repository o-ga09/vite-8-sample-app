# マルチテナント 家計簿システム — 実装プラン

## Context

layerx.go #4 イベント向けのデモアプリ。bob ORM の主要機能（CRUD、JOIN、型安全クエリビルド、queries plugin = sqlc記法、テナントフック）を PostgreSQL 上で実演する。
ogen でスキーマ駆動の REST API を生成し、Go 1.24+ の `go tool` を活用する。

家計簿システムとして、ワークスペース（家計単位）ごとに口座・カテゴリ・取引を管理し、収支レポートを提供する。

**追加要件**:
- **Repository パターン**: データアクセス層を分離し、Service層はビジネスロジックに専念
- **トランザクション管理**: Service層で `WithTx()` パターンによる明示的なトランザクション境界
- **golangci-lint**: Go コードの静的解析（gofmt, govet, staticcheck 等）
- **CI/CD**: GitHub Actions による自動 lint, test, build

## アーキテクチャ概要

```
OpenAPI spec (api/openapi.yaml)
    ↓ ogen generate
internal/oas/  (生成コード)
    ↓ implements Handler interface
internal/handler/  (ハンドラ実装)
    ↓ calls
internal/service/  (ビジネスロジック + トランザクション管理)
    ↓ calls
internal/repository/  (データアクセス層、bob クエリ構築)
    ↓ uses
internal/infra/db/  (DB接続, WorkspaceScoped/Global executor, Transaction)
    ↓ hooks fire
internal/infra/dbgen/  (bob 生成モデル・where・joins・loaders)
internal/infra/dbgen/dbenums/
internal/infra/dbgen/factory/
queries/  (sqlc記法 SQL → bob queries plugin 生成)
```

アーキテクチャ設計のポイント:
- テナントキーは `workspace_id`（organization_id / tenant_id は使わない）
- **Repository パターン採用**: Service層はビジネスロジック、Repository層はデータアクセスに専念
- **トランザクション管理**: Service層で `db.WithTx()` によるトランザクション境界を明示
- ScopedExecutor ではなく **WorkspaceScopedExec / GlobalExec** という命名
- Hook 登録は手動の `RegisterHooks()` 関数（custom-bobgen 不要）
- QueryHooks ベース（ContextualMod ではない）。デモ用の簡略化であり、本番では `bob.SkipHooks()` でバイパスされるリスクがある点はトレードオフとして認識

## テーブル設計（5テーブル）

| テーブル | パーティションキー | 主な関連 |
|---|---|---|
| workspaces | — (テナント自体) | — |
| members | workspace_id | workspaces |
| accounts | workspace_id | workspaces |
| categories | workspace_id | workspaces |
| transactions | workspace_id | accounts(from/to), categories, members(created_by) |

### accounts（口座）
- `id` UUID PRIMARY KEY
- `workspace_id` UUID NOT NULL
- `name` VARCHAR(100) NOT NULL
- `account_type` account_type ENUM NOT NULL
- `initial_balance` DECIMAL(15,2) NOT NULL DEFAULT 0
- `description` TEXT
- `is_active` BOOLEAN NOT NULL DEFAULT true
- `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
- `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP

### categories（カテゴリ）
- `id` UUID PRIMARY KEY
- `workspace_id` UUID NOT NULL
- `name` VARCHAR(100) NOT NULL
- `category_type` category_type ENUM NOT NULL
- `parent_id` UUID (self FK, nullable, 階層化用)
- `description` TEXT
- `is_active` BOOLEAN NOT NULL DEFAULT true
- `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
- `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP

### transactions（取引）
- `id` UUID PRIMARY KEY
- `workspace_id` UUID NOT NULL
- `transaction_type` transaction_type ENUM NOT NULL
- `amount` DECIMAL(15,2) NOT NULL
- `transaction_date` DATE NOT NULL
- `account_id` UUID (FK accounts, nullable, 収入/支出の場合に使用)
- `from_account_id` UUID (FK accounts, nullable, 振替の場合に使用)
- `to_account_id` UUID (FK accounts, nullable, 振替の場合に使用)
- `category_id` UUID (FK categories, nullable, 振替はNULL)
- `description` TEXT
- `created_by` UUID NOT NULL (FK members)
- `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
- `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP

**制約**:
- 収入/支出: `account_id` NOT NULL, `from_account_id` / `to_account_id` NULL, `category_id` NOT NULL
- 振替: `from_account_id` / `to_account_id` NOT NULL, `account_id` NULL, `category_id` NULL
- CHECK制約で上記を保証

ENUMs:
- `member_role` (owner/editor/viewer)
- `account_type` (cash/bank/credit_card/e_money/investment)
- `transaction_type` (income/expense/transfer)
- `category_type` (income/expense)

全テーブルに `updated_at` 自動更新トリガーを設定。

**UUID 方針**: ADR に従い UUIDv7 をアプリ側で生成（`github.com/google/uuid` の `uuid.NewV7()`）。DB の `DEFAULT gen_random_uuid()` は設定しない。これにより時系列ソート可能な ID を実現。

**cross-workspace 参照防止**: テナントスコープのテーブルには `UNIQUE(workspace_id, id)` 制約を追加し、FK は複合FK `(workspace_id, xxx_id) REFERENCES xxx(workspace_id, id)` で定義。これにより DB レベルで workspace を跨いだ参照を防止。

**Insert 時の workspace_id 強制**: Service 層で path パラメータの workspace_id を setter に上書き固定する。リクエストbodyからは受け取らない。

**通貨**: 円（JPY）のみ。複数通貨対応は将来の拡張として保留。

**将来の拡張可能性**（現時点では実装しない）:
- 添付ファイル（レシート画像）
- 予算管理（カテゴリ別月次予算・アラート）
- 定期取引（サブスクリプション、給与等）

## Hooks と queries plugin の関係（重要な知見）

bob の queries plugin で生成されるコードは `ExecQuery` の `Hooks` フィールドが **nil** のまま生成される（`gen/templates/queries/query/01_query.go.tpl` L86-91 で確認済み）。
よって **テナントフック（QueryHooks）は queries plugin 出力には効かない**。

対処方針:
- **CRUD（bob 標準モデル経由）**: フックが自動適用 → `workspace_id` フィルタ自動挿入
- **集計クエリ（queries plugin）**: SQL 自体に `WHERE workspace_id = $1` を明示的に書く（フックに頼らない）
- これは layerone が sqlc で採用しているのと同じ戦略

## ディレクトリ構成

pnpm workspace でモノレポ構成。ルートに Go バックエンド、`webapp/` にフロントエンド。

```
.
├── pnpm-workspace.yaml             # packages: [webapp]
├── package.json                     # root（pnpm, scripts）
├── api/
│   └── openapi.yaml              # OpenAPI 3.0 spec（バックエンド・フロントエンド共有）
├── cmd/
│   ├── server/main.go            # HTTPサーバ起動
│   └── seed/main.go              # サンプルデータ投入
├── db/
│   └── migrations/               # dbmate マイグレーション（7本）
├── queries/
│   └── reports.sql               # sqlc記法の集計クエリ
├── internal/
│   ├── oas/                      # [生成] ogen サーバコード
│   ├── handler/                  # ogen Handler 実装
│   │   ├── handler.go            # 共通構造体・エラーハンドリング
│   │   ├── workspace.go
│   │   ├── member.go
│   │   ├── account.go
│   │   ├── category.go
│   │   ├── transaction.go
│   │   └── report.go
│   ├── service/                  # ビジネスロジック（トランザクション管理）
│   │   ├── workspace.go
│   │   ├── member.go
│   │   ├── account.go
│   │   ├── category.go
│   │   ├── transaction.go
│   │   └── report.go
│   ├── repository/               # データアクセス層（bob操作）
│   │   ├── workspace.go
│   │   ├── member.go
│   │   ├── account.go
│   │   ├── category.go
│   │   ├── transaction.go
│   │   └── report.go            # queries plugin 生成関数を呼ぶ
│   └── infra/
│       ├── db/
│       │   ├── conn.go           # pgx 接続 + bob.NewDB
│       │   ├── scope.go          # WorkspaceScopedExec / GlobalExec
│       │   └── tx.go             # トランザクション管理（WithTx）
│       ├── dbgen/                # [生成] bob モデル
│       │   ├── dbenums/          # [生成] enum型
│       │   └── factory/          # [生成] テストファクトリ
│       └── hook/
│           ├── workspace.go      # WorkspaceSelectHook / Update / Delete
│           └── register.go       # RegisterHooks()
├── webapp/                         # Next.js フロントエンド（Vite+ 統合）
│   ├── package.json
│   ├── vite.config.ts              # Vite+ 設定（oxlint/oxfmt は自動）
│   ├── next.config.ts
│   ├── tsconfig.json               # strict 設定
│   ├── tailwind.config.ts
│   ├── components.json             # shadcn/ui config
│   ├── src/
│   │   ├── app/                    # Next.js App Router
│   │   │   ├── layout.tsx
│   │   │   ├── page.tsx            # ワークスペース選択
│   │   │   └── workspaces/
│   │   │       └── [wsId]/
│   │   │           ├── layout.tsx  # サイドバー付きレイアウト
│   │   │           ├── page.tsx    # ダッシュボード
│   │   │           ├── transactions/
│   │   │           ├── accounts/
│   │   │           ├── categories/
│   │   │           ├── members/
│   │   │           └── reports/
│   │   ├── components/
│   │   │   ├── ui/                 # [生成] shadcn/ui コンポーネント
│   │   │   └── ...                 # アプリ固有コンポーネント
│   │   ├── lib/
│   │   │   ├── api/                # [生成] openapi-fetch クライアント
│   │   │   └── utils.ts            # tailwind-merge 等
│   │   └── .env.local              # NEXT_PUBLIC_API_URL
│   └── .env                        # デフォルト env
├── mise.toml                       # env定義（ポート、DSN）
├── compose.yaml
├── .golangci.yml                   # golangci-lint 設定
├── .github/
│   └── workflows/
│       └── ci.yml                  # CI/CD（lint, test, build）
├── docs/
│   └── plan.md                     # 本ドキュメント
├── docker/
│   ├── grafana/
│   │   └── datasource.yaml
│   └── promtail/
│       └── config.yaml
├── bobgen.yaml
├── go.mod
└── go.sum
```

## 実装順序（チーム分担前提）

### Phase 0: インフラ・スキーマ基盤
1. `mise.toml` — ツールチェーン（go 1.26, node 22, pnpm 10）+ env 定義（ポート、DSN等）。`mise trust` で有効化。`mise.local.toml` は既存を維持（GH_TOKEN 等）
1b. `compose.yaml` — PostgreSQL 17 + Loki + Grafana + Promtail。全ポートを env で変更可能にし既存サービスと干渉しない（デフォルト: PG=15432, Grafana=13000, Loki=13100）
2. `mise` — up, down, migrate, bobgen, ogen, lint, test, build, seed, clean ターゲット
3. `go.mod` — tool ディレクティブ追加（bobgen-psql, dbmate, ogen, golangci-lint）+ 依存取得
4. `.golangci.yml` — golangci-lint 設定（gofmt, goimports, govet, staticcheck, errcheck, gosec, revive 等）
5. `.github/workflows/ci.yml` — GitHub Actions CI（lint, test, build, frontend）
6. `db/migrations/` — 全マイグレーションファイル（7本）
   - `001_create_enums.sql` — 4つのenum型定義
   - `002_create_workspaces.sql` — workspaces テーブル
   - `003_create_members.sql` — members テーブル
   - `004_create_accounts.sql` — accounts テーブル
   - `005_create_categories.sql` — categories テーブル（parent_id self FK含む）
   - `006_create_transactions.sql` — transactions テーブル（CHECK制約含む）
   - `007_create_triggers.sql` — updated_at 自動更新トリガー
7. `bobgen.yaml` — プラグイン設定
8. `queries/reports.sql` — 集計SQL（家計簿向け3クエリ）
9. **実行**: `make up && make migrate && make bobgen` で生成コード確認、`make lint` で静的解析確認

### Phase 1: DB層・フック
10. `internal/infra/db/conn.go` — DB接続
11. `internal/infra/db/scope.go` — WorkspaceScopedExec / GlobalExec
12. `internal/infra/db/tx.go` — トランザクション管理（WithTx パターン）
13. `internal/infra/hook/workspace.go` — 3つのQueryHook実装
14. `internal/infra/hook/register.go` — フック登録

### Phase 2: Repository 層
15. `internal/repository/workspace.go` — CRUD (GlobalExec使用)
16. `internal/repository/member.go` — CRUD
17. `internal/repository/account.go` — CRUD + 残高計算クエリ
18. `internal/repository/category.go` — CRUD + 階層取得
19. `internal/repository/transaction.go` — CRUD + JOIN(Preload/ThenLoad) + フィルタクエリ
20. `internal/repository/report.go` — queries plugin 生成関数を呼ぶ

### Phase 3: OpenAPI・ogen
21. `api/openapi.yaml` — 全エンドポイント定義
22. **実行**: `make ogen` で `internal/oas/` 生成

### Phase 4: Service 層
23. `internal/service/workspace.go` — ビジネスロジック（Repository呼び出し）
24. `internal/service/member.go` — ビジネスロジック
25. `internal/service/account.go` — ビジネスロジック + トランザクション（残高整合性）
26. `internal/service/category.go` — ビジネスロジック
27. `internal/service/transaction.go` — ビジネスロジック + トランザクション（振替の原子性保証）
28. `internal/service/report.go` — Repository 呼び出し

### Phase 5: Handler 層
29. `internal/handler/handler.go` — 共通構造体・エラーハンドリング
30. `internal/handler/workspace.go`, `member.go`, `account.go`, `category.go`, `transaction.go`, `report.go` — ogen 型変換 + service呼び出し

### Phase 6: サーバ起動・シード
31. `cmd/server/main.go` — wiring（CORS 設定含む: webapp のポートを許可）
32. `cmd/seed/main.go` — factory プラグインでサンプルデータ（口座・カテゴリ・取引）
33. Grafana ダッシュボード・datasource設定

### Phase 7: フロントエンド基盤
34. `pnpm-workspace.yaml` + ルート `package.json` — モノレポ設定
35. `webapp/` — Vite+ 環境で Next.js 初期化（App Router, TypeScript, Tailwind）
36. `webapp/tsconfig.json` — strict 設定（ncr-orchestrator 参考: strict, noUnusedLocals, noUnusedParameters, noFallthroughCasesInSwitch）
37. `webapp/vite.config.ts` — Vite+ 設定（oxlint/oxfmt は Vite+ にバンドル済み）
38. shadcn/ui 初期化（`vp dlx shadcn@latest init`）→ button, card, table, badge, dialog, input, select, dropdown-menu, calendar, date-picker 等

**Vite+ によるリント・フォーマット**:
- **リント**: `vp lint` — Vite+ バンドルの oxlint を実行（ESLint 不要）
- **フォーマット**: `vp fmt` — Vite+ バンドルの oxfmt を実行（Prettier 不要）
- **型チェック**: `vp check` — oxlint + oxfmt + TypeScript の一括チェック

**重要**: ESLint/Prettier は**インストール不要**。Vite+ が oxlint/oxfmt をバンドル。
- `eslint.config.mjs`, `.prettierrc.mjs` は作成しない
- `package.json` の devDependencies に `eslint`, `prettier` を追加しない

### Phase 8: OpenAPI クライアント生成 + 画面実装
39. `vp dlx openapi-typescript ../api/openapi.yaml -o src/lib/api/schema.d.ts` → OpenAPI クライアント型生成
40. `openapi-fetch` で型安全な API クライアント（`webapp/src/lib/api/client.ts`）
41. 画面実装:
    - `/` — ワークスペース一覧・選択
    - `/workspaces/[wsId]` — ダッシュボード（月次収支、口座残高サマリ、カテゴリ別支出グラフ）
    - `/workspaces/[wsId]/transactions` — 取引一覧（日付範囲、カテゴリ、口座フィルタ）+ CRUD
    - `/workspaces/[wsId]/accounts` — 口座一覧（残高表示）+ CRUD
    - `/workspaces/[wsId]/categories` — カテゴリ一覧（階層表示）+ CRUD
    - `/workspaces/[wsId]/members` — メンバー一覧 + CRUD
    - `/workspaces/[wsId]/reports` — 月次推移、カテゴリ別集計、口座残高推移

### Phase 9: ドキュメント
42. 本ドキュメント（`docs/plan.md`）を最新状態に維持

## チーム分担案

TeamCreate で `kakeibo` チームを作成し、以下のように Agent を spawn して分担:

| エージェント | 担当 | Phase | 依存 |
|---|---|---|---|
| **infra** | mise.toml, compose.yaml, mise, go.mod, .golangci.yml, .github/workflows/ci.yml, migrations, bobgen.yaml, queries/reports.sql, DB層, フック | 0, 1 | なし（最初に実行） |
| **repository** | Repository層全体（データアクセス層） | 2 | infra 完了後 |
| **api** | openapi.yaml, ogen生成, handler層, CORS設定 | 3, 5 | repository 完了後 |
| **service** | Service層全体（ビジネスロジック、トランザクション管理）、cmd/server, cmd/seed | 4, 6 | repository 完了後（api と並行可） |
| **frontend** | pnpm workspace, Next.js セットアップ（Vite+統合）, tsconfig 堅い設定, shadcn/ui, OpenAPI クライアント生成, 全画面実装 | 7, 8 | api 完了後（openapi.yaml 必要） |

依存関係:
```
infra → repository → (api, service 並行) → server wiring
                          api → frontend
```

**リード（自分）の役割**: TeamCreate → TaskCreate → 各エージェント spawn → 進捗監視 → server wiring 統合 → 検証

## 技術詳細

### mise 主要ターゲット
```mise
.PHONY: up down migrate migrate-down bobgen ogen lint test build seed clean

up:
	docker compose up -d

down:
	docker compose down

migrate:
	go run github.com/amacneil/dbmate/v2 up

migrate-down:
	go run github.com/amacneil/dbmate/v2 down

bobgen:
	go run github.com/stephenafamo/bob/gen/bobgen-psql

ogen:
	go run github.com/ogen-go/ogen/cmd/ogen --target internal/oas --clean api/openapi.yaml

lint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint run

test:
	go test -v -race ./...

build:
	go build -o bin/server ./cmd/server
	go build -o bin/seed ./cmd/seed

seed:
	go run ./cmd/seed

clean:
	rm -rf bin/
	docker compose down -v
```

### go.mod tool ディレクティブ
```go
tool (
    github.com/stephenafamo/bob/gen/bobgen-psql
    github.com/amacneil/dbmate/v2
    github.com/ogen-go/ogen/cmd/ogen
    github.com/golangci/golangci-lint/cmd/golangci-lint
)
```

### .golangci.yml 主要設定
```yaml
run:
  timeout: 5m
  go: '1.26'

linters:
  enable:
    - gofmt
    - goimports
    - govet
    - staticcheck
    - errcheck
    - gosec
    - revive
    - unconvert
    - unparam
    - unused
    - ineffassign
    - misspell
    - dupl
    - gocritic

linters-settings:
  govet:
    enable-all: true
  staticcheck:
    checks: ["all"]
  errcheck:
    check-blank: true
  revive:
    rules:
      - name: exported
        severity: warning
        disabled: false

issues:
  exclude-use-default: false
  max-issues-per-linter: 0
  max-same-issues: 0
```

### .github/workflows/ci.yml 主要設定
```yaml
name: CI

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.26'
      - name: golangci-lint
        run: go run github.com/golangci/golangci-lint/cmd/golangci-lint run --timeout 5m

  test:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:17
        env:
          POSTGRES_PASSWORD: password
          POSTGRES_DB: kakeibo_test
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
        ports:
          - 5432:5432
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.26'
      - name: Run migrations
        run: go run github.com/amacneil/dbmate/v2 up
        env:
          DATABASE_URL: postgres://postgres:password@localhost:5432/kakeibo_test?sslmode=disable
      - name: Run tests
        run: go test -v -race -coverprofile=coverage.out ./...
        env:
          DATABASE_URL: postgres://postgres:password@localhost:5432/kakeibo_test?sslmode=disable
      - name: Upload coverage
        uses: codecov/codecov-action@v4
        with:
          files: ./coverage.out

  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.26'
      - name: Build
        run: go build -v ./...

  frontend:
    runs-on: ubuntu-latest
    defaults:
      run:
        working-directory: webapp
    steps:
      - uses: actions/checkout@v4
      - uses: jdx/mise-action@v2
      - name: Install dependencies
        run: vp install
      - name: Check
        run: vp check
      - name: Build
        run: vp build
```

### bobgen.yaml 主要設定
```yaml
struct_tag_casing: snake
tags: [json]
plugins:
  models:
    destination: "internal/infra/dbgen"
    pkgname: "dbgen"
  enums:
    destination: "internal/infra/dbgen/dbenums"
  factory:
    destination: "internal/infra/dbgen/factory"
  where: {}
  loaders: {}
  joins: {}
  dberrors: {}
  dbinfo:
    disabled: true
psql:
  dsn: "postgres://postgres:password@localhost:15432/kakeibo?sslmode=disable"
  # 公式仕様: bobgen は koanf + env.Provider で PSQL_ プレフィックスの環境変数をサポート
  # PSQL_DSN 環境変数で psql.dsn を上書き可能（公式ドキュメント: code-generation/psql.md）
  # mise.toml で PSQL_DSN を定義すれば OK
  uuid_pkg: google
  driver: "github.com/jackc/pgx/v5/stdlib"
  queries:
    - "./queries"
  except:
    "schema_migrations": {}
    "*":
      columns: [created_at, updated_at]
```

### queries/reports.sql（sqlc記法の実演）

家計簿向け集計クエリ3本:

```sql
-- name: GetCategoryExpenseSummary :many
-- カテゴリ別支出集計（指定期間）
SELECT
    c.id, c.name, c.category_type,
    COALESCE(SUM(t.amount), 0) AS total_amount,
    COUNT(t.id) AS transaction_count
FROM categories c
LEFT JOIN transactions t ON t.category_id = c.id 
    AND t.workspace_id = c.workspace_id
    AND t.transaction_type = 'expense'
    AND t.transaction_date >= $2
    AND t.transaction_date <= $3
WHERE c.workspace_id = $1 AND c.category_type = 'expense'
GROUP BY c.id, c.name, c.category_type
ORDER BY total_amount DESC;

-- name: GetAccountBalanceSummary :many
-- 口座別残高サマリ（初期残高 + 取引増減）
SELECT
    a.id, a.name, a.account_type, a.initial_balance,
    COALESCE(SUM(CASE 
        WHEN t.transaction_type = 'income' THEN t.amount
        WHEN t.transaction_type = 'expense' THEN -t.amount
        WHEN t.transaction_type = 'transfer' AND t.to_account_id = a.id THEN t.amount
        WHEN t.transaction_type = 'transfer' AND t.from_account_id = a.id THEN -t.amount
        ELSE 0
    END), 0) AS balance_change,
    (a.initial_balance + COALESCE(SUM(CASE 
        WHEN t.transaction_type = 'income' THEN t.amount
        WHEN t.transaction_type = 'expense' THEN -t.amount
        WHEN t.transaction_type = 'transfer' AND t.to_account_id = a.id THEN t.amount
        WHEN t.transaction_type = 'transfer' AND t.from_account_id = a.id THEN -t.amount
        ELSE 0
    END), 0)) AS current_balance
FROM accounts a
LEFT JOIN transactions t ON (t.account_id = a.id OR t.from_account_id = a.id OR t.to_account_id = a.id)
    AND t.workspace_id = a.workspace_id
WHERE a.workspace_id = $1
GROUP BY a.id, a.name, a.account_type, a.initial_balance
ORDER BY a.name;

-- name: GetWorkspaceDashboard :one
-- ダッシュボード概要（指定月の収支サマリ）
SELECT
    (SELECT COUNT(*) FROM accounts WHERE workspace_id = $1 AND is_active = true) AS active_accounts,
    (SELECT COUNT(*) FROM members WHERE workspace_id = $1) AS member_count,
    (SELECT COALESCE(SUM(amount), 0) FROM transactions WHERE workspace_id = $1 AND transaction_type = 'income' 
        AND transaction_date >= $2 AND transaction_date <= $3) AS total_income,
    (SELECT COALESCE(SUM(amount), 0) FROM transactions WHERE workspace_id = $1 AND transaction_type = 'expense' 
        AND transaction_date >= $2 AND transaction_date <= $3) AS total_expense,
    (SELECT COUNT(*) FROM transactions WHERE workspace_id = $1 
        AND transaction_date >= $2 AND transaction_date <= $3) AS transaction_count;
```

### WorkspaceSelectHook の仕組み（概念）

**注意**: 以下は設計意図を示す擬似コード。実装時に bob の実際の Hook 型シグネチャ（`bob.Hook[T] = func(context.Context, bob.Executor, T) (context.Context, error)` 等）を確認して正確に合わせること。また JOIN/Preload 時に列名が曖昧にならないよう、テーブル名を明示（`psql.Quote("members", "workspace_id")`）する。

```go
// 概念的な hook 設計（実装時に正確なシグネチャに合わせる）
// SELECT クエリに WHERE <table>.workspace_id = ? を自動注入
// executor から scope 情報を取得し、ModeGlobal なら skip
```

### トランザクションパターン（Service 層）

**概念的な実装パターン**（実装時に正確な bob の Transaction API に合わせる）:

```go
// internal/infra/db/tx.go
type TxFunc func(ctx context.Context, exec bob.Executor) error

func WithTx(ctx context.Context, db *bob.DB, fn TxFunc) error {
    tx, err := db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    
    defer func() {
        if p := recover(); p != nil {
            _ = tx.Rollback()
            panic(p)
        }
    }()
    
    if err := fn(ctx, tx); err != nil {
        _ = tx.Rollback()
        return err
    }
    
    return tx.Commit()
}

// internal/service/transaction.go
// 振替取引の例: from_account と to_account の原子性を保証
func (s *TransactionService) CreateTransfer(ctx context.Context, wsID uuid.UUID, req CreateTransferRequest) error {
    return db.WithTx(ctx, s.db, func(ctx context.Context, exec bob.Executor) error {
        // Repository に exec を渡してトランザクション内で実行
        tx := &dbgen.Transaction{
            ID:             uuid.NewV7(),
            WorkspaceID:    wsID,
            TransactionType: dbenums.TransactionTypeTransfer,
            Amount:         req.Amount,
            FromAccountID:  req.FromAccountID,
            ToAccountID:    req.ToAccountID,
            // ...
        }
        
        // 1. 振替取引を作成
        if err := s.repo.CreateTransaction(ctx, exec, tx); err != nil {
            return err
        }
        
        // 2. from_account の残高チェック（オプション: ビジネスロジック）
        balance, err := s.accountRepo.GetBalance(ctx, exec, wsID, req.FromAccountID)
        if err != nil {
            return err
        }
        if balance < req.Amount {
            return errors.New("insufficient balance")
        }
        
        // すべて成功すれば commit（WithTx が自動実行）
        return nil
    })
}
```

**Repository 層での対応**:
- Repository メソッドは `bob.Executor` インターフェースを受け取る
- `bob.Executor` は `*bob.DB` と `*bob.Tx` の両方を抽象化
- これにより同じメソッドがトランザクション内外で使える

```go
// internal/repository/transaction.go
func (r *TransactionRepository) CreateTransaction(ctx context.Context, exec bob.Executor, tx *dbgen.Transaction) error {
    // exec は *bob.DB または *bob.Tx（トランザクション境界を意識しない）
    _, err := tx.Insert(ctx, exec)
    return err
}
```

### transaction service での JOIN 実演（Preload + ThenLoad）（概念）

**注意**: 以下は概念的な記法。bob の生成コードは `Preload.<SingularTable>.<Rel>()` / `SelectThenLoad.<SingularTable>.<Rel>()` 形式で生成される。正確な API 名は bobgen 実行後の生成コードを参照すること。bob バージョンは固定して使う。

```go
// 概念: transaction を account, category, created_by(member) と共に取得
// Preload → LEFT JOIN（to-one）
// 振替の場合は from_account, to_account も Preload
```

### フロントエンド技術スタック

**フレームワーク**: Next.js (App Router) + TypeScript
**UIライブラリ**: Tailwind CSS + shadcn/ui (new-york スタイル)
**リント・フォーマット**: Vite+ の oxlint/oxfmt（バンドル済み、ESLint/Prettier 不要）
**API クライアント**: openapi-typescript + openapi-fetch（OpenAPI spec から型安全クライアント自動生成）
**パッケージマネージャ**: pnpm（workspace でモノレポ）
**ポート**: `${WEBAPP_PORT:-13001}`（next.config.ts で `--port` 指定）

**厳格な TypeScript 設定**（ncr-orchestrator 参考）:
```json
{
  "compilerOptions": {
    "strict": true,
    "noUnusedLocals": true,
    "noUnusedParameters": true,
    "noFallthroughCasesInSwitch": true,
    "isolatedModules": true
  }
}
```

**Vite+ によるリント・フォーマット**:
AGENTS.md の Vite+ セクション参照。webapp/ ディレクトリで以下を実行:
- `vp lint` — oxlint によるコードリント（TypeScript 型情報ベースのリントにも対応）
- `vp fmt` — oxfmt によるコードフォーマット（Prettier より高速）
- `vp check` — lint + fmt + TypeScript 型チェックの一括実行

**重要**: ESLint/Prettier は**インストール不要**。Vite+ が oxlint/oxfmt をバンドルして提供。

**OpenAPI クライアント生成**:
```bash
# webapp/ で実行
vp dlx openapi-typescript ../api/openapi.yaml -o src/lib/api/schema.d.ts
```
→ `openapi-fetch` の `createClient<paths>()` で型安全 API 呼び出し

**画面構成**:
| パス | 内容 | 使用 API |
|---|---|---|
| `/` | ワークスペース一覧・新規作成 | GET/POST /workspaces |
| `/workspaces/[wsId]` | ダッシュボード（月次収支、残高サマリ） | GET /workspaces/{wsId}/reports/dashboard |
| `/workspaces/[wsId]/transactions` | 取引一覧（日付範囲、カテゴリ、口座フィルタ）+ CRUD | GET/POST/PUT/DELETE /workspaces/{wsId}/transactions |
| `/workspaces/[wsId]/accounts` | 口座一覧（残高表示）+ CRUD | GET/POST/PUT/DELETE /workspaces/{wsId}/accounts |
| `/workspaces/[wsId]/categories` | カテゴリ一覧（階層表示）+ CRUD | GET/POST/PUT/DELETE /workspaces/{wsId}/categories |
| `/workspaces/[wsId]/members` | メンバー一覧 + CRUD | GET/POST/PUT/DELETE /workspaces/{wsId}/members |
| `/workspaces/[wsId]/reports` | 月次推移、カテゴリ別集計、口座残高推移 | GET /workspaces/{wsId}/reports/* |

### compose.yaml + mise.toml による環境分離

全ポートを環境変数で変更可能にし、既存サービスと干渉しない設計。

**mise.toml**（リポジトリにコミット）— ツールチェーン + デフォルト env:
```toml
[tools]
go = "1.26.1"
node = "24.14.0"
pnpm = "10.32.1"

[env]
POSTGRES_PORT = "15432"
POSTGRES_USER = "postgres"
POSTGRES_PASSWORD = "password"
POSTGRES_DB = "kakeibo"
GRAFANA_PORT = "13000"
LOKI_PORT = "13100"
APP_PORT = "18080"
WEBAPP_PORT = "13001"
NEXT_PUBLIC_API_URL = "http://localhost:{{env.APP_PORT}}"
DATABASE_URL = "postgres://{{env.POSTGRES_USER}}:{{env.POSTGRES_PASSWORD}}@localhost:{{env.POSTGRES_PORT}}/{{env.POSTGRES_DB}}?sslmode=disable"
PSQL_DSN = "{{env.DATABASE_URL}}"  # bobgen が koanf 経由で psql.dsn を上書き
```

**mise.local.toml**（.gitignore済み、個人PC向け上書き）— 既存の GH_TOKEN 等に加え、ポート競合時の上書き等:
```toml
[env]
# 既にこのPCでは GH_TOKEN, MISE_GITHUB_TOKEN, DEVIN_API_KEY が定義済み
# ポート競合時はここで上書き:
# POSTGRES_PORT = "25432"
```

**compose.yaml** で `${POSTGRES_PORT:-15432}` 等の env 展開を使用:
- `postgres:17` — port `${POSTGRES_PORT:-15432}`:5432
- `grafana/loki:3.4.3` — port `${LOKI_PORT:-13100}`:3100
- `grafana/grafana:12.0.0` — port `${GRAFANA_PORT:-13000}`:3000, anonymous auth
- `grafana/promtail:3.4.3` — Docker logs → Loki

デフォルトで本業の 5432/3000/3100 と干渉しないポート（15432/13000/13100）を使用。
bobgen.yaml の DSN も `DATABASE_URL` 環境変数から取る or mise.toml 経由で注入。

## 検証手順

1. `docker compose up -d` → PostgreSQL(15432) + Grafana(13000) + Loki(13100) 起動（ポート干渉なし）
2. `make migrate` → テーブル作成
3. `make bobgen` → モデル・クエリ生成、コンパイル確認
4. `make ogen` → API サーバコード生成
5. `make lint` → golangci-lint 実行（エラーなし確認）
6. `go test ./queries/...` → 生成クエリのテスト実行（bob 公式推奨）
7. `go test -v -race ./...` → 全テスト実行（race detector 有効）
8. `go build ./...` → 全体コンパイル
9. `go run ./cmd/seed/` → サンプルデータ投入（口座・カテゴリ・取引）
10. `go run ./cmd/server/` → サーバ起動 (`:${APP_PORT:-18080}`)
11. curl でエンドポイント確認（ポートは `${APP_PORT:-18080}`）:
   - `POST /workspaces` → ワークスペース作成
   - `POST /workspaces/{wsId}/members` → メンバー追加
   - `POST /workspaces/{wsId}/accounts` → 口座作成（initial_balance設定）
   - `POST /workspaces/{wsId}/categories` → カテゴリ作成
   - `POST /workspaces/{wsId}/transactions` → 取引作成（収入/支出/振替）
   - `GET /workspaces/{wsId}/transactions?start_date=2026-03-01&end_date=2026-03-31&type=expense` → フィルタ付き一覧
   - `GET /workspaces/{wsId}/reports/category-summary?start_date=2026-03-01&end_date=2026-03-31` → カテゴリ別集計（queries plugin）
   - `GET /workspaces/{wsId}/reports/account-balance` → 口座残高サマリ
   - `GET /workspaces/{wsId}/reports/dashboard?year=2026&month=3` → ダッシュボード
12. テナント分離確認: WS-A のデータが WS-B から見えないこと
13. 口座残高計算確認: initial_balance + 取引増減 = current_balance
14. 振替取引確認: from_account 減、to_account 増が正しく反映（トランザクション原子性）
15. トランザクションロールバック確認: エラー時に振替が巻き戻されること
16. Grafana (`localhost:${GRAFANA_PORT:-13000}`) でアプリログ確認
17. CI確認: `.github/workflows/ci.yml` で lint, test, build が通ること
18. `cd webapp && vp install && vp dev` → フロントエンド起動 (`:${WEBAPP_PORT:-13001}`)
19. `vp dlx openapi-typescript ../api/openapi.yaml -o src/lib/api/schema.d.ts` → OpenAPI クライアント型生成
20. ブラウザで `localhost:${WEBAPP_PORT:-13001}` → ワークスペース一覧 → 取引登録 → ダッシュボード確認
21. `vp check` → oxlint + oxfmt + TypeScript 型チェック（エラーなし確認）
