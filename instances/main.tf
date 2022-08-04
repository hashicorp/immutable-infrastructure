terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "4.23.0" //this version is depecrated, had issues on mac w m1 chip 
      //use 3.63.0 if this version doesn't work when running terraform apply && terraform init 
    }
    hcp = {
      source  = "hashicorp/hcp"
      version = "0.17.0"
    }
  }

  #Change 

  cloud { 
  organization = "im-infr-proj" //Enter your own workspace form TFC

  workspaces {
   name = "immutable" //Enter your own workspace name  
   }
  }
}
variable "iteration_id" {
  description = "HCP Packer Iteration ID"
}
provider "hcp" {}

provider "aws" {
  region = var.region
}

data "hcp_packer_iteration" "ubuntu" {
  bucket_name = var.hcp_bucket_ubuntu
  channel     = var.hcp_channel
}

data "hcp_packer_image" "ubuntu-us-east-1" {
  bucket_name    = data.hcp_packer_iteration.ubuntu.bucket_name
  cloud_provider = "aws"
  iteration_id   = var.iteration_id 
  region         = var.region
}

resource "aws_vpc" "vpc" {
  cidr_block           = var.cidr_vpc
  enable_dns_support   = true
  enable_dns_hostnames = true
}

resource "aws_internet_gateway" "igw" {
  vpc_id = aws_vpc.vpc.id
}

resource "aws_subnet" "subnet_public" {
  vpc_id     = aws_vpc.vpc.id
  cidr_block = var.cidr_subnet
}

resource "aws_route_table" "rtb_public" {
  vpc_id = aws_vpc.vpc.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.igw.id
  }
}

resource "aws_route_table_association" "rta_subnet_public" {
  subnet_id      = aws_subnet.subnet_public.id
  route_table_id = aws_route_table.rtb_public.id
}

resource "aws_security_group" "sg_22_80" {
  name   = "sg_22"
  vpc_id = aws_vpc.vpc.id

  # SSH access from the VPC
  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 8080
    to_port     = 8080
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_instance" "web" {
  ami                         = data.hcp_packer_image.ubuntu-us-east-1.cloud_image_id
  instance_type               = "t2.micro"
  subnet_id                   = aws_subnet.subnet_public.id
  vpc_security_group_ids      = [aws_security_group.sg_22_80.id]
  associate_public_ip_address = true
  user_data = file("${path.module}/userdata.sh")

  lifecycle {  # This allow you to manip how a resource is built or destroyed within terraform , allows us to create the new resource before destroying the old one this allows us to do zero downtime deployment when building ami 
    create_before_destroy = true 
  }

  tags = {
    Name = "immutable-infrastructure"
  }
}

resource "aws_autoscaling_group" "example" {
  name = "${var.cluster_name}-${aws_instance.web.ami}"
  launch_configuration = aws_instance.web.ami 
  vpc_zone_identifier = aws_instance.web.subnet_id 
  health_check_type = "EC2"
  min_size = var.min_size
  max_size = 3
  min_elb_capacity = var.min_size
  lifecycle {
    create_before_destroy = true
  }
  tag {
    key = "Name"
    value  = var.cluster_name
    propagate_at_launch = true 
  }

}



output "public_ip" {
  description = "URL for check endpoint"
  value = "http://${aws_instance.web.public_ip}:8080"
}
