# 🌐 API Service - Presentation Layer

[![Go](https://img.shields.io/badge/Go-1.21-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Chart.js](https://img.shields.io/badge/Chart.js-Dashboard-FF6384?style=flat&logo=chartdotjs)](https://www.chartjs.org/)

**Components**: REST API Server, Web Dashboard  
**AWS EC2**: `social-insight-api` (13.214.56.2)  
**Status**: ✅ Production Running

---

## 🌐 Live Demo

> 🔗 **Dashboard**: [http://13.214.56.2:8888](http://13.214.56.2:8888)

---

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────┐
│                   API Service                       │
│                                                     │
│  ┌─────────────────────────────────────────────┐   │
│  │              REST API Server                │   │
│  │                 (:8888)                     │   │
│  │                                             │   │
│  │  /api/health    /api/stats    /api/recent  │   │
│  │  /api/topics    /api/sentiment             │   │
│  │  /api/trending  /api/compare               │   │
│  └─────────────────────────────────────────────┘   │
│                        │                            │
│                        ▼                            │
│  ┌─────────────────────────────────────────────┐   │
│  │            Web Dashboard                    │   │
│  │         (HTML + CSS + Chart.js)             │   │
│  │                                             │   │
│  │  📊 Topic Trends    📈 Sentiment Analysis  │   │
│  │  👤 Top Authors     📝 Recent Posts        │   │
│  └─────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────┘
           │                    │
           ▼                    ▼
    ┌────────────┐       ┌────────────┐
    │ PostgreSQL │       │   Redis    │
    │   (read)   │       │  (cache)   │
    └────────────┘       └────────────┘
```

---

## 🚀 Quick Start

```bash
# Ensure Data Service is running
cd ../data-service && docker-compose ps

# Start API Service
docker-compose up -d --build

# Access Dashboard
open http://localhost:8888
```

---

## ☁️ AWS Production

| Property | Value |
|----------|-------|
| **EC2 Instance** | social-insight-api |
| **Public IP** | 13.214.56.2 |
| **Private IP** | 172.31.8.70 |
| **Instance Type** | t3.small |
| **Security Group** | api-sg |

### Production URLs
| Endpoint | URL |
|----------|-----|
| Dashboard | http://13.214.56.2:8888 |
| API Health | http://13.214.56.2:8888/api/health |
| API Stats | http://13.214.56.2:8888/api/stats |

### Running Containers
```
api_server - REST API :8888 (Up)
```

---

## 📡 API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/health` | Health check |
| GET | `/api/stats` | Overall statistics |
| GET | `/api/recent` | Recent posts (limit 50) |
| GET | `/api/topics` | Topic distribution |
| GET | `/api/sentiment` | Sentiment analysis |
| GET | `/api/trending` | Top trending posts |
| GET | `/api/compare` | Today vs Yesterday |

### Examples

```bash
# Health Check
curl http://localhost:8888/api/health
# {"status":"ok","time":"2026-01-31T10:00:00Z"}

# Statistics
curl http://localhost:8888/api/stats | jq .
# {
#   "total_posts": 5432,
#   "by_topic": {"ai": 1200, "cloud": 800, ...},
#   "by_sentiment": {"positive": 3200, "negative": 800, ...}
# }

# Recent Posts
curl http://localhost:8888/api/recent | jq '.[0]'
# {
#   "id": "abc123",
#   "author": "John Doe",
#   "content": "AI Trends 2024...",
#   "topic": "ai",
#   "sentiment": "positive",
#   "likes": 120
# }
```

---

## 📊 Dashboard Features

| Feature | Description |
|---------|-------------|
| **Topic Distribution** | Pie chart of posts by topic |
| **Sentiment Analysis** | Bar chart (positive/neutral/negative) |
| **Recent Posts** | Live feed of latest posts |
| **Top Authors** | Leaderboard by post count |
| **Statistics Cards** | Total posts, today's count |

---

## ⚙️ Environment Variables

```env
# API Server
API_PORT=:8888

# Redis (from Data Service)
REDIS_ADDR=redis:6379           # Local
REDIS_ADDR=172.31.16.144:6379   # AWS

# PostgreSQL (from Data Service)
PG_HOST=postgres                # Local
PG_HOST=172.31.16.144           # AWS
PG_PORT=5432
PG_USER=postgres
PG_PASSWORD=postgres123
PG_DBNAME=social_insight
```

---

## 🔧 Common Commands

```bash
# View logs
docker logs api_server -f

# Restart
docker-compose restart

# Rebuild after code changes
docker-compose up -d --build

# Test API
curl http://localhost:8888/api/health
```

---

## 🛠️ Development

### Modify Dashboard
1. Edit `web/index.html`
2. Rebuild: `docker-compose up -d --build`
3. Refresh browser

### Add New Endpoint
1. Edit `cmd/api/main.go`
2. Add handler function
3. Register route
4. Rebuild and test

---

## 📈 Performance

| Endpoint | Response Time |
|----------|---------------|
| `/api/health` | < 10ms |
| `/api/stats` | < 100ms |
| `/api/recent` | < 50ms |

---

> **Last Updated**: January 31, 2026
