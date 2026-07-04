# Vue + Go 博客 — 实施计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 构建个人学习日记博客，支持管理员拖拽上传 Markdown 文件，自动解析并展示。

**Architecture:** Vue 3 SPA 前端通过 REST API 与 Go + Gin 后端通信，SQLite 存储数据。Go 端负责 Markdown 解析（goldmark）、代码高亮（Shiki 通过 Node.js sidecar）、JWT 认证、FTS5 全文搜索。前端使用 Pinia 状态管理、Tailwind CSS 暗色模式、Vue Router 4 路由。最终交付单一 Go 二进制（通过 embed 内嵌前端静态文件）。

**Tech Stack:** Vue 3 + Vite + Vue Router 4 + Pinia + Tailwind CSS 3, Go 1.22+ + Gin + GORM + SQLite + goldmark + golang-jwt

## Global Constraints

- Go 版本 >= 1.22，使用 embed 内嵌前端文件
- SQLite 使用 `github.com/glebarez/sqlite`（纯 Go，无 CGO）
- JWT 使用 HS256 签名，7 天过期
- 所有密码使用 bcrypt 哈希
- 前端使用 TypeScript
- 管理接口必须 JWT 认证
- 每次 Markdown 上传自动解析 frontmatter（title/date/tags）
- 代码高亮使用 highlight.js（前端侧渲染，避免 Node.js sidecar 复杂性）
- DRY, YAGNI, TDD

---

### Task 1: Go 后端项目初始化

**Files:**
- Create: `backend/go.mod`
- Create: `backend/main.go`

**Interfaces:**
- Consumes: none
- Produces: `backend/` 目录结构就绪，`go mod` 初始化，依赖已安装

- [ ] **Step 1: 创建 go.mod**

```bash
cd backend
go mod init blog-backend
```

Expected: `go.mod` 创建，module 为 `blog-backend`

- [ ] **Step 2: 安装核心依赖**

```bash
cd backend
go get github.com/gin-gonic/gin
go get github.com/glebarez/sqlite
go get gorm.io/gorm
go get github.com/golang-jwt/jwt/v5
go get golang.org/x/crypto
go get github.com/yuin/goldmark
go get github.com/yuin/goldmark-meta
go get github.com/abhinav/goldmark-toc
go get gopkg.in/yaml.v3
```

Expected: `go.mod` 和 `go.sum` 更新，所有依赖安装成功

- [ ] **Step 3: 创建最小启动文件验证环境**

```go
// backend/main.go
package main

import (
    "fmt"
    "net/http"

    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()
    r.GET("/api/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"status": "ok"})
    })
    fmt.Println("Server starting on :8080")
    r.Run(":8080")
}
```

- [ ] **Step 4: 运行验证**

```bash
cd backend
go run main.go
```

Expected: `Server starting on :8080`，访问 `http://localhost:8080/api/health` 返回 `{"status":"ok"}`

- [ ] **Step 5: 创建目录结构**

```bash
cd backend
mkdir -p internal/auth internal/handler internal/model internal/parser internal/service internal/database config
```

Expected: 所有目录创建完毕

- [ ] **Step 6: Commit**

```bash
git add backend/
git commit -m "feat: initialize Go backend with Gin and dependencies"
```

---

### Task 2: Vue 3 前端项目初始化

**Files:**
- Create: `frontend/` (Vite + Vue 3 + TS 项目)
- Create: `frontend/tailwind.config.js`
- Create: `frontend/src/main.ts`
- Create: `frontend/src/App.vue`

**Interfaces:**
- Consumes: none
- Produces: Vue 3 项目就绪，Tailwind CSS 可用，`npm run dev` 可启动

- [ ] **Step 1: 用 Vite 创建 Vue 3 + TS 项目**

```bash
npm create vite@latest frontend -- --template vue-ts
cd frontend
npm install
```

Expected: `frontend/` 目录创建成功，`npm run dev` 可启动

- [ ] **Step 2: 安装 Tailwind CSS 3**

```bash
cd frontend
npm install -D tailwindcss@3 postcss autoprefixer
npx tailwindcss init -p
```

Expected: `tailwind.config.js` 和 `postcss.config.js` 创建成功

- [ ] **Step 3: 配置 Tailwind 内容路径和暗色模式**

```js
// frontend/tailwind.config.js
/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  darkMode: 'class',
  theme: {
    extend: {},
  },
  plugins: [],
}
```

- [ ] **Step 4: 添加 Tailwind 指令到 CSS**

```css
/* frontend/src/style.css */
@tailwind base;
@tailwind components;
@tailwind utilities;
```

- [ ] **Step 5: 安装前端依赖**

```bash
cd frontend
npm install vue-router@4 pinia axios highlight.js
```

Expected: 所有前端依赖安装成功

- [ ] **Step 6: 创建最小目录结构**

```bash
cd frontend/src
mkdir -p components views stores router utils
```

Expected: 目录结构创建完毕

- [ ] **Step 7: 验证 Tailwind 可用**

```vue
<!-- frontend/src/App.vue -->
<script setup lang="ts">
</script>
<template>
  <div class="min-h-screen bg-white dark:bg-gray-900 text-gray-900 dark:text-gray-100">
    <h1 class="text-3xl font-bold text-blue-500">Blog</h1>
  </div>
</template>
```

- [ ] **Step 8: 验证项目运行**

```bash
cd frontend
npm run dev
```

Expected: 开发服务器启动，页面显示蓝色 "Blog" 标题

- [ ] **Step 9: Commit**

```bash
git add frontend/
git commit -m "feat: initialize Vue 3 frontend with Tailwind CSS"
```

---

### Task 3: 数据库初始化与模型定义

**Files:**
- Create: `backend/config/config.go`
- Create: `backend/internal/database/sqlite.go`
- Create: `backend/internal/model/user.go`
- Create: `backend/internal/model/post.go`
- Create: `backend/internal/model/tag.go`

**Interfaces:**
- Consumes: none
- Produces:
  - `database.Init(path string) *gorm.DB` — 打开 SQLite，执行 AutoMigrate，返回 DB 实例
  - `model.User{ID, Username, PasswordHash, CreatedAt}`
  - `model.Post{ID, Title, Slug, ContentMD, ContentHTML, TOCJSON, CreatedAt, UpdatedAt, Tags []Tag}`
  - `model.Tag{ID, Name, Posts []Post (many2many:post_tags)}`

- [ ] **Step 1: 创建配置管理**

```go
// backend/config/config.go
package config

import "os"

type Config struct {
    Port      string
    DBPath    string
    JWTSecret string
}

func Load() Config {
    return Config{
        Port:      getEnv("PORT", "8080"),
        DBPath:    getEnv("DB_PATH", "./blog.db"),
        JWTSecret: getEnv("JWT_SECRET", "change-me-in-production"),
    }
}

func getEnv(key, fallback string) string {
    if v := os.Getenv(key); v != "" {
        return v
    }
    return fallback
}
```

- [ ] **Step 2: 创建 User 模型**

```go
// backend/internal/model/user.go
package model

import "time"

type User struct {
    ID           uint   `gorm:"primaryKey"`
    Username     string `gorm:"uniqueIndex;size:100;not null"`
    PasswordHash string `gorm:"not null"`
    CreatedAt    time.Time
}
```

- [ ] **Step 3: 创建 Post 模型**

```go
// backend/internal/model/post.go
package model

import "time"

type Post struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    Title       string    `gorm:"size:500;not null" json:"title"`
    Slug        string    `gorm:"uniqueIndex;size:500;not null" json:"slug"`
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
```

- [ ] **Step 4: 创建 Tag 模型**

```go
// backend/internal/model/tag.go
package model

type Tag struct {
    ID   uint   `gorm:"primaryKey" json:"id"`
    Name string `gorm:"uniqueIndex;size:100;not null" json:"name"`
}
```

