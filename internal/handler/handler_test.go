package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	return r
}

func TestHealth(t *testing.T) {
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

func TestRegister_InvalidBody(t *testing.T) {
	r := setupRouter()
	r.POST("/auth/register", Register)

	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBufferString("invalid"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status code = %v, want %v", w.Code, http.StatusBadRequest)
	}
}

func TestLogin_InvalidBody(t *testing.T) {
	r := setupRouter()
	r.POST("/auth/login", Login)

	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBufferString("invalid"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status code = %v, want %v", w.Code, http.StatusBadRequest)
	}
}

func TestGetProfile_NoAuth(t *testing.T) {
	r := setupRouter()
	r.GET("/user/profile", GetProfile)

	req, _ := http.NewRequest("GET", "/user/profile", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status code = %v, want %v", w.Code, http.StatusUnauthorized)
	}
}

func TestDebugAPI_InvalidID(t *testing.T) {
	r := setupRouter()
	r.POST("/apis/:id/debug", DebugAPI)

	req, _ := http.NewRequest("POST", "/apis/abc/debug", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status code = %v, want %v", w.Code, http.StatusBadRequest)
	}
}

func TestGetDynamicFunctions(t *testing.T) {
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
	if !ok || len(data) == 0 {
		t.Error("expected non-empty data array")
	}
}

func TestExportTestReport_InvalidID(t *testing.T) {
	r := setupRouter()
	r.GET("/test-runs/:id/export", ExportTestReport)

	req, _ := http.NewRequest("GET", "/test-runs/abc/export", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status code = %v, want %v", w.Code, http.StatusBadRequest)
	}
}
