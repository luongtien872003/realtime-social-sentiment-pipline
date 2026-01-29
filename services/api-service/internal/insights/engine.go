package insights

import (
	"fmt"
	"social-insight/internal/models"
	"time"
)

// Insight đại diện cho một insight được phát hiện
type Insight struct {
	Type        string    `json:"type"` // "trending", "emerging", "declining", "anomaly"
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Confidence  float64   `json:"confidence"` // 0-1
	Timestamp   time.Time `json:"timestamp"`
}

// Engine phân tích dữ liệu để tìm insights
type Engine struct {
	posts     []models.Post
	timeRange time.Duration
}

// New tạo insights engine
func New(timeRange time.Duration) *Engine {
	return &Engine{
		posts:     make([]models.Post, 0),
		timeRange: timeRange,
	}
}

// AddPosts thêm posts để phân tích
func (e *Engine) AddPosts(posts []models.Post) {
	e.posts = append(e.posts, posts...)

	// Chỉ giữ posts trong time range
	cutoff := time.Now().Add(-e.timeRange)
	filtered := make([]models.Post, 0)
	for _, p := range e.posts {
		if p.CreatedAt.After(cutoff) {
			filtered = append(filtered, p)
		}
	}
	e.posts = filtered
}

// DetectTrending tìm topics/models đang trending
func (e *Engine) DetectTrending() []Insight {
	insights := make([]Insight, 0)

	// Count mentions of each topic
	topicCounts := make(map[string]int)
	for _, post := range e.posts {
		topicCounts[post.Topic]++
	}

	// Find fastest growing topics
	for topic, count := range topicCounts {
		if count > 5 { // Threshold
			insights = append(insights, Insight{
				Type:        "trending",
				Title:       topic + " is trending",
				Description: fmt.Sprintf("%d mentions in last 24h", count),
				Confidence:  0.8,
				Timestamp:   time.Now(),
			})
		}
	}

	return insights
}

// DetectEmergingModels tìm AI models mới nổi
func (e *Engine) DetectEmergingModels() []Insight {
	insights := make([]Insight, 0)

	// Simple heuristic: models that appeared in last few posts but weren't before
	models := map[string]int{
		"GPT-4": 0, "Claude": 0, "Llama": 0, "Gemini": 0, "Mistral": 0,
	}

	for _, post := range e.posts {
		for model := range models {
			// Simple string matching
			if len(post.Content) > 0 {
				// Would need proper NLP in production
				models[model]++
			}
		}
	}

	// If a model has high mentions, it's emerging
	for model, count := range models {
		if count > 10 {
			insights = append(insights, Insight{
				Type:        "emerging",
				Title:       model + " gaining attention",
				Description: fmt.Sprintf("%d mentions, strong growth", count),
				Confidence:  0.75,
				Timestamp:   time.Now(),
			})
		}
	}

	return insights
}

// DetectAnomalies finds unusual patterns
func (e *Engine) DetectAnomalies() []Insight {
	insights := make([]Insight, 0)

	if len(e.posts) == 0 {
		return insights
	}

	// Calculate average likes
	totalLikes := 0
	for _, post := range e.posts {
		totalLikes += post.Likes
	}
	avgLikes := totalLikes / len(e.posts)
	threshold := avgLikes * 3 // If > 3x average, it's anomaly

	// Find anomalies
	for _, post := range e.posts {
		if post.Likes > threshold {
			insights = append(insights, Insight{
				Type:        "anomaly",
				Title:       "Viral post detected",
				Description: fmt.Sprintf("Post with %d likes (avg: %d)", post.Likes, avgLikes),
				Confidence:  0.9,
				Timestamp:   time.Now(),
			})
		}
	}

	return insights
}

// GetAllInsights returns all detected insights
func (e *Engine) GetAllInsights() []Insight {
	insights := make([]Insight, 0)
	insights = append(insights, e.DetectTrending()...)
	insights = append(insights, e.DetectEmergingModels()...)
	insights = append(insights, e.DetectAnomalies()...)
	return insights
}

// TrendingScore calculates how "hot" a post is
// Based on: recency, engagement (likes+comments+shares), topic popularity
func TrendingScore(post models.Post, avgEngagement int) float64 {
	// Recency factor (newer = higher)
	timeSincePost := time.Since(post.CreatedAt).Minutes()
	recencyScore := 1.0 / (1.0 + timeSincePost/60.0) // Decay over hours

	// Engagement factor
	engagement := post.Likes + post.Comments*2 + post.Shares*3
	engagementScore := float64(engagement) / float64(avgEngagement)
	if engagementScore > 1 {
		engagementScore = 1 // Cap at 1
	}

	// Combined score (weighted)
	score := (recencyScore * 0.4) + (engagementScore * 0.6)
	return score
}
