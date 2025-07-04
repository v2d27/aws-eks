# ---------------------------------------------------------------------------------------------------------------------
# TERRAGRUNT CONFIGURATION
# ---------------------------------------------------------------------------------------------------------------------

include "env" {
  path = find_in_parent_folders("env.hcl")
  expose = true
}

# Generate an AWS provider block
generate "provider" {
  path      = "provider.tf"
  if_exists = "overwrite_terragrunt"
  contents  = <<EOF
provider "aws" {
  region = "${inlcude.env.locals.region}"

  # Only these AWS Account IDs may be operated on by this template
  # allowed_account_ids = ["${include.env.locals.account}"]
}
EOF
}

# Remote state configuration
remote_state {
  backend = "s3"
  config = {
    bucket         = "aws-eks-terraform-state" # Update this with your actual bucket name
    key            = "${path_relative_to_include()}/terraform.tfstate"
    region         = "${inlcude.env.locals.region}" # Region to store the state only
    encrypt        = true
    use_lockfile = {
      enabled = true
      lock_table = "aws-eks-terraform-state-lock" # Update this with your actual DynamoDB table name
    }
  }
  generate = {
    path      = "backend.tf"
    if_exists = "overwrite_terragrunt"
  }
}

# Configure what repos to search when you run 'terragrunt catalog'
catalog {}
