#!/bin/bash
#
# Perform Go binary build for Amazon Linux/Linux OS distros and zip it

err() {
  echo "[$(date +'%Y-%m-%dT%H:%M:%S%z')]: $*" >&2
}

#######################################
# Build Go binary file
# Globals:
#   GOPATH
# Arguments:
#   $1 = Module name (e.g. category)
#   $2 = Function name (e.g. add-category)
#######################################
function build() {
  if [ "$1" != "" ] || [ "$2" != "" ]; then
    FILE_NAME="$1/$2"
    echo "building go binary file for module: $FILE_NAME"
    LT_PATH="$GOPATH/src/github.com/neutrinocorp/lifetrack-api"
    GOOS=linux GOARCH=amd64 go build -o "$LT_PATH/build/bin/server/$FILE_NAME" "$LT_PATH/cmd/serverless/$FILE_NAME/main.go"
    echo "binary was successfully built"
  else
    err "empty argument(s)"
    exit 1
  fi
}

#######################################
# Zip (compress) Go binary file
# Globals:
#   GOPATH
# Arguments:
#   $1 = Module name (e.g. category)
#   $2 = Function name (e.g. add-category)
#######################################
function compress() {
  if [ "$1" != "" ] || [ "$2" != "" ]; then
    FILE_NAME="$1/$2"
    echo "compressing binary file for module: $FILE_NAME"
    LT_PATH="$GOPATH/src/github.com/neutrinocorp/lifetrack-api"
    mkdir "$LT_PATH/build/release/$1"
    zip -D -j "$LT_PATH/build/release/$FILE_NAME".zip "$LT_PATH/build/bin/$FILE_NAME"
    echo "compression completed, output file can be found at $LT_PATH/build/release/$FILE_NAME.zip"
  else
    err "empty argument(s)"
    exit 1
  fi
}

build "$1" "$2"
compress "$1" "$2"
