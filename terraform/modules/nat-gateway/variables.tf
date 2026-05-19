variable "project_name" {
  type = string
}

variable "public_subnet_id" {
  description = "The public subnet ID where the NAT Gateway will be placed"
  type        = string
}

variable "tags" {
  type    = map(string)
  default = {}
}
