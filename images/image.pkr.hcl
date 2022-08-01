
packer {
  required_plugins {
    amazon = {
      version = ">= 1.1.1"
      source  = "github.com/hashicorp/amazon"
    }
  }
}


variable "region" {
  type    = string
  default = "us-east-1"
}

variable "version" {      //fix this so we can make other new versions instead of manually changing it here
  type    = string
  default = "1.0.3" 
}


variable "aws_tags" {
    type = map(string)
    default = {
        "Name" = "immutable-infrastructure-ubuntu-us-east-1"
        "Environment" = "Hashicorp Workshop"
        "Developer" = "Immutable Infrastructure Interns"
        "Owner" = "production"
        "OS" = "Ubuntu"
        "Version" = "Focal 20.04"
    }
}

locals { timestamp = regex_replace(timestamp(), "[- TZ:]", "") }


# source blocks are generated from your builders; a source can be referenced in
# build blocks. A build block runs provisioners and post-processors on a
# source.
source "amazon-ebs" "immutable-infrastructure" {
  ami_name      = "immutable-infrastructure-app{{timestamp}}_v${var.version}"
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
  tags = var.aws_tags
}




# a build block invokes sources and runs provisioning steps on them.
# this block does not work in github actions error: hcp_packer registry build 
build {
  
  hcp_packer_registry {
    bucket_name = "immutable-infrastructure"  //{your bucket-name}
    description = "Immutable-Infrastructure Demo"
    bucket_labels = var.aws_tags
    build_labels = {
      "build-time" = timestamp(),
      "build-source" = basename(path.cwd)
    }
  }
  
  sources = ["source.amazon-ebs.immutable-infrastructure"]

  // Create directories
  provisioner "shell" {
    inline = ["sudo mkdir /opt/webapp/"]
  }

   // Copy binary to tmp
  provisioner "file" {
    source      = "../app/bin/server"
    destination = "/tmp/"
  }

   // move binary to desired directory
  provisioner "shell" {
    inline = ["sudo mv /tmp/server /opt/webapp/"]
  }

  provisioner "file" {
    source      = "../tf-packer.pub"
    destination = "/tmp/tf-packer.pub"
  }
  provisioner "shell" { //this is useless now since we are building the go app within this whole repo 
    script = "../scripts/setup.sh"
  }

   post-processor "manifest" {
    output     = "packer_manifest.json"
    strip_path = true
    custom_data = {
      iteration_id = packer.iterationID
    }
  }
}

