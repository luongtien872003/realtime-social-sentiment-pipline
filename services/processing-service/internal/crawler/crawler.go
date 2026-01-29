// =====================================================
// CRAWLER INTERFACE - Base interface cho all crawlers
// =====================================================
// Mô tả: Shared interface và utilities cho crawlers
// Features: Dedup check, Kafka send, Redis tracking
// =====================================================

package crawler

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"social-insight/internal/kafka"
	"social-insight/internal/models"
	"social-insight/internal/redis"
	"social-insight/internal/validation"
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
	producer   *kafka.Producer
	redis      *redis.Client
	source     string
	kafkaTopic string
	validator  *validation.Validator
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
		validator:  validation.New(),
	}
}

// ProcessAndSend kiểm tra dedup, gửi Kafka, track Redis
// input: posts từ crawler
// output: số posts gửi thành công, số skip (duplicate)
func (b *BaseCrawler) ProcessAndSend(posts []models.Post) (sent, skipped int, err error) {
	// Always update last crawl time even if no posts were found
	if err := b.redis.SetLastCrawl(b.source, time.Now()); err != nil {
		fmt.Printf("⚠️  Redis set last crawl error for %s: %v\n", b.source, err)
	}
	if len(posts) == 0 {
		return 0, 0, nil
	}

	for _, post := range posts {
		// Validate and sanitize post first
		if b.validator != nil {
			ok, verrs := b.validator.ValidatePost(&post)
			if !ok {
				fmt.Printf("❌ Validation failed for %s: %v\n", post.ID, verrs)
				skipped++
				continue
			}
		}
		// Compute content hash (content + author) to detect duplicates across sources
		h := sha256.Sum256([]byte(post.Content + "|" + post.Author))
		hashStr := hex.EncodeToString(h[:])

		// Check dedup by content hash first
		seenHash, err := b.redis.CheckIfSeen("content_hash", hashStr)
		if err != nil {
			fmt.Printf("❌ Redis hash check error for %s: %v\n", hashStr, err)
			// Fall back to ID-based check
		}
		if seenHash {
			skipped++
			fmt.Printf("⏭️  Skipped duplicate by content hash [%s]\n", hashStr)
			continue
		}

		// Fallback: check dedup by source-specific ID
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
		if err := b.producer.SendPost(post); err != nil {
			fmt.Printf("❌ Kafka send error for %s: %v\n", post.ID, err)
			skipped++
			continue
		}

		// Mark as seen in Redis: both by source ID (short TTL) and by content hash (long TTL)
		if err := b.redis.MarkAsSeen(b.source, post.ID, 7*24*time.Hour); err != nil {
			fmt.Printf("⚠️  Redis mark error for %s: %v\n", post.ID, err)
		}
		if err := b.redis.MarkAsSeen("content_hash", hashStr, 365*24*time.Hour); err != nil {
			fmt.Printf("⚠️  Redis mark error for content hash %s: %v\n", hashStr, err)
		}

		sent++
		fmt.Printf("✅ [%s] Sent post %s (hash=%s)\n", b.source, post.ID, hashStr)
	}

	// Update last crawl time in Redis
	if err := b.redis.SetLastCrawl(b.source, time.Now()); err != nil {
		fmt.Printf("⚠️  Redis set last crawl error for %s: %v\n", b.source, err)
	}

	return sent, skipped, nil
}
