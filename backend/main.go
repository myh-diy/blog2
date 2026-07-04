package main

import (
	"blog-backend/config"
	"blog-backend/internal/database"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	database.Init(cfg)

	r := gin.Default()
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	fmt.Println("Server starting on :" + cfg.Port)
	r.Run(":" + cfg.Port)
}
