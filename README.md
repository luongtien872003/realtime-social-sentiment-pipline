# 🎯 Social Insight - Real-time Social Media Sentiment Analysis

**Production Ready** | **AWS Compatible** | **3-Service Architecture**

---

## 📦 Project Structure

This project is organized into two separate environments:

```
social-insight/
├── local/                             # 🔵 Local Development
│   ├── README.md                      # Development guide (detailed)
│   ├── docker-compose.yml             # Full-stack local setup
│   ├── docker-compose.local.yml       # Infrastructure only
│   └── [source code & config]
│
├── production/                        # 🔴 AWS Production
│   ├── README.md                      # AWS deployment guide
│   ├── docker-compose.*.yml           # Service compositions (4 files)
│   ├── .env.aws-*                     # AWS configs (4 files)
│   ├── AWS_*.md                       # AWS documentation (3 files)
│   └── [source code & config]
│
├── docs/                              # Shared documentation
├── go.mod & go.sum                    # Go dependencies
└── 📖 THIS FILE - Project overview
```

---

## 🚀 Quick Start

### Choose Your Environment:

#### 🔵 **Local Development** (Recommended for Development)
```bash
cd local/
# See local/README.md for complete setup
docker-compose up -d
# Open http://localhost:8888
```

**Best for:**
- 5-15 minute setup
- Development & testing
- Full 3-layer architecture
- Easy debugging

👉 **[Read local/README.md for complete guide](local/README.md)**

---

#### 🔴 **AWS Production** (For Deployment)
```bash
cd production/
# See production/README.md for complete AWS setup
# Deploy 3 services to separate EC2 instances
```

**Best for:**
- 30-45 minute deployment
- Production workloads
- Scalable architecture
- AWS Student Account optimization

👉 **[Read production/README.md for complete guide](production/README.md)**

---

## 📋 Architecture Overview

### 3-Layer Microservices

```
Layer 1: INGESTION        Layer 2: PROCESSING      Layer 3: PRESENTATION
┌──────────────────┐     ┌──────────────────┐     ┌──────────────────┐
│ HN Crawler       │     │ Consumer         │     │ API Gateway      │
│ Medium Crawler   │────▶│ ML Service       │────▶│ Frontend         │
│ DevTo Crawler    │     │                  │     │ Dashboard        │
└──────────────────┘     └──────────────────┘     └──────────────────┘
```

### Local vs Production Comparison

| Aspect | Local | Production |
|--------|-------|-----------|
| **Folder** | `local/` | `production/` |
| **Setup** | Docker Compose (all-in-one) | 3 AWS EC2 instances |
| **Time** | 5-15 minutes | 30-45 minutes |
| **Cost** | Free | $50-180/month (student) |
| **Best for** | Development | Production |
| **Scaling** | Limited | Full AWS scaling |

---

## 🛠️ Technology Stack

| Component | Technology | Purpose |
|-----------|-----------|---------|
| Crawlers | Go | Scrape HN, Medium, DevTo |
| Consumer | Go | Read Kafka, save to DB |
| ML Service | Python (FastAPI) | Sentiment, trends, AI detection |
| API | Go | REST API endpoints |
| Frontend | HTML/CSS/JS | Real-time dashboard |
| Database | PostgreSQL | Data persistence |
| Cache | Redis | Caching layer |
| Message Queue | Kafka | Async pipeline |
| Container | Docker | Application packaging |

---

## 📖 Documentation Map

| Document | Purpose | Audience |
|----------|---------|----------|
| **[local/README.md](local/README.md)** | Complete local dev guide | Developers |
| **[production/README.md](production/README.md)** | AWS deployment with student account guide | DevOps/Developers |
| **[production/AWS_DEPLOYMENT.md](production/AWS_DEPLOYMENT.md)** | Step-by-step AWS setup | DevOps |
| **[production/AWS_CHECKLIST.md](production/AWS_CHECKLIST.md)** | Verification & troubleshooting | DevOps |
| **[docs/00_START_HERE.md](docs/00_START_HERE.md)** | Project overview | Everyone |
| **[docs/AWS_READY.md](docs/AWS_READY.md)** | AWS readiness checklist | DevOps |
| **[docs/README_AWS.md](docs/README_AWS.md)** | AWS architecture details | DevOps |

---

## 🎯 Getting Started (3 Steps)

### Step 1: Choose Your Path
- **Development?** → Go to `local/`
- **Production?** → Go to `production/`

### Step 2: Read the README
- **Local Dev:** [local/README.md](local/README.md) (370+ lines of guidance)
- **AWS:** [production/README.md](production/README.md) (comprehensive AWS guide)

### Step 3: Follow the Instructions
- Copy/adjust environment files
- Run docker-compose commands
- Access dashboard or API endpoints

---

## 🔍 Key Features

- ✅ Multi-source web crawlers (HN, Medium, DevTo)
- ✅ Real-time data pipeline (Kafka)
- ✅ Sentiment analysis (Python ML)
- ✅ Trend detection & analysis
- ✅ AI model detection
- ✅ REST API with caching
- ✅ Real-time dashboard
- ✅ Docker containerization
- ✅ AWS deployment ready
- ✅ Production/Local separation
- ✅ Student account optimization

