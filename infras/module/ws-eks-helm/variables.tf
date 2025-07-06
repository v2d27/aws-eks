variable "release_name" {
  type        = string
  description = "Helm release name"
}

variable "namespace" {
  description = "Namespace to install helm chart into"
  type        = string
}

variable "chart_name" {
  description = "Name of the helm chart to install"
  type        = string
}



variable "repository" {
  description = "Repository of helm chart"
  type        = string
  default     = null
}

variable "chart_version" {
  description = "Version of helm chart to install"
  type        = string
  default     = null
}

variable "atomic" {
  description = "If set, installation process purges chart on fail. The wait flag will be set automatically if atomic is used"
  type        = bool
  default     = false
}

variable "cleanup_on_fail" {
  description = "Allow deletion of new resources created in this upgrade when upgrade fails"
  type        = bool
  default     = false
}

variable "create_namespace" {
  description = "Create the namespace if it does not exist"
  type        = bool
  default     = false
}

variable "dependency_update" {
  description = "Run helm dependency update before installing the chart"
  type        = bool
  default     = false
}

variable "description" {
  description = "Add a custom description"
  type        = string
  default     = null
}

variable "devel" {
  description = "Use chart development versions, too. Equivalent to version '>0.0.0-0'. If version is set, this is ignored"
  type        = bool
  default     = false
}

variable "disable_crd_hooks" {
  description = "Prevent CRD hooks from running, but run other hooks"
  type        = bool
  default     = false
}

variable "disable_openapi_validation" {
  description = "If set, the installation process will not validate rendered templates against the Kubernetes OpenAPI Schema"
  type        = bool
  default     = false
}

variable "disable_webhooks" {
  description = "Prevent hooks from running"
  type        = bool
  default     = false
}

variable "force_update" {
  description = "Force resource update through delete/recreate if needed"
  type        = bool
  default     = false
}

variable "keyring" {
  description = "Location of public keys used for verification. Used only if verify is true"
  type        = string
  default     = null
}

variable "lint" {
  description = "Run helm lint when planning"
  type        = bool
  default     = false
}

variable "max_history" {
  description = "Limit the maximum number of revisions saved per release. Use 0 for no limit"
  type        = number
  default     = 0
}

variable "pass_credentials" {
  description = "Pass credentials to all domains"
  type        = bool
  default     = false
}

variable "postrender" {
  description = "Postrender command configuration"
  type = object({
    binary_path = string
    args        = optional(list(string))
  })
  default = null
}

variable "recreate_pods" {
  description = "Perform pods restart during upgrade/rollback"
  type        = bool
  default     = false
}

variable "render_subchart_notes" {
  description = "If set, render subchart notes along with the parent"
  type        = bool
  default     = true
}

variable "replace" {
  description = "Re-use the given name, even if that name is already used. This is unsafe in production"
  type        = bool
  default     = false
}

variable "repository_ca_file" {
  description = "The Repositories CA File"
  type        = string
  default     = null
}

variable "repository_cert_file" {
  description = "The repositories cert file"
  type        = string
  default     = null
}

variable "repository_key_file" {
  description = "The repositories cert key file"
  type        = string
  default     = null
}

variable "repository_password" {
  description = "Password for HTTP basic authentication"
  type        = string
  default     = null
  sensitive   = true
}

variable "repository_username" {
  description = "Username for HTTP basic authentication"
  type        = string
  default     = null
}

variable "reset_values" {
  description = "When upgrading, reset the values to the ones built into the chart"
  type        = bool
  default     = false
}

variable "reuse_values" {
  description = "When upgrading, reuse the last release's values and merge in any overrides"
  type        = bool
  default     = false
}

# variable "set" {
#   description = "Custom values to be merged with the values"
#   type = list(object({
#     name  = string
#     value = string
#     type  = optional(string)
#   }))
#   default = []
# }

# variable "set_list" {
#   description = "Custom list values to be merged with the values"
#   type = list(object({
#     name  = string
#     value = list(string)
#   }))
#   default = []
# }

# variable "set_sensitive" {
#   description = "Custom sensitive values to be merged with the values"
#   type = list(object({
#     name  = string
#     value = string
#     type  = optional(string)
#   }))
#   default = []
# }

variable "skip_crds" {
  description = "If set, no CRDs will be installed. By default, CRDs are installed if not already present"
  type        = bool
  default     = false
}

variable "timeout" {
  description = "Time in seconds to wait for any individual kubernetes operation"
  type        = number
  default     = 300
}

# variable "upgrade_install" {
#   description = "If true, the provider will install the release at the specified version even if a release not controlled by the provider is present"
#   type        = bool
#   default     = false
# }

# variable "values" {
#   description = "List of values in raw yaml format to pass to helm"
#   type        = list(string)
#   default     = []
# }

variable "verify" {
  description = "Verify the package before installing it"
  type        = bool
  default     = false
}

variable "wait" {
  description = "Will wait until all resources are in a ready state before marking the release as successful"
  type        = bool
  default     = true
}

variable "wait_for_jobs" {
  description = "If wait is enabled, will wait until all Jobs have been completed before marking the release as successful"
  type        = bool
  default     = false
}

variable "values_file" {
  description = "The name of the helm chart values file to use"
  type        = string
  default     = "values.yaml"
}

variable "values_variables" {
  description = "A map of variables to pass to the helm chart values file"
  type        = map(any)
  default     = {}
}

variable "local_balancer_name" {
  description = "Name of the load balancer to be used in the helm chart values"
  type        = string
  default     = null
}

variable "security_group" {
  description = "Security group ID to be used in the helm chart values"
  type        = string
  default     = null
}
