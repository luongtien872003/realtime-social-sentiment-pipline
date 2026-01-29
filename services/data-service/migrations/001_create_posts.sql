-- =====================================================
-- MIGRATION: Tạo bảng posts
-- =====================================================
-- Mô tả: Tạo table để lưu trữ tất cả social media posts
-- File này sẽ tự động chạy khi PostgreSQL khởi động lần đầu
-- =====================================================

-- Tạo extension cho UUID (tùy chọn)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- =====================================================
-- BẢNG POSTS - Lưu trữ các bài viết social media
-- =====================================================
CREATE TABLE IF NOT EXISTS posts (
    -- ID duy nhất cho mỗi post
    id VARCHAR(50) PRIMARY KEY,
    
    -- Thông tin tác giả
    author VARCHAR(255) NOT NULL,
    
    -- Nội dung bài viết
    content TEXT NOT NULL,
    
    -- Chủ đề: ai, cloud, devops, programming, startup
    topic VARCHAR(50) NOT NULL,
    
    -- Cảm xúc: positive, negative, neutral
    sentiment VARCHAR(20) NOT NULL,
    
    -- Metrics
    likes INTEGER DEFAULT 0,
    comments INTEGER DEFAULT 0,
    shares INTEGER DEFAULT 0,
    
    -- Nền tảng: twitter, linkedin, reddit, hackernews
    platform VARCHAR(50) NOT NULL,
    
    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    processed_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- =====================================================
-- INDEXES - Tăng tốc query
-- =====================================================

-- Index cho topic (để lọc theo chủ đề)
CREATE INDEX IF NOT EXISTS idx_posts_topic ON posts(topic);

-- Index cho sentiment (để lọc theo cảm xúc)
CREATE INDEX IF NOT EXISTS idx_posts_sentiment ON posts(sentiment);

-- Index cho platform (để lọc theo nền tảng)
CREATE INDEX IF NOT EXISTS idx_posts_platform ON posts(platform);

-- Index cho created_at (để sắp xếp theo thời gian)
CREATE INDEX IF NOT EXISTS idx_posts_created_at ON posts(created_at DESC);

-- =====================================================
-- BẢNG ANALYTICS - Lưu kết quả phân tích từ Spark
-- =====================================================
CREATE TABLE IF NOT EXISTS analytics (
    id SERIAL PRIMARY KEY,
    
    -- Loại phân tích: hourly, daily, weekly
    analysis_type VARCHAR(50) NOT NULL,
    
    -- Thời gian phân tích
    period_start TIMESTAMP WITH TIME ZONE NOT NULL,
    period_end TIMESTAMP WITH TIME ZONE NOT NULL,
    
    -- Kết quả dạng JSON
    -- Ví dụ: {"total_posts": 1000, "positive": 400, "negative": 200}
    metrics JSONB NOT NULL,
    
    -- Timestamp
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Index cho analytics
CREATE INDEX IF NOT EXISTS idx_analytics_type ON analytics(analysis_type);
CREATE INDEX IF NOT EXISTS idx_analytics_period ON analytics(period_start DESC);

-- =====================================================
-- IN RA THÔNG BÁO KHI MIGRATION HOÀN TẤT
-- =====================================================
DO $$
BEGIN
    RAISE NOTICE '✅ Migration completed: Tables posts and analytics created successfully!';
END $$;
