output "vpc_id" {
  description = "The ID of the VPC"
  value       = module.vpc.vpc_id
}

output "public_subnet_ids" {
  description = "Public subnet IDs"
  value       = module.vpc.public_subnet_ids
}

output "private_subnet_ids" {
  description = "Private subnet IDs"
  value       = module.vpc.private_subnet_ids
}

output "nat_gateway_id" {
  description = "The NAT Gateway ID"
  value       = module.nat_gateway.nat_gateway_id
}

output "public_instance_id" {
  description = "Public EC2 instance ID"
  value       = module.ec2.public_instance_id
}

output "public_instance_public_ip" {
  description = "Public IP of the public EC2 instance"
  value       = module.ec2.public_instance_public_ip
}

output "private_instance_id" {
  description = "Private EC2 instance ID"
  value       = module.ec2.private_instance_id
}

output "private_instance_private_ip" {
  description = "Private IP of the private EC2 instance"
  value       = module.ec2.private_instance_private_ip
}
