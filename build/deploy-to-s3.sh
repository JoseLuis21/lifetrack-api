#!/bin/bash

if [ "$1" != "" ]; then
  /bin/bash ./build.sh "$1"
  aws s3 cp ./release/"$1".zip s3://life-track-serverless/v1.0.0/"$1".zip
  else
    echo "empty argument"
fi