# 🔵 Data Service (Infrastructure Layer)

**Components**: PostgreSQL, Redis, Kafka, Zookeeper  
**Network**: `data_network`  
**Status**: ✅ Stateless Infrastructure  

---

## Quick Start

```bash
docker-compose up -d
```

Wait for all containers to be healthy (~10 seconds):
```bash
docker-compose ps
```

---

## Components

### PostgreSQL (Port 5432)
- **Purpose**: Primary data storage for all posts
- **Container**: `data_postgres`
- **Volumes**: `postgres_data`
- **Health Check**: Connected within 10 seconds

```bash
# Connect to PostgreSQL
psql -h localhost -U postgres -d social_insight

# Check tables
\dt

# Count posts
SELECT COUNT(*) FROM posts;
```

### Redis (Port 6379)
- **Purpose**: Cache layer for fast data access
- **Container**: `data_redis`
- **Volumes**: `redis_data`
- **Persistence**: AOF (Append-Only File)

```bash
# Connect to Redis
redis-cli

# Check connection
PING

# View cache keys
KEYS *

# Monitor activity
MONITOR
```

### Kafka (Port 9092)
- **Purpose**: Message queue for data flow
- **Container**: `data_kafka`
- **Volumes**: `kafka_data`
- **Internal Port**: 29092 (for containers within network)

```bash
# List topics
docker exec data_kafka kafka-topics --bootstrap-server localhost:9092 --list

# Describe topic
docker exec data_kafka kafka-topics --bootstrap-server localhost:9092 --describe --topic raw_posts

# Monitor messages
docker exec data_kafka kafka-console-consumer --bootstrap-server localhost:9092 --topic raw_posts --from-beginning
```

### Zookeeper (Port 2181)
- **Purpose**: Kafka cluster coordination
- **Container**: `data_zookeeper`
- **Volumes**: `zookeeper_data`, `zookeeper_log`

```bash
# Check Zookeeper status
docker exec data_zookeeper zookeeper-shell localhost:2181 stat
```

### Kafka UI (Port 8080)
- **Purpose**: Web interface for Kafka management
- **Container**: `data_kafka_ui`
- **Access**: http://localhost:8080

---

## Environment Configuration

Edit `.env` to customize:

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

## Database Schema

Migrations run automatically on startup via `/migrations` volume.

### Posts Table
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
```

---

## Health Checks

All containers have health checks configured:

```bash
docker-compose ps

# Expected output:
# NAME              STATUS
# data_postgres     Up X seconds (healthy)
# data_zookeeper    Up X seconds
# data_kafka        Up X seconds (healthy)
# data_redis        Up X seconds (healthy)
# data_kafka_ui     Up X seconds
```

---

## Commands

### View Logs
```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f postgres
docker-compose logs -f kafka
docker-compose logs -f redis
```

### Stop Services
```bash
docker-compose down
```

### Remove All Data (Reset)
```bash
docker-compose down -v
```

### Restart Services
```bash
docker-compose restart
```

### Rebuild (if Dockerfiles changed)
```bash
docker-compose build --no-cache
docker-compose up -d
```

---

## Monitoring

### Container Stats
```bash
# CPU, Memory, Network
docker stats data_postgres data_redis data_kafka

# Detailed stats
docker container stats --no-stream
```

### Kafka Monitoring
```bash
# Open in browser
# http://localhost:8080

# Command line
docker exec data_kafka kafka-broker-api-versions --bootstrap-server localhost:9092
```

### PostgreSQL Monitoring
```bash
# Connect and check statistics
psql -h localhost -U postgres -d social_insight

# Check table sizes
SELECT tablename, pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) 
FROM pg_tables WHERE schemaname = 'public';

# Check slow queries (if logging enabled)
SELECT query, mean_time FROM pg_stat_statements ORDER BY mean_time DESC LIMIT 10;
```

### Redis Monitoring
```bash
redis-cli

# Memory usage
INFO memory

