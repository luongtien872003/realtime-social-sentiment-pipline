// =====================================================
// MEDIUM CRAWLER
// =====================================================
// Mô tả: Crawl posts từ Medium.com bằng RSS feed
// RSS Feed: https://medium.com/feed/tag/{topic}
// Topics: machine-learning, ai, cloud-computing, devops, startups
// Rate: 50 posts mỗi giờ
// =====================================================

package crawler

import (
	"encoding/xml"
	"fmt"
	"html"
	"regexp"
	httpclient "social-insight/internal/http"
	"social-insight/internal/models"
	"strings"
	"time"
)

// MediumItem là item từ RSS feed
type MediumItem struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
	Author      string `xml:"creator"`
	PubDate     string `xml:"pubDate"`
}

// MediumRSS là RSS feed structure
type MediumRSS struct {
	XMLName xml.Name      `xml:"rss"`
	Channel MediumChannel `xml:"channel"`
}

// MediumChannel là channel trong RSS
type MediumChannel struct {
	Items []MediumItem `xml:"item"`
}

// MediumCrawler crawl từ Medium.com
type MediumCrawler struct {
	*BaseCrawler
	client *httpclient.Client
	topics []string // Topics to crawl
	limit  int      // Number of posts per topic
}

// NewMediumCrawler tạo Medium crawler
func NewMediumCrawler(base *BaseCrawler, limit int) *MediumCrawler {
	return &MediumCrawler{
		BaseCrawler: base,
		client:      httpclient.NewClient(10 * time.Second),
		topics: []string{
			"machine-learning",
			"artificial-intelligence",
			"cloud-computing",
			"devops",
			"startups",
		},
		limit: limit,
	}
}

// Name return tên crawler
func (m *MediumCrawler) Name() string {
	return "medium"
}

// Fetch lấy posts từ Medium
func (m *MediumCrawler) Fetch() ([]models.Post, error) {
	fmt.Println("📡 Fetching Medium posts...")

	posts := make([]models.Post, 0)

	for _, topic := range m.topics {
		fmt.Printf("   📌 Crawling topic: %s\n", topic)

		topicPosts, err := m.fetchFromTopic(topic)
		if err != nil {
			fmt.Printf("   ⚠️  Topic error: %v\n", err)
			continue
		}

		posts = append(posts, topicPosts...)

		// Rate limiting: delay 200ms giữa mỗi topic
		time.Sleep(200 * time.Millisecond)
	}

	fmt.Printf("✅ Fetched %d posts from Medium\n", len(posts))
	return posts, nil
}

// fetchFromTopic fetch posts từ một topic
func (m *MediumCrawler) fetchFromTopic(topic string) ([]models.Post, error) {
	url := fmt.Sprintf("https://medium.com/feed/tag/%s", topic)

	data, err := m.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetch feed error: %w", err)
	}

	var rss MediumRSS
	if err := xml.Unmarshal(data, &rss); err != nil {
		return nil, fmt.Errorf("unmarshal RSS error: %w", err)
	}

	posts := make([]models.Post, 0)
	for i, item := range rss.Channel.Items {
		if i >= m.limit {
			break
		}

		post := m.parseItem(item, topic)
		if post != nil {
			posts = append(posts, *post)
		}
	}

	return posts, nil
}

// parseItem parse RSS item to Post
func (m *MediumCrawler) parseItem(item MediumItem, topic string) *models.Post {
	// Extract story ID từ link (Medium URL format)
	// https://medium.com/@author/title-abc123def456
	link := item.Link
	if link == "" {
		return nil
	}

	// Extract ID from URL using hex-like tail (robust to dashes in title)
	// e.g. https://medium.com/@author/title-with-dashes-abc123def456
	idRe := regexp.MustCompile(`([a-f0-9]{8,})$`)
	idMatch := idRe.FindStringSubmatch(link)
	if len(idMatch) < 2 {
		// fallback: use whole link hashed later by dedup
		return nil
	}
	id := idMatch[1]

	// Clean up description
	content := item.Title
	if len(item.Description) > 0 {
		// Remove HTML tags and unescape entities
		re := regexp.MustCompile(`<[^>]*>`)
		content = re.ReplaceAllString(item.Description, "")
		content = html.UnescapeString(content)
		content = strings.TrimSpace(content)
	}

	// Map topic
	topicName := m.mapTopic(topic)

	// Author extraction fallback: try to extract from link if empty
	author := strings.TrimSpace(item.Author)
	if author == "" {
		// look for @username in link
		authRe := regexp.MustCompile(`@([A-Za-z0-9_-]+)`)
		am := authRe.FindStringSubmatch(link)
		if len(am) > 1 {
			author = am[1]
		} else {
			author = "Unknown"
		}
	}

	// Parse publication date from item.PubDate if possible
	createdAt := time.Now()
	if item.PubDate != "" {
		// try common formats
		layouts := []string{time.RFC1123Z, time.RFC1123, time.RFC3339, time.RFC822}
		parsed := false
		for _, l := range layouts {
			if t, err := time.Parse(l, strings.TrimSpace(item.PubDate)); err == nil {
				createdAt = t
				parsed = true
				break
			}
		}
		if !parsed {
			// leave createdAt as now but log will be helpful elsewhere
		}
	}

	post := &models.Post{
		ID:        id,
		Author:    author,
		Content:   content,
		Topic:     topicName,
		Platform:  "medium",
		CreatedAt: createdAt,
		Sentiment: "neutral", // Will be set by ML service
	}

	return post
}

// mapTopic map Medium topic to our topics
func (m *MediumCrawler) mapTopic(topic string) string {
	switch topic {
	case "machine-learning", "artificial-intelligence":
		return "ai"
	case "cloud-computing":
		return "cloud"
	case "devops":
		return "devops"
	case "startups":
		return "startup"
	default:
		return "programming"
	}
}
