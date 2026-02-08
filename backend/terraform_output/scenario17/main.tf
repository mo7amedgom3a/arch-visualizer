provider "aws" {
  region = "us-east-1"
}

resource "aws_autoscaling_group" "webasg" {
  depends_on = [aws_subnet.publicsubnet,
  aws_lb_target_group.webtargetgroup]
  desired_capacity    = 2
  max_size            = 3
  min_size            = 1
  name                = "WebASG"
  target_group_arns   = [aws_lb_target_group.webtargetgroup.arn]
  vpc_zone_identifier = [aws_subnet.publicsubnet.id]
  launch_template {
    id      = aws_launch_template.weblaunchtemplate.id
    version = "$Latest"
  }
  tag {
    key                 = "Name"
    propagate_at_launch = true
    value               = "WebASG"
  }
}

resource "aws_lb" "webalb" {
  depends_on = [aws_subnet.publicsubnet,
  aws_security_group.websg]
  internal           = false
  load_balancer_type = "application"
  name               = "WebALB"
  security_groups    = [aws_security_group.websg.id]
  subnets            = [aws_subnet.publicsubnet.id]
  tags = {
    Name  =  "WebALB"
  }
}

resource "aws_vpc" "mainvpc" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true
  enable_dns_support   = true
  tags = {
    Name  =  "MainVPC"
  }
}

resource "aws_launch_template" "weblaunchtemplate" {
  depends_on    = [aws_security_group.websg]
  image_id      = "ami-0123456789"
  instance_type = "t3.micro"
  name_prefix   = "weblaunchtemplate-"
  tags = {
    Name  =  "WebLaunchTemplate"
  }
  update_default_version = true
  vpc_security_group_ids = [aws_security_group.websg.id]
}

resource "aws_lb_target_group" "webtargetgroup" {
  name     = "WebTargetGroup"
  port     = 80
  protocol = "HTTP"
  tags = {
    Name  =  "WebTargetGroup"
  }
  target_type = "instance"
  vpc_id      = aws_vpc.mainvpc.id
  health_check {
    path = "/"
  }
}

resource "aws_lb_listener" "weblistener" {
  load_balancer_arn = aws_lb.webalb.arn
  port              = 80
  protocol          = "HTTP"
  tags = {
    Name  =  "WebListener"
  }
  default_action {
    target_group_arn = aws_lb_target_group.webtargetgroup.arn
    type             = "forward"
  }
}

resource "aws_subnet" "publicsubnet" {
  availability_zone       = "us-east-1a"
  cidr_block              = "10.0.1.0/24"
  map_public_ip_on_launch = true
  tags = {
    Name  =  "PublicSubnet"
  }
  vpc_id = aws_vpc.mainvpc.id
}

resource "aws_security_group" "websg" {
  description = "Allow HTTP"
  name        = "WebSG"
  tags = {
    Name  =  "WebSG"
  }
  vpc_id = aws_vpc.mainvpc.id
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

