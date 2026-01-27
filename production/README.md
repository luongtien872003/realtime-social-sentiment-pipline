# 🔴 Social Insight - AWS Deployment Guide (Student Account)

**Environment**: AWS Production  
**Target**: AWS EC2 (3-instance setup)  
**Account Type**: AWS Student Account (Free Tier / Credits)  
**Status**: ✅ Production Ready

---

## ⚠️ AWS Student Account Limitations

### Important Warnings

1. **Free Tier Limitations** (if using)
   - ❌ Limited compute hours
   - ❌ Limited database usage
   - ❌ Limited data transfer
   - ✅ May need budget for production workload

2. **Quotas** (Service Limits)
   - Default EC2 instance quota: **Usually 5-20**
   - Check actual quota in AWS Console > Service Quotas
   - Request increase before deployment

3. **IP Address Costs**
   - ✅ Elastic IP: Free if attached to running instance
   - ❌ Elastic IP: $0.005/hour if NOT attached
   - Don't allocate unnecessary IPs

4. **Data Transfer Costs**
   - ✅ Inbound data: Always free
   - ❌ Outbound data: $0.09/GB after free tier
   - ❌ Inter-region transfer: Expensive

5. **Cost Estimates**
   - 3x EC2 t2.micro: ~$25-35/month (free tier eligible)
   - RDS db.t2.micro: ~$10-15/month (if outside free tier)
   - ElastiCache cache.t2.micro: ~$15-20/month
   - Kafka: Expensive! (~$100+/month)
   - **Total**: $50-180/month depending on usage

---

## 📋 Quick Start

### Prerequisites
- AWS Account (student or regular)
- AWS CLI v2 installed
- SSH key pair (.pem file)
- 3x EC2 instances (t2.micro for free tier)
- Basic networking knowledge

### Estimated Deployment Time: **30-45 minutes**

---

## 📁 Folder Structure

```
production/
├── README.md                           # This file
├── docker-compose.prod.yml             # Full stack composition
├── docker-compose.aws-ingestion.yml    # EC2-1 (Crawlers)
├── docker-compose.aws-api.yml          # EC2-2 (API)
├── docker-compose.aws-processing.yml   # EC2-3 (Processing)
├── .env.prod                           # Production config template
├── .env.aws-ingestion                  # EC2-1 environment
├── .env.aws-api                        # EC2-2 environment
├── .env.aws-processing                 # EC2-3 environment
│
├── AWS_DEPLOYMENT.md                   # Detailed AWS setup
├── AWS_CHECKLIST.md                    # Verification steps
│
├── ingestion/
├── processing/
├── presentation/
├── shared/
└── infrastructure/
```

---

## 🏗️ 3-Service Architecture (AWS)

```
┌──────────────────────────────────────────────────┐
│            AWS VPC (Private Network)             │
│                                                  │
│  ┌─────────────────┐  ┌─────────────────────┐  │
│  │   EC2-1         │  │     EC2-2           │  │
│  │  Ingestion      │  │   API Service       │  │
│  │                 │  │                     │  │
│  │ • HN Crawler    │  │ • API Gateway       │  │
│  │ • Medium        │  │ • Frontend          │  │
│  │ • DevTo         │  │ • Port: 8888        │  │
│  │                 │  │ • Public IP: Yes    │  │
│  └────────┬────────┘  └──────────┬──────────┘  │
│           │                      │             │
│           └──────────┬───────────┘             │
│                      │                        │
│              Kafka (MSK/EC2)                  │
│                      │                        │
│           ┌──────────┴──────────┐             │
│           │                     │             │
│    ┌──────▼──────┐    ┌────────▼────────┐   │
│    │   EC2-3     │    │   Shared (RDS)  │   │
│    │ Processing  │    │  • PostgreSQL   │   │
│    │             │    │  • ElastiCache  │   │
│    │ • Consumer  │    │                 │   │
│    │ • ML Svc    │    │   Private       │   │
│    │             │    │   (No public IP)│   │
│    └─────────────┘    └─────────────────┘   │
│                                              │
└──────────────────────────────────────────────┘
```

