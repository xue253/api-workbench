package handler

import (
	"fmt"
	"net/http"
	"strings"

	"api-workbench/internal/model"
	"api-workbench/internal/repository"

	"github.com/gin-gonic/gin"
)

// ExportTestReport 导出测试报告
func ExportTestReport(c *gin.Context) {
	id, ok := getUintParam(c, "id")
	if !ok {
		return
	}

	format := c.DefaultQuery("format", "md")

	run := &model.TestRun{}
	if err := repository.GetTestRunByID(id, run); err != nil {
		errorResp(c, 404, "测试运行不存在")
		return
	}

	var details []model.TestRunDetail
	repository.GetTestRunDetails(id, &details)

	var content string
	if format == "html" {
		content = generateHTMLReport(run, details)
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"report-%d.html\"", id))
	} else {
		content = generateMDReport(run, details)
		c.Header("Content-Type", "text/markdown; charset=utf-8")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"report-%d.md\"", id))
	}

	c.String(http.StatusOK, content)
}

func generateMDReport(run *model.TestRun, details []model.TestRunDetail) string {
	var sb strings.Builder

	sb.WriteString("# 测试报告\n\n")
	sb.WriteString(fmt.Sprintf("**执行时间**: %s\n", run.StartedAt.Format("2006-01-02 15:04:05")))
	sb.WriteString(fmt.Sprintf("**执行类型**: %s\n", run.TargetType))
	sb.WriteString(fmt.Sprintf("**状态**: %s\n", run.Status))
	sb.WriteString(fmt.Sprintf("**总耗时**: %dms\n\n", run.DurationMs))

	sb.WriteString("## 汇总\n\n")
	sb.WriteString(fmt.Sprintf("- 总数: %d\n", run.Total))
	sb.WriteString(fmt.Sprintf("- 通过: %d\n", run.Passed))
	sb.WriteString(fmt.Sprintf("- 失败: %d\n", run.Failed))
	sb.WriteString(fmt.Sprintf("- 跳过: %d\n", run.Skipped))

	if run.Total > 0 {
		passRate := float64(run.Passed) / float64(run.Total) * 100
		sb.WriteString(fmt.Sprintf("- 通过率: %.1f%%\n", passRate))
	}

	sb.WriteString("\n## 详情\n\n")
	sb.WriteString("| # | 状态 | 状态码 | 耗时 | 错误 |\n")
	sb.WriteString("|---|------|--------|------|------|\n")

	for i, d := range details {
		status := "✅"
		if d.Status == "failed" {
			status = "❌"
		}
		errMsg := d.ErrorMessage
		if errMsg == "" {
			errMsg = "-"
		}
		if len(errMsg) > 50 {
			errMsg = errMsg[:50] + "..."
		}
		sb.WriteString(fmt.Sprintf("| %d | %s | %d | %dms | %s |\n", i+1, status, d.StatusCode, d.DurationMs, errMsg))
	}

	return sb.String()
}

func generateHTMLReport(run *model.TestRun, details []model.TestRunDetail) string {
	passRate := 0.0
	if run.Total > 0 {
		passRate = float64(run.Passed) / float64(run.Total) * 100
	}

	statusColor := "#52c41a"
	if run.Status == "failed" {
		statusColor = "#ff4d4f"
	}

	var rows strings.Builder
	for i, d := range details {
		rowColor := ""
		if d.Status == "failed" {
			rowColor = "#fff2f0"
		}
		errMsg := d.ErrorMessage
		if errMsg == "" {
			errMsg = "-"
		}
		rows.WriteString(fmt.Sprintf(`<tr style="background:%s"><td>%d</td><td>%s</td><td>%d</td><td>%dms</td><td>%s</td></tr>`,
			rowColor, i+1, d.Status, d.StatusCode, d.DurationMs, errMsg))
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<title>测试报告</title>
<style>
body{font-family:-apple-system,BlinkMacSystemFont,sans-serif;max-width:900px;margin:40px auto;padding:0 20px;color:#1d1d1f}
h1{font-size:24px;font-weight:600}
.stats{display:flex;gap:24px;margin:20px 0}
.stat{text-align:center;padding:20px;background:#f5f5f7;border-radius:12px;flex:1}
.stat-value{font-size:32px;font-weight:600}
.stat-label{font-size:14px;color:#86868b;margin-top:4px}
table{width:100%%;border-collapse:collapse;margin-top:20px}
th,td{padding:12px;text-align:left;border-bottom:1px solid #e5e5e5}
th{background:#f5f5f7;font-weight:600}
</style>
</head>
<body>
<h1>测试报告</h1>
<p>执行时间: %s | 状态: <span style="color:%s;font-weight:600">%s</span> | 耗时: %dms</p>
<div class="stats">
<div class="stat"><div class="stat-value">%d</div><div class="stat-label">总数</div></div>
<div class="stat"><div class="stat-value" style="color:#52c41a">%d</div><div class="stat-label">通过</div></div>
<div class="stat"><div class="stat-value" style="color:#ff4d4f">%d</div><div class="stat-label">失败</div></div>
<div class="stat"><div class="stat-value">%.1f%%</div><div class="stat-label">通过率</div></div>
</div>
<table>
<tr><th>#</th><th>状态</th><th>状态码</th><th>耗时</th><th>错误</th></tr>
%s
</table>
</body>
</html>`,
		run.StartedAt.Format("2006-01-02 15:04:05"),
		statusColor,
		run.Status,
		run.DurationMs,
		run.Total,
		run.Passed,
		run.Failed,
		passRate,
		rows.String(),
	)
}
