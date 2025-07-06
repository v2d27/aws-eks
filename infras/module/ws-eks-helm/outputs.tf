output "status" {
    description = "The status of the Helm release"
    value = helm_release.this.status
}

output "metadata" {
    description = "The metadata of the Helm release"
    value = helm_release.this.metadata
}