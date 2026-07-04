package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"blog-backend/internal/database"
	"blog-backend/internal/service"
)

func ListTags() gin.HandlerFunc {
	return func(c *gin.Context) {
		tags, err := service.GetAllTags(database.DB)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"tags": tags})
	}
}
