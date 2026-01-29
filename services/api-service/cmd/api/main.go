// =====================================================
// API SERVER - REST API cho Dashboard
// =====================================================
// Mô tả: HTTP server cung cấp API cho Web UI
// Endpoints: /stats, /posts, /topics, /sentiment
//
// Cách chạy: go run cmd/api/main.go
// =====================================================

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"social-insight/config"
	"social-insight/internal/database"
	"social-insight/internal/models"
	redisclient "social-insight/internal/redis"
)

// Server chứa các dependencies
type Server struct {
	redis *redisclient.Client
	db    *database.DB
}

// =====================================================
// MIDDLEWARE
// =====================================================

// enableCORS thêm CORS headers cho mọi response
func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

// jsonResponse helper để trả về JSON
func jsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// =====================================================
// HANDLERS
// =====================================================

// handleHealth kiểm tra health của API
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, map[string]string{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	})
}

// handleOverallStats trả về thống kê tổng quan
func (s *Server) handleOverallStats(w http.ResponseWriter, r *http.Request) {
	// Thử lấy từ Redis cache trước
	stats, err := s.redis.GetStats()
	if err != nil || stats["posts:total"] == 0 {
		// Fallback: lấy từ PostgreSQL
		count, _ := s.db.GetPostCount()
		stats = map[string]int64{"posts:total": count}
	}

	jsonResponse(w, map[string]interface{}{
		"total_posts": stats["posts:total"],
		"by_topic": map[string]int64{
			"ai":          stats["posts:ai"],
			"cloud":       stats["posts:cloud"],
			"devops":      stats["posts:devops"],
			"programming": stats["posts:programming"],
			"startup":     stats["posts:startup"],
		},
		"by_sentiment": map[string]int64{
			"positive": stats["sentiment:positive"],
			"negative": stats["sentiment:negative"],
			"neutral":  stats["sentiment:neutral"],
		},
	})
}

// handleTopicStats trả về thống kê theo topic
func (s *Server) handleTopicStats(w http.ResponseWriter, r *http.Request) {
	stats, err := s.db.GetStatsByTopic()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, stats)
}

// handleSentimentStats trả về thống kê theo sentiment
func (s *Server) handleSentimentStats(w http.ResponseWriter, r *http.Request) {
	stats, err := s.db.GetStatsBySentiment()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, stats)
}

// handleTopAuthors trả về top tác giả
func (s *Server) handleTopAuthors(w http.ResponseWriter, r *http.Request) {
	authors, err := s.db.GetTopAuthors(10)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, authors)
}

// handleRecentPosts trả về posts mới nhất từ Redis
func (s *Server) handleRecentPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := s.redis.GetRecentPosts(20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, posts)
}

// handleCrawlers trả về trạng thái last crawl cho các source
func (s *Server) handleCrawlers(w http.ResponseWriter, r *http.Request) {
	sources := []string{"hn", "medium", "devto"}
	result := make(map[string]string)
	for _, src := range sources {
		t, err := s.redis.GetLastCrawl(src)
		if err != nil {
			result[src] = "unknown"
			continue
		}
		if t.IsZero() {
			result[src] = "never"
		} else {
			result[src] = t.Format(time.RFC3339)
		}
	}
	jsonResponse(w, result)
}

// handleInsights trả về insights phát hiện được
func (s *Server) handleInsights(w http.ResponseWriter, r *http.Request) {
	// Lấy posts từ 24h trước
	posts, err := s.db.GetPostsSince(time.Now().Add(-24 * time.Hour))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(posts) == 0 {
		jsonResponse(w, map[string]interface{}{
			"insights": []interface{}{},
			"message":  "Không đủ dữ liệu để phân tích",
		})
		return
	}

	// Đếm topics
	topicCounts := make(map[string]int)
	for _, post := range posts {
		topicCounts[post.Topic]++
	}

	insights := make([]map[string]interface{}, 0)

	// Trending topics
	for topic, count := range topicCounts {
		if count > 3 {
			insights = append(insights, map[string]interface{}{
				"type":        "trending",
				"title":       topic + " is trending",
				"description": fmt.Sprintf("%d mentions in last 24h", count),
				"confidence":  0.85,
				"timestamp":   time.Now().Format(time.RFC3339),
			})
		}
	}

	jsonResponse(w, map[string]interface{}{
		"insights": insights,
		"total":    len(insights),
	})
}

