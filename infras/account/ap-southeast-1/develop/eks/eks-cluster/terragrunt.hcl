include "env" {
  path = find_in_parent_folders("env.hcl")
  expose = true
}

terraform {
  source = "${get_repo_root()}/infras/module/eks"
}

dependency "vpc" {
  config_path = "../vpc"
}

dependency "current_user" {
  config_path = "../iam/current-user"
}

locals {
  admin_users = [] # Enter more admin users here, e.g. ["arn:aws:iam::123456789012:user/JohnDoe"]
}

inputs = {
  cluster_name                   = "${include.env.locals.system}-${include.env.locals.environment}-eks"
  cluster_version                = "1.33" # Docs: https://docs.aws.amazon.com/eks/latest/userguide/platform-versions.html
  cluster_endpoint_public_access = true

  vpc_id                   = dependency.vpc.outputs.vpc_id
  subnet_ids               = [dependency.vpc.outputs.public_subnets[0], dependency.vpc.outputs.public_subnets[1]]
  control_plane_subnet_ids = [dependency.vpc.outputs.public_subnets[0], dependency.vpc.outputs.public_subnets[1]]

  cluster_addons = {
    eks-pod-identity-agent = {}
    kube-proxy             = {}
    coredns                = {}
    vpc-cni = {
      before_compute = true
      configuration_values = jsonencode({
        env = {
          ENABLE_PREFIX_DELEGATION = "true"
          WARM_PREFIX_TARGET       = "1"
          WARM_ENI_TARGET          = "1"
        }
      })
    }
  }

  eks_managed_node_group_defaults = {
    iam_role_additional_policies = {
      AmazonEBSCSIDriverPolicy           = "arn:aws:iam::aws:policy/service-role/AmazonEBSCSIDriverPolicy"
      AmazonEC2ContainerRegistryReadOnly = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"
      AmazonEKSWorkerNodePolicy          = "arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy"
      AmazonEKS_CNI_Policy               = "arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy"
      AmazonSSMManagedInstanceCore       = "arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore"
    }
  }

  # EKS Managed Node Group(s)
  # AMI types: https://docs.aws.amazon.com/eks/latest/userguide/retrieve-ami-id.html
  # EKS Instance: https://docs.aws.amazon.com/eks/latest/userguide/automode-learn-instances.html
  eks_managed_node_groups = {
    develop-on-demand = {
      ami_type       = "AL2023_x86_64_STANDARD"
      instance_types = ["t3.medium"]
      min_size     = 0
      max_size     = 2
      desired_size = 1
      autoscaling_group_schedule_arns = {}
    }

    develop-spot = {
      ami_type       = "AL2023_x86_64_STANDARD"
      instance_types = ["t3.medium"]
      capacity_type  = "SPOT"
      min_size     = 0
      max_size     = 2
      desired_size = 1
      autoscaling_group_schedule_arns = {}
    }
  }

  # Cluster access entry
  access_entries = {
    for user in concat(local.admin_users, [dependency.current_user.outputs.current_account_arn]) : user => {
      principal_arn = user
      type          = "STANDARD"

      policy_associations = {
        admin = {
          policy_arn = "arn:aws:eks::aws:cluster-access-policy/AmazonEKSClusterAdminPolicy"
          access_scope = {
            type = "cluster"
          }
        }
      }
    }
  }

  create_cloudwatch_log_group    = false
}