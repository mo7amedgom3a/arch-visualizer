provider "aws" {
  region = "us-east-1"
}

resource "aws_vpc" "main_vpc" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true
  enable_dns_support   = true
  tags = {
    Name  =  "main-vpc"
  }
}

resource "aws_subnet" "public_subnet_1" {
  availability_zone       = "us-east-1a"
  cidr_block              = "10.0.1.0/24"
  map_public_ip_on_launch = true
  tags = {
    Name  =  "public-subnet-1"
  }
  vpc_id = aws_vpc.main_vpc.id
}

resource "aws_subnet" "public_subnet_2" {
  availability_zone       = "us-east-1b"
  cidr_block              = "10.0.2.0/24"
  map_public_ip_on_launch = true
  tags = {
    Name  =  "public-subnet-2"
  }
  vpc_id = aws_vpc.main_vpc.id
}

resource "aws_subnet" "private_subnet_1" {
  availability_zone = "us-east-1a"
  cidr_block        = "10.0.3.0/24"
  tags = {
    Name  =  "private-subnet-1"
  }
  vpc_id = aws_vpc.main_vpc.id
}

resource "aws_subnet" "private_subnet_2" {
  availability_zone = "us-east-1b"
  cidr_block        = "10.0.4.0/24"
  tags = {
    Name  =  "private-subnet-2"
  }
  vpc_id = aws_vpc.main_vpc.id
}

resource "aws_internet_gateway" "main_igw" {
  tags = {
    Name  =  "main-igw"
  }
  vpc_id = aws_vpc.main_vpc.id
}

resource "aws_route_table" "public_rtb" {
  tags = {
    Name  =  "public-rtb"
  }
  vpc_id = aws_vpc.main_vpc.id
}

resource "aws_route_table" "private_rtb" {
  tags = {
    Name  =  "private-rtb"
  }
  vpc_id = aws_vpc.main_vpc.id
}

resource "aws_security_group" "ec2_sg" {
  description = "Allow HTTP"
  name        = "ec2-sg"
  tags = {
    Name  =  "ec2-sg"
  }
  vpc_id = aws_vpc.main_vpc.id
  egress {
    cidr_blocks = ["0.0.0.0/0"]
    description = "Allow all outbound traffic"
    from_port   = 0
    protocol    = "-1"
    to_port     = 0
  }
  ingress {
    cidr_blocks = ["0.0.0.0/0"]
    description = "HTTP"
    from_port   = 80
    protocol    = "tcp"
    to_port     = 80
  }
}

resource "aws_security_group" "rds_sg" {
  description = "Allow EC2"
  name        = "rds-sg"
  tags = {
    Name  =  "rds-sg"
  }
  vpc_id = aws_vpc.main_vpc.id
  egress {
    cidr_blocks = ["0.0.0.0/0"]
    description = "Allow all outbound traffic"
    from_port   = 0
    protocol    = "-1"
    to_port     = 0
  }
  ingress {
    cidr_blocks = ["10.0.0.0/16"]
    description = "PostgreSQL"
    from_port   = 5432
    protocol    = "tcp"
    to_port     = 5432
  }
}

resource "aws_eip" "bec58f77_5b81_4e3d_8501_b18f96f4795f_eip" {
  domain = "vpc"
  tags = {
    Name  =  "nat-gw-eip"
  }
}

resource "aws_nat_gateway" "nat_gw" {
  allocation_id = aws_eip.bec58f77_5b81_4e3d_8501_b18f96f4795f_eip.id
  subnet_id     = aws_subnet.public_subnet_1.id
  tags = {
    Name  =  "nat-gw"
  }
}

resource "aws_db_instance" "primary_db" {
  allocated_storage       = 20
  backup_retention_period = 7
  db_name                 = "mydb"
  engine                  = "postgres"
  engine_version          = "13.7"
  identifier              = "b87add4d_110a_4bf9_b801_ab6d3240f492"
  instance_class          = "db.t3.micro"
  multi_az                = true
  password                = "password"
  skip_final_snapshot     = true
  tags = {
    Name  =  "primary-db"
  }
  username = "admin"
}

resource "aws_instance" "web_server" {
  ami                         = "ami-12345678"
  associate_public_ip_address = false
  depends_on                  = [aws_db_instance.primary_db]
  instance_type               = "t3.micro"
  key_name                    = "my-key"
  subnet_id                   = aws_subnet.private_subnet_1.id
  tags = {
    Name  =  "web-server"
  }
}

