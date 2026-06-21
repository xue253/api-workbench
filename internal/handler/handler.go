package handler

import (
	"net/http"
	"strconv"

	"api-workbench/internal/model"
	"api-workbench/internal/repository"
	"api-workbench/internal/scheduler"

	"github.com/gin-gonic/gin"
)

func success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func errorResp(c *gin.Context, code int, msg string) {
	c.JSON(code, gin.H{"error": msg})
}

// ---- Environment ----
func ListEnvironments(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		return
	}
	pid, err := strconv.Atoi(c.Param("pid"))
	if err != nil || pid <= 0 {
		errorResp(c, 400, "无效的项目ID")
		return
	}
	if !repository.IsProjectOwnedByUser(uint(pid), userID) {
		errorResp(c, 403, "无权访问此项目")
		return
	}
	var list []model.Environment
	if err := repository.GetEnvironmentsByProject(uint(pid), &list); err != nil {
		errorResp(c, 500, err.Error())
		return
	}
	success(c, list)
}

func CreateEnvironment(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		return
	}
	pid, err := strconv.Atoi(c.Param("pid"))
	if err != nil || pid <= 0 {
		errorResp(c, 400, "无效的项目ID")
		return
	}
	if !repository.IsProjectOwnedByUser(uint(pid), userID) {
		errorResp(c, 403, "无权操作此项目")
		return
	}
	var e model.Environment
	if err := c.ShouldBindJSON(&e); err != nil {
		errorResp(c, 400, err.Error())
		return
	}
	e.ProjectID = uint(pid)
	if err := repository.CreateEnvironment(&e); err != nil {
		errorResp(c, 500, err.Error())
		return
	}
	success(c, e)
}

func UpdateEnvironment(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	if !repository.IsEnvironmentOwnedByUser(uint(id), userID) {
		errorResp(c, 403, "无权操作此环境")
		return
	}
	var e model.Environment
	if err := c.ShouldBindJSON(&e); err != nil {
		errorResp(c, 400, err.Error())
		return
	}
	e.ID = uint(id)
	if err := repository.UpdateEnvironment(&e); err != nil {
		errorResp(c, 500, err.Error())
		return
	}
	success(c, e)
}

func DeleteEnvironment(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	if !repository.IsEnvironmentOwnedByUser(uint(id), userID) {
		errorResp(c, 403, "无权操作此环境")
		return
	}
	if err := repository.DeleteEnvironment(uint(id)); err != nil {
		errorResp(c, 500, err.Error())
		return
	}
	success(c, nil)
}

// ---- Environment Variables ----
func ListEnvVars(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	if !repository.IsEnvironmentOwnedByUser(uint(id), userID) {
		errorResp(c, 403, "无权操作此环境")
		return
	}
	var list []model.EnvironmentVariable
	if err := repository.GetEnvVarsByEnvID(uint(id), &list); err != nil {
		errorResp(c, 500, err.Error())
		return
	}
	success(c, list)
}

func SaveEnvVars(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	if !repository.IsEnvironmentOwnedByUser(uint(id), userID) {
		errorResp(c, 403, "无权操作此环境")
		return
	}
	var vars []model.EnvironmentVariable
	if err := c.ShouldBindJSON(&vars); err != nil {
		errorResp(c, 400, err.Error())
		return
	}
	if err := repository.SaveEnvVars(uint(id), vars); err != nil {
		errorResp(c, 500, err.Error())
		return
	}
	success(c, vars)
}

// ---- Collection ----
func ListCollections(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		return
	}
	pid, err := strconv.Atoi(c.Param("pid"))
	if err != nil || pid <= 0 {
		errorResp(c, 400, "无效的项目ID")
		return
	}
	if !repository.IsProjectOwnedByUser(uint(pid), userID) {
		errorResp(c, 403, "无权访问此项目")
		return
	}
	var list []model.Collection
	if err := repository.GetCollectionsByProject(uint(pid), &list); err != nil {
		errorResp(c, 500, err.Error())
		return
	}
	success(c, list)
}

func CreateCollection(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		return
	}
	pid, err := strconv.Atoi(c.Param("pid"))
	if err != nil || pid <= 0 {
		errorResp(c, 400, "无效的项目ID")
		return
	}
	if !repository.IsProjectOwnedByUser(uint(pid), userID) {
		errorResp(c, 403, "无权操作此项目")
		return
	}
	var col model.Collection
	if err := c.ShouldBindJSON(&col); err != nil {
		errorResp(c, 400, err.Error())
		return
	}
	col.ProjectID = uint(pid)
	if err := repository.CreateCollection(&col); err != nil {
		errorResp(c, 500, err.Error())
		return
	}
	success(c, col)
}

