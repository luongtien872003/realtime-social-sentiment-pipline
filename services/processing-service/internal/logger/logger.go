package logger

import (
	"fmt"
	"time"
)

// Level định nghĩa mức độ logging
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

var levelNames = []string{"DEBUG", "INFO", "WARN", "ERROR"}

// Logger là wrapper cho structured logging
type Logger struct {
	name     string
	minLevel Level
}

// New tạo logger mới
func New(name string) *Logger {
	return &Logger{
		name:     name,
		minLevel: LevelDebug,
	}
}

// SetLevel đặt mức độ minimum logging
func (l *Logger) SetLevel(level Level) {
	l.minLevel = level
}

// log ghi log với context
func (l *Logger) log(level Level, msg string, fields map[string]interface{}) {
	if level < l.minLevel {
		return
	}

	timestamp := time.Now().Format("15:04:05.000")
	fmt.Printf("[%s] %s [%s] %s", timestamp, levelNames[level], l.name, msg)

	if len(fields) > 0 {
		fmt.Print(" {")
		first := true
		for k, v := range fields {
			if !first {
				fmt.Print(", ")
			}
			fmt.Printf("%s=%v", k, v)
			first = false
		}
		fmt.Print("}")
	}
	fmt.Println()
}

// Debug logs at DEBUG level
func (l *Logger) Debug(msg string, fields map[string]interface{}) {
	l.log(LevelDebug, msg, fields)
}

// Info logs at INFO level
func (l *Logger) Info(msg string, fields map[string]interface{}) {
	l.log(LevelInfo, msg, fields)
}

// Warn logs at WARN level
func (l *Logger) Warn(msg string, fields map[string]interface{}) {
	l.log(LevelWarn, msg, fields)
}

// Error logs at ERROR level
func (l *Logger) Error(msg string, fields map[string]interface{}) {
	l.log(LevelError, msg, fields)
}

// Shorthand methods
func (l *Logger) Debugf(msg string) {
	l.log(LevelDebug, msg, nil)
}

func (l *Logger) Infof(msg string) {
	l.log(LevelInfo, msg, nil)
}

func (l *Logger) Warnf(msg string) {
	l.log(LevelWarn, msg, nil)
}

func (l *Logger) Errorf(msg string) {
	l.log(LevelError, msg, nil)
}
