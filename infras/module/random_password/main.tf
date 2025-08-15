resource "random_password" "this" {
  length           = 32
  special          = true
  override_special = "#$%&-_=+@"

  keepers = {
    # Use a fixed value to ensure the password is always the same
    id = "v2d27-demo-password"
  }
}