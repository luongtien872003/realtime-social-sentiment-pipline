// =====================================================
// DEV.TO CRAWLER
// =====================================================
// Mô tả: Crawl posts từ Dev.to bằng REST API
// API: https://dev.to/api/articles
// Tags: ai, cloud, devops, programming, startups
// Rate: 30 posts mỗi 10 phút
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

// DevToArticle là response từ Dev.to API
type DevToArticle struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Author      struct {
		Name string `json:"name"`
	} `json:"user"`
	PublishedAt string   `json:"published_at"`
	Tags        []string `json:"tag_list"`
	BodyHTML    string   `json:"body_html"`
}

// DevToCrawler crawl từ Dev.to
type DevToCrawler struct {
	*BaseCrawler
	client *httpclient.Client
	tags   []string // Tags to crawl
	limit  int      // Number of posts per tag
}

// NewDevToCrawler tạo Dev.to crawler
func NewDevToCrawler(base *BaseCrawler, limit int) *DevToCrawler {
	return &DevToCrawler{
		BaseCrawler: base,
		client:      httpclient.NewClient(10 * time.Second),
		tags: []string{
			"ai",
			"machine-learning",
			"cloud",
			"devops",
			"startups",
		},
		limit: limit,
	}
}

// Name return tên crawler
func (d *DevToCrawler) Name() string {
	return "devto"
}

// Fetch lấy posts từ Dev.to
func (d *DevToCrawler) Fetch() ([]models.Post, error) {
	fmt.Println("📡 Fetching Dev.to posts...")

	posts := make([]models.Post, 0)

	for _, tag := range d.tags {
		fmt.Printf("   📌 Crawling tag: %s\n", tag)

		tagPosts, err := d.fetchFromTag(tag)
		if err != nil {
			fmt.Printf("   ⚠️  Tag error: %v\n", err)
			continue
		}

		posts = append(posts, tagPosts...)

		// Rate limiting: delay 200ms giữa mỗi tag
		time.Sleep(200 * time.Millisecond)
	}

	fmt.Printf("✅ Fetched %d posts from Dev.to\n", len(posts))
	return posts, nil
}

// fetchFromTag fetch posts từ một tag
func (d *DevToCrawler) fetchFromTag(tag string) ([]models.Post, error) {
	url := fmt.Sprintf("https://dev.to/api/articles?tag=%s&per_page=%d", tag, d.limit)

	data, err := d.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetch articles error: %w", err)
	}

	var articles []DevToArticle
	if err := json.Unmarshal(data, &articles); err != nil {
		return nil, fmt.Errorf("unmarshal articles error: %w", err)
	}

	posts := make([]models.Post, 0)
	for _, article := range articles {
		post := d.parseArticle(article, tag)
		if post != nil {
			posts = append(posts, *post)
		}
	}

	return posts, nil
}

// parseArticle parse Dev.to article to Post
func (d *DevToCrawler) parseArticle(article DevToArticle, tag string) *models.Post {
	if article.Title == "" {
		return nil
	}

	// Extract ID from URL or use article ID
	id := fmt.Sprintf("%d", article.ID)

	// Map topic
	topicName := d.mapTopic(tag)

	// Parse published date
	createdAt := time.Now()
	if article.PublishedAt != "" {
		if t, err := time.Parse(time.RFC3339, article.PublishedAt); err == nil {
			createdAt = t
		}
	}

	post := &models.Post{
		ID:        id,
		Author:    article.Author.Name,
		Content:   article.Title,
		Topic:     topicName,
		URL:       article.URL,
		Platform:  "devto",
		Source:    "devto",
		CreatedAt: createdAt,
		Sentiment: "neutral", // Will be set by ML service
		AIModel:   "none",    // Will be set by ML service
	}

	return post
}

// mapTopic map Dev.to tag to our topics
func (d *DevToCrawler) mapTopic(tag string) string {
	tag = strings.ToLower(tag)

	switch tag {
	case "ai", "machine-learning", "ml", "llm":
		return "ai"
	case "cloud", "aws", "gcp", "azure":
		return "cloud"
	case "devops", "kubernetes", "docker":
		return "devops"
	case "startups", "entrepreneurship":
		return "startup"
	default:
		return "programming"
	}
}
