provider "aws" {
  region = "us-east-1"
}

resource "aws_lb" "ecs_alb" {
  depends_on = [aws_subnet.public_subnet_1,
    aws_subnet.public_subnet_2,
  aws_security_group.alb_sg]
  internal           = false
  load_balancer_type = "application"
  name               = "ecs-alb"
  security_groups    = [aws_security_group.alb_sg.id]
  subnets = [aws_subnet.public_subnet_1.id,
  aws_subnet.public_subnet_2.id]
  tags = {
    Name  =  "ecs-alb"
  }
}

resource "aws_ecs_cluster" "main-cluster" {
  name = "main-cluster"
}

resource "aws_lb_target_group" "ecs_tg" {
  name     = "ecs-tg"
  port     = 80
  protocol = "HTTP"
  tags = {
    Name  =  "ecs-tg"
  }
  target_type = "ip"
  vpc_id      = aws_vpc.ecs_vpc.id
  health_check {
    path     = "/health"
    protocol = "HTTP"
  }
}

resource "aws_ecs_task_definition" "web-app" {
  container_definitions    = jsonencode([
  {
    "cpu": 256,
    "essential": true,
    "image": "nginx:latest",
    "memory": 512,
    "name": "web-container"
  }
])
  cpu                      = "256"
  family                   = "web-app"
  memory                   = "512"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
}

resource "aws_ecs_service" "web-service" {
  name = "web-service"
}

resource "aws_vpc" "ecs_vpc" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true
  enable_dns_support   = true
  tags = {
    Name  =  "ecs-vpc"
  }
}

resource "aws_lb_listener" "ecs_listener" {
  load_balancer_arn = aws_lb.ecs_alb.arn
  port              = 80
  protocol          = "HTTP"
  tags = {
    Name  =  "ecs-listener"
  }
  default_action {
    target_group_arn = aws_lb_target_group.ecs_tg.arn
    type             = "forward"
  }
}

resource "aws_internet_gateway" "ecs_igw" {
  tags = {
    Name  =  "ecs-igw"
  }
  vpc_id = aws_vpc.ecs_vpc.id
}

resource "aws_subnet" "public_subnet_1" {
  availability_zone       = "us-east-1a"
  cidr_block              = "10.0.1.0/24"
  map_public_ip_on_launch = true
  tags = {
    Name  =  "public-subnet-1"
  }
  vpc_id = aws_vpc.ecs_vpc.id
}

resource "aws_subnet" "public_subnet_2" {
  availability_zone       = "us-east-1b"
  cidr_block              = "10.0.2.0/24"
  map_public_ip_on_launch = true
  tags = {
    Name  =  "public-subnet-2"
  }
  vpc_id = aws_vpc.ecs_vpc.id
}

resource "aws_subnet" "private_subnet_1" {
  availability_zone = "us-east-1a"
  cidr_block        = "10.0.3.0/24"
  tags = {
    Name  =  "private-subnet-1"
  }
  vpc_id = aws_vpc.ecs_vpc.id
}

resource "aws_subnet" "private_subnet_2" {
  availability_zone = "us-east-1b"
  cidr_block        = "10.0.4.0/24"
  tags = {
    Name  =  "private-subnet-2"
  }
  vpc_id = aws_vpc.ecs_vpc.id
}

resource "aws_security_group" "alb_sg" {
  description = "Allow HTTP from internet"
  name        = "alb-sg"
  tags = {
    Name  =  "alb-sg"
  }
  vpc_id = aws_vpc.ecs_vpc.id
  egress {
    cidr_blocks = ["0.0.0.0/0"]
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

resource "aws_security_group" "ecs_tasks_sg" {
  description = "Allow HTTP from ALB"
  name        = "ecs-tasks-sg"
  tags = {
    Name  =  "ecs-tasks-sg"
  }
  vpc_id = aws_vpc.ecs_vpc.id
  egress {
    cidr_blocks = ["0.0.0.0/0"]
    from_port   = 0
    protocol    = "-1"
    to_port     = 0
  }
  ingress {
    description     = "HTTP from ALB"
    from_port       = 80
    protocol        = "tcp"
    security_groups = [aws_security_group.alb_sg.id]
    to_port         = 80
  }
}

