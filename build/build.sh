#!/bin/bash

if [ "$1" != "" ]; then
  LT_PATH="$GOPATH/src/github.com/neutrinocorp/life-track-api"
  GOOS=linux GOARCH=amd64 go build -o "$LT_PATH"/build/bin/"$1" "$LT_PATH"/cmd/"$1"/main.go
  zip -D -j "$LT_PATH"/build/release/"$1".zip "$LT_PATH"/build/bin/"$1"
  else
    echo "empty argument"
fi
