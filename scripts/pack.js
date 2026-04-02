/**
 * pack.js — bundles both packs into a single .mcaddon file in dist/.
 *
 * Usage:
 *   npm run pack
 *
 * Output:
 *   dist/hello-world.mcaddon
 */

"use strict";

const archiver = require("archiver");
const fs = require("fs");
const path = require("path");

const ROOT = path.resolve(__dirname, "..");
const DIST = path.join(ROOT, "dist");

if (!fs.existsSync(DIST)) {
  fs.mkdirSync(DIST, { recursive: true });
}

const outPath = path.join(DIST, "catmob-madness.mcaddon");
const output = fs.createWriteStream(outPath);
const archive = archiver("zip", { zlib: { level: 9 } });

output.on("close", () => {
  console.log(`✅  Packed ${archive.pointer()} bytes → ${outPath}`);
});

archive.on("error", (err) => {
  throw err;
});

archive.pipe(output);

// Each top-level folder inside the .mcaddon becomes its own pack.
archive.directory(path.join(ROOT, "behavior_pack"), "behavior_pack");
archive.directory(path.join(ROOT, "resource_pack"), "resource_pack");

archive.finalize();
