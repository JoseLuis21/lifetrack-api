resource "aws_s3_bucket" "serverless" {
  bucket        = "${var.app_short_name}-serverless"
  acl           = "private"
  region        = "us-east-1"
  force_destroy = true

  versioning {
    enabled    = true
    mfa_delete = false
  }

  tags = {
    Name : var.app_name
    Version : var.app_version
    Environment : var.app_stage
  }
}
