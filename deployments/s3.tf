resource "aws_s3_bucket" "serverless" {
  bucket = "${var.app_short_name}-serverless"
  acl    = "private"

  tags = {
    Name : var.app_name
    Version : var.app_version
    Environment : var.app_stage
  }
}
