// =====================================================
// ENV LOADER - Load .env file
// =====================================================
// Mô tả: Load environment variables từ .env file
// Sử dụng thư viện godotenv
// =====================================================

package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// LoadEnvFile load .env file từ thư mục chỉ định
func LoadEnvFile(envPath string) error {
	// Nếu không chỉ định, tìm .env từ current directory
	if envPath == "" {
		envPath = ".env"
	}

	// Check if file exists
	if _, err := os.Stat(envPath); err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("⚠️  .env file not found at %s, using environment variables\n", envPath)
			return nil // Không phải lỗi - có thể được cấu hình bằng env vars
		}
		return fmt.Errorf("cannot access .env file: %w", err)
	}

	// Load .env file
	if err := godotenv.Load(envPath); err != nil {
		return fmt.Errorf("error loading .env file: %w", err)
	}

	fmt.Printf("✅ Loaded environment from %s\n", envPath)
	return nil
}

// LoadEnvFileFromDir find và load .env từ directory
func LoadEnvFileFromDir(dir string) error {
	envPath := filepath.Join(dir, ".env")
	return LoadEnvFile(envPath)
}
