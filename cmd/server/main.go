package main

import (
	"fmt"
	"log"

	"api-workbench/internal/config"
	"api-workbench/internal/db"
	"api-workbench/internal/middleware"
	"api-workbench/internal/router"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := config.Load("config.yaml"); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	gin.SetMode(config.AppConfig.Server.Mode)

	db.Init()

	r := gin.Default()
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())

	router.Setup(r)

	addr := fmt.Sprintf(":%d", config.AppConfig.Server.Port)
	log.Printf("server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
