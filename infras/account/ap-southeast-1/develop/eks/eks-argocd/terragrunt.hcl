include "env" {
  path = find_in_parent_folders("env.hcl")
  expose = true
}

terraform {
  source = "${get_repo_root()}/infras/module/eks-helm"
}

dependency "eks_cluster" {
  config_path = "../eks-cluster"
}

dependency "argocd_admin_pass" {
  config_path = "../argocd-admin-pass"
}

generate "helm_provider" {
  path      = "helm_provider.tf"
  if_exists = "overwrite_terragrunt"
  contents  = <<EOF
    provider "helm" {
      kubernetes {
        host     = "${dependency.eks_cluster.outputs.cluster_endpoint}"
        cluster_ca_certificate = base64decode("${dependency.eks_cluster.outputs.cluster_certificate_authority_data}")
        exec {
          api_version = "client.authentication.k8s.io/v1beta1"
          args        = ["eks", "get-token", "--cluster-name", "${dependency.eks_cluster.outputs.cluster_name}"]
          command     = "aws"
        }
      }
    }
EOF
}

locals {
  environment = "${include.env.locals.environment}"
  system     = "${include.env.locals.system}"
  domain = "${include.env.locals.domain}"

  # ArgoCD domain
  argocd_domain = "argocd-${local.system}.${local.domain}"
  
  # Generate bcrypt hashes using local execution
  argocd_admin_password = dependency.argocd_admin_pass.outputs.random_password_bcrypt
}

inputs = {
  release_name  = "argocd-${local.environment}"
  repository    = "https://argoproj.github.io/argo-helm"
  chart_name    = "argo-cd"
  chart_version = "8.3.0"
  namespace     = "argocd"
  values_file   = "values.yaml"
  values_variables = {
    argocd_domain = "${local.argocd_domain}"
    argocd_admin_password = "${local.argocd_admin_password}"
  }
}
