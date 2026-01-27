// =====================================================
// HTTP CLIENT - Reusable HTTP client with retry logic
// =====================================================
// Mô tả: Shared HTTP client cho crawlers
// Features: Timeout, Retry, User-Agent, Error handling
// =====================================================

package http

import (
	"fmt"
	"io"
	nethttp "net/http"
	"time"
)

// Client là HTTP client reusable với retry logic
type Client struct {
	client      *nethttp.Client
	maxRetries  int
	retryDelay  time.Duration
	userAgent   string
	timeout     time.Duration
}

// NewClient tạo HTTP client mới
func NewClient(timeout time.Duration) *Client {
	return &Client{
		client: &nethttp.Client{
			Timeout: timeout,
		},
		maxRetries: 3,
		retryDelay: 1 * time.Second,
		userAgent:  "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
		timeout:    timeout,
	}
}

// Get fetch dữ liệu từ URL với retry
func (c *Client) Get(url string) ([]byte, error) {
	var lastErr error

	for attempt := 0; attempt < c.maxRetries; attempt++ {
		resp, err := c.doRequest("GET", url, nil)
		if err != nil {
			lastErr = err
			fmt.Printf("⚠️  Attempt %d failed: %v, retrying...\n", attempt+1, err)
			time.Sleep(c.retryDelay)
			continue
		}
		defer resp.Body.Close()

		// Check status code
		if resp.StatusCode != nethttp.StatusOK {
			lastErr = fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
			fmt.Printf("⚠️  HTTP %d, retrying...\n", resp.StatusCode)
			time.Sleep(c.retryDelay)
			continue
		}

		// Read body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			lastErr = fmt.Errorf("read body error: %w", err)
			time.Sleep(c.retryDelay)
			continue
		}

		return body, nil
	}

	return nil, fmt.Errorf("failed after %d retries: %w", c.maxRetries, lastErr)
}

// doRequest thực hiện HTTP request
func (c *Client) doRequest(method, url string, body io.Reader) (*nethttp.Response, error) {
	req, err := nethttp.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("create request error: %w", err)
	}

	// Set headers
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Charset", "utf-8")

	// Send request
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}

	return resp, nil
}

// SetMaxRetries set số lần retry
func (c *Client) SetMaxRetries(retries int) {
	c.maxRetries = retries
}

// SetRetryDelay set delay giữa các retries
func (c *Client) SetRetryDelay(delay time.Duration) {
	c.retryDelay = delay
}

// SetUserAgent set custom user agent
func (c *Client) SetUserAgent(ua string) {
	c.userAgent = ua
}
