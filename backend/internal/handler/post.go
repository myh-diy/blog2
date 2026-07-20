package handler

import (
	"net/http"
	"os"
	"strconv"
	"strings"

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

func ExportPost() gin.HandlerFunc {
	return func(c *gin.Context) {
		post, err := service.GetPostBySlug(database.DB, c.Param("slug"))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
			return
		}

		filename := post.Slug
		if filename == "" {
			filename = "post"
		}
		filename = strings.ReplaceAll(filename, "\"", "") + ".md"
		c.Header("Content-Disposition", `attachment; filename="`+filename+`"`)
		c.Data(http.StatusOK, "text/markdown; charset=utf-8", []byte(post.ContentMD))
	}
}

func GetPostSource() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		post, err := service.GetPostByID(database.DB, uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"content": post.ContentMD})
	}
}

func UpdatePostContent() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		var req struct {
			Content string `json:"content" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "markdown content is required"})
			return
		}
		post, err := service.UpdatePostContent(database.DB, uint(id), req.Content)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
		form, err := c.MultipartForm()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid form data"})
			return
		}

		updates := make(map[string]interface{})
		title := c.PostForm("title")
		if title != "" {
			updates["title"] = title
		}

		if coverFile, err := c.FormFile("cover_image"); err == nil {
			os.MkdirAll(uploadDir, 0755)
			savedName := saveUploadedFile(coverFile)
			if savedName == "" {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save cover image"})
				return
			}
			updates["cover_image"] = "/uploads/" + savedName
		}

		// Handle tags separately for many-to-many
		if tagNames, ok := form.Value["tags[]"]; ok {
			if err := service.UpdatePostTags(database.DB, uint(id), tagNames); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
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
