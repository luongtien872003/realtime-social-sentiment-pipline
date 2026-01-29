package circuitbreaker

import (
	"fmt"
	"sync"
	"time"
)

// State của circuit breaker
type State int

const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)

// CircuitBreaker kiểm soát việc gửi requests đến một service
type CircuitBreaker struct {
	mu              sync.RWMutex
	state           State
	failureCount    int
	lastFailureTime time.Time

	// Config
	maxFailures      int           // Số lần fail trước khi mở (default 5)
	resetTimeout     time.Duration // Thời gian chờ trước khi half-open (default 30s)
	successThreshold int           // Số success cần để đóng lại (default 2)
	successCount     int           // Đếm success hiện tại

	// Callbacks
	onStateChange func(from, to State)
}

// New tạo circuit breaker mới
func New(maxFailures int, resetTimeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		state:            StateClosed,
		maxFailures:      maxFailures,
		resetTimeout:     resetTimeout,
		successThreshold: 2,
	}
}

// Call thực hiện action, nếu fail tăng counter
// Trả về error nếu circuit open
func (cb *CircuitBreaker) Call(fn func() error) error {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	// Nếu open, check xem có nên half-open không
	if cb.state == StateOpen {
		if time.Since(cb.lastFailureTime) > cb.resetTimeout {
			cb.setState(StateHalfOpen)
		} else {
			return fmt.Errorf("circuit breaker is OPEN (retry after %v)", cb.resetTimeout-time.Since(cb.lastFailureTime))
		}
	}

	// Thực hiện request
	err := fn()
	if err != nil {
		cb.onFailure()
		return err
	}

	// Success
	cb.onSuccess()
	return nil
}

// onFailure xử lý khi request fail
func (cb *CircuitBreaker) onFailure() {
	cb.failureCount++
	cb.lastFailureTime = time.Now()
	cb.successCount = 0

	fmt.Printf("⚠️  Circuit breaker: failure %d/%d\n", cb.failureCount, cb.maxFailures)

	if cb.failureCount >= cb.maxFailures {
		cb.setState(StateOpen)
	}
}

// onSuccess xử lý khi request success
func (cb *CircuitBreaker) onSuccess() {
	if cb.state == StateHalfOpen {
		cb.successCount++
		if cb.successCount >= cb.successThreshold {
			cb.setState(StateClosed)
		}
	}
	cb.failureCount = 0
}

// setState đổi trạng thái và callback
func (cb *CircuitBreaker) setState(newState State) {
	oldState := cb.state
	cb.state = newState

	stateNames := []string{"CLOSED", "OPEN", "HALF-OPEN"}
	fmt.Printf("🔄 Circuit breaker state: %s → %s\n", stateNames[oldState], stateNames[newState])

	if cb.onStateChange != nil {
		cb.onStateChange(oldState, newState)
	}

	if newState == StateClosed {
		cb.failureCount = 0
		cb.successCount = 0
	}
}

// GetState trả về trạng thái hiện tại
func (cb *CircuitBreaker) GetState() State {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

// SetStateChangeCallback đặt callback khi state đổi
func (cb *CircuitBreaker) SetStateChangeCallback(fn func(from, to State)) {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.onStateChange = fn
}
