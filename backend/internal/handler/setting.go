package handler

import (
	"net/http"
	"strconv"

	"blog-backend/internal/database"
	"blog-backend/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetSettings() gin.HandlerFunc {
	return func(c *gin.Context) {
		var settings []model.Setting
		if err := database.DB.Find(&settings).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		result := make(map[string]string, len(settings))
		for _, s := range settings {
			result[s.Key] = s.Value
		}
		c.JSON(http.StatusOK, gin.H{"settings": result})
	}
}

func UpdateSettings() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload map[string]interface{}
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := database.DB.Transaction(func(tx *gorm.DB) error {
			for key, value := range payload {
				var v string
				switch val := value.(type) {
				case string:
					v = val
				case float64:
					v = formatFloat(val)
				case bool:
					v = "false"
					if val {
						v = "true"
					}
				default:
					v = ""
				}

				var existing model.Setting
				if err := tx.Where("`key` = ?", key).First(&existing).Error; err != nil {
					if err == gorm.ErrRecordNotFound {
						if err := tx.Create(&model.Setting{Key: key, Value: v}).Error; err != nil {
							return err
						}
						continue
					}
					return err
				}

				existing.Value = v
				if err := tx.Save(&existing).Error; err != nil {
					return err
				}
			}
			return nil
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"ok": true})
	}
}

func formatFloat(f float64) string {
	if f == float64(int64(f)) {
		return strconv.FormatInt(int64(f), 10)
	}
	return strconv.FormatFloat(f, 'f', -1, 64)
}
