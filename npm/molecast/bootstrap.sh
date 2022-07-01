#!/usr/bin/env bash

VERSION=0.0.30

OS=$(uname -s | awk '{print tolower($0)}')
ARCH=$(uname -m | awk '{print tolower($0)}')

if [ "$ARCH" = "x86_64" ]; then
  ARCH="amd64"
fi

molecast="molecast-$OS-$ARCH"
molecast_ver="$molecast-$VERSION"
tgz="$molecast_ver.tgz"
tgz_link="https://registry.npmjs.org/$molecast/-/$tgz"

dir=$(dirname $0)
molecast_bin="$dir/$molecast_ver"

info() {
  echo -e "\033[1;36m$1\033[0m"
}

ok() {
  echo -e "\033[1;32m$1\033[0m"
}

download() {
  info "\nThe executable binary does not exist, downloading it from $tgz_link\n"

  if which curl >/dev/null; then
    curl $tgz_link -o /tmp/$tgz
  elif which wget >/dev/null; then
    wget -O /tmp/$tgz $tgz_link
  fi

  d=/tmp/$molecast
  mkdir -p $d
  tar -xzf /tmp/$tgz -C $d && mv $d/package/mole $dir/$molecast_ver

  if [ $? -eq 0 ]; then
    ok "\nSucceeded, please molecast again.\n"
  fi
}

if [ -x $molecast_bin ]; then
  $molecast_bin "$@"
else
  download
fi
