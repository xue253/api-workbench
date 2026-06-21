package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSetup(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	Setup(r)

	routes := r.Routes()
	if len(routes) == 0 {
		t.Error("no routes registered")
	}
}

func TestRoutesExist(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	Setup(r)

	routes := r.Routes()

	routeMap := make(map[string]bool)
	for _, route := range routes {
		routeMap[route.Method+":"+route.Path] = true
	}

	expectedRoutes := []string{
		"POST:/api/v1/auth/register",
		"POST:/api/v1/auth/login",
		"GET:/api/v1/user/profile",
		"GET:/api/v1/projects",
		"POST:/api/v1/projects",
		"GET:/api/v1/test-runs",
	}

	for _, route := range expectedRoutes {
		if !routeMap[route] {
			t.Errorf("route %v not found", route)
		}
	}
}

func TestHealthEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	Setup(r)

	req, _ := http.NewRequest("GET", "/api/v1/health", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status code = %v, want %v", w.Code, http.StatusOK)
	}
}
