locals {
  create = var.create
}

module "sg" {
    count = local.create ? 1 : 0

    source  = "terraform-aws-modules/security-group/aws"
    version = "5.3.0"

    name        = "sg_${var.name}"
    description = var.description

    vpc_id      = var.vpc_id

    ingress_with_source_security_group_id = var.ingress_with_source_security_group_id
    ingress_with_cidr_blocks = var.ingress_with_cidr_blocks
    egress_with_cidr_blocks = var.egress_with_cidr_blocks

    tags = var.tags
}