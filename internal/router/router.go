package router

import (
	"api-workbench/internal/handler"

	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine) {
	api := r.Group("/api/v1")

	api.GET("/health", handler.Health)

	// Project
	api.GET("/projects", handler.ListProjects)
	api.POST("/projects", handler.CreateProject)
	api.PUT("/projects/:id", handler.UpdateProject)
	api.DELETE("/projects/:id", handler.DeleteProject)

	// Environment
	api.GET("/projects/:pid/environments", handler.ListEnvironments)
	api.POST("/projects/:pid/environments", handler.CreateEnvironment)
	api.PUT("/environments/:id", handler.UpdateEnvironment)
	api.DELETE("/environments/:id", handler.DeleteEnvironment)

	// Environment Variables
	api.GET("/environments/:eid/variables", handler.ListEnvVars)
	api.PUT("/environments/:eid/variables", handler.SaveEnvVars)

	// Collection
	api.GET("/projects/:pid/collections", handler.ListCollections)
	api.POST("/projects/:pid/collections", handler.CreateCollection)
	api.PUT("/collections/:id", handler.UpdateCollection)
	api.DELETE("/collections/:id", handler.DeleteCollection)
	api.POST("/collections/:id/move", handler.MoveCollection)

	// API
	api.GET("/collections/:cid/apis", handler.ListAPIs)
	api.POST("/collections/:cid/apis", handler.CreateAPI)
	api.GET("/apis/:id", handler.GetAPI)
	api.PUT("/apis/:id", handler.UpdateAPI)
	api.DELETE("/apis/:id", handler.DeleteAPI)
	api.PUT("/apis/:id/assertions", handler.SaveAssertions)
	api.GET("/apis/:id/assertions", handler.ListAssertions)

	// TestCase
	api.GET("/projects/:pid/test-cases", handler.ListTestCases)
	api.POST("/projects/:pid/test-cases", handler.CreateTestCase)
	api.PUT("/test-cases/:id", handler.UpdateTestCase)
	api.DELETE("/test-cases/:id", handler.DeleteTestCase)
	api.PUT("/test-cases/:id/apis", handler.SaveTestCaseAPIs)

	// TestDataSet
	api.GET("/test-case-apis/:id/datasets", handler.ListDataSets)
	api.PUT("/test-case-apis/:id/datasets", handler.SaveDataSets)

	// TestSuite
	api.GET("/projects/:pid/test-suites", handler.ListTestSuites)
	api.POST("/projects/:pid/test-suites", handler.CreateTestSuite)
	api.PUT("/test-suites/:id", handler.UpdateTestSuite)
	api.DELETE("/test-suites/:id", handler.DeleteTestSuite)
	api.PUT("/test-suites/:id/cases", handler.SaveTestSuiteCases)

	// ScheduledTask
	api.GET("/projects/:pid/schedules", handler.ListScheduledTasks)
	api.POST("/projects/:pid/schedules", handler.CreateScheduledTask)
	api.PUT("/schedules/:id", handler.UpdateScheduledTask)
	api.DELETE("/schedules/:id", handler.DeleteScheduledTask)
}
