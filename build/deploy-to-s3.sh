#!/bin/bash

if [ "$1" != "" ]; then
  LT_PATH="$GOPATH/src/github.com/neutrinocorp/life-track-api"
  /bin/bash "$LT_PATH"/build/build.sh "$1"
  aws s3 cp "$LT_PATH"/build/release/"$1".zip s3://life-track-serverless/v1.0.0/"$1".zip
  else
    echo "empty argument"
fi