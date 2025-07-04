include "env" {
  path = find_in_parent_folders("env.hcl")
  expose = true
}

terraform {
  source = "${get_repo_root()}/infras/module/vpc"
}

locals {
  debug = include.env
}


inputs = {
  name = "${include.env.locals.system}-${include.env.locals.environment}-vpc"
  vpc_cidr = "10.0.0.0/16"
  azs      = ["ap-southeast-1a", "ap-southeast-1b"]

  public_subnets  = ["10.0.0.0/24", "10.0.1.0/24"]
  private_subnets = ["10.0.112.0/20", "10.0.128.0/20"]

  enable_dns_hostnames = true
  enable_dns_support   = true
  enable_flow_logs     = false

  tags = include.env.locals.tags
}