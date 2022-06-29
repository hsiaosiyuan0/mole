#!/usr/bin/env bash

arch=$(node -e "console.log(require('os').arch())")
os=$(node -e "console.log(require('os').platform())")

if [ "$arch" = "x64" ]; then
  arch="amd64"
fi

cwd=$(dirname $(realpath $0))
mole="$cwd/mole-$os-$arch"

if [ ! -f "$mole" ]; then
  echo "Unable to run molecast at $mole"
else
  $mole "$@"
fi
