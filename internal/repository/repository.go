package repository

import (
	"api-workbench/internal/db"
	"api-workbench/internal/model"

	"gorm.io/gorm"
)

// User
func CreateUser(u *model.User) error { return db.DB.Create(u).Error }
func GetUserByUsername(username string, user *model.User) error {
	return db.DB.Where("username = ?", username).First(user).Error
}
func GetUserByID(id uint, user *model.User) error { return db.DB.First(user, id).Error }
func UpdateUser(u *model.User) error { return db.DB.Save(u).Error }
func DeleteUser(id uint) error { return db.DB.Delete(&model.User{}, id).Error }
func GetUsers(list *[]model.User) error { return db.DB.Find(list).Error }

// Project
func CreateProject(p *model.Project) error { return db.DB.Create(p).Error }
func GetProjectsByUser(uid uint, list *[]model.Project) error {
	return db.DB.Where("user_id = ?", uid).Find(list).Error
}
func GetProjectByID(id uint, p *model.Project) error { return db.DB.First(p, id).Error }
func UpdateProject(p *model.Project) error { return db.DB.Save(p).Error }
func DeleteProject(id uint) error { return db.DB.Delete(&model.Project{}, id).Error }

// Environment
func CreateEnvironment(e *model.Environment) error { return db.DB.Create(e).Error }
func GetEnvironmentsByProject(pid uint, list *[]model.Environment) error {
	return db.DB.Where("project_id = ?", pid).Order("sort_order").Find(list).Error
}
func GetEnvironmentByID(id uint, e *model.Environment) error { return db.DB.First(e, id).Error }
func UpdateEnvironment(e *model.Environment) error { return db.DB.Save(e).Error }
func DeleteEnvironment(id uint) error { return db.DB.Delete(&model.Environment{}, id).Error }

// EnvironmentVariable
func GetEnvVarsByEnvID(eid uint, list *[]model.EnvironmentVariable) error {
	return db.DB.Where("environment_id = ?", eid).Find(list).Error
}
func SaveEnvVars(eid uint, vars []model.EnvironmentVariable) error {
	db.DB.Where("environment_id = ?", eid).Delete(&model.EnvironmentVariable{})
	for i := range vars {
		vars[i].EnvironmentID = eid
	}
	return db.DB.Create(&vars).Error
}

// Collection
func CreateCollection(c *model.Collection) error { return db.DB.Create(c).Error }
func GetCollectionsByProject(pid uint, list *[]model.Collection) error {
	return db.DB.Where("project_id = ?", pid).Order("sort_order").Find(list).Error
}
func GetCollectionByID(id uint, c *model.Collection) error { return db.DB.First(c, id).Error }
func UpdateCollection(c *model.Collection) error { return db.DB.Save(c).Error }
func DeleteCollection(id uint) error { return db.DB.Delete(&model.Collection{}, id).Error }
func MoveCollection(id uint, parentID *uint) error {
	return db.DB.Model(&model.Collection{}).Where("id = ?", id).Update("parent_id", parentID).Error
}

// API
func CreateAPI(a *model.API) error { return db.DB.Create(a).Error }
func GetAPIsByCollection(cid uint, list *[]model.API) error {
	return db.DB.Where("collection_id = ?", cid).Find(list).Error
}
func GetAPIByID(id uint, a *model.API) error { return db.DB.First(a, id).Error }
func UpdateAPI(a *model.API) error { return db.DB.Save(a).Error }
func DeleteAPI(id uint) error { return db.DB.Delete(&model.API{}, id).Error }

// Assertion
func GetAssertionsByAPIID(aid uint, list *[]model.Assertion) error {
	return db.DB.Where("api_id = ?", aid).Find(list).Error
}
func SaveAssertions(aid uint, assertions []model.Assertion) error {
	db.DB.Where("api_id = ?", aid).Delete(&model.Assertion{})
	for i := range assertions {
		assertions[i].APIID = aid
	}
	return db.DB.Create(&assertions).Error
}

