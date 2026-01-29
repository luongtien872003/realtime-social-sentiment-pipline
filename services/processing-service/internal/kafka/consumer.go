// =====================================================
// KAFKA CONSUMER - Đọc dữ liệu từ Kafka
// =====================================================
// Mô tả: Consumer để đọc posts từ Kafka topic
// Lưu vào Redis (cache) và PostgreSQL (storage)
// =====================================================

package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"social-insight/internal/models"

	"github.com/IBM/sarama"
)

// PostHandler là interface để xử lý posts nhận được
// Cho phép inject database và redis handlers
type PostHandler interface {
	HandlePost(post models.Post) error
}

// Consumer là struct wrapper cho Kafka consumer group
type Consumer struct {
	// consumerGroup là Sarama consumer group
	consumerGroup sarama.ConsumerGroup

	// topic là tên topic để đọc messages
	topic string

	// handler xử lý posts nhận được
	handler PostHandler

	// messageCount đếm số messages đã xử lý
	messageCount int64
}

// consumerGroupHandler implement sarama.ConsumerGroupHandler
type consumerGroupHandler struct {
	consumer *Consumer
}

// NewConsumer tạo một Kafka consumer mới
// brokers: danh sách Kafka brokers
// groupID: ID của consumer group
// topic: tên topic để subscribe
// handler: interface xử lý posts
func NewConsumer(brokers []string, groupID, topic string, handler PostHandler) (*Consumer, error) {
	// Cấu hình consumer
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest // Đọc từ message cũ nhất (từ đầu)

	// Tạo consumer group
	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, fmt.Errorf("không thể tạo consumer group: %w", err)
	}

	return &Consumer{
		consumerGroup: consumerGroup,
		topic:         topic,
		handler:       handler,
	}, nil
}

// Start bắt đầu consume messages
// Chạy trong goroutine riêng, dừng khi context bị cancel
func (c *Consumer) Start(ctx context.Context) error {
	handler := &consumerGroupHandler{consumer: c}

	for {
		// Check context cancellation
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Consume messages
		err := c.consumerGroup.Consume(ctx, []string{c.topic}, handler)
		if err != nil {
			fmt.Printf("❌ Consumer error: %v\n", err)
		}
	}
}

// GetMessageCount trả về số messages đã xử lý
func (c *Consumer) GetMessageCount() int64 {
	return c.messageCount
}

// Close đóng consumer
func (c *Consumer) Close() error {
	return c.consumerGroup.Close()
}

// =====================================================
// CONSUMER GROUP HANDLER IMPLEMENTATION
// =====================================================

// Setup được gọi khi consumer group session bắt đầu
func (h *consumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	fmt.Println("✅ Consumer group session started")
	return nil
}

// Cleanup được gọi khi consumer group session kết thúc
func (h *consumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	fmt.Println("🔄 Consumer group session ended")
	return nil
}

// ConsumeClaim xử lý messages từ partition
func (h *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		// Parse JSON thành Post
		var post models.Post
		if err := json.Unmarshal(message.Value, &post); err != nil {
			fmt.Printf("❌ Cannot unmarshal message: %v\n", err)
			continue
		}

		// Xử lý post (lưu vào DB, cache, etc.)
		if err := h.consumer.handler.HandlePost(post); err != nil {
			fmt.Printf("❌ Cannot handle post %s: %v\n", post.ID, err)
			continue
		}

		// Đánh dấu message đã xử lý
		session.MarkMessage(message, "")
		h.consumer.messageCount++

		// Log progress mỗi 10000 messages
		if h.consumer.messageCount%10000 == 0 {
			fmt.Printf("📊 Processed %d messages\n", h.consumer.messageCount)
		}
	}

	return nil
}
