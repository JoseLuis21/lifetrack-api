terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
  }
}

provider "aws" {
  region = "us-east-1"
}

resource "aws_lambda_function" "add-category" {
  function_name = "lt-add-category"
  s3_bucket = "life-track-serverless"
  s3_key = "v1.0.0/add-category"
  handler = "main"
  role = aws_iam_role.lambda-exec-write-db.arn
  runtime = "go1.x"
}

data "aws_iam_policy_document" "lambda-dynamo-read" {
  statement {
    actions = ["sts:AssumeRole", "dynamodb:GetItem",
      "dynamodb:Scan", "dynamodb:Query"]
    resources = [aws_dynamodb_table.lt-category.arn]
    effect = "Allow"
    principals {
      identifiers = ["lambda.amazonaws.com"]
      type = ""
    }
  }
}

data "aws_iam_policy_document" "lambda-dynamo-write" {
  statement {
    actions = ["sts:AssumeRole", "dynamodb:PutItem",
      "dynamodb:DeleteItem", "dynamodb:UpdateItem"]
    resources = [aws_dynamodb_table.lt-category.arn]
    effect = "Allow"
    principals {
      identifiers = ["lambda.amazonaws.com"]
      type = ""
    }
  }
}

resource "aws_iam_role" "lambda-exec-write-db" {
  name = "lt-lambda-write-db"
  assume_role_policy = data.aws_iam_policy_document.lambda-dynamo-write.json
}
