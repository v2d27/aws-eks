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

# https://github.com/kubernetes-sigs/aws-ebs-csi-driver/blob/master/charts/aws-ebs-csi-driver/values.yaml

inputs = {
  release_name  = "ebs-csi-driver"
  repository    = "https://kubernetes-sigs.github.io/aws-ebs-csi-driver"
  chart_name    = "aws-ebs-csi-driver"
  chart_version = "2.45.1" # https://github.com/kubernetes-sigs/aws-ebs-csi-driver/releases
  namespace     = "kube-system"
  values_file   = "values.yaml"
}
