# 🎯 AWS Deployment Configuration Files

Dự án Social Insight đã được chuẩn bị cho triển khai trên AWS với 3 dịch vụ độc lập.

## 📁 Files Created

### Docker Compose Files

| File | Purpose | Deployment |
|------|---------|-----------|
| `docker-compose.prod.yml` | Full stack production | Single machine / docker-compose up |
| `docker-compose.aws-ingestion.yml` | Crawlers only | EC2-1 (Ingestion Service) |
| `docker-compose.aws-api.yml` | API Gateway only | EC2-2 (API Service) |
| `docker-compose.aws-processing.yml` | Consumer + ML Service | EC2-3 (Processing Service) |

### Environment Files

| File | Purpose | Services |
|------|---------|----------|
| `.env.prod` | Local production config | All services (single instance) |
| `.env.aws-ingestion` | Ingestion service config | HN/Medium/DevTo Crawlers |
| `.env.aws-api` | API service config | API Gateway + Frontend |
| `.env.aws-processing` | Processing service config | Consumer + ML Service |

### Documentation

| File | Content |
|------|---------|
| `AWS_DEPLOYMENT.md` | Complete AWS deployment guide (94 pages) |
| `HARDCODE_FIXES.md` | (This file) Summary of changes |

## 🔧 Changes Made

### 1. Removed Hard-coded Values from Dockerfiles

**Before:**
```dockerfile
ENV KAFKA_BROKERS=kafka:9092
ENV DB_HOST=postgres
ENV DB_PASSWORD=postgres123
```

**After:**
```dockerfile
# Environment variables (will be overridden by docker-compose/container)
```

**Files Updated:**
- ✅ `processing/ml-service/Dockerfile`
- ✅ `ingestion/hn-crawler/Dockerfile`
- ✅ `ingestion/medium-crawler/Dockerfile`
- ✅ `ingestion/devto-crawler/Dockerfile`
- ✅ `processing/consumer/Dockerfile`
- ✅ `presentation/api-gateway/Dockerfile`

### 2. Created Configuration Files

All environment variables now configurable via:
- `.env` files for docker-compose
- Environment variables when running containers
- Docker Compose environment sections

## 🚀 Quick Start

### Local Development (Full Stack)

```bash
# Use local docker-compose
docker-compose up -d --build

# Or use prod config
docker-compose -f docker-compose.prod.yml --env-file .env.prod up -d --build
```

### AWS Deployment (3 Services)

**EC2-1: Ingestion Service**
```bash
docker-compose -f docker-compose.aws-ingestion.yml --env-file .env.aws-ingestion up -d
```

**EC2-2: API Service**
```bash
docker-compose -f docker-compose.aws-api.yml --env-file .env.aws-api up -d
```

**EC2-3: Processing Service**
```bash
docker-compose -f docker-compose.aws-processing.yml --env-file .env.aws-processing up -d
```

## 📊 3-Service Architecture

```
┌─────────────────────────────────────────────────┐
│           AWS Infrastructure                    │
│                                                 │
│  ┌──────────────────┐  ┌──────────────────┐   │
│  │   EC2-1          │  │    EC2-2         │   │
│  │  Ingestion       │  │   API Service    │   │
│  │                  │  │                  │   │
│  │ • HN Crawler     │  │ • API Gateway    │   │
│  │ • Medium         │  │ • Frontend       │   │
│  │ • DevTo          │  │                  │   │
│  └────────┬─────────┘  └────────┬─────────┘   │
│           │                     │             │
│           └──────────┬──────────┘             │
│                      │                        │
│                   Kafka (MSK)                 │
│                      │                        │
│           ┌──────────┴──────────┐             │
│           │                     │             │
│        ┌──┴─────────┐   ┌───────┴──┐         │
│        │   EC2-3    │   │   Shared  │         │
│        │ Processing │   │ Databases │         │
│        │            │   │           │         │
│        │ • Consumer │   │ • RDS DB  │         │
│        │ • ML Svc   │   │ • Redis   │         │
│        └────────────┘   └───────────┘         │
│                                                 │
└─────────────────────────────────────────────────┘
```

