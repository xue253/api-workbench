package variable

import (
	"testing"
)

func TestReplace(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		envVars  map[string]string
		tempVars map[string]string
		expected string
	}{
		{
			name:     "无变量",
			text:     "hello world",
			envVars:  nil,
			tempVars: nil,
			expected: "hello world",
		},
		{
			name:     "环境变量替换",
			text:     "http://{{host}}/api",
			envVars:  map[string]string{"host": "localhost:8080"},
			tempVars: nil,
			expected: "http://localhost:8080/api",
		},
		{
			name:     "临时变量优先",
			text:     "{{key}}",
			envVars:  map[string]string{"key": "env_value"},
			tempVars: map[string]string{"key": "temp_value"},
			expected: "temp_value",
		},
		{
			name:     "多个变量",
			text:     "{{proto}}://{{host}}:{{port}}/api",
			envVars:  map[string]string{"proto": "http", "host": "localhost", "port": "8080"},
			tempVars: nil,
			expected: "http://localhost:8080/api",
		},
		{
			name:     "变量不存在保留原样",
			text:     "http://{{host}}:{{port}}",
			envVars:  map[string]string{"host": "localhost"},
			tempVars: nil,
			expected: "http://localhost:{{port}}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Replace(tt.text, tt.envVars, tt.tempVars)
			if result != tt.expected {
				t.Errorf("Replace() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestReplaceMap(t *testing.T) {
	input := map[string]string{
		"Content-Type": "application/{{format}}",
		"Authorization": "Bearer {{token}}",
	}
	envVars := map[string]string{"format": "json", "token": "abc123"}

	result := ReplaceMap(input, envVars, nil)

	if result["Content-Type"] != "application/json" {
		t.Errorf("Content-Type = %v, want application/json", result["Content-Type"])
	}
	if result["Authorization"] != "Bearer abc123" {
		t.Errorf("Authorization = %v, want Bearer abc123", result["Authorization"])
	}
}

func TestGenerateDynamic(t *testing.T) {
	tests := []struct {
		name     string
		funcName string
		check    func(string) bool
	}{
		{"timestamp", "timestamp", func(s string) bool { return len(s) == 13 }},
		{"unix", "unix", func(s string) bool { return len(s) == 10 }},
		{"uuid", "uuid", func(s string) bool { return len(s) == 36 }},
		{"date", "date", func(s string) bool { return len(s) == 10 }},
		{"datetime", "datetime", func(s string) bool { return len(s) == 19 }},
		{"random_string", "random_string", func(s string) bool { return len(s) == 8 }},
		{"random_int", "random_int", func(s string) bool { return len(s) > 0 && len(s) <= 4 }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateDynamic(tt.funcName, "")
			if !tt.check(result) {
				t.Errorf("GenerateDynamic(%v) = %v, 验证失败", tt.funcName, result)
			}
		})
	}
}

func TestGenerateDynamicUnknown(t *testing.T) {
	result := GenerateDynamic("unknown_func", "")
	if result != "{{unknown_func}}" {
		t.Errorf("GenerateDynamic(unknown_func) = %v, want {{unknown_func}}", result)
	}
}