# Client connections
CLIENT LIST

# Key statistics
DBSIZE

# Memory fragmentation
INFO stats
```

---

## Troubleshooting

### PostgreSQL Won't Start
```bash
# Check logs
docker logs data_postgres

# Verify port
netstat -an | findstr 5432

# Check disk space
df -h /var/lib/postgresql
```

### Kafka Health Check Failed
```bash
# Verify Zookeeper is running
docker logs data_zookeeper

# Check Kafka startup
docker logs data_kafka

# Test connection manually
nc -zv localhost 9092
```

### Redis Memory Issues
```bash
# Check memory usage
redis-cli INFO memory

# Clear cache if needed
redis-cli FLUSHDB

# Increase memory if available
# Edit docker-compose.yml and increase memory limit
```

### Data Persistence Issues
```bash
# Verify volumes exist
docker volume ls

# Check volume usage
docker volume inspect postgres_data
docker volume inspect redis_data
docker volume inspect kafka_data

# Backup data before cleanup
tar -czf postgres_backup.tar.gz /var/lib/postgresql/data
```

---

## Performance Tuning

### PostgreSQL
```bash
# Increase connections if needed
# Edit docker-compose.yml environment:
#   POSTGRES_INIT_ARGS: "-c max_connections=200"

# Enable query logging
#   POSTGRES_INIT_ARGS: "-c log_statement=all -c log_duration=on"
```

### Kafka
Configured in docker-compose.yml:
- `KAFKA_MESSAGE_MAX_BYTES=10485760` (10MB)
- `KAFKA_REPLICA_FETCH_MAX_BYTES=10485760` (10MB)
- `KAFKA_NUM_NETWORK_THREADS=8` (default)

### Redis
```bash
# Monitor performance
redis-cli SLOWLOG GET 10

# Optimize persistence
# Edit docker-compose.yml:
#   command: redis-server --appendonly yes --appendfsync everysec
```

---

## Backup & Recovery

### Backup PostgreSQL
```bash
docker exec data_postgres pg_dump -U postgres -d social_insight > backup.sql
```

### Restore PostgreSQL
```bash
docker exec -i data_postgres psql -U postgres -d social_insight < backup.sql
```

### Backup Kafka Topics
```bash
# Topics are in zookeeper, no direct backup needed
# Config is defined in docker-compose.yml and can be version controlled
```

### Backup Redis
```bash
docker exec data_redis redis-cli BGSAVE
docker cp data_redis:/data/dump.rdb ./redis_backup.rdb
```

---

## Network Details

### Docker Network
- **Name**: `data_network`
- **Type**: Bridge
- **Internal DNS**: Container names are resolved

### Internal Addresses (for containers)
- PostgreSQL: `postgres:5432`
- Redis: `redis:6379`
- Kafka: `kafka:29092` (internal) or `localhost:9092` (external)
- Zookeeper: `zookeeper:2181`

### External Addresses (from host)
- PostgreSQL: `localhost:5432`
- Redis: `localhost:6379`
- Kafka: `localhost:9092`
- Zookeeper: `localhost:2181`
- Kafka UI: `http://localhost:8080`

---

## Security

### For Development (Current)
✅ Acceptable - all defaults used

### For Production
⚠️ Required changes:
1. Change PostgreSQL credentials
2. Enable Redis password
3. Enable Kafka authentication
4. Use SSL/TLS for all connections
5. Restrict network access
6. Enable PostgreSQL WAL archiving for backup
7. Monitor resource limits

---

## Next Steps

Once Data Service is running:
1. Start Processing Service: `cd ../processing-service && docker-compose up -d`
2. Monitor logs for crawlers and consumer
3. Start API Service: `cd ../api-service && docker-compose up -d`
4. Access dashboard: http://localhost:8888

---

**Last Updated**: January 28, 2025  
**Components**: PostgreSQL 15, Redis 7, Kafka 7.5, Zookeeper 7.5
