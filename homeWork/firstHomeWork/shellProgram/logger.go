package shellProgram

import (
	"encoding/json"
	"fmt"
	"time"
)

type LogEntry struct {
	User      string `json:"user"`
	Command   string `json:"command"`
	Result    string `json:"result"`
	Timestamp string `json:"timestamp"`
}

func (s *Shell) logAction(command, result string) {
	entry := LogEntry{
		User:      s.Config,
		Command:   command,
		Result:    result,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	logData, err := json.Marshal(entry)
	if err != nil {
		fmt.Println("Error logging:", err)
		return
	}

	_, err = s.logFile.Write(logData)
	if err != nil {
		fmt.Println("Error writing log:", err)
	}
	s.logFile.WriteString("\n")
}
