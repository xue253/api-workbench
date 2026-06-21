package engine

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
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

	if err := validateURL(req.URL); err != nil {
		resp.Error = err.Error()
		return resp
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

	bodyBytes, err := io.ReadAll(io.LimitReader(httpResp.Body, 10<<20))
	if err != nil {
		resp.Error = fmt.Sprintf("读取响应失败: %v", err)
		return resp
	}
	resp.Body = string(bodyBytes)

	return resp
}

func validateURL(rawURL string) error {
	u, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("无效的 URL: %v", err)
	}
	host := u.Hostname()
	if host == "" {
		return fmt.Errorf("URL 缺少主机名")
	}
	addrs, err := net.LookupHost(host)
	if err != nil {
		return fmt.Errorf("无法解析主机: %v", err)
	}
	for _, addr := range addrs {
		ip := net.ParseIP(addr)
		if ip == nil {
			continue
		}
		if isPrivateIP(ip) {
			return fmt.Errorf("禁止访问内网地址: %s", addr)
		}
	}
	return nil
}

func isPrivateIP(ip net.IP) bool {
	privateRanges := []struct {
		start net.IP
		end   net.IP
	}{
		{net.ParseIP("10.0.0.0"), net.ParseIP("10.255.255.255")},
		{net.ParseIP("172.16.0.0"), net.ParseIP("172.31.255.255")},
		{net.ParseIP("192.168.0.0"), net.ParseIP("192.168.255.255")},
		{net.ParseIP("127.0.0.0"), net.ParseIP("127.255.255.255")},
		{net.ParseIP("169.254.0.0"), net.ParseIP("169.254.255.255")},
		{net.ParseIP("::1"), net.ParseIP("::1")},
	}
	for _, r := range privateRanges {
		if bytesInRange(ip, r.start, r.end) {
			return true
		}
	}
	return ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast()
}

func bytesInRange(ip, start, end net.IP) bool {
	ip4 := ip.To4()
	start4 := start.To4()
	end4 := end.To4()
	if ip4 == nil || start4 == nil || end4 == nil {
		return false
	}
	for i := 0; i < 4; i++ {
		if ip4[i] < start4[i] {
			return false
		}
		if ip4[i] > end4[i] {
			return false
		}
		if ip4[i] > start4[i] && ip4[i] < end4[i] {
			return true
		}
	}
	return true
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
