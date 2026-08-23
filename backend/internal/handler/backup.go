package handler

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"

	"blog-backend/config"
	"blog-backend/internal/database"
)

func DownloadBackup(cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		database.DB.Exec("PRAGMA wal_checkpoint(FULL)")
		filename := "blog-backup-" + time.Now().Format("20060102-150405") + ".zip"
		c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
		c.Header("Content-Type", "application/zip")

		zw := zip.NewWriter(c.Writer)
		if err := addBackupFile(zw, cfg.Database.Path, "data/blog.db"); err != nil {
			_ = zw.Close()
			return
		}
		_ = filepath.Walk(cfg.Storage.UploadDir, func(path string, info os.FileInfo, err error) error {
			if err != nil || info == nil || info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
				return nil
			}
			rel, relErr := filepath.Rel(cfg.Storage.UploadDir, path)
			if relErr != nil {
				return relErr
			}
			return addBackupFile(zw, path, filepath.ToSlash(filepath.Join("uploads", rel)))
		})
		_ = zw.Close()
	}
}

func addBackupFile(zw *zip.Writer, path, name string) error {
	src, err := os.Open(path)
	if err != nil {
		return err
	}
	defer src.Close()
	dst, err := zw.Create(name)
	if err != nil {
		return err
	}
	_, err = io.Copy(dst, src)
	return err
}
