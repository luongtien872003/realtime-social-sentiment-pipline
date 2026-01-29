// =====================================================
// CONFIG LOADER - Quản lý tất cả cấu hình
// =====================================================
// Mô tả: Load config từ .env file hoặc default values
// Hỗ trợ multiple environments (local, production)
// =====================================================

package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config chứa tất cả cấu hình ứng dụng
type Config struct {
	// Kafka
	KafkaBrokers  []string
	KafkaTopic    string
	ConsumerGroup string

	// Redis
	RedisAddr string

	// PostgreSQL
	PGHost     string
	PGPort     int
	PGUser     string
	PGPassword string
	PGDBName   string

	// API
	APIPort string

	// Crawler - HackerNews
	HNCrawlInterval time.Duration
	HNStoriesLimit  int

	// Crawler - Dev.to
	DevtoCrawlInterval time.Duration
	DevtoPostsPerTag   int
	DevtoTags          []string

	// Crawler - Medium
	MediumCrawlInterval time.Duration
	MediumPostsPerTopic int
	MediumTopics        []string

	// Consumer
	ConsumerBatchSize     int
	ConsumerFlushInterval time.Duration

	// HTTP Client
	HTTPClientTimeout time.Duration
	HTTPMaxRetries    int
	HTTPRetryDelay    time.Duration
}

// Load tải config từ environment variables hoặc default values
func Load() (*Config, error) {
	// Helper to get Kafka brokers with fallback
	kafkaBrokersEnv := os.Getenv("KAFKA_BROKERS")
	var kafkaBrokers []string
	if kafkaBrokersEnv != "" {
		kafkaBrokers = parseStringSlice(kafkaBrokersEnv, ",")
	} else {
		kafkaBrokers = parseStringSlice("localhost:9092", ",")
	}

	cfg := &Config{
		// Default values (Local development)
		KafkaBrokers:          kafkaBrokers,
		KafkaTopic:            getEnv("KAFKA_TOPIC", "raw_posts"),
		ConsumerGroup:         getEnv("CONSUMER_GROUP", "social_insight_consumer"),
		RedisAddr:             getEnv("REDIS_ADDR", "localhost:6379"),
		PGHost:                getEnv("PG_HOST", "localhost"),
		PGPort:                getEnvInt("PG_PORT", 5432),
		PGUser:                getEnv("PG_USER", "postgres"),
		PGPassword:            getEnv("PG_PASSWORD", "postgres123"),
		PGDBName:              getEnv("PG_DBNAME", "social_insight"),
		APIPort:               getEnv("API_PORT", ":8888"),
		HNCrawlInterval:       parseDuration(getEnv("HN_CRAWL_INTERVAL", "5m")),
		HNStoriesLimit:        getEnvInt("HN_STORIES_LIMIT", 30),
		DevtoCrawlInterval:    parseDuration(getEnv("DEVTO_CRAWL_INTERVAL", "10m")),
		DevtoPostsPerTag:      getEnvInt("DEVTO_POSTS_PER_TAG", 6),
		DevtoTags:             parseStringSlice(getEnv("DEVTO_TAGS", "ai,machine-learning,cloud,devops,startups"), ""),
		MediumCrawlInterval:   parseDuration(getEnv("MEDIUM_CRAWL_INTERVAL", "10m")),
		MediumPostsPerTopic:   getEnvInt("MEDIUM_POSTS_PER_TOPIC", 10),
		MediumTopics:          parseStringSlice(getEnv("MEDIUM_TOPICS", "machine-learning,artificial-intelligence,cloud-computing,devops,startups"), ""),
		ConsumerBatchSize:     getEnvInt("CONSUMER_BATCH_SIZE", 500),
		ConsumerFlushInterval: parseDuration(getEnv("CONSUMER_FLUSH_INTERVAL", "2s")),
		HTTPClientTimeout:     parseDuration(getEnv("HTTP_CLIENT_TIMEOUT", "10s")),
		HTTPMaxRetries:        getEnvInt("HTTP_MAX_RETRIES", 3),
		HTTPRetryDelay:        parseDuration(getEnv("HTTP_RETRY_DELAY", "1s")),
	}

	return cfg, nil
}

// =====================================================
// HELPER FUNCTIONS
// =====================================================

// getEnv lấy environment variable với default value
func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

// getEnvInt lấy int environment variable
func getEnvInt(key string, defaultVal int) int {
	valStr := getEnv(key, "")
	if valStr == "" {
		return defaultVal
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		fmt.Printf("⚠️  Invalid int value for %s: %v, using default: %d\n", key, err, defaultVal)
		return defaultVal
	}
	return val
}

// parseDuration parse duration string
func parseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		fmt.Printf("⚠️  Invalid duration: %s, using 5s\n", s)
		return 5 * time.Second
	}
	return d
}

// parseStringSlice parse comma-separated string to slice
func parseStringSlice(s, separator string) []string {
	if s == "" {
		return []string{}
	}
	if separator == "" {
		separator = ","
	}
	parts := strings.Split(s, separator)
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// LogConfig in ra config (với password masked)
func (c *Config) LogConfig() {
	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║              APPLICATION CONFIGURATION                     ║")
	fmt.Println("╠════════════════════════════════════════════════════════════╣")
	fmt.Printf("║ Kafka Brokers: %s\n", strings.Join(c.KafkaBrokers, ", "))
	fmt.Printf("║ Kafka Topic: %s\n", c.KafkaTopic)
	fmt.Printf("║ Redis: %s\n", c.RedisAddr)
	fmt.Printf("║ PostgreSQL: %s:%d/%s (user: %s)\n", c.PGHost, c.PGPort, c.PGDBName, c.PGUser)
	fmt.Printf("║ API Port: %s\n", c.APIPort)
	fmt.Println("║")
	fmt.Printf("║ HN Crawler: interval=%v, limit=%d\n", c.HNCrawlInterval, c.HNStoriesLimit)
	fmt.Printf("║ Devto Crawler: interval=%v, %d posts/tag\n", c.DevtoCrawlInterval, c.DevtoPostsPerTag)
	fmt.Printf("║ Medium Crawler: interval=%v, %d posts/topic\n", c.MediumCrawlInterval, c.MediumPostsPerTopic)
	fmt.Printf("║ Consumer: batch=%d, flush=%v\n", c.ConsumerBatchSize, c.ConsumerFlushInterval)
	fmt.Println("╚════════════════════════════════════════════════════════════╝")
}

// Validate check configuration values
func (c *Config) Validate() error {
	if len(c.KafkaBrokers) == 0 {
		return fmt.Errorf("kafka brokers not configured")
	}
	if c.KafkaTopic == "" {
		return fmt.Errorf("kafka topic not configured")
	}
	if c.RedisAddr == "" {
		return fmt.Errorf("redis address not configured")
	}
	if c.PGHost == "" {
		return fmt.Errorf("postgresql host not configured")
	}
	return nil
}
