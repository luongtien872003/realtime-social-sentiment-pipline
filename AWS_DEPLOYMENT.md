# =====================================================
# AWS DEPLOYMENT GUIDE - 3 Service Architecture
# =====================================================
# Deploy Social Insight on AWS with 3 separate EC2 instances
# =====================================================

## 📋 Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                         AWS VPC                             │
│  ┌──────────────────────────────────────────────────────┐  │
│  │                                                      │  │
│  │  EC2-1: Ingestion Service      EC2-2: API Service  │  │
│  │  ├─ HN Crawler                 ├─ API Gateway      │  │
│  │  ├─ Medium Crawler             └─ Frontend         │  │
│  │  └─ DevTo Crawler                                   │  │
│  │          │                             │             │  │
│  │          └──────────────┬──────────────┘             │  │
│  │                         │                            │  │
│  │                    Kafka (MSK)                       │  │
│  │                         │                            │  │
│  │         EC2-3: Processing Service                   │  │
│  │         ├─ Consumer                                 │  │
│  │         └─ ML Service                               │  │
│  │                │                                     │  │
│  │        ┌───────┼────────┐                           │  │
│  │        ▼       ▼        ▼                           │  │
│  │    PostgreSQL (RDS)  Redis (ElastiCache)           │  │
│  │                                                      │  │
│  └──────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

## 🚀 Prerequisites

### AWS Services Required:
- **EC2**: 3 instances (t3.medium or larger)
- **RDS**: PostgreSQL 15 (db.t3.micro or larger)
- **ElastiCache**: Redis 7 (cache.t3.micro or larger)
- **MSK or Self-managed Kafka**: For message queue
- **VPC & Security Groups**: For network communication

### Local Requirements:
- Docker & Docker Compose
- AWS CLI v2
- SSH Key Pair for EC2 access

## 1️⃣ Setup AWS Infrastructure

### Step 1: Create VPC & Security Groups

```bash
# Create VPC (or use default)
aws ec2 create-vpc --cidr-block 10.0.0.0/16

# Create Security Groups
aws ec2 create-security-group \
  --group-name social-insight-sg \
  --description "Social Insight application security group"

# Allow inter-service communication (internal)
aws ec2 authorize-security-group-ingress \
  --group-id sg-xxxxx \
  --protocol tcp \
  --port 0-65535 \
  --cidr 10.0.0.0/16

# Allow SSH (for management)
aws ec2 authorize-security-group-ingress \
  --group-id sg-xxxxx \
  --protocol tcp \
  --port 22 \
  --cidr YOUR_IP/32

# Allow HTTP/HTTPS (for API)
aws ec2 authorize-security-group-ingress \
  --group-id sg-xxxxx \
  --protocol tcp \
  --port 80 \
  --cidr 0.0.0.0/0

aws ec2 authorize-security-group-ingress \
  --group-id sg-xxxxx \
  --protocol tcp \
  --port 8888 \
  --cidr 0.0.0.0/0
```

### Step 2: Create RDS PostgreSQL

```bash
aws rds create-db-instance \
  --db-instance-identifier social-insight-db \
  --db-instance-class db.t3.micro \
  --engine postgres \
  --engine-version 15.3 \
  --master-username postgres \
  --master-user-password "YOUR_STRONG_PASSWORD" \
  --allocated-storage 20 \
  --storage-type gp3 \
  --db-name social_insight \
  --vpc-security-group-ids sg-xxxxx \
  --publicly-accessible false
```

Wait for RDS to be available:
```bash
aws rds describe-db-instances \
  --db-instance-identifier social-insight-db \
  --query 'DBInstances[0].DBInstanceStatus'
```

Get RDS Endpoint:
```bash
aws rds describe-db-instances \
  --db-instance-identifier social-insight-db \
  --query 'DBInstances[0].Endpoint.Address'
```

### Step 3: Create ElastiCache Redis

```bash
aws elasticache create-cache-cluster \
  --cache-cluster-id social-insight-redis \
  --cache-node-type cache.t3.micro \
  --engine redis \
  --engine-version 7.0 \
  --num-cache-nodes 1 \
  --security-group-ids sg-xxxxx \
  --auto-failover-enabled
```

Get Redis Endpoint:
```bash
aws elasticache describe-cache-clusters \
  --cache-cluster-id social-insight-redis \
  --show-cache-node-info \
  --query 'CacheClusters[0].CacheNodes[0].Endpoint'
```

### Step 4: Create MSK (Managed Streaming Kafka) or Self-managed Kafka

**Option A: AWS MSK (Recommended)**

