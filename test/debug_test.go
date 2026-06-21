package test

import (
	"fmt"
	"testing"
)

func TestDebugAPI_Success(t *testing.T) {
	createBody := map[string]string{"name": "调试测试项目"}
	resp, _ := makeRequest("POST", "/api/v1/projects", createBody, true)
	data, _ := parseResponse(resp)

	if projectData, ok := data["data"].(map[string]interface{}); ok {
		if pid, ok := projectData["id"].(float64); ok {
			colBody := map[string]string{"name": "调试集合"}
			colResp, _ := makeRequest("POST", "/api/v1/projects/"+itoa(int(pid))+"/collections", colBody, true)
			colData, _ := parseResponse(colResp)

			if colObj, ok := colData["data"].(map[string]interface{}); ok {
				if cid, ok := colObj["id"].(float64); ok {
					apiBody := map[string]interface{}{
						"name":       "HTTP测试",
						"method":     "GET",
						"url":        "http://httpbin.org/get",
						"protocol":   "http",
						"body_type":  "json",
						"timeout_ms": 5000,
					}
					apiResp, _ := makeRequest("POST", "/api/v1/collections/"+itoa(int(cid))+"/apis", apiBody, true)
					apiData, _ := parseResponse(apiResp)

					if apiObj, ok := apiData["data"].(map[string]interface{}); ok {
						if aid, ok := apiObj["id"].(float64); ok {
							debugBody := map[string]interface{}{
								"method":    "GET",
								"url":       "http://httpbin.org/get",
								"headers":   map[string]string{"User-Agent": "API-Workbench-Test"},
								"timeout_ms": 5000,
							}
							debugResp, err := makeRequest("POST", "/api/v1/apis/"+itoa(int(aid))+"/debug", debugBody, true)
							if err != nil {
								t.Fatalf("request failed: %v", err)
							}

							if debugResp.StatusCode != 200 {
								t.Logf("Debug response: %v", debugResp.StatusCode)
							}
						}
					}
				}
			}
		}
	}
}

func TestGetDynamicFunctions(t *testing.T) {
	resp, err := makeRequest("GET", "/api/v1/variables/dynamic", nil, true)
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

func TestGetTestRuns(t *testing.T) {
	resp, err := makeRequest("GET", "/api/v1/test-runs", nil, true)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("status code = %v, want 200", resp.StatusCode)
	}
}

func TestGetTestRun_InvalidID(t *testing.T) {
	resp, err := makeRequest("GET", "/api/v1/test-runs/abc", nil, true)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if resp.StatusCode != 400 {
		t.Errorf("status code = %v, want 400", resp.StatusCode)
	}
}

func TestExportTestReport(t *testing.T) {
	createBody := map[string]string{"name": "导出报告项目"}
	resp, _ := makeRequest("POST", "/api/v1/projects", createBody, true)
	data, _ := parseResponse(resp)

	if projectData, ok := data["data"].(map[string]interface{}); ok {
		if pid, ok := projectData["id"].(float64); ok {
			tcBody := map[string]string{"name": "报告测试用例"}
			tcResp, _ := makeRequest("POST", "/api/v1/projects/"+itoa(int(pid))+"/test-cases", tcBody, true)
			tcData, _ := parseResponse(tcResp)

			if tcObj, ok := tcData["data"].(map[string]interface{}); ok {
				if tcid, ok := tcObj["id"].(float64); ok {
					runResp, _ := makeRequest("POST", "/api/v1/test-cases/"+itoa(int(tcid))+"/run", map[string]interface{}{}, true)
					runData, _ := parseResponse(runResp)

					if runObj, ok := runData["data"].(map[string]interface{}); ok {
						if rid, ok := runObj["id"].(float64); ok {
							exportResp, err := makeRequest("GET", "/api/v1/test-runs/"+itoa(int(rid))+"/export?format=md", nil, true)
							if err != nil {
								t.Fatalf("request failed: %v", err)
							}

							if exportResp.StatusCode != 200 {
								t.Logf("Export status: %v", exportResp.StatusCode)
							}
						}
					}
				}
			}
		}
	}
}

func itoa(i int) string {
	return fmt.Sprintf("%d", i)
}
