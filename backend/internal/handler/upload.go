package handler

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"blog-backend/internal/database"
	"blog-backend/internal/parser"
	"blog-backend/internal/service"
)

const uploadDir = "uploads"

func UploadPost() gin.HandlerFunc {
	return func(c *gin.Context) {
		form, err := c.MultipartForm()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid form data"})
			return
		}

		files := form.File["file"]
		if len(files) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "at least one file is required"})
			return
		}

		// Separate .md file from image files
		var mdFile *multipart.FileHeader
		var imgFiles []*multipart.FileHeader
		for _, f := range files {
			ext := strings.ToLower(filepath.Ext(f.Filename))
			if ext == ".md" || ext == ".markdown" {
				if mdFile == nil {
					mdFile = f
				}
			} else {
				imgFiles = append(imgFiles, f)
			}
		}

		if mdFile == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no .md file found"})
			return
		}

		// Read markdown content
		f, err := mdFile.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open md file"})
			return
		}
		content, err := io.ReadAll(f)
		f.Close()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read md file"})
			return
		}

		// Build image map: original filename -> saved filename
		imageMap := make(map[string]string)
		os.MkdirAll(uploadDir, 0755)

		for _, img := range imgFiles {
			savedName := saveUploadedFile(img)
			if savedName != "" {
				// Map by original filename AND basename for flexible matching
				imageMap[img.Filename] = savedName
				imageMap[filepath.Base(img.Filename)] = savedName
			}
		}

		// Replace local image paths in markdown
		mdStr := string(content)
		mdStr = replaceImagePaths(mdStr, imageMap)

		result, err := parser.ParseMarkdown([]byte(mdStr))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse markdown: " + err.Error()})
			return
		}

		post, err := service.CreatePost(database.DB, result, mdStr, mdFile.Filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create post: " + err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"post":   post,
			"images": len(imageMap),
		})
	}
}

// UploadImage handles standalone image upload, returns the URL
func UploadImage() gin.HandlerFunc {
	return func(c *gin.Context) {
		file, err := c.FormFile("image")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "image is required"})
			return
		}
		os.MkdirAll(uploadDir, 0755)
		savedName := saveUploadedFile(file)
		if savedName == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save image"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"url": "/uploads/" + savedName,
		})
	}
}

func saveUploadedFile(file *multipart.FileHeader) string {
	src, err := file.Open()
	if err != nil {
		return ""
	}
	defer src.Close()

	ext := filepath.Ext(file.Filename)
	name := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), sanitizeName(file.Filename), ext)
	dst, err := os.Create(filepath.Join(uploadDir, name))
	if err != nil {
		log.Printf("Failed to create file %s: %v", name, err)
		return ""
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		log.Printf("Failed to save file %s: %v", name, err)
		return ""
	}
	return name
}

// replaceImagePaths replaces local image references with /uploads/ paths.
// Handles both Markdown ![alt](path) and HTML <img src="path"> syntax.
func replaceImagePaths(md string, imageMap map[string]string) string {
	if len(imageMap) == 0 {
		return md
	}

	// 1. Replace Markdown syntax: ![alt](path)
	mdRe := regexp.MustCompile(`!\[([^\]]*)\]\(([^)]+)\)`)
	md = mdRe.ReplaceAllStringFunc(md, func(match string) string {
		sub := mdRe.FindStringSubmatch(match)
		if len(sub) < 3 {
			return match
		}
		alt := sub[1]
		path := sub[2]
		if newPath, ok := tryMatchPath(path, imageMap); ok {
			return fmt.Sprintf("![%s](%s)", alt, "/uploads/"+newPath)
		}
		return match
	})

	// 2. Replace HTML <img> syntax: <img src="path" ...>
	imgRe := regexp.MustCompile(`<img\s+[^>]*src="([^"]+)"[^>]*>`)
	md = imgRe.ReplaceAllStringFunc(md, func(match string) string {
		sub := imgRe.FindStringSubmatch(match)
		if len(sub) < 2 {
			return match
		}
		path := sub[1]
		if newPath, ok := tryMatchPath(path, imageMap); ok {
			// Replace src and remove non-standard style attributes
			result := strings.Replace(match, `src="`+path+`"`, `src="/uploads/`+newPath+`"`, 1)
			// Remove zoom style
			styleRe := regexp.MustCompile(`\s*style="[^"]*zoom\s*:[^"]*"[^>]*`)
			result = styleRe.ReplaceAllString(result, "")
			// Clean up double spaces
			result = strings.ReplaceAll(result, "  ", " ")
			return result
		}
		return match
	})

	return md
}

// tryMatchPath attempts to match a local image path against uploaded files.
func tryMatchPath(path string, imageMap map[string]string) (string, bool) {
	// Skip external URLs and already-absolute paths
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") || strings.HasPrefix(path, "/uploads/") {
		return "", false
	}

	candidates := []string{
		path,                           // exact: "images/photo.png"
		filepath.Base(path),             // basename: "photo.png"
		strings.TrimPrefix(path, "./"),  // strip ./: "images/photo.png"
	}

	for _, c := range candidates {
		if saved, ok := imageMap[c]; ok {
			return saved, true
		}
	}
	return "", false
}

func sanitizeName(name string) string {
	name = strings.TrimSuffix(name, filepath.Ext(name))
	re := regexp.MustCompile(`[^a-zA-Z0-9_\-.]+`)
	return re.ReplaceAllString(name, "_")
}
