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
	level, ok := NormalizeLevel(value)
	if !ok {
		return LevelError
	}
	switch level {
	case "debug":
		return LevelDebug
	case "info":
		return LevelInfo
	case "warn":
		return LevelWarn
	default:
		return LevelError
	}
}

func NormalizeLevel(value string) (string, bool) {
	level := strings.ToLower(strings.TrimSpace(value))
	if level == "warning" {
		level = "warn"
	}
	switch level {
	case "debug", "info", "warn", "error":
		return level, true
	default:
		return "", false
	}
}

func ResolveLevel(cliLevel string, cliSet bool, trailerLevel string, envLevel string, defaultLevel string) string {
	if cliSet {
		if level, ok := NormalizeLevel(cliLevel); ok {
			return level
		}
	}
	if level, ok := NormalizeLevel(trailerLevel); ok {
		return level
	}
	if level, ok := NormalizeLevel(envLevel); ok {
		return level
	}
	if level, ok := NormalizeLevel(defaultLevel); ok {
		return level
	}
	return "error"
}

func (l Logger) Debugf(format string, args ...any) {
	l.Logf(LevelDebug, format, args...)
}

func (l Logger) Infof(format string, args ...any) {
	l.Logf(LevelInfo, format, args...)
}

func (l Logger) Warnf(format string, args ...any) {
	l.Logf(LevelWarn, format, args...)
}

func (l Logger) Errorf(format string, args ...any) {
	l.Logf(LevelError, format, args...)
}

func (l Logger) Logf(level Level, format string, args ...any) {
	if l.level > level {
		return
	}
	prefix := levelPrefix(level)
	if level >= LevelWarn {
		writef(os.Stderr, prefix, format, args...)
		return
	}
	writef(os.Stdout, prefix, format, args...)
}

type bufferedEntry struct {
	level  Level
	format string
	args   []any
}

type BufferedLogger struct {
	entries []bufferedEntry
}

func NewBufferedLogger() *BufferedLogger {
	return &BufferedLogger{entries: []bufferedEntry{}}
}

func (b *BufferedLogger) Debugf(format string, args ...any) {
	b.entries = append(b.entries, bufferedEntry{level: LevelDebug, format: format, args: args})
}

func (b *BufferedLogger) Infof(format string, args ...any) {
	b.entries = append(b.entries, bufferedEntry{level: LevelInfo, format: format, args: args})
}

func (b *BufferedLogger) Warnf(format string, args ...any) {
	b.entries = append(b.entries, bufferedEntry{level: LevelWarn, format: format, args: args})
}

func (b *BufferedLogger) Errorf(format string, args ...any) {
	b.entries = append(b.entries, bufferedEntry{level: LevelError, format: format, args: args})
}

func (b *BufferedLogger) Flush(logger Logger) {
	for _, entry := range b.entries {
		logger.Logf(entry.level, entry.format, entry.args...)
	}
	b.entries = nil
}

func writef(out io.Writer, prefix string, format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	fmt.Fprintf(out, "%s: %s\n", prefix, message)
}

func levelPrefix(level Level) string {
	switch level {
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelWarn:
		return "warn"
	default:
		return "error"
	}
}
