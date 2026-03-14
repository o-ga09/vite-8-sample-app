---
name: implement-issue-and-create-pr
description: >
  GitHubのIssueを取得し、t-wadaのTDD（Red→Green→Refactor）サイクルで実装してPRを作成するスキル。
  「Issueを実装して」「#123のIssueを対応して」「TDDで実装して」「PRを作って」「Issueを取ってきて実装して」
  と言われたときに必ずこのスキルを使用すること。Go・フロントエンド（TypeScript/React/Vite）の両方に対応している。
---

# Issue実装 & PR作成スキル

GitHubのIssueを取得し、t-wadaのTDDで実装し、PRを作成するエンドツーエンドのワークフロー。

## プロジェクト概要

- **Goバックエンド**: `github.com/o-ga09/vite-8-sample-app`（`cmd/`, `internal/` 配下）
- **フロントエンド**: `webapp/` 配下（React + TypeScript + Vite+）
- **テストツール**: Go → `go test`, フロントエンド → `vp test`（Vite+経由のVitest）

---

## Step 1: Issueの取得と分析

GitHub MCPが利用可能な場合はMCPを優先し、利用できない場合は `gh` コマンドを使用する。

**GitHub MCPを使う場合:**
```
mcp_github_get_issue(owner="o-ga09", repo="vite-8-sample-app", issue_number=<番号>)
```

**`gh` コマンドを使う場合:**
```bash
gh issue view <issue_number> --json title,body,labels,assignees
```

取得後に以下を整理する:
- 実装すべき機能・修正の内容
- 変更対象が **Go（バックエンド）** か **フロントエンド** か（両方の場合もある）
- 受け入れ条件（完了の定義）

---

## Step 2: ブランチの作成

```bash
git-wt feature/<issue_number>-<short-description>
```

例:
- `feature/42-add-user-authentication`
- `fix/17-correct-calculation-error`

---

## Step 3: TDDで実装する（t-wada's TDD 3原則）

> **TDDの3ルール（Robert C. Martin）をt-wadaスタイルで徹底する:**
> 1. 失敗するテストを1つだけ書くまで、プロダクションコードを書いてはならない
> 2. コンパイルエラーもテストの失敗とみなす
> 3. 現在失敗しているテスト1つを通過させる以上のプロダクションコードを書いてはならない

### 🔴 Red — 失敗するテストを書く

- まず**テストファイルを先に作成**する
- テストコードだけ書き、プロダクションコードはまだ書かない
- テストが「コンパイルエラー」または「アサーション失敗」で落ちることを確認する
- テスト名は「何をすべきか」が明確にわかる名前にする

### 🟢 Green — テストを通過させる最小限のコードを書く

- テストが通るための**最小限**のプロダクションコードを書く
- この段階では美しさより「動くこと」を優先する
- テスト実行して全件パスを確認する

### 🔵 Refactor — テストを通したままコードを改善する

- テストがグリーンの状態を維持しながらリファクタリングする
- 重複排除・命名改善・責務分離などを行う
- リファクタリング後も全テストがパスすることを確認する

このサイクルを**機能の最小単位ごと**に繰り返す。

---

## Step 4: テストの実行

### Goの場合

```bash
# 実装したパッケージのみテスト
go test ./internal/<package>/...

# 全体テスト
go test ./...

# 静的解析も実行
go vet ./...
```

### フロントエンドの場合

```bash
# webapp ディレクトリで実行
cd webapp

# 特定のテストファイルを指定
vp test <test-file-path>

# 全テスト実行
vp test

# lint / type check も実施
vp check
```

**テストが全件パスしてから次のステップに進む。**

---

## Step 5: コミットとプッシュ

```bash
git add .
git commit -m "<type>(<scope>): <subject>

Closes #<issue_number>"
git push origin <branch-name>
```

コミットメッセージの type: `feat` / `fix` / `refactor` / `test` / `docs` / `chore`

例:
- `feat(api): add user authentication endpoint`
- `fix(webapp): correct total calculation logic`

---

## Step 6: PRの作成

### PRテンプレートの確認

`.github/PULL_REQUEST_TEMPLATE.md` が存在する場合はテンプレートを使用する。

### GitHub MCPを使う場合

```
mcp_github_create_pull_request(
  owner="o-ga09",
  repo="vite-8-sample-app",
  title="<PRタイトル>",
  body="<テンプレートに沿ったPR本文>",
  head="<branch-name>",
  base="main",
  reviewers=["<自分のGitHubユーザー名>"]
)
```

### `gh` コマンドを使う場合

```bash
gh pr create \
  --title "<PRタイトル>" \
  --body-file .github/PULL_REQUEST_TEMPLATE.md \
  --reviewer "@me" \
  --base main
```

PRの body には以下を必ず記入する:
- `Closes #<issue_number>` — Issueとの紐付け
- 変更内容の箇条書き
- 実行したテストの結果
- フロントエンドの変更がある場合はスクリーンショット

**レビュワーには必ず `@me`（自分自身）を設定すること。**

---

## フロントエンド固有の追加手順

フロントエンド（`webapp/`）に変更がある場合:

1. **`ui-screenshot` スキルを使用**してスクリーンショットを撮影する（`.agents/skills/ui-screenshot/SKILL.md` を参照）
2. 取得したスクリーンショットをPRのbodyに添付する（After形式、可能であればBefore/After形式）

---

## 対応言語・パターンの判定

| 変更ファイルパス | 言語 | テストコマンド |
|---|---|---|
| `cmd/`, `internal/` | Go | `go test ./...` |
| `webapp/src/` | フロントエンド | `cd webapp && vp test` |
| 両方含む | 両方 | 両方を順に実行 |