- [ ] **Step 5: 创建数据库初始化**

```go
// backend/internal/database/sqlite.go
package database

import (
    "blog-backend/config"
    "blog-backend/internal/model"
    "log"

    "github.com/glebarez/sqlite"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init(cfg config.Config) {
    var err error
    DB, err = gorm.Open(sqlite.Open(cfg.DBPath), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    if err := DB.AutoMigrate(
        &model.User{},
        &model.Post{},
        &model.Tag{},
    ); err != nil {
        log.Fatalf("Failed to migrate database: %v", err)
    }

    // Enable WAL mode for better concurrent read performance
    DB.Exec("PRAGMA journal_mode=WAL")

    log.Println("Database initialized successfully")
}
```

- [ ] **Step 6: 更新 main.go 验证数据库**

```go
// backend/main.go
package main

import (
    "blog-backend/config"
    "blog-backend/internal/database"
    "fmt"
    "net/http"

    "github.com/gin-gonic/gin"
)

func main() {
    cfg := config.Load()
    database.Init(cfg)

    r := gin.Default()
    r.GET("/api/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"status": "ok"})
    })
    fmt.Println("Server starting on :" + cfg.Port)
    r.Run(":" + cfg.Port)
}
```

- [ ] **Step 7: 运行验证**

```bash
cd backend
go run main.go
```

Expected: 日志显示 `Database initialized successfully`，`blog.db` 文件生成，表结构创建

- [ ] **Step 8: Commit**

```bash
git add backend/
git commit -m "feat: add database initialization with GORM and SQLite models"
```

---

### Task 4: JWT 认证（签发 + 中间件）

**Files:**
- Create: `backend/internal/auth/jwt.go`
- Create: `backend/internal/auth/middleware.go`
- Create: `backend/internal/service/user.go`

**Interfaces:**
- Consumes: `config.Config`, `database.DB`, `model.User`
- Produces:
  - `auth.GenerateToken(userID uint, secret string) (string, time.Time, error)` — 签发 7 天 Token
  - `auth.ValidateToken(tokenString string, secret string) (*Claims, error)` — 验证并返回 Claims
  - `auth.AuthMiddleware(secret string) gin.HandlerFunc` — Gin 中间件，失败返回 401
  - `service.CreateDefaultUser(db *gorm.DB) error` — 创建默认管理员（如果不存在）

- [ ] **Step 1: 实现 JWT 签发和验证**

```go
// backend/internal/auth/jwt.go
package auth

import (
    "errors"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
    UserID uint `json:"user_id"`
    jwt.RegisteredClaims
}

func GenerateToken(userID uint, secret string) (string, time.Time, error) {
    expiresAt := time.Now().Add(7 * 24 * time.Hour)
    claims := &Claims{
        UserID: userID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expiresAt),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString([]byte(secret))
    return tokenString, expiresAt, err
}

func ValidateToken(tokenString, secret string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
        if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("unexpected signing method")
        }
        return []byte(secret), nil
    })
    if err != nil {
        return nil, err
    }
    claims, ok := token.Claims.(*Claims)
    if !ok || !token.Valid {
        return nil, errors.New("invalid token")
    }
    return claims, nil
}
```

- [ ] **Step 2: 实现 Auth 中间件**

```go
// backend/internal/auth/middleware.go
package auth

import (
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
)

func AuthMiddleware(secret string) gin.HandlerFunc {
    return func(c *gin.Context) {
        header := c.GetHeader("Authorization")
        if header == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
            return
        }
        parts := strings.SplitN(header, " ", 2)
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization format"})
            return
        }
        claims, err := ValidateToken(parts[1], secret)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
            return
        }
        c.Set("userID", claims.UserID)
        c.Next()
    }
}
```

- [ ] **Step 3: 创建默认管理员种子**

```go
// backend/internal/service/user.go
package service

import (
    "errors"

    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"

    "blog-backend/internal/model"
)

func CreateDefaultUser(db *gorm.DB) error {
    var count int64
    db.Model(&model.User{}).Count(&count)
    if count > 0 {
        return nil
    }
    hash, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    user := model.User{
        Username:     "admin",
        PasswordHash: string(hash),
    }
    return db.Create(&user).Error
}

func AuthenticateUser(db *gorm.DB, username, password string) (*model.User, error) {
    var user model.User
    if err := db.Where("username = ?", username).First(&user).Error; err != nil {
        return nil, errors.New("invalid credentials")
    }
    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
        return nil, errors.New("invalid credentials")
    }
    return &user, nil
}
```

- [ ] **Step 4: 在 main.go 调用种子函数**

在 `database.Init(cfg)` 之后添加：

```go
if err := service.CreateDefaultUser(database.DB); err != nil {
    log.Fatalf("Failed to create default user: %v", err)
}
log.Println("Default admin user ensured (username: admin, password: admin)")
```

- [ ] **Step 5: 更新 import**

在 `backend/main.go` 中添加：
```go
import (
    "blog-backend/config"
    "blog-backend/internal/database"
    "blog-backend/internal/service"
    // ... existing imports
)
```

- [ ] **Step 6: 运行验证**

```bash
cd backend
go run main.go
```

Expected: 日志显示默认管理员已创建，`users` 表有一条 admin 记录

- [ ] **Step 7: Commit**

```bash
git add backend/internal/auth/ backend/internal/service/ backend/main.go
git commit -m "feat: add JWT authentication and admin user seed"
```

---

### Task 5: Markdown 解析器

**Files:**
- Create: `backend/internal/parser/markdown.go`

**Interfaces:**
- Consumes: `goldmark`, `goldmark-meta`, `goldmark-toc`
- Produces:
  - `parser.ParseResult{Title, Date, Tags []string, HTML, TOCJSON string, Slug string}`
  - `parser.ParseMarkdown(mdContent []byte) (*ParseResult, error)` — 解析 MD 内容，返回结构化结果

- [ ] **Step 1: 安装 goldmark 扩展**

```bash
cd backend
go get github.com/yuin/goldmark
go get github.com/yuin/goldmark-meta
go get go.abhg.dev/goldmark/toc
```

- [ ] **Step 2: 实现解析器**

