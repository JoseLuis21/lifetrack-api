/* Category */
// Builder - Deploy to bucket
resource "null_resource" "add-category-upload" {
  provisioner "local-exec" {
    command     = "/bin/bash ../build/deploy-to-s3.sh add-category ${var.app_version}"
    interpreter = ["/bin/bash", "-c"]
  }
}

resource "null_resource" "list-category-upload" {
  provisioner "local-exec" {
    command     = "/bin/bash ../build/deploy-to-s3.sh list-category ${var.app_version}"
    interpreter = ["/bin/bash", "-c"]
  }
}

resource "null_resource" "get-category-upload" {
  provisioner "local-exec" {
    command     = "/bin/bash ../build/deploy-to-s3.sh get-category ${var.app_version}"
    interpreter = ["/bin/bash", "-c"]
  }
}

/* Activity */

/* Occurrence */