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

  tags = {
    Environment: "prod",
    Name: "neutrino-lifetrack"
  }
}

