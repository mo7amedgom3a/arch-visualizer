provider "aws" {
  region = var.aws_region
}

resource "aws_vpc" "vpc_2" {
  cidr_block           = var.vpc_cidr
  enable_dns_hostnames = true
  enable_dns_support   = true
  instance_tenancy     = "default"
  tags = {
    Name  =  "project-vpc"
  }
}

resource "aws_subnet" "subnet_4" {
  availability_zone       = "us-east-1a"
  cidr_block              = "10.0.1.0/24"
  map_public_ip_on_launch = true
  tags = {
    Name  =  "public"
  }
  vpc_id = aws_vpc.vpc_2.id
}

resource "aws_route_table" "route_table_10" {
  tags = {
    Name  =  "private-rtb"
  }
  vpc_id = aws_vpc.vpc_2.id
}

resource "aws_route" "route_table_10_route_1" {
  destination_cidr_block = "0.0.0.0/0"
  nat_gateway_id         = aws_nat_gateway.nat_gateway_5.id
  route_table_id         = aws_route_table.route_table_10.id
}

resource "aws_route_table_association" "route_table_10_assoc_0" {
  route_table_id = aws_route_table.route_table_10.id
  subnet_id      = aws_subnet.subnet_6.id
}

resource "aws_subnet" "subnet_6" {
  availability_zone       = "us-east-1b"
  cidr_block              = "10.0.2.0/24"
  map_public_ip_on_launch = false
  tags = {
    Name  =  "private subnet"
  }
  vpc_id = aws_vpc.vpc_2.id
}

resource "aws_security_group" "security_group_1" {
  description = "security group for http"
  name        = "http-sg"
  tags = {
    Name  =  "http-sg"
  }
  vpc_id = aws_vpc.vpc_2.id
  egress {
    cidr_blocks = ["0.0.0.0/0"]
    description = "Allow all outbound traffic"
    from_port   = 0
    protocol    = "-1"
    to_port     = 0
  }
  ingress {
    cidr_blocks = ["0.0.0.0/0"]
    from_port   = 22
    protocol    = "tcp"
    to_port     = 22
  }
  ingress {
    cidr_blocks = ["0.0.0.0/0"]
    from_port   = 80
    protocol    = "tcp"
    to_port     = 80
  }
}

resource "aws_route_table" "route_table_8" {
  tags = {
    Name  =  "public-rtb"
  }
  vpc_id = aws_vpc.vpc_2.id
}

resource "aws_route" "route_table_8_route_1" {
  destination_cidr_block = "0.0.0.0/0"
  gateway_id             = aws_internet_gateway.igw_9.id
  route_table_id         = aws_route_table.route_table_8.id
}

resource "aws_route_table_association" "route_table_8_assoc_0" {
  route_table_id = aws_route_table.route_table_8.id
  subnet_id      = aws_subnet.subnet_4.id
}

resource "aws_internet_gateway" "igw_9" {
  tags = {
    Name  =  "project-igw"
  }
  vpc_id = aws_vpc.vpc_2.id
}

resource "aws_eip" "nat_gateway_5_eip" {
  domain = "vpc"
  tags = {
    Name  =  "my-nat-gateway-eip"
  }
}

resource "aws_nat_gateway" "nat_gateway_5" {
  allocation_id = aws_eip.nat_gateway_5_eip.id
  subnet_id     = aws_subnet.subnet_4.id
  tags = {
    Name  =  "my-nat-gateway"
  }
}

resource "aws_instance" "ec2_7" {
  ami                         = "ami-0123456789"
  associate_public_ip_address = false
  instance_type               = var.instance_type
  subnet_id                   = aws_subnet.subnet_6.id
  tags = {
    Name  =  "web-server"
  }
  vpc_security_group_ids = [aws_security_group.security_group_1.id]
}

