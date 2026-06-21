package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
)

type GRPCRequest struct {
	Address   string            `json:"address"`
	Service   string            `json:"service"`
	Method    string            `json:"method"`
	Headers   map[string]string `json:"headers"`
	Body      string            `json:"body"`
	TimeoutMs int               `json:"timeout_ms"`
}

type GRPCResponse struct {
	StatusCode  int               `json:"status_code"`
	Headers     map[string]string `json:"headers"`
	Body        string            `json:"body"`
	DurationMs  int64             `json:"duration_ms"`
	Error       string            `json:"error,omitempty"`
}

func ExecuteGRPC(req *GRPCRequest) *GRPCResponse {
	resp := &GRPCResponse{}

	if req.TimeoutMs <= 0 {
		req.TimeoutMs = 30000
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(req.TimeoutMs)*time.Millisecond)
	defer cancel()

	conn, err := grpc.NewClient(req.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		resp.Error = fmt.Sprintf("连接失败: %v", err)
		return resp
	}
	defer conn.Close()

	fullMethod := fmt.Sprintf("/%s/%s", req.Service, req.Method)

	requestMsg := &structpb.Struct{}
	if req.Body != "" {
		if err := protojson.Unmarshal([]byte(req.Body), requestMsg); err != nil {
			resp.Error = fmt.Sprintf("请求体解析失败: %v", err)
			return resp
		}
	}

	start := time.Now()

	respBody := &structpb.Struct{}
	err = conn.Invoke(ctx, fullMethod, requestMsg, respBody)
	resp.DurationMs = time.Since(start).Milliseconds()

	if err != nil {
		resp.Error = fmt.Sprintf("调用失败: %v", err)
		return resp
	}

	bodyBytes, _ := protojson.MarshalOptions{Indent: "  "}.Marshal(respBody)
	resp.Body = string(bodyBytes)

	resp.StatusCode = 0
	resp.Headers = make(map[string]string)

	return resp
}

func FormatGRPCResponse(body string) string {
	var v interface{}
	if err := json.Unmarshal([]byte(body), &v); err != nil {
		return body
	}
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return body
	}
	return string(b)
}
