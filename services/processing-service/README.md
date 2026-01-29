# 🔄 Processing Service (Crawlers & Consumer)

**Components**: HN Crawler, DevTo Crawler, Medium Crawler, Kafka Consumer  
**Dependency**: Data Service must be running  
**Network**: `processing_network`  

---

## Quick Start

### Prerequisites
```bash
# Ensure Data Service is running
cd ../data-service
docker-compose ps

# Expected: All containers healthy
```

### Start Processing Service
```bash
docker-compose up -d --build
```

### Monitor Services
```bash
docker-compose logs -f
```

---

## Components

### HackerNews Crawler
- **Container**: `processing_hn_crawler`
- **Purpose**: Scrape top stories from HackerNews
- **Interval**: Configurable (default: 5 minutes)
- **Output**: Posts to Kafka `raw_posts` topic

```bash
# View logs
docker logs processing_hn_crawler -f

# Expected output:
# ✅ Connected to Kafka
# 📊 Fetching top 30 stories from HackerNews
# 📤 Sent X posts to Kafka
```

### DevTo Crawler
- **Container**: `processing_devto_crawler`
- **Purpose**: Scrape latest articles from Dev.to
- **Interval**: Configurable (default: 10 minutes)
- **Tags**: AI, Machine Learning, Cloud, DevOps, Startups

```bash
# View logs
docker logs processing_devto_crawler -f

# Expected output:
# ✅ Connected to Kafka
# 📊 Fetching posts for tags: ai, machine-learning, ...
# 📤 Sent X posts to Kafka
```

### Medium Crawler
- **Container**: `processing_medium_crawler`
- **Purpose**: Scrape articles from Medium.com
- **Interval**: Configurable (default: 10 minutes)
- **Topics**: AI, Cloud Computing, DevOps

```bash
# View logs
docker logs processing_medium_crawler -f

# Expected output:
# ✅ Connected to Kafka
# 📊 Fetching posts for topics: machine-learning, ...
# 📤 Sent X posts to Kafka
```

### Kafka Consumer
- **Container**: `processing_consumer`
- **Purpose**: Read posts from Kafka → Save to PostgreSQL + Redis
- **Batch Size**: 500 posts (configurable)
- **Flush Interval**: 2 seconds

```bash
# View logs
docker logs processing_consumer -f

# Expected output:
# 👂 Listening to Kafka topic: raw_posts
# 💾 Batch 1: Saved 500 posts to PostgreSQL
# 🎯 Total processed: 1500 posts
```

---

## Environment Configuration

Edit `.env` to customize behavior:

```env
# Kafka Connection (from Data Service)
KAFKA_HOST=kafka:29092              # Internal Docker network
KAFKA_TOPIC=raw_posts

# Redis Connection (from Data Service)
REDIS_HOST=redis:6379

# PostgreSQL Connection (from Data Service)
PG_HOST=postgres
PG_PORT=5432
PG_USER=postgres
PG_PASSWORD=postgres123
PG_DBNAME=social_insight

# HackerNews Crawler
HN_CRAWL_INTERVAL=5m               # Crawl every 5 minutes
HN_STORIES_LIMIT=30                # Get top 30 stories

# DevTo Crawler
DEVTO_CRAWL_INTERVAL=10m           # Crawl every 10 minutes
DEVTO_POSTS_PER_TAG=6              # 6 posts per tag
DEVTO_TAGS=ai,machine-learning,cloud,devops,startups

# Medium Crawler
MEDIUM_CRAWL_INTERVAL=10m          # Crawl every 10 minutes
MEDIUM_POSTS_PER_TOPIC=10          # 10 posts per topic
MEDIUM_TOPICS=machine-learning,artificial-intelligence,cloud-computing,devops,startups

# Consumer Configuration
CONSUMER_GROUP=social_insight_consumer
CONSUMER_BATCH_SIZE=500            # Process 500 posts at a time
CONSUMER_FLUSH_INTERVAL=2s         # Flush to DB every 2 seconds
```

