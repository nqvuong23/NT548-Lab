locals {
  tags = {
    Project     = var.project_name
    Environment = var.environment
    ManagedBy   = "Terraform"
  }
}

# ------- Module: VPC -------
module "vpc" {
  source = "./modules/vpc"

  project_name         = var.project_name
  vpc_cidr             = var.vpc_cidr
  public_subnet_cidrs  = var.public_subnet_cidrs
  private_subnet_cidrs = var.private_subnet_cidrs
  availability_zones   = var.availability_zones
  tags                 = local.tags
}

# ------- Module: NAT Gateway -------
module "nat_gateway" {
  source = "./modules/nat-gateway"

  project_name     = var.project_name
  public_subnet_id = module.vpc.public_subnet_ids[0]
  tags             = local.tags
}

# ------- Module: Route Tables -------
module "route_tables" {
  source = "./modules/route-tables"

  project_name        = var.project_name
  vpc_id              = module.vpc.vpc_id
  public_subnet_ids   = module.vpc.public_subnet_ids
  private_subnet_ids  = module.vpc.private_subnet_ids
  internet_gateway_id = module.vpc.internet_gateway_id
  nat_gateway_id      = module.nat_gateway.nat_gateway_id
  tags                = local.tags
}

# ------- Module: Security Groups -------
module "security_groups" {
  source = "./modules/security-groups"

  project_name = var.project_name
  vpc_id       = module.vpc.vpc_id
  tags         = local.tags
}

# ------- Module: EC2 -------
module "ec2" {
  source = "./modules/ec2"

  project_name                   = var.project_name
  ec2_ami_id                     = var.ec2_ami_id
  ec2_instance_type              = var.ec2_instance_type
  ec2_volume_type                = var.ec2_volume_type
  ec2_volume_size                = var.ec2_volume_size
  ec2_enable_detailed_monitoring = var.ec2_enable_detailed_monitoring

  public_subnet_id  = module.vpc.public_subnet_ids[0]
  private_subnet_id = module.vpc.private_subnet_ids[0]

  public_sg_ids  = [module.security_groups.public_instance_sg_id]
  private_sg_ids = [module.security_groups.private_instance_sg_id]

  tags = local.tags
}
