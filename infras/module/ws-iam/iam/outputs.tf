output "current_account_id" {
  value = data.aws_caller_identity.current.account_id
}

output "current_user_id" {
  value = data.aws_caller_identity.current.user_id
}

output "current_account_arn" {
  value = data.aws_caller_identity.current.arn
}