variable "environment_name" {}

variable "vpc" {}

variable "sg" {}

variable "db_config" {
  type = object(
    {
      user     = string
      password = string
      database = string
      hostname = string
      port     = string
    }
  )
}

variable "container_url" {}
