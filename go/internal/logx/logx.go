package logx

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

type Logger struct {
	level Level
}

func New(level string) Logger {
	return Logger{level: ParseLevel(level)}
}

func ParseLevel(value string) Level {
	level := strings.ToLower(strings.TrimSpace(value))
	switch level {
	case "debug":
		return LevelDebug
	case "info":
		return LevelInfo
	case "warn", "warning":
		return LevelWarn
	case "error":
		return LevelError
	default:
		return LevelError
	}
}

func (l Logger) Debugf(format string, args ...any) {
	if l.level <= LevelDebug {
		writef(os.Stdout, "debug", format, args...)
	}
}

func (l Logger) Infof(format string, args ...any) {
	if l.level <= LevelInfo {
		writef(os.Stdout, "info", format, args...)
	}
}

func (l Logger) Warnf(format string, args ...any) {
	if l.level <= LevelWarn {
		writef(os.Stderr, "warn", format, args...)
	}
}

func (l Logger) Errorf(format string, args ...any) {
	if l.level <= LevelError {
		writef(os.Stderr, "error", format, args...)
	}
}

func writef(out io.Writer, prefix string, format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	fmt.Fprintf(out, "%s: %s\n", prefix, message)
}
