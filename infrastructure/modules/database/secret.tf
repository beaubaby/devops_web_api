resource "random_password" "password" {
  length           = 17
  special          = true
  override_special = "!#"
}

resource "aws_secretsmanager_secret" "core_db_master_password" {
  name                    = "${var.environment_name}/coreplatform-db-secrets"
  recovery_window_in_days = 0
}

resource "aws_secretsmanager_secret_version" "core_db_master_password" {
  secret_id     = aws_secretsmanager_secret.core_db_master_password.id
  secret_string = random_password.password.result
}

data "aws_secretsmanager_secret" "core_db_master_password" {
  arn = aws_secretsmanager_secret.core_db_master_password.arn
}

data "aws_secretsmanager_secret_version" "core_db_master_password" {
  secret_id = data.aws_secretsmanager_secret.core_db_master_password.id
}