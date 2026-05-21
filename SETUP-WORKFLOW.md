# Huong dan Setup GitHub Actions - Cho Vuong

## Buoc 1: Sua file workflow

File: .github/workflows/terraform.yml

Sua cac dong sau:

Dong 10: AWS_REGION: "ap-southeast-1"
Dong 11: TF_STATE_BUCKET: "nqvuong23-terraform-project"

Trong step "Create terraform.tfvars":
project_name = "nt548-lab"
aws_region = "ap-southeast-1"

Trong step "Create Backend Config":
region = "ap-southeast-1"

(Lap lai cho ca job terraform-apply)

## Buoc 2: Them GitHub Secrets

1. AWS Academy -> Start Lab -> doi den XANH
2. AWS Details -> Show -> copy 3 gia tri
3. Repo Settings -> Secrets -> New secret

Them 3 secrets:

- AWS_ACCESS_KEY_ID
- AWS_SECRET_ACCESS_KEY
- AWS_SESSION_TOKEN

## Buoc 3: Kiem tra S3 Bucket

aws s3 ls | grep nqvuong23-terraform-project

Neu chua co:
aws s3 mb s3://nqvuong23-terraform-project --region ap-southeast-1

## Buoc 4: Chay workflow

Actions tab -> Terraform CI/CD with Checkov -> Run workflow
