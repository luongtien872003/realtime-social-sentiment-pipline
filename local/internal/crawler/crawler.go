// =====================================================
// CRAWLER INTERFACE - Base interface cho all crawlers
// =====================================================
// Mô tả: Shared interface và utilities cho crawlers
// Features: Dedup check, Kafka send, Redis tracking
// =====================================================

package crawler

import (
	"fmt"
	"social-insight/internal/kafka"
	"social-insight/internal/models"
	"social-insight/internal/redis"
	"time"
)

// Crawler là interface cho all crawlers
type Crawler interface {
	// Fetch lấy posts từ source và return
	Fetch() ([]models.Post, error)
	// Name return tên crawler ("hn", "medium", "devto")
	Name() string
}

// BaseCrawler chứa shared logic cho all crawlers
type BaseCrawler struct {
	producer  *kafka.Producer
	redis     *redis.Client
	source    string
	kafkaTopic string
}

// NewBaseCrawler tạo BaseCrawler instance
func NewBaseCrawler(
	producer *kafka.Producer,
	redis *redis.Client,
	source string,
	kafkaTopic string,
) *BaseCrawler {
	return &BaseCrawler{
		producer:   producer,
		redis:      redis,
		source:     source,
		kafkaTopic: kafkaTopic,
	}
}

// ProcessAndSend kiểm tra dedup, gửi Kafka, track Redis
// input: posts từ crawler
// output: số posts gửi thành công, số skip (duplicate)
func (b *BaseCrawler) ProcessAndSend(posts []models.Post) (sent, skipped int, err error) {
	if len(posts) == 0 {
		return 0, 0, nil
	}

	for _, post := range posts {
		// Check dedup
		seen, err := b.redis.CheckIfSeen(b.source, post.ID)
		if err != nil {
			fmt.Printf("❌ Redis check error for %s: %v\n", post.ID, err)
			skipped++
			continue
		}

		if seen {
			skipped++
			continue
		}

		// Send to Kafka
		// Send post via producer
		if err := b.producer.SendPost(post); err != nil {
			fmt.Printf("❌ Kafka send error for %s: %v\n", post.ID, err)
			skipped++
			continue
		}

		// Mark as seen in Redis (TTL 7 days)
		if err := b.redis.MarkAsSeen(b.source, post.ID, 7*24*time.Hour); err != nil {
			fmt.Printf("⚠️  Redis mark error for %s: %v\n", post.ID, err)
			// Don't fail, just warn
		}

		sent++
		fmt.Printf("✅ [%s] Sent post %s\n",
			b.source, post.ID)
	}

	return sent, skipped, nil
}
