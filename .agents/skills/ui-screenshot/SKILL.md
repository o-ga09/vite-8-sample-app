---
name: ui-screenshot
description: >
  UIの変更を実装した後、開発サーバーを起動してPlaywrightのヘッドレスブラウザでスクリーンショットを撮影するスキル。
  「スクリーンショットを撮って」「UIを確認して」「見た目を確認して」「ブラウザで確認して」「見せて」
  と言われたとき、またはwebapp/src/配下のReact/TypeScript/CSSファイルを変更した直後に必ずこのスキルを使用すること。
  implement-issue-and-create-prスキルでフロントエンドの変更を実装した後にも自動的に呼び出すこと。
---

# UIスクリーンショット撮影スキル

UIの変更を実装した後、変更されたページのスクリーンショットをPlaywright（ヘッドレスChromium）で自動撮影し、`webapp/screenshots/` に保存する。

---

## Step 1: 変更されたUIファイルを特定し、撮影URLを決定する

```bash
# 変更されたwebapp/src/配下のファイルを取得
git diff --name-only HEAD -- webapp/src/ 2>/dev/null
git diff --name-only --cached -- webapp/src/ 2>/dev/null
```

変更ファイルから撮影対象URLを決定する:

| 変更ファイルのパターン | 撮影URL |
|---|---|
| `webapp/src/App.tsx`, `webapp/src/style.css` など | `http://localhost:5173/` |
| `webapp/src/pages/Foo.tsx` | `http://localhost:5173/foo` |
| `webapp/src/components/Bar.tsx` | Barコンポーネントを使うページのURL |

**ルーター定義がある場合:**
`webapp/src/router.tsx`（または`App.tsx`内のroutes定義）を読み込み、変更コンポーネントにマッチするルートパスを特定する。

**ルーターがない場合（現在のシングルページアプリ）:**
`http://localhost:5173/` だけを撮影する。

**明示的に指定された場合:** ユーザーが「`/about`を撮って」のように言った場合はそのURLを使う。

---

## Step 2: gitignoreの確認（初回のみ）

```bash
grep -q "screenshots/" .gitignore || echo "" >> .gitignore && echo "webapp/screenshots/" >> .gitignore
```

---

## Step 3: Playwrightのセットアップ（未インストールの場合のみ）

```bash
cd webapp

# Playwrightがインストール済みか確認
if ! ls node_modules/.bin/playwright &>/dev/null; then
  vp add -D playwright
  vp exec playwright install chromium
fi
```

> `vp add -D playwright` は `webapp/` ディレクトリで実行すること（ルートではない）。
> ブラウザのインストース(`playwright install chromium`)は初回のみ必要。

---

## Step 4: 開発サーバーの起動

```bash
cd webapp

# サーバー起動（バックグラウンド）
vp dev &
DEV_SERVER_PID=$!

# 起動を待つ（最大30秒）
for i in $(seq 1 30); do
  curl -sf http://localhost:5173 > /dev/null 2>&1 && break
  sleep 1
done
```

> `vp dev` は `webapp/` ディレクトリで実行すること。

---

## Step 5: スクリーンショット撮影

撮影にはこのスキルに同梱された `scripts/screenshot.mjs` を使う:

```bash
TIMESTAMP=$(date +%Y%m%d-%H%M%S)
OUTPUT_DIR="webapp/screenshots/${TIMESTAMP}"

node .agents/skills/ui-screenshot/scripts/screenshot.mjs \
  --urls "http://localhost:5173/,http://localhost:5173/about" \
  --output "$OUTPUT_DIR"
```

`--urls` には Step 1 で決定したURLをカンマ区切りで渡す。出力先は `webapp/screenshots/<タイムスタンプ>/` とする。

---

## Step 6: サーバーの停止と結果の報告

```bash
# サーバーを停止
kill $DEV_SERVER_PID 2>/dev/null
wait $DEV_SERVER_PID 2>/dev/null
```

撮影完了後、保存されたスクリーンショットのパス一覧をユーザーに報告する:

```
✅ スクリーンショット撮影完了
   webapp/screenshots/20260314-120000/index.png
   webapp/screenshots/20260314-120000/about.png
```

---

## implement-issue-and-create-prとの連携

`implement-issue-and-create-pr` スキルでフロントエンド（`webapp/src/`）の変更を実装した後:

1. このスキルを呼び出してスクリーンショットを撮影する
2. 撮影したスクリーンショットのパスをPRのbody（Before/After形式）に添付する

PRのbodyへの記載例:
```markdown
## スクリーンショット
| Before | After |
|--------|-------|
| （変更前の画像がある場合） | ![after](webapp/screenshots/20260314-120000/index.png) |
```

---

## トラブルシューティング

| 症状 | 対処 |
|---|---|
| `curl` が30秒たっても成功しない | `vp dev` の出力を確認。ポートが5173以外なら `vite.config.ts` の `server.port` を確認 |
| `playwright` コマンドが見つからない | `cd webapp && vp add -D playwright && vp exec playwright install chromium` を再実行 |
| スクリーンショットが真っ白 | `waitUntil: 'networkidle'` で待機しているが、APIリクエストが完了しない場合は `--wait-ms 3000` オプションを追加 |
| ポートがすでに使用中 | `lsof -ti:5173 | xargs kill -9` でポートを解放してから再試行 |
