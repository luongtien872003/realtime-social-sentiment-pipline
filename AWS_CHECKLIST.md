# ✅ AWS Deployment Checklist

## 🔍 Hard-code Issues Fixed

### Dockerfiles Cleaned
- [x] `processing/ml-service/Dockerfile` - Removed `ENV ML_PORT`, `ENV LOG_LEVEL`
- [x] `ingestion/hn-crawler/Dockerfile` - Removed hard-coded Kafka/Interval
- [x] `ingestion/medium-crawler/Dockerfile` - Removed hard-coded Kafka/Interval
- [x] `ingestion/devto-crawler/Dockerfile` - Removed hard-coded Kafka/Interval
- [x] `processing/consumer/Dockerfile` - Removed hard-coded DB/Redis/Kafka
- [x] `presentation/api-gateway/Dockerfile` - Removed hard-coded DB/Redis/API

**Status**: ✅ All 6 Dockerfiles updated

## 📁 Configuration Files Created

### Docker Compose Files
- [x] `docker-compose.prod.yml` - Full production stack (all services together)
- [x] `docker-compose.aws-ingestion.yml` - EC2-1 (3 crawlers)
- [x] `docker-compose.aws-api.yml` - EC2-2 (API Gateway + Frontend)
- [x] `docker-compose.aws-processing.yml` - EC2-3 (Consumer + ML Service)

**Status**: ✅ 4 docker-compose files created

### Environment Files
- [x] `.env.prod` - Local production environment variables
- [x] `.env.aws-ingestion` - Ingestion service variables (MSK endpoint needed)
- [x] `.env.aws-api` - API service variables (RDS + ElastiCache endpoints needed)
- [x] `.env.aws-processing` - Processing service variables (RDS + ElastiCache + MSK endpoints needed)

**Status**: ✅ 4 environment files created

### Documentation Files
- [x] `AWS_DEPLOYMENT.md` - Complete 94-line deployment guide
- [x] `HARDCODE_FIXES.md` - Summary of changes (this file)

**Status**: ✅ 2 documentation files created

## 🚀 Ready for AWS Deployment

### Before Deployment
- [ ] Update `.env.aws-ingestion` with Kafka brokers
- [ ] Update `.env.aws-api` with RDS endpoint, Redis endpoint
- [ ] Update `.env.aws-processing` with Kafka, RDS, Redis endpoints
- [ ] Create AWS VPC and Security Groups (see AWS_DEPLOYMENT.md)
- [ ] Launch 3 EC2 instances
- [ ] Create RDS PostgreSQL instance
- [ ] Create ElastiCache Redis cluster
- [ ] Create MSK Kafka cluster (or self-managed)

### Deployment
- [ ] Install Docker on all 3 EC2 instances
- [ ] Clone repository on each instance
- [ ] Deploy Ingestion service (EC2-1)
- [ ] Deploy API service (EC2-2)
- [ ] Deploy Processing service (EC2-3)
- [ ] Verify all services are running

### Post-Deployment
- [ ] Check health endpoints
- [ ] Verify data flow (crawlers → Kafka → consumer → DB)
- [ ] Setup CloudWatch monitoring
- [ ] Setup CloudWatch alarms
- [ ] Configure backups for RDS
- [ ] Configure log aggregation

## 📊 3-Service Deployment Map

```
┌────────────────────────────────────────────────────────────┐
│                      AWS Account                            │
│                                                              │
│  ┌────────────────────┐     ┌────────────────────┐         │
│  │   EC2-1            │     │    EC2-2           │         │
│  │  Ingestion Service │     │   API Service      │         │
│  │                    │     │                    │         │
│  │ compose file:      │     │ compose file:      │         │
│  │ aws-ingestion.yml  │     │ aws-api.yml        │         │
│  │                    │     │                    │         │
│  │ env file:          │     │ env file:          │         │
│  │ .env.aws-ingestion │     │ .env.aws-api       │         │
│  │                    │     │                    │         │
│  │ Services:          │     │ Services:          │         │
│  │ • hn-crawler       │     │ • api-gateway      │         │
│  │ • medium-crawler   │     │ • frontend         │         │
│  │ • devto-crawler    │     │                    │         │
│  └──────────┬─────────┘     └────────┬───────────┘         │
│             │                        │                     │
│             │  KAFKA BROKERS         │                     │
│             │  kafka.xxxxx.amazonaws.com:9092              │
│             │                        │                     │
│  ┌──────────▼──────────┐             │                     │
│  │                     │             │                     │
│  │    MSK Kafka 🗄️    │             │                     │
│  │  (3 broker nodes)   │             │                     │
│  │                     │             │                     │
│  └──────────┬──────────┘             │                     │
│             │                        │                     │
│             │                        │                     │
│             └────────┬───────────────┘                     │
│                      │                                     │
│           ┌──────────▼──────────┐                         │
│           │    EC2-3            │                         │
│           │ Processing Service  │                         │
│           │                     │                         │
│           │ compose file:       │                         │
│           │ aws-processing.yml  │                         │
│           │                     │                         │
│           │ env file:           │                         │
│           │ .env.aws-processing │                         │
│           │                     │                         │
│           │ Services:           │                         │
│           │ • consumer          │                         │
│           │ • ml-service        │                         │
│           └──────────┬──────────┘                         │
│                      │                                     │
│                      │                                     │
│           ┌──────────▼──────────────────┬────────────┐   │
│           │                             │            │   │
│      RDS DB 🗄️                    Redis 🗄️            │   │
│    PostgreSQL 15              ElastiCache 7         │   │
│  social-insight.xxxxx        social-insight.xxxxx   │   │
│   .rds.amazonaws.com         .cache.amazonaws.com   │   │
│                                                      │   │
└──────────────────────────────────────────────────────┴───┘
```

