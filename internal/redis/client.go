// =====================================================
// REDIS CLIENT - Cache Layer
// =====================================================
// Mô tả: Redis client để cache dữ liệu
// Cache các posts mới nhất và thống kê realtime
// =====================================================

package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"social-insight/internal/models"
	"time"

	"github.com/go-redis/redis/v8"
)

// Client là wrapper cho Redis client
type Client struct {
	// rdb là Redis client instance
	rdb *redis.Client

	// ctx là context mặc định
	ctx context.Context
}

// NewClient tạo Redis client mới
// addr: địa chỉ Redis (ví dụ: "localhost:6379")
func NewClient(addr string) (*Client, error) {
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",  // Không có password cho local dev
		DB:       0,   // Database mặc định
		PoolSize: 100, // Connection pool size
	})

	// Test connection
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("không thể kết nối Redis: %w", err)
	}

	return &Client{
		rdb: rdb,
		ctx: ctx,
	}, nil
}

// CachePost lưu post vào cache với TTL
// key format: post:{id}
func (c *Client) CachePost(post models.Post, ttl time.Duration) error {
	key := fmt.Sprintf("post:%s", post.ID)

	data, err := json.Marshal(post)
	if err != nil {
		return fmt.Errorf("không thể marshal post: %w", err)
	}

	return c.rdb.Set(c.ctx, key, data, ttl).Err()
}

// GetPost lấy post từ cache
func (c *Client) GetPost(id string) (*models.Post, error) {
	key := fmt.Sprintf("post:%s", id)

	data, err := c.rdb.Get(c.ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, err
	}

	var post models.Post
	if err := json.Unmarshal(data, &post); err != nil {
		return nil, err
	}

	return &post, nil
}

// IncrementCounter tăng counter (dùng cho thống kê realtime)
// Ví dụ: posts:total, posts:ai, sentiment:positive
func (c *Client) IncrementCounter(key string) error {
	return c.rdb.Incr(c.ctx, key).Err()
}

// GetCounter lấy giá trị counter
func (c *Client) GetCounter(key string) (int64, error) {
	return c.rdb.Get(c.ctx, key).Int64()
}

// AddToRecentPosts thêm post vào danh sách recent (sorted set)
// Score là timestamp, giữ 1000 posts mới nhất
func (c *Client) AddToRecentPosts(post models.Post) error {
	data, err := json.Marshal(post)
	if err != nil {
		return err
	}

	// Thêm vào sorted set với score là timestamp
	err = c.rdb.ZAdd(c.ctx, "recent_posts", &redis.Z{
		Score:  float64(post.CreatedAt.Unix()),
		Member: data,
	}).Err()
	if err != nil {
		return err
	}

	// Giữ chỉ 1000 posts mới nhất
	return c.rdb.ZRemRangeByRank(c.ctx, "recent_posts", 0, -1001).Err()
}

// GetRecentPosts lấy n posts mới nhất
func (c *Client) GetRecentPosts(count int64) ([]models.Post, error) {
	// Lấy từ sorted set, order by score descending
	results, err := c.rdb.ZRevRange(c.ctx, "recent_posts", 0, count-1).Result()
	if err != nil {
		return nil, err
	}

	posts := make([]models.Post, 0, len(results))
	for _, data := range results {
		var post models.Post
		if err := json.Unmarshal([]byte(data), &post); err != nil {
			continue
		}
		posts = append(posts, post)
	}

	return posts, nil
}

// GetStats lấy thống kê từ cache
func (c *Client) GetStats() (map[string]int64, error) {
	stats := make(map[string]int64)

	// Danh sách các keys cần lấy
	keys := []string{
		"posts:total",
		"posts:ai", "posts:cloud", "posts:devops", "posts:programming", "posts:startup",
		"sentiment:positive", "sentiment:negative", "sentiment:neutral",
	}

	for _, key := range keys {
		val, err := c.rdb.Get(c.ctx, key).Int64()
		if err != nil && err != redis.Nil {
			continue
		}
		stats[key] = val
	}

	return stats, nil
}

// Close đóng Redis connection
func (c *Client) Close() error {
	return c.rdb.Close()
}
