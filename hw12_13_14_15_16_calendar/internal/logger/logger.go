// Package logger предоставляет простую реализацию логгера с уровнями логирования.
// Поддерживает уровни: ERROR, WARN, INFO, DEBUG.
package logger

import (
	"fmt"
	"strings"
)

type Level int

const (
	ErrorLevel Level = iota // 0 - только ошибки
	WarnLevel               // 1 - предупреждения и ошибки
	InfoLevel               // 2 - информация, предупреждения и ошибки
	DebugLevel              // 3 - все сообщения включая отладочные
)

func parseLevel(level string) Level {
	switch strings.ToLower(level) {
	case "error":
		return ErrorLevel
	case "warn":
		return WarnLevel
	case "info":
		return InfoLevel
	case "debug":
		return DebugLevel
	default:
		return InfoLevel
	}
}

type Logger struct {
	level Level
}

func New(level string) *Logger {
	return &Logger{level: parseLevel(level)}
}

func (l *Logger) Error(msg string) {
	if l.level >= ErrorLevel {
		fmt.Println("[ERROR]", msg)
	}
}

func (l *Logger) Warn(msg string) {
	if l.level >= WarnLevel {
		fmt.Println("[WARN]", msg)
	}
}

func (l *Logger) Info(msg string) {
	if l.level >= InfoLevel {
		fmt.Println("[INFO]", msg)
	}
}

func (l *Logger) Debug(msg string) {
	if l.level >= DebugLevel {
		fmt.Println("[DEBUG]", msg)
	}
}
