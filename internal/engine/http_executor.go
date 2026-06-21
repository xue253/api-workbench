package engine

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type DebugRequest struct {
	Method      string            `json:"method"`
	URL         string            `json:"url"`
	Headers     map[string]string `json:"headers"`
	QueryParams map[string]string `json:"query_params"`
	BodyType    string            `json:"body_type"`
	Body        string            `json:"body"`
	TimeoutMs   int               `json:"timeout_ms"`
}

type DebugResponse struct {
	StatusCode    int               `json:"status_code"`
	Headers       map[string]string `json:"headers"`
	Body          string            `json:"body"`
	DurationMs    int64             `json:"duration_ms"`
	ContentLength int64             `json:"content_length"`
	Error         string            `json:"error,omitempty"`
}

func ExecuteHTTP(req *DebugRequest) *DebugResponse {
	resp := &DebugResponse{}

	if req.TimeoutMs <= 0 {
		req.TimeoutMs = 30000
	}

	client := &http.Client{
		Timeout: time.Duration(req.TimeoutMs) * time.Millisecond,
	}

	method := strings.ToUpper(req.Method)
	if method == "" {
		method = "GET"
	}

	var bodyReader io.Reader
	if req.Body != "" && method != "GET" && method != "HEAD" {
		bodyReader = strings.NewReader(req.Body)
	}

	httpReq, err := http.NewRequest(method, req.URL, bodyReader)
	if err != nil {
		resp.Error = fmt.Sprintf("请求构建失败: %v", err)
		return resp
	}

	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}

	if req.BodyType == "json" && httpReq.Header.Get("Content-Type") == "" {
		httpReq.Header.Set("Content-Type", "application/json")
	}

	q := httpReq.URL.Query()
	for k, v := range req.QueryParams {
		q.Set(k, v)
	}
	httpReq.URL.RawQuery = q.Encode()

	start := time.Now()
	httpResp, err := client.Do(httpReq)
	resp.DurationMs = time.Since(start).Milliseconds()

	if err != nil {
		resp.Error = fmt.Sprintf("请求失败: %v", err)
		return resp
	}
	defer httpResp.Body.Close()

	resp.StatusCode = httpResp.StatusCode
	resp.ContentLength = httpResp.ContentLength

	resp.Headers = make(map[string]string)
	for k, v := range httpResp.Header {
		if len(v) > 0 {
			resp.Headers[k] = v[0]
		}
	}

	bodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		resp.Error = fmt.Sprintf("读取响应失败: %v", err)
		return resp
	}
	resp.Body = string(bodyBytes)

	return resp
}

// FormatJSON 尝试格式化 JSON 字符串
func FormatJSON(s string) string {
	var v interface{}
	if err := json.Unmarshal([]byte(s), &v); err != nil {
		return s
	}
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return s
	}
	return string(b)
}
