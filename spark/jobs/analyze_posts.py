# =====================================================
# SPARK STREAMING JOB - Phân tích realtime
# =====================================================
# Mô tả: Spark job chạy LIÊN TỤC, phân tích dữ liệu mỗi 30 giây
# Mode: Micro-batch streaming (đọc PostgreSQL định kỳ)
#
# Cách chạy: python spark/jobs/analyze_posts.py
# Dừng: Ctrl+C
# =====================================================

from datetime import datetime
import time
import json
import signal
import sys

# Flag để dừng gracefully
running = True

def signal_handler(sig, frame):
    global running
    print("\n⚠️  Nhận tín hiệu dừng, đang shutdown...")
    running = False

signal.signal(signal.SIGINT, signal_handler)
signal.signal(signal.SIGTERM, signal_handler)

# =====================================================
# CẤU HÌNH
# =====================================================

# PostgreSQL config
PG_CONFIG = {
    "host": "localhost",
    "port": "5432",
    "database": "social_insight",
    "user": "postgres",
    "password": "postgres123"
}

# Streaming config
BATCH_INTERVAL = 30  # Xử lý mỗi 30 giây
OUTPUT_DIR = "analytics_output"

# =====================================================
# DATABASE CONNECTION
# =====================================================

import psycopg2
from psycopg2.extras import RealDictCursor

def get_db_connection():
    """Tạo kết nối database"""
    return psycopg2.connect(
        host=PG_CONFIG["host"],
        port=PG_CONFIG["port"],
        database=PG_CONFIG["database"],
        user=PG_CONFIG["user"],
        password=PG_CONFIG["password"]
    )

# =====================================================
# ANALYTICS FUNCTIONS
# =====================================================

def get_overall_stats(conn):
    """Lấy thống kê tổng quan"""
    with conn.cursor(cursor_factory=RealDictCursor) as cur:
        cur.execute("""
            SELECT 
                COUNT(*) as total_posts,
                COALESCE(SUM(likes), 0) as total_likes,
                COALESCE(SUM(comments), 0) as total_comments,
                COALESCE(AVG(likes), 0) as avg_likes
            FROM posts
        """)
        return dict(cur.fetchone())

def get_topic_stats(conn):
    """Thống kê theo topic"""
    with conn.cursor(cursor_factory=RealDictCursor) as cur:
        cur.execute("""
            SELECT topic, COUNT(*) as count, SUM(likes) as total_likes
            FROM posts
            GROUP BY topic
            ORDER BY count DESC
        """)
        return {row['topic']: {'count': row['count'], 'likes': row['total_likes']} for row in cur.fetchall()}

def get_sentiment_stats(conn):
    """Thống kê theo sentiment"""
    with conn.cursor(cursor_factory=RealDictCursor) as cur:
        cur.execute("""
            SELECT sentiment, COUNT(*) as count
            FROM posts
            GROUP BY sentiment
        """)
        total = sum(row['count'] for row in cur.fetchall())
        
        cur.execute("""
            SELECT sentiment, COUNT(*) as count
            FROM posts
            GROUP BY sentiment
        """)
        return {
            row['sentiment']: {
                'count': row['count'],
                'percentage': round(row['count'] / total * 100, 1) if total > 0 else 0
            }
            for row in cur.fetchall()
        }

def get_platform_stats(conn):
    """Thống kê theo platform"""
    with conn.cursor(cursor_factory=RealDictCursor) as cur:
        cur.execute("""
            SELECT platform, COUNT(*) as count, SUM(likes) as total_likes
            FROM posts
            GROUP BY platform
            ORDER BY count DESC
        """)
        return {row['platform']: {'count': row['count'], 'likes': row['total_likes']} for row in cur.fetchall()}

