package model

import "time"

type Post struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"size:500;not null" json:"title"`
	Slug        string    `gorm:"uniqueIndex;size:500;not null" json:"slug"`
	SourceFile  string    `gorm:"size:500" json:"source_file"`
	CoverImage  string    `gorm:"size:1000" json:"cover_image"`
	ContentMD   string    `gorm:"type:text;not null" json:"-"`
	ContentHTML string    `gorm:"type:text;not null" json:"content_html"`
	TOCJSON     string    `gorm:"type:text" json:"toc"`
	Published   bool      `gorm:"not null;default:true;index" json:"published"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Tags        []Tag     `gorm:"many2many:post_tags;" json:"tags"`
}

type PostRevision struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	PostID    uint      `gorm:"not null;index" json:"post_id"`
	Title     string    `gorm:"size:500;not null" json:"title"`
	ContentMD string    `gorm:"type:text;not null" json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

type TimelineEntry struct {
	Year  int    `json:"year"`
	Month int    `json:"month"`
	Posts []Post `json:"posts"`
}
