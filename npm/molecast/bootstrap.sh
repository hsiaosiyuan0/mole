#!/usr/bin/env bash

arch=$(node -e "console.log(require('os').arch())")
os=$(node -e "console.log(require('os').platform())")

if [ "$arch" = "x64" ]; then
  arch="amd64"
fi

mole="molecast-$os-$arch"
$mole "$@"
