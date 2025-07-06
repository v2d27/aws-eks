locals {
    global_vars = read_terragrunt_config(find_in_parent_folders("global_vars.hcl"))
    system = local.global_vars.locals.system
    allowed_ip = local.global_vars.locals.allowed_ip

    region_data = read_terragrunt_config(find_in_parent_folders("region.hcl"))
    region = local.region_data.locals.region
    region_name = local.region_data.locals.region_name

    account = read_terragrunt_config(find_in_parent_folders("account.hcl")).locals.account
    environment = "develop"

    tags = {
        system  = local.system
        region  = local.region
        account = local.account
        environment = local.environment
    }
}