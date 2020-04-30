resource "template_dir" "db_config" {
  source_dir      = "template"
  destination_dir = "output"

  vars = {
    DB_PASSWORD = var.db_password
    DB_USER = var.db_user
    DB_CONNECTION_STRING = var.db_connection_string
    IMAGE_URL = var.image_url
  }
}