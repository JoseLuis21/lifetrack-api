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
#   $1 = Module name (e.g. category)
#   $2 = Function name (e.g. add-category)
#   $3 = Application version (e.g. 1.0.1)
#     Uses v1.0.0 as default value
#   $4 = AWS S3 bucket name (e.g. my-lifetrack-bucket)
#     Uses lifetrack-serverless as default value
#######################################
function upload_to_s3() {
  APP_VERSION="v1.0.0"
  # Legacy bucket name - life-track-serverless
  BUCKET_NAME="lifetrack-serverless"

  if [ "$3" != "" ]; then
      APP_VERSION=v"$3"
  fi

  if [ "$4" != "" ]; then
      BUCKET_NAME="$4"
  fi

  if [ "$1" != "" ] || [ "$2" != "" ]; then
    FILE_NAME="$1/$2"
    LT_PATH="$GOPATH/src/github.com/neutrinocorp/life-track-api"
    /bin/bash "$LT_PATH/build/build.sh" "$1" "$2"
    aws s3 cp "$LT_PATH/build/release/$FILE_NAME".zip s3://"$BUCKET_NAME/$APP_VERSION/$FILE_NAME".zip
    else
      err "empty argument(s)"
      exit 1
  fi
}

upload_to_s3 "$1" "$2" "$3" "$4"
