# NT548.Q21 - Lab

**Nhóm:** 10

# Hướng dẫn triển khai và kiểm thử hạ tầng AWS bằng Terraform

## Các test case kiểm tra dịch vụ

### VPC
1. Kiểm tra VPC được tạo thành công với đúng CIDR block đã khai báo.
2. Kiểm tra DNS hostnames và DNS support đều được bật.
3. Kiểm tra số lượng Public Subnet khớp với danh sách CIDR đầu vào.
4. Kiểm tra số lượng Private Subnet khớp với danh sách CIDR đầu vào.
5. Kiểm tra Public Subnet có `map_public_ip_on_launch = true`.
6. Kiểm tra Private Subnet có `map_public_ip_on_launch = false`.
7. Kiểm tra Internet Gateway được tạo và gắn đúng vào VPC.
8. Kiểm tra `length(availability_zones) >= length(public_subnet_cidrs)` để tránh index out of range khi tạo subnet.

### Route Tables
1. Kiểm tra Public Route Table có route `0.0.0.0/0` trỏ đến Internet Gateway.
2. Kiểm tra Private Route Table có route `0.0.0.0/0` trỏ đến NAT Gateway.
3. Kiểm tra tất cả Public Subnet được associate với Public Route Table.
4. Kiểm tra tất cả Private Subnet được associate với Private Route Table.

### NAT Gateway
1. Kiểm tra Elastic IP được tạo thành công trong domain `vpc`.
2. Kiểm tra NAT Gateway được đặt trong Public Subnet chứ không phải Private Subnet.
3. Kiểm tra NAT Gateway được gắn đúng Elastic IP (`allocation_id` khớp với EIP).

### EC2
1. Kiểm tra Public Instance được tạo trong đúng Public Subnet.
2. Kiểm tra subnet của Public Instance có `map_public_ip_on_launch = true`.
3. Kiểm tra Private Instance được tạo trong đúng Private Subnet.
4. Kiểm tra Private Instance không có public IP (không associate public IP).
5. Kiểm tra cả hai instance có root volume được mã hóa (`encrypted = true`).
6. Kiểm tra cả hai instance dùng IMDSv2 (`http_tokens = "required"`).

### Security Groups
1. Kiểm tra Public SG có ingress rule cho phép TCP port 22 từ `0.0.0.0/0`.
2. Kiểm tra Public SG có egress rule cho phép toàn bộ traffic ra ngoài.
3. Kiểm tra Private SG có ingress rule SSH (port 22) chỉ từ Public SG (dùng `referenced_security_group_id`, không phải CIDR).
4. Kiểm tra Private SG có egress rule cho phép toàn bộ traffic ra ngoài.
5. Kiểm tra Private SG không có bất kỳ ingress rule nào từ `0.0.0.0/0`.

## Hướng dẫn chạy terraform_test.go

### Bước 1: Cấu hình khóa xác thực AWS (AWS Credentials)
Để Terraform có thể giao tiếp với AWS, cần cung cấp Access Key và Secret Key.

1. Truy cập AWS Management Console.
2. Điều hướng đến IAM -> IAM users -> Chọn User hiện có. Nếu chưa có thì tạo IAM users mới, chọn Attach policies directly -> AdministratorAccess ở bước Permissions.
3. Chuyển sang tab Security credentials, tìm mục Access keys và nhấn Create access key (chọn loại Command Line Interface (CLI)).
4. Mở Terminal/PowerShell trên máy tính và chạy lệnh:

```cmd
aws configure
```

Nhập lần lượt các thông tin:

- AWS Access Key ID: (dán Access Key vừa tạo)
- AWS Secret Access Key: (dán Secret Key vừa tạo)
- Default region name: ap-southeast-1
- Default output format: json

### Bước 2: Build Docker Image cho môi trường test
Môi trường test yêu cầu Golang và Terraform. Để tránh xung đột phiên bản, chúng ta sẽ đóng gói tất cả vào Docker.

