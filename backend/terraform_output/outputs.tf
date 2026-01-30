output "vpc_id" {
  value       = aws_vpc.vpc_2.id
  description = "The ID of the VPC"
}

output "public_subnet_id" {
  value       = aws_subnet.subnet_4.id
  description = "The ID of the public subnet"
}

output "private_subnet_id" {
  value       = aws_subnet.subnet_6.id
  description = "The ID of the private subnet"
}

output "nat_gateway_id" {
  value       = aws_nat_gateway.nat_gateway_5.id
  description = "The ID of the NAT Gateway"
}

output "web_server_private_ip" {
  value       = aws_instance.ec2_7.private_ip
  description = "Private IP address of the web server"
}

output "security_group_id" {
  value       = aws_security_group.security_group_1.id
  description = "The ID of the security group"
}

