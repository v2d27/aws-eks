include "env" {
  path = find_in_parent_folders("env.hcl")
  expose = true
}

terraform {
  source = "${get_repo_root()}/infras/module/ws-eks-helm"
}

dependency "eks_cluster" {
  config_path = "../eks-cluster"
}

dependency "sg" {
  config_path = "../security-groups"
}

generate "helm_provider" {
  path      = "helm_provider.tf"
  if_exists = "overwrite_terragrunt"
  contents  = <<EOF
    provider "helm" {
      kubernetes = {
        host     = "${dependency.eks_cluster.outputs.cluster_endpoint}"
        cluster_ca_certificate = base64decode("${dependency.eks_cluster.outputs.cluster_certificate_authority_data}")
        exec = {
          api_version = "client.authentication.k8s.io/v1beta1"
          args        = ["eks", "get-token", "--cluster-name", "${dependency.eks_cluster.outputs.cluster_name}"]
          command     = "aws"
        }
      }
    }
EOF
}

locals {
  system                  = "${include.env.locals.system}"
  environment             = "${include.env.locals.environment}"
  values_file             = "values.yaml"
}

inputs = {
  release_name  = "ingress-nginx-${local.environment}"
  repository    = "https://kubernetes.github.io/ingress-nginx/"
  chart_name    = "ingress-nginx"
  chart_version = "4.12.1" #"4.11.3"
  namespace     = "8ten-${local.environment}"

  values_file         = "${local.values_file}"
  values_variables    = {
    local_balancer_name = "${local.system}-${local.environment}-nlb"
    security_group      = dependency.sg.outputs.security_group_id
  }

  wait_for_completion = false
  timeout             = 300
}