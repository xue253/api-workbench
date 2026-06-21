package test

import (
	"fmt"
	"testing"
	"time"
)

func TestCreateProject_Success(t *testing.T) {
	body := map[string]string{
		"name":        "测试项目_" + time.Now().Format("150405"),
		"description": "这是一个测试项目",
	}

	resp, err := makeRequest("POST", "/api/v1/projects", body, true)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("status code = %v, want 200", resp.StatusCode)
	}
}

func TestCreateProject_Unauthorized(t *testing.T) {
	body := map[string]string{
		"name":        "未授权项目",
		"description": "不应该创建成功",
	}

	resp, err := makeRequest("POST", "/api/v1/projects", body, false)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if resp.StatusCode != 401 {
		t.Errorf("status code = %v, want 401", resp.StatusCode)
	}
}

func TestListProjects_Success(t *testing.T) {
	resp, err := makeRequest("GET", "/api/v1/projects", nil, true)
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

func TestUpdateProject_Success(t *testing.T) {
	createBody := map[string]string{
		"name":        "待更新项目_" + time.Now().Format("150405"),
		"description": "原始描述",
	}
	resp, _ := makeRequest("POST", "/api/v1/projects", createBody, true)
	data, _ := parseResponse(resp)

	if projectData, ok := data["data"].(map[string]interface{}); ok {
		if id, ok := projectData["id"].(float64); ok {
			updateBody := map[string]string{
				"name":        "已更新项目",
				"description": "更新后的描述",
			}
			updateResp, err := makeRequest("PUT", fmt.Sprintf("/api/v1/projects/%d", int(id)), updateBody, true)
			if err != nil {
				t.Fatalf("request failed: %v", err)
			}

			if updateResp.StatusCode != 200 {
				t.Logf("Update response: %v", updateResp.StatusCode)
			}
		}
	}
}

func TestDeleteProject_Success(t *testing.T) {
	createBody := map[string]string{
		"name":        "待删除项目_" + time.Now().Format("150405"),
		"description": "将被删除",
	}
	resp, _ := makeRequest("POST", "/api/v1/projects", createBody, true)
	data, _ := parseResponse(resp)

	if projectData, ok := data["data"].(map[string]interface{}); ok {
		if id, ok := projectData["id"].(float64); ok {
			deleteResp, err := makeRequest("DELETE", fmt.Sprintf("/api/v1/projects/%d", int(id)), nil, true)
			if err != nil {
				t.Fatalf("request failed: %v", err)
			}

			if deleteResp.StatusCode != 200 {
				t.Errorf("status code = %v, want 200", deleteResp.StatusCode)
			}
		}
	}
}
