output "db_config" {
  value = {
//    user     = aws_rds_cluster.devops_database_cluster.master_username
//    password = aws_rds_cluster.devops_database_cluster.master_password
//    database = aws_rds_cluster.devops_database_cluster.database_name
    user     = aws_db_instance.database.username
    password = aws_db_instance.database.password
    database = aws_db_instance.database.name
    hostname = aws_db_instance.database.address
    port     = aws_db_instance.database.port
  }
}