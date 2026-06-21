package test

import (
	"fmt"
	"testing"
	"time"
)

func TestCreateEnvironment_Success(t *testing.T) {
	createBody := map[string]string{"name": "测试项目"}
	resp, _ := makeRequest("POST", "/api/v1/projects", createBody, true)
	data, _ := parseResponse(resp)

	if projectData, ok := data["data"].(map[string]interface{}); ok {
		if pid, ok := projectData["id"].(float64); ok {
			envBody := map[string]string{
				"name":        "测试环境_" + time.Now().Format("150405"),
				"description": "开发环境",
			}
			envResp, err := makeRequest("POST", fmt.Sprintf("/api/v1/projects/%d/environments", int(pid)), envBody, true)
			if err != nil {
				t.Fatalf("request failed: %v", err)
			}

			if envResp.StatusCode != 200 {
				t.Errorf("status code = %v, want 200", envResp.StatusCode)
			}
		}
	}
}

func TestListEnvironments_Success(t *testing.T) {
	createBody := map[string]string{"name": "环境测试项目"}
	resp, _ := makeRequest("POST", "/api/v1/projects", createBody, true)
	data, _ := parseResponse(resp)

	if projectData, ok := data["data"].(map[string]interface{}); ok {
		if pid, ok := projectData["id"].(float64); ok {
			listResp, err := makeRequest("GET", fmt.Sprintf("/api/v1/projects/%d/environments", int(pid)), nil, true)
			if err != nil {
				t.Fatalf("request failed: %v", err)
			}

			if listResp.StatusCode != 200 {
				t.Errorf("status code = %v, want 200", listResp.StatusCode)
			}
		}
	}
}

func TestSaveEnvVars_Success(t *testing.T) {
	createBody := map[string]string{"name": "变量测试项目"}
	resp, _ := makeRequest("POST", "/api/v1/projects", createBody, true)
	data, _ := parseResponse(resp)

	if projectData, ok := data["data"].(map[string]interface{}); ok {
		if pid, ok := projectData["id"].(float64); ok {
			envBody := map[string]string{"name": "变量环境"}
			envResp, _ := makeRequest("POST", fmt.Sprintf("/api/v1/projects/%d/environments", int(pid)), envBody, true)
			envData, _ := parseResponse(envResp)

			if envObj, ok := envData["data"].(map[string]interface{}); ok {
				if eid, ok := envObj["id"].(float64); ok {
					varsBody := []map[string]string{
						{"key": "host", "value": "localhost"},
						{"key": "port", "value": "8080"},
					}
					varsResp, err := makeRequest("PUT", fmt.Sprintf("/api/v1/environments/%d/variables", int(eid)), varsBody, true)
					if err != nil {
						t.Fatalf("request failed: %v", err)
					}

					if varsResp.StatusCode != 200 {
						t.Errorf("status code = %v, want 200", varsResp.StatusCode)
					}
				}
			}
		}
	}
}
