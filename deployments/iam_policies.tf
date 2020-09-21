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

// -- IAM Policy Statements --

data "aws_iam_policy_document" "lambda-exec" {
  statement {
    actions = ["sts:AssumeRole"]
    effect = "Allow"
    principals {
      type = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}

data "aws_iam_policy_document" "dynamo-category-read" {
  statement {
    actions = ["dynamodb:GetItem", "dynamodb:Scan", "dynamodb:Query"]
    resources = [aws_dynamodb_table.lt-category.arn]
    effect = "Allow"
  }
}

data "aws_iam_policy_document" "dynamo-category-write" {
  statement {
    actions = ["dynamodb:PutItem", "dynamodb:DeleteItem", "dynamodb:UpdateItem"]
    resources = [aws_dynamodb_table.lt-category.arn]
    effect = "Allow"
  }
}

// -- IAM Policies --

resource "aws_iam_policy" "dynamo-category-read" {
  name = "lt-dynamo-category-read"
  description = "Allow read operations to DynamoDB category service table."
  policy = data.aws_iam_policy_document.dynamo-category-read.json
}

resource "aws_iam_policy" "dynamo-category-write" {
  name = "lt-dynamo-category-write"
  description = "Allow write operations to DynamoDB category service table."
  policy = data.aws_iam_policy_document.dynamo-category-write.json
}


// -- IAM Roles --

resource "aws_iam_role" "category-lambda-exec-write-db" {
  name = "lt-category-lambda-write-db"
  description = "Allows Neutrino's LifeTrack lambda functions to call DynamoDB write operations."
  path = "/"
  assume_role_policy = data.aws_iam_policy_document.lambda-exec.json
  tags = {
    Environment: "prod",
    Name: "neutrino-lifetrack"
  }
}

resource "aws_iam_role" "category-lambda-exec-read-db" {
  name = "lt-category-lambda-read-db"
  description = "Allows Neutrino's LifeTrack lambda functions to call DynamoDB read operations."
  path = "/"
  assume_role_policy = data.aws_iam_policy_document.lambda-exec.json
  tags = {
    Environment: "prod",
    Name: "neutrino-lifetrack"
  }
}

resource "aws_iam_role" "category-lambda-exec-full-db" {
  name = "lt-category-lambda-full-db"
  description = "Allows Neutrino's LifeTrack lambda functions to call DynamoDB read-write operations."
  path = "/"
  assume_role_policy = data.aws_iam_policy_document.lambda-exec.json
  tags = {
    Environment: "prod",
    Name: "neutrino-lifetrack"
  }
}

// -- IAM Policy attachment --

// Write Lambda-Dynamo-XRay-CloudWatch
resource "aws_iam_role_policy_attachment" "category-write-role-policy-attachment" {
  role = aws_iam_role.category-lambda-exec-write-db.name
  policy_arn = aws_iam_policy.dynamo-category-write.arn
}

resource "aws_iam_role_policy_attachment" "category-write-xray-write-only-access" {
  role       = aws_iam_role.category-lambda-exec-write-db.name
  policy_arn = "arn:aws:iam::aws:policy/AWSXrayWriteOnlyAccess"
}

resource "aws_iam_role_policy_attachment" "category-write-lambda-write-only-access" {
  role = aws_iam_role.category-lambda-exec-write-db.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

// Read Lambda-Dynamo-XRay-CloudWatch
resource "aws_iam_role_policy_attachment" "category-read-role-policy-attachment" {
  role = aws_iam_role.category-lambda-exec-read-db.name
  policy_arn = aws_iam_policy.dynamo-category-read.arn
}

resource "aws_iam_role_policy_attachment" "category-read-xray-write-only-access" {
  role       = aws_iam_role.category-lambda-exec-read-db.name
  policy_arn = "arn:aws:iam::aws:policy/AWSXrayWriteOnlyAccess"
}

resource "aws_iam_role_policy_attachment" "category-read-lambda-write-only-access" {
  role = aws_iam_role.category-lambda-exec-read-db.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

// Full (Read/Write) Lambda-Dynamo-XRay-CloudWatch
resource "aws_iam_role_policy_attachment" "category-full-read-role-policy-attachment" {
  role = aws_iam_role.category-lambda-exec-full-db.name
  policy_arn = aws_iam_policy.dynamo-category-read.arn
}

resource "aws_iam_role_policy_attachment" "category-full-write-role-policy-attachment" {
  role = aws_iam_role.category-lambda-exec-full-db.name
  policy_arn = aws_iam_policy.dynamo-category-write.arn
}

resource "aws_iam_role_policy_attachment" "category-full-xray-write-only-access" {
  role       = aws_iam_role.category-lambda-exec-full-db.name
  policy_arn = "arn:aws:iam::aws:policy/AWSXrayWriteOnlyAccess"
}

resource "aws_iam_role_policy_attachment" "category-full-lambda-write-only-access" {
  role = aws_iam_role.category-lambda-exec-full-db.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}
