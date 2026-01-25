// =====================================================
// FAKE DATA GENERATOR - Sinh dữ liệu giả lập về Tech/AI
// =====================================================
// Mô tả: Sinh các bài viết giả lập về chủ đề công nghệ
// Sử dụng goroutines để sinh 100k posts nhanh chóng
// =====================================================

package generator

import (
	"fmt"
	"math/rand"
	"social-insight/internal/models"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Generator là struct chính để sinh dữ liệu giả lập
type Generator struct {
	// Templates cho từng chủ đề
	aiTemplates          []string
	cloudTemplates       []string
	devopsTemplates      []string
	programmingTemplates []string
	startupTemplates     []string

	// Từ khóa positive/negative để xác định sentiment
	positiveWords []string
	negativeWords []string

	// Tên tác giả ngẫu nhiên
	authors []string
}

// New tạo một Generator mới với các templates có sẵn
func New() *Generator {
	return &Generator{
		// === TEMPLATES CHO CHỦ ĐỀ AI ===
		aiTemplates: []string{
			"ChatGPT vừa ra bản mới, %s quá! Tính năng reasoning cải thiện đáng kể",
			"Mình vừa thử Claude 3.5 Sonnet, %s cho việc code",
			"Google Gemini 2.0 Flash %s, xử lý multimodal cực nhanh",
			"OpenAI vừa công bố o3, benchmark %s hơn cả con người",
			"Hugging Face ra model mới, open source %s",
			"LLM local chạy trên laptop giờ %s rồi các bạn ơi",
			"Fine-tuning với LoRA %s, tiết kiệm GPU đáng kể",
			"RAG + Vector DB combo %s cho enterprise AI",
			"AI Agent tự động hóa workflow, productivity %s",
			"Stable Diffusion 3 generate ảnh %s, chi tiết cực kỳ",
		},

		// === TEMPLATES CHO CHỦ ĐỀ CLOUD ===
		cloudTemplates: []string{
			"AWS Lambda cold start giờ %s, dưới 100ms",
			"GCP Cloud Run autoscaling %s, zero to hero trong 2 giây",
			"Azure OpenAI Service %s cho enterprise deployment",
			"Serverless architecture %s, giảm 70%% chi phí",
			"Multi-cloud strategy %s, không phụ thuộc vendor",
			"Cloud cost optimization %s, tiết kiệm $10k/tháng",
			"AWS Bedrock %s cho AI workloads",
			"Kubernetes trên EKS %s, managed control plane",
			"CloudFlare Workers %s cho edge computing",
			"Terraform IaC %s, infrastructure reproducible 100%%",
		},

		// === TEMPLATES CHO CHỦ ĐỀ DEVOPS ===
		devopsTemplates: []string{
			"GitHub Actions CI/CD %s, setup trong 5 phút",
			"Docker multi-stage build %s, image giảm 80%% size",
			"Kubernetes HPA %s, auto-scale theo CPU/memory",
			"ArgoCD GitOps %s, deployment declarative",
			"Prometheus + Grafana stack %s cho monitoring",
			"Helm charts %s, package Kubernetes apps",
			"Service mesh với Istio %s, observability tuyệt vời",
			"Infrastructure as Code %s, goodbye manual config",
			"Feature flags với LaunchDarkly %s cho progressive delivery",
			"Chaos engineering %s, tìm bug trước production",
		},

		// === TEMPLATES CHO CHỦ ĐỀ PROGRAMMING ===
		programmingTemplates: []string{
			"Golang concurrency %s, goroutines xử lý millions requests",
			"Rust memory safety %s, no more segfaults",
			"Python 3.12 %s, type hints ngày càng mạnh",
			"TypeScript 5.x %s, type inference cải thiện",
			"Next.js 14 App Router %s, React Server Components",
			"Go 1.22 %s, range over integers cuối cùng cũng có",
			"Bun runtime %s, nhanh hơn Node.js 3 lần",
			"Rust async/await %s, ecosystem mature rồi",
			"Zig language %s, low-level nhưng safe",
			"HTMX %s, hypermedia back to basics",
		},

		// === TEMPLATES CHO CHỦ ĐỀ STARTUP ===
		startupTemplates: []string{
			"Startup AI vừa raise Series A $50M, %s quá!",
			"Y Combinator batch mới %s, toàn AI startups",
			"Founder chia sẻ journey từ 0 đến $1M ARR, %s",
			"Remote-first culture %s cho productivity",
			"Product-market fit %s sau 6 tháng iterate",
			"Bootstrapped to $10M ARR, story %s",
			"Developer tools startup %s, B2D model works",
			"SaaS metrics chuẩn %s cho early stage",
			"Fundraising climate 2024 %s hơn năm trước",
			"AI-native startup playbook %s, must read",
		},

		// === TỪ KHÓA SENTIMENT ===
		positiveWords: []string{
			"tuyệt vời", "xuất sắc", "amazing", "great", "awesome",
			"impressive", "incredible", "fantastic", "brilliant", "superb",
			"game-changer", "revolutionary", "best-in-class", "top-tier",
		},
		negativeWords: []string{
			"tệ", "disappointing", "underwhelming", "buggy", "slow",
			"frustrating", "terrible", "awful", "worst", "broken",
			"unstable", "unreliable", "overpriced", "overhyped",
		},

		// === DANH SÁCH TÁC GIẢ ===
		authors: []string{
			"TechBro_VN", "AI_Enthusiast", "CloudNinja", "DevOps_Master",
			"Golang_Fan", "Rust_Advocate", "Python_Guru", "JS_Developer",
			"Startup_Founder", "VC_Investor", "Senior_Engineer", "Tech_Lead",
			"ML_Researcher", "Data_Scientist", "Backend_Dev", "Frontend_Dev",
			"Full_Stack", "Architect", "CTO_Startup", "Engineering_Manager",
			"Open_Source_Contributor", "Indie_Hacker", "Solo_Founder",
		},
	}
}

// GenerateOne tạo một bài viết ngẫu nhiên
func (g *Generator) GenerateOne() models.Post {
	// Chọn ngẫu nhiên topic
	topic := models.Topics[rand.Intn(len(models.Topics))]

	// Chọn template theo topic
	var templates []string
	switch topic {
	case "ai":
		templates = g.aiTemplates
	case "cloud":
		templates = g.cloudTemplates
	case "devops":
		templates = g.devopsTemplates
	case "programming":
		templates = g.programmingTemplates
	case "startup":
		templates = g.startupTemplates
	}

	// Random sentiment với tỷ lệ: 50% positive, 30% neutral, 20% negative
	r := rand.Float32()
	var sentiment string
	var word string

	if r < 0.5 {
		sentiment = "positive"
		word = g.positiveWords[rand.Intn(len(g.positiveWords))]
	} else if r < 0.8 {
		sentiment = "neutral"
		word = "ổn" // Từ trung tính
	} else {
		sentiment = "negative"
		word = g.negativeWords[rand.Intn(len(g.negativeWords))]
	}

	// Tạo nội dung từ template
	template := templates[rand.Intn(len(templates))]
	content := fmt.Sprintf(template, word)

	// Tạo post
	post := models.Post{
		ID:        strings.Split(uuid.New().String(), "-")[0], // UUID rút gọn 8 ký tự
		Author:    g.authors[rand.Intn(len(g.authors))],
		Content:   content,
		Topic:     topic,
		Sentiment: sentiment,
		Likes:     rand.Intn(10001), // 0-10000
		Comments:  rand.Intn(501),   // 0-500
		Shares:    rand.Intn(201),   // 0-200
		Platform:  models.Platforms[rand.Intn(len(models.Platforms))],
		CreatedAt: time.Now(),
	}

	return post
}

// GenerateBatch tạo nhiều posts cùng lúc (không dùng goroutines)
// Dành cho trường hợp đơn giản, số lượng nhỏ
func (g *Generator) GenerateBatch(count int) []models.Post {
	posts := make([]models.Post, count)
	for i := 0; i < count; i++ {
		posts[i] = g.GenerateOne()
	}
	return posts
}

// GenerateBatchConcurrent tạo nhiều posts sử dụng goroutines
// workers: số goroutines chạy song song
// count: tổng số posts cần tạo
// Trả về channel để nhận posts theo batch
func (g *Generator) GenerateBatchConcurrent(count int, workers int, batchSize int) <-chan []models.Post {
	// Channel để gửi posts ra ngoài
	results := make(chan []models.Post, workers*2)

	// Tính số posts mỗi worker cần sinh
	postsPerWorker := count / workers
	remainder := count % workers

	go func() {
		// WaitGroup để đợi tất cả workers hoàn thành
		done := make(chan bool, workers)

		for w := 0; w < workers; w++ {
			// Worker cuối cùng nhận thêm phần dư
			workerCount := postsPerWorker
			if w == workers-1 {
				workerCount += remainder
			}

			// Khởi chạy goroutine cho mỗi worker
			go func(workerID, totalPosts int) {
				// Sinh posts theo batch
				for i := 0; i < totalPosts; i += batchSize {
					// Tính kích thước batch thực tế
					currentBatchSize := batchSize
					if i+batchSize > totalPosts {
						currentBatchSize = totalPosts - i
					}

					// Sinh batch posts
					batch := make([]models.Post, currentBatchSize)
					for j := 0; j < currentBatchSize; j++ {
						batch[j] = g.GenerateOne()
					}

					// Gửi batch vào channel
					results <- batch
				}
				done <- true
			}(w, workerCount)
		}

		// Đợi tất cả workers hoàn thành
		for i := 0; i < workers; i++ {
			<-done
		}

		// Đóng channel khi hoàn tất
		close(results)
	}()

	return results
}
