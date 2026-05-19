variable "project_name" {
  type = string
}

variable "vpc_id" {
  description = "The VPC ID"
  type        = string
}

variable "public_subnet_ids" {
  description = "List of public subnet IDs to associate with the public route table"
  type        = list(string)
}

variable "private_subnet_ids" {
  description = "List of private subnet IDs to associate with the private route table"
  type        = list(string)
}

variable "internet_gateway_id" {
  description = "The Internet Gateway ID for public route table"
  type        = string
}

variable "nat_gateway_id" {
  description = "The NAT Gateway ID for private route table"
  type        = string
}

variable "tags" {
  type    = map(string)
  default = {}
}
