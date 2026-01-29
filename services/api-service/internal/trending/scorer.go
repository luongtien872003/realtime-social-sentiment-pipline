package trending

import (
	"social-insight/internal/models"
	"time"
)

// TrendingPost represents a post with its trending score
type TrendingPost struct {
	Post  models.Post `json:"post"`
	Score float64     `json:"score"`
	Rank  int         `json:"rank"`
}

// Scorer calculates trending scores for posts
type Scorer struct {
	avgEngagement int
	recentPosts   []models.Post
}

// New creates a new Scorer
func New(avgEngagement int) *Scorer {
	return &Scorer{
		avgEngagement: avgEngagement,
		recentPosts:   make([]models.Post, 0),
	}
}

// AddPost adds a post for scoring
func (s *Scorer) AddPost(post models.Post) {
	s.recentPosts = append(s.recentPosts, post)
}

// AddPosts adds multiple posts
func (s *Scorer) AddPosts(posts []models.Post) {
	s.recentPosts = append(s.recentPosts, posts...)
}

// CalculateScore calculates trending score for a single post
// Score factors:
// - Recency (50%): Newer posts rank higher
// - Engagement (30%): Likes, comments, shares
// - Virality (20%): Above average engagement spike
func (s *Scorer) CalculateScore(post models.Post) float64 {
	// Recency factor (0-1)
	// Decay over 24 hours
	timeSincePost := time.Since(post.CreatedAt).Hours()
	recencyScore := 1.0 / (1.0 + timeSincePost/24.0)

	// Engagement factor (0-1)
	engagement := post.Likes + post.Comments*2 + post.Shares*3
	var engagementScore float64
	if s.avgEngagement > 0 {
		engagementScore = float64(engagement) / float64(s.avgEngagement)
		if engagementScore > 1 {
			engagementScore = 1 // Cap at 1
		}
	}

	// Virality bonus: if engagement is 2x average
	viralityScore := 0.0
	if engagement > (s.avgEngagement * 2) {
		viralityScore = 0.5 // 50% bonus
	}

	// Weighted combination
	score := (recencyScore * 0.5) + (engagementScore * 0.3) + (viralityScore * 0.2)
	return score
}

// GetTrending returns top N trending posts
func (s *Scorer) GetTrending(limit int) []TrendingPost {
	// Calculate scores for all posts
	scoredPosts := make([]TrendingPost, len(s.recentPosts))
	for i, post := range s.recentPosts {
		scoredPosts[i] = TrendingPost{
			Post:  post,
			Score: s.CalculateScore(post),
			Rank:  0,
		}
	}

	// Sort by score (bubble sort for simplicity, ok for small datasets)
	for i := 0; i < len(scoredPosts); i++ {
		for j := i + 1; j < len(scoredPosts); j++ {
			if scoredPosts[j].Score > scoredPosts[i].Score {
				scoredPosts[i], scoredPosts[j] = scoredPosts[j], scoredPosts[i]
			}
		}
	}

	// Assign ranks
	for i := range scoredPosts {
		scoredPosts[i].Rank = i + 1
	}

	// Return top N
	if limit > len(scoredPosts) {
		limit = len(scoredPosts)
	}
	return scoredPosts[:limit]
}

// GetTrendingByTopic returns trending posts for a specific topic
func (s *Scorer) GetTrendingByTopic(topic string, limit int) []TrendingPost {
	// Filter posts by topic
	topicPosts := make([]models.Post, 0)
	for _, post := range s.recentPosts {
		if post.Topic == topic {
			topicPosts = append(topicPosts, post)
		}
	}

	// Create temporary scorer for topic posts
	topicScorer := &Scorer{
		avgEngagement: s.avgEngagement,
		recentPosts:   topicPosts,
	}

	return topicScorer.GetTrending(limit)
}

// GetMostEngaging returns most engaged posts
func (s *Scorer) GetMostEngaging(limit int) []TrendingPost {
	// Sort by engagement
	scoredPosts := make([]TrendingPost, len(s.recentPosts))
	for i, post := range s.recentPosts {
		engagement := post.Likes + post.Comments*2 + post.Shares*3
		scoredPosts[i] = TrendingPost{
			Post:  post,
			Score: float64(engagement),
			Rank:  0,
		}
	}

	// Bubble sort
	for i := 0; i < len(scoredPosts); i++ {
		for j := i + 1; j < len(scoredPosts); j++ {
			if scoredPosts[j].Score > scoredPosts[i].Score {
				scoredPosts[i], scoredPosts[j] = scoredPosts[j], scoredPosts[i]
			}
		}
	}

	// Assign ranks
	for i := range scoredPosts {
		scoredPosts[i].Rank = i + 1
	}

	// Return top N
	if limit > len(scoredPosts) {
		limit = len(scoredPosts)
	}
	return scoredPosts[:limit]
}
