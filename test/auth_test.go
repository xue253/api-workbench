package test

import (
	"testing"
	"time"
)

func TestRegister_Success(t *testing.T) {
	body := map[string]string{
		"username": "testuser_" + time.Now().Format("150405"),
		"password": "123456",
		"email":    "test@example.com",
	}

	resp, err := makeRequest("POST", "/api/v1/auth/register", body, false)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("status code = %v, want 200", resp.StatusCode)
	}

	data, _ := parseResponse(resp)
	if _, ok := data["data"]; !ok {
		t.Error("response should contain data field")
	}
}

func TestRegister_DuplicateUsername(t *testing.T) {
	username := "duplicate_user_" + time.Now().Format("150405")

	body1 := map[string]string{"username": username, "password": "123456"}
	makeRequest("POST", "/api/v1/auth/register", body1, false)

	body2 := map[string]string{"username": username, "password": "123456"}
	resp, err := makeRequest("POST", "/api/v1/auth/register", body2, false)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if resp.StatusCode != 409 {
		t.Errorf("status code = %v, want 409", resp.StatusCode)
	}
}

func TestRegister_InvalidBody(t *testing.T) {
	body := map[string]string{}

	resp, err := makeRequest("POST", "/api/v1/auth/register", body, false)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if resp.StatusCode != 400 {
		t.Errorf("status code = %v, want 400", resp.StatusCode)
	}
}

func TestLogin_Success(t *testing.T) {
	username := "login_test_" + time.Now().Format("150405")
	regBody := map[string]string{"username": username, "password": "123456"}
	makeRequest("POST", "/api/v1/auth/register", regBody, false)

	loginBody := map[string]string{"username": username, "password": "123456"}
	resp, err := makeRequest("POST", "/api/v1/auth/login", loginBody, false)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("status code = %v, want 200", resp.StatusCode)
	}

	data, _ := parseResponse(resp)
	if data["data"] != nil {
		if tokenData, ok := data["data"].(map[string]interface{}); ok {
			if t, ok := tokenData["token"].(string); ok {
				token = t
			}
		}
	}

	if token == "" {
		t.Error("token should not be empty")
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	username := "login_wrong_" + time.Now().Format("150405")
	regBody := map[string]string{"username": username, "password": "123456"}
	makeRequest("POST", "/api/v1/auth/register", regBody, false)

	loginBody := map[string]string{"username": username, "password": "wrongpassword"}
	resp, err := makeRequest("POST", "/api/v1/auth/login", loginBody, false)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if resp.StatusCode != 401 {
		t.Errorf("status code = %v, want 401", resp.StatusCode)
	}
}

func TestLogin_NonexistentUser(t *testing.T) {
	loginBody := map[string]string{"username": "nonexistent_user", "password": "123456"}
	resp, err := makeRequest("POST", "/api/v1/auth/login", loginBody, false)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if resp.StatusCode != 401 {
		t.Errorf("status code = %v, want 401", resp.StatusCode)
	}
}