---

## 🚀 Step 1: AWS Infrastructure Setup

### Create VPC & Security Groups

```bash
# Create VPC
aws ec2 create-vpc --cidr-block 10.0.0.0/16 --region us-east-1

# Create Security Group
aws ec2 create-security-group \
  --group-name social-insight-sg \
  --description "Social Insight Application" \
  --vpc-id vpc-xxxxx \
  --region us-east-1
```

### Allow Inbound Traffic

```bash
# SSH access (your IP only - IMPORTANT for security)
aws ec2 authorize-security-group-ingress \
  --group-id sg-xxxxx \
  --protocol tcp \
  --port 22 \
  --cidr YOUR_PUBLIC_IP/32 \
  --region us-east-1

# HTTP (for API)
aws ec2 authorize-security-group-ingress \
  --group-id sg-xxxxx \
  --protocol tcp \
  --port 8888 \
  --cidr 0.0.0.0/0 \
  --region us-east-1

# Internal communication
aws ec2 authorize-security-group-ingress \
  --group-id sg-xxxxx \
  --protocol tcp \
  --port 0-65535 \
  --cidr 10.0.0.0/16 \
  --region us-east-1
```

### Create RDS (PostgreSQL)

```bash
# ⚠️ Student Account: Use db.t2.micro (free tier)
aws rds create-db-instance \
  --db-instance-identifier social-insight-db \
  --db-instance-class db.t2.micro \
  --engine postgres \
  --engine-version 15.3 \
  --master-username postgres \
  --master-user-password "STRONG_PASSWORD_HERE" \
  --allocated-storage 20 \
  --storage-type gp2 \
  --publicly-accessible false \
  --db-name social_insight \
  --region us-east-1
```

⚠️ **Wait 5-10 minutes for RDS to be ready!**

### Create ElastiCache (Redis)

```bash
# ⚠️ Student Account: Use cache.t2.micro (free tier)
aws elasticache create-cache-cluster \
  --cache-cluster-id social-insight-redis \
  --cache-node-type cache.t2.micro \
  --engine redis \
  --engine-version 7.0 \
  --num-cache-nodes 1 \
  --region us-east-1
```

⚠️ **Wait 3-5 minutes for ElastiCache to be ready!**

---

## 🚀 Step 2: Launch EC2 Instances

### Create EC2 Key Pair

```bash
aws ec2 create-key-pair \
  --key-name social-insight-key \
  --region us-east-1 \
  --query 'KeyMaterial' \
  --output text > ~/.ssh/social-insight-key.pem

chmod 400 ~/.ssh/social-insight-key.pem
```

### Launch 3 Instances

```bash
# EC2-1: Ingestion Service (t2.micro for free tier)
aws ec2 run-instances \
  --image-id ami-0c55b159cbfafe1f0 \
  --instance-type t2.micro \
  --key-name social-insight-key \
  --security-group-ids sg-xxxxx \
  --tag-specifications 'ResourceType=instance,Tags=[{Key=Name,Value=social-insight-ingestion}]' \
  --region us-east-1

# EC2-2: API Service
aws ec2 run-instances \
  --image-id ami-0c55b159cbfafe1f0 \
  --instance-type t2.micro \
  --key-name social-insight-key \
  --security-group-ids sg-xxxxx \
  --tag-specifications 'ResourceType=instance,Tags=[{Key=Name,Value=social-insight-api}]' \
  --region us-east-1

# EC2-3: Processing Service
aws ec2 run-instances \
  --image-id ami-0c55b159cbfafe1f0 \
  --instance-type t2.micro \
  --key-name social-insight-key \
  --security-group-ids sg-xxxxx \
  --tag-specifications 'ResourceType=instance,Tags=[{Key=Name,Value=social-insight-processing}]' \
  --region us-east-1
```

---

## 🐳 Step 3: Install Docker on EC2 Instances

SSH to each instance:
```bash
ssh -i ~/.ssh/social-insight-key.pem ec2-user@INSTANCE_IP
```

