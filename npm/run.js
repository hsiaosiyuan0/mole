#!/usr/bin/env node

const os = require("os");
const path = require("path");
const fs = require("fs");
const { execSync } = require("child_process");

const OS = os.platform();
let ARCH = os.arch();

if (ARCH === "x64") {
  ARCH = "amd64";
}

const mole = path.join(__dirname, `mole-${OS}-${ARCH}`);

if (fs.existsSync(mole)) {
  execSync(mole + " " + process.argv.slice(2).join(" "));
} else {
  console.log(`Unable to run molecast at ${mole}`);
}
