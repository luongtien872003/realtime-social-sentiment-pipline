package orchestrator

import (
	"fmt"
	"sync"
	"time"

	"social-insight/internal/crawler"
)

// CrawlerOrchestrator quản lý việc chạy đồng thời các crawlers
type CrawlerOrchestrator struct {
	crawlers map[string]crawler.Crawler
	mu       sync.RWMutex

	// Metrics
	runs     int64
	failures map[string]int64
	lastRun  time.Time
}

// New tạo orchestrator mới
func New() *CrawlerOrchestrator {
	return &CrawlerOrchestrator{
		crawlers: make(map[string]crawler.Crawler),
		failures: make(map[string]int64),
	}
}

// Register đăng ký crawler
func (o *CrawlerOrchestrator) Register(name string, c crawler.Crawler) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.crawlers[name] = c
	o.failures[name] = 0
}

// RunParallel chạy tất cả crawlers đồng thời
// Trả về kết quả từ mỗi crawler
func (o *CrawlerOrchestrator) RunParallel() map[string]CrawlResult {
	o.mu.RLock()
	crawlers := o.crawlers
	o.mu.RUnlock()

	results := make(map[string]CrawlResult)
	var wg sync.WaitGroup
	var mu sync.Mutex

	// Chạy mỗi crawler trong goroutine riêng
	for name, c := range crawlers {
		wg.Add(1)
		go func(crawlerName string, crawler crawler.Crawler) {
			defer wg.Done()

			startTime := time.Now()
			posts, err := crawler.Fetch()
			duration := time.Since(startTime)

			result := CrawlResult{
				CrawlerName: crawlerName,
				PostCount:   len(posts),
				Duration:    duration,
				Timestamp:   time.Now(),
			}

			if err != nil {
				result.Error = err.Error()
				mu.Lock()
				o.failures[crawlerName]++
				mu.Unlock()
				fmt.Printf("❌ [%s] Error: %v (took %v)\n", crawlerName, err, duration)
			} else {
				fmt.Printf("✅ [%s] Fetched %d posts (took %v)\n", crawlerName, len(posts), duration)
			}

			mu.Lock()
			results[crawlerName] = result
			mu.Unlock()
		}(name, c)
	}

	// Chờ tất cả goroutines hoàn tất
	wg.Wait()

	o.mu.Lock()
	o.runs++
	o.lastRun = time.Now()
	o.mu.Unlock()

	return results
}

// CrawlResult chứa kết quả từ một lần crawl
type CrawlResult struct {
	CrawlerName string
	PostCount   int
	Error       string
	Duration    time.Duration
	Timestamp   time.Time
}

// GetMetrics trả về metrics hiện tại
func (o *CrawlerOrchestrator) GetMetrics() map[string]interface{} {
	o.mu.RLock()
	defer o.mu.RUnlock()

	failures := make(map[string]int64)
	for name, count := range o.failures {
		failures[name] = count
	}

	return map[string]interface{}{
		"total_runs": o.runs,
		"failures":   failures,
		"last_run":   o.lastRun,
	}
}