```bash
aws kafka create-cluster \
  --cluster-name social-insight-kafka \
  --broker-node-group-info \
    BrokerAZDistribution=DEFAULT,\
    ClientSubnets=subnet-xxxxx,subnet-yyyyy,\
    SecurityGroups=sg-xxxxx,\
    InstanceType=kafka.t3.small,\
    StorageInfo="{EbsStorageInfo={VolumeSize=100}}" \
  --kafka-version 3.4.0 \
  --number-of-broker-nodes 3
```

**Option B: Self-managed EC2 Kafka**
- Use a separate t3.small EC2 instance
- Install Kafka using standard Docker setup
- Expose port 9092 to VPC

## 2️⃣ Launch EC2 Instances

### Launch 3 Instances:

```bash
# Instance 1: Ingestion Service
aws ec2 run-instances \
  --image-id ami-0c55b159cbfafe1f0 \
  --instance-type t3.medium \
  --key-name your-key-pair \
  --security-group-ids sg-xxxxx \
  --tag-specifications 'ResourceType=instance,Tags=[{Key=Name,Value=social-insight-ingestion}]'

# Instance 2: API Service
aws ec2 run-instances \
  --image-id ami-0c55b159cbfafe1f0 \
  --instance-type t3.medium \
  --key-name your-key-pair \
  --security-group-ids sg-xxxxx \
  --tag-specifications 'ResourceType=instance,Tags=[{Key=Name,Value=social-insight-api}]'

# Instance 3: Processing Service
aws ec2 run-instances \
  --image-id ami-0c55b159cbfafe1f0 \
  --instance-type t3.medium \
  --key-name your-key-pair \
  --security-group-ids sg-xxxxx \
  --tag-specifications 'ResourceType=instance,Tags=[{Key=Name,Value=social-insight-processing}]'
```

Get Instance IPs:
```bash
aws ec2 describe-instances \
  --filters "Name=tag:Name,Values=social-insight-*" \
  --query 'Reservations[*].Instances[*].[Tags[?Key==`Name`].Value[0],PublicIpAddress]'
```

## 3️⃣ Setup Docker on EC2 Instances

SSH to each instance and run:

```bash
#!/bin/bash
# Update system
sudo yum update -y

# Install Docker
sudo amazon-linux-extras install docker -y
sudo systemctl start docker
sudo systemctl enable docker
sudo usermod -a -G docker ec2-user

# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" \
  -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Install Git
sudo yum install git -y

# Verify
docker --version
docker-compose --version
```

## 4️⃣ Deploy Services

### On EC2-1 (Ingestion Service):

```bash
# SSH to instance
ssh -i your-key.pem ec2-user@INGESTION_IP

# Clone repository
git clone https://github.com/your-repo/social-insight.git
cd social-insight

# Create and edit .env.aws-ingestion
cat > .env.aws-ingestion << 'EOF'
KAFKA_BROKERS=kafka.xxxxx.kafka.us-east-1.amazonaws.com:9092
KAFKA_TOPIC=raw_posts
HN_CRAWL_INTERVAL=30s
MEDIUM_CRAWL_INTERVAL=60s
DEVTO_CRAWL_INTERVAL=60s
ENV=production
EOF

# Start crawlers
docker-compose -f docker-compose.aws-ingestion.yml --env-file .env.aws-ingestion up -d

# View logs
docker-compose -f docker-compose.aws-ingestion.yml logs -f
```

### On EC2-2 (API Service):

```bash
# SSH to instance
ssh -i your-key.pem ec2-user@API_IP

# Clone repository
git clone https://github.com/your-repo/social-insight.git
cd social-insight

# Create and edit .env.aws-api
cat > .env.aws-api << 'EOF'
DB_HOST=social-insight.xxxxx.us-east-1.rds.amazonaws.com
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=YOUR_STRONG_PASSWORD
DB_NAME=social_insight
REDIS_ADDR=social-insight.xxxxx.ng.0001.use1.cache.amazonaws.com:6379
API_PORT=8888
ENV=production
EOF

# Start API Gateway
docker-compose -f docker-compose.aws-api.yml --env-file .env.aws-api up -d

# View logs
docker-compose -f docker-compose.aws-api.yml logs -f

# Access application
# http://API_IP:8888
```

### On EC2-3 (Processing Service):

```bash
# SSH to instance
ssh -i your-key.pem ec2-user@PROCESSING_IP

# Clone repository
git clone https://github.com/your-repo/social-insight.git
cd social-insight

# Create and edit .env.aws-processing
cat > .env.aws-processing << 'EOF'
KAFKA_BROKERS=kafka.xxxxx.kafka.us-east-1.amazonaws.com:9092
KAFKA_TOPIC=raw_posts
CONSUMER_GROUP=social_insight_processor
DB_HOST=social-insight.xxxxx.us-east-1.rds.amazonaws.com
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=YOUR_STRONG_PASSWORD
DB_NAME=social_insight
REDIS_ADDR=social-insight.xxxxx.ng.0001.use1.cache.amazonaws.com:6379
ML_PORT=8001
ML_SERVICE_URL=http://localhost:8001
ENV=production
EOF

# Start Consumer + ML Service
docker-compose -f docker-compose.aws-processing.yml --env-file .env.aws-processing up -d

# View logs
docker-compose -f docker-compose.aws-processing.yml logs -f
```

