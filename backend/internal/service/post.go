package service

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"gorm.io/gorm"

	"blog-backend/internal/model"
	"blog-backend/internal/parser"
)

// CreatePost creates a new post from parsed Markdown, associates tags, and
// syncs the FTS index. Returns the created post or an error.
func CreatePost(db *gorm.DB, result *parser.ParseResult, rawMD, sourceFile string) (*model.Post, error) {
	// Resolve slug collisions by trying incrementing suffixes until unique
	slug := result.Slug
	base := slug
	for i := 1; ; i++ {
		var count int64
		db.Model(&model.Post{}).Where("slug = ?", slug).Count(&count)
		if count == 0 {
			break
		}
		slug = fmt.Sprintf("%s-%d", base, i+1)
	}

	// Parse date, fall back to time.Now() if malformed
	parsedDate, err := parseDate(result.Date)
	if err != nil {
		parsedDate = time.Now()
		log.Printf("WARNING: invalid date format '%s', falling back to %s: %v",
			result.Date, parsedDate.Format("2006-01-02"), err)
	}

	// Find or create tags
	var tags []model.Tag
	for _, tagName := range result.Tags {
		var tag model.Tag
		db.Where("name = ?", tagName).FirstOrCreate(&tag, model.Tag{Name: tagName})
		tags = append(tags, tag)
	}

	post := &model.Post{
		Title:       result.Title,
		Slug:        slug,
		SourceFile:  sourceFile,
		ContentMD:   rawMD,
		ContentHTML: result.HTML,
		TOCJSON:     result.TOCJSON,
		CreatedAt:   parsedDate,
		UpdatedAt:   parsedDate,
		Tags:        tags,
	}

	if err := db.Create(post).Error; err != nil {
		return nil, err
	}

	// Sync FTS index
	syncFTS(db, post.ID, post.Title, rawMD)

	return post, nil
}

// GetPosts returns a paginated list of posts, optionally filtered by tag.
func GetPosts(db *gorm.DB, page, perPage int, tag string) ([]model.Post, int64, error) {
	var posts []model.Post
	var total int64

	query := db.Model(&model.Post{}).Preload("Tags")
	if tag != "" {
		query = query.Joins("JOIN post_tags ON post_tags.post_id = posts.id").
			Joins("JOIN tags ON tags.id = post_tags.tag_id").
			Where("tags.name = ?", tag)
	}
	query.Count(&total).
		Order("created_at DESC").
		Offset((page - 1) * perPage).
		Limit(perPage).
		Find(&posts)

	return posts, total, query.Error
}

