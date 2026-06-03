# Hướng dẫn cài đặt và cấu hình SonarQube

## Cài đặt SonarQube

```
echo "vm.max_map_count=524288" | sudo tee -a /etc/sysctl.d/99-sonarqube.conf
sudo sysctl --system
```

```
# Tạo file docker-compose.yaml
touch docker-compose.yaml

# Mở file
nano docker-compose.yaml

# Dán nội dung file docker-compose-sonarqube.yaml trong thư mục dự án vào file docker-compose.yaml đang mở --> lưu và thoát

# chạy docker compose
docker compose up -d
```

---

## Cấu hình SonarQube

### Bước 1: Đăng nhập

1. Truy cập URL: `http://<IP của EC2 hoặc localhost>:9000`
2. Đăng nhập bằng tài khoản demo: Username `admin`, Password `admin` --> sau đó sonarqube sẽ yêu cầu đổi mật khẩu.

---

### Bước 2: Tạo Project

1. Chọn tab **Projects** trên thanh menu trên cùng.
2. Bấm **Create a local project**.
3. Tại ô *Project display name* và *Project key*, nhập chính xác tên: `NT548-Lab-Nhom10`.
4. Chọn **Follows the instance's default**
5. Bấm **Creat project**.

---

### Bước 3: Tạo Token để cấp quyền cho Jenkins

1. Chọn dự án **DevSecOps_Nhom10** vừa tạo.
2. Tại **Project onboarding** chọn **Locally**.
3. Ở phần *Generate Tokens*:
   - Name: Nhập tên tùy ý (VD: `jenkins-token`).
   - Type: Chọn `Global Analysis Token` (hoặc Project Token).
   - Expires in: Chọn `No expiration`.
4. Bấm **Generate**. 
5. **QUAN TRỌNG:** Copy ngay đoạn mã Token vừa hiện ra và lưu tạm ra Notepad (vì nó chỉ hiện 1 lần duy nhất).

---

### Bước 4: Thiết lập Webhook (Quality Gate)

1. Quay lại trang chủ SonarQube, bấm vào dự án `NT548-Lab-Nhom10` vừa tạo.
2. Chọn menu **Project Settings** > **Webhooks**.
3. Bấm **Create** và điền thông tin:
   - Name: `Jenkins Webhook`
   - URL: `http://<IP của EC2 hoặc localhost>:8080/sonarqube-webhook/`
4. Bấm **Create** để lưu lại.
