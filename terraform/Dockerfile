# Sử dụng base image là Golang 1.21 trên nền Alpine linux 
FROM golang:1.21-alpine

# Cài đặt các gói cần thiết: git (để tải go modules), wget/unzip (để tải terraform), bash
RUN apk add --no-cache git bash wget unzip

# Khai báo phiên bản Terraform
ENV TERRAFORM_VERSION=1.10.5

# Tải và cài đặt Terraform trực tiếp từ HashiCorp
RUN wget https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip && \
    unzip terraform_${TERRAFORM_VERSION}_linux_amd64.zip -d /usr/local/bin/ && \
    rm terraform_${TERRAFORM_VERSION}_linux_amd64.zip

# Đặt thư mục làm việc mặc định trong container
WORKDIR /app