func UpdateCollection(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	if !repository.IsCollectionOwnedByUser(uint(id), userID) {
		errorResp(c, 403, "无权操作此接口分组")
		return
	}
	var col model.Collection
	if err := c.ShouldBindJSON(&col); err != nil {
		errorResp(c, 400, err.Error())
		return
	}
	col.ID = uint(id)
	if err := repository.UpdateCollection(&col); err != nil {
		errorResp(c, 500, err.Error())
		return
	}
	success(c, col)
}

func DeleteCollection(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	if !repository.IsCollectionOwnedByUser(uint(id), userID) {
		errorResp(c, 403, "无权操作此接口分组")
		return
	}
	if err := repository.DeleteCollection(uint(id)); err != nil {
		errorResp(c, 500, err.Error())
		return
	}
	success(c, nil)
}

func MoveCollection(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	if !repository.IsCollectionOwnedByUser(uint(id), userID) {
		errorResp(c, 403, "无权操作此接口分组")
		return
	}
	var body struct {
		ParentID *uint `json:"parent_id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		errorResp(c, 400, err.Error())
		return
	}
	if err := repository.MoveCollection(uint(id), body.ParentID); err != nil {
		errorResp(c, 500, err.Error())
		return
	}
	success(c, nil)
}

// ---- API ----
func ListAPIsByCollection(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		errorResp(c, 400, "无效的分组ID")
		return
	}
	if !repository.IsCollectionOwnedByUser(uint(id), userID) {
		errorResp(c, 403, "无权访问此分组")
		return
	}
	var list []model.API
	if err := repository.GetAPIsByCollection(uint(id), &list); err != nil {
		errorResp(c, 500, err.Error())
		return
	}
	success(c, list)
}

func CreateAPIByCollection(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		errorResp(c, 400, "无效的分组ID")
		return
	}
	if !repository.IsCollectionOwnedByUser(uint(id), userID) {
		errorResp(c, 403, "无权操作此分组")
		return
	}
	var a model.API
	if err := c.ShouldBindJSON(&a); err != nil {
		errorResp(c, 400, err.Error())
		return
	}
	a.CollectionID = uint(id)
	if err := repository.CreateAPI(&a); err != nil {
		errorResp(c, 500, err.Error())
		return
	}
	success(c, a)
}

func GetAPI(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	if !repository.IsAPIOwnedByUser(uint(id), userID) {
		errorResp(c, 403, "无权访问此接口")
		return
	}
	var a model.API
	if err := repository.GetAPIByID(uint(id), &a); err != nil {
		errorResp(c, 404, "not found")
		return
	}
	success(c, a)
}

func UpdateAPI(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	if !repository.IsAPIOwnedByUser(uint(id), userID) {
		errorResp(c, 403, "无权操作此接口")
		return
	}
	var a model.API
	if err := c.ShouldBindJSON(&a); err != nil {
		errorResp(c, 400, err.Error())
		return
	}
	a.ID = uint(id)
	if err := repository.UpdateAPI(&a); err != nil {
		errorResp(c, 500, err.Error())
		return
	}
	success(c, a)
}

func DeleteAPI(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	if !repository.IsAPIOwnedByUser(uint(id), userID) {
		errorResp(c, 403, "无权操作此接口")
		return
	}
	if err := repository.DeleteAPI(uint(id)); err != nil {
		errorResp(c, 500, err.Error())
		return
	}
	success(c, nil)
}

// ---- Assertions ----
func ListAssertions(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		return
	}
	aid, _ := strconv.Atoi(c.Param("id"))
	if !repository.IsAPIOwnedByUser(uint(aid), userID) {
		errorResp(c, 403, "无权访问此接口")
		return
	}
	var list []model.Assertion
	if err := repository.GetAssertionsByAPIID(uint(aid), &list); err != nil {
		errorResp(c, 500, err.Error())
		return
	}
	success(c, list)
}

func SaveAssertions(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		return
	}
	aid, _ := strconv.Atoi(c.Param("id"))
	if !repository.IsAPIOwnedByUser(uint(aid), userID) {
		errorResp(c, 403, "无权操作此接口")
		return
	}
	var assertions []model.Assertion
	if err := c.ShouldBindJSON(&assertions); err != nil {
		errorResp(c, 400, err.Error())
		return
	}
	if err := repository.SaveAssertions(uint(aid), assertions); err != nil {
		errorResp(c, 500, err.Error())
		return
	}
	success(c, assertions)
}

// ---- TestCase ----
func ListTestCases(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		return
	}
	pid, err := strconv.Atoi(c.Param("pid"))
	if err != nil || pid <= 0 {
		errorResp(c, 400, "无效的项目ID")
		return
	}
	if !repository.IsProjectOwnedByUser(uint(pid), userID) {
		errorResp(c, 403, "无权访问此项目")
		return
	}
	var list []model.TestCase
	if err := repository.GetTestCasesByProject(uint(pid), &list); err != nil {
		errorResp(c, 500, err.Error())
		return
	}
	success(c, list)
}

func CreateTestCase(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		return
	}
	pid, err := strconv.Atoi(c.Param("pid"))
	if err != nil || pid <= 0 {
		errorResp(c, 400, "无效的项目ID")
		return
	}
	if !repository.IsProjectOwnedByUser(uint(pid), userID) {
		errorResp(c, 403, "无权操作此项目")
		return
	}
	var tc model.TestCase
	if err := c.ShouldBindJSON(&tc); err != nil {
		errorResp(c, 400, err.Error())
		return
	}
	tc.ProjectID = uint(pid)
	if err := repository.CreateTestCase(&tc); err != nil {
		errorResp(c, 500, err.Error())
		return
	}
	success(c, tc)
}

func UpdateTestCase(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	if !repository.IsTestCaseOwnedByUser(uint(id), userID) {
		errorResp(c, 403, "无权操作此测试用例")
		return
	}
	var tc model.TestCase
	if err := c.ShouldBindJSON(&tc); err != nil {
		errorResp(c, 400, err.Error())
		return
	}
	tc.ID = uint(id)
	if err := repository.UpdateTestCase(&tc); err != nil {
		errorResp(c, 500, err.Error())
		return
	}
	success(c, tc)
}

func DeleteTestCase(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	if !repository.IsTestCaseOwnedByUser(uint(id), userID) {
		errorResp(c, 403, "无权操作此测试用例")
		return
	}
	if err := repository.DeleteTestCase(uint(id)); err != nil {
		errorResp(c, 500, err.Error())
		return
	}
	success(c, nil)
}

func SaveTestCaseAPIs(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	if !repository.IsTestCaseOwnedByUser(uint(id), userID) {
		errorResp(c, 403, "无权操作此测试用例")
		return
	}
	var apis []model.TestCaseAPI
	if err := c.ShouldBindJSON(&apis); err != nil {
		errorResp(c, 400, err.Error())
		return
	}
	if err := repository.SaveTestCaseAPIs(uint(id), apis); err != nil {
		errorResp(c, 500, err.Error())
		return
	}
	success(c, apis)
}

// ---- TestDataSet ----
func ListDataSets(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		return
	}
	tcaID, _ := strconv.Atoi(c.Param("id"))
	if !repository.IsTestCaseAPIOwnedByUser(uint(tcaID), userID) {
		errorResp(c, 403, "无权访问此测试接口")
		return
	}
	var list []model.TestDataSet
	if err := repository.GetTestDataSets(uint(tcaID), &list); err != nil {
		errorResp(c, 500, err.Error())
		return
	}
	success(c, list)
}

func SaveDataSets(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		return
	}
	tcaID, _ := strconv.Atoi(c.Param("id"))
	if !repository.IsTestCaseAPIOwnedByUser(uint(tcaID), userID) {
		errorResp(c, 403, "无权操作此测试接口")
		return
	}
	var datasets []model.TestDataSet
	if err := c.ShouldBindJSON(&datasets); err != nil {
		errorResp(c, 400, err.Error())
		return
	}
	if err := repository.SaveTestDataSets(uint(tcaID), datasets); err != nil {
		errorResp(c, 500, err.Error())
		return
	}
	success(c, datasets)
}

// ---- TestSuite ----
func ListTestSuites(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		return
	}
	pid, err := strconv.Atoi(c.Param("pid"))
	if err != nil || pid <= 0 {
		errorResp(c, 400, "无效的项目ID")
		return
	}
	if !repository.IsProjectOwnedByUser(uint(pid), userID) {
		errorResp(c, 403, "无权访问此项目")
		return
	}
	var list []model.TestSuite
	if err := repository.GetTestSuitesByProject(uint(pid), &list); err != nil {
		errorResp(c, 500, err.Error())
		return
	}
	success(c, list)
}

func CreateTestSuite(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		return
	}
	pid, err := strconv.Atoi(c.Param("pid"))
	if err != nil || pid <= 0 {
		errorResp(c, 400, "无效的项目ID")
		return
	}
	if !repository.IsProjectOwnedByUser(uint(pid), userID) {
		errorResp(c, 403, "无权操作此项目")
		return
	}
	var ts model.TestSuite
	if err := c.ShouldBindJSON(&ts); err != nil {
		errorResp(c, 400, err.Error())
		return
	}
	ts.ProjectID = uint(pid)
	if err := repository.CreateTestSuite(&ts); err != nil {
		errorResp(c, 500, err.Error())
		return
	}
	success(c, ts)
}

func UpdateTestSuite(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	if !repository.IsTestSuiteOwnedByUser(uint(id), userID) {
		errorResp(c, 403, "无权操作此测试套件")
		return
	}
	var ts model.TestSuite
	if err := c.ShouldBindJSON(&ts); err != nil {
		errorResp(c, 400, err.Error())
		return
	}
	ts.ID = uint(id)
	if err := repository.UpdateTestSuite(&ts); err != nil {
		errorResp(c, 500, err.Error())
		return
	}
	success(c, ts)
}

func DeleteTestSuite(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	if !repository.IsTestSuiteOwnedByUser(uint(id), userID) {
		errorResp(c, 403, "无权操作此测试套件")
		return
	}
	if err := repository.DeleteTestSuite(uint(id)); err != nil {
		errorResp(c, 500, err.Error())
		return
	}
	success(c, nil)
}

func SaveTestSuiteCases(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	if !repository.IsTestSuiteOwnedByUser(uint(id), userID) {
		errorResp(c, 403, "无权操作此测试套件")
		return
	}
	var cases []model.TestSuiteCase
	if err := c.ShouldBindJSON(&cases); err != nil {
		errorResp(c, 400, err.Error())
		return
	}
	if err := repository.SaveTestSuiteCases(uint(id), cases); err != nil {
		errorResp(c, 500, err.Error())
		return
	}
	success(c, cases)
}

// ---- ScheduledTask ----
func ListScheduledTasks(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		return
	}
	pid, err := strconv.Atoi(c.Param("pid"))
	if err != nil || pid <= 0 {
		errorResp(c, 400, "无效的项目ID")
		return
	}
	if !repository.IsProjectOwnedByUser(uint(pid), userID) {
		errorResp(c, 403, "无权访问此项目")
		return
	}
	var list []model.ScheduledTask
	if err := repository.GetScheduledTasksByProject(uint(pid), &list); err != nil {
		errorResp(c, 500, err.Error())
		return
	}
	success(c, list)
}

func CreateScheduledTask(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		return
	}
	pid, err := strconv.Atoi(c.Param("pid"))
	if err != nil || pid <= 0 {
		errorResp(c, 400, "无效的项目ID")
		return
	}
	if !repository.IsProjectOwnedByUser(uint(pid), userID) {
		errorResp(c, 403, "无权操作此项目")
		return
	}
	var st model.ScheduledTask
	if err := c.ShouldBindJSON(&st); err != nil {
		errorResp(c, 400, err.Error())
		return
	}
	st.ProjectID = uint(pid)
	if err := repository.CreateScheduledTask(&st); err != nil {
		errorResp(c, 500, err.Error())
		return
	}
	scheduler.AddTask(st)
	success(c, st)
}

func UpdateScheduledTask(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	if !repository.IsScheduledTaskOwnedByUser(uint(id), userID) {
		errorResp(c, 403, "无权操作此调度任务")
		return
	}
	var st model.ScheduledTask
	if err := c.ShouldBindJSON(&st); err != nil {
		errorResp(c, 400, err.Error())
		return
	}
	st.ID = uint(id)
	if err := repository.UpdateScheduledTask(&st); err != nil {
		errorResp(c, 500, err.Error())
		return
	}
	scheduler.UpdateTask(st)
	success(c, st)
}

func DeleteScheduledTask(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	if !repository.IsScheduledTaskOwnedByUser(uint(id), userID) {
		errorResp(c, 403, "无权操作此调度任务")
		return
	}
	scheduler.RemoveTask(uint(id))
	if err := repository.DeleteScheduledTask(uint(id)); err != nil {
		errorResp(c, 500, err.Error())
		return
	}
	success(c, nil)
}

// ---- Health ----
func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
