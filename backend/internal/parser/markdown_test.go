package parser

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestParseMarkdown(t *testing.T) {
	input := `---
title: Learning Go
date: 2026-07-04
tags: [Go, backend]
---

# Learning Go

Some content here.

## Setup

Install Go from the website.
`
	result, err := ParseMarkdown([]byte(input))
	if err != nil {
		t.Fatalf("ParseMarkdown error: %v", err)
	}
	if result.Title != "Learning Go" {
		t.Errorf("expected title 'Learning Go', got '%s'", result.Title)
	}
	if result.Date != "2026-07-04" {
		t.Errorf("expected date '2026-07-04', got '%s'", result.Date)
	}
	if len(result.Tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(result.Tags))
	}
	if result.Slug != "learning-go" {
		t.Errorf("expected slug 'learning-go', got '%s'", result.Slug)
	}
	if !strings.HasPrefix(result.HTML, "<") {
		t.Errorf("expected HTML output starting with '<', got: %s", result.HTML[:min(50, len(result.HTML))])
	}

	// TOC should contain entries for the headings.
	if result.TOCJSON == "" {
		t.Errorf("expected non-empty TOCJSON")
	}
	var tocItems []TOCItem
	if err := json.Unmarshal([]byte(result.TOCJSON), &tocItems); err != nil {
		t.Fatalf("failed to unmarshal TOCJSON: %v", err)
	}
	if len(tocItems) == 0 {
		t.Errorf("expected at least one TOC entry, got none")
	}
}

func TestParseMarkdown_NoFrontmatter(t *testing.T) {
	input := `# Hello World

This post has no frontmatter.
`
	result, err := ParseMarkdown([]byte(input))
	if err != nil {
		t.Fatalf("ParseMarkdown error: %v", err)
	}
	if result.Title != "Hello World" {
		t.Errorf("expected title 'Hello World' (from h1), got '%s'", result.Title)
	}
	if result.Date == "" {
		t.Errorf("expected a fallback date, got empty")
	}
	if len(result.Tags) != 0 {
		t.Errorf("expected 0 tags, got %d", len(result.Tags))
	}
	if result.Slug != "hello-world" {
		t.Errorf("expected slug 'hello-world', got '%s'", result.Slug)
	}
}

func TestParseMarkdown_NoH1NoFrontmatterTitle(t *testing.T) {
	input := `Some content without a heading.

## Subsection

Just content.
`
	result, err := ParseMarkdown([]byte(input))
	if err != nil {
		t.Fatalf("ParseMarkdown error: %v", err)
	}
	if result.Title != "Untitled" {
		t.Errorf("expected title 'Untitled', got '%s'", result.Title)
	}
}

func TestParseMarkdown_CommaSeparatedTags(t *testing.T) {
	input := `---
title: Tag Test
date: 2026-07-04
tags: Go, backend, testing
---

# Tag Test

Content.
`
	result, err := ParseMarkdown([]byte(input))
	if err != nil {
		t.Fatalf("ParseMarkdown error: %v", err)
	}
	if len(result.Tags) != 3 {
		t.Errorf("expected 3 tags, got %d: %v", len(result.Tags), result.Tags)
	}
}

func TestParseMarkdown_SlugSpecialChars(t *testing.T) {
	input := `---
title: "Learning Go: A Beginner's Guide!"
date: 2026-07-04
---

# Learning Go: A Beginner's Guide!

Content.
`
	result, err := ParseMarkdown([]byte(input))
	if err != nil {
		t.Fatalf("ParseMarkdown error: %v", err)
	}
	if result.Slug != "learning-go-a-beginners-guide" {
		t.Errorf("expected slug 'learning-go-a-beginners-guide', got '%s'", result.Slug)
	}
}

func TestParseMarkdown_TOCLevels(t *testing.T) {
	input := `---
title: TOC Test
date: 2026-07-04
---

# Section 1

## Subsection 1.1

### Sub-subsection 1.1.1

## Subsection 1.2

# Section 2

## Subsection 2.1
`
	result, err := ParseMarkdown([]byte(input))
	if err != nil {
		t.Fatalf("ParseMarkdown error: %v", err)
	}

	var tocItems []TOCItem
	if err := json.Unmarshal([]byte(result.TOCJSON), &tocItems); err != nil {
		t.Fatalf("failed to unmarshal TOCJSON: %v", err)
	}

	// Should have 2 top-level items (Section 1, Section 2)
	if len(tocItems) != 2 {
		t.Errorf("expected 2 top-level TOC items, got %d", len(tocItems))
	}

	if len(tocItems) > 0 {
		if tocItems[0].Level != 1 {
			t.Errorf("expected Level 1 for top-level item, got %d", tocItems[0].Level)
		}
		if tocItems[0].Text != "Section 1" {
			t.Errorf("expected 'Section 1', got '%s'", tocItems[0].Text)
		}
		// Section 1 should have 2 children
		if len(tocItems[0].Children) != 2 {
			t.Errorf("expected 2 children under Section 1, got %d", len(tocItems[0].Children))
		}
		// First child should be level 2
		if len(tocItems[0].Children) > 0 && tocItems[0].Children[0].Level != 2 {
			t.Errorf("expected Level 2 for subsection, got %d", tocItems[0].Children[0].Level)
		}
	}
}
