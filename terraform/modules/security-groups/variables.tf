variable "project_name" {
  type = string
}

variable "vpc_id" {
  description = "The VPC ID"
  type        = string
}

variable "tags" {
  type    = map(string)
  default = {}
}
