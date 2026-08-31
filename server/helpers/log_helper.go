package helpers

import (
	"sync"
	"time"
)

type LogEntry struct {
	Index     int       `json:"index"`
	AgentID   string    `json:"agentId"`
	Hostname  string    `json:"hostname"`
	Command   string    `json:"command"`
	Response  string    `json:"response"`
	Timestamp time.Time `json:"timestamp"`
}

var (
	outputLog  []LogEntry
	outputMu   sync.Mutex
	logCounter int
)

func AppendLog(agentID, hostname, command, response string) {
	outputMu.Lock()
	defer outputMu.Unlock()
	logCounter++
	entry := LogEntry{
		Index:     logCounter,
		AgentID:   agentID,
		Hostname:  hostname,
		Command:   command,
		Response:  response,
		Timestamp: time.Now(),
	}
	if len(outputLog) >= 500 {
		outputLog = outputLog[1:]
	}
	outputLog = append(outputLog, entry)
}

func GetLogSince(afterIndex int) []LogEntry {
	outputMu.Lock()
	defer outputMu.Unlock()
	result := []LogEntry{}
	for _, e := range outputLog {
		if e.Index > afterIndex {
			result = append(result, e)
		}
	}
	return result
}