1. Mở Terminal (Command Prompt hoặc PowerShell).
2. Di chuyển vào thư mục terraform (nơi chứa file Dockerfile và mã nguồn).
3. Chạy lệnh sau để build image (chỉ cần chạy 1 lần duy nhất):

```cmd
docker build -t terratest-env .
```

### Bước 3: Chạy kịch bản tự động hóa (Automated Tests)
Khi image đã sẵn sàng, chạy lệnh dưới đây để kích hoạt chu trình: Khởi tạo hạ tầng -> Kiểm tra 26 tiêu chí bảo mật & kiến trúc -> Dọn dẹp tài nguyên.

Lưu ý: Lệnh dưới đây dành cho Command Prompt (CMD) trên Windows.

```cmd
docker run --rm -it ^
   -v "%cd%":/app ^
   -w /app/test ^
   -v "%USERPROFILE%\.aws":/root/.aws:ro ^
   terratest-env ^
   sh -c "rm -rf /app/.terraform && go mod tidy && go test -v -timeout 30m"
```
Kết quả sẽ là:

```plaintext
TestAll 2026-05-30T05:21:16Z logger.go:66: Destroy complete! Resources: 22 destroyed.
TestAll 2026-05-30T05:21:16Z logger.go:66:
--- PASS: TestAll (290.29s)
   --- PASS: TestAll/TC01_VPC_CIDR (0.57s)
   --- PASS: TestAll/TC02_VPC_DNS (0.23s)
   --- PASS: TestAll/TC03_PublicSubnet_Count (0.00s)
   --- PASS: TestAll/TC04_PrivateSubnet_Count (0.00s)
   --- PASS: TestAll/TC05_PublicSubnet_MapPublicIP (0.21s)
   --- PASS: TestAll/TC06_PrivateSubnet_NoMapPublicIP (0.14s)
   --- PASS: TestAll/TC07_InternetGateway_AttachedToVPC (0.07s)
   --- PASS: TestAll/TC08_AZ_Count_GTE_Subnet_Count (0.00s)
   --- PASS: TestAll/TC09_PublicRouteTable_HasIGWRoute (0.09s)
   --- PASS: TestAll/TC10_PrivateRouteTable_HasNATRoute (0.10s)
   --- PASS: TestAll/TC11_PublicSubnets_AssociatedWith_PublicRT (0.08s)
   --- PASS: TestAll/TC12_PrivateSubnets_AssociatedWith_PrivateRT (0.13s)
   --- PASS: TestAll/TC13_EIP_Domain_VPC (0.15s)
   --- PASS: TestAll/TC14_NATGateway_InPublicSubnet (0.06s)
   --- PASS: TestAll/TC15_NATGateway_HasCorrectEIP (0.06s)
   --- PASS: TestAll/TC16_PublicInstance_InPublicSubnet (0.15s)
   --- PASS: TestAll/TC17_PublicSubnet_MapPublicIP_ForPublicInstance (0.21s)
   --- PASS: TestAll/TC18_PrivateInstance_InPrivateSubnet (0.15s)
   --- PASS: TestAll/TC19_PrivateInstance_NoPublicIP (0.09s)
   --- PASS: TestAll/TC20_BothInstances_RootVolume_Encrypted (0.41s)
   --- PASS: TestAll/TC21_BothInstances_IMDSv2_Required (0.20s)
   --- PASS: TestAll/TC22_PublicSG_SSH_From_Anywhere (0.06s)
   --- PASS: TestAll/TC23_PublicSG_AllowAllEgress (0.06s)
   --- PASS: TestAll/TC24_PrivateSG_SSH_OnlyFrom_PublicSG (0.07s)
   --- PASS: TestAll/TC25_PrivateSG_AllowAllEgress (0.06s)
   --- PASS: TestAll/TC26_PrivateSG_NoIngress_FromAnywhere (0.07s)
PASS
ok      github.com/nt548-lab/terraform-test     290.943s
```


# Hướng dẫn triển khai và kiểm thử hạ tầng AWS bằng CloudFormation

