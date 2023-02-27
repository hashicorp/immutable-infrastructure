# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

variable "region" {
  type    = string
  default = "us-east-1"
}

variable "version" {
  type    = string
  default = "1.0.2"
}
//make 

locals { timestamp = regex_replace(timestamp(), "[- TZ:]", "") }


# source blocks are generated from your builders; a source can be referenced in
# build blocks. A build block runs provisioners and post-processors on a
# source.
source "amazon-ebs" "immutable-infrastructure" {
  ami_name      = "immutable-infrastructure-app${var.version}"
  instance_type = "t2.micro"
  region        = var.region
  source_ami_filter {
    filters = {
      name                = "ubuntu/images/*ubuntu-xenial-16.04-amd64-server-*"
      root-device-type    = "ebs"
      virtualization-type = "hvm"
    }
    most_recent = true
    owners      = ["099720109477"]
  }
  ssh_username = "ubuntu"
}

# a build block invokes sources and runs provisioning steps on them.
build {

  hcp_packer_registry {
    bucket_name = "immutable-infrastructure"
    description = <<EOT
      Immutable Infrastructure version 1.
    EOT
    bucket_labels = {
      "owner"          = "immutable-infrastructure-team"
      "os"             = "Ubuntu",
      "ubuntu-version" = "Focal 20.04",
    }
  }
  sources = ["source.amazon-ebs.immutable-infrastructure"]
  provisioner "file" {
    source      = "../tf-packer.pub"
    destination = "/tmp/tf-packer.pub"
  }
  provisioner "shell" {
    script = "../scripts/setup.sh"
  }

}

