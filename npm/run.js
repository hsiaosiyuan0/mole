#!/usr/bin/env node

const os = require("os");
const path = require("path");
const fs = require("fs");
const { execSync } = require("child_process");

function fetchVer() {
  const info = {
    latest: "",
    yarn: false,
    npm: false,
  };

  try {
    info.latest = execSync("npm view molecast --json");
    info.npm = true;
  } catch (error) {}

  if (!info) {
    try {
      info.latest = execSync("yarn info molecast --json");
      info.yarn = true;
    } catch (error) {}
  }

  if (info) {
    try {
      info.latest = JSON.parse(info.latest)["dist-tags"]["latest"];
    } catch (error) {}
  }

  return info;
}

const selfVersion = require(path.join(__dirname, "./package.json")).version;
const info = fetchVer();

let args = process.argv.slice(2);
const isQuirks = args.includes("--quirks");

if (info.latest && info.latest !== selfVersion && !isQuirks) {
  console.log(
    `Please upgrade to the latest version: ${info.latest}, use \`--quirks\` force to run`
  );
  process.exit(0);
}

const OS = os.platform();
let ARCH = os.arch();
if (ARCH === "x64") {
  ARCH = "amd64";
}

const mole = path.join(__dirname, `mole-${OS}-${ARCH}`);

if (fs.existsSync(mole)) {
  if (isQuirks) {
    args = args.filter((arg) => arg !== "--quirks");
  }

  try {
    execSync(mole + " " + args.join(" "), { shell: true, stdio: "inherit" });
  } catch (error) {
    // discard the `Command failed` error raised from nodejs to keep the error output clean
  }
} else {
  console.log(`Unable to run molecast at ${mole}`);
}
