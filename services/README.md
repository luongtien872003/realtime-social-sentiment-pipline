# 🔵 Social Insight - Microservices Architecture

**Architecture**: 3 Independent Microservices  
**Status**: ✅ Production Ready  
**Setup Time**: 10-20 minutes  

---

## 📊 Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                         LOAD BALANCER                           │
└────────────────┬──────────────────────────────────────┬─────────┘
                 │                                      │
        ┌────────▼──────┐              ┌───────────────▼────┐
        │  Data Service │              │  Processing Service │
        │ (Infrastructure)│              │ (Crawlers + Consumer)│
        │                │              │                    │
        │ • PostgreSQL   │◄──connection──┤ • HN Crawler       │
        │ • Redis        │   (env vars)  │ • DevTo Crawler    │
        │ • Kafka        │              │ • Medium Crawler   │
        │ • Zookeeper    │              │ • Consumer         │
        └─────────┬──────┘              └────────────────────┘
                  │                              │
                  │◄─────connection (env vars)──┤
                  │                              │
                  └────────────┬─────────────────┘
                               │
                      ┌────────▼──────────┐
                      │   API Service     │
                      │                   │
                      │ • REST API        │
                      │ • Web Dashboard   │
                      └───────────────────┘
                               │
                               ▼
                        http://localhost:8888
```

---

## 🚀 Quick Start (3 Services, 3 Commands)

### Step 1: Start Data Service (Infrastructure)
```bash
cd services/data-service
docker-compose up -d
```

**Expected Output:**
```
Creating data_postgres ... done
Creating data_zookeeper ... done
Creating data_kafka ... done
Creating data_redis ... done
Creating data_kafka_ui ... done
```

**Verify:**
- Kafka UI: http://localhost:8080
- Redis: Check with `redis-cli ping`
- PostgreSQL: Port 5432 ready
- Kafka: Port 9092 ready

---

### Step 2: Start Processing Service (Crawlers + Consumer)
```bash
cd services/processing-service
docker-compose up -d
```

**Expected Output:**
```
Building hn-crawler
Building devto-crawler
Building medium-crawler
Building consumer
Creating processing_hn_crawler ... done
Creating processing_devto_crawler ... done
Creating processing_medium_crawler ... done
Creating processing_consumer ... done
```

**Monitor:**
```bash
docker-compose logs -f
```

You should see:
- Crawlers fetching data from HN, DevTo, Medium
- Consumer reading from Kafka
- Data being saved to PostgreSQL

---

### Step 3: Start API Service (REST API + Dashboard)
```bash
cd services/api-service
docker-compose up -d
```

**Expected Output:**
```
Building api
Creating api_server ... done
```

**Access Dashboard:**
- 🌐 http://localhost:8888
- 📊 API: http://localhost:8888/api/stats
- 🏥 Health: http://localhost:8888/api/health

---

## 📁 Project Structure

```
services/
│
├── data-service/                 # Data Layer (Infrastructure)
│   ├── docker-compose.yml
│   ├── .env                      # Configuration
│   ├── .env.example
│   ├── config/                   # Shared config package
│   ├── internal/                 # Shared Go packages
│   ├── migrations/               # Database migrations
│   └── README.md
│
├── processing-service/           # Processing Layer (Crawlers + Consumer)
│   ├── docker-compose.yml
│   ├── .env                      # Configuration
│   ├── .env.example
│   ├── Dockerfile.hn-crawler     # HackerNews crawler
│   ├── Dockerfile.devto-crawler  # DevTo crawler
│   ├── Dockerfile.medium-crawler # Medium crawler
│   ├── Dockerfile.consumer       # Kafka consumer
│   ├── config/                   # Shared config package
│   ├── internal/                 # Shared Go packages
│   ├── cmd/
│   │   ├── consumer/             # Kafka consumer source
│   │   └── crawlers/
│   │       ├── hn/               # HN crawler source
│   │       ├── devto/            # DevTo crawler source
│   │       └── medium/           # Medium crawler source
│   └── README.md
│
├── api-service/                  # API Layer (REST API + Web)
│   ├── docker-compose.yml
│   ├── .env                      # Configuration
│   ├── .env.example
│   ├── Dockerfile.api            # API server
│   ├── config/                   # Shared config package
│   ├── internal/                 # Shared Go packages
│   ├── cmd/
│   │   └── api/                  # API server source
│   ├── web/                      # HTML/CSS/JS dashboard
│   │   └── index.html
│   └── README.md
│
└── README.md                     # This file
```

---

## 🔌 Service Communication

All services communicate through **environment variables**:

### Data Service → Processing Service
```env
KAFKA_HOST=kafka:29092
REDIS_HOST=redis:6379
PG_HOST=postgres
```

### Data Service → API Service
```env
REDIS_HOST=redis:6379
PG_HOST=postgres
```

**No hardcoded values - all configurable via .env files**

---

## 🛠️ Commands Reference

### View Logs
```bash
# Data service
cd services/data-service && docker-compose logs -f

# Processing service
cd services/processing-service && docker-compose logs -f

