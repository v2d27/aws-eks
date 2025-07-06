resource "helm_release" "this" {
  # Required
  name      = var.release_name
  chart     = var.chart_name
  namespace = var.namespace

  # Optional basic configuration
  repository       = var.repository
  version          = var.chart_version
  create_namespace = var.create_namespace
  description      = var.description

  # Behavior flags
  atomic                     = var.atomic
  cleanup_on_fail            = var.cleanup_on_fail
  dependency_update          = var.dependency_update
  devel                      = var.devel
  disable_crd_hooks          = var.disable_crd_hooks
  disable_openapi_validation = var.disable_openapi_validation
  disable_webhooks           = var.disable_webhooks
  force_update               = var.force_update
  lint                       = var.lint
  max_history                = var.max_history
  pass_credentials           = var.pass_credentials
  recreate_pods              = var.recreate_pods
  render_subchart_notes      = var.render_subchart_notes
  replace                    = var.replace
  reset_values               = var.reset_values
  reuse_values               = var.reuse_values
  skip_crds                  = var.skip_crds
  timeout                    = var.timeout
  verify                     = var.verify
  wait                       = var.wait
  wait_for_jobs              = var.wait_for_jobs

  # Repository authentication
  repository_ca_file   = var.repository_ca_file
  repository_cert_file = var.repository_cert_file
  repository_key_file  = var.repository_key_file
  repository_password  = var.repository_password
  repository_username  = var.repository_username
  keyring              = var.keyring

  # Values configuration
  values = [
    fileexists("${path.root}/${var.values_file}") ? templatefile("${path.root}/${var.values_file}", var.values_variables) : ""
  ]
}