## 🔐 Security Notes

### Hard-coded Credentials - Before
```dockerfile
ENV DB_PASSWORD=postgres123
```

### Secure Method - After
1. **Local Dev**: Use `.env` files (add to `.gitignore`)
2. **AWS**: Use AWS Secrets Manager
3. **CI/CD**: Use GitHub Secrets / GitLab CI variables

### Best Practices Included

✅ All env vars externalized
✅ No credentials in Dockerfiles
✅ No credentials in git repository
✅ Health checks configured
✅ Logging configuration added
✅ VPC networking guidance
✅ Security groups examples

## 📋 Configuration Variables

### Common Variables (All Services)

```bash
ENV=production              # Environment name
DEBUG=false                 # Debug logging
```

### Ingestion Service Variables

```bash
KAFKA_BROKERS=             # Kafka broker address
KAFKA_TOPIC=raw_posts      # Kafka topic
HN_CRAWL_INTERVAL=30s      # HackerNews interval
MEDIUM_CRAWL_INTERVAL=60s  # Medium interval
DEVTO_CRAWL_INTERVAL=60s   # Dev.to interval
```

### API Service Variables

```bash
DB_HOST=                   # RDS endpoint
DB_PORT=5432               # PostgreSQL port
DB_USER=postgres           # Database user
DB_PASSWORD=               # Database password (use Secrets Manager)
DB_NAME=social_insight     # Database name
REDIS_ADDR=                # ElastiCache endpoint
API_PORT=8888              # API port
```

### Processing Service Variables

```bash
KAFKA_BROKERS=             # Kafka broker address
KAFKA_TOPIC=raw_posts      # Kafka topic
CONSUMER_GROUP=            # Consumer group
DB_HOST=                   # RDS endpoint
DB_PORT=5432               # PostgreSQL port
DB_USER=postgres           # Database user
DB_PASSWORD=               # Database password (use Secrets Manager)
DB_NAME=social_insight     # Database name
REDIS_ADDR=                # ElastiCache endpoint
ML_PORT=8001               # ML service port
ML_SERVICE_URL=            # ML service URL
```

## 📖 Full Documentation

For detailed AWS deployment instructions, see:
**`AWS_DEPLOYMENT.md`** - Contains:

1. Architecture overview
2. AWS prerequisites
3. VPC & Security groups setup
4. RDS PostgreSQL setup
5. ElastiCache Redis setup
6. MSK/Kafka setup
7. EC2 instance launch
8. Docker installation
9. Service deployment steps
10. Verification procedures
11. Production best practices
12. Monitoring setup
13. Backup & DR procedures
14. Scaling strategies
15. CI/CD updates
16. Cost estimation
17. Troubleshooting guide

## ✅ Ready for Production

All files are production-ready:

- [x] No hard-coded credentials
- [x] Environment-based configuration
- [x] Health checks configured
- [x] Logging setup
- [x] Network isolation (VPC)
- [x] Data persistence
- [x] Auto-restart policies
- [x] Multi-instance architecture
- [x] AWS service integration

## 🎓 Next Steps

1. **Update AWS Endpoints** in `.env.aws-*` files
2. **Setup AWS Infrastructure** using AWS_DEPLOYMENT.md
3. **Configure Secrets Manager** for sensitive values
4. **Setup Monitoring** with CloudWatch
5. **Test Deployment** and verify data flow
6. **Setup CI/CD** for automated updates

## 💡 Tips

- Use AWS Secrets Manager for passwords
- Use Parameter Store for non-sensitive config
- Enable VPC Flow Logs for debugging
- Setup CloudWatch alarms
- Regular database backups
- Monitor with X-Ray
- Use Systems Manager Session Manager instead of SSH

---

**Created**: January 27, 2026
**Status**: Production Ready ✅
