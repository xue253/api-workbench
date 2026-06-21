package engine

import (
	"testing"
)

func TestWSClientManagement(t *testing.T) {
	runID1 := uint(1)
	runID2 := uint(2)

	wsMu.Lock()
	wsClients[runID1] = nil
	wsClients[runID2] = nil
	wsMu.Unlock()

	wsMu.RLock()
	count1 := len(wsClients[runID1])
	count2 := len(wsClients[runID2])
	wsMu.RUnlock()

	if count1 != 0 {
		t.Errorf("runID1 clients = %v, want 0", count1)
	}
	if count2 != 0 {
		t.Errorf("runID2 clients = %v, want 0", count2)
	}

	wsMu.Lock()
	delete(wsClients, runID1)
	wsMu.Unlock()

	wsMu.RLock()
	_, exists := wsClients[runID1]
	wsMu.RUnlock()

	if exists {
		t.Error("runID1 should be deleted")
	}
}

func TestWSMessageFormat(t *testing.T) {
	msg := WSMessage{
		Type:  "progress",
		RunID: 1,
		Data: map[string]interface{}{
			"total":  10,
			"passed": 5,
			"failed": 5,
		},
	}

	if msg.Type != "progress" {
		t.Errorf("Type = %v, want progress", msg.Type)
	}

	if msg.RunID != 1 {
		t.Errorf("RunID = %v, want 1", msg.RunID)
	}

	data, ok := msg.Data.(map[string]interface{})
	if !ok {
		t.Error("Data should be map[string]interface{}")
	}

	if data["total"] != 10 {
		t.Errorf("total = %v, want 10", data["total"])
	}
}

func TestReplaceVars(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		envVars  map[string]string
		expected string
	}{
		{
			name:     "替换变量",
			text:     "http://{{host}}/api",
			envVars:  map[string]string{"host": "localhost"},
			expected: "http://localhost/api",
		},
		{
			name:     "无变量",
			text:     "http://example.com",
			envVars:  nil,
			expected: "http://example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ReplaceVars(tt.text, tt.envVars)
			if result != tt.expected {
				t.Errorf("ReplaceVars() = %v, want %v", result, tt.expected)
			}
		})
	}
}
