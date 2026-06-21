package handler

import (
	"api-workbench/internal/engine"
	"api-workbench/internal/model"
	"api-workbench/internal/repository"

	"github.com/gin-gonic/gin"
)

type RunRequest struct {
	EnvID uint `json:"env_id"`
}

// RunTestCase 执行单个测试用例
func RunTestCase(c *gin.Context) {
	id, ok := getUintParam(c, "id")
	if !ok {
		return
	}

	var req RunRequest
	_ = c.ShouldBindJSON(&req)

	run, err := engine.RunTestCase(id, &engine.RunOptions{EnvID: req.EnvID})
	if err != nil {
		errorResp(c, 500, err.Error())
		return
	}
	success(c, run)
}

// RunTestSuite 执行测试套件
func RunTestSuite(c *gin.Context) {
	id, ok := getUintParam(c, "id")
	if !ok {
		return
	}

	var req RunRequest
	_ = c.ShouldBindJSON(&req)

	run, err := engine.RunTestSuite(id, &engine.RunOptions{EnvID: req.EnvID})
	if err != nil {
		errorResp(c, 500, err.Error())
		return
	}
	success(c, run)
}

// GetTestRun 获取测试运行详情
func GetTestRun(c *gin.Context) {
	id, ok := getUintParam(c, "id")
	if !ok {
		return
	}

	run := &model.TestRun{}
	err := repository.GetTestRunByID(id, run)
	if err != nil {
		errorResp(c, 404, "测试运行不存在")
		return
	}
	success(c, run)
}

// GetTestRunReport 获取测试报告
func GetTestRunReport(c *gin.Context) {
	id, ok := getUintParam(c, "id")
	if !ok {
		return
	}

	run := &model.TestRun{}
	err := repository.GetTestRunByID(id, run)
	if err != nil {
		errorResp(c, 404, "测试运行不存在")
		return
	}

	var details []model.TestRunDetail
	repository.GetTestRunDetails(id, &details)

	report := map[string]interface{}{
		"run":     run,
		"details": details,
	}
	success(c, report)
}

// GetTestRuns 获取测试运行列表
func GetTestRuns(c *gin.Context) {
	targetType := c.Query("target_type")
	targetID := c.Query("target_id")

	var runs []model.TestRun
	err := repository.GetTestRunsByFilter(targetType, targetID, &runs)
	if err != nil {
		errorResp(c, 500, err.Error())
		return
	}

	success(c, runs)
}
