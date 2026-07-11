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

//go:embed frontend-dist/* frontend-dist/**/*
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
	r.Static("/uploads", "./uploads")
	r.GET("/api/quotes", handler.GetQuotes())
	r.GET("/api/settings", handler.GetSettings())
	r.GET("/api/system/metrics", handler.GetSystemMetrics(cfg.ExporterURL))

	// Admin routes (protected)
	admin := r.Group("/api/admin")
	admin.Use(auth.AuthMiddleware(cfg.JWTSecret))
	{
		admin.POST("/upload", handler.UploadPost())
		admin.POST("/upload-image", handler.UploadImage())
		admin.PUT("/posts/:id", handler.UpdatePost())
		admin.DELETE("/posts/:id", handler.DeletePost())
		admin.GET("/quotes", handler.ListQuotes())
		admin.POST("/quotes", handler.CreateQuote())
		admin.DELETE("/quotes/:id", handler.DeleteQuote())
		admin.PUT("/settings", handler.UpdateSettings())
	}

	// Serve SPA static files with fallback to index.html
	frontendFS, err := fs.Sub(staticFiles, "frontend-dist")
	if err != nil {
		log.Fatalf("Failed to load frontend files: %v", err)
	}

	// Pre-read index.html for SPA fallback
	indexHTML, err := fs.ReadFile(frontendFS, "index.html")
	if err != nil {
		log.Fatalf("Failed to read index.html: %v", err)
	}

	httpFS := http.FS(frontendFS)

	// SPA routes: serve file if exists, otherwise serve index.html
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/api/") {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		// Try to serve exact file
		trimmedPath := strings.TrimPrefix(path, "/")
		if _, err := frontendFS.Open(trimmedPath); err == nil {
			c.FileFromFS(path, httpFS)
			return
		}
		// SPA fallback: send index.html with proper content type
		c.Data(http.StatusOK, "text/html; charset=utf-8", indexHTML)
	})

	fmt.Println("Server starting on :" + cfg.Port)
	r.Run(":" + cfg.Port)
}
