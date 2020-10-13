/* Category */

resource "aws_lambda_function" "add-category" {
  function_name = "${var.app_short_name}-add-category"
  description   = "Neutrino LifeTrack - Add category command"
  s3_bucket     = aws_s3_bucket.category.bucket
  s3_key        = "v${var.app_version}/category/add-category.zip"
  handler       = "add-category"
  role          = aws_iam_role.lambda-exec-full-db.arn
  timeout       = 15
  runtime       = "go1.x"
  environment {
    variables = {
      "LT_DYNAMO_TABLE_NAME" : aws_dynamodb_table.lifetrack-prod.name
      "LT_DYNAMO_TABLE_REGION" : "us-east-1"
      "LT_DYNAMO_EVENT_AWS_REGION" : "us-east-1"
    }
  }
  tags = {
    Name : var.app_name
    Version : var.app_version
    Environment : var.app_stage
  }
  tracing_config {
    mode = "Active"
  }
  depends_on = [
    null_resource.add-category-upload
  ]
}

resource "aws_lambda_permission" "add-category-apigw" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.add-category.function_name
  principal     = "apigateway.amazonaws.com"

  source_arn = "${aws_api_gateway_rest_api.lifeTrack.execution_arn}/*/*"
}

resource "aws_lambda_function" "list-category" {
  function_name = "${var.app_short_name}-list-category"
  description   = "Neutrino LifeTrack - List categories query"
  s3_bucket     = aws_s3_bucket.category.bucket
  s3_key        = "v${var.app_version}/category/list-category.zip"
  handler       = "list-category"
  role          = aws_iam_role.lambda-exec-read-db.arn
  timeout       = 15
  runtime       = "go1.x"
  environment {
    variables = {
      "LT_DYNAMO_TABLE_NAME" : aws_dynamodb_table.lifetrack-prod.name
      "LT_DYNAMO_TABLE_REGION" : "us-east-1"
    }
  }
  tags = {
    Name : var.app_name
    Version : var.app_version
    Environment : var.app_stage
  }
  tracing_config {
    mode = "Active"
  }
  depends_on = [
    null_resource.list-category-upload
  ]
}

resource "aws_lambda_permission" "list-category-apigw" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.list-category.function_name
  principal     = "apigateway.amazonaws.com"

  source_arn = "${aws_api_gateway_rest_api.lifeTrack.execution_arn}/*/*"
}

resource "aws_lambda_function" "get-category" {
  function_name = "${var.app_short_name}-get-category"
  description   = "Neutrino LifeTrack - Get category query"
  s3_bucket     = aws_s3_bucket.category.bucket
  s3_key        = "v${var.app_version}/category/get-category.zip"
  handler       = "get-category"
  role          = aws_iam_role.lambda-exec-read-db.arn
  timeout       = 15
  runtime       = "go1.x"
  environment {
    variables = {
      "LT_DYNAMO_TABLE_NAME" : aws_dynamodb_table.lifetrack-prod.name
      "LT_DYNAMO_TABLE_REGION" : "us-east-1"
    }
  }
  tags = {
    Name : var.app_name
    Version : var.app_version
    Environment : var.app_stage
  }
  tracing_config {
    mode = "Active"
  }
  depends_on = [
    null_resource.get-category-upload
  ]
}

resource "aws_lambda_permission" "get-category-apigw" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.get-category.function_name
  principal     = "apigateway.amazonaws.com"

  source_arn = "${aws_api_gateway_rest_api.lifeTrack.execution_arn}/*/*"
}

resource "aws_lambda_function" "edit-category" {
  function_name = "${var.app_short_name}-edit-category"
  description   = "Neutrino LifeTrack - Edit category command"
  s3_bucket     = aws_s3_bucket.category.bucket
  s3_key        = "v${var.app_version}/category/edit-category.zip"
  handler       = "edit-category"
  role          = aws_iam_role.lambda-exec-full-db.arn
  timeout       = 15
  runtime       = "go1.x"
  environment {
    variables = {
      "LT_DYNAMO_TABLE_NAME" : aws_dynamodb_table.lifetrack-prod.name
      "LT_DYNAMO_TABLE_REGION" : "us-east-1"
      "LT_DYNAMO_EVENT_AWS_REGION" : "us-east-1"
    }
  }
  tags = {
    Name : var.app_name
    Version : var.app_version
    Environment : var.app_stage
  }
  tracing_config {
    mode = "Active"
  }
  depends_on = [
    null_resource.edit-category-upload
  ]
}

resource "aws_lambda_permission" "edit-category-apigw" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.edit-category.function_name
  principal     = "apigateway.amazonaws.com"

  source_arn = "${aws_api_gateway_rest_api.lifeTrack.execution_arn}/*/*"
}

resource "aws_lambda_function" "change-state-category" {
  function_name = "${var.app_short_name}-change-state-category"
  description   = "Neutrino LifeTrack - Change state category command"
  s3_bucket     = aws_s3_bucket.category.bucket
  s3_key        = "v${var.app_version}/category/change-state-category.zip"
  handler       = "change-state-category"
  role          = aws_iam_role.lambda-exec-full-db.arn
  timeout       = 15
  runtime       = "go1.x"
  environment {
    variables = {
      "LT_DYNAMO_TABLE_NAME" : aws_dynamodb_table.lifetrack-prod.name
      "LT_DYNAMO_TABLE_REGION" : "us-east-1"
      "LT_DYNAMO_EVENT_AWS_REGION" : "us-east-1"
    }
  }
  tags = {
    Name : var.app_name
    Version : var.app_version
    Environment : var.app_stage
  }
  tracing_config {
    mode = "Active"
  }
  depends_on = [
    null_resource.change-state-category-upload
  ]
}

resource "aws_lambda_permission" "change-state-category-apigw" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.change-state-category.function_name
  principal     = "apigateway.amazonaws.com"

  source_arn = "${aws_api_gateway_rest_api.lifeTrack.execution_arn}/*/*"
}

resource "aws_lambda_function" "remove-category" {
  function_name = "${var.app_short_name}-remove-category"
  description   = "Neutrino LifeTrack - Remove category command"
  s3_bucket     = aws_s3_bucket.category.bucket
  s3_key        = "v${var.app_version}/category/remove-category.zip"
  handler       = "remove-category"
  role          = aws_iam_role.lambda-exec-full-db.arn
  timeout       = 15
  runtime       = "go1.x"
  environment {
    variables = {
      "LT_DYNAMO_TABLE_NAME" : aws_dynamodb_table.lifetrack-prod.name
      "LT_DYNAMO_TABLE_REGION" : "us-east-1"
      "LT_DYNAMO_EVENT_AWS_REGION" : "us-east-1"
    }
  }
  tags = {
    Name : var.app_name
    Version : var.app_version
    Environment : var.app_stage
  }
  tracing_config {
    mode = "Active"
  }
  depends_on = [
    null_resource.remove-category-upload
  ]
}

resource "aws_lambda_permission" "remove-category-apigw" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.remove-category.function_name
  principal     = "apigateway.amazonaws.com"

  source_arn = "${aws_api_gateway_rest_api.lifeTrack.execution_arn}/*/*"
}

/* Activity */

/* Occurrence */
