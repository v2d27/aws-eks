generate "provider" {
  path      = "provider.tf"
  if_exists = "overwrite_terragrunt"
  contents  = <<EOF
provider "aws" {}
EOF
}


# Remote state configuration
remote_state {
  backend = "s3"
  config = {
    bucket         = "aws-eks-terraform-state" # Update this with your actual bucket name
    key            = "${path_relative_to_include()}/terraform.tfstate"
    region         = "ap-southeast-1" # Region to store the state only
    encrypt        = true 
    dynamodb_table = "aws-eks-terraform-state-lock"
  }
  generate = {
    path      = "backend.tf"
    if_exists = "overwrite_terragrunt"
  }
}

# Terraform version constraints should be defined in versions.tf in your Terraform code, not in Terragrunt.
# Remove this block and add the following to your Terraform versions.tf file:

# terraform {
#   required_version = "= 1.12.2"
#   required_providers {
#     aws = {
#       source  = "hashicorp/aws"
#       version = ">= 6.0"
#     }
#   }
# Note: With the latest Terragrunt, the "generate" block is only valid at the root level, not inside other blocks.
# Your usage is correct for the provider block above.
# For remote_state, the "generate" argument is deprecated and should be removed.
# Terragrunt will automatically generate the backend config file.

# Remove the following from your remote_state block:
# generate = {
#   path      = "backend.tf"
#   if_exists = "overwrite_terragrunt"
# }
