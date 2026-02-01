# 🔵 Data Service - Infrastructure Layer

[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-4169E1?style=flat&logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![Redis](https://img.shields.io/badge/Redis-7-DC382D?style=flat&logo=redis&logoColor=white)](https://redis.io/)
[![Kafka](https://img.shields.io/badge/Kafka-7.5-231F20?style=flat&logo=apachekafka)](https://kafka.apache.org/)

**Components**: PostgreSQL, Redis, Kafka, Zookeeper, Kafka UI  
**AWS EC2**: `social-insight-data` (13.250.47.65)  
**Status**: ✅ Production Running

---

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────┐
│                  Data Service                       │
│                                                     │
│  ┌───────────┐  ┌───────────┐  ┌───────────────┐   │
│  │PostgreSQL │  │   Redis   │  │    Kafka      │   │
│  │  :5432    │  │   :6379   │  │    :9092      │   │
│  └───────────┘  └───────────┘  └───────────────┘   │
│                                       │             │
│                               ┌───────────────┐    │
│                               │  Zookeeper    │    │
│                               │    :2181      │    │
│                               └───────────────┘    │
│                                                     │
│  ┌─────────────────────────────────────────────┐   │
│  │              Kafka UI (:8080)               │   │
│  └─────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────┘
```

---

## 🚀 Quick Start

```bash
# Create network (required)
docker network create social_insight_network

# Start all services
docker-compose up -d

# Verify health
docker-compose ps
```

---

## ☁️ AWS Production

| Property | Value |
|----------|-------|
| **EC2 Instance** | social-insight-data |
| **Public IP** | 13.250.47.65 |
| **Private IP** | 172.31.16.144 |
| **Instance Type** | m7i-flex.large |
| **Security Group** | data-sg |

### Production URLs
- Kafka UI: http://13.250.47.65:8080

### Running Containers
```
data_postgres     - PostgreSQL :5432 (healthy)
data_redis        - Redis :6379 (healthy)
data_kafka        - Kafka :9092 (healthy)
data_zookeeper    - Zookeeper :2181
data_kafka_ui     - Kafka UI :8080
```

---

## 📦 Components

| Component | Port | Container | Purpose |
|-----------|------|-----------|---------|
| PostgreSQL | 5432 | data_postgres | Primary data storage |
| Redis | 6379 | data_redis | Cache layer |
| Kafka | 9092 | data_kafka | Message queue |
| Zookeeper | 2181 | data_zookeeper | Kafka coordination |
| Kafka UI | 8080 | data_kafka_ui | Web monitoring |

---

## ⚙️ Environment Variables

```env
# PostgreSQL
DB_USER=postgres
DB_PASSWORD=postgres123
DB_NAME=social_insight
DB_PORT=5432

# Kafka
KAFKA_PORT=9092

# Redis
REDIS_PORT=6379

# Zookeeper
ZOOKEEPER_PORT=2181

# Kafka UI
KAFKA_UI_PORT=8080
```

---

## 🗄️ Database Schema

```sql
CREATE TABLE posts (
    id VARCHAR(255) PRIMARY KEY,
    author VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    topic VARCHAR(50) NOT NULL,
    sentiment VARCHAR(20) NOT NULL,
    likes INTEGER DEFAULT 0,
    comments INTEGER DEFAULT 0,
    shares INTEGER DEFAULT 0,
    platform VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_posts_topic ON posts(topic);
CREATE INDEX idx_posts_sentiment ON posts(sentiment);
CREATE INDEX idx_posts_created_at ON posts(created_at);
```

---

## 🔧 Common Commands

```bash
# View logs
docker-compose logs -f

# Check PostgreSQL
docker exec data_postgres psql -U postgres -d social_insight -c "SELECT COUNT(*) FROM posts;"

# Check Kafka topics
docker exec data_kafka kafka-topics --bootstrap-server kafka:29092 --list

# Check Redis
docker exec data_redis redis-cli PING
```

---

## 🌐 Network Configuration

### For Local Development
```
PostgreSQL: localhost:5432
Redis: localhost:6379
Kafka: localhost:9092
```

### For AWS Multi-Host (Private IP)
```env
# On Processing/API instances, use:
PG_HOST=172.31.16.144
REDIS_HOST=172.31.16.144
KAFKA_HOST=172.31.16.144:9092
```

---

> **Last Updated**: January 31, 2026
