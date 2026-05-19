# ------- Elastic IP for NAT Gateway -------
resource "aws_eip" "nat" {
  domain = "vpc"

  tags = merge(var.tags, {
    Name = "${var.project_name}-nat-eip"
  })
}

# ------- NAT Gateway -------
resource "aws_nat_gateway" "main" {
  allocation_id = aws_eip.nat.id
  subnet_id     = var.public_subnet_id

  tags = merge(var.tags, {
    Name = "${var.project_name}-nat-gw"
  })

  depends_on = [aws_eip.nat]
}