Run setup script:
```bash
#!/bin/bash
sudo yum update -y
sudo amazon-linux-extras install docker -y
sudo systemctl start docker
sudo systemctl enable docker
sudo usermod -aG docker ec2-user
sudo curl -L https://github.com/docker/compose/releases/download/v2.20.0/docker-compose-$(uname -s)-$(uname -m) -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
sudo yum install git -y
```

---

## 🚀 Step 4: Deploy Services

### EC2-1: Ingestion Service

```bash
ssh -i ~/.ssh/social-insight-key.pem ec2-user@EC2_1_IP
git clone <your-repo> && cd social-insight/production
vim .env.aws-ingestion  # Update KAFKA_BROKERS
docker-compose -f docker-compose.aws-ingestion.yml --env-file .env.aws-ingestion up -d
docker-compose logs -f
```

### EC2-2: API Service

```bash
ssh -i ~/.ssh/social-insight-key.pem ec2-user@EC2_2_IP
git clone <your-repo> && cd social-insight/production
vim .env.aws-api  # Update DB_HOST, REDIS_ADDR
docker-compose -f docker-compose.aws-api.yml --env-file .env.aws-api up -d
curl http://localhost:8888/api/health
```

### EC2-3: Processing Service

```bash
ssh -i ~/.ssh/social-insight-key.pem ec2-user@EC2_3_IP
git clone <your-repo> && cd social-insight/production
vim .env.aws-processing  # Update all endpoints
docker-compose -f docker-compose.aws-processing.yml --env-file .env.aws-processing up -d
docker-compose logs -f
```

---

## 🔍 Verification

### Check Services
```bash
# On each instance
docker-compose ps

# All services should show "Up"
```

### Test API
```bash
curl http://EC2_2_IP:8888/api/health
curl http://EC2_2_IP:8888/api/stats
```

### Monitor Data Flow
```bash
# EC2-1: Check crawlers
docker-compose logs crawler -f

# EC2-3: Check consumer
docker-compose logs consumer -f

# EC2-2: Check API
curl http://localhost:8888/api/stats
```

---

## ⚠️ Common Errors & Solutions

### "Connection Refused"
```
Solution:
  - Check Security Group allows ports
  - Verify RDS/Redis is running (aws rds describe-db-instances)
  - Check endpoints are correct in .env files
  - Wait 5-10 minutes for services to fully start
```

### "Out of Memory"
```
Solution:
  - Use t2.small instead of t2.micro (costs more)
  - Reduce batch sizes in consumer
  - Monitor: docker stats
```

### "Disk Space Full"
```
Solution:
  - docker image prune -a
  - df -h  (check space)
  - Use larger instance type
```

### "Kafka Connection Timeout"
```
Solution:
  - Check KAFKA_BROKERS is correct IP:9092
  - Verify Security Group allows 9092
  - docker-compose logs kafka
```

### "RDS Connection Failed"
```
Solution:
  - Verify security group allows 5432 from EC2
  - Check RDS is "available" status
  - Verify DB_HOST is correct endpoint
  - Test: psql -h <endpoint> -U postgres
```

---

## 💰 Cost Management

### Monitor Costs
```bash
aws ce get-cost-and-usage \
  --time-period Start=2026-01-01,End=2026-01-31 \
  --granularity DAILY \
  --metrics "BlendedCost"
```

### Stop Instances When Not Needed
```bash
aws ec2 stop-instances --instance-ids i-xxxxx i-yyyyy i-zzzzz
aws ec2 start-instances --instance-ids i-xxxxx i-yyyyy i-zzzzz
```

### Cost Optimization
- Use t2.micro (free tier eligible)
- Stop when not in use
- Delete unused Elastic IPs
- Use smaller RDS instance
- Consider alternatives to MSK (use EC2 + Kafka)

---

## 📞 Support & Documentation

- **AWS_DEPLOYMENT.md** - Detailed step-by-step guide
- **AWS_CHECKLIST.md** - Verification checklist
- **../local/README.md** - Local development guide

---

**Status**: ✅ Production Ready for AWS Student Account  
**Last Updated**: January 27, 2026  
**Version**: 1.0
