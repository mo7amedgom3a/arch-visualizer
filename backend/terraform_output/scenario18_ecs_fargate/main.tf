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

resource "aws_ecs_cluster_capacity_providers" "main-cluster" {
  capacity_providers = ["ecs-cp"]
  cluster_name       = "main-cluster"
  default_capacity_provider_strategy {
    base              = 1
    capacity_provider = "ecs-cp"
    weight            = 1
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
  execution_role_arn       = "ecs-execution-role"
  family                   = "web-app"
  memory                   = "512"
  network_mode             = "bridge"
  requires_compatibilities = ["EC2"]
}

resource "aws_ecs_capacity_provider" "ecs-ec2-provider" {
  name = "ecs-ec2-provider"
  auto_scaling_group_provider {
    auto_scaling_group_arn         = "ecs-asg"
    managed_termination_protection = "DISABLED"
    managed_scaling {
      status          = "ENABLED"
      target_capacity = 80
    }
  }
}

resource "aws_iam_instance_profile" "ecs-instance-profile" {
  name = "ecs-instance-profile"
  role = "ecs-instance-role"
}

resource "aws_iam_role" "ecs-task-execution-role" {
  assume_role_policy  = "{}"
  managed_policy_arns = ["arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"]
  name                = "ecs-task-execution-role"
}

resource "aws_vpc" "ecs_vpc" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true
  enable_dns_support   = true
  tags = {
    Name  =  "ecs-vpc"
  }
}

resource "aws_lb_target_group" "ecs_tg" {
  name     = "ecs-tg"
  port     = 80
  protocol = "HTTP"
  tags = {
    Name  =  "ecs-tg"
  }
  target_type = "instance"
  vpc_id      = aws_vpc.ecs_vpc.id
  health_check {
    path     = "/health"
    protocol = "HTTP"
  }
}

resource "aws_autoscaling_group" "ecs_asg" {
  depends_on = [aws_subnet.private_subnet_1,
  aws_subnet.private_subnet_2]
  desired_capacity = 2
  max_size         = 4
  min_size         = 1
  name             = "ecs-asg"
  launch_template {
    id      = aws_launch_template.ecs_template.id
    version = "$Latest"
  }
  tag {
    key                 = "Name"
    propagate_at_launch = true
    value               = "ecs-asg"
  }
}

resource "aws_iam_role" "ecs-instance-role" {
  assume_role_policy  = "{}"
  managed_policy_arns = ["arn:aws:iam::aws:policy/service-role/AmazonEC2ContainerServiceforEC2Role"]
  name                = "ecs-instance-role"
}

resource "aws_ecr_repository" "app-repo" {
  force_delete         = true
  image_tag_mutability = "MUTABLE"
  name                 = "app-repo"
  image_scanning_configuration {
    scan_on_push = true
  }
}

resource "aws_launch_template" "ecs_template" {
  image_id      = "ami-0123456789"
  instance_type = "t3.micro"
  name_prefix   = "ecs_template-"
  tags = {
    Name  =  "ecs-template"
  }
  update_default_version = true
  vpc_security_group_ids = [aws_security_group.ecs_tasks_sg.id]
}

resource "aws_ecs_cluster" "main-cluster" {
  name = "main-cluster"
}

resource "aws_ecs_service" "web-service" {
  name = "web-service"
  capacity_provider_strategy {
    capacity_provider = "ecs-ec2-provider"
    weight            = 1
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

