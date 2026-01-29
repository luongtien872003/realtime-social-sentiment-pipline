// =====================================================
// KAFKA PRODUCER - Gửi dữ liệu vào Kafka
// =====================================================
// Mô tả: Producer để gửi posts vào Kafka topic
// Sử dụng thư viện IBM/sarama (fork của Shopify/sarama)
// =====================================================

package kafka

import (
	"encoding/json"
	"fmt"
	"social-insight/internal/models"
	"time"

	"github.com/IBM/sarama"
)

// Producer là struct wrapper cho Kafka producer
type Producer struct {
	// producer là Sarama async producer
	producer sarama.AsyncProducer

	// topic là tên topic để gửi messages
	topic string

	// successCount đếm số messages gửi thành công
	successCount int64

	// errorCount đếm số messages gửi thất bại
	errorCount int64
}

// NewProducer tạo một Kafka producer mới
// brokers: danh sách Kafka brokers (ví dụ: ["localhost:9092"])
// topic: tên topic để gửi messages
func NewProducer(brokers []string, topic string) (*Producer, error) {
	// Cấu hình producer
	config := sarama.NewConfig()

	// Bật async producer (gửi không đợi response)
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	// Cấu hình retry
	config.Producer.Retry.Max = 3
	config.Producer.Retry.Backoff = 100 * time.Millisecond

	// Cấu hình batch để tăng throughput
	config.Producer.Flush.Frequency = 100 * time.Millisecond // Flush mỗi 100ms
	config.Producer.Flush.Messages = 1000                    // Hoặc khi đủ 1000 messages

	// Tạo async producer
	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("không thể tạo Kafka producer: %w", err)
	}

	p := &Producer{
		producer: producer,
		topic:    topic,
	}

	// Goroutine xử lý success responses
	go func() {
		for range producer.Successes() {
			p.successCount++
		}
	}()

	// Goroutine xử lý errors
	go func() {
		for err := range producer.Errors() {
			p.errorCount++
			fmt.Printf("❌ Kafka error: %v\n", err)
		}
	}()

	return p, nil
}

// SendPost gửi một post vào Kafka
func (p *Producer) SendPost(post models.Post) error {
	// Chuyển post thành JSON
	data, err := json.Marshal(post)
	if err != nil {
		return fmt.Errorf("không thể marshal post: %w", err)
	}

	// Tạo message
	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.StringEncoder(post.ID), // Dùng post ID làm key
		Value: sarama.ByteEncoder(data),
	}

	// Gửi message (async)
	p.producer.Input() <- msg

	return nil
}

// SendPosts gửi nhiều posts cùng lúc
func (p *Producer) SendPosts(posts []models.Post) error {
	for _, post := range posts {
		if err := p.SendPost(post); err != nil {
			return err
		}
	}
	return nil
}

// GetStats trả về thống kê
func (p *Producer) GetStats() (success int64, errors int64) {
	return p.successCount, p.errorCount
}

// Close đóng producer
func (p *Producer) Close() error {
	return p.producer.Close()
}
