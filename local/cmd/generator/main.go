// =====================================================
// GENERATOR MAIN - Streaming Mode
// =====================================================
// Mô tả: Generator chạy LIÊN TỤC, sinh 10k posts mỗi phút
// Rate: ~167 posts/giây
//
// Cách chạy: go run cmd/generator/main.go
// Dừng: Ctrl+C
// =====================================================

package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"social-insight/internal/generator"
	"social-insight/internal/kafka"
)

// =====================================================
// CẤU HÌNH - Có thể thay đổi theo nhu cầu
// =====================================================
const (
	// Kafka config
	KafkaBroker = "localhost:9092"
	KafkaTopic  = "raw_posts"

	// Generator config
	PostsPerMinute = 10000 // 10k posts mỗi phút
	Workers        = 10    // Số goroutines gửi song song
)

// Tính toán
var (
	// PostsPerSecond = 10000 / 60 ≈ 167 posts/giây
	PostsPerSecond = float64(PostsPerMinute) / 60.0

	// Delay giữa mỗi batch (ms)
	// Gửi 100 posts mỗi batch, 1.67 batches/giây
	BatchSize     = 100
	BatchesPerSec = PostsPerSecond / float64(BatchSize)
	BatchDelay    = time.Duration(1000/BatchesPerSec) * time.Millisecond
)

func main() {
	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║     SOCIAL INSIGHT - STREAMING GENERATOR                   ║")
	fmt.Println("║     10,000 posts/phút | Chạy liên tục                      ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Printf("📊 Cấu hình:\n")
	fmt.Printf("   - Rate: %d posts/phút (%.0f posts/giây)\n", PostsPerMinute, PostsPerSecond)
	fmt.Printf("   - Batch size: %d posts\n", BatchSize)
	fmt.Printf("   - Batch delay: %v\n", BatchDelay)
	fmt.Printf("   - Workers: %d goroutines\n", Workers)
	fmt.Println()

	// Bắt signal để graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// ====== BƯỚC 1: Tạo Kafka Producer ======
	fmt.Println("📡 Đang kết nối Kafka...")
	producer, err := kafka.NewProducer([]string{KafkaBroker}, KafkaTopic)
	if err != nil {
		fmt.Printf("❌ Lỗi: %v\n", err)
		fmt.Println("   Hãy chắc chắn Kafka đang chạy: docker-compose up -d")
		os.Exit(1)
	}
	defer producer.Close()
	fmt.Printf("✅ Đã kết nối Kafka, topic: %s\n\n", KafkaTopic)

	// ====== BƯỚC 2: Tạo Generator ======
	gen := generator.New()

	// ====== BƯỚC 3: Chạy streaming loop ======
	fmt.Println("🚀 Bắt đầu streaming...")
	fmt.Println("   Nhấn Ctrl+C để dừng")
	fmt.Println()
	fmt.Println("════════════════════════════════════════════════════════════")

	// Counters
	var totalSent int64
	var minuteSent int64
	startTime := time.Now()

	// Channel để dừng workers
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
				total := atomic.LoadInt64(&totalSent)
				elapsed := time.Since(startTime).Seconds()
				rate := float64(total) / elapsed

				success, errors := producer.GetStats()
				fmt.Printf("📊 [%s] Total: %d | Rate: %.0f/s | Success: %d | Errors: %d\n",
					time.Now().Format("15:04:05"),
					total, rate, success, errors)
			}
		}
	}()

	// Goroutine reset minute counter
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-stopChan:
				return
			case <-ticker.C:
				minCount := atomic.SwapInt64(&minuteSent, 0)
				fmt.Printf("📈 [MINUTE] Sent %d posts in last minute\n", minCount)
			}
		}
	}()

	// Main streaming loop
	running := true
	go func() {
		for running {
			// Sinh một batch
			for i := 0; i < BatchSize && running; i++ {
				post := gen.GenerateOne()
				if err := producer.SendPost(post); err != nil {
					// Ignore errors, continue streaming
					continue
				}
				atomic.AddInt64(&totalSent, 1)
				atomic.AddInt64(&minuteSent, 1)
			}

			// Delay để đạt đúng rate
			time.Sleep(BatchDelay)
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
	finalCount := atomic.LoadInt64(&totalSent)
	success, errors := producer.GetStats()

	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║                     KẾT QUẢ STREAMING                      ║")
	fmt.Println("╠════════════════════════════════════════════════════════════╣")
	fmt.Printf("║  📝 Tổng posts đã gửi: %d\n", finalCount)
	fmt.Printf("║  ✅ Thành công: %d\n", success)
	fmt.Printf("║  ❌ Lỗi: %d\n", errors)
	fmt.Printf("║  ⏱️  Thời gian chạy: %s\n", elapsed.Round(time.Second))
	fmt.Printf("║  🚀 Tốc độ trung bình: %.0f posts/giây\n", float64(finalCount)/elapsed.Seconds())
	fmt.Println("╚════════════════════════════════════════════════════════════╝")
}
