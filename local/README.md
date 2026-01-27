# 🔵 Social Insight - Local Development Guide

**Environment**: Local Development with Docker Compose  
**Status**: ✅ Ready to Use  
**Setup Time**: 5-15 minutes

---

## ⚡ Quick Start (5 minutes)

### Prerequisites
- Docker & Docker Compose installed
- 4GB RAM available
- Port 5432, 6379, 9092, 8888 available

### 1. Start All Services
```bash
cd local/
docker-compose up -d --build
```

### 2. Access Dashboard
```
🌐 Open: http://localhost:8888
📊 API: http://localhost:8888/api/stats
```

### 3. View Logs
```bash
docker-compose logs -f
```

### 4. Stop Services
```bash
docker-compose down
```

---

## 📁 Folder Structure

```
local/
├── README.md                    # This file
├── docker-compose.yml           # Full-stack local compose
├── docker-compose.local.yml     # Infrastructure only (DB, Redis, Kafka)
│
├── ingestion/                   # Layer 1: Data Crawlers
│   ├── hn-crawler/              # HackerNews crawler
│   ├── medium-crawler/          # Medium.com crawler
│   └── devto-crawler/           # Dev.to crawler
│
├── processing/                  # Layer 2: Processing & ML
│   ├── consumer/                # Kafka consumer
│   └── ml-service/              # Python ML service (sentiment, trends, AI)
│
├── presentation/                # Layer 3: API & UI
│   ├── api-gateway/             # REST API gateway
│   └── frontend/                # HTML/CSS/JS dashboard
│
├── shared/                      # Shared Code
│   ├── config/                  # Configuration
│   ├── database/                # PostgreSQL helpers
│   ├── kafka/                   # Kafka utilities
│   ├── logger/                  # Logging
│   ├── models/                  # Data models
│   └── redis/                   # Redis cache
│
└── infrastructure/              # Database Setup
    ├── init.sql                 # Initial schema
    └── migrations/
```

---

## 🏗️ 3-Layer Architecture

```
┌─────────────────────────────────────────────────────┐
│              INGESTION LAYER (Layer 1)              │
│                                                     │
│  HN Crawler      Medium Crawler      DevTo Crawler │
│        │                │                  │        │
│        └────────────────┬──────────────────┘        │
│                         │                          │
└─────────────────────────┼──────────────────────────┘
                          │ Kafka Topic: posts
                          ▼
┌─────────────────────────────────────────────────────┐
│             PROCESSING LAYER (Layer 2)              │
│                                                     │
│        Consumer → ML Service → PostgreSQL           │
│        • Read Kafka messages                        │
│        • Run sentiment analysis                     │
│        • Detect AI models                           │
│        • Save to database                           │
│                                                     │
└─────────────────────────────────────────────────────┘
                          │
                PostgreSQL + Redis
                          │
                          ▼
┌─────────────────────────────────────────────────────┐
│            PRESENTATION LAYER (Layer 3)             │
│                                                     │
│          API Gateway ←→ Frontend Dashboard         │
│          • REST API endpoints                       │
│          • Real-time dashboard                      │
│          • Statistics & analytics                   │
│                                                     │
└─────────────────────────────────────────────────────┘
```

---

## 🐳 Docker Services

| Service | Port | Language | Purpose |
|---------|------|----------|---------|
| **hn-crawler** | - | Go | Crawl HackerNews top stories |
| **medium-crawler** | - | Go | Crawl Medium AI articles |
| **devto-crawler** | - | Go | Crawl Dev.to posts |
| **kafka** | 9092 | Java | Message queue |
| **postgres** | 5432 | SQL | Data storage |
| **redis** | 6379 | C | Caching layer |
| **consumer** | - | Go | Process Kafka messages |
| **ml-service** | 8000 | Python | ML inference |
| **api-gateway** | 8888 | Go | REST API |

---

## 📝 Environment Configuration

### Default .env for Local
```bash
# Database
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=social_insight

# Redis
REDIS_ADDR=redis:6379
REDIS_DB=0

# Kafka
KAFKA_BROKERS=kafka:9092
KAFKA_TOPIC=posts
KAFKA_CONSUMER_GROUP=social-insight-group

# Services
API_PORT=8888
ML_SERVICE_URL=http://ml-service:8000
LOG_LEVEL=info

# Crawlers
CRAWL_INTERVAL=30s
FETCH_LIMIT=50
```