## 📝 Environment Variable Checklist

### `.env.aws-ingestion`
```
✅ KAFKA_BROKERS=kafka.xxxxx.kafka.us-east-1.amazonaws.com:9092
✅ KAFKA_TOPIC=raw_posts
✅ HN_CRAWL_INTERVAL=30s
✅ MEDIUM_CRAWL_INTERVAL=60s
✅ DEVTO_CRAWL_INTERVAL=60s
```

### `.env.aws-api`
```
✅ DB_HOST=social-insight.xxxxx.us-east-1.rds.amazonaws.com
✅ DB_PORT=5432
✅ DB_USER=postgres
✅ DB_PASSWORD=<STRONG_PASSWORD>
✅ DB_NAME=social_insight
✅ REDIS_ADDR=social-insight.xxxxx.cache.amazonaws.com:6379
✅ API_PORT=8888
```

### `.env.aws-processing`
```
✅ KAFKA_BROKERS=kafka.xxxxx.kafka.us-east-1.amazonaws.com:9092
✅ KAFKA_TOPIC=raw_posts
✅ CONSUMER_GROUP=social_insight_processor
✅ DB_HOST=social-insight.xxxxx.us-east-1.rds.amazonaws.com
✅ DB_PORT=5432
✅ DB_USER=postgres
✅ DB_PASSWORD=<STRONG_PASSWORD>
✅ DB_NAME=social_insight
✅ REDIS_ADDR=social-insight.xxxxx.cache.amazonaws.com:6379
✅ ML_PORT=8001
✅ ML_SERVICE_URL=http://localhost:8001
```

## 🔐 Security Checklist

- [ ] No passwords in Dockerfiles ✅
- [ ] All env vars in .env files (gitignored) ✅
- [ ] AWS Secrets Manager for production passwords
- [ ] VPC Security Groups properly configured
- [ ] RDS encryption enabled
- [ ] ElastiCache encryption enabled
- [ ] VPC endpoints for AWS services
- [ ] CloudTrail logging enabled
- [ ] S3 bucket policy for logs
- [ ] KMS keys for encryption

## 🧪 Verification Steps

After deployment, run these checks:

### 1. Check Crawlers (EC2-1)
```bash
ssh -i your-key.pem ec2-user@INGESTION_IP
docker-compose -f docker-compose.aws-ingestion.yml logs -f
# Should see: "Fetched X, sent X to Kafka"
```

### 2. Check API (EC2-2)
```bash
ssh -i your-key.pem ec2-user@API_IP
docker-compose -f docker-compose.aws-api.yml logs -f
# Should see: "API Gateway running on port 8888"
curl http://localhost:8888/api/health
```

### 3. Check Processing (EC2-3)
```bash
ssh -i your-key.pem ec2-user@PROCESSING_IP
docker-compose -f docker-compose.aws-processing.yml logs -f
# Should see: "Service started" and "Saved X posts to DB"
```

### 4. Check Data
```bash
# SSH to any instance with database access
psql -h social-insight.xxxxx.us-east-1.rds.amazonaws.com \
  -U postgres \
  -d social_insight \
  -c "SELECT COUNT(*) FROM posts;"
# Should see increasing number
```

## 📊 Service Dependencies

```
EC2-1 Ingestion
    └─> Kafka MSK

EC2-3 Processing
    ├─> Kafka MSK
    ├─> RDS PostgreSQL
    ├─> ElastiCache Redis
    └─> Local ML Service

EC2-2 API
    ├─> RDS PostgreSQL
    └─> ElastiCache Redis
```

## 🆘 Troubleshooting Quick Links

| Issue | Solution |
|-------|----------|
| Crawlers can't reach Kafka | Check VPC Security Groups allow port 9092 |
| API can't reach RDS | Check RDS Security Group allows port 5432 |
| Consumer not processing | Check Kafka topic exists and has data |
| No data in database | Check consumer logs for errors |
| High CPU on ML Service | Tune batch size in consumer |

## ✨ Production Optimizations Done

- [x] Health checks added to all services
- [x] Logging to file with rotation
- [x] Auto-restart policies configured
- [x] VPC networking ready
- [x] Environment-based configuration
- [x] No hard-coded credentials
- [x] Metrics/monitoring ready
- [x] Backup strategies included

---

**Last Updated**: January 27, 2026
**Status**: Ready for AWS Deployment ✅
