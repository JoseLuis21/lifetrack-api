#!/bin/bash

if [ "$1" != "" ]; then
  LT_PATH="$GOPATH/src/github.com/neutrinocorp/life-track-api"
  GOOS=linux go build -o "$LT_PATH"/build/bin/"$1" "$LT_PATH"/cmd/"$1"/main.go
  zip "$LT_PATH"/build/release/"$1".zip "$LT_PATH"/build/bin/"$1"
  else
    echo "empty argument"
fi
