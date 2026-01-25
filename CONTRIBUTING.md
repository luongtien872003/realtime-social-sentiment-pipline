# 🤝 Contributing Guide

Hướng dẫn đóng góp cho dự án Social Insight.

---

## 📋 Branch Strategy

```
main (production)
  │
  └── develop (staging/integration)
        │
        ├── feature/xyz  (tính năng mới)
        ├── bugfix/xyz   (sửa bug)
        └── hotfix/xyz   (sửa lỗi khẩn cấp production)
```

### Branch Rules

| Branch | Mục đích | Merge từ | Deploy đến |
|--------|----------|----------|------------|
| `main` | Production code | develop, hotfix | Production |
| `develop` | Integration branch | feature, bugfix | Staging |
| `feature/*` | Tính năng mới | - | - |
| `bugfix/*` | Sửa bug | - | - |
| `hotfix/*` | Sửa lỗi khẩn cấp | - | - |

---

## 🔄 Workflow

### 1. Tạo Feature Branch

```bash
# Từ develop branch
git checkout develop
git pull origin develop
git checkout -b feature/my-feature
```

### 2. Phát triển

```bash
# Code, test locally
go test ./...

# Commit với message rõ ràng
git add .
git commit -m "feat: add new sentiment analysis model"
```

### 3. Push và Tạo PR

```bash
git push origin feature/my-feature
# Tạo Pull Request vào develop trên GitHub
```

### 4. Code Review

- Ít nhất 1 reviewer approve
- CI pipeline phải pass
- No conflicts với develop

### 5. Merge

- Squash and merge vào develop
- Delete feature branch sau khi merge

---

## 📝 Commit Convention

Format: `<type>: <description>`

| Type | Mô tả |
|------|-------|
| `feat` | Tính năng mới |
| `fix` | Sửa bug |
| `docs` | Cập nhật documentation |
| `style` | Format code (không thay đổi logic) |
| `refactor` | Refactor code |
| `test` | Thêm/sửa tests |
| `chore` | Tasks khác (CI, deps, etc.) |

Ví dụ:
```
feat: add kafka producer with batch support
fix: resolve race condition in consumer
docs: update README with streaming guide
```

---

## 🧪 Testing Locally

```bash
# Chạy tất cả tests
go test ./...

# Chạy với coverage
go test -cover ./...

# Chạy race detector
go test -race ./...
```

---

## 🐳 Running Locally

```bash
# 1. Khởi động infrastructure
docker-compose up -d

# 2. Chạy services
go run cmd/consumer/main.go  # Terminal 1
go run cmd/generator/main.go # Terminal 2
go run cmd/api/main.go       # Terminal 3

# 3. Xem dashboard
open http://localhost:8888
```

---

## 📋 Checklist Trước Khi Tạo PR

- [ ] Code chạy được locally
- [ ] Đã viết/update tests
- [ ] `go test ./...` pass
- [ ] `golangci-lint run` không có errors
- [ ] Commit messages theo convention
- [ ] Update documentation nếu cần