Phần này hướng dẫn chi tiết cách sử dụng các template CloudFormation (`vpc.yaml`, `routing.yaml`, `ec2.yaml`) để tự động khởi tạo các dịch vụ VPC, Route Tables, NAT Gateway, EC2, Security Groups và cách chạy 4 test case để kiểm tra kết quả.

---

## Phần 1: Cấu hình và Triển khai trên AWS

**Điều kiện tiên quyết:** Cần tạo một bộ khóa (Key Pair) trước khi chạy mã nguồn.

1. Truy cập vào giao diện AWS Console, tìm dịch vụ **EC2**.
2. Ở menu bên trái, chọn **Key pairs** -> Nhấn nút **Create key pair**.
3. Điền thông tin cấu hình như sau:
   * **Name**: `lab-key`
   * **Key pair type**: `RSA`
   * **Private key file format**: `.pem`
4. Nhấn **Create key pair**. File `lab-key.pem` sẽ được tải về máy.

### Bước 1: Triển khai Module mạng (VPC)
1. Chuyển sang dịch vụ **CloudFormation** trên AWS Console.
2. Nhấn nút **Create stack** -> Chọn **With new resources (standard)**.
3. Tại phần **Prerequisite - Prepare template**, chọn **Choose an existing template**.
4. Tại phần **Specify template**, chọn **Upload a template file** -> Nhấn **Choose file** và chọn file `vpc.yaml`. Nhấn **Next**.
5. Tại trang **Specify stack details**:
   * **Stack name**: Nhập `Stack-VPC`.
   * Các thông số trong phần **Parameters** giữ nguyên mặc định.
6. Nhấn **Next** cho đến trang cuối cùng và nhấn **Submit**. 
7. Chờ cho đến khi Status của stack chuyển sang `CREATE_COMPLETE`.

### Bước 2: Triển khai Module định tuyến (Routing & NAT Gateway)
1. Quay lại trang chủ CloudFormation, nhấn **Create stack** -> Chọn **With new resources (standard)**.
2. Chọn **Upload a template file** -> Chọn file `routing.yaml` -> Nhấn **Next**.
3. Tại trang **Specify stack details**:
   * **Stack name**: Nhập `Stack-Routing`.
4. Nhấn **Next** qua các bước và nhấn **Submit**. 
5. Chờ vài phút để khởi tạo NAT Gateway cho đến khi Status báo `CREATE_COMPLETE`.

### Bước 3: Triển khai Module máy chủ (EC2 & Security Groups)
1. Tiếp tục nhấn **Create stack** -> Chọn **With new resources (standard)**.
2. Chọn **Upload a template file** -> Chọn file `ec2.yaml` -> Nhấn **Next**.
3. Tại trang **Specify stack details**, cấu hình như sau:
   * **Stack name**: Nhập `Stack-EC2`.
   * **InstanceType**: Giữ nguyên `t2.micro`.
   * **KeyName**: Nhấn vào menu thả xuống và chọn `lab-key` đã tạo ở phần trước.
   * **LatestAmiId**: Giữ nguyên.
   * **YourIpAddress**: **[LƯU Ý QUAN TRỌNG]** Nhập địa chỉ Public IP của mạng đang sử dụng, bắt buộc kèm theo hậu tố `/32` (Ví dụ: `14.226.x.x/32` hoặc `115.76.x.x/32`). 
4. Nhấn **Next** qua các bước và nhấn **Submit**.
5. Đợi Status báo `CREATE_COMPLETE`. Sau đó, chọn vào tên `Stack-EC2`, chuyển sang tab **Outputs** và ghi lại 2 giá trị:
   * `PublicInstanceIP`: Địa chỉ Public IP của Public Instance.
   * `PrivateInstanceIP`: Địa chỉ Private IP của Private Instance.

---

## Phần 2: Các Test Case kiểm tra dịch vụ

