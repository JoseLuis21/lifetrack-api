variable "gsi-index" {
  default     = "GSIPK-index"
  description = "Global Secondary Index name"
}

resource "aws_dynamodb_table" "lifetrack-prod" {
  name           = "${var.app_short_name}-prod"
  billing_mode   = "PROVISIONED"
  read_capacity  = 5
  write_capacity = 5
  hash_key       = "PK"
  range_key      = "SK"

  attribute {
    name = "PK"
    type = "S"
  }
  attribute {
    name = "SK"
    type = "S"
  }

  attribute {
    name = "GSIPK"
    type = "S"
  }
  attribute {
    name = "GSISK"
    type = "S"
  }

  global_secondary_index {
    name            = var.gsi-index
    projection_type = "ALL"
    read_capacity   = 5
    write_capacity  = 5
    hash_key        = "GSIPK"
    range_key       = "GSISK"
  }

  tags = {
    Name : var.app_name
    Version : var.app_version
    Environment : var.app_stage
  }
}

resource "aws_appautoscaling_target" "dynamodb_table_read_target" {
  max_capacity       = 10
  min_capacity       = 5
  resource_id        = "table/${aws_dynamodb_table.lifetrack-prod.name}"
  scalable_dimension = "dynamodb:table:ReadCapacityUnits"
  service_namespace  = "dynamodb"
}

resource "aws_appautoscaling_target" "dynamodb_table_write_target" {
  max_capacity       = 10
  min_capacity       = 5
  resource_id        = "table/${aws_dynamodb_table.lifetrack-prod.name}"
  scalable_dimension = "dynamodb:table:WriteCapacityUnits"
  service_namespace  = "dynamodb"
}

resource "aws_appautoscaling_target" "dynamodb_index_read_target" {
  max_capacity       = 10
  min_capacity       = 5
  resource_id        = "table/${aws_dynamodb_table.lifetrack-prod.name}/index/${var.gsi-index}"
  scalable_dimension = "dynamodb:index:ReadCapacityUnits"
  service_namespace  = "dynamodb"
}

resource "aws_appautoscaling_target" "dynamodb_index_write_target" {
  max_capacity       = 10
  min_capacity       = 5
  resource_id        = "table/${aws_dynamodb_table.lifetrack-prod.name}/index/${var.gsi-index}"
  scalable_dimension = "dynamodb:index:WriteCapacityUnits"
  service_namespace  = "dynamodb"
}
