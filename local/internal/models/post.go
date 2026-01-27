// =====================================================
// POST MODEL - Cấu trúc dữ liệu cho bài viết
// =====================================================
// Mô tả: Định nghĩa struct Post để lưu trữ thông tin bài viết
// Sử dụng trong toàn bộ pipeline: Generator → Kafka → Storage
// =====================================================

package models

import (
	"time"
)

// Post là cấu trúc chính lưu thông tin một bài viết social media
// Struct này được sử dụng xuyên suốt pipeline
type Post struct {
	// ID là mã định danh duy nhất cho mỗi bài viết
	// Format: UUID rút gọn (8 ký tự)
	ID string `json:"id"`

	// Author là tên người đăng bài
	Author string `json:"author"`

	// Content là nội dung bài viết
	// Chủ đề xoay quanh công nghệ và AI
	Content string `json:"content"`

	// Topic là chủ đề của bài viết
	// Giá trị: "ai", "cloud", "devops", "programming", "startup"
	Topic string `json:"topic"`

	// Sentiment là cảm xúc của bài viết
	// Giá trị: "positive", "negative", "neutral"
	Sentiment string `json:"sentiment"`

	// Likes là số lượt thích (0-10000)
	Likes int `json:"likes"`

	// Comments là số bình luận (0-500)
	Comments int `json:"comments"`

	// Shares là số lượt chia sẻ (0-200)
	Shares int `json:"shares"`

	// Platform là nền tảng đăng bài
	// Giá trị: "twitter", "linkedin", "reddit", "hackernews"
	Platform string `json:"platform"`

	// CreatedAt là thời điểm tạo bài viết
	CreatedAt time.Time `json:"created_at"`
}

// Topics là danh sách các chủ đề công nghệ
var Topics = []string{
	"ai",          // Trí tuệ nhân tạo, Machine Learning
	"cloud",       // Cloud Computing (AWS, GCP, Azure)
	"devops",      // DevOps, Kubernetes, Docker
	"programming", // Ngôn ngữ lập trình, frameworks
	"startup",     // Startup, tech news, funding
}

// Platforms là danh sách các nền tảng mạng xã hội
var Platforms = []string{
	"twitter",    // Twitter/X
	"linkedin",   // LinkedIn
	"reddit",     // Reddit (r/programming, r/machinelearning)
	"hackernews", // Hacker News
}

// Sentiments là danh sách các loại cảm xúc
var Sentiments = []string{
	"positive", // Tích cực
	"negative", // Tiêu cực
	"neutral",  // Trung tính
}
