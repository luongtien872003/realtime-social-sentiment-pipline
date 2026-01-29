# 🌐 API Service (REST API + Web Dashboard)

**Components**: REST API Server, Web Dashboard (HTML/CSS/JS)  
**Dependency**: Data Service must be running  
**Network**: `api_network`  
**Access**: http://localhost:8888  

---

## Quick Start

### Prerequisites
```bash
# Ensure Data Service is running
cd ../data-service
docker-compose ps

# Expected: All containers healthy
```

### Start API Service
```bash
docker-compose up -d --build
```

### Access Dashboard
```
🌐 http://localhost:8888
```

### Monitor Logs
```bash
docker-compose logs -f
```

---

## Components

### REST API Server
- **Port**: 8888
- **Container**: `api_server`
- **Purpose**: Serve JSON API endpoints for frontend
- **Technology**: Go HTTP Server (No framework)

```bash
# View logs
docker logs api_server -f

# Expected output:
# ✅ Connected to PostgreSQL
# ✅ Connected to Redis
# 🚀 API Server running on http://localhost:8888
```

### Web Dashboard
- **Location**: `/web/index.html`
- **Technology**: HTML + CSS + JavaScript + Chart.js
- **Features**: 
  - Real-time statistics
  - Topic trends
  - Sentiment analysis
  - Author rankings
  - Recent posts feed

```bash
# Access directly
open http://localhost:8888
```

---

## Environment Configuration

Edit `.env` to customize:

```env
# API Server Configuration
API_PORT=:8888                      # Listen on port 8888

# Redis Connection (from Data Service)
REDIS_HOST=redis:6379

# PostgreSQL Connection (from Data Service)
PG_HOST=postgres
PG_PORT=5432
PG_USER=postgres
PG_PASSWORD=postgres123
PG_DBNAME=social_insight

# Kafka Configuration (for future extensions)
KAFKA_HOST=kafka:29092
KAFKA_TOPIC=raw_posts
```

---

## API Endpoints

### 1. Health Check
```bash
curl http://localhost:8888/api/health

# Response:
# {
#   "status": "ok",
#   "time": "2024-01-28T10:30:00Z"
# }
```

### 2. Overall Statistics
```bash
curl http://localhost:8888/api/stats

# Response:
# {
#   "total_posts": 5432,
#   "by_topic": {
#     "ai": 1200,
#     "cloud": 800,
#     "devops": 950,
#     "programming": 1100,
#     "startup": 1382
#   },
#   "by_sentiment": {
#     "positive": 3200,
#     "negative": 800,
#     "neutral": 1432
#   }
# }
```

### 3. Statistics by Topic
```bash
curl http://localhost:8888/api/topics

# Response:
# {
#   "ai": {"count": 1200, "avg_likes": 45, "avg_comments": 12},
#   "cloud": {"count": 800, "avg_likes": 38, "avg_comments": 10},
#   ...
# }
```

### 4. Statistics by Sentiment
```bash
curl http://localhost:8888/api/sentiment

# Response:
# {
#   "positive": {"count": 3200, "posts": [...]},
#   "negative": {"count": 800, "posts": [...]},
#   "neutral": {"count": 1432, "posts": [...]}
# }
```

### 5. Top Authors
```bash
curl http://localhost:8888/api/authors

# Response:
# [
#   {"author": "John Doe", "posts": 45, "avg_likes": 120},
#   {"author": "Jane Smith", "posts": 38, "avg_likes": 95},
#   ...
# ]
```

### 6. Recent Posts
```bash
curl http://localhost:8888/api/recent

# Response:
# [
#   {
#     "id": "abc123",
#     "author": "John Doe",
#     "title": "AI Trends 2024",
#     "content": "...",
#     "topic": "ai",
#     "sentiment": "positive",
#     "likes": 120,
#     "created_at": "2024-01-28T10:00:00Z"
#   },
#   ...
# ]
```

---

## Dashboard Features

### Live Statistics
- Total posts count
- Posts by topic
- Posts by sentiment
- Real-time updates (refreshed from API)

### Charts & Visualization
- **Topic Distribution**: Pie/Doughnut chart
- **Sentiment Analysis**: Bar chart
- **Trends Over Time**: Line chart
- **Author Leaderboard**: Table with rankings

### Data Sources
- **Real-time Cache**: Redis (for recent posts)
- **Historical Data**: PostgreSQL (for statistics)

### Performance
- Caching strategies to reduce database load
- Pagination on large result sets
- Responsive design for mobile/desktop

---

## Commands

### Build & Start
```bash
# Build image
docker-compose build

# Start service
docker-compose up -d

# Start with logs
docker-compose up
```

### View Logs
```bash
# All logs
docker-compose logs -f

# Specific service
docker-compose logs -f api

# Follow and tail
docker-compose logs -f --tail=100
```

### Stop Service
```bash
# Stop (containers remain)
docker-compose stop

# Stop and remove
docker-compose down

# Remove everything
docker-compose down -v
```

### Restart
```bash
docker-compose restart api
```

### Rebuild Images
```bash
# Build without cache
docker-compose build --no-cache

# Full rebuild
docker-compose up -d --build
```

---

## Testing

### Test Connectivity
```bash
# API is responding
curl http://localhost:8888/api/health

# Dashboard loads
curl http://localhost:8888/ | head -20

# API returns data
curl http://localhost:8888/api/stats | jq .
```

### Test with Sample Data
```bash
# Insert test data
psql -h localhost -U postgres -d social_insight << EOF
INSERT INTO posts (id, author, content, topic, sentiment, likes, comments, shares, platform, created_at)
VALUES 
  ('test1', 'Test User', 'Test post about AI', 'ai', 'positive', 10, 2, 1, 'twitter', NOW()),
  ('test2', 'Test User', 'Test post about Cloud', 'cloud', 'neutral', 5, 1, 0, 'linkedin', NOW());
EOF

# Verify in API
curl http://localhost:8888/api/stats | jq '.total_posts'
```

