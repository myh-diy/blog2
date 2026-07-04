package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"blog-backend/internal/database"
	"blog-backend/internal/model"
	"blog-backend/internal/service"
)

func GetQuotes() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Try quote table first, fall back to extracting from posts
		quotes := service.GetRandomQuotes(database.DB, 20)
		if quotes == nil {
			quotes = service.GetRandomPostQuotes(database.DB, 20)
		}
		if quotes == nil {
			quotes = []string{}
		}
		c.JSON(http.StatusOK, gin.H{"quotes": quotes})
	}
}

// Admin quote management

func ListQuotes() gin.HandlerFunc {
	return func(c *gin.Context) {
		quotes, err := service.GetAllQuotes(database.DB)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if quotes == nil {
			quotes = []model.Quote{}
		}
		c.JSON(http.StatusOK, gin.H{"quotes": quotes})
	}
}

func CreateQuote() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Text string `json:"text" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "text is required"})
			return
		}
		quote, err := service.CreateQuote(database.DB, req.Text)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"quote": quote})
	}
}

func DeleteQuote() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		if err := service.DeleteQuote(database.DB, uint(id)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusNoContent)
	}
}
