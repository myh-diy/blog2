package handler

import (
	"encoding/xml"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"blog-backend/internal/database"
	"blog-backend/internal/model"
)

type rss struct {
	XMLName xml.Name   `xml:"rss"`
	Version string     `xml:"version,attr"`
	Channel rssChannel `xml:"channel"`
}

type rssChannel struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Items       []rssItem `xml:"item"`
}

type rssItem struct {
	Title   string `xml:"title"`
	Link    string `xml:"link"`
	GUID    string `xml:"guid"`
	PubDate string `xml:"pubDate"`
}

type sitemapURLSet struct {
	XMLName xml.Name     `xml:"urlset"`
	Xmlns   string       `xml:"xmlns,attr"`
	URLs    []sitemapURL `xml:"url"`
}

type sitemapURL struct {
	Location     string `xml:"loc"`
	LastModified string `xml:"lastmod,omitempty"`
}

func RSSFeed(publicURL string) gin.HandlerFunc {
	return func(c *gin.Context) {
		posts, err := publishedPosts()
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		items := make([]rssItem, 0, len(posts))
		for _, post := range posts {
			link := publicURL + "/post/" + post.Slug
			items = append(items, rssItem{Title: post.Title, Link: link, GUID: link, PubDate: post.CreatedAt.Format(time.RFC1123Z)})
		}
		writeXML(c, rss{Version: "2.0", Channel: rssChannel{
			Title: "Blog", Link: publicURL, Description: "Latest published posts", Items: items,
		}})
	}
}

func Sitemap(publicURL string) gin.HandlerFunc {
	return func(c *gin.Context) {
		posts, err := publishedPosts()
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		urls := []sitemapURL{{Location: publicURL}}
		for _, post := range posts {
			urls = append(urls, sitemapURL{
				Location: publicURL + "/post/" + post.Slug, LastModified: post.UpdatedAt.Format("2006-01-02"),
			})
		}
		writeXML(c, sitemapURLSet{Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9", URLs: urls})
	}
}

func publishedPosts() ([]model.Post, error) {
	var posts []model.Post
	err := database.DB.Where("published = ?", true).Order("created_at DESC").Find(&posts).Error
	return posts, err
}

func writeXML(c *gin.Context, value any) {
	data, err := xml.MarshalIndent(value, "", "  ")
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Data(http.StatusOK, "application/xml; charset=utf-8", append([]byte(strings.TrimSpace(xml.Header)), data...))
}
