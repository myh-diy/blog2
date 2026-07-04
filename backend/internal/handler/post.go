package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"blog-backend/internal/database"
	"blog-backend/internal/service"
)

func ListPosts() gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
		tag := c.Query("tag")

		posts, total, err := service.GetPosts(database.DB, page, perPage, tag)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"posts":    posts,
			"total":    total,
			"page":     page,
			"per_page": perPage,
		})
	}
}

func GetPost() gin.HandlerFunc {
	return func(c *gin.Context) {
		slug := c.Param("slug")
		post, err := service.GetPostBySlug(database.DB, slug)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"post": post})
	}
}

func UpdatePost() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		var updates map[string]interface{}
		if err := c.ShouldBindJSON(&updates); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}
		post, err := service.UpdatePost(database.DB, uint(id), updates)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"post": post})
	}
}

func DeletePost() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		if err := service.DeletePost(database.DB, uint(id)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusNoContent)
	}
}
