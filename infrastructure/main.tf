terraform {
  required_version = "~> 0.12"
  required_providers {
    aws      = "~> 2.13"
    local    = "~> 1.2"
    random   = "~> 2.1"
    template = "~> 2.1"
  }
}

provider "aws" {
  region = var.aws_region
}

data "aws_caller_identity" "current" {}

module "autoscaling" {
  source           = "./modules/autoscaling"
  key_name         = var.key_name

  environment_name = var.environment_name
  vpc              = module.networking.vpc
  sg               = module.networking.sg
  db_config        = module.database.db_config
}

module "database" {
  source = "./modules/database"

  environment_name = var.environment_name
  vpc              = module.networking.vpc
  sg               = module.networking.sg
  account_id       = var.account_id
  database_name    = var.database_name
}

module "networking" {
  source = "./modules/networking"

  environment_name = var.environment_name
}