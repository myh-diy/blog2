package service

import (
	"gorm.io/gorm"

	"blog-backend/internal/model"
)

// TagWithCount pairs a tag name with the number of posts that use it.
type TagWithCount struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// GetAllTags returns all tags ordered by post count descending.
func GetAllTags(db *gorm.DB) ([]TagWithCount, error) {
	var tags []TagWithCount
	err := db.Raw(`
		SELECT tags.name, COUNT(post_tags.post_id) as count
		FROM tags
		JOIN post_tags ON post_tags.tag_id = tags.id
		GROUP BY tags.id
		ORDER BY count DESC
	`).Scan(&tags).Error
	return tags, err
}

// EnsureTags finds existing tags by name or creates them if missing.
// Returns the complete set of Tag model objects.
func EnsureTags(db *gorm.DB, tagNames []string) ([]model.Tag, error) {
	var tags []model.Tag
	for _, name := range tagNames {
		var tag model.Tag
		db.Where("name = ?", name).FirstOrCreate(&tag, model.Tag{Name: name})
		tags = append(tags, tag)
	}
	return tags, nil
}