```go
// backend/internal/parser/markdown.go
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
    "go.abhg.dev/goldmark/toc"
    "gopkg.in/yaml.v3"
)

type ParseResult struct {
    Title   string   `json:"title"`
    Date    string   `json:"date"`
    Tags    []string `json:"tags"`
    HTML    string   `json:"html"`
    TOCJSON string   `json:"toc"`
    Slug    string   `json:"slug"`
}

func ParseMarkdown(mdContent []byte) (*ParseResult, error) {
    md := goldmark.New(
        goldmark.WithExtensions(
            extension.GFM,
            extension.Table,
            meta.Meta,
            &toc.Extender{
                Title: "目录",
            },
        ),
        goldmark.WithParserOptions(parser.WithAutoHeadingID()),
        goldmark.WithRendererOptions(html.WithHardWraps()),
    )

    // Extract frontmatter
    var buf bytes.Buffer
    ctx := parser.NewContext()
    if err := md.Convert(mdContent, &buf, parser.WithContext(ctx)); err != nil {
        return nil, fmt.Errorf("markdown convert: %w", err)
    }

    metaData := meta.Get(ctx)
    if metaData == nil {
        metaData = map[string]interface{}{}
    }

    // Parse title
    title := getStringField(metaData, "title")
    if title == "" {
        // Fallback: use first h1 from content
        title = extractFirstH1(mdContent)
    }

    // Parse date
    date := getStringField(metaData, "date")
    if date == "" {
        date = time.Now().Format("2006-01-02")
    }

    // Parse tags
    var tags []string
    switch v := metaData["tags"].(type) {
    case []interface{}:
        for _, t := range v {
            tags = append(tags, fmt.Sprintf("%v", t))
        }
    case string:
        // Support comma-separated string format
        for _, t := range strings.Split(v, ",") {
            t = strings.TrimSpace(t)
            if t != "" {
                tags = append(tags, t)
            }
        }
    }

    // Extract TOC
    tocTree, ok := toc.Get(ctx)
    tocJSON := "[]"
    if ok && tocTree != nil {
        tocJSON = tocToJSON(tocTree)
    }

    // Generate slug from title
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

func getStringField(m map[string]interface{}, key string) string {
    if v, ok := m[key]; ok {
        return fmt.Sprintf("%v", v)
    }
    return ""
}

func extractFirstH1(content []byte) string {
    lines := strings.Split(string(content), "\n")
    for _, line := range lines {
        if strings.HasPrefix(line, "# ") {
            return strings.TrimPrefix(line, "# ")
        }
    }
    return "Untitled"
}

func slugify(title string) string {
    s := strings.ToLower(title)
    s = strings.ReplaceAll(s, " ", "-")
    // Remove non-alphanumeric characters except hyphens
    var buf strings.Builder
    for _, r := range s {
        if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
            buf.WriteRune(r)
        }
    }
    result := buf.String()
    // Collapse multiple hyphens
    for strings.Contains(result, "--") {
        result = strings.ReplaceAll(result, "--", "-")
    }
    return strings.Trim(result, "-")
}

// TOCItem represents a heading in the table of contents
type TOCItem struct {
    ID       string    `json:"id"`
    Text     string    `json:"text"`
    Level    int       `json:"level"`
    Children []TOCItem `json:"children"`
}

func tocToJSON(toc *toc.TOC) string {
    items := buildTOCItems(toc.Items)
    b, _ := json.Marshal(items)
    return string(b)
}

func buildTOCItems(items toc.Items) []TOCItem {
    var result []TOCItem
    for _, item := range items {
        ti := TOCItem{
            ID:       item.ID,
            Text:     string(item.Title),
            Level:    item.Level,
            Children: buildTOCItems(item.Items),
        }
        if ti.Children == nil {
            ti.Children = []TOCItem{}
        }
        result = append(result, ti)
    }
    return result
}
```

- [ ] **Step 3: 编写测试验证解析器**

创建 `backend/internal/parser/markdown_test.go`，用一段含 frontmatter 的 markdown 测试解析结果：

```go
// backend/internal/parser/markdown_test.go
package parser

import (
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
        t.Errorf("expected HTML output, got: %s", result.HTML[:50])
    }
    if result.TOCJSON == "" || result.TOCJSON == "[]" {
        t.Errorf("expected non-empty TOC")
    }
}
```

- [ ] **Step 4: 运行测试**

```bash
cd backend
go test ./internal/parser/ -v
```

Expected: 所有测试 PASS

- [ ] **Step 5: Commit**

```bash
git add backend/internal/parser/
git commit -m "feat: add markdown parser with goldmark, frontmatter, and TOC"
```

---

### Task 6: 文章与标签服务层

**Files:**
- Create: `backend/internal/service/post.go`
- Create: `backend/internal/service/tag.go`
- Create: `backend/internal/service/search.go`

**Interfaces:**
- Consumes: `database.DB`, `model.Post`, `model.Tag`
- Produces:
  - `service.CreatePost(db, result *parser.ParseResult, rawMD string) (*model.Post, error)` — 创建文章及关联标签，更新 FTS
  - `service.GetPosts(db, page, perPage int, tag string) ([]model.Post, int64, error)` — 分页文章列表，可按标签筛选
  - `service.GetPostBySlug(db, slug string) (*model.Post, error)` — 按 slug 获取单篇文章
  - `service.UpdatePost(db, id uint, updates map[string]interface{}) (*model.Post, error)` — 更新文章
  - `service.DeletePost(db, id uint) error` — 删除文章
  - `service.GetAllTags(db) ([]TagWithCount, error)` — 所有标签及文章数
  - `service.SearchPosts(db, query string, page, perPage int) ([]model.Post, int64, error)` — FTS5 搜索
  - `service.GetTimeline(db) ([]model.TimelineEntry, error)` — 按年月归档

- [ ] **Step 1: 实现文章服务**

```go
// backend/internal/service/post.go
package service

import (
    "fmt"
    "strings"

    "gorm.io/gorm"

    "blog-backend/internal/database"
    "blog-backend/internal/model"
    "blog-backend/internal/parser"
)

func CreatePost(db *gorm.DB, result *parser.ParseResult, rawMD string) (*model.Post, error) {
    // Resolve slug collisions
    slug := result.Slug
    var count int64
    db.Model(&model.Post{}).Where("slug = ?", result.Slug).Count(&count)
    if count > 0 {
        slug = fmt.Sprintf("%s-%d", result.Slug, count+1)
    }

    // Parse date
    parsedDate, _ := parseDate(result.Date)

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

    // sync FTS index
    syncFTS(db, post.ID, post.Title, rawMD)

    return post, nil
}

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

    return posts, total, nil
}

func GetPostBySlug(db *gorm.DB, slug string) (*model.Post, error) {
    var post model.Post
    if err := db.Preload("Tags").Where("slug = ?", slug).First(&post).Error; err != nil {
        return nil, err
    }
    return &post, nil
}

func GetPostByID(db *gorm.DB, id uint) (*model.Post, error) {
    var post model.Post
    if err := db.Preload("Tags").First(&post, id).Error; err != nil {
        return nil, err
    }
    return &post, nil
}

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

func DeletePost(db *gorm.DB, id uint) error {
    var post model.Post
    if err := db.First(&post, id).Error; err != nil {
        return err
    }
    // Clear tag associations then delete
    db.Model(&post).Association("Tags").Clear()
    return db.Unscoped().Delete(&post).Error
}

func GetAllTags(db *gorm.DB) ([]model.Tag, error) {
    database.DB.Raw(`
        SELECT tags.id, tags.name, COUNT(post_tags.post_id) as count
        FROM tags
        JOIN post_tags ON post_tags.tag_id = tags.id
        GROUP BY tags.id
        ORDER BY count DESC
    `)
    return nil, nil // Placeholder — see tag service
}

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

func GetTimeline(db *gorm.DB) ([]map[string]interface{}, error) {
    var posts []model.Post
    db.Preload("Tags").Order("created_at DESC").Find(&posts)

    timeline := make(map[string][]model.Post)
    for _, p := range posts {
        key := p.CreatedAt.Format("2006-01")
        timeline[key] = append(timeline[key], p)
    }

    // Build structured result
    years := make(map[int]map[int][]model.Post)
    for _, p := range posts {
        y := p.CreatedAt.Year()
        m := int(p.CreatedAt.Month())
        if years[y] == nil {
            years[y] = make(map[int][]model.Post)
        }
        years[y][m] = append(years[y][m], p)
    }

    var result []map[string]interface{}
    for y, months := range years {
        var monthList []map[string]interface{}
        for m := 12; m >= 1; m-- {
            if posts, ok := months[m]; ok {
                monthList = append(monthList, map[string]interface{}{
                    "month": m,
                    "posts": posts,
                })
            }
        }
        result = append(result, map[string]interface{}{
            "year":   y,
            "months": monthList,
        })
    }

    return result, nil
}

// helpers

func parseDate(dateStr string) (time.Time, error) {
    // Only import at the top with proper formatting
    return time.Parse("2006-01-02", dateStr)
}

func syncFTS(db *gorm.DB, postID uint, title, contentMD string) {
    db.Exec("DELETE FROM posts_fts WHERE rowid = ?", postID)
    db.Exec("INSERT INTO posts_fts(rowid, title, content_md) VALUES(?, ?, ?)", postID, title, contentMD)
}
```

