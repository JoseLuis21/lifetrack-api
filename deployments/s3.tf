resource "aws_s3_bucket" "category" {
  bucket = "lifetrack-serverless"
  acl    = "private"

  tags = {
    Name : var.app_name
    Version : var.app_version
    Environment : var.app_stage
  }
}
