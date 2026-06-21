package handler

import (
	"encoding/json"
	"net/http"

	"api-workbench/internal/engine"
	"api-workbench/internal/model"
	"api-workbench/internal/repository"
	"api-workbench/internal/variable"

	"github.com/gin-gonic/gin"
)

type DebugRequestPayload struct {
	Method      string            `json:"method"`
	URL         string            `json:"url"`
	Headers     map[string]string `json:"headers"`
	QueryParams map[string]string `json:"query_params"`
	BodyType    string            `json:"body_type"`
	Body        string            `json:"body"`
	TimeoutMs   int               `json:"timeout_ms"`
	EnvID       uint              `json:"env_id"`
}

// DebugAPI 发送调试请求
func DebugAPI(c *gin.Context) {
	id, ok := getUintParam(c, "id")
	if !ok {
		return
	}

	var payload DebugRequestPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		errorResp(c, 400, err.Error())
		return
	}

	// 获取 API 定义
	var apiDef model.API
	if err := repository.GetAPIByID(id, &apiDef); err != nil {
		errorResp(c, 404, "接口不存在")
		return
	}

	// 合并请求参数（payload 覆盖 apiDef）
	method := payload.Method
	if method == "" {
		method = apiDef.Method
	}
	url := payload.URL
	if url == "" {
		url = apiDef.URL
	}
	bodyType := payload.BodyType
	if bodyType == "" {
		bodyType = apiDef.BodyType
	}
	body := payload.Body
	if body == "" {
		body = apiDef.Body
	}
	timeoutMs := payload.TimeoutMs
	if timeoutMs <= 0 {
		timeoutMs = apiDef.TimeoutMs
	}

	// 加载环境变量
	envVars := make(map[string]string)
	if payload.EnvID > 0 {
		var vars []model.EnvironmentVariable
		if err := repository.GetEnvVarsByEnvID(payload.EnvID, &vars); err == nil {
			for _, v := range vars {
				envVars[v.Key] = v.Value
			}
		}
	}

	// 根据协议类型执行
	if apiDef.Protocol == "grpc" {
		body = variable.Replace(body, envVars, nil)

		req := &engine.GRPCRequest{
			Address:   url,
			Service:   apiDef.ProtoService,
			Method:    method,
			Body:      body,
			TimeoutMs: timeoutMs,
		}

		resp := engine.ExecuteGRPC(req)
		resp.Body = engine.FormatGRPCResponse(resp.Body)
		success(c, resp)
		return
	}

	// HTTP 协议
	// 解析 headers
	headers := make(map[string]string)
	if apiDef.Headers != "" {
		_ = json.Unmarshal([]byte(apiDef.Headers), &headers)
	}
	for k, v := range payload.Headers {
		headers[k] = v
	}

	// 解析 query params
	queryParams := make(map[string]string)
	if apiDef.QueryParams != "" {
		_ = json.Unmarshal([]byte(apiDef.QueryParams), &queryParams)
	}
	for k, v := range payload.QueryParams {
		queryParams[k] = v
	}

	// 替换变量
	url = variable.Replace(url, envVars, nil)
	body = variable.Replace(body, envVars, nil)
	headers = variable.ReplaceMap(headers, envVars, nil)
	queryParams = variable.ReplaceMap(queryParams, envVars, nil)

	// 执行请求
	req := &engine.DebugRequest{
		Method:      method,
		URL:         url,
		Headers:     headers,
		QueryParams: queryParams,
		BodyType:    bodyType,
		Body:        body,
		TimeoutMs:   timeoutMs,
	}

	resp := engine.ExecuteHTTP(req)
	resp.Body = engine.FormatJSON(resp.Body)

	success(c, resp)
}

// GetDebugHistory 获取调试历史（最近20条）
func GetDebugHistory(c *gin.Context) {
	// 暂用内存存储，后续可改为数据库
	c.JSON(http.StatusOK, gin.H{"data": []interface{}{}})
}

// GetDynamicFunctions 获取可用的动态变量函数列表
func GetDynamicFunctions(c *gin.Context) {
	functions := []map[string]string{
		{"name": "timestamp", "desc": "当前时间戳（毫秒）", "example": "{{timestamp}}"},
		{"name": "unix", "desc": "当前 Unix 时间戳（秒）", "example": "{{unix}}"},
		{"name": "uuid", "desc": "生成 UUID", "example": "{{uuid}}"},
		{"name": "random_string", "desc": "随机字符串", "example": "{{random_string}}"},
		{"name": "random_int", "desc": "随机整数", "example": "{{random_int}}"},
		{"name": "date", "desc": "当前日期 YYYY-MM-DD", "example": "{{date}}"},
		{"name": "datetime", "desc": "当前日期时间", "example": "{{datetime}}"},
		{"name": "year", "desc": "当前年份", "example": "{{year}}"},
		{"name": "month", "desc": "当前月份", "example": "{{month}}"},
		{"name": "day", "desc": "当前日期", "example": "{{day}}"},
	}
	c.JSON(http.StatusOK, gin.H{"data": functions})
}
