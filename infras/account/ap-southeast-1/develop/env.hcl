locals {
    system = read_terragrunt_config(find_in_parent_folders("global_vars.hcl")).locals.system
    region = read_terragrunt_config(find_in_parent_folders("region.hcl")).locals.region
    account = read_terragrunt_config(find_in_parent_folders("account.hcl")).locals.account
    environment = "develop"

    tags = {
        system  = local.system
        region  = local.region
        account = local.account
        environment = local.environment
    }
}