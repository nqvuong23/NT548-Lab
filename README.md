# NT548.Q21 - Lab

**Nhóm:** 10

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

