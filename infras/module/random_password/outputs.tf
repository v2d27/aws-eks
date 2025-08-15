output "random_password" {
  value     = random_password.this.result
  sensitive = true
}

output "random_password_bcrypt" {
  value     = bcrypt(random_password.this.result)
  sensitive = true
}