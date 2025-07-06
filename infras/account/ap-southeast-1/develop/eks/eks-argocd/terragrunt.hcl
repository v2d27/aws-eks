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
  argocd_domain = "${local.system}-argocd.${local.domain}"
  
  # Generate bcrypt hashes using local execution
  admin_password_plain = "admin-password-plain"
  devuser_password_plain = "devuser-password-plain"
  webhook_secret_plain = "webhook-secret-plain"
  
  # Generate bcrypt hashes using htpasswd command
  admin_password_bcrypt = run_cmd("--terragrunt-quiet", "bash", "-c", "htpasswd -bnBC 10 '' '${local.admin_password_plain}' | tr -d ':\\n' | cut -d: -f2")
  devuser_password_bcrypt = run_cmd("--terragrunt-quiet", "bash", "-c", "htpasswd -bnBC 10 '' '${local.devuser_password_plain}' | tr -d ':\\n' | cut -d: -f2")
  
  # Generate current timestamp for passwordMtime
  current_timestamp = run_cmd("--terragrunt-quiet", "date", "-u", "+%Y-%m-%dT%H:%M:%SZ")
}

inputs = {
  release_name  = "argocd-${local.environment}"
  repository    = "https://argoproj.github.io/argo-helm"
  chart_name    = "argo-cd"
  chart_version = "8.1.2"
  namespace     = "argocd"
  values_file   = "values.yaml"
  values_variables = {
    argocd_domain = "${local.argocd_domain}"
    admin_password_bcrypt = "${local.admin_password_bcrypt}"
    devuser_password_bcrypt = "${local.devuser_password_bcrypt}"
    webhook_secret = "${local.webhook_secret_plain}"
  }
}
