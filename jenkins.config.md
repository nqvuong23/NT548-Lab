# Hướng dẫn cài đặt và cấu hình Jenkins

## Cài đặt Jenkins

```
sudo apt update
sudo apt install fontconfig openjdk-21-jre
java -version
```

```
sudo wget -O /etc/apt/keyrings/jenkins-keyring.asc \
  https://pkg.jenkins.io/debian-stable/jenkins.io-2026.key

echo "deb [signed-by=/etc/apt/keyrings/jenkins-keyring.asc]" \
  https://pkg.jenkins.io/debian-stable binary/ | sudo tee \
  /etc/apt/sources.list.d/jenkins.list > /dev/null

sudo apt update
sudo apt install jenkins
sudo systemctl enable jenkins
sudo systemctl start jenkins
```

---

## Cài đặt Docker

```
sudo apt remove $(dpkg --get-selections docker.io docker-compose docker-compose-v2 docker-doc podman-docker containerd runc | cut -f1)
```

```
sudo apt update
sudo apt install ca-certificates curl
sudo install -m 0755 -d /etc/apt/keyrings
sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
sudo chmod a+r /etc/apt/keyrings/docker.asc
```

```
# Copy toàn bộ các dòng trong block code
sudo tee /etc/apt/sources.list.d/docker.sources <<EOF
Types: deb
URIs: https://download.docker.com/linux/ubuntu
Suites: $(. /etc/os-release && echo "${UBUNTU_CODENAME:-$VERSION_CODENAME}")
Components: stable
Architectures: $(dpkg --print-architecture)
Signed-By: /etc/apt/keyrings/docker.asc
EOF
```

```
sudo apt update
sudo apt install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
```

```
sudo systemctl enable docker
sudo systemctl start docker
sudo usermod -aG docker $USER
```

---

## Cài đặt Trivy

```
sudo apt-get install wget gnupg
wget -qO - https://aquasecurity.github.io/trivy-repo/deb/public.key | gpg --dearmor | sudo tee /usr/share/keyrings/trivy.gpg > /dev/null
echo "deb [signed-by=/usr/share/keyrings/trivy.gpg] https://aquasecurity.github.io/trivy-repo/deb generic main" | sudo tee -a /etc/apt/sources.list.d/trivy.list
sudo apt-get update
sudo apt-get install trivy
```

---

## Cấu hình Jenkins

### Bước 1: Đăng nhập

1. Mở trình duyệt web và truy cập thẳng vào: `http://<IP của EC2 hoặc localhost>:8080`
2. Đăng nhập bằng tài khoản admin: Username `admin`, Password `admin`.

---

### Bước 2: Cấu hình Jenkins sử dụng SonarQube

1. Tải Plugin của SonarQube: Vào Manage Jenkins (biểu tượng răng cưa) > Plugins > vào Tab "Available plugins" > Tìm và tải plugin "SonarQube Scanner"
2. Vào Manage Jenkins (biểu tượng răng cưa) > Tools > Tìm mục "SonarQube Scanner installations" > Nhấn `Add SonarQube Scanner` > Nhập Name = "sonar-scanner" và tich vào ô `Install automatically` > Apply và Save
3. Vào Manage Jenkins (biểu tượng răng cưa) > Credentials > Add Credentials. Điền các thông tin sau:
   - Kind: Chọn `Secret text`.
   - Secret: Dán token của SonarQube.
   - ID: Nhập `sonarqube-token`.
   - Bấm Create
4. Vào Manage Jenkins > System 
5. Tìm đến mục `SonarQube servers` > Nhấn Add SonarQube. Nhập các thông tin sau:
   - Name: Nhập `sonarqube-server`.
   - Server URL: Nhập `http://<IP của EC2 hoặc localhost>:9000`
   - Server authentication token: chọn Credentials `sonarqube-token` vừa tạo
6. Apply and Save

---

### Bước 3: Tạo Pipeline Job

1. Quay ra trang chủ Jenkins, bấm **New Item**.
2. Nhập tên Job: `NT548-Lab-Nhom10-Pipeline`.
3. Chọn loại **Pipeline** và bấm **OK**.
4. Cuộn xuống mục **Pipeline**, thiết lập như sau:
   - Definition: Chọn `Pipeline script from SCM`.
   - SCM: Chọn `Git`.
   - Repository URL: Dán link GitHub repo của nhóm vào (không phải HTTPS URL, dùng SSH).
   - Branch Specifier: `*/main`.
   - Script Path: Nhập `Jenkinsfile`.
5. Bấm **Save**.
