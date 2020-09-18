#!/bin/bash

APP_VERSION="v1.0.0"

if [ "$2" != "" ]; then
    APP_VERSION=v"$2"
fi

if [ "$1" != "" ]; then
  LT_PATH="$GOPATH/src/github.com/neutrinocorp/life-track-api"
  /bin/bash "$LT_PATH"/build/build.sh "$1"
  aws s3 cp "$LT_PATH"/build/release/"$1".zip s3://life-track-serverless/"$APP_VERSION"/"$1".zip
  else
    echo "empty argument"
fi