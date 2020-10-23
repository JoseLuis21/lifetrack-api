/* Category */
// Builder - Deploy to bucket
resource "null_resource" "add-category-upload" {
  provisioner "local-exec" {
    command     = "/bin/bash ../build/deploy-to-s3.sh category add-category ${var.app_version}"
    interpreter = ["bash", "-c"]
  }
}

resource "null_resource" "list-category-upload" {
  provisioner "local-exec" {
    command     = "/bin/bash ../build/deploy-to-s3.sh category list-category ${var.app_version}"
    interpreter = ["bash", "-c"]
  }
}

resource "null_resource" "get-category-upload" {
  provisioner "local-exec" {
    command     = "/bin/bash ../build/deploy-to-s3.sh category get-category ${var.app_version}"
    interpreter = ["bash", "-c"]
  }
}

resource "null_resource" "edit-category-upload" {
  provisioner "local-exec" {
    command     = "/bin/bash ../build/deploy-to-s3.sh category edit-category ${var.app_version}"
    interpreter = ["bash", "-c"]
  }
}

resource "null_resource" "change-state-category-upload" {
  provisioner "local-exec" {
    command     = "/bin/bash ../build/deploy-to-s3.sh category change-state-category ${var.app_version}"
    interpreter = ["bash", "-c"]
  }
}

resource "null_resource" "remove-category-upload" {
  provisioner "local-exec" {
    command     = "/bin/bash ../build/deploy-to-s3.sh category remove-category ${var.app_version}"
    interpreter = ["bash", "-c"]
  }
}

/* Activity */

/* Occurrence */