include "env" {
  path = find_in_parent_folders("env.hcl")
  expose = true
}

terraform {
  source = "tfr://v2d27/argocd/argocd?version=1.0.1"
}

dependency "eks_cluster" {
  config_path = "../eks-cluster"
}

dependency "argocd_admin_pass" {
  config_path = "../eks-argocd/argocd-admin-pass"
}

generate "argocd_provider" {
  path      = "argocd_provider.tf"
  if_exists = "overwrite_terragrunt"
  contents  = <<EOF
    provider "argocd" {
      server_addr = "${local.argocd_domain}"
      username   = "admin"
      password   = "${local.argocd_admin_password}"
      insecure = true # Set to false in production environments
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
  argocd_admin_password = dependency.argocd_admin_pass.outputs.random_password
}

inputs = {
  repo_url        = "git@github.com:v2d27/aws-eks.git"
  git_ssh_private_key = file("~/.ssh/id_rsa") # Path to your SSH private key
  
  app_name        = "chat-app-argocd"
  project_name    = "aws-eks"
  kustomize_path  = "manifest/argocd-appsets/overlays/main"
  app_namespace   = "chat-main"

  project_destinations = [
    {
        server    = "https://kubernetes.default.svc"
        namespace = "argocd"
    },
    {
        server    = "https://kubernetes.default.svc"
        namespace = "kube-system"
    }
  ]
  project_cluster_resource_whitelist = [
    {
        group = "rbac.authorization.k8s.io"
        kind  = "ClusterRole"
    },
    {
        group = "rbac.authorization.k8s.io"
        kind  = "ClusterRoleBinding"
    },
    {
        group = ""
        kind  = "Namespace"
    }
  ]
}
