# ------- Security Group for Public Instance -------
resource "aws_security_group" "public_instance" {
  name        = "${var.project_name}-public-instance-sg"
  description = "Allow inbound SSH from the internet"
  vpc_id      = var.vpc_id

  tags = merge(var.tags, {
    Name = "${var.project_name}-public-instance-sg"
  })
}

# Allow inbound SSH (port 22) from anywhere
resource "aws_vpc_security_group_ingress_rule" "public_ssh" {
  security_group_id = aws_security_group.public_instance.id
  description       = "Allow SSH from internet"
  cidr_ipv4         = "0.0.0.0/0"
  from_port         = 22
  to_port           = 22
  ip_protocol       = "tcp"
}

# Allow all outbound traffic
resource "aws_vpc_security_group_egress_rule" "public_all_out" {
  security_group_id = aws_security_group.public_instance.id
  description       = "Allow all outbound traffic"
  cidr_ipv4         = "0.0.0.0/0"
  ip_protocol       = "-1"
}


# ------- Security Group for Private Instance -------
resource "aws_security_group" "private_instance" {
  name        = "${var.project_name}-private-instance-sg"
  description = "Allow inbound SSH only from the public EC2 instance"
  vpc_id      = var.vpc_id

  tags = merge(var.tags, {
    Name = "${var.project_name}-private-instance-sg"
  })
}

# Allow inbound SSH (port 22) only from the public instance security group
resource "aws_vpc_security_group_ingress_rule" "private_ssh_from_public" {
  security_group_id            = aws_security_group.private_instance.id
  description                  = "Allow SSH from public instance"
  referenced_security_group_id = aws_security_group.public_instance.id
  from_port                    = 22
  to_port                      = 22
  ip_protocol                  = "tcp"
}

# Allow all outbound traffic
resource "aws_vpc_security_group_egress_rule" "private_all_out" {
  security_group_id = aws_security_group.private_instance.id
  description       = "Allow all outbound traffic"
  cidr_ipv4         = "0.0.0.0/0"
  ip_protocol       = "-1"
}