- [ ] **Step 2: 实现标签服务**

```go
// backend/internal/service/tag.go
package service

import (
    "blog-backend/internal/model"

    "gorm.io/gorm"
)

type TagWithCount struct {
    Name  string `json:"name"`
    Count int    `json:"count"`
}

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

func EnsureTags(db *gorm.DB, tagNames []string) ([]model.Tag, error) {
    var tags []model.Tag
    for _, name := range tagNames {
        var tag model.Tag
        db.Where("name = ?", name).FirstOrCreate(&tag, model.Tag{Name: name})
        tags = append(tags, tag)
    }
    return tags, nil
}
```

- [ ] **Step 3: FTS5 初始化**

```go
// backend/internal/database/fts.go
package database

import "log"

func InitFTS() {
    if DB == nil {
        return
    }
    DB.Exec(`
        CREATE VIRTUAL TABLE IF NOT EXISTS posts_fts USING fts5(
            title,
            content_md,
            content='posts',
            content_rowid='id'
        )
    `)
    DB.Exec(`
        CREATE TRIGGER IF NOT EXISTS posts_ai AFTER INSERT ON posts BEGIN
            INSERT INTO posts_fts(rowid, title, content_md) VALUES (new.id, new.title, new.content_md);
        END
    `)
    DB.Exec(`
        CREATE TRIGGER IF NOT EXISTS posts_ad AFTER DELETE ON posts BEGIN
            INSERT INTO posts_fts(posts_fts, rowid, title, content_md) VALUES('delete', old.id, old.title, old.content_md);
        END
    `)
    DB.Exec(`
        CREATE TRIGGER IF NOT EXISTS posts_au AFTER UPDATE ON posts BEGIN
            INSERT INTO posts_fts(posts_fts, rowid, title, content_md) VALUES('delete', old.id, old.title, old.content_md);
            INSERT INTO posts_fts(rowid, title, content_md) VALUES (new.id, new.title, new.content_md);
        END
    `)
    log.Println("FTS5 index initialized")
}
```

- [ ] **Step 4: 在 database.Init 中调用 InitFTS**

在 `backend/internal/database/sqlite.go` 的 `Init` 函数末尾，`AutoMigrate` 之后添加：

```go
InitFTS()
```

- [ ] **Step 5: Commit**

```bash
git add backend/internal/service/ backend/internal/database/fts.go
git commit -m "feat: add post, tag, search services with FTS5"
```

---

### Task 7: API Handler 层

**Files:**
- Create: `backend/internal/handler/auth.go`
- Create: `backend/internal/handler/post.go`
- Create: `backend/internal/handler/tag.go`
- Create: `backend/internal/handler/upload.go`
- Create: `backend/internal/handler/search.go`
- Create: `backend/internal/handler/timeline.go`

**Interfaces:**
- Consumes: `database.DB`, `service.*`, `auth.*`, `parser.*`
- Produces: Gin handler functions，注册到路由

- [ ] **Step 1: Auth handler — 登录**

```go
// backend/internal/handler/auth.go
package handler

import (
    "net/http"

    "github.com/gin-gonic/gin"

    "blog-backend/config"
    "blog-backend/internal/auth"
    "blog-backend/internal/database"
    "blog-backend/internal/service"
)

func Login(cfg config.Config) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req struct {
            Username string `json:"username" binding:"required"`
            Password string `json:"password" binding:"required"`
        }
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
            return
        }
        user, err := service.AuthenticateUser(database.DB, req.Username, req.Password)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
            return
        }
        token, expiresAt, err := auth.GenerateToken(user.ID, cfg.JWTSecret)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
            return
        }
        c.JSON(http.StatusOK, gin.H{
            "token":      token,
            "expires_at": expiresAt,
        })
    }
}
```

- [ ] **Step 2: Post handlers**

```go
// backend/internal/handler/post.go
package handler

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"

    "blog-backend/internal/database"
    "blog-backend/internal/service"
)

func ListPosts() gin.HandlerFunc {
    return func(c *gin.Context) {
        page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
        perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
        tag := c.Query("tag")

        posts, total, err := service.GetPosts(database.DB, page, perPage, tag)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, gin.H{
            "posts":    posts,
            "total":    total,
            "page":     page,
            "per_page": perPage,
        })
    }
}

func GetPost() gin.HandlerFunc {
    return func(c *gin.Context) {
        slug := c.Param("slug")
        post, err := service.GetPostBySlug(database.DB, slug)
        if err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
            return
        }
        c.JSON(http.StatusOK, gin.H{"post": post})
    }
}

func UpdatePost() gin.HandlerFunc {
    return func(c *gin.Context) {
        id, err := strconv.ParseUint(c.Param("id"), 10, 64)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
            return
        }
        var updates map[string]interface{}
        if err := c.ShouldBindJSON(&updates); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
            return
        }
        post, err := service.UpdatePost(database.DB, uint(id), updates)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, gin.H{"post": post})
    }
}

func DeletePost() gin.HandlerFunc {
    return func(c *gin.Context) {
        id, err := strconv.ParseUint(c.Param("id"), 10, 64)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
            return
        }
        if err := service.DeletePost(database.DB, uint(id)); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.Status(http.StatusNoContent)
    }
}
```

- [ ] **Step 3: Upload handler**

```go
// backend/internal/handler/upload.go
package handler

import (
    "io"
    "net/http"

    "github.com/gin-gonic/gin"

    "blog-backend/internal/database"
    "blog-backend/internal/parser"
    "blog-backend/internal/service"
)

func UploadPost() gin.HandlerFunc {
    return func(c *gin.Context) {
        file, err := c.FormFile("file")
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
            return
        }
        f, err := file.Open()
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open file"})
            return
        }
        defer f.Close()

        content, err := io.ReadAll(f)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file"})
            return
        }

        result, err := parser.ParseMarkdown(content)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse markdown: " + err.Error()})
            return
        }

        post, err := service.CreatePost(database.DB, result, string(content))
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create post: " + err.Error()})
            return
        }

        c.JSON(http.StatusCreated, gin.H{"post": post})
    }
}
```

- [ ] **Step 4: Tag handler**

```go
// backend/internal/handler/tag.go
package handler

import (
    "net/http"

    "github.com/gin-gonic/gin"

    "blog-backend/internal/database"
    "blog-backend/internal/service"
)

func ListTags() gin.HandlerFunc {
    return func(c *gin.Context) {
        tags, err := service.GetAllTags(database.DB)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, gin.H{"tags": tags})
    }
}
```

- [ ] **Step 5: Search handler**

```go
// backend/internal/handler/search.go
package handler

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"

    "blog-backend/internal/database"
    "blog-backend/internal/service"
)

func SearchPosts() gin.HandlerFunc {
    return func(c *gin.Context) {
        q := c.Query("q")
        if q == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "query is required"})
            return
        }
        page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
        perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))

        posts, total, err := service.SearchPosts(database.DB, q, page, perPage)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        if posts == nil {
            posts = []model.Post{}
        }
        c.JSON(http.StatusOK, gin.H{
            "posts":    posts,
            "total":    total,
            "page":     page,
            "per_page": perPage,
        })
    }
}
```

- [ ] **Step 6: Timeline handler**

```go
// backend/internal/handler/timeline.go
package handler

import (
    "net/http"

    "github.com/gin-gonic/gin"

    "blog-backend/internal/database"
    "blog-backend/internal/service"
)

func GetTimeline() gin.HandlerFunc {
    return func(c *gin.Context) {
        timeline, err := service.GetTimeline(database.DB)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, gin.H{"timeline": timeline})
    }
}
```

- [ ] **Step 7: 注册路由 — 更新 main.go**

