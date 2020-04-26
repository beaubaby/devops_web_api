data "aws_kms_alias" "db_secrets" { name = "alias/${var.environment_name}/db-secrets" }
locals {
  loan_eligibility_db_secret = {
    db_user     = "LoanDBUser",
    db_password = "${random_string.password.result}"
  }
}

resource "random_string" "password" {
  length           = 17
  special          = true
  override_special = "!#"
}

resource "aws_secretsmanager_secret" "loan_eligibility_db_secret" {
  count         = var.create_secret ? 1 : 0
  name          = "${var.environment_name}/db-secrets"
  kms_key_id    = "${data.aws_kms_alias.db_secrets.target_key_id}"
}

resource "aws_secretsmanager_secret_version" "secret_version" {
  count         = var.create_secret ? 1 : 0
  secret_id     = aws_secretsmanager_secret.loan_eligibility_db_secret[0].id
  secret_string = random_string.password.result
}

# resource "aws_secretsmanager_secret" "loan_eligibility_db_secret" {
#   name = "${var.environment_name}/loan-eligibility/db-secrets"
#   rotation_lambda_arn = "${aws_lambda_function.secretmanagerlambda.arn}"
#   rotation_rules {
#     automatically_after_days = 30
#   }
#   kms_key_id = "${data.aws_kms_alias.db_secrets.target_key_id}"
# }

# resource "aws_secretsmanager_secret_version" "loan_eligibility_db_secret_version" {
#   secret_id = "${aws_secretsmanager_secret.loan_eligibility_db_secret.id}"
#   secret_string = "${jsonencode(map("username", aws_rds_cluster.loan_eligibility_database_cluster.master_username,
#     "password", aws_rds_cluster.loan_eligibility_database_cluster.master_password,
#     "dbname", aws_rds_cluster.loan_eligibility_database_cluster.database_name,
#     "host", aws_rds_cluster.loan_eligibility_database_cluster.endpoint,
#   "engine", "mysql"))}"
#   lifecycle {
#     ignore_changes = ["secret_string"]
#   }
# }

##lamdba function
# data "aws_iam_role" "iam_for_lambda" {
#   name = "${var.environment_name}-SecretsManagerRDSMySQLRotation"
# }

# resource "aws_lambda_function" "secretmanagerlambda" {
#   filename = "SecretsManagerRDSMySQLRotationSingleUser.zip"
#   function_name = "${var.environment_name}-SecretsManagerRDSMySQLRotationSingleUser"
#   handler = "lambda_function.lambda_handler"
#   timeout = 30
#   role = "${data.aws_iam_role.iam_for_lambda.arn}"
#   runtime = "python2.7"
#   source_code_hash = "${base64sha256(filesha256("SecretsManagerRDSMySQLRotationSingleUser.zip"))}"
#   vpc_config {
#     security_group_ids = ["${aws_security_group.loan_eligibility_backend.id}", "${aws_security_group.lambda_to_rds.id}"]
#     subnet_ids = "${aws_db_subnet_group.aurora_subnet_group.subnet_ids}"
#   }
#   environment {
#     variables = {
#       SECRETS_MANAGER_ENDPOINT = "https://secretsmanager.ap-southeast-1.amazonaws.com/"
#       AWS_DATA_PATH = "models"
#     }
#   }
#   description = "Conducts an AWS SecretsManager secret rotation for RDS MySQL using single user rotation scheme"

#   tags = {
#     Name = "loan-eligibility DB Secret Rotation Lambda"
#     App = "loan-eligibility-backend"
#     Component = "loan-eligibility-secrets"
#     Environment = "${var.environment_name}"
#   }
# }

# resource "aws_lambda_permission" "secretmanagerlambda" {
#   action = "lambda:InvokeFunction"
#   function_name = "${aws_lambda_function.secretmanagerlambda.function_name}"
#   principal = "secretsmanager.amazonaws.com"
#   statement_id = "AllowSecretManagerToCallLambda"
# }

