# Build
To build an AWS Lambda function, just run the `build.sh` file with a single argument as the name of the Lambda function.

The shell script will use your `GOPATH` _environment variable_ and will lookup for Neutrino's _Life 
Track_ API project. For example, a correct path would be `~/go/src/github.com/neutrinocorp/life-track-api`.

_E.g._ `/bin/bash ./build.sh add-category`

## Requirements

- You **must** have the following  `GOPATH` _environment variable defined_.

- You **must** have _Go 1.15_ version installed.

# Deploy to AWS S3
If you need to deploy to AWS S3 Life Track with a single command, you may run the `deploy-to-s3.sh` file with 
two arguments as the name of the Lambda function and desired version (_using [Semantic Versioning](https://semver.org)_).

Likewise, the `deploy-to-s3.sh` shell script will use your `GOPATH` _environment variable_ to lookup for 
required path(s).

_E.g._ `/bin/bash ./deploy-to-s3.sh add-category`

## Create target bucket from CLI
In order to deploy a Lambda function, you must create an S3 bucket first to upload any compiled binary and subsequently, 
deploy it.

Run the following command into your favorite terminal.

`aws s3api create-bucket --bucket=life-track-serverless --region=us-east-1`

_You may change either the bucket's name or region, but have in mind you would need to modify the 
`deploy-to-s3.sh` script to match your new bucket name. Also, you would need to modify your AWS CLI 
default region (at `~/.aws/config`)._

## Requirements

- You **must** have the following  `GOPATH` _environment variable defined_.

- You **must** have _Go 1.15_ version installed.

- You **must** have _AWS CLI v2_ installed.

- You **must** have _AWS required policies_ to write/read to an S3 Bucket defined at `~/.aws/credentials` file.