---

## Data Flow

```
HN Crawler
  │
  ├──► Kafka (raw_posts topic)
  │                    │
DevTo Crawler          │
  │                    ├──► Kafka Consumer
  ├──► Kafka           │        │
  │                    │        ├──► PostgreSQL (storage)
Medium Crawler         │        │
  │                    │        └──► Redis (cache)
  └──► Kafka ──────────┘
```

---

## Monitoring

### Crawler Status
```bash
# Check if crawlers are running
docker-compose ps

# View crawler logs
docker logs processing_hn_crawler -f
docker logs processing_devto_crawler -f
docker logs processing_medium_crawler -f
```

### Consumer Status
```bash
# Monitor consumer progress
docker logs processing_consumer -f

# Expected every 5 seconds:
# 📊 Total processed: X posts
```

### Kafka Monitoring
```bash
# Check topic messages
docker exec data_kafka kafka-topics --bootstrap-server kafka:29092 --describe --topic raw_posts

# Monitor consumer group
docker exec data_kafka kafka-consumer-groups --bootstrap-server kafka:29092 --group social_insight_consumer --describe

# Check offset lag
docker exec data_kafka kafka-consumer-groups --bootstrap-server kafka:29092 --group social_insight_consumer --describe | grep raw_posts
```

---

## Performance Metrics

Monitor via logs:

```bash
# All services
docker-compose logs -f | grep -E "(Sent|Saved|processed)"

# Consumer only
docker-compose logs consumer -f | grep "Total processed"
```

Expected performance:
- **HN Crawler**: ~30 posts per crawl interval
- **DevTo Crawler**: ~30 posts per crawl interval (6 posts × 5 tags)
- **Medium Crawler**: ~30 posts per crawl interval (10 posts × 3 topics)
- **Consumer**: ~500 posts per batch (configurable)

---

## Commands

### Build & Start
```bash
# Build images first
docker-compose build

# Start in background
docker-compose up -d

# Start with logs visible
docker-compose up
```

### View Logs
```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f hn-crawler
docker-compose logs -f devto-crawler
docker-compose logs -f medium-crawler
docker-compose logs -f consumer

# Tail last 100 lines
docker-compose logs -f --tail=100
```

### Stop Services
```bash
# Stop all (containers remain)
docker-compose stop

# Stop and remove
docker-compose down

# Remove everything including volumes
docker-compose down -v
```

### Restart
```bash
# Restart all
docker-compose restart

# Restart specific service
docker-compose restart consumer
```

### Rebuild Images
```bash
# Build without cache
docker-compose build --no-cache

# Build specific service
docker-compose build --no-cache hn-crawler
```

---

## Troubleshooting

### Crawlers Not Starting
```bash
# Check logs for errors
docker logs processing_hn_crawler

# Verify Kafka is reachable
docker exec processing_hn_crawler nc -zv kafka 29092

# Check if Data Service is running
docker ps | grep data_kafka
```

### Consumer Not Processing
```bash
# Check logs
docker logs processing_consumer -f

# Verify Kafka has messages
docker exec data_kafka kafka-console-consumer --bootstrap-server kafka:29092 --topic raw_posts --max-messages 5

# Check consumer group offset
docker exec data_kafka kafka-consumer-groups --bootstrap-server kafka:29092 --group social_insight_consumer --describe
```

### Connection Refused Errors
```bash
# Verify network connectivity
docker exec processing_consumer ping kafka
docker exec processing_consumer redis-cli -h redis ping
docker exec processing_consumer psql -h postgres -U postgres -d social_insight -c "SELECT 1"

# Check environment variables in container
docker exec processing_consumer env | grep -E "KAFKA_|REDIS_|PG_"
```

