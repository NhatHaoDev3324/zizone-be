# goAuth - Toà diện hệ thống Xác thực và Mail Service (Gin Gonic + Go)

**goAuth** là một dự án backend được xây dựng bằng ngôn ngữ Go, tập trung vào việc cung cấp hệ thống xác thực (Authentication) mạnh mẽ, bảo mật và dịch vụ gửi email hiệu suất cao. Dự án này tuân thủ các kiến trúc phát triển phần mềm hiện đại, có tính mở rộng cao và hiệu năng tối ưu.

## 🚀 Tính năng nổi bật

### 1. Hệ thống Xác thực (Authentication)
*   **Đăng ký/Đăng nhập bằng Email**: Bảo mật mật khẩu bằng thuật toán Bcrypt.
*   **Xác thực 2 lớp (OTP)**: Quy trình xác thực email bắt buộc sau khi đăng ký để kích hoạt tài khoản.
*   **Google OAuth2**: Hỗ trợ đăng nhập nhanh thông qua tài khoản Google.
*   **Quản lý phiên (Session)**: Sử dụng JWT (JSON Web Token) để quản lý truy cập API.

### 2. Dịch vụ Mail chuyên nghiệp (Professional Mail Service)
*   **Worker Pool**: Sử dụng cơ chế hàng đợi (Channel) và nhiều công nhân (Workers) chạy ngầm để gửi email số lượng lớn mà không làm nghẽn API.
*   **Persistent SMTP**: Duy trì kết nối SMTP liên tục để tăng tốc độ gửi mail gấp 5-10 lần so với cách thông thường.
*   **Template linh hoạt**: Toàn bộ email được lưu dưới dạng file `template/email.html`, dễ dàng chỉnh sửa giao diện mà không cần can thiệp vào mã nguồn Go.

### 3. Bảo mật mã OTP theo tiêu chuẩn ngành
*   **Crypto Random**: Mã 6 số được tạo từ bộ sinh số ngẫu nhiên cấp độ mã hóa (Cryptographically Secure).
*   **Rate Limiting (Throttling)**: Giới hạn tần suất gửi OTP (1 phút/lần) để chống spam.
*   **Chống Brute-force**: Tự động vô hiệu hóa mã OTP sau 5 lần nhập sai để bảo vệ tài khoản khỏi các cuộc tấn công dò mã.
*   **Redis Caching**: Lưu trữ OTP trong Redis với thời hạn 5 phút, tự động dọn dẹp sau khi sử dụng hoặc hết hạn.

### 4. Cơ sở hạ tầng & Hiệu năng
*   **Cơ sở dữ liệu**: Sử dụng **PostgreSQL** cùng với **GORM** để quản lý dữ liệu người dùng.
*   **Bộ nhớ đệm**: Sử dụng **Redis** để tăng tốc độ truy vấn thông tin người dùng và quản lý OTP.
*   **Graceful Shutdown**: Đảm bảo server đóng các kết nối an toàn và hoàn tất các tác vụ đang chạy (như gửi email) trước khi tắt hoàn toàn.
*   **Custom Logging**: Hệ thống log phân loại màu sắc (Success, Info, Error, Warn) giúp theo dõi hệ thống dễ dàng.

## 🏗️ Cấu trúc dự án (Architecture)

```text
goAuth/
├── config/             # Cấu hình Database, Redis
├── constant/           # Định nghĩa các hằng số (Role, Provider, Color...)
├── factory/            # Các hàm hỗ trợ (Logging...)
├── internal/           # Mã nguồn logic cốt lõi
│   ├── modules/
│   │   └── auth/       # Module Xác thực (Handler, Service, Repository, Model)
│   └── router/         # Cấu hình định tuyến API
├── pkg/                # Các thư viện tiện ích dùng chung (Response, Validator)
├── template/           # Chứa các file HTML Email
├── utils/              # Các công cụ hỗ trợ (Mail, OTP, JWT, Hash)
├── .env                # Biến môi trường
├── main.go             # Điểm bắt đầu của ứng dụng
└── README.md
```

## 🛠️ Hướng dẫn cài đặt

### 1. Yêu cầu hệ thống
*   Go (v1.20+)
*   PostgreSQL
*   Redis

### 2. Cấu hình biến môi trường
Tạo file `.env` tại thư mục gốc và điền các thông số sau:

```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=goauth

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379

# JWT
ACCESS_JWT_SECRET=your_jwt_secret

# Gmail Service (Dùng App Password)
MAIL_USER=your_email@gmail.com
MAIL_PASS=your_app_password

# Google OAuth2
GOOGLE_CLIENT_ID=your_id
GOOGLE_CLIENT_SECRET=your_secret
GOOGLE_REDIRECT_URI=http://localhost:8080/api/v1/users/register-by-google
```

### 3. Chạy ứng dụng
```bash
go mod tidy
go run .
```

## 📝 Danh sách API chính

| Phương thức | Endpoint | Chức năng |
| :--- | :--- | :--- |
| **POST** | `/api/v1/users/register-by-email` | Đăng ký tài khoản mới & gửi OTP |
| **POST** | `/api/v1/users/verify-otp` | Xác thực OTP để kích hoạt tài khoản |
| **POST** | `/api/v1/users/login-by-email` | Đăng nhập bằng Email/Pass |
| **POST** | `/api/v1/users/register-by-google` | Đăng ký/Đăng nhập bằng Google |
| **GET** | `/api/v1/users/` | Lấy danh sách người dùng (Admin) |

---
*Dự án được phát triển bởi **NhatHaoDev3324***
