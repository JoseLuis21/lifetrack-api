resource "aws_dynamodb_table" "lt-category" {
  hash_key = "category_id"
  name = "lt-category"
  billing_mode = "PROVISIONED"
  read_capacity = 5
  write_capacity = 5
  attribute {
    name = "category_id"
    type = "S"
  }

  ttl {
    attribute_name = "TimeToExist"
    enabled = false
  }

  tags = {
    Name = "life-track"
    Environment = "production"
  }
}

