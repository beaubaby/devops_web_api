variable "environment_name" {}

variable "key_name" {
  default = "my_key_pair"
}

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
