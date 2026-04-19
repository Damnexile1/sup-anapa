package logger

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type ZapLogger struct {
	mu   sync.Mutex
	file *os.File
}

func Init(logFile string) (*ZapLogger, error) {
	if err := os.MkdirAll(filepath.Dir(logFile), 0o755); err != nil {
		return nil, err
	}

	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, err
	}

	log.SetOutput(f)
	log.SetFlags(0)

	return &ZapLogger{file: f}, nil
}

func (l *ZapLogger) Close() error {
	if l == nil || l.file == nil {
		return nil
	}
	return l.file.Close()
}

func (l *ZapLogger) Info(msg string, fields map[string]interface{}) {
	l.write("info", msg, fields)
}

func (l *ZapLogger) Error(msg string, fields map[string]interface{}) {
	l.write("error", msg, fields)
}

func (l *ZapLogger) Fatal(msg string, fields map[string]interface{}) {
	l.write("fatal", msg, fields)
	os.Exit(1)
}

func (l *ZapLogger) write(level, msg string, fields map[string]interface{}) {
	if l == nil {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()

	payload := map[string]interface{}{
		"ts":    time.Now().Format(time.RFC3339Nano),
		"level": level,
		"msg":   msg,
	}
	for k, v := range fields {
		payload[k] = v
	}

	b, _ := json.Marshal(payload)
	log.Print(string(b))
}
