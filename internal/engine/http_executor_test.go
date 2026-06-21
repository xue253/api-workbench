package engine

import (
	"encoding/json"
	"testing"
)

func TestExecuteHTTP_InvalidURL(t *testing.T) {
	req := &DebugRequest{
		Method:    "GET",
		URL:       "http://invalid-url-that-does-not-exist.local",
		TimeoutMs: 1000,
	}

	resp := ExecuteHTTP(req)

	if resp.Error == "" {
		t.Error("期望连接错误，但没有错误")
	}
}

func TestFormatJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "格式化JSON",
			input:    `{"key":"value","num":123}`,
			expected: "{\n  \"key\": \"value\",\n  \"num\": 123\n}",
		},
		{
			name:     "非JSON返回原文",
			input:    "not json",
			expected: "not json",
		},
		{
			name:     "已格式化JSON",
			input:    "{\n  \"key\": \"value\"\n}",
			expected: "{\n  \"key\": \"value\"\n}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatJSON(tt.input)
			if result != tt.expected {
				t.Errorf("FormatJSON() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDebugRequest_JSON(t *testing.T) {
	req := &DebugRequest{
		Method:      "POST",
		URL:         "http://example.com/api",
		Headers:     map[string]string{"Content-Type": "application/json"},
		QueryParams: map[string]string{"page": "1"},
		BodyType:    "json",
		Body:        `{"name":"test"}`,
		TimeoutMs:   5000,
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var decoded DebugRequest
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if decoded.Method != "POST" {
		t.Errorf("Method = %v, want POST", decoded.Method)
	}
	if decoded.URL != "http://example.com/api" {
		t.Errorf("URL = %v, want http://example.com/api", decoded.URL)
	}
	if decoded.Body != `{"name":"test"}` {
		t.Errorf("Body = %v, want {\"name\":\"test\"}", decoded.Body)
	}
}

func TestDebugResponse_JSON(t *testing.T) {
	resp := &DebugResponse{
		StatusCode: 200,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       `{"status":"ok"}`,
		DurationMs: 123,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var decoded DebugResponse
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if decoded.StatusCode != 200 {
		t.Errorf("StatusCode = %v, want 200", decoded.StatusCode)
	}
	if decoded.DurationMs != 123 {
		t.Errorf("DurationMs = %v, want 123", decoded.DurationMs)
	}
}
