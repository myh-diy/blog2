package parser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"go.abhg.dev/goldmark/toc"
)

// ParseResult holds the structured output of parsing a Markdown document.
type ParseResult struct {
	Title   string   `json:"title"`
	Date    string   `json:"date"`
	Tags    []string `json:"tags"`
	HTML    string   `json:"html"`
	TOCJSON string   `json:"toc"`
	Slug    string   `json:"slug"`
}

// ParseMarkdown parses a Markdown document with YAML frontmatter and returns a
// structured ParseResult. It extracts frontmatter metadata (title, date, tags),
// renders the body to HTML, builds a table-of-contents JSON, and generates a
// URL-safe slug from the title.
func ParseMarkdown(mdContent []byte) (*ParseResult, error) {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Table,
			meta.Meta,
		),
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
		goldmark.WithRendererOptions(
				html.WithHardWraps(),
				html.WithUnsafe(), // allow raw HTML like <img> tags
			),
	)

	// Parse the AST once — share it for HTML rendering, meta extraction, and
	// TOC inspection.
	ctx := parser.NewContext()
	doc := md.Parser().Parse(text.NewReader(mdContent), parser.WithContext(ctx))

	// Render HTML from the parsed AST.
	var buf bytes.Buffer
	if err := md.Renderer().Render(&buf, mdContent, doc); err != nil {
		return nil, fmt.Errorf("markdown render: %w", err)
	}

	// Extract frontmatter metadata from the parser context.
	metaData := meta.Get(ctx)
	if metaData == nil {
		metaData = map[string]interface{}{}
	}

	// Parse title — prefer frontmatter, fall back to first h1 in content.
	title := getStringField(metaData, "title")
	if title == "" {
		title = extractFirstH1(mdContent)
	}

	// Parse date — prefer frontmatter, fall back to today.
	date := getStringField(metaData, "date")
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	// Parse tags — supports both YAML list and comma-separated string formats.
	var tags []string
	switch v := metaData["tags"].(type) {
	case []interface{}:
		for _, t := range v {
			tags = append(tags, fmt.Sprintf("%v", t))
		}
	case string:
		for _, t := range strings.Split(v, ",") {
			t = strings.TrimSpace(t)
			if t != "" {
				tags = append(tags, t)
			}
		}
	}

	// Extract table of contents from the AST.
	tocTree, err := toc.Inspect(doc, mdContent)
	tocJSON := "[]"
	if err == nil && tocTree != nil && len(tocTree.Items) > 0 {
		tocJSON = tocToJSON(tocTree, 0)
	}

	// Generate a URL-safe slug from the title.
	slug := slugify(title)

	return &ParseResult{
		Title:   title,
		Date:    date,
		Tags:    tags,
		HTML:    buf.String(),
		TOCJSON: tocJSON,
		Slug:    slug,
	}, nil
}

// getStringField returns the string representation of a value stored under key
// in the frontmatter map, or an empty string if the key is absent.
func getStringField(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		return fmt.Sprintf("%v", v)
	}
	return ""
}

// extractFirstH1 scans the raw Markdown content for the first level-1 heading
// ("# ...") and returns its text. Returns "Untitled" when no h1 is found.
func extractFirstH1(content []byte) string {
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "# ") {
			return strings.TrimPrefix(line, "# ")
		}
	}
	return "Untitled"
}

// slugify converts a title string into a URL-safe slug by lowercasing,
// replacing spaces with hyphens, and stripping non-alphanumeric characters
// (except hyphens). Consecutive hyphens are collapsed.
func slugify(title string) string {
	s := strings.ToLower(title)
	s = strings.ReplaceAll(s, " ", "-")
	var buf strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			buf.WriteRune(r)
		}
	}
	result := buf.String()
	for strings.Contains(result, "--") {
		result = strings.ReplaceAll(result, "--", "-")
	}
	return strings.Trim(result, "-")
}

// TOCItem represents a single heading entry in the table-of-contents JSON
// output.
type TOCItem struct {
	ID       string    `json:"id"`
	Text     string    `json:"text"`
	Level    int       `json:"level"`
	Children []TOCItem `json:"children"`
}

// tocToJSON converts a toc.TOC tree into its JSON string representation.
// depth tracks the current heading level, starting from 1 for top-level items.
func tocToJSON(t *toc.TOC, depth int) string {
	items := buildTOCItems(t.Items, depth+1)
	b, _ := json.Marshal(items)
	return string(b)
}

// buildTOCItems recursively converts toc.Items into a slice of TOCItem,
// assigning level based on the current nesting depth.
func buildTOCItems(items toc.Items, level int) []TOCItem {
	var result []TOCItem
	for _, item := range items {
		ti := TOCItem{
			ID:       string(item.ID),
			Text:     string(item.Title),
			Level:    level,
			Children: buildTOCItems(item.Items, level+1),
		}
		if ti.Children == nil {
			ti.Children = []TOCItem{}
		}
		result = append(result, ti)
	}
	return result
}
