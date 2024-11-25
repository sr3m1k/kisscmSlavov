package logger

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Username  string    `json:"username"`
	Command   string    `json:"command"`
	Args      []string  `json:"args"`
}

type Logger struct {
	mu      sync.Mutex
	logFile string
	entries []LogEntry
}

func NewLogger(logFile string) *Logger {
	return &Logger{
		logFile: logFile,
		entries: []LogEntry{},
	}
}

func (l *Logger) Log(username, command string, args []string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	entry := LogEntry{
		Timestamp: time.Now(),
		Username:  username,
		Command:   command,
		Args:      args,
	}
	l.entries = append(l.entries, entry)

	return l.writeToFile()
}

func (l *Logger) writeToFile() error {
	file, err := json.MarshalIndent(l.entries, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(l.logFile, file, 0644)
}
