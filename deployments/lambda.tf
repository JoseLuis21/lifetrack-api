resource "aws_lambda_function" "add-category" {
  function_name = "lt-add-category"
  s3_bucket = "life-track-serverless"
  s3_key = "v1.0.0/add-category.zip"
  handler = "add-category"
  role = aws_iam_role.category-lambda-exec-full-db.arn
  timeout = 15
  runtime = "go1.x"
  environment {
    variables = {
      "LT_TABLE_NAME": aws_dynamodb_table.lt-category.name,
      "LT_TABLE_REGION": "us-east-1"
    }
  }
}
