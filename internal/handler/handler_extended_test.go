package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetDynamicFunctions_List(t *testing.T) {
	r := setupRouter()
	r.GET("/variables/dynamic", GetDynamicFunctions)

	req, _ := http.NewRequest("GET", "/variables/dynamic", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status code = %v, want %v", w.Code, http.StatusOK)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	data, ok := resp["data"].([]interface{})
	if !ok {
		t.Fatal("data should be an array")
	}

	expectedFuncs := []string{"timestamp", "unix", "uuid", "random_string", "random_int", "date", "datetime"}
	foundFuncs := make(map[string]bool)
	for _, item := range data {
		fn, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		if name, ok := fn["name"].(string); ok {
			foundFuncs[name] = true
		}
	}

	for _, name := range expectedFuncs {
		if !foundFuncs[name] {
			t.Errorf("function %v not found in response", name)
		}
	}
}

func TestHealthCheck(t *testing.T) {
	r := setupRouter()
	r.GET("/health", Health)

	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status code = %v, want %v", w.Code, http.StatusOK)
	}

	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)

	if resp["status"] != "ok" {
		t.Errorf("status = %v, want ok", resp["status"])
	}
}
