/* Category */

resource "aws_dynamodb_table" "lifetrack-prod" {
  hash_key       = "PK"
  range_key      = "SK"
  name           = "lifetrack-prod"
  billing_mode   = "PROVISIONED"
  read_capacity  = 5
  write_capacity = 5
  attribute {
    name = "PK"
    type = "S"
  }
  attribute {
    name = "SK"
    type = "S"
  }

  global_secondary_index {
    hash_key        = "GSIPK"
    range_key       = "GSISK"
    name            = "GSIPK-index"
    projection_type = "ALL"
    read_capacity   = 5
    write_capacity  = 5
  }

  tags = {
    Name : var.app_name
    Version : var.app_version
    Environment : var.app_stage
  }
}
