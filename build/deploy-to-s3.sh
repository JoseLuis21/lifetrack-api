#!/bin/bash
#
# Perform Go binary build for Amazon Linux/Linux OS based distros and subsequently upload it to AWS S3
# Requires AWS CLi and proper configured credentials

err() {
  echo "[$(date +'%Y-%m-%dT%H:%M:%S%z')]: $*" >&2
}

#######################################
# Run build.sh executable shell script and upload output compressed file to AWS S3
# Important Note: Requires AWS CLi and proper credentials
# Globals:
#   GOPATH
# Arguments:
#   $1 = Module name (e.g. add-category)
#   $2 = Application version (e.g. 1.0.1)
#     Uses v1.0.0 as default value
#   $3 = AWS S3 bucket name (e.g. my-lifetrack-bucket)
#     Uses lifetrack-serverless as default value
#######################################
function upload_to_s3() {
  APP_VERSION="v1.0.0"
  # Legacy bucket name - life-track-serverless
  BUCKET_NAME="lifetrack-serverless"

  if [ "$2" != "" ]; then
      APP_VERSION=v"$2"
  fi

  if [ "$3" != "" ]; then
      BUCKET_NAME="$3"
  fi

  if [ "$1" != "" ]; then
    LT_PATH="$GOPATH/src/github.com/neutrinocorp/life-track-api"
    /bin/bash "$LT_PATH"/build/build.sh "$1"
    aws s3 cp "$LT_PATH"/build/release/"$1".zip s3://"$BUCKET_NAME"/"$APP_VERSION"/"$1".zip
    else
      err "empty argument"
      exit 1
  fi
}

upload_to_s3 "$1" "$2" "$3"
