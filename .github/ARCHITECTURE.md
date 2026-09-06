# Zizone - Backend Architecture & Technical Specs

Backend system for the **Zizone** platform — a Chinese vocabulary learning application providing REST APIs for both the user mobile app and the admin panel.

---

## 🛠️ Tech Stack

- **Language:** Go 1.25
- **Framework:** Gin v1.12
- **ORM:** GORM v1.31 + PostgreSQL
- **Cache:** Redis (go-redis v9)
- **Auth:** JWT (golang-jwt v5) + Google OAuth 2.0
- **Upload:** Cloudflare R2 (aws-sdk-go-v2) + WebP processing (`deepteams/webp`, `disintegration/imaging`)
- **Mail:** SMTP (`net/smtp`, worker pool)
- **Password:** bcrypt (`golang.org/x/crypto`)
- **UUID:** `google/uuid`

---

## 📁 Project Structure

```text
.
├── main.go                     # Entry point, initializes DB/Redis/R2/Mail
├── config/                     # Connection initializations
│   ├── database.go             # PostgreSQL + GORM + AutoMigrate (User, Word)
│   ├── redis.go                # Redis client (global Redis + Ctx)
│   ├── cloudflareR2.go         # Cloudflare R2 client (R2Client, GetR2BucketName, GetR2PublicURL)
│   └── googleAuth.go           # Google OAuth (authorization_code → userinfo)
├── constant/                   # Enum constants
│   ├── app.go                  # NoAvatar (default avatar)
│   ├── color.go                # Terminal color codes
│   ├── provider.go             # ProviderAdminCreate, ProviderEmail, ProviderGoogle
│   └── role.go                 # RoleUser, RoleAdmin
├── internal/
│   ├── middleware/
│   │   ├── middleware.go        # ParseJWT (global, optional) + RequireAuth (gated) + RequireRole
│   │   ├── cors.go              # CORS — uses ALLOWED_ORIGINS from env
│   │   └── ...
│   ├── model/
│   │   ├── user.go              # User model (UUID PK, soft delete via gorm.DeletedAt)
│   │   ├── word.go              # Word model + Character + Example (JSONB)
│   │   └── ...
│   ├── modules/                 # Feature modules (handler → service → repository, SAME PACKAGE)
│   │   ├── auth/                # handler.go, service.go, repository.go, routes.go
│   │   ├── word/                # handler.go, service.go, repository.go, routes.go, dto.go
│   │   └── ...
│   └── router/
│       └── router.go            # Route registration & middleware chain
├── pkg/
│   ├── log/                     # Colorized console log (LogSuccess/Info/Error/Warn)
│   ├── response/                # JSON response helpers (Success/Fail + variants)
│   └── ...
├── tdo/                         # DTO (meta.go, profile.go)
├── template/                    # HTML email templates
│   ├── verifyEmail.html
│   ├── sendPassword.html
│   └── ...
└── utils/                       # Utility functions
    ├── generatePass.go          # GeneratePassword + SendPassword (for admin account creation)
    ├── hash.go                  # HashPassword / CheckPasswordHash (bcrypt)
    ├── jwt.go                   # GenerateAccessToken / GenerateResetPasswordToken / ParseAccessToken
    ├── mail.go                  # MailService (worker pool with 5 goroutines, SMTP Gmail)
    ├── otp.go                   # SendOTP / VerifyOTP (stores OTP in Redis, rate limit, max fails)
    ├── uploadImage.go           # UploadR2Image / UploadMultipleR2Images / UploadR2Video (webp, resize)
    └── ...