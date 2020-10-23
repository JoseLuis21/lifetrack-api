// -- IAM Policy Statements --
data "aws_iam_policy_document" "lambda-exec" {
  statement {
    actions = ["sts:AssumeRole"]
    effect  = "Allow"
    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}

data "aws_iam_policy_document" "dynamo-read" {
  statement {
    actions   = ["dynamodb:GetItem", "dynamodb:Scan", "dynamodb:Query"]
    resources = [aws_dynamodb_table.lifetrack-prod.arn, "${aws_dynamodb_table.lifetrack-prod.arn}/index/*"]
    effect    = "Allow"
  }
}

data "aws_iam_policy_document" "dynamo-write" {
  statement {
    actions   = ["dynamodb:PutItem", "dynamodb:DeleteItem", "dynamodb:UpdateItem"]
    resources = [aws_dynamodb_table.lifetrack-prod.arn, "${aws_dynamodb_table.lifetrack-prod.arn}/index/*"]
    effect    = "Allow"
  }
}

data "aws_iam_policy_document" "sns-read" {
  statement {
    actions = ["sns:GetTopicAttributes", "sns:GetSubscriptionAttributes", "sns:ListSubscriptions",
    "sns:ListSubscriptionsByTopic", "sns:ListTopics", "sns:Subscribe", "sns:Unsubscribe"]
    resources = [aws_sns_topic.add-category.arn, aws_sns_topic.hard-remove-category.arn,
    aws_sns_topic.remove-category.arn, aws_sns_topic.restore-category.arn, aws_sns_topic.update-category.arn]
    effect = "Allow"
  }
}

data "aws_iam_policy_document" "sns-write" {
  statement {
    actions = ["sns:CreateTopic", "sns:DeleteTopic", "sns:Publish", "sns:SetSubscriptionAttributes",
    "sns:SetTopicAttributes", "sns:ConfirmSubscription"]
    resources = [aws_sns_topic.add-category.arn, aws_sns_topic.hard-remove-category.arn,
    aws_sns_topic.remove-category.arn, aws_sns_topic.restore-category.arn, aws_sns_topic.update-category.arn]
    effect = "Allow"
  }
}

// -- IAM Policies --

resource "aws_iam_policy" "dynamo-read" {
  name        = "lt-dynamo-read"
  description = "Allow read operations to DynamoDB service table"
  policy      = data.aws_iam_policy_document.dynamo-read.json
}

resource "aws_iam_policy" "dynamo-write" {
  name        = "lt-dynamo-write"
  description = "Allow write operations to DynamoDB service table"
  policy      = data.aws_iam_policy_document.dynamo-write.json
}

resource "aws_iam_policy" "sns-read" {
  name        = "lt-sns-read"
  description = "Allow read operations to SNS Neutrino LifeTrack topics"
  policy      = data.aws_iam_policy_document.dynamo-write.json
}

resource "aws_iam_policy" "sns-write" {
  name        = "lt-sns-write"
  description = "Allow write operations to SNS Neutrino LifeTrack topics"
  policy      = data.aws_iam_policy_document.dynamo-write.json
}


// -- IAM Roles --
resource "aws_iam_role" "lambda-exec-write" {
  name               = "lt-lambda-write-db"
  description        = "Allows Neutrino's LifeTrack lambda functions to call DynamoDB, SNS write operations"
  path               = "/"
  assume_role_policy = data.aws_iam_policy_document.lambda-exec.json
  tags = {
    Name : var.app_name
    Version : var.app_version
    Environment : var.app_stage
  }
}

resource "aws_iam_role" "lambda-exec-read" {
  name               = "lt-lambda-read-db"
  description        = "Allows Neutrino's LifeTrack lambda functions to call DynamoDB, SNS read operations"
  path               = "/"
  assume_role_policy = data.aws_iam_policy_document.lambda-exec.json
  tags = {
    Name : var.app_name
    Version : var.app_version
    Environment : var.app_stage
  }
}

resource "aws_iam_role" "lambda-exec-full" {
  name               = "lt-lambda-full-db"
  description        = "Allows Neutrino's LifeTrack lambda functions to call DynamoDB, SNS read-write operations"
  path               = "/"
  assume_role_policy = data.aws_iam_policy_document.lambda-exec.json
  tags = {
    Name : var.app_name
    Version : var.app_version
    Environment : var.app_stage
  }
}

// -- IAM Policy attachment --

// Write Lambda-Dynamo-SNS-XRay-CloudWatch
resource "aws_iam_role_policy_attachment" "write-db-role-policy-attachment" {
  role       = aws_iam_role.lambda-exec-write.name
  policy_arn = aws_iam_policy.dynamo-write.arn
}

resource "aws_iam_role_policy_attachment" "write-sns-role-policy-attachment" {
  role       = aws_iam_role.lambda-exec-write.name
  policy_arn = aws_iam_policy.sns-write.arn
}

resource "aws_iam_role_policy_attachment" "write-xray-write-only-access" {
  role       = aws_iam_role.lambda-exec-write.name
  policy_arn = "arn:aws:iam::aws:policy/AWSXrayWriteOnlyAccess"
}

resource "aws_iam_role_policy_attachment" "write-lambda-write-only-access" {
  role       = aws_iam_role.lambda-exec-write.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

// Read Lambda-Dynamo-SNS-XRay-CloudWatch
resource "aws_iam_role_policy_attachment" "read-db-role-policy-attachment" {
  role       = aws_iam_role.lambda-exec-read.name
  policy_arn = aws_iam_policy.dynamo-read.arn
}

resource "aws_iam_role_policy_attachment" "read-sns-role-policy-attachment" {
  role       = aws_iam_role.lambda-exec-read.name
  policy_arn = aws_iam_policy.sns-read.arn
}

resource "aws_iam_role_policy_attachment" "read-xray-write-only-access" {
  role       = aws_iam_role.lambda-exec-read.name
  policy_arn = "arn:aws:iam::aws:policy/AWSXrayWriteOnlyAccess"
}

resource "aws_iam_role_policy_attachment" "read-lambda-write-only-access" {
  role       = aws_iam_role.lambda-exec-read.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

// Full (Read/Write) Lambda-Dynamo-SNS-XRay-CloudWatch
resource "aws_iam_role_policy_attachment" "full-read-db-role-policy-attachment" {
  role       = aws_iam_role.lambda-exec-full.name
  policy_arn = aws_iam_policy.dynamo-read.arn
}

resource "aws_iam_role_policy_attachment" "full-write-db-role-policy-attachment" {
  role       = aws_iam_role.lambda-exec-full.name
  policy_arn = aws_iam_policy.dynamo-write.arn
}

resource "aws_iam_role_policy_attachment" "full-read-sns-role-policy-attachment" {
  role       = aws_iam_role.lambda-exec-full.name
  policy_arn = aws_iam_policy.sns-read.arn
}

resource "aws_iam_role_policy_attachment" "full-write-sns-role-policy-attachment" {
  role       = aws_iam_role.lambda-exec-full.name
  policy_arn = aws_iam_policy.sns-write.arn
}

resource "aws_iam_role_policy_attachment" "full-xray-write-only-access" {
  role       = aws_iam_role.lambda-exec-full.name
  policy_arn = "arn:aws:iam::aws:policy/AWSXrayWriteOnlyAccess"
}

resource "aws_iam_role_policy_attachment" "full-lambda-write-only-access" {
  role       = aws_iam_role.lambda-exec-full.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

