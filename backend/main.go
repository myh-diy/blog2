package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"blog-backend/config"
	"blog-backend/internal/auth"
	"blog-backend/internal/database"
	"blog-backend/internal/handler"
	"blog-backend/internal/service"
)

//go:embed frontend-dist/*
var staticFiles embed.FS

func main() {
	cfg := config.Load()
	database.Init(cfg)

	if err := service.CreateDefaultUser(database.DB); err != nil {
		log.Fatalf("Failed to create default user: %v", err)
	}
	log.Println("Default admin user ensured (username: admin, password: admin)")

	r := gin.Default()

	// Public routes
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.POST("/api/auth/login", handler.Login(cfg))

	r.GET("/api/posts", handler.ListPosts())
	r.GET("/api/posts/:slug", handler.GetPost())
	r.GET("/api/tags", handler.ListTags())
	r.GET("/api/search", handler.SearchPosts())
	r.GET("/api/timeline", handler.GetTimeline())

	// Admin routes (protected)
	admin := r.Group("/api/admin")
	admin.Use(auth.AuthMiddleware(cfg.JWTSecret))
	{
		admin.POST("/upload", handler.UploadPost())
		admin.PUT("/posts/:id", handler.UpdatePost())
		admin.DELETE("/posts/:id", handler.DeletePost())
	}

	// Serve SPA static files with fallback to index.html
	frontendFS, err := fs.Sub(staticFiles, "frontend-dist")
	if err != nil {
		log.Fatalf("Failed to load frontend files: %v", err)
	}

	// Explicit root route to avoid Gin's 301 redirect
	r.GET("/", func(c *gin.Context) {
		c.FileFromFS("index.html", http.FS(frontendFS))
	})

	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		// Don't intercept API routes
		if strings.HasPrefix(path, "/api/") {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		// Try to serve the requested file, fallback to index.html
		f, err := frontendFS.Open(strings.TrimPrefix(path, "/"))
		if err == nil {
			f.Close()
			c.FileFromFS(path, http.FS(frontendFS))
			return
		}
		c.FileFromFS("index.html", http.FS(frontendFS))
	})

	fmt.Println("Server starting on :" + cfg.Port)
	r.Run(":" + cfg.Port)
}
