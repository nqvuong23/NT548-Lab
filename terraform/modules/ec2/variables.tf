variable "project_name" {
  type = string
}

variable "ec2_ami_id" {
  type    = string
}

variable "ec2_instance_type" {
  type = string
}

variable "public_subnet_id" {
  type = string
}

variable "private_subnet_id" {
  type = string
}

variable "public_sg_ids" {
  type = list(string)
}

variable "private_sg_ids" {
  type = list(string)
}

variable "ec2_volume_type" {
  type    = string
}

variable "ec2_volume_size" {
  type    = number
}

variable "ec2_enable_detailed_monitoring" {
  type    = bool
}

variable "tags" {
  type    = map(string)
  default = {}
}
