# ------- EC2 Instance -------
resource "aws_instance" "public" {
  ami                    = var.ec2_ami_id
  instance_type          = var.ec2_instance_type
  subnet_id              = var.public_subnet_id
  vpc_security_group_ids = var.public_sg_ids

  root_block_device {
    volume_type           = var.ec2_volume_type
    volume_size           = var.ec2_volume_size
    delete_on_termination = true
    encrypted             = true

    tags = merge(var.tags, {
      Name = "${var.project_name}-data-processing-root-volume"
    })
  }

  metadata_options {
    http_endpoint               = "enabled"
    http_tokens                 = "required"
    http_put_response_hop_limit = 2
  }

  monitoring = var.ec2_enable_detailed_monitoring

  tags = merge(var.tags, {
    Name = "${var.project_name}-public-instance"
  })
}

# ------- EC2 Instance -------
resource "aws_instance" "private" {
  ami                    = var.ec2_ami_id
  instance_type          = var.ec2_instance_type
  subnet_id              = var.private_subnet_id
  vpc_security_group_ids = var.private_sg_ids

  root_block_device {
    volume_type           = var.ec2_volume_type
    volume_size           = var.ec2_volume_size
    delete_on_termination = true
    encrypted             = true

    tags = merge(var.tags, {
      Name = "${var.project_name}-data-processing-root-volume"
    })
  }

  metadata_options {
    http_endpoint               = "enabled"
    http_tokens                 = "required"
    http_put_response_hop_limit = 2
  }

  monitoring = var.ec2_enable_detailed_monitoring

  tags = merge(var.tags, {
    Name = "${var.project_name}-private-instance"
  })
}
