package test

import (
	"fmt"
	"testing"
	"time"
)

func TestCreateCollection_Success(t *testing.T) {
	createBody := map[string]string{"name": "接口库测试项目"}
	resp, _ := makeRequest("POST", "/api/v1/projects", createBody, true)
	data, _ := parseResponse(resp)

	if projectData, ok := data["data"].(map[string]interface{}); ok {
		if pid, ok := projectData["id"].(float64); ok {
			colBody := map[string]string{
				"name":        "测试集合_" + time.Now().Format("150405"),
				"description": "用户相关接口",
			}
			colResp, err := makeRequest("POST", fmt.Sprintf("/api/v1/projects/%d/collections", int(pid)), colBody, true)
			if err != nil {
				t.Fatalf("request failed: %v", err)
			}

			if colResp.StatusCode != 200 {
				t.Errorf("status code = %v, want 200", colResp.StatusCode)
			}
		}
	}
}

func TestListCollections_Success(t *testing.T) {
	createBody := map[string]string{"name": "集合列表项目"}
	resp, _ := makeRequest("POST", "/api/v1/projects", createBody, true)
	data, _ := parseResponse(resp)

	if projectData, ok := data["data"].(map[string]interface{}); ok {
		if pid, ok := projectData["id"].(float64); ok {
			listResp, err := makeRequest("GET", fmt.Sprintf("/api/v1/projects/%d/collections", int(pid)), nil, true)
			if err != nil {
				t.Fatalf("request failed: %v", err)
			}

			if listResp.StatusCode != 200 {
				t.Errorf("status code = %v, want 200", listResp.StatusCode)
			}
		}
	}
}

func TestCreateAPI_Success(t *testing.T) {
	createBody := map[string]string{"name": "API测试项目"}
	resp, _ := makeRequest("POST", "/api/v1/projects", createBody, true)
	data, _ := parseResponse(resp)

	if projectData, ok := data["data"].(map[string]interface{}); ok {
		if pid, ok := projectData["id"].(float64); ok {
			colBody := map[string]string{"name": "API集合"}
			colResp, _ := makeRequest("POST", fmt.Sprintf("/api/v1/projects/%d/collections", int(pid)), colBody, true)
			colData, _ := parseResponse(colResp)

			if colObj, ok := colData["data"].(map[string]interface{}); ok {
				if cid, ok := colObj["id"].(float64); ok {
					apiBody := map[string]interface{}{
						"name":       "获取用户列表",
						"method":     "GET",
						"url":        "https://api.example.com/users",
						"protocol":   "http",
						"body_type":  "json",
						"timeout_ms": 30000,
					}
					apiResp, err := makeRequest("POST", fmt.Sprintf("/api/v1/collections/%d/apis", int(cid)), apiBody, true)
					if err != nil {
						t.Fatalf("request failed: %v", err)
					}

					if apiResp.StatusCode != 200 {
						t.Errorf("status code = %v, want 200", apiResp.StatusCode)
					}
				}
			}
		}
	}
}

func TestGetAPI_Success(t *testing.T) {
	createBody := map[string]string{"name": "GetAPI项目"}
	resp, _ := makeRequest("POST", "/api/v1/projects", createBody, true)
	data, _ := parseResponse(resp)

	if projectData, ok := data["data"].(map[string]interface{}); ok {
		if pid, ok := projectData["id"].(float64); ok {
			colBody := map[string]string{"name": "GetAPI集合"}
			colResp, _ := makeRequest("POST", fmt.Sprintf("/api/v1/projects/%d/collections", int(pid)), colBody, true)
			colData, _ := parseResponse(colResp)

			if colObj, ok := colData["data"].(map[string]interface{}); ok {
				if cid, ok := colObj["id"].(float64); ok {
					apiBody := map[string]interface{}{
						"name":       "获取用户",
						"method":     "GET",
						"url":        "https://api.example.com/user",
						"protocol":   "http",
						"body_type":  "json",
						"timeout_ms": 30000,
					}
					apiResp, _ := makeRequest("POST", fmt.Sprintf("/api/v1/collections/%d/apis", int(cid)), apiBody, true)
					apiData, _ := parseResponse(apiResp)

					if apiObj, ok := apiData["data"].(map[string]interface{}); ok {
						if aid, ok := apiObj["id"].(float64); ok {
							getResp, err := makeRequest("GET", fmt.Sprintf("/api/v1/apis/%d", int(aid)), nil, true)
							if err != nil {
								t.Fatalf("request failed: %v", err)
							}

							if getResp.StatusCode != 200 {
								t.Errorf("status code = %v, want 200", getResp.StatusCode)
							}
						}
					}
				}
			}
		}
	}
}