def get_recent_trend(conn, minutes=5):
    """Trend trong N phút gần nhất"""
    with conn.cursor(cursor_factory=RealDictCursor) as cur:
        cur.execute("""
            SELECT 
                COUNT(*) as posts_count,
                COALESCE(SUM(likes), 0) as total_likes,
                COALESCE(AVG(likes), 0) as avg_likes
            FROM posts
            WHERE created_at >= NOW() - INTERVAL '%s minutes'
        """, (minutes,))
        return dict(cur.fetchone())

def save_analytics(analytics, batch_num):
    """Lưu kết quả phân tích ra JSON"""
    import os
    os.makedirs(OUTPUT_DIR, exist_ok=True)
    
    filename = f"{OUTPUT_DIR}/analytics_{batch_num:06d}.json"
    with open(filename, 'w', encoding='utf-8') as f:
        json.dump(analytics, f, indent=2, ensure_ascii=False)
    
    return filename

# =====================================================
# STREAMING LOOP
# =====================================================

def run_streaming():
    """Main streaming loop"""
    print("╔════════════════════════════════════════════════════════════╗")
    print("║     SOCIAL INSIGHT - SPARK STREAMING ANALYTICS            ║")
    print("║     Phân tích realtime mỗi 30 giây                        ║")
    print("╚════════════════════════════════════════════════════════════╝")
    print()
    print(f"📊 Batch interval: {BATCH_INTERVAL} giây")
    print(f"📁 Output directory: {OUTPUT_DIR}/")
    print("   Nhấn Ctrl+C để dừng")
    print()
    print("════════════════════════════════════════════════════════════")
    
    batch_num = 0
    start_time = datetime.now()
    
    while running:
        batch_num += 1
        batch_start = datetime.now()
        
        try:
            # Kết nối database
            conn = get_db_connection()
            
            # Thu thập analytics
            analytics = {
                "batch_number": batch_num,
                "timestamp": datetime.now().isoformat(),
                "overall": get_overall_stats(conn),
                "by_topic": get_topic_stats(conn),
                "by_sentiment": get_sentiment_stats(conn),
                "by_platform": get_platform_stats(conn),
                "recent_5min": get_recent_trend(conn, 5)
            }
            
            conn.close()
            
            # Lưu kết quả
            filename = save_analytics(analytics, batch_num)
            
            # In summary
            total = analytics['overall']['total_posts']
            recent = analytics['recent_5min']['posts_count']
            
            print(f"\n📈 [Batch {batch_num}] {batch_start.strftime('%H:%M:%S')}")
            print(f"   Total posts: {total:,}")
            print(f"   Recent 5min: {recent:,} posts")
            
            if analytics['by_sentiment']:
                sentiments = analytics['by_sentiment']
                pos = sentiments.get('positive', {}).get('percentage', 0)
                neg = sentiments.get('negative', {}).get('percentage', 0)
                print(f"   Sentiment: 😊 {pos}% | 😢 {neg}%")
            
            print(f"   Saved: {filename}")
            
        except Exception as e:
            print(f"❌ [Batch {batch_num}] Error: {e}")
        
        # Đợi đến batch tiếp theo
        elapsed = (datetime.now() - batch_start).total_seconds()
        sleep_time = max(0, BATCH_INTERVAL - elapsed)
        
        if running and sleep_time > 0:
            time.sleep(sleep_time)
    
    # Kết thúc
    total_time = datetime.now() - start_time
    print()
    print("════════════════════════════════════════════════════════════")
    print(f"✅ Streaming kết thúc sau {batch_num} batches")
    print(f"   Thời gian chạy: {total_time}")
    print("════════════════════════════════════════════════════════════")

# =====================================================
# MAIN
# =====================================================

if __name__ == "__main__":
    try:
        # Test connection trước
        print("📡 Kiểm tra kết nối PostgreSQL...")
        conn = get_db_connection()
        conn.close()
        print("✅ Kết nối thành công!\n")
        
        run_streaming()
        
    except psycopg2.OperationalError as e:
        print(f"❌ Không thể kết nối PostgreSQL: {e}")
        print("   Hãy chắc chắn PostgreSQL đang chạy: docker-compose up -d")
        sys.exit(1)
