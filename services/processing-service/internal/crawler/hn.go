// =====================================================
// HACKERNEWS CRAWLER
// =====================================================
// Mô tả: Crawl top stories từ HackerNews API
// API: https://hacker-news.firebaseio.com/v0/topstories.json
// Rate: 30 stories mỗi 5 phút (tránh rate limit)
// =====================================================

package crawler

import (
	"encoding/json"
	"fmt"
	httpclient "social-insight/internal/http"
	"social-insight/internal/models"
	"strings"
	"time"
)

// HNStory là response từ HN API
type HNStory struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	URL   string `json:"url"`
	By    string `json:"by"`
	Score int    `json:"score"`
	Text  string `json:"text"`
	Type  string `json:"type"`
}

// HNTopStoriesResponse là array của IDs
// type HNTopStoriesResponse []int

// HackerNewsCrawler crawl từ HackerNews
type HackerNewsCrawler struct {
	*BaseCrawler
	client       *httpclient.Client
	storiesLimit int
}

// NewHackerNewsCrawler tạo HN crawler
func NewHackerNewsCrawler(base *BaseCrawler, storiesLimit int) *HackerNewsCrawler {
	return &HackerNewsCrawler{
		BaseCrawler:  base,
		client:       httpclient.NewClient(10 * time.Second),
		storiesLimit: storiesLimit,
	}
}

// Name return tên crawler
func (h *HackerNewsCrawler) Name() string {
	return "hn"
}

// Fetch lấy top stories từ HN
func (h *HackerNewsCrawler) Fetch() ([]models.Post, error) {
	fmt.Println("📡 Fetching HackerNews top stories...")

	// Fetch list of top story IDs
	url := "https://hacker-news.firebaseio.com/v0/topstories.json"
	data, err := h.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetch top stories error: %w", err)
	}

	var storyIDs []int
	if err := json.Unmarshal(data, &storyIDs); err != nil {
		return nil, fmt.Errorf("unmarshal story IDs error: %w", err)
	}

	if len(storyIDs) > h.storiesLimit {
		storyIDs = storyIDs[:h.storiesLimit]
	}

	fmt.Printf("✅ Got %d story IDs, fetching details...\n", len(storyIDs))

	// Fetch details cho từng story (parallel)
	posts := make([]models.Post, 0, len(storyIDs))
	for i, id := range storyIDs {
		post, err := h.fetchStory(id)
		if err != nil {
			fmt.Printf("⚠️  Skip story %d: %v\n", id, err)
			continue
		}

		if post != nil {
			posts = append(posts, *post)
		}

		// Rate limiting: mỗi 10 requests delay 1 giây
		if (i+1)%10 == 0 {
			time.Sleep(100 * time.Millisecond)
		}
	}

	fmt.Printf("✅ Fetched %d stories from HN\n", len(posts))
	return posts, nil
}

// fetchStory fetch chi tiết một story
func (h *HackerNewsCrawler) fetchStory(id int) (*models.Post, error) {
	url := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%d.json", id)

	data, err := h.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetch story %d error: %w", id, err)
	}

	var story HNStory
	if err := json.Unmarshal(data, &story); err != nil {
		return nil, fmt.Errorf("unmarshal story %d error: %w", id, err)
	}

	// Filter: chỉ lấy "story" type, skip "poll", "job", etc.
	if story.Type != "story" {
		return nil, fmt.Errorf("skip non-story type: %s", story.Type)
	}

	// Filter: phải có title hoặc URL
	if story.Title == "" && story.URL == "" {
		return nil, fmt.Errorf("no title or url")
	}

	// Detect topic (AI, cloud, devops, programming, startup)
	topic := h.detectTopic(story.Title)

	post := &models.Post{
		ID:        fmt.Sprintf("%d", story.ID),
		Author:    story.By,
		Content:   story.Title,
		Topic:     topic,
		Platform:  "hackernews",
		Likes:     story.Score,
		CreatedAt: time.Now(),
		Sentiment: "neutral", // Will be set by ML service
	}

	return post, nil
}

// detectTopic detect topic từ title
func (h *HackerNewsCrawler) detectTopic(title string) string {
	title = strings.ToLower(title)

	if strings.Contains(title, "ai") || strings.Contains(title, "llm") ||
		strings.Contains(title, "machine learning") || strings.Contains(title, "chatgpt") ||
		strings.Contains(title, "gpt") || strings.Contains(title, "neural") {
		return "ai"
	}
	if strings.Contains(title, "cloud") || strings.Contains(title, "aws") ||
		strings.Contains(title, "gcp") || strings.Contains(title, "azure") {
		return "cloud"
	}
	if strings.Contains(title, "devops") || strings.Contains(title, "kubernetes") ||
		strings.Contains(title, "docker") || strings.Contains(title, "ci/cd") {
		return "devops"
	}
	if strings.Contains(title, "startup") || strings.Contains(title, "funding") ||
		strings.Contains(title, "venture") || strings.Contains(title, "ipo") {
		return "startup"
	}

	return "programming" // default
}
