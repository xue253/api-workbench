package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	r.GET("/api/items", func(c *gin.Context) {
		items := []map[string]interface{}{
			{"id": 1, "name": "Item 1", "value": 100},
			{"id": 2, "name": "Item 2", "value": 200},
			{"id": 3, "name": "Item 3", "value": 300},
		}
		c.JSON(http.StatusOK, items)
	})

	r.Run(":8080")
}
