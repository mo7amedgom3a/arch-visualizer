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
  cidr_block        = "10.0.144.0/24"
  tags = {
    Name  =  "public"
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
}

resource "aws_security_group_rule" "security_group_1_rule_0" {
  cidr_blocks       = ["0.0.0.0/0"]
  from_port         = 22
  protocol          = "tcp"
  security_group_id = aws_security_group.security_group_1.id
  to_port           = 22
  type              = "ingress"
}

resource "aws_security_group_rule" "security_group_1_rule_1" {
  cidr_blocks       = ["0.0.0.0/0"]
  from_port         = 80
  protocol          = "tcp"
  security_group_id = aws_security_group.security_group_1.id
  to_port           = 80
  type              = "ingress"
}

resource "aws_subnet" "subnet_6" {
  availability_zone = "us-east-1b"
  cidr_block        = "10.0.128.0/24"
  tags = {
    Name  =  "private subnet"
  }
  vpc_id = aws_vpc.vpc_2.id
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

resource "aws_nat_gateway" "nat_gateway_5" {
  subnet_id = aws_subnet.subnet_4.id
  tags = {
    Name  =  "my-nat-gateway"
  }
}

resource "aws_instance" "ec2_7" {
  ami           = "ami-0123456789"
  instance_type = "t3.micro"
  subnet_id     = aws_subnet.subnet_6.id
  tags = {
    Name  =  "web-server"
  }
  vpc_security_group_ids = [aws_security_group.security_group_1.id]
}