---

## 📊 Project Status

| Component | Status | Notes |
|-----------|--------|-------|
| **Architecture** | ✅ Complete | 3-layer microservices |
| **Local Setup** | ✅ Ready | Docker Compose fully configured |
| **AWS Setup** | ✅ Ready | Student account optimized |
| **Data Pipeline** | ✅ Working | Kafka consumer fixed (OffsetOldest) |
| **ML Service** | ✅ Complete | Sentiment, trends, AI detection |
| **API Gateway** | ✅ Ready | Full REST API |
| **Dashboard** | ✅ Live | Real-time visualization |
| **Documentation** | ✅ Complete | Local + AWS + troubleshooting |

---

## 💡 Why This Structure?

**Before:** All source code and configs in root directory
- ❌ Confusing for developers
- ❌ Hard to maintain separate environments
- ❌ Difficult to onboard new team members

**After:** Separated into local/ and production/
- ✅ Clear development/production separation
- ✅ Each has identical 3-layer structure
- ✅ Easy to understand project organization
- ✅ Simpler to deploy both environments
- ✅ Better for team collaboration

---

## 🚀 Next Steps

### If You're New to This Project:
1. Read this file (you're here!)
2. Choose your environment (local or production)
3. Go to that folder and read its README
4. Follow the setup instructions

### If You Want to Develop Locally:
```bash
cd local/
# Read local/README.md
docker-compose up -d
```

### If You Want to Deploy to AWS:
```bash
cd production/
# Read production/README.md
# Follow AWS deployment steps
```

### If You Want to Learn the Architecture:
- Read [docs/00_START_HERE.md](docs/00_START_HERE.md)
- Check the 3-layer architecture diagram above
- Review source code in ingestion/, processing/, presentation/

---

## 🔐 Security Notes

### Local Development
⚠️ **For development only!**
- Default weak credentials
- No encryption
- No authentication
- Not suitable for production

### Production (AWS)
✅ **Production-ready security:**
- AWS Security Groups
- VPC network isolation
- RDS encryption
- ElastiCache encryption
- Strong password policies

---

## 📞 Need Help?

### Local Issues?
→ See [local/README.md - Troubleshooting Section](local/README.md#-troubleshooting)

### AWS Issues?
→ See [production/README.md - Common Errors & Solutions](production/README.md#️-common-errors--solutions)

### General Questions?
→ See [docs/00_START_HERE.md](docs/00_START_HERE.md)

---

## 📝 Recent Changes

**Version 2.0 - Project Restructuring** (January 27, 2026)
- ✅ Reorganized into local/ and production/ folders
- ✅ Created comprehensive local development guide
- ✅ Created AWS production deployment guide with student account optimization
- ✅ Identical 3-layer architecture in both environments
- ✅ Improved documentation and onboarding
- ✅ Fixed Kafka consumer offset issue (previous version)
- ✅ Removed all hard-coded configuration values (previous version)

---

## 🎓 Learning Path

1. **Understand** → Read this README
2. **Develop** → Follow local/README.md
3. **Deploy** → Follow production/README.md
4. **Optimize** → Review AWS cost & security guides
5. **Scale** → Explore AWS auto-scaling options

---

## 📄 Repository

**GitHub:** https://github.com/luongtien872003/realtime-social-sentiment-pipline  
**Branch:** develop (for active development)  
**License:** MIT

---

## 🤝 Contributing

1. Clone the repository
2. Choose local/ for development
3. Make your changes
4. Test thoroughly
5. Commit with clear messages
6. Push to develop branch

---

**Status**: ✅ Production Ready  
**Last Updated**: January 27, 2026  
**Version**: 2.0 (Restructured)

👉 **[Start with local/README.md](local/README.md) or [production/README.md](production/README.md)**
docker-compose up -d
# Đợi 30s cho Kafka ready
```

### 3. Run Pipeline (2 terminals)
```bash
# Terminal 1: Consumer
go run cmd/consumer/main.go

# Terminal 2: API + Dashboard
go run cmd/api/main.go
```

### 4. View Dashboard
Open: **http://localhost:8888**

---

## 📊 Features

- ✅ **Kafka Message Queue**: High-throughput message processing
- ✅ **Redis Cache**: Realtime stats và recent posts
- ✅ **PostgreSQL Storage**: Batch insert với indexes
- ✅ **Crawler Pipeline**: Ingest data từ HN, Medium, DevTo
- ✅ **Web Dashboard**: Charts realtime với Chart.js

---

## 🔗 URLs

| Service | URL |
|---------|-----|
| Dashboard | http://localhost:8888 |
| Kafka UI | http://localhost:8080 |
| API Health | http://localhost:8888/api/health |

---

## 📁 Project Structure

```
├── cmd/
│   ├── consumer/    # Kafka to DB consumer
│   └── api/         # REST API server
├── internal/
│   ├── kafka/       # Producer & Consumer
│   ├── redis/       # Cache layer
│   └── database/    # PostgreSQL
├── web/             # Dashboard UI
├── processing/      # ML service
├── migrations/      # Database schema
└── .github/workflows/ # CI/CD
```

---

## 🤝 Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for branch strategy and workflow.

---

## 📄 License

MIT License
