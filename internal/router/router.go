package router

import (
	"api-workbench/internal/handler"
	"api-workbench/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine) {
	api := r.Group("/api/v1")

	api.GET("/health", handler.Health)

	// Auth (public)
	api.POST("/auth/register", handler.Register)
	api.POST("/auth/login", handler.Login)

	// Protected routes
	auth := api.Group("")
	auth.Use(middleware.Auth())
	{
		// User
		auth.GET("/user/profile", handler.GetProfile)
		auth.PUT("/user/profile", handler.UpdateProfile)
		auth.PUT("/user/password", handler.ChangePassword)
		auth.DELETE("/user/account", handler.DeleteAccount)

		// Project
		auth.GET("/projects", handler.ListProjects)
		auth.POST("/projects", handler.CreateProject)

		// Project sub-resources
		auth.GET("/projects/:pid/environments", handler.ListEnvironments)
		auth.POST("/projects/:pid/environments", handler.CreateEnvironment)
		auth.GET("/projects/:pid/collections", handler.ListCollections)
		auth.POST("/projects/:pid/collections", handler.CreateCollection)
		auth.GET("/projects/:pid/test-cases", handler.ListTestCases)
		auth.POST("/projects/:pid/test-cases", handler.CreateTestCase)
		auth.GET("/projects/:pid/test-suites", handler.ListTestSuites)
		auth.POST("/projects/:pid/test-suites", handler.CreateTestSuite)
		auth.GET("/projects/:pid/schedules", handler.ListScheduledTasks)
		auth.POST("/projects/:pid/schedules", handler.CreateScheduledTask)

		// Project CRUD
		auth.PUT("/projects/:id", handler.UpdateProject)
		auth.DELETE("/projects/:id", handler.DeleteProject)

		// Environment
		env := auth.Group("/environments/:id")
		{
			env.PUT("", handler.UpdateEnvironment)
			env.DELETE("", handler.DeleteEnvironment)
			env.GET("/variables", handler.ListEnvVars)
			env.PUT("/variables", handler.SaveEnvVars)
		}

		// Collection
		col := auth.Group("/collections/:id")
		{
			col.PUT("", handler.UpdateCollection)
			col.DELETE("", handler.DeleteCollection)
			col.POST("/move", handler.MoveCollection)
			col.GET("/apis", handler.ListAPIsByCollection)
			col.POST("/apis", handler.CreateAPIByCollection)
		}

		// API
		apir := auth.Group("/apis/:id")
		{
			apir.GET("", handler.GetAPI)
			apir.PUT("", handler.UpdateAPI)
			apir.DELETE("", handler.DeleteAPI)
			apir.GET("/assertions", handler.ListAssertions)
			apir.PUT("/assertions", handler.SaveAssertions)
		}

		// TestCase
		tc := auth.Group("/test-cases/:id")
		{
			tc.PUT("", handler.UpdateTestCase)
			tc.DELETE("", handler.DeleteTestCase)
			tc.PUT("/apis", handler.SaveTestCaseAPIs)
		}

		// TestDataSet
		auth.GET("/test-case-apis/:id/datasets", handler.ListDataSets)
		auth.PUT("/test-case-apis/:id/datasets", handler.SaveDataSets)

		// TestSuite
		ts := auth.Group("/test-suites/:id")
		{
			ts.PUT("", handler.UpdateTestSuite)
			ts.DELETE("", handler.DeleteTestSuite)
			ts.PUT("/cases", handler.SaveTestSuiteCases)
		}

		// ScheduledTask
		auth.PUT("/schedules/:id", handler.UpdateScheduledTask)
		auth.DELETE("/schedules/:id", handler.DeleteScheduledTask)
	}
}
