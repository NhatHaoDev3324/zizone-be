# Go Clean Architecture Template (Gin, GORM, PostgreSQL)

Đây là một bản mẫu (template) dự án Golang được xây dựng với mục tiêu: **Dễ mở rộng (Scalable), Dễ bảo trì (Maintainable) và Tái sử dụng cao (Reusable)**. Dự án sử dụng cấu trúc Modular kết hợp với Clean Architecture rút gọn.

---

## 🚀 Công nghệ sử dụng
- **Framework:** [Gin Gonic](https://gin-gonic.com/) (High-performance HTTP web framework)
- **ORM:** [GORM](https://gorm.io/) (Object Relational Mapping cho Golang)
- **Database:** PostgreSQL
- **Config:** [Godotenv](https://github.com/joho/godotenv)
- **Middleware:** CORS tùy chỉnh, Recovery, Logger.
- **CI/CD:** GitHub Actions (hỗ trợ SSH và Docker Hub).

---

## 📂 Cấu trúc thư mục (Directory Structure)

```text
├── main.go                 # Entry point của ứng dụng
├── config/                 # Cấu hình Database & Environment
├── internal/               # Chứa logic cốt lõi (không cho phép import từ ngoài)
│   ├── middleware/         # Các middleware (CORS, Auth, Logger...)
│   ├── modules/            # Chia theo từng tính năng (Module)
│   │   └── user/           # Ví dụ Module User
│   │       ├── handler/    # Tiếp nhận HTTP request, validate dữ liệu
│   │       ├── service/    # Logic nghiệp vụ (Business Logic)
│   │       ├── repository/ # Thao tác với Database (SQL queries)
│   │       ├── model/      # Định nghĩa struct (Schema)
│   │       └── routes.go   # Định nghĩa route riêng cho module
│   └── router/             # Khởi tạo router chính, gộp các module
├── pkg/                    # Các công cụ/tiện ích dùng chung (Utils)
├── .github/workflows/      # Cấu hình CI/CD (GitHub Actions)
├── Dockerfile              # Cấu hình Docker image
├── docker-compose.yml      # Orchestration cho Docker
└── .env                    # Biến môi trường (DB, Secret Key...)
```

---

## 🛠️ Cách cài đặt và sử dụng

### 1. Yêu cầu hệ thống
- Go version 1.20+
- PostgreSQL đang chạy.

### 2. Cài đặt
1. **Clone dự án:**
   ```bash
   git clone <link-repo>
   cd <folder-name>
   ```

2. **Cấu hình môi trường:**
   Tạo file `.env` từ mẫu:
   ```env
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=yourpassword
   DB_NAME=yourdatabase
   ALLOWED_ORIGINS=http://localhost:3000,http://127.0.0.1:3000
   ```

3. **Cài đặt thư viện:**
   ```bash
   go mod tidy
   ```

4. **Chạy dự án:**
   ```bash
   go run .
   ```

---

## 🔄 Cách hoạt động (Data Flow)

Dự án tuân thủ luồng dữ liệu 1 chiều để đảm bảo tính minh bạch:

1. **Client** gửi Request đến các Endpoint.
2. **Router** định tuyến request đến đúng **Handler**.
3. **Handler** kiểm tra dữ liệu đầu vào (Binding/Validation) và gọi **Service**.
4. **Service** xử lý logic nghiệp vụ (tính toán, kiểm tra điều kiện...) và gọi **Repository**.
5. **Repository** thực hiện truy vấn xuống **Database** thông qua GORM.
6. Kết quả được trả ngược lại theo đúng thứ tự để Handler phản hồi cho Client.

---

## 🚢 Triển khai (Deployment)

Dự án đi kèm sẵn file `.github/workflows/main.yml` hỗ trợ các luồng deploy tự động:
- **Deploy qua SSH:** Tự động pull code mới về VPS và khởi chạy qua Docker Compose.
- **Deploy qua Docker Hub:** Build image, push lên Docker Hub và lệnh cho server pull image mới về chạy.

---

## ⚖️ Ưu và Nhược điểm

### ✅ Ưu điểm:
- **Tính đóng gói cao:** Mỗi module (user, product,...) là một thực thể độc lập. Khi bạn muốn sửa User, bạn không cần quan tâm đến Product.
- **Dễ Unit Test:** Vì tách biệt Service và Repository dưới dạng Interface, bạn có thể dễ dàng viết Mock Test.
- **Dễ mở rộng:** Chỉ cần copy cấu trúc một module hiện có để tạo module mới trong vài phút.
- **Cấu hình linh hoạt:** Mọi thứ từ Database đến CORS đều được quản lý qua `.env`.

### ❌ Nhược điểm:
- **Boilerplate:** Đối với các dự án siêu nhỏ (vài API), cấu trúc này có vẻ hơi "cồng kềnh" do phải chia nhiều lớp (Handler-Service-Repo).
- **Độ dốc học tập:** Những người mới bắt đầu có thể thấy khó hiểu về lý do tại sao phải chia nhỏ file như vậy.

---

## 📝 Ghi chú
- Dự án này đã tích hợp sẵn **Auto-Migration** của GORM, giúp tự động tạo bảng dữ liệu khi bạn khai báo Model mới.
- Đừng bao giờ commit file `.env` lên Github (đã có trong `.gitignore`).
