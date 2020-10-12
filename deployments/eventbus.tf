/* Category */
resource "aws_sns_topic" "add-category" {
  name                          = "${var.app_short_name}_category_added"
  display_name                  = "LifeTrack Category added"
  sqs_failure_feedback_role_arn = "arn:aws:iam::aws:policy/service-role/AmazonSNSRole"
  tags = {
    Name : var.app_name
    Version : var.app_version
    Environment : var.app_stage
  }
}

resource "aws_sns_topic" "update-category" {
  name                          = "${var.app_short_name}_category_updated"
  display_name                  = "LifeTrack Category updated"
  sqs_failure_feedback_role_arn = "arn:aws:iam::aws:policy/service-role/AmazonSNSRole"
  tags = {
    Name : var.app_name
    Version : var.app_version
    Environment : var.app_stage
  }
}

resource "aws_sns_topic" "remove-category" {
  name                          = "${var.app_short_name}_category_removed"
  display_name                  = "LifeTrack Category removed"
  sqs_failure_feedback_role_arn = "arn:aws:iam::aws:policy/service-role/AmazonSNSRole"
  tags = {
    Name : var.app_name
    Version : var.app_version
    Environment : var.app_stage
  }
}

resource "aws_sns_topic" "restore-category" {
  name                          = "${var.app_short_name}_category_restored"
  display_name                  = "LifeTrack Category restored"
  sqs_failure_feedback_role_arn = "arn:aws:iam::aws:policy/service-role/AmazonSNSRole"
  tags = {
    Name : var.app_name
    Version : var.app_version
    Environment : var.app_stage
  }
}

resource "aws_sns_topic" "hard-remove-category" {
  name                          = "${var.app_short_name}_category_hard_removed"
  display_name                  = "LifeTrack Category hard removed"
  sqs_failure_feedback_role_arn = "arn:aws:iam::aws:policy/service-role/AmazonSNSRole"
  tags = {
    Name : var.app_name
    Version : var.app_version
    Environment : var.app_stage
  }
}

/* Activity */

/* Occurrence */