// handleCompare trả về so sánh hôm nay vs hôm qua
func (s *Server) handleCompare(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	today := now.Truncate(24 * time.Hour)
	yesterday := today.Add(-24 * time.Hour)

	// Lấy posts hôm nay
	todayPosts, _ := s.db.GetPostsInRange(today, now)

	// Lấy posts hôm qua
	yesterdayPosts, _ := s.db.GetPostsInRange(yesterday, today)

	// Tính metrics
	todayCount := len(todayPosts)
	yesterdayCount := len(yesterdayPosts)

	var percentChange float64
	if yesterdayCount > 0 {
		percentChange = float64(todayCount-yesterdayCount) / float64(yesterdayCount) * 100
	} else if todayCount > 0 {
		percentChange = 100
	}

	// Calculate engagement
	todayEngagement := 0
	yesterdayEngagement := 0

	for _, p := range todayPosts {
		todayEngagement += p.Likes + p.Comments + p.Shares
	}
	for _, p := range yesterdayPosts {
		yesterdayEngagement += p.Likes + p.Comments + p.Shares
	}

	jsonResponse(w, map[string]interface{}{
		"today": map[string]interface{}{
			"posts":      todayCount,
			"engagement": todayEngagement,
		},
		"yesterday": map[string]interface{}{
			"posts":      yesterdayCount,
			"engagement": yesterdayEngagement,
		},
		"comparison": map[string]interface{}{
			"posts_change":      todayCount - yesterdayCount,
			"posts_percent":     percentChange,
			"engagement_change": todayEngagement - yesterdayEngagement,
		},
	})
}

// handleTrending trả về trending posts
func (s *Server) handleTrending(w http.ResponseWriter, r *http.Request) {
	// Lấy posts từ 7 ngày trước
	posts, err := s.db.GetPostsSince(time.Now().Add(-7 * 24 * time.Hour))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(posts) == 0 {
		jsonResponse(w, map[string]interface{}{
			"trending": []interface{}{},
		})
		return
	}

	// Calculate average engagement
	totalEngagement := 0
	for _, p := range posts {
		totalEngagement += p.Likes + p.Comments + p.Shares
	}
	avgEngagement := totalEngagement / len(posts)
	if avgEngagement == 0 {
		avgEngagement = 1
	}

	// Simple trending algorithm: recent + highly engaged
	type TrendingItem struct {
		Post    models.Post `json:"post"`
		Score   float64     `json:"score"`
		Hotness string      `json:"hotness"` // 🔥 level
	}

	trending := make([]TrendingItem, 0)

	for _, post := range posts {
		// Recency factor
		hoursSince := time.Since(post.CreatedAt).Hours()
		recencyScore := 1.0 / (1.0 + hoursSince/24.0) // Decay over 24 hours

		// Engagement factor
		engagement := post.Likes + post.Comments*2 + post.Shares*3
		engagementScore := float64(engagement) / float64(avgEngagement)
		if engagementScore > 2 {
			engagementScore = 2 // Cap
		}

		// Combined score
		score := (recencyScore * 0.4) + (engagementScore * 0.6)

		hotness := "🔥"
		if score > 1.0 {
			hotness = "🔥🔥"
		}
		if score > 1.5 {
			hotness = "🔥🔥🔥"
		}

		trending = append(trending, TrendingItem{
			Post:    post,
			Score:   score,
			Hotness: hotness,
		})
	}

	// Sort by score (descending)
	for i := 0; i < len(trending); i++ {
		for j := i + 1; j < len(trending); j++ {
			if trending[j].Score > trending[i].Score {
				trending[i], trending[j] = trending[j], trending[i]
			}
		}
	}

	// Return top 10
	limit := 10
	if len(trending) < limit {
		limit = len(trending)
	}

	jsonResponse(w, map[string]interface{}{
		"trending": trending[:limit],
		"total":    limit,
	})
}

