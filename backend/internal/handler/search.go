package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"blog-backend/internal/database"
	"blog-backend/internal/model"
	"blog-backend/internal/service"
)

func SearchPosts() gin.HandlerFunc {
	return func(c *gin.Context) {
		q := c.Query("q")
		if q == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "query is required"})
			return
		}
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))

		posts, total, err := service.SearchPosts(database.DB, q, page, perPage)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if posts == nil {
			posts = []model.Post{}
		}
		c.JSON(http.StatusOK, gin.H{
			"posts":    posts,
			"total":    total,
			"page":     page,
			"per_page": perPage,
		})
	}
}
