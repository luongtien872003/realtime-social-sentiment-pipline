# 📁 Social Insight - Project Structure

> Kiến trúc chuẩn cho dự án Social Insight

---

## 🏗 Kiến Trúc Tổng Quan

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          SOCIAL INSIGHT ARCHITECTURE                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│   ┌──────────────┐                                                          │
│   │   Sources    │  Facebook, Twitter, Reddit, TikTok, News...              │
│   └──────┬───────┘                                                          │
│          │                                                                   │
│          ▼                                                                   │
│   ┌──────────────────────────────────────────────────────────────┐          │
│   │                     CRAWLER SERVICE (Go)                      │          │
│   │   cmd/crawler/main.go                                         │          │
│   │   internal/crawler/                                           │          │
│   └──────────────────────────┬───────────────────────────────────┘          │
│                              │                                               │
│                              ▼                                               │
│   ┌──────────────────────────────────────────────────────────────┐          │
│   │                    KAFKA MESSAGE QUEUE                        │          │
│   │   Topics: raw_data, processed_data, alerts                    │          │
│   └──────────────────────────┬───────────────────────────────────┘          │
│                              │                                               │
│              ┌───────────────┴───────────────┐                              │
│              │                               │                               │
│              ▼                               ▼                               │
│   ┌─────────────────────┐     ┌─────────────────────────────────┐          │
│   │  SPARK PROCESSING   │     │      ML PIPELINE (Python)        │          │
│   │  spark/jobs/        │     │      ml/sentiment/               │          │
│   │  - ETL              │     │      ml/trend/                   │          │
│   │  - Aggregation      │     │      ml/anomaly/                 │          │
│   └──────────┬──────────┘     └──────────────┬──────────────────┘          │
│              │                               │                               │
│              └───────────────┬───────────────┘                              │
│                              │                                               │
│                              ▼                                               │
│   ┌──────────────────────────────────────────────────────────────┐          │
│   │                       DATA STORAGE                            │          │
│   │   ┌─────────────┐  ┌─────────────┐  ┌─────────────┐          │          │
│   │   │    S3       │  │ PostgreSQL  │  │   Redis     │          │          │
│   │   │  (Raw/ML)   │  │ (Metrics)   │  │  (Cache)    │          │          │
│   │   └─────────────┘  └─────────────┘  └─────────────┘          │          │
│   └──────────────────────────┬───────────────────────────────────┘          │
│                              │                                               │
│                              ▼                                               │
│   ┌──────────────────────────────────────────────────────────────┐          │
│   │                   API GATEWAY (Go)                            │          │
│   │   cmd/api/main.go                                             │          │
│   │   internal/api/                                               │          │
│   └──────────────────────────┬───────────────────────────────────┘          │
│                              │                                               │
│                              ▼                                               │
│   ┌──────────────────────────────────────────────────────────────┐          │
│   │                    DASHBOARD (React)                          │          │
│   │   web/                                                        │          │
│   └──────────────────────────────────────────────────────────────┘          │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 📂 Cấu Trúc Thư Mục

