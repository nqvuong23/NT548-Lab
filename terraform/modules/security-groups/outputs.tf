output "public_instance_sg_id" {
  description = "The ID of the security group for the public EC2 instance"
  value       = aws_security_group.public_instance.id
}

output "private_instance_sg_id" {
  description = "The ID of the security group for the private EC2 instance"
  value       = aws_security_group.private_instance.id
}
