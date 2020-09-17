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
  role = aws_iam_role.lambda-exec-full-db.arn
  runtime = "go1.x"
  environment {
    variables = {
      "LT_TABLE_NAME": aws_dynamodb_table.lt-category.name,
      "LT_TABLE_REGION": "us-east-1"
    }
  }
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

data "aws_iam_policy_document" "lambda-dynamo-full" {
  statement {
    actions = ["sts:AssumeRole", "dynamodb:PutItem",
      "dynamodb:DeleteItem", "dynamodb:UpdateItem", "dynamodb:GetItem",
      "dynamodb:Scan", "dynamodb:Query"]
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

resource "aws_iam_role" "lambda-exec-read-db" {
  name = "lt-lambda-read-db"
  assume_role_policy = data.aws_iam_policy_document.lambda-dynamo-read.json
}

resource "aws_iam_role" "lambda-exec-full-db" {
  name = "lt-lambda-full-db"
  assume_role_policy = data.aws_iam_policy_document.lambda-dynamo-full.json
}