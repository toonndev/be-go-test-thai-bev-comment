# Backend Requirements — IT 08-1 (Comment Feed)

## Tech Stack

- **Language:** Go (latest stable version, ≥ 1.23)
- **Framework:** Gin (`github.com/gin-gonic/gin`)
- **ORM:** GORM (`gorm.io/gorm`) + PostgreSQL driver (`gorm.io/driver/postgres`)
- **Architecture:** Hexagonal Architecture (Ports & Adapters)
- **Config:** Environment variables via `.env` (use `godotenv`)
- **Testing:** `testify` + `testify/mock`

---

## Project Structure

```
backend/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── domain/
│   │   └── comment/
│   │       ├── comment.go        # Entity
│   │       └── repository.go     # Port (interface)
│   ├── usecase/
│   │   └── comment/
│   │       ├── usecase.go
│   │       └── usecase_test.go
│   ├── adapter/
│   │   ├── handler/
│   │   │   └── comment_handler.go
│   │   └── repository/
│   │       └── postgres_comment_repo.go
│   └── infrastructure/
│       └── database/
│           └── postgres.go
├── .env.example
├── go.mod
└── go.sum
```

---

## Database Schema

### Table: `comments`

| Column       | Type         | Constraints               |
|--------------|--------------|---------------------------|
| id           | BIGSERIAL    | PRIMARY KEY               |
| author_name  | VARCHAR(255) | NOT NULL                  |
| content      | TEXT         | NOT NULL                  |
| created_at   | TIMESTAMPTZ  | NOT NULL, DEFAULT NOW()   |
| updated_at   | TIMESTAMPTZ  | NOT NULL                  |

> หมายเหตุ: ไม่มี soft delete (ไม่ต้องการ deleted_at)

---

## Domain Entity

```go
// internal/domain/comment/comment.go
type Comment struct {
    ID          uint      `json:"id"`
    AuthorName  string    `json:"author_name"`
    Content     string    `json:"content"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

---

## Repository Port (Interface)

```go
// internal/domain/comment/repository.go
type Repository interface {
    Create(ctx context.Context, comment *Comment) (*Comment, error)
    FindAll(ctx context.Context) ([]*Comment, error)
}
```

---

## API Endpoints

### Base URL: `/api/v1`

---

### 1. GET `/api/v1/comments`

ดึง comment ทั้งหมด เรียงจากใหม่ไปเก่า

**Response 200:**
```json
{
  "data": [
    {
      "id": 2,
      "author_name": "Blend 285",
      "content": "have a good day",
      "created_at": "2021-10-16T16:00:00Z",
      "updated_at": "2021-10-16T16:00:00Z"
    }
  ]
}
```

---

### 2. POST `/api/v1/comments`

สร้าง comment ใหม่

**Request Body:**
```json
{
  "author_name": "Blend 285",
  "content": "have a good day"
}
```

**Validation:**
- `author_name`: required, ไม่ใช่ empty string
- `content`: required, ไม่ใช่ empty string

**Response 201:**
```json
{
  "data": {
    "id": 2,
    "author_name": "Blend 285",
    "content": "have a good day",
    "created_at": "2021-10-16T16:00:00Z",
    "updated_at": "2021-10-16T16:00:00Z"
  }
}
```

**Response 400 (validation error):**
```json
{
  "error": "author_name and content are required"
}
```

---

## CORS

เปิด CORS สำหรับทุก origin (เพื่อรองรับ Vue frontend ที่รันบน port อื่น)

```go
// Allow: *
// Methods: GET, POST, OPTIONS
// Headers: Content-Type
```

---

## Environment Variables (`.env.example`)

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=it08_1
SERVER_PORT=8080
```

---

## Unit Tests

เขียน unit test สำหรับ **usecase layer** โดยใช้ mock repository

### ไฟล์: `internal/usecase/comment/usecase_test.go`

**Test Cases:**

| Test Name | Scenario | Expected |
|---|---|---|
| `TestCreateComment_Success` | input ถูกต้อง | return comment, no error |
| `TestCreateComment_EmptyContent` | content เป็น empty string | return error |
| `TestCreateComment_EmptyAuthorName` | author_name เป็น empty string | return error |
| `TestGetAllComments_Success` | มี comment ใน DB | return slice of comments |
| `TestGetAllComments_Empty` | ไม่มี comment ใน DB | return empty slice, no error |

---

## Auto Migration

เมื่อ server start ให้ GORM AutoMigrate สร้าง table `comments` อัตโนมัติ

---

## Error Handling

- Bind/Validation error → HTTP 400
- Database error → HTTP 500
- Not Found → HTTP 404
- ทุก error response ใช้ format: `{"error": "message"}`

---

## Run

```bash
go mod tidy
go run cmd/server/main.go
```