include {
    path = find_in_parent_folders("root.hcl")
}

include {
    path = find_in_parent_folders("region.hcl")
}

include {
    path = find_in_parent_folders("account.hcl")
}

include {
    path = find_in_parent_folders("global_vars.hcl")
}

locals {
    tags = {
        "system" = local.system
        "region" = local.region
        "account" = local.account
        "env" = "develop"
    }
}