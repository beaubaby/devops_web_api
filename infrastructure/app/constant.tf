locals {
  log_shipping_lambda_arn = "arn:aws:lambda:ap-southeast-1:${var.account_id}:function:datadog_log_shipping"
}