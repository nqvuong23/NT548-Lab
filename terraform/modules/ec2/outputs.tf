output "public_instance_id" {
  description = "The ID of the public EC2 instance"
  value       = aws_instance.public.id
}

output "public_instance_public_ip" {
  description = "The public IP of the public EC2 instance"
  value       = aws_instance.public.public_ip
}

output "public_instance_private_ip" {
  description = "The private IP of the public EC2 instance"
  value       = aws_instance.public.private_ip
}

output "private_instance_id" {
  description = "The ID of the private EC2 instance"
  value       = aws_instance.private.id
}

output "private_instance_private_ip" {
  description = "The private IP of the private EC2 instance"
  value       = aws_instance.private.private_ip
}
