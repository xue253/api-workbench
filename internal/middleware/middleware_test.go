package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"api-workbench/internal/config"

	"github.com/gin-gonic/gin"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	config.AppConfig.JWT.Secret = "api-workbench-jwt-secret-2026"
	return r
}

func TestCORS(t *testing.T) {
	r := setupTestRouter()
	r.Use(CORS())

	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf("Access-Control-Allow-Origin = %v, want *", w.Header().Get("Access-Control-Allow-Origin"))
	}

	if w.Header().Get("Access-Control-Allow-Methods") == "" {
		t.Error("Access-Control-Allow-Methods should not be empty")
	}

	if w.Header().Get("Access-Control-Allow-Headers") == "" {
		t.Error("Access-Control-Allow-Headers should not be empty")
	}
}

func TestCORS_Preflight(t *testing.T) {
	r := setupTestRouter()
	r.Use(CORS())

	r.OPTIONS("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	req, _ := http.NewRequest("OPTIONS", "/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	req.Header.Set("Access-Control-Request-Method", "POST")
	req.Header.Set("Access-Control-Request-Headers", "Content-Type")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent && w.Code != http.StatusOK {
		t.Errorf("status code = %v, want %v or %v", w.Code, http.StatusNoContent, http.StatusOK)
	}
}

func TestLogger(t *testing.T) {
	r := setupTestRouter()
	r.Use(Logger())

	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status code = %v, want %v", w.Code, http.StatusOK)
	}
}

func TestAuth_NoHeader(t *testing.T) {
	r := setupTestRouter()
	r.Use(Auth())

	r.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	req, _ := http.NewRequest("GET", "/protected", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status code = %v, want %v", w.Code, http.StatusUnauthorized)
	}
}

func TestAuth_InvalidFormat(t *testing.T) {
	r := setupTestRouter()
	r.Use(Auth())

	r.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "InvalidToken")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status code = %v, want %v", w.Code, http.StatusUnauthorized)
	}
}

func TestAuth_InvalidToken(t *testing.T) {
	r := setupTestRouter()
	r.Use(Auth())

	r.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status code = %v, want %v", w.Code, http.StatusUnauthorized)
	}
}

func TestAuth_ValidToken(t *testing.T) {
	token, err := generateTestToken(1)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	r := setupTestRouter()
	r.Use(Auth())

	r.GET("/protected", func(c *gin.Context) {
		uid, exists := c.Get("user_id")
		if !exists {
			t.Error("user_id should exist in context")
			return
		}
		if uid.(uint) != 1 {
			t.Errorf("user_id = %v, want 1", uid)
		}
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Logf("Response: %s", w.Body.String())
		t.Errorf("status code = %v, want %v", w.Code, http.StatusOK)
	}
}