# resource "aws_security_group" "lambda_to_rds" {
#   vpc_id = "${data.aws_vpc.vpc_data.id}"
#   name_prefix = "lambda_to_rds"
#   egress {
#     from_port = 3306
#     protocol = "tcp"
#     to_port = 3306
#     cidr_blocks = ["${data.aws_vpc.vpc_data.cidr_block}"]
#     description = "allow lambda to connect to database rds"
#   }

#   tags = {
#     Name = "loan-eligibility Secret Rotation Lambda -> RDS"
#     App = "loan-eligibility-backend"
#     Component = "loan-eligibility-secrets"
#     Environment = "${var.environment_name}"
#   }
# }

resource "aws_rds_cluster_instance" "cluster_instances" {
  count                = 2
  identifier_prefix    = "${var.environment_name}-aurora-instance-${count.index}"
  cluster_identifier   = "${aws_rds_cluster.loan_eligibility_database_cluster.id}"
  instance_class       = "db.t3.small"
  db_subnet_group_name = "${aws_db_subnet_group.aurora_subnet_group.name}"
  apply_immediately    = true

  tags = {
    Name        = "Loan Eligibility RdsDatabase Instance"
    App         = "Loan Eligibility"
    Component   = "DB"
    Environment = "${var.environment_name}"
  }
  lifecycle {
    prevent_destroy = false                               // TODO: set back to true after rolling out scale-down in prod
    ignore_changes  = ["identifier_prefix", "identifier"] // to be able to gradually roll out change from fixed identifier to identifier prefix
  }
}


resource "aws_rds_cluster" "loan_eligibility_database_cluster" {
  cluster_identifier_prefix = "${var.environment_name}-aurora-cluster"
  database_name             = "LoanEligibilityDB"
  master_username           = "${local.loan_eligibility_db_secret["db_user"]}"
  master_password           = "${local.loan_eligibility_db_secret["db_password"]}"
  db_subnet_group_name      = "${aws_db_subnet_group.aurora_subnet_group.name}"
  vpc_security_group_ids    = ["${aws_security_group.allow-loan-eligibility-to-database.id}"]
  skip_final_snapshot       = true
  storage_encrypted         = true
  kms_key_id                = "${aws_kms_alias.rds.target_key_arn}"
  backup_retention_period   = 30
  snapshot_identifier       = "${var.rds_snapshot_to_restore}"
  backtrack_window          = "${24 * 60 * 60}"

  enabled_cloudwatch_logs_exports = [
    "audit",
    "error",
    "general",
    "slowquery",
  ]
  lifecycle {
    prevent_destroy = false
    // to be able to gradually roll out change from fixed identifier to identifier prefix and adding backtrack_window
    ignore_changes = [
      "backtrack_window",
      "cluster_identifier",
      "cluster_identifier_prefix"
    ]
  }
  tags = {
    Name        = "Loan Eligibility RdsDatabase Cluster"
    App         = "Loan Eligibility"
    Component   = "DB"
    Environment = "${var.environment_name}"
  }
}

resource "aws_db_subnet_group" "aurora_subnet_group" {
  name        = "${var.environment_name}_aurora_db_subnet_group"
  description = "Allowed subnets for Aurora DB cluster instances"
  subnet_ids  = "${data.aws_subnet_ids.private_subnets.ids}"

  tags = {
    Name        = "Loan Eligibility RdsDatabase"
    App         = "Loan Eligibility"
    Component   = "DB"
    Environment = "${var.environment_name}"
  }
}

resource "aws_security_group" "allow-loan-eligibility-to-database" {
  vpc_id      = "${data.aws_vpc.vpc_data.id}"
  name_prefix = "loan-eligibility_instace_to_db"

  ingress {
    from_port = 3306
    protocol  = "tcp"
    to_port   = 3306
    self      = true
  }
  # ingress {
  #   from_port = 3306
  #   protocol  = "tcp"
  #   to_port   = 3306
  #   security_groups = ["${aws_security_group.lambda_to_rds.id}"]
  # }
  egress {
    from_port = 3306
    protocol  = "tcp"
    to_port   = 3306
    self      = true
  }

  tags = {
    Name        = "Loan Eligibility RdsDatabase SecurityGroup"
    App         = "Loan Eligibility"
    Environment = "${var.environment_name}"
  }
  lifecycle {
    create_before_destroy = true
  }
}
