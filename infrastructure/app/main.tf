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
    key            = "loan-eligibility-service-ecr.tfstate" // this key needs to be unique for each terraform project.
    region         = "ap-southeast-1"
    role_arn       = "arn:aws:iam::688318228301:role/Access-Terraform-State"
  }
}
