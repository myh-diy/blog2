package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"blog-backend/internal/database"
	"blog-backend/internal/service"
)

func GetTimeline() gin.HandlerFunc {
	return func(c *gin.Context) {
		timeline, err := service.GetTimeline(database.DB)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"timeline": timeline})
	}
}
