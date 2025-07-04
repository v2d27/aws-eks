include "env" {
  path = find_in_parent_folders("env.hcl")
  expose = true
}

terraform {
  source = "${get_repo_root()}/infras/module/vpc"
}

inputs = {
  name = "${include.env.locals.system}-${include.env.locals.environment}-vpc"
  vpc_cidr = "10.0.0.0/16"

  azs = [
    "ap-southeast-1a",
    "ap-southeast-1b",
    "ap-southeast-1a",
    "ap-southeast-1b"
  ]

  public_subnets = [
    "10.0.112.0/20",  # public-subnet-1 (AZ1)
    "10.0.128.0/20",  # public-subnet-2 (AZ2)
    "10.0.144.0/20",  # public-subnet-3 (AZ1 again)
    "10.0.160.0/20"   # public-subnet-4 (AZ2 again)
  ]

  # private_subnets = [
  #   "10.0.0.0/24", # private-subnet-1 (AZ1)
  #   "10.0.1.0/24", # private-subnet-2 (AZ2)
  #   "10.0.2.0/24", # private-subnet-3 (AZ1 again)
  #   "10.0.3.0/24", # private-subnet-4 (AZ2 again)
  # ]

  enable_dns_hostnames = true
  enable_dns_support   = true
  enable_flow_logs     = false
  map_public_ip_on_launch = true

  tags = include.env.locals.tags
}