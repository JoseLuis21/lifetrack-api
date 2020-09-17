#!/bin/bash

if [ "$1" != "" ]; then
  GOOS=linux go build -o ./bin/"$1" ../cmd/"$1"/main.go
  zip ./release/"$1".zip ./bin/"$1"
  else
    echo "empty argument"
fi