// =====================================================
// MAIN
// =====================================================

func main() {
	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║     SOCIAL INSIGHT - API SERVER                            ║")
	fmt.Println("║     REST API cho Dashboard                                 ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Load configuration
	if err := config.LoadEnvFile(".env"); err != nil {
		fmt.Printf("⚠️  Warning: %v\n", err)
	}
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("❌ Config error: %v\n", err)
		os.Exit(1)
	}
	if err := cfg.Validate(); err != nil {
		fmt.Printf("❌ Validation error: %v\n", err)
		os.Exit(1)
	}
	cfg.LogConfig()

	// Bắt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// ====== Kết nối Redis ======
	fmt.Println("📡 Đang kết nối Redis...")
	redisClient, err := redisclient.NewClient(cfg.RedisAddr)
	if err != nil {
		log.Printf("⚠️  Redis không available, một số features sẽ bị giới hạn: %v", err)
	} else {
		fmt.Println("✅ Đã kết nối Redis")
	}

	// ====== Kết nối PostgreSQL ======
	fmt.Println("📡 Đang kết nối PostgreSQL...")
	db, err := database.NewDB(database.Config{
		Host:     cfg.PGHost,
		Port:     cfg.PGPort,
		User:     cfg.PGUser,
		Password: cfg.PGPassword,
		DBName:   cfg.PGDBName,
	})
	if err != nil {
		log.Fatalf("❌ Lỗi PostgreSQL: %v", err)
	}
	fmt.Println("✅ Đã kết nối PostgreSQL")

	// Tạo server
	server := &Server{
		redis: redisClient,
		db:    db,
	}

	// ====== Đăng ký routes ======
	http.HandleFunc("/api/health", enableCORS(server.handleHealth))
	http.HandleFunc("/api/stats", enableCORS(server.handleOverallStats))
	http.HandleFunc("/api/topics", enableCORS(server.handleTopicStats))
	http.HandleFunc("/api/sentiment", enableCORS(server.handleSentimentStats))
	http.HandleFunc("/api/authors", enableCORS(server.handleTopAuthors))
	http.HandleFunc("/api/recent", enableCORS(server.handleRecentPosts))
	http.HandleFunc("/api/crawlers", enableCORS(server.handleCrawlers))
	http.HandleFunc("/api/insights", enableCORS(server.handleInsights))
	http.HandleFunc("/api/compare", enableCORS(server.handleCompare))
	http.HandleFunc("/api/trending", enableCORS(server.handleTrending))

	// Serve static files cho web dashboard
	fs := http.FileServer(http.Dir("web"))
	http.Handle("/", fs)

	// ====== Start server ======
	fmt.Printf("\n🚀 API Server đang chạy tại http://localhost%s\n", cfg.APIPort)
	fmt.Println("   Dashboard: http://localhost:8888")
	fmt.Println("   API Endpoints:")
	fmt.Println("   - GET /api/health    - Health check")
	fmt.Println("   - GET /api/stats     - Thống kê tổng quan")
	fmt.Println("   - GET /api/topics    - Thống kê theo topic")
	fmt.Println("   - GET /api/sentiment - Thống kê theo sentiment")
	fmt.Println("   - GET /api/authors   - Top tác giả")
	fmt.Println("   - GET /api/recent    - Posts mới nhất")
	fmt.Println("   - GET /api/insights  - Phát hiện insights")
	fmt.Println("   - GET /api/compare   - So sánh hôm nay vs hôm qua")
	fmt.Println("   - GET /api/trending  - Top trending posts")
	fmt.Println("\n   Nhấn Ctrl+C để dừng")

	// Goroutine để chạy server
	go func() {
		if err := http.ListenAndServe(cfg.APIPort, nil); err != nil {
			log.Fatalf("❌ Server error: %v", err)
		}
	}()

	// Đợi signal
	<-sigChan
	fmt.Println("\n⚠️  Đang shutdown...")

	if redisClient != nil {
		redisClient.Close()
	}
	db.Close()

	fmt.Println("✅ Server đã dừng")
}
