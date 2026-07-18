package logger

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

type Category string

const (
	CategoryApplication Category = "application"
	CategoryService     Category = "service"
	CategoryRuntime     Category = "runtime"
	CategoryPackage     Category = "package"
	CategoryProcess     Category = "process"
	CategoryConfig      Category = "configuration"
	CategoryDatabase    Category = "database"
)

type Logger struct {
	mu     sync.RWMutex
	level  Level
	output io.Writer
}

func New(level Level, output io.Writer) *Logger {
	if output == nil {
		output = os.Stderr
	}
	return &Logger{
		level:  level,
		output: output,
	}
}

func (l *Logger) SetLevel(level Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

func (l *Logger) Debug(category Category, msg string, args ...any) {
	l.log(LevelDebug, category, msg, args...)
}

func (l *Logger) Info(category Category, msg string, args ...any) {
	l.log(LevelInfo, category, msg, args...)
}

func (l *Logger) Warn(category Category, msg string, args ...any) {
	l.log(LevelWarn, category, msg, args...)
}

func (l *Logger) Error(category Category, msg string, args ...any) {
	l.log(LevelError, category, msg, args...)
}

func (l *Logger) log(level Level, category Category, msg string, args ...any) {
	l.mu.RLock()
	currentLevel := l.level
	l.mu.RUnlock()

	if level < currentLevel {
		return
	}

	var b strings.Builder
	fmt.Fprintf(&b, "[%s] [%s] %s", level, category, msg)

	if len(args) > 0 {
		for i := 0; i < len(args); i += 2 {
			b.WriteString(" ")
			if i+1 < len(args) {
				fmt.Fprintf(&b, "%v=%v", args[i], args[i+1])
			} else {
				fmt.Fprintf(&b, "%v", args[i])
			}
		}
	}

	b.WriteString("\n")

	l.mu.RLock()
	defer l.mu.RUnlock()
	_, _ = l.output.Write([]byte(b.String()))
}

func (l *Logger) WithPrefix(prefix string) *PrefixLogger {
	return &PrefixLogger{
		logger:   l,
		prefix:   prefix,
		category: Category(prefix),
	}
}

type PrefixLogger struct {
	logger   *Logger
	prefix   string
	category Category
}

func (p *PrefixLogger) Debug(msg string, args ...any) {
	p.logger.Debug(p.category, msg, args...)
}

func (p *PrefixLogger) Info(msg string, args ...any) {
	p.logger.Info(p.category, msg, args...)
}

func (p *PrefixLogger) Warn(msg string, args ...any) {
	p.logger.Warn(p.category, msg, args...)
}

func (p *PrefixLogger) Error(msg string, args ...any) {
	p.logger.Error(p.category, msg, args...)
}
