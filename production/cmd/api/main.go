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

	"social-insight/internal/database"
	redisclient "social-insight/internal/redis"
)

// Cấu hình
const (
	APIPort = ":8888"

	// Redis config
	RedisAddr = "localhost:6379"

	// PostgreSQL config
	PGHost     = "localhost"
	PGPort     = 5432
	PGUser     = "postgres"
	PGPassword = "postgres123"
	PGDBName   = "social_insight"
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

// =====================================================
// MAIN
// =====================================================

func main() {
	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║     SOCIAL INSIGHT - API SERVER                            ║")
	fmt.Println("║     REST API cho Dashboard                                 ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Bắt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// ====== Kết nối Redis ======
	fmt.Println("📡 Đang kết nối Redis...")
	redisClient, err := redisclient.NewClient(RedisAddr)
	if err != nil {
		log.Printf("⚠️  Redis không available, một số features sẽ bị giới hạn: %v", err)
	} else {
		fmt.Println("✅ Đã kết nối Redis")
	}

	// ====== Kết nối PostgreSQL ======
	fmt.Println("📡 Đang kết nối PostgreSQL...")
	db, err := database.NewDB(database.Config{
		Host:     PGHost,
		Port:     PGPort,
		User:     PGUser,
		Password: PGPassword,
		DBName:   PGDBName,
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

	// Serve static files cho web dashboard
	fs := http.FileServer(http.Dir("web"))
	http.Handle("/", fs)

	// ====== Start server ======
	fmt.Printf("\n🚀 API Server đang chạy tại http://localhost%s\n", APIPort)
	fmt.Println("   Dashboard: http://localhost:8888")
	fmt.Println("   API Endpoints:")
	fmt.Println("   - GET /api/health    - Health check")
	fmt.Println("   - GET /api/stats     - Thống kê tổng quan")
	fmt.Println("   - GET /api/topics    - Thống kê theo topic")
	fmt.Println("   - GET /api/sentiment - Thống kê theo sentiment")
	fmt.Println("   - GET /api/authors   - Top tác giả")
	fmt.Println("   - GET /api/recent    - Posts mới nhất")
	fmt.Println("\n   Nhấn Ctrl+C để dừng")

	// Goroutine để chạy server
	go func() {
		if err := http.ListenAndServe(APIPort, nil); err != nil {
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
