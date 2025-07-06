# Terraform Helm Chart Module - v2d27

This module provides a comprehensive Terraform configuration for deploying Helm charts with extensive customization options.

## Features

- Support for all Helm release configuration options
- Repository authentication
- Custom values via files or inline
- Set values, set_list, and set_sensitive support
- Postrender configuration
- Comprehensive timeout and retry settings
- Namespace management

## Usage

### Basic Usage

```hcl
module "helm_chart" {
  source = "./path/to/this/module"

  release_name = "my-app"
  chart_name   = "nginx"
  namespace    = "default"
  repository   = "https://charts.bitnami.com/bitnami"
  chart_version = "13.2.0"
}
```

### Advanced Usage

```hcl
module "helm_chart" {
  source = "./path/to/this/module"

  release_name = "my-app"
  chart_name   = "nginx"
  namespace    = "production"
  repository   = "https://charts.bitnami.com/bitnami"
  chart_version = "13.2.0"

  # Namespace management
  create_namespace = true

  # Behavior flags
  atomic          = true
  cleanup_on_fail = true
  wait            = true
  wait_for_jobs   = true
  timeout         = 600

  # Custom values
  values = [
    yamlencode({
      service = {
        type = "LoadBalancer"
      }
      ingress = {
        enabled = true
        hosts   = ["example.com"]
      }
    })
  ]

  # Set individual values
  set = [
    {
      name  = "service.type"
      value = "LoadBalancer"
    },
    {
      name  = "replicas"
      value = "3"
      type  = "string"
    }
  ]

  # Set sensitive values
  set_sensitive = [
    {
      name  = "auth.password"
      value = "supersecret"
    }
  ]

  # Repository authentication
  repository_username = "myuser"
  repository_password = "mypassword"
}
```

## Variables

### Required Variables

- `release_name` (string) - Helm release name
- `chart_name` (string) - Name of the helm chart to install

### Optional Variables

#### Basic Configuration
- `namespace` (string) - Namespace to install helm chart into (default: "default")
- `repository` (string) - Repository of helm chart
- `chart_version` (string) - Version of helm chart to install
- `create_namespace` (bool) - Create the namespace if it does not exist (default: false)
- `description` (string) - Add a custom description

#### Behavior Flags
- `atomic` (bool) - If set, installation process purges chart on fail (default: false)
- `cleanup_on_fail` (bool) - Allow deletion of new resources created in this upgrade when upgrade fails (default: false)
- `dependency_update` (bool) - Run helm dependency update before installing the chart (default: false)
- `devel` (bool) - Use chart development versions (default: false)
- `disable_crd_hooks` (bool) - Prevent CRD hooks from running (default: false)
- `disable_openapi_validation` (bool) - Skip OpenAPI Schema validation (default: false)
- `disable_webhooks` (bool) - Prevent hooks from running (default: false)
- `force_update` (bool) - Force resource update through delete/recreate (default: false)
- `lint` (bool) - Run helm lint when planning (default: false)
- `max_history` (number) - Limit maximum number of revisions (default: 0)
- `pass_credentials` (bool) - Pass credentials to all domains (default: false)
- `recreate_pods` (bool) - Perform pods restart during upgrade/rollback (default: false)
- `render_subchart_notes` (bool) - Render subchart notes with parent (default: true)
- `replace` (bool) - Re-use the given name (default: false)
- `reset_values` (bool) - Reset values to chart defaults when upgrading (default: false)
- `reuse_values` (bool) - Reuse last release's values when upgrading (default: false)
- `skip_crds` (bool) - Skip CRD installation (default: false)
- `timeout` (number) - Time in seconds to wait for operations (default: 300)
- `verify` (bool) - Verify the package before installing (default: false)
- `wait` (bool) - Wait until all resources are ready (default: true)
- `wait_for_jobs` (bool) - Wait for all Jobs to complete (default: false)

#### Repository Authentication
- `repository_ca_file` (string) - The Repositories CA File
- `repository_cert_file` (string) - The repositories cert file
- `repository_key_file` (string) - The repositories cert key file
- `repository_password` (string, sensitive) - Password for HTTP basic authentication
- `repository_username` (string) - Username for HTTP basic authentication
- `keyring` (string) - Location of public keys for verification

#### Values Configuration
- `values` (list of strings) - List of values in raw yaml format
- `values_file` (string) - Path to values file (default: "values.yaml")
- `values_variables` (map) - Variables to pass to values file template

## Outputs

- `status` - The status of the Helm release
- `metadata` - The metadata of the Helm release

## Requirements

- Terraform >= 1.0
- Helm provider >= 2.0