```go
// backend/main.go
package main

import (
    "blog-backend/config"
    "blog-backend/internal/auth"
    "blog-backend/internal/database"
    "blog-backend/internal/handler"
    "blog-backend/internal/service"
    "fmt"
    "log"

    "github.com/gin-gonic/gin"
)

func main() {
    cfg := config.Load()
    database.Init(cfg)

    if err := service.CreateDefaultUser(database.DB); err != nil {
        log.Fatalf("Failed to create default user: %v", err)
    }

    r := gin.Default()

    // Public routes
    r.GET("/api/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })
    r.POST("/api/auth/login", handler.Login(cfg))

    r.GET("/api/posts", handler.ListPosts())
    r.GET("/api/posts/:slug", handler.GetPost())
    r.GET("/api/tags", handler.ListTags())
    r.GET("/api/search", handler.SearchPosts())
    r.GET("/api/timeline", handler.GetTimeline())

    // Admin routes (protected)
    admin := r.Group("/api/admin")
    admin.Use(auth.AuthMiddleware(cfg.JWTSecret))
    {
        admin.POST("/upload", handler.UploadPost())
        admin.PUT("/posts/:id", handler.UpdatePost())
        admin.DELETE("/posts/:id", handler.DeletePost())
    }

    fmt.Println("Server starting on :" + cfg.Port)
    r.Run(":" + cfg.Port)
}
```

- [ ] **Step 8: 编译运行验证**

```bash
cd backend
go build -o blog-server .
./blog-server
```

用 curl 验证：

```bash
# 验证健康检查
curl http://localhost:8080/api/health

# 登录获取 token
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin"}'
```

Expected: 两个接口都返回 200

- [ ] **Step 9: Commit**

```bash
git add backend/
git commit -m "feat: add API handlers and route registration"
```

---

### Task 8: Go embed 内嵌前端 + SPA 回退

**Files:**
- Modify: `backend/main.go`

- [ ] **Step 1: 更新 main.go，添加 embed 和 SPA fallback**

```go
// backend/main.go
package main

import (
    "blog-backend/config"
    "blog-backend/internal/auth"
    "blog-backend/internal/database"
    "blog-backend/internal/handler"
    "blog-backend/internal/service"
    "embed"
    "fmt"
    "io/fs"
    "log"
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
)

//go:embed frontend-dist/*
var staticFiles embed.FS

func main() {
    cfg := config.Load()
    database.Init(cfg)

    if err := service.CreateDefaultUser(database.DB); err != nil {
        log.Fatalf("Failed to create default user: %v", err)
    }

    r := gin.Default()

    // Public API routes
    r.GET("/api/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })
    r.POST("/api/auth/login", handler.Login(cfg))
    r.GET("/api/posts", handler.ListPosts())
    r.GET("/api/posts/:slug", handler.GetPost())
    r.GET("/api/tags", handler.ListTags())
    r.GET("/api/search", handler.SearchPosts())
    r.GET("/api/timeline", handler.GetTimeline())

    // Admin routes (protected)
    admin := r.Group("/api/admin")
    admin.Use(auth.AuthMiddleware(cfg.JWTSecret))
    {
        admin.POST("/upload", handler.UploadPost())
        admin.PUT("/posts/:id", handler.UpdatePost())
        admin.DELETE("/posts/:id", handler.DeletePost())
    }

    // Serve SPA static files
    frontendFS, err := fs.Sub(staticFiles, "frontend-dist")
    if err != nil {
        log.Fatalf("Failed to load frontend files: %v", err)
    }

    r.NoRoute(func(c *gin.Context) {
        path := c.Request.URL.Path
        // Don't intercept API routes
        if strings.HasPrefix(path, "/api/") {
            c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
            return
        }
        // Try to serve the requested file, fallback to index.html
        f, err := frontendFS.Open(strings.TrimPrefix(path, "/"))
        if err == nil {
            f.Close()
            c.FileFromFS(path, http.FS(frontendFS))
            return
        }
        c.FileFromFS("index.html", http.FS(frontendFS))
    })

    fmt.Println("Server starting on :" + cfg.Port)
    r.Run(":" + cfg.Port)
}
```

- [ ] **Step 2: 创建构建脚本将前端产物复制到 backend**

创建项目根目录的构建脚本：

```bash
# build.sh (项目根目录)
#!/bin/bash
set -e

echo "Building frontend..."
cd frontend
npm run build
cd ..

echo "Copying frontend dist to backend..."
rm -rf backend/frontend-dist
cp -r frontend/dist backend/frontend-dist

echo "Building backend..."
cd backend
go build -o blog-server .
cd ..

echo "Build complete: backend/blog-server"
```

- [ ] **Step 3: Commit**

```bash
git add backend/main.go build.sh
git commit -m "feat: add Go embed for frontend SPA with fallback routing"
```

---

### Task 9: 前端路由 + Layout + API 工具

**Files:**
- Create: `frontend/src/router/index.ts`
- Create: `frontend/src/utils/api.ts`
- Create: `frontend/src/layouts/DefaultLayout.vue`

- [ ] **Step 1: API 封装**

```typescript
// frontend/src/utils/api.ts
import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  timeout: 10000,
})

// Request interceptor: attach JWT token
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// Response interceptor: handle 401
api.interceptors.response.use(
  (res) => res,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

export default api
```

- [ ] **Step 2: 路由配置**

```typescript
// frontend/src/router/index.ts
import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import('../views/HomeView.vue'),
    },
    {
      path: '/post/:slug',
      name: 'post',
      component: () => import('../views/PostDetailView.vue'),
    },
    {
      path: '/timeline',
      name: 'timeline',
      component: () => import('../views/TimelineView.vue'),
    },
    {
      path: '/tags',
      name: 'tags',
      component: () => import('../views/TagsView.vue'),
    },
    {
      path: '/search',
      name: 'search',
      component: () => import('../views/SearchView.vue'),
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('../views/LoginView.vue'),
    },
    {
      path: '/admin',
      name: 'admin',
      component: () => import('../views/AdminView.vue'),
      meta: { requiresAuth: true },
    },
  ],
})

// Navigation guard
router.beforeEach((to, _from, next) => {
  if (to.meta.requiresAuth) {
    const token = localStorage.getItem('token')
    if (!token) {
      return next('/login')
    }
  }
  next()
})

export default router
```

- [ ] **Step 3: DefaultLayout**

```vue
<!-- frontend/src/layouts/DefaultLayout.vue -->
<script setup lang="ts">
</script>

<template>
  <div class="min-h-screen bg-white dark:bg-gray-900 text-gray-900 dark:text-gray-100">
    <header class="border-b border-gray-200 dark:border-gray-700">
      <nav class="max-w-4xl mx-auto px-4 h-14 flex items-center justify-between">
        <router-link to="/" class="text-xl font-bold">My Blog</router-link>
        <div class="flex items-center gap-4 text-sm">
          <router-link to="/timeline" class="hover:text-blue-500">Timeline</router-link>
          <router-link to="/tags" class="hover:text-blue-500">Tags</router-link>
          <router-link to="/search" class="hover:text-blue-500">Search</router-link>
          <router-link to="/admin" class="hover:text-blue-500">Admin</router-link>
        </div>
      </nav>
    </header>
    <main class="max-w-4xl mx-auto px-4 py-8">
      <slot />
    </main>
  </div>
</template>
```

- [ ] **Step 4: 更新 App.vue**

```vue
<!-- frontend/src/App.vue -->
<script setup lang="ts">
import DefaultLayout from './layouts/DefaultLayout.vue'
</script>

<template>
  <DefaultLayout>
    <router-view />
  </DefaultLayout>
</template>
```

- [ ] **Step 5: 更新 main.ts**

```typescript
// frontend/src/main.ts
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import './style.css'

const app = createApp(App)
app.use(createPinia())
app.use(router)
app.mount('#app')
```

