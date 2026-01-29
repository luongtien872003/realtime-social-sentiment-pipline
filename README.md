# Social Insight

Real-time social media trend analysis platform with independent microservices. Aggregates posts from HackerNews, Dev.to, and Medium.

## Quick Start

```bash
git clone https://github.com/yourusername/Social_Insight
cd Social_Insight

# Start services (order matters)
cd services/data-service && docker-compose up -d
cd ../processing-service && docker-compose up -d
cd ../api-service && docker-compose up -d

# Dashboard at http://localhost:8888
```

## Architecture

Three independent services with clear separation:

| Service | Runs | Exposes | Depends On |
|---------|------|---------|-----------|
| **data-service** | PostgreSQL, Redis, Kafka, Zookeeper | Database, Cache, Queue | Nothing |
| **processing-service** | 3 Crawlers (HN, Dev.to, Medium), Consumer | Posts via Kafka → DB | data-service |
| **api-service** | REST API, Dashboard | HTTP:8888 | data-service |

Each service has its own `docker-compose.yml` and `config/env.go` with environment-based configuration (no hardcoding).

## API

```
GET /api/health              # Status check
GET /api/stats               # Overall stats
GET /api/recent              # Recent posts
GET /api/topics              # Topic distribution
GET /api/sentiment           # Sentiment distribution
GET /api/insights            # Detected insights
GET /api/trending            # Top trending posts
GET /api/compare             # Today vs yesterday
```

## Stack

- **Frontend**: HTML5, CSS3, JavaScript, Chart.js
- **API**: Go 1.21
- **Queue**: Kafka 7.5
- **Database**: PostgreSQL 15
- **Cache**: Redis 7
- **Container**: Docker & Compose

## Configuration

All services read from environment variables. Example:

```yaml
# data-service/docker-compose.yml
environment:
  POSTGRES_PASSWORD: ${DB_PASSWORD:-postgres123}
  KAFKA_BROKER: kafka:29092
```

No hardcoded secrets or IPs. Safe for AWS EC2 deployment.

## Deployment

### AWS EC2 (Multi-Instance)

Each instance gets ONE service:

**Instance 1** (data-service):
```bash
cd services/data-service && docker-compose up -d
```

**Instance 2** (processing-service):
```bash
# Update config to point to Instance 1 database
export DB_CONN=postgresql://user:pass@instance1-ip:5432/social_insight
cd services/processing-service && docker-compose up -d
```

**Instance 3** (api-service):
```bash
# Update config to point to Instance 1 database
export DB_CONN=postgresql://user:pass@instance1-ip:5432/social_insight
cd services/api-service && docker-compose up -d
```

## Project Structure

```
services/
├── data-service/
│   ├── docker-compose.yml
│   ├── Dockerfile.postgres (if custom)
│   ├── config/env.go
│   └── migrations/001_create_posts.sql
├── processing-service/
│   ├── docker-compose.yml
│   ├── Dockerfile.consumer
│   ├── config/env.go
│   └── cmd/crawlers/
└── api-service/
    ├── docker-compose.yml
    ├── Dockerfile.api
    ├── config/env.go
    ├── web/index.html
    └── cmd/api/main.go
```

## Development

### Edit Dashboard (No Rebuild)
```bash
# Edit HTML/CSS/JS
nano services/api-service/web/index.html

# Refresh browser
http://localhost:8888
```

### Add API Endpoint
```bash
# Create handler
nano services/api-service/internal/myfeature/handler.go

# Register in main.go
# Rebuild
cd services/api-service && docker-compose up -d --build api
```

### Add Crawler
```bash
# Create crawler
nano services/processing-service/internal/crawler/newsource.go

# Register in docker-compose.yml
# Rebuild
cd services/processing-service && docker-compose up -d --build
```

## Verify Setup

```bash
# Health check
curl http://localhost:8888/api/health

# Database
docker exec data_postgres psql -U postgres -d social_insight -c "SELECT COUNT(*) FROM posts;"

# Crawlers
docker logs processing_hn_crawler

# API logs
docker logs api_server -f
```

## Production Checklist

- [ ] Change database password
- [ ] Enable Redis authentication
- [ ] Use Kubernetes instead of Compose
- [ ] Add HTTPS certificates
- [ ] Configure firewall rules
- [ ] Setup monitoring (Prometheus/Grafana)
- [ ] Enable database backups
- [ ] Add rate limiting
- [ ] Implement request logging

## Troubleshooting

**Services won't start?**
```bash
docker-compose logs --tail=50
```

**Database empty?**
```bash
# Wait 5-10 minutes for crawlers
docker logs processing_hn_crawler
```

**API errors?**
```bash
# Check DB connection
docker exec api_server curl http://localhost:8888/api/health
```

## License

MIT

## Contact

Issues: GitHub Issues | Email: support@example.com

---

**Version**: 1.0.0 | **Status**: Production Ready
