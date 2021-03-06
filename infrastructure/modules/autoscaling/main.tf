module "iam_instance_profile" {
  source  = "scottwinkler/iip/aws"
  actions = ["logs:*", "rds:*"]
}

//data "template_cloudinit_config" "config" {
//  gzip          = true
//  base64_encode = true
//  part {
//    content_type = "text/cloud-config"
//    content      = templatefile("${path.module}/cloud_config.yaml", var.db_config)
//  }
//}

data "template_file" "deploy_script" {
  template = file("${path.module}/deploy.sh")
  vars = {
    deploy_environment     = var.environment_name
    container_url          = var.container_url
  }
}


data "template_cloudinit_config" "deploy_script" {
  base64_encode = true
  gzip          = false

  part {
    content      = data.template_file.deploy_script.rendered
    content_type = "text/x-shellscript"
  }
}

data "aws_ami" "ubuntu" {
  most_recent = true
  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-bionic-18.04-amd64-server-*"]
  }
  owners = ["099720109477"]
}

resource "aws_launch_template" "webserver" {
  name_prefix   = "${var.environment_name}-devops-web-server"
  image_id      = data.aws_ami.ubuntu.id
  instance_type = "t2.micro"
  user_data     = data.template_cloudinit_config.deploy_script.rendered
  key_name      = aws_key_pair.webserver_key.key_name
  iam_instance_profile {
    name = module.iam_instance_profile.name
  }
  vpc_security_group_ids = [var.sg.websvr]

  tags = {
    Name        = "Devops Launch template"
    Site        = "Devops Web APP"
    Environment = "${var.environment_name}"
  }

  depends_on = ["aws_key_pair.webserver_key"]
}

resource "tls_private_key" "webserver_key" {
  algorithm   =  "RSA"
  rsa_bits    =  4096
}

resource "local_file" "private_key" {
  content         =  tls_private_key.webserver_key.private_key_pem
  filename        =  "webserver.pem"
  file_permission =  0400
}

resource "aws_key_pair" "webserver_key" {
  key_name   = "my_key_pair"
  public_key = tls_private_key.webserver_key.public_key_openssh
}

resource "aws_autoscaling_group" "webserver" {
  name                = "${var.environment_name}-devops-asg"
  min_size            = 1
  max_size            = 3
  vpc_zone_identifier = var.vpc.private_subnets
  target_group_arns   = module.alb.target_group_arns
  launch_template {
    id      = aws_launch_template.webserver.id
    version = aws_launch_template.webserver.latest_version
  }
}

module "alb" {
  source                   = "terraform-aws-modules/alb/aws"
  version                  = "~> 4.0"
  load_balancer_name       = "${var.environment_name}-alb"
  security_groups          = [var.sg.lb]
  subnets                  = var.vpc.public_subnets
  vpc_id                   = var.vpc.vpc_id
  logging_enabled          = false
  http_tcp_listeners       = [{ port = 80, protocol = "HTTP" }]
  http_tcp_listeners_count = "1"
  target_groups            = [{ name = "websvr", backend_protocol = "HTTP", backend_port = 8080 }]
  target_groups_count      = "1"
}