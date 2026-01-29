provider "aws" {
  region = "us-east-1"
}

resource "aws_vpc" "vpc_2" {
  cidr_block = "10.0.0.0/16"
  tags = {
    Name  =  "project-vpc"
  }
}

resource "aws_subnet" "subnet_4" {
  availability_zone = "us-east-1b"
  cidr_block        = "10.0.144.0/16"
  tags = {
    Name  =  "project-subnet-private2-us-east-1b"
  }
  vpc_id = aws_vpc.vpc_2.id
}

resource "aws_security_group" "security_group_6" {
  description = "Default SG"
  name        = "default-sg"
  tags = {
    Name  =  "default-sg"
  }
  vpc_id = aws_vpc.vpc_2.id
}

resource "aws_route_table" "route_table_3" {
  tags = {
    Name  =  "Main Route Table"
  }
  vpc_id = aws_vpc.vpc_2.id
}

resource "aws_instance" "ec2_5" {
  ami           = "ami-0123456789"
  instance_type = "t3.micro"
  subnet_id     = aws_subnet.subnet_4.id
  tags = {
    Name  =  "web-server"
  }
}

