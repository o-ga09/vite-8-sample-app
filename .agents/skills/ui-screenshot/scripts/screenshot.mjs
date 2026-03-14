#!/usr/bin/env node
/**
 * screenshot.mjs
 * Playwright を使ってヘッドレスChromiumでスクリーンショットを撮影する。
 *
 * 使い方:
 *   node screenshot.mjs --urls "http://localhost:5173/,http://localhost:5173/about" --output "webapp/screenshots/20260314"
 *
 * オプション:
 *   --urls      カンマ区切りのURL一覧（デフォルト: http://localhost:5173）
 *   --output    出力ディレクトリ（デフォルト: webapp/screenshots/<ISO timestamp>）
 *   --wait-ms   ページ読み込み後の追加待機時間ms（デフォルト: 0）
 *   --width     ビューポート幅（デフォルト: 1280）
 *   --height    ビューポート高さ（デフォルト: 720）
 */

import { chromium } from "playwright";
import { mkdirSync } from "fs";
import { resolve, join } from "path";

// --- 引数パース ---
const args = process.argv.slice(2);
const getArg = (name, defaultValue) => {
  const idx = args.indexOf(`--${name}`);
  return idx >= 0 ? args[idx + 1] : defaultValue;
};

const urlsArg = getArg("urls", "http://localhost:5173");
const outputDir = resolve(
  getArg(
    "output",
    `webapp/screenshots/${new Date().toISOString().replace(/[:.]/g, "-")}`,
  ),
);
const waitMs = parseInt(getArg("wait-ms", "0"), 10);
const viewportWidth = parseInt(getArg("width", "1280"), 10);
const viewportHeight = parseInt(getArg("height", "720"), 10);

const urls = urlsArg
  .split(",")
  .map((u) => u.trim())
  .filter(Boolean);

// --- 出力ディレクトリを作成 ---
mkdirSync(outputDir, { recursive: true });

console.log(`📸 スクリーンショット撮影開始`);
console.log(`   対象URL: ${urls.join(", ")}`);
console.log(`   出力先: ${outputDir}`);

const browser = await chromium.launch({ headless: true });
const context = await browser.newContext({
  viewport: { width: viewportWidth, height: viewportHeight },
});

for (const url of urls) {
  const page = await context.newPage();
  try {
    await page.goto(url, { waitUntil: "networkidle", timeout: 30000 });

    if (waitMs > 0) {
      await page.waitForTimeout(waitMs);
    }

    // URL パスからファイル名を生成（例: / → index, /about → about）
    const pathname = new URL(url).pathname;
    const slug =
      pathname === "/"
        ? "index"
        : pathname.replace(/^\//, "").replace(/\//g, "_");
    const filename = join(outputDir, `${slug}.png`);

    await page.screenshot({ path: filename, fullPage: true });
    console.log(`   ✅ ${url} → ${filename}`);
  } catch (err) {
    console.error(`   ❌ ${url} の撮影に失敗: ${err.message}`);
  } finally {
    await page.close();
  }
}

await browser.close();
console.log(`\n✅ スクリーンショット撮影完了: ${outputDir}`);
