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
}

inputs = {
  release_name  = "argocd-${local.environment}"
  repository    = "https://argoproj.github.io/argo-helm"
  chart_name    = "argo-cd"
  chart_version = "8.1.2"
  namespace     = "argocd"
  values_file   = "values.yaml"
}