---

## 🚀 Common Commands

### View Status
```bash
# List all containers
docker-compose ps

# Show service logs
docker-compose logs <service-name>

# Follow logs in real-time
docker-compose logs -f

# Show logs for specific service
docker-compose logs -f consumer
```

### Access Services
```bash
# PostgreSQL (from your machine)
psql -h localhost -U postgres -d social_insight

# Redis CLI
redis-cli -h localhost

# Kafka topics
docker-compose exec kafka kafka-topics.sh --list --bootstrap-server localhost:9092

# View Kafka messages
docker-compose exec kafka kafka-console-consumer.sh \
  --bootstrap-server localhost:9092 \
  --topic posts \
  --from-beginning
```

### Development Workflow
```bash
# Start fresh (clean slate)
docker-compose down -v
docker-compose up -d --build

# Rebuild a specific service
docker-compose build api-gateway
docker-compose up -d api-gateway

# View container details
docker-compose exec postgres psql -U postgres -d social_insight -c "SELECT COUNT(*) FROM posts;"
```

### Database Operations
```bash
# Connect to database
docker-compose exec postgres psql -U postgres -d social_insight

# View posts
SELECT * FROM posts LIMIT 10;

# Count posts
SELECT COUNT(*) FROM posts;

# Check sentiment distribution
SELECT sentiment, COUNT(*) FROM posts GROUP BY sentiment;

# View latest posts
SELECT title, source, created_at FROM posts ORDER BY created_at DESC LIMIT 20;
```

---

## 🔍 Testing the Pipeline

### 1. Check Crawlers are Running
```bash
docker-compose logs hn-crawler | tail -20
docker-compose logs medium-crawler | tail -20
docker-compose logs devto-crawler | tail -20
```

Should see messages like:
```
[INFO] Successfully scraped 15 posts from HN
[INFO] Sending posts to Kafka topic: posts
```

### 2. Check Consumer is Processing
```bash
docker-compose logs consumer | tail -20
```

Should see:
```
[INFO] Received 50 messages from Kafka
[INFO] Saved 50 posts to database
```

### 3. Check API is Working
```bash
curl http://localhost:8888/api/health
curl http://localhost:8888/api/stats
curl http://localhost:8888/api/posts
```

Should return JSON responses.

### 4. Check Dashboard
Open http://localhost:8888 in your browser

Should see:
- Live post count
- Sentiment distribution
- Trending topics
- Recent posts list

### 5. Verify Database
```bash
docker-compose exec postgres psql -U postgres -d social_insight

# Check table exists
\dt

# Count posts
SELECT COUNT(*) FROM posts;

# Show sample posts
SELECT id, title, source, sentiment FROM posts LIMIT 5;
```

---

## 🎯 Data Flow Verification

```
Step 1: Crawlers scrape data
  └─ docker-compose logs hn-crawler

Step 2: Data published to Kafka
  └─ docker-compose exec kafka kafka-console-consumer.sh \
     --bootstrap-server kafka:9092 \
     --topic posts \
     --from-beginning \
     --max-messages 5

Step 3: Consumer reads and processes
  └─ docker-compose logs consumer

Step 4: Data saved to PostgreSQL
  └─ docker-compose exec postgres psql -U postgres -d social_insight -c \
     "SELECT COUNT(*) FROM posts;"

Step 5: API serves data
  └─ curl http://localhost:8888/api/stats

Step 6: Dashboard displays data
  └─ Open http://localhost:8888
```

---

## 🐛 Troubleshooting

### Issue: "Port Already in Use"
```
Error: bind: address already in use

Solution:
  # Find which service uses the port
  lsof -i :8888
  
  # Kill the process or change port in docker-compose.yml
  kill -9 <PID>
  
  # Or change the port mapping
  # Change "8888:8888" to "8889:8888" in docker-compose.yml
```

### Issue: "Cannot Connect to Docker"
```
Error: Cannot connect to Docker daemon

Solution:
  # Make sure Docker is running
  docker ps
  
  # On Windows, ensure Docker Desktop is running
  # On Mac, restart Docker Desktop
  # On Linux, start Docker service:
  sudo systemctl start docker
```

