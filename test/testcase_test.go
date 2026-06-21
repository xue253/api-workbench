package test

import (
	"fmt"
	"testing"
	"time"
)

func TestCreateTestCase_Success(t *testing.T) {
	createBody := map[string]string{"name": "用例测试项目"}
	resp, _ := makeRequest("POST", "/api/v1/projects", createBody, true)
	data, _ := parseResponse(resp)

	if projectData, ok := data["data"].(map[string]interface{}); ok {
		if pid, ok := projectData["id"].(float64); ok {
			tcBody := map[string]string{
				"name":        "登录测试用例_" + time.Now().Format("150405"),
				"description": "测试用户登录流程",
			}
			tcResp, err := makeRequest("POST", fmt.Sprintf("/api/v1/projects/%d/test-cases", int(pid)), tcBody, true)
			if err != nil {
				t.Fatalf("request failed: %v", err)
			}

			if tcResp.StatusCode != 200 {
				t.Errorf("status code = %v, want 200", tcResp.StatusCode)
			}
		}
	}
}

func TestListTestCases_Success(t *testing.T) {
	createBody := map[string]string{"name": "用例列表项目"}
	resp, _ := makeRequest("POST", "/api/v1/projects", createBody, true)
	data, _ := parseResponse(resp)

	if projectData, ok := data["data"].(map[string]interface{}); ok {
		if pid, ok := projectData["id"].(float64); ok {
			listResp, err := makeRequest("GET", fmt.Sprintf("/api/v1/projects/%d/test-cases", int(pid)), nil, true)
			if err != nil {
				t.Fatalf("request failed: %v", err)
			}

			if listResp.StatusCode != 200 {
				t.Errorf("status code = %v, want 200", listResp.StatusCode)
			}
		}
	}
}

func TestCreateTestSuite_Success(t *testing.T) {
	createBody := map[string]string{"name": "套件测试项目"}
	resp, _ := makeRequest("POST", "/api/v1/projects", createBody, true)
	data, _ := parseResponse(resp)

	if projectData, ok := data["data"].(map[string]interface{}); ok {
		if pid, ok := projectData["id"].(float64); ok {
			tsBody := map[string]interface{}{
				"name":             "完整测试套件_" + time.Now().Format("150405"),
				"description":      "包含所有测试用例",
				"run_mode":         "sequential",
				"max_concurrency":  5,
			}
			tsResp, err := makeRequest("POST", fmt.Sprintf("/api/v1/projects/%d/test-suites", int(pid)), tsBody, true)
			if err != nil {
				t.Fatalf("request failed: %v", err)
			}

			if tsResp.StatusCode != 200 {
				t.Errorf("status code = %v, want 200", tsResp.StatusCode)
			}
		}
	}
}

func TestListTestSuites_Success(t *testing.T) {
	createBody := map[string]string{"name": "套件列表项目"}
	resp, _ := makeRequest("POST", "/api/v1/projects", createBody, true)
	data, _ := parseResponse(resp)

	if projectData, ok := data["data"].(map[string]interface{}); ok {
		if pid, ok := projectData["id"].(float64); ok {
			listResp, err := makeRequest("GET", fmt.Sprintf("/api/v1/projects/%d/test-suites", int(pid)), nil, true)
			if err != nil {
				t.Fatalf("request failed: %v", err)
			}

			if listResp.StatusCode != 200 {
				t.Errorf("status code = %v, want 200", listResp.StatusCode)
			}
		}
	}
}

func TestRunTestCase_Success(t *testing.T) {
	createBody := map[string]string{"name": "运行用例项目"}
	resp, _ := makeRequest("POST", "/api/v1/projects", createBody, true)
	data, _ := parseResponse(resp)

	if projectData, ok := data["data"].(map[string]interface{}); ok {
		if pid, ok := projectData["id"].(float64); ok {
			tcBody := map[string]string{"name": "待运行用例"}
			tcResp, _ := makeRequest("POST", fmt.Sprintf("/api/v1/projects/%d/test-cases", int(pid)), tcBody, true)
			tcData, _ := parseResponse(tcResp)

			if tcObj, ok := tcData["data"].(map[string]interface{}); ok {
				if tcid, ok := tcObj["id"].(float64); ok {
					runResp, err := makeRequest("POST", fmt.Sprintf("/api/v1/test-cases/%d/run", int(tcid)), map[string]interface{}{}, true)
					if err != nil {
						t.Fatalf("request failed: %v", err)
					}

					if runResp.StatusCode != 200 {
						t.Errorf("status code = %v, want 200", runResp.StatusCode)
					}
				}
			}
		}
	}
}
