package main

import (
	"blog-backend/config"
	"blog-backend/internal/database"
	"blog-backend/internal/service"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	database.Init(cfg)

	if err := service.CreateDefaultUser(database.DB); err != nil {
		log.Fatalf("Failed to create default user: %v", err)
	}
	log.Println("Default admin user ensured (username: admin, password: admin)")

	r := gin.Default()
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	fmt.Println("Server starting on :" + cfg.Port)
	r.Run(":" + cfg.Port)
}