### Issue: "Out of Memory"
```
Error: cannot allocate memory

Solution:
  # Stop other containers
  docker-compose down
  
  # Check memory usage
  docker stats
  
  # If persistent, increase Docker memory limit
  # Docker Desktop → Preferences → Resources → Memory
```

### Issue: "Kafka Connection Error"
```
Error: kafka: Failed to connect to broker

Solution:
  # Check if Kafka is running
  docker-compose ps | grep kafka
  
  # View Kafka logs
  docker-compose logs kafka
  
  # Restart Kafka
  docker-compose restart kafka
  
  # Wait 10-15 seconds for Kafka to fully start
  sleep 15
```

### Issue: "Database Connection Error"
```
Error: psql: could not connect to server

Solution:
  # Check PostgreSQL is running
  docker-compose ps | grep postgres
  
  # Check environment variables
  grep DB_ .env
  
  # View PostgreSQL logs
  docker-compose logs postgres
  
  # Wait for DB to initialize (first startup)
  docker-compose logs postgres | grep "ready to accept"
```

### Issue: "ML Service Not Responding"
```
Error: curl: (7) Failed to connect to ml-service:8000

Solution:
  # Check ML Service is running
  docker-compose ps | grep ml-service
  
  # View logs
  docker-compose logs ml-service
  
  # Check if dependencies are installed
  docker-compose exec ml-service pip list
  
  # Rebuild the image
  docker-compose build ml-service
  docker-compose up -d ml-service
```

### Issue: "No Posts Appearing"
```
Problem: Dashboard shows 0 posts

Solution:
  # 1. Check if crawlers are running
  docker-compose logs hn-crawler | grep -i "error"
  
  # 2. Check if consumer is processing
  docker-compose logs consumer | grep "Saved"
  
  # 3. Check database directly
  docker-compose exec postgres psql -U postgres -d social_insight -c \
    "SELECT COUNT(*) FROM posts;"
  
  # 4. Check API is returning data
  curl http://localhost:8888/api/stats
  
  # 5. Check Kafka has messages
  docker-compose exec kafka kafka-console-consumer.sh \
    --bootstrap-server kafka:9092 \
    --topic posts \
    --from-beginning \
    --max-messages 1
  
  # 6. If still empty, restart all services
  docker-compose down
  docker-compose up -d
  
  # 7. Wait 2-3 minutes for data to flow
```

### Issue: "High CPU Usage"
```
Problem: Docker consuming lots of CPU

Solution:
  # Check which container
  docker stats
  
  # View logs to see if there's an error loop
  docker-compose logs <service-name>
  
  # Reduce crawl frequency in .env
  CRAWL_INTERVAL=60s  (instead of 30s)
  
  # Reduce fetch limit
  FETCH_LIMIT=20      (instead of 50)
  
  # Restart service
  docker-compose restart <service-name>
```

### Issue: "Cannot Delete Containers"
```
Error: Error response from daemon: conflict

Solution:
  # Force remove containers
  docker-compose down -f
  
  # Remove Docker volumes (will delete data!)
  docker-compose down -v
  
  # Clean up all stopped containers
  docker system prune
  
  # Remove all unused Docker resources
  docker system prune -a
```

---

## 📊 Monitoring

### Real-time Container Stats
```bash
docker stats

# Watch specific container
docker stats consumer
```

### View Resource Usage
```bash
# Show container size
docker compose ps -s

# Show image size
docker images --format "table {{.Repository}}\t{{.Size}}"
```

### Logs Analysis
```bash
# Show last 100 lines
docker-compose logs --tail 100

# Follow in real-time
docker-compose logs -f

# Show only errors
docker-compose logs | grep -i error

# Show timestamps
docker-compose logs --timestamps
```

---

## 💡 Pro Tips

### 1. Quick Restart
```bash
# Restart all services
docker-compose restart

# Restart one service
docker-compose restart consumer
```

### 2. Access Container Shell
```bash
# Get bash shell in container
docker-compose exec consumer bash

# Run Python in ML service
docker-compose exec ml-service python -c "print('Hello')"
```

### 3. View Recent Changes
```bash
# Show recent git changes
git diff

# Show git status
git status

# View git log
git log --oneline -10
```

