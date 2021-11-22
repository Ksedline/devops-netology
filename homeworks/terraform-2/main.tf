provider "aws" {
  region = "us-east-1"
}

resource "aws_vpc" "netology_vpc" {
  cidr_block = "10.0.0.0/16"
}

data "aws_ami" "ubuntu" {
  most_recent = true

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-focal-20.04-amd64-server-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  filter {
    name   = "root-device-type"
    values = ["ebs"]
  }

  owners = [""] # Скрыл, с точки зрения безопасности аккаунта
}

resource "aws_instance" "netology_instance" {
  ami                         = data.aws_ami.ubuntu.id
  instance_type               = "t3.micro"
  count                       = 1
  cpu_core_count              = 1
  cpu_threads_per_core        = 1
  disable_api_termination     = false
  hibernation                 = true
  monitoring                  = true
  associate_public_ip_address = true
  instance_initiated_shutdown_behavior = "stop"
  
  tags = {
    Name = "Hello_Netology"
  }
}

data "aws_caller_identity" "current" {}
data "aws_region" "current" {}
