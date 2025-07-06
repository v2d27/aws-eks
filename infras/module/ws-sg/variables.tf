variable "create" {
    description = "Flag to create the security group"
    type        = bool
    default     = true
}

variable "name" {
    description = "Name of the security group"
    type        = string
}

variable "tags" {
    description = "Tags to apply to the security group"
    type        = map(string)
    default     = {}
}

variable "vpc_id" {
    description = "VPC ID to create the security group in"
    type        = string
}

variable "ingress_with_source_security_group_id" {
  description = "List of ingress rules to create where 'source_security_group_id' is used"
  type        = list(map(string))
  default     = []
}

variable "ingress_with_cidr_blocks" {
  description = "List of ingress rules to create where 'cidr_blocks' is used"
  type        = list(map(string))
}

variable "egress_with_cidr_blocks" {
  description = "List of egress rules to create where 'cidr_blocks' is used"
  type        = list(map(string))
}

variable "description" {
  description = "Description of the security group"
  type        = string
  default     = ""
}