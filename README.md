# 🚀 Realtime Social Sentiment Pipeline

> Hệ thống phân tích cảm xúc mạng xã hội realtime với Golang, Kafka, Spark

[![CI](https://github.com/luongtien872003/realtime-social-sentiment-pipline/actions/workflows/ci.yml/badge.svg)](https://github.com/luongtien872003/realtime-social-sentiment-pipline/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Kafka](https://img.shields.io/badge/Kafka-3.x-231F20?style=flat&logo=apache-kafka)](https://kafka.apache.org)

---

## 📋 Kiến Trúc

```
┌────────────────────┐     ┌─────────────┐     ┌─────────────┐
│  Golang Generator  │ ──► │   KAFKA     │ ──► │   REDIS     │
│  10k posts/min     │     │   Topic     │     │   Cache     │
│  Goroutines        │     │             │     │             │
└────────────────────┘     └──────┬──────┘     └──────┬──────┘
                                  │                   │
                                  ▼                   │
                           ┌─────────────┐            │
                           │ POSTGRESQL  │◄───────────┘
                           │   Storage   │
                           └──────┬──────┘
                                  │
                                  ▼
                           ┌─────────────┐     ┌─────────────┐
                           │   SPARK     │ ──► │   WEB UI    │
                           │  Streaming  │     │   Charts    │
                           └─────────────┘     └─────────────┘
```

---

## 🛠 Quick Start

### Prerequisites
- Go 1.21+
- Docker + Docker Compose
- Python 3.10+ (cho Spark)

### 1. Clone & Setup
```bash
git clone https://github.com/luongtien872003/realtime-social-sentiment-pipline.git
cd realtime-social-sentiment-pipline
go mod download
```

### 2. Start Infrastructure
```bash
docker-compose up -d
# Đợi 30s cho Kafka ready
```

### 3. Run Pipeline (3 terminals)
```bash
# Terminal 1: Consumer
go run cmd/consumer/main.go

# Terminal 2: Generator (10k posts/min streaming)
go run cmd/generator/main.go

# Terminal 3: API + Dashboard
go run cmd/api/main.go
```

### 4. View Dashboard
Open: **http://localhost:8888**

---

## 📊 Features

- ✅ **Streaming Generator**: 10k posts/min liên tục với goroutines
- ✅ **Kafka Message Queue**: High-throughput message processing
- ✅ **Redis Cache**: Realtime stats và recent posts
- ✅ **PostgreSQL Storage**: Batch insert với indexes
- ✅ **Spark Streaming**: Phân tích mỗi 30 giây
- ✅ **Web Dashboard**: Charts realtime với Chart.js

---

## 🔗 URLs

| Service | URL |
|---------|-----|
| Dashboard | http://localhost:8888 |
| Kafka UI | http://localhost:8080 |
| API Health | http://localhost:8888/api/health |

---

## 📁 Project Structure

```
├── cmd/
│   ├── generator/   # Streaming data generator
│   ├── consumer/    # Kafka to DB consumer
│   └── api/         # REST API server
├── internal/
│   ├── generator/   # Fake data logic
│   ├── kafka/       # Producer & Consumer
│   ├── redis/       # Cache layer
│   └── database/    # PostgreSQL
├── spark/jobs/      # Spark streaming analytics
├── web/             # Dashboard UI
├── deploy/docker/   # Dockerfiles
└── .github/workflows/ # CI/CD
```

---

## 🤝 Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for branch strategy and workflow.

---

## 📄 License

MIT License
