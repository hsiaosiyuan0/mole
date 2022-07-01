#!/usr/bin/env bash

VERSION=0.0.29

OS=$(uname -s | awk '{print tolower($0)}')
ARCH=$(uname -m | awk '{print tolower($0)}')

if [ "$ARCH" = "x86_64" ]; then
  ARCH="amd64"
fi

molecast="molecast-$OS-$ARCH"
tgz="$molecast-$VERSION.tgz"
tgz_link="https://registry.npmjs.org/$molecast/-/$tgz"

dir=$(dirname $0)
mole="$dir/mole"

download() {
  echo "downloading $tgz_link"

  if which curl >/dev/null; then
    curl $tgz_link -o /tmp/$tgz
  elif which wget >/dev/null; then
    wget -O /tmp/$tgz $tgz_link
  fi

  d=/tmp/$molecast-$VERSION
  mkdir -p $d
  tar -xzf /tmp/$tgz -C $d
  mv $d/package/mole $dir/
}

if [ -x $mole ]; then
  $mole "$@"
else
  download
fi
