locals {
  db_secret = {
    db_user = "devops_user",
  }
}

resource "aws_db_instance" "database" {
  depends_on = [aws_secretsmanager_secret_version.core_db_master_password]

  allocated_storage      = 10
  engine                 = "mysql"
  engine_version         = "5.7"
  instance_class         = "db.t2.micro"
  identifier             = "${var.environment_name}-devops-db-instance"
  username               = local.db_secret.db_user
  password               = aws_secretsmanager_secret_version.core_db_master_password.secret_string
  db_subnet_group_name   = var.vpc.database_subnet_group
  vpc_security_group_ids = [var.sg.db]
  skip_final_snapshot    = true

  tags = {
    Name        = "Devops RDS DB"
    Site        = "Devops DB"
    Environment = "${var.environment_name}"
  }
}
