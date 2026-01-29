// =====================================================
// POSTGRESQL DATABASE - Lưu trữ dữ liệu
// =====================================================
// Mô tả: Database client cho PostgreSQL
// Lưu trữ tất cả posts và hỗ trợ batch insert
// =====================================================

package database

import (
	"database/sql"
	"fmt"
	"social-insight/internal/models"
	"strings"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// DB là wrapper cho database connection
type DB struct {
	// conn là SQL connection
	conn *sql.DB
}

// Config chứa cấu hình kết nối database
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

// NewDB tạo database connection mới
func NewDB(cfg Config) (*DB, error) {
	// Connection string
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName,
	)

	// Mở connection
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("không thể mở connection: %w", err)
	}

	// Cấu hình connection pool
	conn.SetMaxOpenConns(50)                 // Tối đa 50 connections
	conn.SetMaxIdleConns(10)                 // Giữ 10 idle connections
	conn.SetConnMaxLifetime(5 * time.Minute) // Tái tạo connection mỗi 5 phút

	// Test connection
	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("không thể kết nối database: %w", err)
	}

	return &DB{conn: conn}, nil
}

// InsertPost chèn một post vào database
func (db *DB) InsertPost(post models.Post) error {
	query := `
		INSERT INTO posts (id, author, content, topic, sentiment, likes, comments, shares, platform, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (id) DO NOTHING
	`

	_, err := db.conn.Exec(query,
		post.ID,
		post.Author,
		post.Content,
		post.Topic,
		post.Sentiment,
		post.Likes,
		post.Comments,
		post.Shares,
		post.Platform,
		post.CreatedAt,
	)

	return err
}

// InsertPosts chèn nhiều posts cùng lúc (batch insert)
// Tối ưu performance với bulk insert
func (db *DB) InsertPosts(posts []models.Post) error {
	if len(posts) == 0 {
		return nil
	}

	// Xây dựng query với nhiều VALUES
	// INSERT INTO posts VALUES ($1...), ($2...), ...
	valueStrings := make([]string, 0, len(posts))
	valueArgs := make([]interface{}, 0, len(posts)*10)

	for i, post := range posts {
		offset := i * 10
		valueStrings = append(valueStrings, fmt.Sprintf(
			"($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			offset+1, offset+2, offset+3, offset+4, offset+5,
			offset+6, offset+7, offset+8, offset+9, offset+10,
		))
		valueArgs = append(valueArgs,
			post.ID,
			post.Author,
			post.Content,
			post.Topic,
			post.Sentiment,
			post.Likes,
			post.Comments,
			post.Shares,
			post.Platform,
			post.CreatedAt,
		)
	}

	query := fmt.Sprintf(`
		INSERT INTO posts (id, author, content, topic, sentiment, likes, comments, shares, platform, created_at)
		VALUES %s
		ON CONFLICT (id) DO NOTHING
	`, strings.Join(valueStrings, ","))

	_, err := db.conn.Exec(query, valueArgs...)
	return err
}

// GetPostCount trả về tổng số posts trong database
func (db *DB) GetPostCount() (int64, error) {
	var count int64
	err := db.conn.QueryRow("SELECT COUNT(*) FROM posts").Scan(&count)
	return count, err
}

// GetStatsByTopic trả về thống kê theo topic
func (db *DB) GetStatsByTopic() (map[string]int64, error) {
	query := `
		SELECT topic, COUNT(*) as count
		FROM posts
		GROUP BY topic
	`

	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make(map[string]int64)
	for rows.Next() {
		var topic string
		var count int64
		if err := rows.Scan(&topic, &count); err != nil {
			return nil, err
		}
		stats[topic] = count
	}

	return stats, nil
}

// GetStatsBySentiment trả về thống kê theo sentiment
func (db *DB) GetStatsBySentiment() (map[string]int64, error) {
	query := `
		SELECT sentiment, COUNT(*) as count
		FROM posts
		GROUP BY sentiment
	`

	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make(map[string]int64)
	for rows.Next() {
		var sentiment string
		var count int64
		if err := rows.Scan(&sentiment, &count); err != nil {
			return nil, err
		}
		stats[sentiment] = count
	}

	return stats, nil
}

// GetTopAuthors trả về top n tác giả có nhiều posts nhất
func (db *DB) GetTopAuthors(limit int) ([]map[string]interface{}, error) {
	query := `
		SELECT author, COUNT(*) as post_count, SUM(likes) as total_likes
		FROM posts
		GROUP BY author
		ORDER BY post_count DESC
		LIMIT $1
	`

	rows, err := db.conn.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]map[string]interface{}, 0)
	for rows.Next() {
		var author string
		var postCount, totalLikes int64
		if err := rows.Scan(&author, &postCount, &totalLikes); err != nil {
			return nil, err
		}
		results = append(results, map[string]interface{}{
			"author":      author,
			"post_count":  postCount,
			"total_likes": totalLikes,
		})
	}

	return results, nil
}

// Close đóng database connection
func (db *DB) Close() error {
	return db.conn.Close()
}
