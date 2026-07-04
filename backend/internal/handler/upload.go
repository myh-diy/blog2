package handler

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"blog-backend/internal/database"
	"blog-backend/internal/parser"
	"blog-backend/internal/service"
)

func UploadPost() gin.HandlerFunc {
	return func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
			return
		}
		f, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open file"})
			return
		}
		defer f.Close()

		content, err := io.ReadAll(f)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file"})
			return
		}

		result, err := parser.ParseMarkdown(content)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse markdown: " + err.Error()})
			return
		}

		post, err := service.CreatePost(database.DB, result, string(content))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create post: " + err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"post": post})
	}
}