### High Memory/CPU Usage
```bash
# Monitor resource usage
docker stats processing_hn_crawler processing_devto_crawler processing_medium_crawler processing_consumer

# Adjust batch size if needed
# Edit .env: CONSUMER_BATCH_SIZE=250

# Reduce crawl frequency
# Edit .env: HN_CRAWL_INTERVAL=10m (instead of 5m)
```

### Messages Stuck in Kafka
```bash
# Check consumer lag
docker exec data_kafka kafka-consumer-groups --bootstrap-server kafka:29092 --group social_insight_consumer --describe

# Reset offset to latest
docker exec data_kafka kafka-consumer-groups --bootstrap-server kafka:29092 --group social_insight_consumer --reset-offsets --to-latest --execute --topic raw_posts

# Or reset to beginning (reprocess all)
docker exec data_kafka kafka-consumer-groups --bootstrap-server kafka:29092 --group social_insight_consumer --reset-offsets --to-earliest --execute --topic raw_posts
```

---

## Development

### Adding New Crawler

1. **Create Dockerfile**
   ```dockerfile
   # Dockerfile.new-crawler
   FROM golang:1.21-alpine AS builder
   WORKDIR /app
   COPY . .
   RUN CGO_ENABLED=0 GOOS=linux go build -o /app/new-crawler ./cmd/crawlers/new
   FROM alpine:latest
   WORKDIR /root/
   COPY --from=builder /app/new-crawler .
   ENTRYPOINT ["./new-crawler"]
   ```

2. **Update docker-compose.yml**
   ```yaml
   new-crawler:
     build:
       context: .
       dockerfile: Dockerfile.new-crawler
     container_name: processing_new_crawler
     environment:
       KAFKA_BROKERS: ${KAFKA_HOST:-kafka:29092}
       KAFKA_TOPIC: ${KAFKA_TOPIC:-raw_posts}
   ```

3. **Update .env**
   ```env
   NEW_CRAWLER_INTERVAL=10m
   NEW_CRAWLER_CONFIG=value
   ```

### Testing Crawlers Locally

```bash
# Build and run single crawler
go build -o hn-crawler ./cmd/crawlers/hn
KAFKA_BROKERS=localhost:9092 REDIS_ADDR=localhost:6379 ./hn-crawler

# Or with Docker
docker build -f Dockerfile.hn-crawler -t hn-crawler .
docker run -e KAFKA_BROKERS=kafka:9092 hn-crawler
```

---

## Performance Optimization

### Increase Throughput
```env
# Increase batch size
CONSUMER_BATCH_SIZE=1000

# Increase crawl frequency
HN_CRAWL_INTERVAL=2m
DEVTO_CRAWL_INTERVAL=5m
MEDIUM_CRAWL_INTERVAL=5m
```

### Reduce Resource Usage
```env
# Decrease batch size
CONSUMER_BATCH_SIZE=100

# Decrease crawl frequency
HN_CRAWL_INTERVAL=15m
DEVTO_CRAWL_INTERVAL=30m
MEDIUM_CRAWL_INTERVAL=30m
```

---

## Network Details

### Docker Network
- **Name**: `processing_network`
- **Type**: Bridge
- **Connects to**: Data Service containers via overlay

### Inter-Service Communication
```
Processing Service ──► Data Service
  (kafka:29092) ──► Kafka (29092 - internal)
  (redis:6379) ──► Redis
  (postgres:5432) ──► PostgreSQL
```

---

## Next Steps

After Processing Service is stable:
1. Check logs: `docker-compose logs -f`
2. Verify data in PostgreSQL:
   ```bash
   psql -h localhost -U postgres -d social_insight -c "SELECT COUNT(*) FROM posts;"
   ```
3. Start API Service: `cd ../api-service && docker-compose up -d`
4. View data in dashboard: http://localhost:8888

---

**Last Updated**: January 28, 2025  
**Go Version**: 1.21  
**Dependencies**: Kafka, Redis, PostgreSQL