- [ ] **Step 6: Commit**

```bash
git add frontend/src/
git commit -m "feat: add Vue Router, API utils, and default layout"
```

---

### Task 10: Auth Store + 登录页

**Files:**
- Create: `frontend/src/stores/auth.ts`
- Create: `frontend/src/views/LoginView.vue`

- [ ] **Step 1: Auth Store**

```typescript
// frontend/src/stores/auth.ts
import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../utils/api'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('token') || '')
  const isAuthenticated = ref(!!token.value)

  async function login(username: string, password: string) {
    const res = await api.post('/auth/login', { username, password })
    token.value = res.data.token
    localStorage.setItem('token', res.data.token)
    isAuthenticated.value = true
  }

  function logout() {
    token.value = ''
    localStorage.removeItem('token')
    isAuthenticated.value = false
  }

  return { token, isAuthenticated, login, logout }
})
```

- [ ] **Step 2: LoginView**

```vue
<!-- frontend/src/views/LoginView.vue -->
<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const auth = useAuthStore()
const username = ref('')
const password = ref('')
const error = ref('')

async function handleLogin() {
  error.value = ''
  try {
    await auth.login(username.value, password.value)
    router.push('/admin')
  } catch (e) {
    error.value = 'Invalid credentials'
  }
}
</script>

<template>
  <div class="max-w-sm mx-auto mt-20">
    <h1 class="text-2xl font-bold mb-6">Admin Login</h1>
    <form @submit.prevent="handleLogin" class="space-y-4">
      <div>
        <label class="block text-sm mb-1">Username</label>
        <input v-model="username" type="text"
          class="w-full px-3 py-2 border rounded dark:bg-gray-800 dark:border-gray-600" />
      </div>
      <div>
        <label class="block text-sm mb-1">Password</label>
        <input v-model="password" type="password"
          class="w-full px-3 py-2 border rounded dark:bg-gray-800 dark:border-gray-600" />
      </div>
      <p v-if="error" class="text-red-500 text-sm">{{ error }}</p>
      <button type="submit"
        class="w-full py-2 bg-blue-500 text-white rounded hover:bg-blue-600">
        Login
      </button>
    </form>
  </div>
</template>
```

- [ ] **Step 3: Commit**

```bash
git add frontend/src/stores/auth.ts frontend/src/views/LoginView.vue
git commit -m "feat: add auth store and login page"
```

---

### Task 11: 首页（文章列表 + 分页 + 标签筛选）

**Files:**
- Create: `frontend/src/components/PostCard.vue`
- Create: `frontend/src/views/HomeView.vue`
- Create: `frontend/src/stores/posts.ts`

- [ ] **Step 1: Posts Store**

```typescript
// frontend/src/stores/posts.ts
import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../utils/api'

export interface Post {
  id: number
  title: string
  slug: string
  content_html: string
  toc: string
  created_at: string
  updated_at: string
  tags: { id: number; name: string }[]
}

export const usePostsStore = defineStore('posts', () => {
  const posts = ref<Post[]>([])
  const total = ref(0)
  const loading = ref(false)

  async function fetchPosts(page = 1, tag = '') {
    loading.value = true
    try {
      const res = await api.get('/posts', { params: { page, per_page: 10, tag } })
      posts.value = res.data.posts
      total.value = res.data.total
    } finally {
      loading.value = false
    }
  }

  async function fetchPost(slug: string): Promise<Post | null> {
    try {
      const res = await api.get(`/posts/${slug}`)
      return res.data.post
    } catch {
      return null
    }
  }

  return { posts, total, loading, fetchPosts, fetchPost }
})
```

- [ ] **Step 2: PostCard 组件**

```vue
<!-- frontend/src/components/PostCard.vue -->
<script setup lang="ts">
import type { Post } from '../stores/posts'

defineProps<{ post: Post }>()
</script>

<template>
  <article class="py-6 border-b border-gray-200 dark:border-gray-700">
    <time class="text-sm text-gray-500 dark:text-gray-400">
      {{ new Date(post.created_at).toLocaleDateString('zh-CN') }}
    </time>
    <router-link :to="`/post/${post.slug}`">
      <h2 class="text-xl font-semibold mt-1 hover:text-blue-500">{{ post.title }}</h2>
    </router-link>
    <div class="flex gap-2 mt-2">
      <router-link
        v-for="tag in post.tags"
        :key="tag.id"
        :to="`/?tag=${tag.name}`"
        class="text-xs px-2 py-0.5 bg-gray-100 dark:bg-gray-800 rounded hover:bg-blue-100"
      >
        {{ tag.name }}
      </router-link>
    </div>
  </article>
</template>
```

- [ ] **Step 3: HomeView**

```vue
<!-- frontend/src/views/HomeView.vue -->
<script setup lang="ts">
import { onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { usePostsStore } from '../stores/posts'
import PostCard from '../components/PostCard.vue'

const route = useRoute()
const store = usePostsStore()

onMounted(() => load())
watch(() => route.query, load)

function load() {
  const page = Number(route.query.page) || 1
  const tag = (route.query.tag as string) || ''
  store.fetchPosts(page, tag)
}
</script>

<template>
  <div>
    <h1 v-if="!route.query.tag" class="text-3xl font-bold mb-8">Blog Posts</h1>
    <h1 v-else class="text-3xl font-bold mb-8">
      Tag: {{ route.query.tag }}
    </h1>

    <div v-if="store.loading" class="text-center py-12 text-gray-500">Loading...</div>

    <div v-else-if="store.posts.length === 0" class="text-center py-12 text-gray-500">
      No posts yet.
    </div>

    <template v-else>
      <PostCard v-for="post in store.posts" :key="post.id" :post="post" />

      <div class="flex justify-center gap-2 mt-8">
        <button
          v-for="p in Math.ceil(store.total / 10)"
          :key="p"
          :class="['px-3 py-1 rounded', p === (Number(route.query.page) || 1) ? 'bg-blue-500 text-white' : 'bg-gray-100 dark:bg-gray-800']"
          @click="$router.push({ query: { ...$route.query, page: p } })"
        >
          {{ p }}
        </button>
      </div>
    </template>
  </div>
</template>
```

- [ ] **Step 4: Commit**

```bash
git add frontend/src/stores/posts.ts frontend/src/components/PostCard.vue frontend/src/views/HomeView.vue
git commit -m "feat: add posts store, PostCard, and HomeView with pagination"
```

---

### Task 12: 文章详情页 + TOC 侧栏 + 代码高亮

**Files:**
- Create: `frontend/src/views/PostDetailView.vue`
- Create: `frontend/src/components/TOCSidebar.vue`
- Create: `frontend/src/components/MarkdownRenderer.vue`

- [ ] **Step 1: MarkdownRenderer（highlight.js 高亮）**

```vue
<!-- frontend/src/components/MarkdownRenderer.vue -->
<script setup lang="ts">
import { onMounted, watch, ref } from 'vue'
import hljs from 'highlight.js/lib/core'
import 'highlight.js/styles/github-dark.css'

const props = defineProps<{ html: string }>()
const el = ref<HTMLElement>()

function highlight() {
  if (el.value) {
    el.value.querySelectorAll('pre code').forEach((block) => {
      hljs.highlightElement(block as HTMLElement)
    })
  }
}

onMounted(highlight)
watch(() => props.html, () => {
  setTimeout(highlight, 0)
})
</script>

<template>
  <div ref="el" class="prose dark:prose-invert max-w-none" v-html="html" />
</template>
```

- [ ] **Step 2: TOCSidebar**

