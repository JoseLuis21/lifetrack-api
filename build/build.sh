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
#   $1 = Module name (e.g. add-category)
#######################################
function build() {
  if [ "$1" != "" ]; then
    echo "building go binary file for module: $1"
    LT_PATH="$GOPATH/src/github.com/neutrinocorp/life-track-api"
    GOOS=linux GOARCH=amd64 go build -o "$LT_PATH"/build/bin/"$1" "$LT_PATH"/cmd/"$1"/main.go
    echo "binary was successfully built"
  else
    err "empty argument"
    exit 1
  fi
}

#######################################
# Zip (compress) Go binary file
# Globals:
#   GOPATH
# Arguments:
#   $1 = Module name (e.g. add-category)
#######################################
function compress() {
  if [ "$1" != "" ]; then
    echo "compressing binary file for module: $1"
    LT_PATH="$GOPATH/src/github.com/neutrinocorp/life-track-api"
    zip -D -j "$LT_PATH"/build/release/"$1".zip "$LT_PATH"/build/bin/"$1"
    echo "compression completed, output file can be found at $LT_PATH/build/release/$1.zip"
  else
    err "empty argument"
    exit 1
  fi
}

build "$1"
compress "$1"