// GetPostBySlug returns a single post by its URL slug, including associated tags.
func GetPostBySlug(db *gorm.DB, slug string) (*model.Post, error) {
	var post model.Post
	if err := db.Preload("Tags").Where("slug = ?", slug).First(&post).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

// GetPostByID returns a single post by its primary key, including associated tags.
func GetPostByID(db *gorm.DB, id uint) (*model.Post, error) {
	var post model.Post
	if err := db.Preload("Tags").First(&post, id).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

// UpdatePost applies partial updates to an existing post and returns the
// updated post. The updates map uses GORM's Updates semantics.
func UpdatePost(db *gorm.DB, id uint, updates map[string]interface{}) (*model.Post, error) {
	post, err := GetPostByID(db, id)
	if err != nil {
		return nil, err
	}
	if err := db.Model(post).Updates(updates).Error; err != nil {
		return nil, err
	}
	return GetPostByID(db, id)
}

// DeletePost removes a post, its tag associations, and its uploaded images.
func DeletePost(db *gorm.DB, id uint) error {
	var post model.Post
	if err := db.First(&post, id).Error; err != nil {
		return err
	}
	// Delete associated uploaded images
	deletePostImages(post.ContentHTML)
	// Clear tag associations then delete
	db.Model(&post).Association("Tags").Clear()
	return db.Unscoped().Delete(&post).Error
}

// SearchPosts performs a full-text search on posts using the FTS5 index.
// Returns matching posts and the total count.
func SearchPosts(db *gorm.DB, query string, page, perPage int) ([]model.Post, int64, error) {
	var ids []uint
	var total int64

	db.Raw("SELECT rowid FROM posts_fts WHERE posts_fts MATCH ? ORDER BY rank LIMIT ? OFFSET ?",
		query, perPage, (page-1)*perPage).Scan(&ids)
	db.Raw("SELECT COUNT(*) FROM posts_fts WHERE posts_fts MATCH ?", query).Scan(&total)

	var posts []model.Post
	if len(ids) > 0 {
		db.Preload("Tags").Where("id IN ?", ids).Order("created_at DESC").Find(&posts)
	}

	return posts, total, nil
}

// GetTimeline builds a blog archive grouped by year and month, ordered from
// most recent to oldest.
func GetTimeline(db *gorm.DB) ([]model.TimelineEntry, error) {
	var posts []model.Post
	if err := db.Preload("Tags").Order("created_at DESC").Find(&posts).Error; err != nil {
		return nil, err
	}

	// Group posts by year-month, preserving descending order from the query
	groups := make(map[string]*model.TimelineEntry)
	var order []string
	for _, p := range posts {
		key := p.CreatedAt.Format("2006-01")
		if groups[key] == nil {
			groups[key] = &model.TimelineEntry{
				Year:  p.CreatedAt.Year(),
				Month: int(p.CreatedAt.Month()),
				Posts: []model.Post{},
			}
			order = append(order, key)
		}
		groups[key].Posts = append(groups[key].Posts, p)
	}

	result := make([]model.TimelineEntry, 0, len(order))
	for _, key := range order {
		result = append(result, *groups[key])
	}
	return result, nil
}

// --- helpers ---

func parseDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}

// UpdatePostTags replaces the tag associations for a post
func UpdatePostTags(db *gorm.DB, postID uint, tagNames []string) error {
	post, err := GetPostByID(db, postID)
	if err != nil {
		return err
	}
	// Clear existing tags
	db.Model(post).Association("Tags").Clear()
	// Set new tags
	tags, err := EnsureTags(db, tagNames)
	if err != nil {
		return err
	}
	return db.Model(post).Association("Tags").Replace(tags)
}

func syncFTS(db *gorm.DB, postID uint, title, contentMD string) {
	db.Exec("DELETE FROM posts_fts WHERE rowid = ?", postID)
	db.Exec("INSERT INTO posts_fts(rowid, title, content_md) VALUES(?, ?, ?)", postID, title, contentMD)
}

// deletePostImages removes uploaded image files referenced in the post HTML.
func deletePostImages(html string) {
	re := regexp.MustCompile(`/uploads/([^")\s]+)`)
	matches := re.FindAllStringSubmatch(html, -1)
	for _, m := range matches {
		if len(m) < 2 {
			continue
		}
		path := filepath.Join("uploads", m[1])
		if err := os.Remove(path); err != nil {
			log.Printf("Failed to delete image %s: %v", path, err)
		}
	}
}

// GetRandomPostQuotes extracts random sentences from all posts' markdown content.
// Returns up to `limit` sentences (roughly 20-80 chars each).
func GetRandomPostQuotes(db *gorm.DB, limit int) []string {
	var posts []model.Post
	if err := db.Select("content_md").Find(&posts).Error; err != nil || len(posts) == 0 {
		return nil
	}

	var sentences []string
	for _, p := range posts {
		for _, s := range splitSentences(p.ContentMD) {
			s = strings.TrimSpace(s)
			// Filter: keep sentences between 10 and 120 chars, skip headings/code
			runeCount := utf8.RuneCountInString(s)
			if runeCount < 10 || runeCount > 120 {
				continue
			}
			if strings.HasPrefix(s, "#") || strings.HasPrefix(s, "```") || strings.HasPrefix(s, "---") {
				continue
			}
			sentences = append(sentences, s)
		}
	}

	if len(sentences) == 0 {
		return nil
	}

	// Shuffle and pick first `limit`
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	rng.Shuffle(len(sentences), func(i, j int) {
		sentences[i], sentences[j] = sentences[j], sentences[i]
	})

	n := limit
	if n > len(sentences) {
		n = len(sentences)
	}
	return sentences[:n]
}

// splitSentences splits text into sentences by common delimiters.
func splitSentences(text string) []string {
	var result []string
	current := strings.Builder{}
	runes := []rune(text)

	for i := 0; i < len(runes); i++ {
		r := runes[i]
		current.WriteRune(r)

		// Split on sentence-ending punctuation followed by space/newline
		if (r == '.' || r == '。' || r == '！' || r == '!' || r == '？' || r == '?' || r == '\n') &&
			(i+1 >= len(runes) || unicode.IsSpace(runes[i+1]) || runes[i+1] == '\n') {
			s := current.String()
			if len(s) > 0 {
				result = append(result, s)
			}
			current.Reset()
		}
	}
	if current.Len() > 0 {
		result = append(result, current.String())
	}
	return result
}