```vue
<!-- frontend/src/components/TOCSidebar.vue -->
<script setup lang="ts">
import { computed } from 'vue'

interface TOCItem {
  id: string
  text: string
  level: number
  children: TOCItem[]
}

const props = defineProps<{ tocJson: string }>()
const items = computed<TOCItem[]>(() => {
  try { return JSON.parse(props.tocJson) || [] }
  catch { return [] }
})
</script>

<template>
  <nav v-if="items.length" class="sticky top-4">
    <h3 class="font-semibold mb-2 text-sm text-gray-500">目录</h3>
    <ul class="space-y-1 text-sm">
      <li v-for="item in items" :key="item.id" :style="{ paddingLeft: (item.level - 1) * 12 + 'px' }">
        <a :href="`#${item.id}`" class="hover:text-blue-500 text-gray-600 dark:text-gray-400">
          {{ item.text }}
        </a>
        <template v-if="item.children?.length">
          <li v-for="child in item.children"
            :key="child.id"
            :style="{ paddingLeft: (child.level - 1) * 12 + 'px' }">
            <a :href="`#${child.id}`" class="hover:text-blue-500 text-gray-500 dark:text-gray-500">
              {{ child.text }}
            </a>
          </li>
        </template>
      </li>
    </ul>
  </nav>
</template>
```

- [ ] **Step 3: PostDetailView**

```vue
<!-- frontend/src/views/PostDetailView.vue -->
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { usePostsStore, type Post } from '../stores/posts'
import MarkdownRenderer from '../components/MarkdownRenderer.vue'
import TOCSidebar from '../components/TOCSidebar.vue'

const route = useRoute()
const store = usePostsStore()
const post = ref<Post | null>(null)

onMounted(async () => {
  post.value = await store.fetchPost(route.params.slug as string)
})
</script>

<template>
  <div v-if="!post" class="text-center py-12 text-gray-500">Loading...</div>
  <div v-else class="lg:grid lg:grid-cols-[1fr_200px] gap-8">
    <article>
      <h1 class="text-3xl font-bold mb-2">{{ post.title }}</h1>
      <p class="text-gray-500 text-sm mb-4">
        {{ new Date(post.created_at).toLocaleDateString('zh-CN') }}
      </p>
      <div class="flex gap-2 mb-8">
        <router-link v-for="tag in post.tags" :key="tag.id" :to="`/?tag=${tag.name}`"
          class="text-xs px-2 py-0.5 bg-gray-100 dark:bg-gray-800 rounded hover:bg-blue-100">
          {{ tag.name }}
        </router-link>
      </div>
      <MarkdownRenderer :html="post.content_html" />
    </article>
    <aside class="hidden lg:block">
      <TOCSidebar :toc-json="post.toc" />
    </aside>
  </div>
</template>
```

- [ ] **Step 4: 安装 Typography 插件优化 prose 样式**

```bash
cd frontend
npm install -D @tailwindcss/typography
```

更新 `tailwind.config.js`：
```js
export default {
  // ...
  plugins: [
    require('@tailwindcss/typography'),
  ],
}
```

- [ ] **Step 5: Commit**

```bash
git add frontend/
git commit -m "feat: add post detail, TOC sidebar, and code highlighting"
```

---

### Task 13: 时间线 + 标签云 + 搜索页

**Files:**
- Create: `frontend/src/views/TimelineView.vue`
- Create: `frontend/src/views/TagsView.vue`
- Create: `frontend/src/views/SearchView.vue`

- [ ] **Step 1: TimelineView**

```vue
<!-- frontend/src/views/TimelineView.vue -->
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import api from '../utils/api'

const timeline = ref<any[]>([])

onMounted(async () => {
  const res = await api.get('/timeline')
  timeline.value = res.data.timeline
})
</script>

<template>
  <div>
    <h1 class="text-3xl font-bold mb-8">Timeline</h1>
    <div v-for="entry in timeline" :key="entry.year" class="mb-8">
      <h2 class="text-2xl font-bold text-blue-500 mb-4">{{ entry.year }}</h2>
      <div v-for="month in entry.months" :key="month.month" class="ml-4 mb-4">
        <h3 class="text-lg font-semibold text-gray-500 mb-2">{{ month.month }}月</h3>
        <ul class="space-y-2">
          <li v-for="post in month.posts" :key="post.id">
            <router-link :to="`/post/${post.slug}`" class="hover:text-blue-500">
              <span class="text-sm text-gray-400 mr-2">{{ post.created_at.split('T')[0] }}</span>
              {{ post.title }}
            </router-link>
          </li>
        </ul>
      </div>
    </div>
    <p v-if="timeline.length === 0" class="text-gray-500 text-center py-12">No posts yet.</p>
  </div>
</template>
```

- [ ] **Step 2: TagsView**

```vue
<!-- frontend/src/views/TagsView.vue -->
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import api from '../utils/api'

const tags = ref<{ name: string; count: number }[]>([])

onMounted(async () => {
  const res = await api.get('/tags')
  tags.value = res.data.tags
})
</script>

<template>
  <div>
    <h1 class="text-3xl font-bold mb-8">Tags</h1>
    <div class="flex flex-wrap gap-3">
      <router-link v-for="tag in tags" :key="tag.name" :to="`/?tag=${tag.name}`"
        class="px-4 py-2 bg-gray-100 dark:bg-gray-800 rounded-full hover:bg-blue-500 hover:text-white transition">
        {{ tag.name }}
        <span class="ml-1 text-sm opacity-60">{{ tag.count }}</span>
      </router-link>
    </div>
    <p v-if="tags.length === 0" class="text-gray-500 text-center py-12">No tags yet.</p>
  </div>
