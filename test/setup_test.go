package test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"api-workbench/internal/config"
	"api-workbench/internal/db"
	"api-workbench/internal/middleware"
	"api-workbench/internal/router"

	"github.com/gin-gonic/gin"
)

var (
	testServer *httptest.Server
	token      string
)

func setupTestEnvironment(t *testing.T) {
	gin.SetMode(gin.TestMode)

	config.AppConfig = config.Config{
		Server: config.ServerConfig{Port: 8080, Mode: "debug"},
		Database: config.DatabaseConfig{
			Host: "127.0.0.1", Port: 3306, User: "root",
			Password: "qazwsx123", DBName: "api_workbench", Charset: "utf8mb4",
		},
		JWT: config.JWTConfig{Secret: "api-workbench-jwt-secret-2026", ExpireHour: 168},
		Log: config.LogConfig{Level: "info", File: "logs/app.log"},
	}

	db.Init()

	r := gin.New()
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())
	router.Setup(r)

	testServer = httptest.NewServer(r)
}

func teardownTestEnvironment() {
	if testServer != nil {
		testServer.Close()
	}
}

func makeRequest(method, path string, body interface{}, auth bool) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonBytes, _ := json.Marshal(body)
		bodyReader = bytes.NewBuffer(jsonBytes)
	}

	req, err := http.NewRequest(method, testServer.URL+path, bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if auth && token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	return http.DefaultClient.Do(req)
}

func parseResponse(resp *http.Response) (map[string]interface{}, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	return result, err
}

func TestMain(m *testing.M) {
	setupTestEnvironment(nil)
	code := m.Run()
	teardownTestEnvironment()
	_ = code
}
