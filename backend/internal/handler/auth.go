package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"blog-backend/config"
	"blog-backend/internal/auth"
	"blog-backend/internal/database"
	"blog-backend/internal/service"
)

func Login(cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}
		user, err := service.AuthenticateUser(database.DB, req.Username, req.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		token, expiresAt, err := auth.GenerateToken(user.ID, cfg.JWTSecret)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"token":      token,
			"expires_at": expiresAt,
		})
	}
}
