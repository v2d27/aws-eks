include "env" {
  path = find_in_parent_folders("env.hcl")
  expose = true
}

terraform {
  source = "${get_repo_root()}/infras/module/ws-sg"
}

dependency "vpc" {
  config_path = "../../../vpc"
}

locals {
  system                 = "${include.env.locals.system}"
  environment            = "${include.env.locals.environment}"
  allowed_ip             = "${include.env.locals.allowed_ip}"

  http_rules = [
    {
      description = "Allow HTTP from anywhere"
      from_port   = 80
      to_port     = 80
      protocol    = "tcp"
      cidr_blocks = "0.0.0.0/0"
    },
    {
      description = "Allow HTTPS from anywhere"
      from_port   = 443
      to_port     = 443
      protocol    = "tcp"
      cidr_blocks = "0.0.0.0/0"
    },
  ]

  tcp_rules = [
    for ip in local.allowed_ip : {
      description = "Allow custom IP"
      from_port   = 5000
      to_port     = 5000
      protocol    = "tcp"
      cidr_blocks = "${ip}/32"
    }
  ]

  ingress_values = concat(local.http_rules, local.tcp_rules)

  egress_values = [
    {
      description = "Allow all outbound traffic"
      from_port   = 0
      to_port     = 0
      protocol    = "-1"
      cidr_blocks = "0.0.0.0/0"
    }
  ]

  tags = "${include.env.locals.tags}"
}

inputs = {
  name = "${local.system}-${local.environment}-ingress-nginx-sg"
  description = "Security group for ingress nginx controller"
  vpc_id = dependency.vpc.outputs.vpc_id
  ingress_with_cidr_blocks = local.ingress_values
  egress_with_cidr_blocks = local.egress_values
  
  tags = {
    system      = local.system
    environment = local.environment
  }
}