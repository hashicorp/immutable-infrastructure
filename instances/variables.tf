variable "cidr_vpc" {
  description = "CIDR block for the VPC"
  default     = "10.1.0.0/16"
}
variable "cidr_subnet" {
  description = "CIDR block for the subnet"
  default     = "10.1.0.0/24"
}

variable "environment_tag" {
  description = "Environment tag"
  default     = "Immutable Infrastructure"
}

variable "region" {
  description = "The region Terraform deploys your instance"
  default     = "us-east-1"
}

variable "hcp_bucket_ubuntu" {
  description = "The Bucket where our AMI is listed."
  default     = "immutable-infrastructure"
}

variable "hcp_channel" {
  description = "HCP Packer channel name"
  default     = "production"
}

