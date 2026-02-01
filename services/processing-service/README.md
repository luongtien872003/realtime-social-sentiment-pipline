# 🔄 Processing Service - Data Pipeline Layer

[![Go](https://img.shields.io/badge/Go-1.21-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Kafka](https://img.shields.io/badge/Kafka-Consumer-231F20?style=flat&logo=apachekafka)](https://kafka.apache.org/)

**Components**: HN Crawler, DevTo Crawler, Medium Crawler, Kafka Consumer  
**AWS EC2**: `social-insight-processing` (13.251.157.185)  
**Status**: ✅ Production Running

---

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────┐
│                 Processing Service                       │
│                                                          │
│  ┌────────────┐  ┌────────────┐  ┌────────────┐        │
│  │HN Crawler  │  │DevTo       │  │Medium      │        │
│  │(5 min)     │  │Crawler     │  │Crawler     │        │
│  │            │  │(10 min)    │  │(10 min)    │        │
│  └─────┬──────┘  └─────┬──────┘  └─────┬──────┘        │
│        │               │               │                │
│        └───────────────┼───────────────┘                │
│                        ▼                                 │
│              ┌─────────────────┐                        │
│              │   Kafka Topic   │                        │
│              │   (raw_posts)   │                        │
│              └────────┬────────┘                        │
│                       ▼                                  │
│              ┌─────────────────┐                        │
│              │    Consumer     │                        │
│              │  (batch: 500)   │                        │
│              └────────┬────────┘                        │
│                       │                                  │
│           ┌───────────┴───────────┐                     │
│           ▼                       ▼                      │
│    ┌────────────┐          ┌────────────┐              │
│    │ PostgreSQL │          │   Redis    │              │
│    │  (Store)   │          │  (Cache)   │              │
│    └────────────┘          └────────────┘              │
└─────────────────────────────────────────────────────────┘
```

---

## 🚀 Quick Start

```bash
# Ensure Data Service is running first
cd ../data-service && docker-compose ps

# Start Processing Service
docker-compose up -d --build

# Monitor logs
docker-compose logs -f
```

---

## ☁️ AWS Production

| Property | Value |
|----------|-------|
| **EC2 Instance** | social-insight-processing |
| **Public IP** | 13.251.157.185 |
| **Private IP** | 172.31.16.49 |
| **Instance Type** | m7i-flex.large |
| **Security Group** | processing-sg |

### Running Containers
```
processing_hn_crawler       - HackerNews crawler
processing_devto_crawler    - Dev.to crawler
processing_medium_crawler   - Medium crawler
processing_consumer         - Kafka consumer
```

---

## 📦 Crawlers

| Crawler | Source | Interval | Posts/Cycle |
|---------|--------|----------|-------------|
| **HN Crawler** | HackerNews API | 5 min | ~30 stories |
| **DevTo Crawler** | Dev.to API | 10 min | ~30 posts |
| **Medium Crawler** | Medium RSS | 10 min | ~30 posts |

### Crawler Output
```bash
docker logs processing_hn_crawler -f

# Expected:
# ✅ Connected to Kafka
# 📊 Fetching top 30 stories from HackerNews
# 📤 Sent 30 posts to Kafka
```

---

## 📥 Consumer

| Property | Value |
|----------|-------|
| **Batch Size** | 500 posts |
| **Flush Interval** | 2 seconds |
| **Consumer Group** | social_insight_consumer |

### Consumer Output
```bash
docker logs processing_consumer -f

# Expected:
# 👂 Listening to Kafka topic: raw_posts
# 💾 Saved 500 posts to PostgreSQL
# 🎯 Total processed: 1500 posts
```

---

## ⚙️ Environment Variables

```env
# Kafka (connect to Data Service)
KAFKA_HOST=kafka:29092          # Local
KAFKA_HOST=172.31.16.144:9092   # AWS
KAFKA_TOPIC=raw_posts

# Redis
REDIS_HOST=redis:6379           # Local
REDIS_HOST=172.31.16.144:6379   # AWS

# PostgreSQL
PG_HOST=postgres                # Local
PG_HOST=172.31.16.144           # AWS
PG_PORT=5432
PG_USER=postgres
PG_PASSWORD=postgres123
PG_DBNAME=social_insight

# Crawler Settings
HN_CRAWL_INTERVAL=5m
HN_STORIES_LIMIT=30
DEVTO_CRAWL_INTERVAL=10m
DEVTO_POSTS_PER_TAG=6
MEDIUM_CRAWL_INTERVAL=10m
MEDIUM_POSTS_PER_TOPIC=10

# Consumer Settings
CONSUMER_BATCH_SIZE=500
CONSUMER_FLUSH_INTERVAL=2s
```

---

## 🔧 Common Commands

```bash
# View all logs
docker-compose logs -f

# View specific crawler
docker logs processing_hn_crawler -f
docker logs processing_devto_crawler -f
docker logs processing_medium_crawler -f

# View consumer
docker logs processing_consumer -f

# Restart all
docker-compose restart

# Rebuild after code changes
docker-compose up -d --build
```

---

## 📊 Monitoring Kafka

```bash
# List topics
docker exec data_kafka kafka-topics --bootstrap-server kafka:29092 --list

# Check consumer group lag
docker exec data_kafka kafka-consumer-groups \
  --bootstrap-server kafka:29092 \
  --group social_insight_consumer \
  --describe

# View messages in topic
docker exec data_kafka kafka-console-consumer \
  --bootstrap-server kafka:29092 \
  --topic raw_posts \
  --from-beginning \
  --max-messages 5
```

---

## 🛠️ Troubleshooting

| Issue | Solution |
|-------|----------|
| Crawler can't connect to Kafka | Check `KAFKA_HOST` env var |
| Consumer not saving to DB | Verify PostgreSQL is healthy |
| No messages in Kafka | Check crawler logs for errors |

---

> **Last Updated**: January 31, 2026