### 4. Performance Optimization
```bash
# Reduce crawl frequency
sed -i 's/CRAWL_INTERVAL=30s/CRAWL_INTERVAL=60s/g' .env

# Reduce consumer batch size
docker-compose exec consumer vi config.go

# Increase Redis cache TTL
docker-compose exec redis redis-cli CONFIG SET timeout 300
```

### 5. Clean Up Space
```bash
# Remove unused images
docker image prune -a

# Remove unused volumes
docker volume prune

# Remove unused networks
docker network prune
```

---

## 🔐 Security for Local Development

⚠️ **For local use only!**

- Default credentials are weak (use for dev only)
- No encryption in local setup
- No authentication on APIs
- All traffic is unencrypted

For production, see [../production/README.md](../production/README.md)

---

## 📦 Tech Stack

| Component | Technology | Version | Purpose |
|-----------|-----------|---------|---------|
| Crawler | Go | 1.21+ | Scrape social media |
| API | Go | 1.21+ | REST endpoints |
| ML Service | Python | 3.10+ | ML inference |
| Consumer | Go | 1.21+ | Message processing |
| Database | PostgreSQL | 15 | Data storage |
| Cache | Redis | 7 | Caching |
| Message Queue | Kafka | 3.x | Async pipeline |
| Container | Docker | 24+ | Packaging |
| Orchestration | Docker Compose | 2.x | Local setup |

---

## 🧪 Development Tips

### Adding New Features
1. Make changes to code
2. Build image: `docker-compose build <service>`
3. Start service: `docker-compose up -d <service>`
4. View logs: `docker-compose logs -f <service>`

### Debugging
```bash
# Print environment variables
docker-compose exec <service> env

# Check file exists in container
docker-compose exec <service> ls -la /path/to/file

# View application config
docker-compose exec <service> cat /app/config.toml

# Debug database connection
docker-compose exec postgres psql -U postgres -d social_insight -c "\conninfo"
```

### Testing Changes
```bash
# Test API locally
curl -X GET http://localhost:8888/api/posts

# Test with headers
curl -X GET http://localhost:8888/api/posts \
  -H "Content-Type: application/json"

# Send data
curl -X POST http://localhost:8888/api/posts \
  -H "Content-Type: application/json" \
  -d '{"title":"test","source":"hn"}'
```

---

## 📞 Getting Help

### Check Logs First
```bash
# All logs
docker-compose logs | grep -i error

# Specific service
docker-compose logs <service-name> | head -50
```

### Verify Setup
```bash
# Check all services running
docker-compose ps

# Check database
docker-compose exec postgres psql -U postgres -c "SELECT version();"

# Check Kafka
docker-compose exec kafka kafka-broker-api-versions.sh --bootstrap-server kafka:9092
```

### Reset Everything
```bash
# Complete reset (WARNING: deletes all data!)
docker-compose down -v
rm -rf ./data/*
docker-compose up -d
```

---

## 🎓 Learning Resources

- **Docker**: https://docs.docker.com
- **Docker Compose**: https://docs.docker.com/compose
- **Kafka**: https://kafka.apache.org/documentation
- **PostgreSQL**: https://www.postgresql.org/docs
- **Go**: https://golang.org/doc
- **Python FastAPI**: https://fastapi.tiangolo.com

---

## 📋 Checklist for First Run

- [ ] Docker & Docker Compose installed
- [ ] Ports 5432, 6379, 9092, 8888 available
- [ ] Enough disk space (~5GB)
- [ ] Enough RAM (~4GB)
- [ ] `docker-compose up -d` runs successfully
- [ ] All containers running: `docker-compose ps`
- [ ] Database accessible: `psql -h localhost -U postgres`
- [ ] API responding: `curl http://localhost:8888/api/stats`
- [ ] Dashboard loads: `http://localhost:8888`
- [ ] Data flowing: See posts in dashboard

---

## 🚀 Next Steps

1. ✅ Local setup complete
2. 👉 **Explore the dashboard**: http://localhost:8888
3. 👉 **Try API endpoints**: curl http://localhost:8888/api/posts
4. 👉 **Check database**: `docker-compose exec postgres psql ...`
5. 👉 **Read code**: Check ingestion/, processing/, presentation/
6. 👉 **Make changes**: Edit code and rebuild services
7. 👉 **Deploy to AWS**: See [../production/README.md](../production/README.md)

---

**Status**: ✅ Ready for Development  
**Last Updated**: January 27, 2026  
**Version**: 1.0