Để chứng minh hệ thống hoạt động đúng yêu cầu, cần thực hiện 4 Test Case sau:
1. **Test Case 1:** Kiểm tra kết nối SSH vào Public EC2 từ IP đã được cấp phép.
2. **Test Case 2:** Kiểm tra tính bảo mật của Public Security Group (từ chối truy cập từ IP lạ).
3. **Test Case 3:** Kiểm tra kết nối SSH vào Private EC2 chỉ từ Public EC2 (Mô hình Bastion Host).
4. **Test Case 4:** Kiểm tra khả năng truy cập Internet của Private EC2 thông qua NAT Gateway.

---

## Phần 3: Hướng dẫn thực hiện Test Cases

### Thực hiện Test Case 1
*Mục đích: Đăng nhập vào Public EC2 bằng quyền truy cập hợp lệ.*

1. Mở Terminal (macOS/Linux) hoặc PowerShell (Windows) và di chuyển đến thư mục chứa file `lab-key.pem` (thường là thư mục Downloads).
2. Cấp quyền bảo mật cho file key:

   ```bash
   chmod 400 lab-key.pem
   ```

3. Chạy lệnh SSH (thay `<PublicInstanceIP>` bằng IP lấy ở tab Outputs):

   ```bash
   ssh -i "lab-key.pem" ec2-user@<PublicInstanceIP>
   ```

4. Khi được hỏi, gõ `yes` và nhấn Enter. Đăng nhập thành công khi màn hình hiển thị logo Amazon Linux và dấu nhắc lệnh `[ec2-user@ip-10-0-1-x ~]$`.

### Thực hiện Test Case 2
*Mục đích: Chứng minh Security Group chỉ cho phép đúng 1 IP cố định.*

1. Ngắt kết nối mạng hiện tại và kết nối máy tính sang một mạng Wifi khác (hoặc phát 4G từ điện thoại) để làm thay đổi Public IP.
2. Mở Terminal/PowerShell và chạy lại lệnh SSH tương tự như Test Case 1:

   ```bash
   ssh -i "lab-key.pem" ec2-user@<PublicInstanceIP>
   ```

3. Do khác IP, hệ thống mạng của AWS sẽ chặn kết nối. Kết quả trả về lỗi `Connection timed out`. Ví dụ kết quả có thể như sau:

   ```plaintext
   PS C:\Users\Lenovo\Downloads> ssh -i "lab-key.pem" ec2-user@54.90.106.48
   ssh: connect to host 54.90.106.48 port 22: Connection timed out
   PS C:\Users\Lenovo\Downloads>
   ```

### Thực hiện Test Case 3
*Mục đích: Chứng minh rằng Private EC2 hoàn toàn cách ly với Internet và chỉ có thể truy cập từ Public instance.*

1. Đảm bảo đang ở trong phiên SSH của Test Case 1 (đứng trong Public EC2).
2. Tạo file chìa khóa mới bên trong máy chủ này:

   ```bash
   nano lab-key.pem
   ```

3. Mở file `lab-key.pem` trên máy tính đã lưu file này trước đó bằng Notepad, copy toàn bộ nội dung và dán vào cửa sổ terminal đang mở. Nhấn Ctrl + O -> Enter để lưu, sau đó Ctrl + X để thoát.
4. Cấp quyền và kết nối vào Private IP (thay `<PrivateInstanceIP>` bằng IP nội bộ):

   ```bash
   chmod 400 lab-key.pem
   ssh -i "lab-key.pem" ec2-user@<PrivateInstanceIP>
   ```

5. Gõ `yes` khi được yêu cầu. Đăng nhập thành công khi dấu nhắc lệnh chuyển thành dải IP của Private Subnet: `[ec2-user@ip-10-0-2-x ~]$`.

### Thực hiện Test Case 4
*Mục đích: Kiểm tra lưu lượng ra Internet của Private Subnet được định tuyến đúng qua NAT Gateway.*

1. Đảm bảo đang đứng bên trong Private EC2 (đã thực hiện xong Test Case 3).
2. Chạy lệnh ping đến server của Google để kiểm tra kết nối mạng:

   ```bash
   ping google.com -c 4
   ```

3. Trả về kết quả thành công với số liệu `4 packets transmitted, 4 received, 0% packet loss`, chứng tỏ luồng dữ liệu đã đi qua NAT Gateway an toàn.