# API service  
cd services/api-service && docker-compose logs -f
```

### Stop Services
```bash
cd services/data-service && docker-compose down
cd services/processing-service && docker-compose down
cd services/api-service && docker-compose down
```

### Remove All Data
```bash
cd services/data-service && docker-compose down -v
cd services/processing-service && docker-compose down -v
cd services/api-service && docker-compose down -v
```

### Rebuild Images
```bash
cd services/processing-service && docker-compose build --no-cache
cd services/api-service && docker-compose build --no-cache
```

### Check Container Status
```bash
docker ps -a | grep data_
docker ps -a | grep processing_
docker ps -a | grep api_
```

---

## 📊 API Endpoints

### Health Check
```bash
curl http://localhost:8888/api/health
```

### Overall Statistics
```bash
curl http://localhost:8888/api/stats
```

### Statistics by Topic
```bash
curl http://localhost:8888/api/topics
```

### Statistics by Sentiment
```bash
curl http://localhost:8888/api/sentiment
```

### Top Authors
```bash
curl http://localhost:8888/api/authors
```

### Recent Posts
```bash
curl http://localhost:8888/api/recent
```

---

## 🧪 Testing Data Flow

### 1. Verify Data Service is Ready
```bash
# Check PostgreSQL
psql -h localhost -U postgres -d social_insight -c "SELECT COUNT(*) FROM posts;"

# Check Redis
redis-cli PING

# Check Kafka
docker exec data_kafka kafka-topics --bootstrap-server localhost:9092 --list
```

### 2. Monitor Processing Service
```bash
cd services/processing-service
docker-compose logs -f | grep -E "(crawler|consumer)"
```

You should see:
- Crawlers sending posts to Kafka
- Consumer processing messages
- Data saved to PostgreSQL

### 3. Verify API Service
```bash
curl http://localhost:8888/api/health
# Expected: {"status":"ok","time":"2024-01-28T..."}

curl http://localhost:8888/api/stats
# Expected: Shows statistics in JSON format
```

---

## ⚙️ Configuration

Each service has its own `.env` file. Copy `.env.example` to `.env` and customize:

### Data Service (.env)
```env
DB_USER=postgres
DB_PASSWORD=postgres123
DB_PORT=5432

KAFKA_PORT=9092
REDIS_PORT=6379
ZOOKEEPER_PORT=2181
```

### Processing Service (.env)
```env
KAFKA_HOST=kafka:29092
REDIS_HOST=redis:6379
PG_HOST=postgres

HN_CRAWL_INTERVAL=5m
DEVTO_CRAWL_INTERVAL=10m
MEDIUM_CRAWL_INTERVAL=10m
```

### API Service (.env)
```env
API_PORT=:8888
REDIS_HOST=redis:6379
PG_HOST=postgres
```

---

## 🐛 Troubleshooting

### Containers Won't Start
```bash
# Check logs
docker-compose logs -f

# Verify ports are available
netstat -an | grep LISTEN

# Check Docker resources
docker system df
```

### Connection Refused
- Ensure Data Service is running first
- Check environment variables match between services
- Verify network: `docker network ls`

### Kafka Connectivity Issues
```bash
# Test Kafka from host
docker exec data_kafka kafka-topics --bootstrap-server localhost:9092 --list

# Test from another container
docker exec processing_consumer nc -zv kafka 29092
```

### PostgreSQL Connection Failed
```bash
# Test PostgreSQL
psql -h localhost -U postgres -d social_insight

# Check PostgreSQL logs
docker logs data_postgres
```

### Redis Connection Failed
```bash
# Test Redis
redis-cli -h localhost ping

# Check Redis logs
docker logs data_redis
```

---

## 📝 Development Workflow

### Adding New Crawler
1. Create `cmd/crawlers/new-crawler/main.go` in processing-service
2. Add `Dockerfile.new-crawler`
3. Update `docker-compose.yml` with new service
4. Update environment variables in `.env`

### Adding New API Endpoint
1. Edit `cmd/api/main.go` in api-service
2. Test locally: `go run cmd/api/main.go`
3. Rebuild Docker image: `docker build -f Dockerfile.api -t api .`

### Database Schema Changes
1. Create new migration in `data-service/migrations/`
2. Migration runs automatically on Docker startup
3. Update models in `internal/models/` if needed

---

## 🔐 Security Notes

### For Production:
1. Change PostgreSQL password
2. Change Redis password (add AUTH)
3. Use environment secrets management (Vault, AWS Secrets Manager)
4. Enable HTTPS for API
5. Use private Docker registry
6. Enable PostgreSQL SSL connections
7. Use network policies/firewall rules

### Current Setup (Development Only):
- Default credentials: `postgres:postgres123`
- No authentication on Redis/Kafka
- HTTP only (no HTTPS)
- Open networks between containers

---

## 📚 Related Documentation

- [Data Service README](./services/data-service/README.md)
- [Processing Service README](./services/processing-service/README.md)
- [API Service README](./services/api-service/README.md)

---

## 🤝 Contributing

When modifying services:
1. Never commit `.env` files (use `.env.example`)
2. Update environment variable documentation
3. Test all 3 services together before committing
4. Update this README if adding new features

---

## 📞 Support

For issues or questions:
1. Check individual service README files
2. Review Docker logs: `docker-compose logs -f`
3. Verify environment variables match `.env.example`
4. Ensure services are started in correct order: Data → Processing → API

---

**Last Updated**: January 28, 2025  
**Architecture Version**: 2.0 (Microservices)  
**Go Version**: 1.21+  
**Docker Version**: 20.10+