// TestCase
func CreateTestCase(tc *model.TestCase) error { return db.DB.Create(tc).Error }
func GetTestCasesByProject(pid uint, list *[]model.TestCase) error {
	return db.DB.Where("project_id = ?", pid).Find(list).Error
}
func GetTestCaseByID(id uint, tc *model.TestCase) error { return db.DB.First(tc, id).Error }
func UpdateTestCase(tc *model.TestCase) error { return db.DB.Save(tc).Error }
func DeleteTestCase(id uint) error { return db.DB.Delete(&model.TestCase{}, id).Error }

// TestCaseAPI
func GetTestCaseAPIs(tcID uint, list *[]model.TestCaseAPI) error {
	return db.DB.Where("test_case_id = ?", tcID).Order("sort_order").Find(list).Error
}
func SaveTestCaseAPIs(tcID uint, apis []model.TestCaseAPI) error {
	db.DB.Where("test_case_id = ?", tcID).Delete(&model.TestCaseAPI{})
	for i := range apis {
		apis[i].TestCaseID = tcID
	}
	return db.DB.Create(&apis).Error
}

// TestDataSet
func GetTestDataSets(tcaID uint, list *[]model.TestDataSet) error {
	return db.DB.Where("test_case_api_id = ?", tcaID).Order("sort_order").Find(list).Error
}
func SaveTestDataSets(tcaID uint, datasets []model.TestDataSet) error {
	db.DB.Where("test_case_api_id = ?", tcaID).Delete(&model.TestDataSet{})
	for i := range datasets {
		datasets[i].TestCaseAPIID = tcaID
	}
	return db.DB.Create(&datasets).Error
}

// TestSuite
func CreateTestSuite(ts *model.TestSuite) error { return db.DB.Create(ts).Error }
func GetTestSuitesByProject(pid uint, list *[]model.TestSuite) error {
	return db.DB.Where("project_id = ?", pid).Find(list).Error
}
func GetTestSuiteByID(id uint, ts *model.TestSuite) error { return db.DB.First(ts, id).Error }
func UpdateTestSuite(ts *model.TestSuite) error { return db.DB.Save(ts).Error }
func DeleteTestSuite(id uint) error { return db.DB.Delete(&model.TestSuite{}, id).Error }

// TestSuiteCase
func GetTestSuiteCases(tsID uint, list *[]model.TestSuiteCase) error {
	return db.DB.Where("test_suite_id = ?", tsID).Order("sort_order").Find(list).Error
}
func SaveTestSuiteCases(tsID uint, cases []model.TestSuiteCase) error {
	db.DB.Where("test_suite_id = ?", tsID).Delete(&model.TestSuiteCase{})
	for i := range cases {
		cases[i].TestSuiteID = tsID
	}
	return db.DB.Create(&cases).Error
}

// ScheduledTask
func CreateScheduledTask(st *model.ScheduledTask) error { return db.DB.Create(st).Error }
func GetScheduledTasksByProject(pid uint, list *[]model.ScheduledTask) error {
	return db.DB.Where("project_id = ?", pid).Find(list).Error
}
func GetScheduledTaskByID(id uint, st *model.ScheduledTask) error { return db.DB.First(st, id).Error }
func UpdateScheduledTask(st *model.ScheduledTask) error { return db.DB.Save(st).Error }
func DeleteScheduledTask(id uint) error { return db.DB.Delete(&model.ScheduledTask{}, id).Error }

// TestRun
func CreateTestRun(tr *model.TestRun) error { return db.DB.Create(tr).Error }
func GetTestRunByID(id uint, tr *model.TestRun) error { return db.DB.First(tr, id).Error }
func UpdateTestRun(tr *model.TestRun) error { return db.DB.Save(tr).Error }
func GetTestRunDB() *gorm.DB { return db.DB }

// TestRunDetail
func CreateTestRunDetail(d *model.TestRunDetail) error { return db.DB.Create(d).Error }
func GetTestRunDetails(runID uint, list *[]model.TestRunDetail) error {
	return db.DB.Where("test_run_id = ?", runID).Find(list).Error
}
func UpdateTestRunDetail(d *model.TestRunDetail) error { return db.DB.Save(d).Error }
