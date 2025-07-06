output "security_group_arn" {
  description = "The ARN of the security group"
  value       = local.create ? module.sg[0].security_group_arn : null
}

output "security_group_description" {
  description = "The description of the security group"
  value       = local.create ? module.sg[0].security_group_description : null
}

output "security_group_id" {
  description = "The ID of the security group"
  value       = local.create ? module.sg[0].security_group_id : null
}

output "security_group_name" {
  description = "The name of the security group"
  value       = local.create ? module.sg[0].security_group_name : null
}

output "security_group_owner_id" {
  description = "The owner ID"
  value       = local.create ? module.sg[0].security_group_owner_id : null
}

output "security_group_vpc_id" {
  description = "The VPC ID"
  value       = local.create ? module.sg[0].security_group_vpc_id : null
}
