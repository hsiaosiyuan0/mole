#!/usr/bin/env bash

cat <<EOF > .env
GOOS=$(node -e "console.log(require('os').platform())")
GOARCH=amd64
EOF
