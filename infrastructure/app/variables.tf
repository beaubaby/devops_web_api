variable "account_id" {}
variable "environment_name" {}
#variable "application_image_url" {}
variable "rds_snapshot_to_restore" {
  default = ""
}

variable "create_secret" {
  default     = true
  description = "If false, this module does nothing (since tf doesn't support conditional modules)"
  type        = bool
}