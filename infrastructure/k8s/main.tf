//provider "aws" {
//  region = "ap-southeast-1"
//
//  assume_role {
//    role_arn = "arn:aws:iam::${local.tools_account_id}:role/Administrator"
//  }
//}
//terraform {
//  backend "s3" {
//    bucket         = "ks-terraform-state-688318228301"
//    dynamodb_table = "terraform_state"
//    key            = "loan-eligibility-service-k8s.tfstate" // this key needs to be unique for each terraform project.
//    region         = "ap-southeast-1"
//    role_arn       = "arn:aws:iam::688318228301:role/Access-Terraform-State"
//  }
//}
//
//provider "kubernetes" {
//  host = "https://4D1E8984AD852724E76BCFA09D17B795.gr7.ap-southeast-1.eks.amazonaws.com"
//}
//
//data "kubernetes_service" "book_service" {
//  metadata {
//    name = "bookservice"
//  }
//}