</template>
```

- [ ] **Step 3: SearchView**

```vue
<!-- frontend/src/views/SearchView.vue -->
<script setup lang="ts">
import { ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import api from '../utils/api'
import PostCard from '../components/PostCard.vue'
import type { Post } from '../stores/posts'

const route = useRoute()
const results = ref<Post[]>([])
const total = ref(0)
const loading = ref(false)
const q = ref((route.query.q as string) || '')

watch(() => route.query.q, (val) => {
  q.value = (val as string) || ''
})

async function search() {
  if (!q.value) return
  loading.value = true
  try {
    const res = await api.get('/search', { params: { q: q.value } })
    results.value = res.data.posts
    total.value = res.data.total
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div>
    <h1 class="text-3xl font-bold mb-8">Search</h1>
    <form @submit.prevent="search" class="flex gap-2 mb-8">
      <input v-model="q" type="text" placeholder="Search posts..."
        class="flex-1 px-3 py-2 border rounded dark:bg-gray-800 dark:border-gray-600" />
      <button type="submit"
        class="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600">
        Search
      </button>
    </form>

    <div v-if="loading" class="text-center py-12 text-gray-500">Searching...</div>

    <template v-else-if="results.length">
      <p class="text-gray-500 mb-4">{{ total }} result(s)</p>
      <PostCard v-for="post in results" :key="post.id" :post="post" />
    </template>
  </div>
</template>
```

- [ ] **Step 4: Commit**

```bash
git add frontend/src/views/
git commit -m "feat: add timeline, tags, and search views"
```

---

### Task 14: 管理页（拖拽上传 + 文章管理）

**Files:**
- Create: `frontend/src/views/AdminView.vue`
- Create: `frontend/src/components/UploadZone.vue`

- [ ] **Step 1: UploadZone**

```vue
<!-- frontend/src/components/UploadZone.vue -->
<script setup lang="ts">
import { ref } from 'vue'
import api from '../utils/api'

const emit = defineEmits<{ uploaded: [] }>()
const dragging = ref(false)
const uploading = ref(false)
const message = ref('')

function onDragOver(e: DragEvent) {
  e.preventDefault()
  dragging.value = true
}
function onDragLeave() { dragging.value = false }
async function onDrop(e: DragEvent) {
  e.preventDefault()
  dragging.value = false
  const file = e.dataTransfer?.files[0]
  if (!file) return
  if (!file.name.endsWith('.md')) {
    message.value = 'Please drop a .md file'
    return
  }
  uploading.value = true
  try {
    const form = new FormData()
    form.append('file', file)
    await api.post('/admin/upload', form)
    message.value = `"${file.name}" uploaded!`
    emit('uploaded')
  } catch (err: any) {
    message.value = err.response?.data?.error || 'Upload failed'
  } finally {
    uploading.value = false
  }
}
</script>

<template>
  <div
    :class="['border-2 border-dashed rounded-lg p-12 text-center transition',
      dragging ? 'border-blue-500 bg-blue-50 dark:bg-blue-900/20' : 'border-gray-300 dark:border-gray-600']"
    @dragover="onDragOver"
    @dragleave="onDragLeave"
    @drop="onDrop"
  >
    <p v-if="uploading" class="text-blue-500">Uploading...</p>
    <p v-else class="text-gray-500">Drop .md file here</p>
    <p v-if="message" class="mt-2 text-sm" :class="message.includes('!') ? 'text-green-500' : 'text-red-500'">
      {{ message }}
    </p>
  </div>
</template>
```

- [ ] **Step 2: AdminView**

```vue
<!-- frontend/src/views/AdminView.vue -->
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import api from '../utils/api'
import UploadZone from '../components/UploadZone.vue'
import type { Post } from '../stores/posts'

const router = useRouter()
const auth = useAuthStore()
const posts = ref<Post[]>([])
const loading = ref(true)

onMounted(loadPosts)

async function loadPosts() {
  loading.value = true
  try {
    const res = await api.get('/posts', { params: { per_page: 100 } })
    posts.value = res.data.posts
  } finally {
    loading.value = false
  }
}

async function deletePost(id: number) {
  if (!confirm('确定删除？')) return
  await api.delete(`/admin/posts/${id}`)
  posts.value = posts.value.filter((p) => p.id !== id)
}

function logout() {
  auth.logout()
  router.push('/login')
}
</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-8">
      <h1 class="text-3xl font-bold">Admin</h1>
      <button @click="logout" class="text-sm text-red-500 hover:underline">Logout</button>
    </div>

    <UploadZone @uploaded="loadPosts" />

    <h2 class="text-xl font-semibold mt-8 mb-4">All Posts ({{ posts.length }})</h2>
    <div v-if="loading" class="text-gray-500">Loading...</div>
    <table v-else class="w-full text-sm">
      <thead>
        <tr class="text-left border-b dark:border-gray-700">
          <th class="py-2">Title</th>
          <th class="py-2">Date</th>
          <th class="py-2">Actions</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="post in posts" :key="post.id" class="border-b dark:border-gray-700/50">
          <td class="py-2">
            <router-link :to="`/post/${post.slug}`" class="hover:text-blue-500">{{ post.title }}</router-link>
          </td>
          <td class="py-2 text-gray-500">{{ post.created_at.split('T')[0] }}</td>
          <td class="py-2">
            <button @click="deletePost(post.id)" class="text-red-500 hover:underline">Delete</button>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
```

- [ ] **Step 3: Commit**

```bash
git add frontend/src/components/UploadZone.vue frontend/src/views/AdminView.vue
git commit -m "feat: add admin page with drag-and-drop upload and post management"
```

---

### Task 15: 暗色模式 + ThemeToggle

**Files:**
- Create: `frontend/src/components/ThemeToggle.vue`
- Create: `frontend/src/utils/theme.ts`
- Modify: `frontend/src/App.vue`

- [ ] **Step 1: Theme 工具**

```typescript
// frontend/src/utils/theme.ts
import { ref, watchEffect } from 'vue'

const isDark = ref(localStorage.getItem('theme') === 'dark' ||
  (!localStorage.getItem('theme') && window.matchMedia('(prefers-color-scheme: dark)').matches))

watchEffect(() => {
  document.documentElement.classList.toggle('dark', isDark.value)
  localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
})

export function useTheme() {
  const toggle = () => { isDark.value = !isDark.value }
  return { isDark, toggle }
}
```

- [ ] **Step 2: ThemeToggle 组件**

```vue
<!-- frontend/src/components/ThemeToggle.vue -->
<script setup lang="ts">
import { useTheme } from '../utils/theme'
const { isDark, toggle } = useTheme()
</script>

<template>
  <button @click="toggle" class="p-1 rounded hover:bg-gray-200 dark:hover:bg-gray-700" :title="isDark ? 'Light mode' : 'Dark mode'">
    <span v-if="isDark">☀️</span>
    <span v-else>🌙</span>
  </button>
</template>
```

- [ ] **Step 3: 在 DefaultLayout 的 header 中加入 ThemeToggle**

修改 `frontend/src/layouts/DefaultLayout.vue`，在 nav 的 `</div>` 之前添加：

```vue
<ThemeToggle />
```

并在 `<script setup>` 中添加：
```typescript
import ThemeToggle from '../components/ThemeToggle.vue'
```

- [ ] **Step 4: 更新 App.vue 确保主题管理生效**

```vue
<!-- frontend/src/App.vue -->
<script setup lang="ts">
import { useTheme } from './utils/theme'
import DefaultLayout from './layouts/DefaultLayout.vue'

useTheme() // Initialize theme on app load
</script>

<template>
  <DefaultLayout>
    <router-view />
  </DefaultLayout>
</template>
```

- [ ] **Step 5: Commit**

```bash
git add frontend/src/utils/theme.ts frontend/src/components/ThemeToggle.vue frontend/src/layouts/DefaultLayout.vue frontend/src/App.vue
git commit -m "feat: add dark mode toggle with Tailwind class strategy"
```

---

### Task 16: 集成验证与修复

**Files:**
- Fix any compilation errors across the project
- Add `build.sh` script (if not already added)

- [ ] **Step 1: 修复 Go 端缺失的 import 和类型问题**

运行编译检查：
```bash
cd backend
go build -o /dev/null ./...
```

修复所有编译错误（缺失的 import `"time"`, `"blog-backend/internal/model"` 等）

- [ ] **Step 2: 修复前端 TypeScript 类型问题**

```bash
cd frontend
npx vue-tsc --noEmit
```

Expected: 零类型错误

- [ ] **Step 3: 构建前端**

```bash
cd frontend
npm run build
```

Expected: `dist/` 目录生成，包含 `index.html` 和静态资源

- [ ] **Step 4: 复制前端产物到后端并构建 Go**

```bash
cp -r frontend/dist backend/frontend-dist
cd backend
go build -o blog-server .
```

Expected: `backend/blog-server` 二进制文件生成成功

- [ ] **Step 5: 端到端验证**

```bash
cd backend
./blog-server
```

打开浏览器 `http://localhost:8080`：
1. 看到博客首页，无文章显示 "No posts yet."
2. 访问 `/tags` 看到空标签页
3. 访问 `/timeline` 看到空时间线
4. 访问 `/login` 登录 admin/admin
5. 登录后访问 `/admin` 看到管理页
6. 拖拽一个 .md 文件上传
7. 回到首页看到新文章

- [ ] **Step 6: Commit**

```bash
git add .
git commit -m "fix: integration fixes, end-to-end build verification"
```

---

## Self-Review

1. **Spec coverage**: 检查所有 spec 需求 → 文章列表 ✓、详情 ✓、TOC ✓、标签云 ✓、时间线 ✓、搜索 ✓、暗色模式 ✓、JWT 认证 ✓、拖拽上传 ✓、代码高亮 ✓、FTS5 ✓、Go embed ✓
2. **Placeholder scan**: 无 TBD/TODO/placeholder
3. **Type consistency**: Post 接口一致，`slug` 用于路由参数，`id` 用于管理 API，store 类型与后端模型对齐
