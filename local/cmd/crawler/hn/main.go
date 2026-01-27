// =====================================================
// HACKERNEWS CRAWLER MAIN
// =====================================================
// Mô tả: Entry point cho HN crawler
// Chạy liên tục, crawl mỗi X phút
// Cách chạy: go run cmd/crawler/hn/main.go
// =====================================================

package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"social-insight/internal/crawler"
	"social-insight/internal/kafka"
	"social-insight/internal/redis"
)

// Cấu hình
const (
	KafkaBroker   = "localhost:9092"
	KafkaTopic    = "raw_posts"
	RedisAddr     = "localhost:6379"
	CrawlInterval = 5 * time.Minute // Crawl mỗi 5 phút
	StoriesLimit  = 30               // Fetch tối đa 30 stories
)

func main() {
	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║     HACKERNEWS CRAWLER                                     ║")
	fmt.Println("║     Crawl top stories mỗi 5 phút                           ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Bắt signal để graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// ====== BƯỚC 1: Tạo Kafka Producer ======
	fmt.Println("📡 Đang kết nối Kafka...")
	producer, err := kafka.NewProducer([]string{KafkaBroker}, KafkaTopic)
	if err != nil {
		fmt.Printf("❌ Kafka error: %v\n", err)
		fmt.Println("   Hãy chắc chắn Kafka đang chạy: docker-compose up -d")
		os.Exit(1)
	}
	defer producer.Close()
	fmt.Printf("✅ Kafka ready\n\n")

	// ====== BƯỚC 2: Tạo Redis Client ======
	fmt.Println("🔴 Đang kết nối Redis...")
	redisClient, err := redis.NewClient(RedisAddr)
	if err != nil {
		fmt.Printf("❌ Redis error: %v\n", err)
		fmt.Println("   Hãy chắc chắn Redis đang chạy: docker-compose up -d")
		os.Exit(1)
	}
	defer redisClient.Close()
	fmt.Printf("✅ Redis ready\n\n")

	// ====== BƯỚC 3: Tạo BaseCrawler ======
	baseCrawler := crawler.NewBaseCrawler(
		producer,
		redisClient,
		"hn",
		KafkaTopic,
	)

	// ====== BƯỚC 4: Tạo HN Crawler ======
	hnCrawler := crawler.NewHackerNewsCrawler(baseCrawler, StoriesLimit)

	// ====== BƯỚC 5: Chạy crawl loop ======
	fmt.Println("🚀 Bắt đầu crawling...")
	fmt.Printf("   Interval: %v, Limit: %d stories\n", CrawlInterval, StoriesLimit)
	fmt.Println("   Nhấn Ctrl+C để dừng")
	fmt.Println()
	fmt.Println("════════════════════════════════════════════════════════════")

	// Counters
	var totalSent int64
	var totalSkipped int64
	startTime := time.Now()

	// Channel để dừng goroutines
	stopChan := make(chan struct{})

	// Goroutine in stats mỗi 10 giây
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-stopChan:
				return
			case <-ticker.C:
				sent := atomic.LoadInt64(&totalSent)
				skipped := atomic.LoadInt64(&totalSkipped)
				elapsed := time.Since(startTime).Seconds()
				rate := float64(sent) / elapsed

				fmt.Printf("📊 [%s] Sent: %d | Skipped: %d | Rate: %.1f/s\n",
					time.Now().Format("15:04:05"),
					sent, skipped, rate)
			}
		}
	}()

	// Main crawl loop
	running := true
	go func() {
		ticker := time.NewTicker(CrawlInterval)
		defer ticker.Stop()

		for running {
			// Crawl
			fmt.Printf("\n🔍 [%s] Starting crawl...\n", time.Now().Format("15:04:05"))

			posts, err := hnCrawler.Fetch()
			if err != nil {
				fmt.Printf("❌ Fetch error: %v\n", err)
				<-ticker.C
				continue
			}

			// Process & send
			sent, skipped, err := baseCrawler.ProcessAndSend(posts)
			if err != nil {
				fmt.Printf("❌ Process error: %v\n", err)
			}

			atomic.AddInt64(&totalSent, int64(sent))
			atomic.AddInt64(&totalSkipped, int64(skipped))

			fmt.Printf("✅ Crawl complete: %d sent, %d skipped (duplicates)\n",
				sent, skipped)

			// Wait for next interval
			<-ticker.C
		}
	}()

	// Đợi signal để dừng
	<-sigChan
	running = false
	close(stopChan)

	// ====== KẾT THÚC ======
	fmt.Println("\n")
	fmt.Println("════════════════════════════════════════════════════════════")
	elapsed := time.Since(startTime)
	sent := atomic.LoadInt64(&totalSent)
	skipped := atomic.LoadInt64(&totalSkipped)

	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║              HACKERNEWS CRAWLER KẾT QUẢ                    ║")
	fmt.Println("╠════════════════════════════════════════════════════════════╣")
	fmt.Printf("║  ✅ Đã gửi: %d posts\n", sent)
	fmt.Printf("║  ⏭️  Đã bỏ qua: %d (duplicates)\n", skipped)
	fmt.Printf("║  ⏱️  Thời gian chạy: %s\n", elapsed.Round(time.Second))
	fmt.Printf("║  🚀 Tốc độ: %.1f posts/phút\n", float64(sent)*60/elapsed.Seconds())
	fmt.Println("╚════════════════════════════════════════════════════════════╝")
}
