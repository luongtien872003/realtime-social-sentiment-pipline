// =====================================================
// CONSUMER MAIN - Đọc data từ Kafka → Redis + PostgreSQL
// =====================================================
// Mô tả: Consumer đọc posts từ Kafka
// Lưu vào Redis (cache) và PostgreSQL (storage)
//
// Cách chạy: go run cmd/consumer/main.go
// =====================================================

package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"social-insight/config"
	"social-insight/internal/database"
	"social-insight/internal/kafka"
	"social-insight/internal/models"
	redisclient "social-insight/internal/redis"
)

// PostHandler xử lý posts nhận từ Kafka
type PostHandler struct {
	redis     *redisclient.Client
	db        *database.DB
	batchSize int

	// Buffer để batch insert
	buffer []models.Post

	// Counter
	processedCount int64
}

// HandlePost xử lý một post
func (h *PostHandler) HandlePost(post models.Post) error {
	// 1. Cache vào Redis (TTL 1 giờ)
	if err := h.redis.CachePost(post, time.Hour); err != nil {
		fmt.Printf("⚠️  Redis cache error: %v\n", err)
	}

	// 2. Cập nhật counters trong Redis
	h.redis.IncrementCounter("posts:total")
	h.redis.IncrementCounter(fmt.Sprintf("posts:%s", post.Topic))
	h.redis.IncrementCounter(fmt.Sprintf("sentiment:%s", post.Sentiment))

	// 3. Thêm vào recent posts
	h.redis.AddToRecentPosts(post)

	// 4. Thêm vào buffer để batch insert
	h.buffer = append(h.buffer, post)

	// 5. Flush nếu đủ batch size
	if len(h.buffer) >= h.batchSize {
		if err := h.flush(); err != nil {
			return err
		}
	}

	atomic.AddInt64(&h.processedCount, 1)
	return nil
}

// flush lưu buffer vào PostgreSQL
func (h *PostHandler) flush() error {
	if len(h.buffer) == 0 {
		return nil
	}

	if err := h.db.InsertPosts(h.buffer); err != nil {
		return fmt.Errorf("batch insert error: %w", err)
	}

	fmt.Printf("💾 Saved %d posts to PostgreSQL\n", len(h.buffer))
	h.buffer = h.buffer[:0] // Clear buffer
	return nil
}

func main() {
	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║     SOCIAL INSIGHT - DATA CONSUMER                         ║")
	fmt.Println("║     Kafka → Redis + PostgreSQL                             ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Load config
	if err := config.LoadEnvFile(".env"); err != nil {
		fmt.Printf("⚠️  Warning: %v\n", err)
	}
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("❌ Config error: %v\n", err)
		os.Exit(1)
	}
	cfg.Validate()
	cfg.LogConfig()

	// Context để graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Bắt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// ====== BƯỚC 1: Kết nối Redis ======
	fmt.Println("📡 Đang kết nối Redis...")
	redisClient, err := redisclient.NewClient(cfg.RedisAddr)
	if err != nil {
		fmt.Printf("❌ Lỗi Redis: %v\n", err)
		os.Exit(1)
	}
	defer redisClient.Close()
	fmt.Println("✅ Đã kết nối Redis")

	// ====== BƯỚC 2: Kết nối PostgreSQL ======
	fmt.Println("📡 Đang kết nối PostgreSQL...")
	db, err := database.NewDB(database.Config{
		Host:     cfg.PGHost,
		Port:     cfg.PGPort,
		User:     cfg.PGUser,
		Password: cfg.PGPassword,
		DBName:   cfg.PGDBName,
	})
	if err != nil {
		fmt.Printf("❌ Lỗi PostgreSQL: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()
	fmt.Println("✅ Đã kết nối PostgreSQL")

	// ====== BƯỚC 3: Tạo Handler ======
	handler := &PostHandler{
		redis:     redisClient,
		db:        db,
		batchSize: cfg.ConsumerBatchSize,
		buffer:    make([]models.Post, 0, cfg.ConsumerBatchSize),
	}

	// ====== BƯỚC 4: Tạo Kafka Consumer ======
	fmt.Println("📡 Đang kết nối Kafka...")
	consumer, err := kafka.NewConsumer(
		cfg.KafkaBrokers,
		cfg.ConsumerGroup,
		cfg.KafkaTopic,
		handler,
	)
	if err != nil {
		fmt.Printf("❌ Lỗi Kafka: %v\n", err)
		os.Exit(1)
	}
	defer consumer.Close()
	fmt.Printf("✅ Đã subscribe topic: %s\n\n", cfg.KafkaTopic)

	// ====== BƯỚC 5: Bắt đầu consume ======
	fmt.Println("👂 Đang lắng nghe messages từ Kafka...")
	fmt.Println("   Nhấn Ctrl+C để dừng")
	fmt.Println()

	// Goroutine để flush định kỳ
	go func() {
		ticker := time.NewTicker(cfg.ConsumerFlushInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				handler.flush()
			}
		}
	}()

	// Goroutine để in stats định kỳ
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				count := atomic.LoadInt64(&handler.processedCount)
				fmt.Printf("📊 Total processed: %d posts\n", count)
			}
		}
	}()

	// Goroutine để consume
	go func() {
		if err := consumer.Start(ctx); err != nil {
			fmt.Printf("❌ Consumer error: %v\n", err)
		}
	}()

	// Đợi signal
	<-sigChan
	fmt.Println("\n⚠️  Nhận tín hiệu dừng, đang shutdown...")

	// Cancel context
	cancel()

	// Flush remaining buffer
	handler.flush()

	// In kết quả cuối
	finalCount := atomic.LoadInt64(&handler.processedCount)
	dbCount, _ := db.GetPostCount()

	fmt.Println()
	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║                     KẾT QUẢ                                ║")
	fmt.Println("╠════════════════════════════════════════════════════════════╣")
	fmt.Printf("║  📝 Tổng posts đã xử lý: %d\n", finalCount)
	fmt.Printf("║  💾 Posts trong PostgreSQL: %d\n", dbCount)
	fmt.Println("╚════════════════════════════════════════════════════════════╝")
}
