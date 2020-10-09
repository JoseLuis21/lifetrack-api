terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
  }
}

provider "aws" {
  region = "us-east-1"
}


/* Variables */
variable "app_name" {
  default = "neutrino-lifetrack"
  description = "<Application Name>"
}

variable "app_version" {
  default = "1.0.0"
  description = "<Application Version>"
}

variable "app_stage" {
  default = "prod"
  description = "<Application Deployment Version>"
}