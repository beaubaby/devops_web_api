provider "aws" {
  region = "ap-southeast-1"

  assume_role {
    role_arn = "arn:aws:iam::${var.account_id}:role/deploy-app"
  }
}

terraform {
  backend "s3" {
    bucket         = "ks-terraform-state-688318228301"
    dynamodb_table = "terraform_state"
    key            = "loan-db-secret.tfstate" // this key needs to be unique for each terraform project.
    region         = "ap-southeast-1"
    role_arn       = "arn:aws:iam::688318228301:role/Access-Terraform-State"
  }
}

locals {
  log_shipping_lambda_arn = "arn:aws:lambda:ap-southeast-1:259510286099:function:datadog_log_shipping"
}
