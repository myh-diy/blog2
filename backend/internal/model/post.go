package model

import "time"

type Post struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"size:500;not null" json:"title"`
	Slug        string    `gorm:"uniqueIndex;size:500;not null" json:"slug"`
	SourceFile  string    `gorm:"size:500" json:"source_file"`
	ContentMD   string    `gorm:"type:text;not null" json:"-"`
	ContentHTML string    `gorm:"type:text;not null" json:"content_html"`
	TOCJSON     string    `gorm:"type:text" json:"toc"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Tags        []Tag     `gorm:"many2many:post_tags;" json:"tags"`
}

type TimelineEntry struct {
	Year  int    `json:"year"`
	Month int    `json:"month"`
	Posts []Post `json:"posts"`
}
