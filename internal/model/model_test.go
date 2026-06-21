package model

import (
	"encoding/json"
	"testing"
	"time"
)

func TestUserJSON(t *testing.T) {
	user := User{
		ID:        1,
		Username:  "testuser",
		Email:     "test@example.com",
		Avatar:    "https://example.com/avatar.jpg",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	data, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var decoded User
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if decoded.ID != 1 {
		t.Errorf("ID = %v, want 1", decoded.ID)
	}
	if decoded.Username != "testuser" {
		t.Errorf("Username = %v, want testuser", decoded.Username)
	}
	if decoded.Email != "test@example.com" {
		t.Errorf("Email = %v, want test@example.com", decoded.Email)
	}

	var raw map[string]interface{}
	json.Unmarshal(data, &raw)
	if _, exists := raw["password"]; exists {
		t.Error("password should not be in JSON output")
	}
}

func TestProjectJSON(t *testing.T) {
	project := Project{
		ID:          1,
		UserID:      10,
		Name:        "My Project",
		Description: "A test project",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	data, err := json.Marshal(project)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var decoded Project
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if decoded.Name != "My Project" {
		t.Errorf("Name = %v, want My Project", decoded.Name)
	}
	if decoded.UserID != 10 {
		t.Errorf("UserID = %v, want 10", decoded.UserID)
	}
}

func TestEnvironmentJSON(t *testing.T) {
	env := Environment{
		ID:          1,
		ProjectID:   10,
		Name:        "Production",
		Description: "Production environment",
		SortOrder:   1,
	}

	data, err := json.Marshal(env)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var decoded Environment
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if decoded.Name != "Production" {
		t.Errorf("Name = %v, want Production", decoded.Name)
	}
}

func TestAPIJSON(t *testing.T) {
	api := API{
		ID:             1,
		CollectionID:   10,
		Name:           "Get Users",
		Description:    "Get all users",
		Protocol:       "http",
		Method:         "GET",
		URL:            "https://api.example.com/users",
		Headers:        `{"Authorization": "Bearer token"}`,
		QueryParams:    `{"page": "1"}`,
		BodyType:       "json",
		Body:           "",
		TimeoutMs:      30000,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	data, err := json.Marshal(api)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var decoded API
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if decoded.Method != "GET" {
		t.Errorf("Method = %v, want GET", decoded.Method)
	}
	if decoded.Protocol != "http" {
		t.Errorf("Protocol = %v, want http", decoded.Protocol)
	}
	if decoded.TimeoutMs != 30000 {
		t.Errorf("TimeoutMs = %v, want 30000", decoded.TimeoutMs)
	}
}

func TestAssertionJSON(t *testing.T) {
	assertion := Assertion{
		ID:         1,
		APIID:      10,
		TargetType: "status_code",
		Operator:   "equals",
		Path:       "",
		Expected:   "200",
		Enabled:    true,
	}

	data, err := json.Marshal(assertion)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var decoded Assertion
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if decoded.TargetType != "status_code" {
		t.Errorf("TargetType = %v, want status_code", decoded.TargetType)
	}
	if !decoded.Enabled {
		t.Error("Enabled should be true")
	}
}

func TestTestCaseJSON(t *testing.T) {
	tc := TestCase{
		ID:          1,
		ProjectID:   10,
		Name:        "Login Test",
		Description: "Test login flow",
		CreatedAt:   time.Now(),
	}

	data, err := json.Marshal(tc)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var decoded TestCase
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if decoded.Name != "Login Test" {
		t.Errorf("Name = %v, want Login Test", decoded.Name)
	}
}

func TestTestSuiteJSON(t *testing.T) {
	ts := TestSuite{
		ID:             1,
		ProjectID:      10,
		Name:           "Full Suite",
		Description:    "Complete test suite",
		RunMode:        "sequential",
		MaxConcurrency: 5,
		CreatedAt:      time.Now(),
	}

	data, err := json.Marshal(ts)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var decoded TestSuite
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if decoded.RunMode != "sequential" {
		t.Errorf("RunMode = %v, want sequential", decoded.RunMode)
	}
	if decoded.MaxConcurrency != 5 {
		t.Errorf("MaxConcurrency = %v, want 5", decoded.MaxConcurrency)
	}
}

func TestScheduledTaskJSON(t *testing.T) {
	st := ScheduledTask{
		ID:            1,
		ProjectID:     10,
		TargetType:    "test_suite",
		TargetID:      20,
		CronExpr:      "0 9 * * *",
		Enabled:       true,
		EnvironmentID: 30,
	}

	data, err := json.Marshal(st)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var decoded ScheduledTask
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if decoded.CronExpr != "0 9 * * *" {
		t.Errorf("CronExpr = %v, want 0 9 * * *", decoded.CronExpr)
	}
	if !decoded.Enabled {
		t.Error("Enabled should be true")
	}
}

func TestTestRunJSON(t *testing.T) {
	now := time.Now()
	tr := TestRun{
		ID:            1,
		TargetType:    "test_case",
		TargetID:      10,
		EnvironmentID: 20,
		Status:        "done",
		TriggerType:   "manual",
		Total:         10,
		Passed:        8,
		Failed:        2,
		Skipped:       0,
		DurationMs:    5000,
		StartedAt:     now,
		FinishedAt:    &now,
	}

	data, err := json.Marshal(tr)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var decoded TestRun
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if decoded.Status != "done" {
		t.Errorf("Status = %v, want done", decoded.Status)
	}
	if decoded.Passed != 8 {
		t.Errorf("Passed = %v, want 8", decoded.Passed)
	}
	if decoded.Failed != 2 {
		t.Errorf("Failed = %v, want 2", decoded.Failed)
	}
}

func TestTestRunDetailJSON(t *testing.T) {
	now := time.Now()
	detail := TestRunDetail{
		ID:              1,
		TestRunID:       10,
		APIID:           20,
		TestCaseID:      30,
		DataIndex:       0,
		Status:          "passed",
		StatusCode:      200,
		ResponseHeaders: `{"Content-Type": "application/json"}`,
		ResponseBody:    `{"data": []}`,
		DurationMs:      123,
		ErrorMessage:    "",
		RetryCount:      0,
		ExecutedAt:      now,
	}

	data, err := json.Marshal(detail)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var decoded TestRunDetail
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if decoded.Status != "passed" {
		t.Errorf("Status = %v, want passed", decoded.Status)
	}
	if decoded.StatusCode != 200 {
		t.Errorf("StatusCode = %v, want 200", decoded.StatusCode)
	}
}