---

## Troubleshooting

### API Server Won't Start
```bash
# Check logs
docker logs api_server

# Verify port is available
netstat -an | findstr 8888

# Check environment variables
docker exec api_server env | grep API_
```

### Cannot Connect to PostgreSQL
```bash
# Verify Data Service is running
docker ps | grep data_postgres

# Test connection from API container
docker exec api_server psql -h postgres -U postgres -d social_insight -c "SELECT 1"

# Check connection string in logs
docker logs api_server | grep PostgreSQL
```

### Cannot Connect to Redis
```bash
# Verify Redis is running
docker ps | grep data_redis

# Test connection from API container
docker exec api_server redis-cli -h redis ping

# Check if Redis is empty
redis-cli DBSIZE
```

### API Returning Empty Data
```bash
# Check if database has data
psql -h localhost -U postgres -d social_insight -c "SELECT COUNT(*) FROM posts;"

# Verify consumer is running and processing
docker logs processing_consumer -f | grep "Saved\|processed"

# Check Redis cache
redis-cli KEYS "*"

# Check API logs
docker logs api_server -f
```

### Dashboard Not Loading
```bash
# Check HTTP response
curl -v http://localhost:8888/

# View browser console for JS errors
# Press F12 in browser → Console tab

# Check file permissions
docker exec api_server ls -la /root/web/
```

---

## Development

### Modify Dashboard
1. Edit `web/index.html`
2. Rebuild Docker image:
   ```bash
   docker-compose build --no-cache
   ```
3. Restart container:
   ```bash
   docker-compose up -d
   ```
4. Refresh browser: http://localhost:8888

### Add New API Endpoint
1. Edit `cmd/api/main.go`
2. Add handler function
3. Register route:
   ```go
   http.HandleFunc("/api/new-endpoint", enableCORS(server.handleNewEndpoint))
   ```
4. Rebuild:
   ```bash
   docker-compose build --no-cache
   docker-compose up -d
   ```

### Update Frontend to Call New Endpoint
1. Edit `web/index.html`
2. Add fetch call:
   ```javascript
   fetch('/api/new-endpoint')
     .then(r => r.json())
     .then(data => console.log(data))
   ```
3. Refresh browser

---

## Performance Optimization

### Caching Strategy
- **Redis Cache**: Recent posts (TTL: 1 hour)
- **API Response Cache**: Statistics (could be cached)
- **Browser Cache**: Static assets (HTML, CSS, JS)

### Database Optimization
- **Indexes**: Created on frequently queried columns
- **Connection Pooling**: Up to 50 connections
- **Query Optimization**: Use appropriate queries in handlers

### API Response Times
Expected:
- `/api/health`: < 10ms
- `/api/stats`: < 100ms (from cache)
- `/api/recent`: < 50ms (from Redis)
- `/api/topics`: < 200ms (from database)

### Monitor Performance
```bash
# Measure response time
time curl http://localhost:8888/api/stats

# Monitor resource usage
docker stats api_server
```

---

## CORS & Security

### CORS Policy
Currently: `Access-Control-Allow-Origin: *` (open)

For production:
```go
w.Header().Set("Access-Control-Allow-Origin", "https://yourdomain.com")
```

### Security Headers (To Add)
```go
w.Header().Set("X-Content-Type-Options", "nosniff")
w.Header().Set("X-Frame-Options", "DENY")
w.Header().Set("Content-Security-Policy", "default-src 'self'")
```

---

## Deployment

### Docker Image
```bash
# Build
docker build -f Dockerfile.api -t social-insight-api:latest .

# Tag for registry
docker tag social-insight-api:latest myregistry/social-insight-api:latest

# Push
docker push myregistry/social-insight-api:latest
```

### Docker Compose Production
```yaml
api:
  image: myregistry/social-insight-api:latest
  ports:
    - "8888:8888"
  environment:
    REDIS_HOST: redis.example.com
    PG_HOST: postgres.example.com
  restart: always
  # Add health check, logging, etc.
```

---

## Network Details

### Docker Network
- **Name**: `api_network`
- **Type**: Bridge
- **Connects to**: Data Service containers

### Service Communication
```
API Server ──► PostgreSQL (port 5432)
           ──► Redis (port 6379)
           ──► Kafka (port 29092) [future]
```

### External Access
- **Host**: http://localhost:8888
- **Docker Network**: http://api_server:8888

---

## Monitoring

### Container Health
```bash
docker-compose ps

# Expected status: Up (healthy)
```

### API Response Monitoring
```bash
# Check health continuously
while true; do curl -s http://localhost:8888/api/health | jq '.status'; sleep 5; done
```

### Resource Usage
```bash
docker stats api_server --no-stream
```

---

## Logs

### View Recent Logs
```bash
docker-compose logs -f --tail=50
```

### Export Logs
```bash
docker-compose logs > api_logs.txt
```

### Grep for Errors
```bash
docker-compose logs | grep -i error
docker-compose logs | grep -i failed
```

---

## Next Steps

1. ✅ Data Service running
2. ✅ Processing Service running
3. ✅ API Service running
4. 📊 Access Dashboard: http://localhost:8888
5. 🧪 Run end-to-end tests
6. 📈 Monitor data flowing through system

---

**Last Updated**: January 28, 2025  
**Go Version**: 1.21  
**Frontend**: HTML5 + CSS3 + JavaScript (Vanilla)  
**Dependencies**: PostgreSQL, Redis, Kafka (optional)
