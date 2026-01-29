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
	"math/rand"
	nethttp "net/http"
	"strconv"
	"time"
)

// Client là HTTP client reusable với retry logic
type Client struct {
	client     *nethttp.Client
	maxRetries int
	retryDelay time.Duration
	userAgent  string
	timeout    time.Duration
	baseDelay  time.Duration
	maxDelay   time.Duration
}

// NewClient tạo HTTP client mới
func NewClient(timeout time.Duration) *Client {
	// seed jitter
	rand.Seed(time.Now().UnixNano())

	return &Client{
		client: &nethttp.Client{
			Timeout: timeout,
		},
		maxRetries: 3,
		retryDelay: 1 * time.Second,
		userAgent:  "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
		timeout:    timeout,
		baseDelay:  500 * time.Millisecond,
		maxDelay:   30 * time.Second,
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
			time.Sleep(c.calculateBackoff(attempt))
			continue
		}

		// Ensure body is closed for all paths
		if resp.Body != nil {
			defer resp.Body.Close()
		}

		// Handle rate-limited responses (429) specially
		if resp.StatusCode == nethttp.StatusTooManyRequests {
			// Try to honor Retry-After header if present
			if ra := c.getRetryAfter(resp); ra > 0 {
				fmt.Printf("⚠️  HTTP 429 received, honoring Retry-After: %v\n", ra)
				time.Sleep(ra)
				continue
			}
			fmt.Printf("⚠️  HTTP 429 received, backing off...\n")
			time.Sleep(c.calculateBackoff(attempt))
			continue
		}

		// Non-200 responses
		if resp.StatusCode != nethttp.StatusOK {
			lastErr = fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
			fmt.Printf("⚠️  HTTP %d, retrying...\n", resp.StatusCode)
			time.Sleep(c.calculateBackoff(attempt))
			continue
		}

		// Read body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			lastErr = fmt.Errorf("read body error: %w", err)
			time.Sleep(c.calculateBackoff(attempt))
			continue
		}

		return body, nil
	}

	return nil, fmt.Errorf("failed after %d retries: %w", c.maxRetries, lastErr)
}

// calculateBackoff returns exponential backoff with small random jitter
func (c *Client) calculateBackoff(attempt int) time.Duration {
	// exponential: baseDelay * 2^attempt
	backoff := c.baseDelay * (1 << attempt)
	if backoff > c.maxDelay {
		backoff = c.maxDelay
	}

	// add up to 20% jitter
	jitterMax := int64(backoff / 5)
	if jitterMax <= 0 {
		return backoff
	}
	jitter := time.Duration(rand.Int63n(jitterMax))
	return backoff + jitter
}

// getRetryAfter parses Retry-After header; supports seconds or HTTP-date
func (c *Client) getRetryAfter(resp *nethttp.Response) time.Duration {
	if resp == nil {
		return 0
	}
	ra := resp.Header.Get("Retry-After")
	if ra == "" {
		return 0
	}
	// try integer seconds
	if secs, err := strconv.Atoi(ra); err == nil {
		return time.Duration(secs) * time.Second
	}
	// try HTTP date formats
	if t, err := time.Parse(time.RFC1123, ra); err == nil {
		delta := time.Until(t)
		if delta > 0 {
			return delta
		}
	}
	if t, err := time.Parse(time.RFC1123Z, ra); err == nil {
		delta := time.Until(t)
		if delta > 0 {
			return delta
		}
	}
	return 0
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