```
project/
├── 📂 cmd/                          # Entry points
│   ├── crawler/
│   │   └── main.go                  # Crawler service entry
│   ├── api/
│   │   └── main.go                  # API Gateway entry
│   └── worker/
│       └── main.go                  # Background workers
│
├── 📂 internal/                     # Private Go packages
│   ├── crawler/
│   │   ├── client.go                # HTTP client
│   │   ├── parser.go                # HTML parser
│   │   └── crawler.go               # Main crawler logic
│   │
│   ├── kafka/
│   │   ├── producer.go              # Kafka producer
│   │   └── consumer.go              # Kafka consumer
│   │
│   ├── api/
│   │   ├── handlers/                # API handlers
│   │   ├── middleware/              # Middlewares
│   │   └── routes.go                # Route definitions
│   │
│   └── models/
│       └── message.go               # Data structures
│
├── 📂 pkg/                          # Public Go packages
│   ├── config/
│   │   └── config.go                # Configuration loader
│   └── utils/
│       └── utils.go                 # Utility functions
│
├── 📂 ml/                           # Machine Learning (Python)
│   ├── sentiment/
│   │   ├── model.py                 # Sentiment model
│   │   ├── train.py                 # Training script
│   │   └── api.py                   # Serving API
│   │
│   ├── trend/
│   │   ├── model.py                 # Trend prediction
│   │   └── train.py
│   │
│   ├── anomaly/
│   │   ├── model.py                 # Anomaly detection
│   │   └── train.py
│   │
│   └── requirements.txt             # Python dependencies
│
├── 📂 spark/                        # Spark Jobs (Python)
│   ├── jobs/
│   │   ├── process_raw_data.py      # Raw data processor
│   │   ├── aggregate_metrics.py     # Metrics aggregation
│   │   └── sentiment_batch.py       # Batch sentiment
│   │
│   └── schemas/
│       └── data_schema.py           # Data schemas
│
├── 📂 deploy/                       # Deployment configs
│   ├── docker/
│   │   ├── Dockerfile.crawler       # Crawler Dockerfile
│   │   ├── Dockerfile.api           # API Dockerfile
│   │   └── Dockerfile.ml            # ML Dockerfile
│   │
│   ├── k8s/                         # Kubernetes manifests
│   │   ├── crawler-deployment.yaml
│   │   ├── api-deployment.yaml
│   │   └── ingress.yaml
│   │
│   └── terraform/                   # Infrastructure as Code
│       ├── main.tf
│       ├── variables.tf
│       └── outputs.tf
│
├── 📂 scripts/                      # Utility scripts
│   ├── setup.sh                     # Initial setup
│   ├── migrate.sh                   # DB migrations
│   └── deploy.sh                    # Deployment script
│
├── 📂 tests/                        # Test files
│   ├── crawler/
│   ├── kafka/
│   ├── spark/
│   ├── api/
│   └── ml/
│
├── 📂 web/                          # Frontend (optional)
│   ├── src/
│   ├── public/
│   └── package.json
│
├── 📂 docs/                         # Documentation
│   ├── api.md
│   ├── architecture.md
│   └── deployment.md
│
├── 📂 migrations/                   # Database migrations
│   └── 001_initial.sql
│
├── 📄 .github/
│   └── workflows/
│       └── ci.yml                   # GitHub Actions CI/CD
│
├── 📄 go.mod                        # Go modules
├── 📄 go.sum
├── 📄 Makefile                      # Build commands
├── 📄 docker-compose.yml            # Local development
├── 📄 docker-compose.kafka.yml      # Kafka setup
├── 📄 .env.example                  # Environment template
├── 📄 README.md
└── 📄 LEARNING_GUIDE.md
```

---

## 📄 Các File Cần Tạo

### Go Files

| Path | Mô tả |
|------|-------|
| `go.mod` | Module definition |
| `cmd/crawler/main.go` | Crawler entry point |
| `cmd/api/main.go` | API entry point |
| `internal/crawler/client.go` | HTTP client |
| `internal/crawler/parser.go` | HTML parser |
| `internal/kafka/producer.go` | Kafka producer |
| `internal/kafka/consumer.go` | Kafka consumer |
| `internal/api/handlers/sentiment.go` | Sentiment API handler |
| `internal/api/routes.go` | API routes |
| `pkg/config/config.go` | Configuration |

### Python Files

| Path | Mô tả |
|------|-------|
| `ml/requirements.txt` | Python dependencies |
| `ml/sentiment/model.py` | Sentiment model class |
| `ml/sentiment/train.py` | Training script |
| `ml/sentiment/api.py` | FastAPI serving |
| `spark/jobs/process_raw_data.py` | Spark ETL job |

### Config Files

| Path | Mô tả |
|------|-------|
| `docker-compose.yml` | Local dev environment |
| `docker-compose.kafka.yml` | Kafka + Zookeeper |
| `Makefile` | Build & run commands |
| `.github/workflows/ci.yml` | CI/CD pipeline |
| `.env.example` | Environment variables template |

### Docker Files

| Path | Mô tả |
|------|-------|
| `deploy/docker/Dockerfile.crawler` | Crawler image |
| `deploy/docker/Dockerfile.api` | API image |
| `deploy/docker/Dockerfile.ml` | ML service image |

---

## 🔧 Makefile Commands

```makefile
# Local Development
make setup          # Install dependencies
make dev            # Start dev environment
make test           # Run all tests

# Docker
make build          # Build all images
make up             # Start containers
make down           # Stop containers

# Database
make migrate-up     # Run migrations
make migrate-down   # Rollback migrations

# Services
make run-crawler    # Run crawler
make run-api        # Run API
make run-spark      # Run Spark job
```

---

## 🌍 Environment Variables

```bash
# .env.example

# Kafka
KAFKA_BROKERS=localhost:9092
KAFKA_TOPIC_RAW=raw_data
KAFKA_TOPIC_PROCESSED=processed_data

# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=social_insight
DB_USER=postgres
DB_PASSWORD=password

# AWS
AWS_REGION=ap-southeast-1
AWS_S3_BUCKET=social-insight-data

# ML
ML_SERVICE_URL=http://localhost:8001
ML_MODEL_PATH=/models/sentiment

# API
API_PORT=8080
API_DEBUG=true
```