## 5️⃣ Verify Deployment

### Check Service Status:

```bash
# On each instance
docker-compose ps
docker-compose logs --tail 50

# Check connectivity
curl -s http://API_IP:8888/api/health | jq
curl -s http://API_IP:8888/api/stats | jq
```

### Check Data Flow:

```bash
# On Processing instance, check if Consumer is processing
docker-compose logs consumer -f

# Check database for posts
psql -h social-insight.xxxxx.us-east-1.rds.amazonaws.com \
  -U postgres \
  -d social_insight \
  -c "SELECT COUNT(*) FROM posts;"
```

## 6️⃣ Production Best Practices

### Security:
- [ ] Use AWS Secrets Manager for passwords
- [ ] Enable VPC encryption
- [ ] Use IAM roles for EC2 instances
- [ ] Enable RDS encryption at rest & transit
- [ ] Use VPC endpoints for AWS service access

### Monitoring:
- [ ] Setup CloudWatch alarms for EC2 CPU/Memory
- [ ] Setup RDS Performance Insights
- [ ] Setup ElastiCache metrics
- [ ] Enable VPC Flow Logs
- [ ] Setup log aggregation with CloudWatch Logs

### Backup & Disaster Recovery:
- [ ] Enable RDS automated backups (7 days retention)
- [ ] Setup cross-region replication for RDS
- [ ] Regular database snapshots
- [ ] Test recovery procedures

### Scaling:
- [ ] Setup Auto Scaling Groups for EC2
- [ ] Use Application Load Balancer for API
- [ ] Horizontal scaling for crawlers
- [ ] RDS read replicas for read-heavy workloads

## 7️⃣ Update Services (CI/CD)

### Update script for each service:

```bash
#!/bin/bash
# File: update-service.sh

SERVICE=$1  # ingestion, api, or processing
ENV_FILE=".env.aws-${SERVICE}"
COMPOSE_FILE="docker-compose.aws-${SERVICE}.yml"

# Pull latest code
git pull origin main

# Rebuild images
docker-compose -f $COMPOSE_FILE --env-file $ENV_FILE build --no-cache

# Restart services
docker-compose -f $COMPOSE_FILE --env-file $ENV_FILE down
docker-compose -f $COMPOSE_FILE --env-file $ENV_FILE up -d

# Verify
docker-compose -f $COMPOSE_FILE ps
```

Usage:
```bash
chmod +x update-service.sh
./update-service.sh ingestion
./update-service.sh api
./update-service.sh processing
```

## 📊 Cost Estimation (Monthly)

| Service | Instance | Cost |
|---------|----------|------|
| EC2 (3x t3.medium) | On-demand | ~$90 |
| RDS PostgreSQL | db.t3.micro | ~$35 |
| ElastiCache Redis | cache.t3.micro | ~$30 |
| MSK Kafka | 3 brokers (t3.small) | ~$150 |
| Data Transfer | ~ 100GB/month | ~$10 |
| **Total** | | **~$315/month** |

## 🆘 Troubleshooting

### Crawlers not sending data:
```bash
# Check Kafka connectivity
docker exec social_insight_hn_crawler nc -zv kafka.xxxxx.amazonaws.com 9092

# Check logs
docker logs social_insight_hn_crawler
```

### API not showing data:
```bash
# Check database connection
docker exec social_insight_consumer psql -h RDS_HOST -U postgres -d social_insight -c "SELECT COUNT(*) FROM posts;"

# Check Redis
docker exec social_insight_api redis-cli -h REDIS_HOST ping
```

### Consumer not processing:
```bash
# Check Kafka topics
docker exec social_insight_kafka kafka-topics --bootstrap-server localhost:9092 --list

# Check consumer groups
docker exec social_insight_kafka kafka-consumer-groups --bootstrap-server localhost:9092 --list

# Check lag
docker exec social_insight_kafka kafka-consumer-groups \
  --bootstrap-server localhost:9092 \
  --group social_insight_processor \
  --describe
```

## 📚 References

- [AWS EC2 Documentation](https://docs.aws.amazon.com/ec2/)
- [AWS RDS Documentation](https://docs.aws.amazon.com/rds/)
- [AWS ElastiCache Documentation](https://docs.aws.amazon.com/elasticache/)
- [AWS MSK Documentation](https://docs.aws.amazon.com/msk/)
- [Docker Documentation](https://docs.docker.com/)
