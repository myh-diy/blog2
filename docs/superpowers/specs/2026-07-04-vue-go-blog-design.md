# Vue + Go 个人学习日记博客 — 设计文档

**日期**: 2026-07-04
**状态**: 设计完成，待用户审阅

---

## 1. 项目概述

个人学习日记博客。管理员通过拖拽上传 Markdown 文件，系统自动解析并发布。前端 Vue 3 SPA，后端 Go + Gin，数据库 SQLite。

## 2. 核心功能

| 模块 | 功能 |
|------|------|
| **Markdown 上传** | 管理员拖拽上传 .md 文件，自动解析标题、日期、标签、正文 |
| **文章解析** | goldmark 解析，提取标题/日期/标签/TOC，代码块 Shiki 高亮 |
| **文章展示** | 博客列表页 + 文章详情页（含 TOC 侧栏） |
| **时间线** | 按年月归档展示 |
| **标签云** | 所有标签及对应文章数 |
| **全文搜索** | SQLite FTS5 实现 |
| **暗色模式** | Tailwind CSS 内置支持 |
| **JWT 认证** | 管理员登录，7 天过期 |

## 3. 技术栈

| 层 | 选型 | 说明 |
|----|------|------|
| 前端框架 | Vue 3 + Composition API | SPA |
| 构建工具 | Vite | 快速开发构建 |
| 路由 | Vue Router 4 | 客户端路由 |
| 状态管理 | Pinia | 用户态、文章数据 |
| CSS | Tailwind CSS 3 | 暗色模式 class 策略 |
| 代码高亮 | Shiki | 服务端渲染，支持所有语言 |
| 后端框架 | Gin | Go 端 HTTP 路由 + 中间件 |
| Markdown 解析 | goldmark | Go 端解析，扩展丰富（TOC、表格等） |
| 数据库 | SQLite + FTS5 | 单文件，零配置，全文搜索 |
| ORM | GORM | Go SQLite 驱动 |
| JWT | golang-jwt/jwt | HS256 签名，7 天过期 |
| 部署 | 单一二进制 + 静态文件目录 | Go embed 内嵌前端静态资源 |

## 4. 数据库设计

### 4.1 users 表

| 字段 | 类型 | 说明 |
|------|------|------|
| id | INTEGER PK AUTOINCREMENT | 主键 |
| username | TEXT UNIQUE NOT NULL | 用户名 |
| password_hash | TEXT NOT NULL | bcrypt 哈希 |
| created_at | DATETIME | 创建时间 |

### 4.2 posts 表

| 字段 | 类型 | 说明 |
|------|------|------|
| id | INTEGER PK AUTOINCREMENT | 主键 |
| title | TEXT NOT NULL | 文章标题 |
| slug | TEXT UNIQUE NOT NULL | URL 友好的标识符 |
| content_md | TEXT NOT NULL | 原始 Markdown |
| content_html | TEXT NOT NULL | 解析后的 HTML |
| toc_json | TEXT | TOC JSON，服务端解析后存储 |
| created_at | DATETIME | 原始创建日期（从 MD 解析或上传时间） |
| updated_at | DATETIME | 最后修改时间 |

### 4.3 tags 表

| 字段 | 类型 | 说明 |
|------|------|------|
| id | INTEGER PK AUTOINCREMENT | 主键 |
| name | TEXT UNIQUE NOT NULL | 标签名 |

### 4.4 post_tags 表（多对多关联）

| 字段 | 类型 | 说明 |
|------|------|------|
| post_id | INTEGER FK | 文章 ID |
| tag_id | INTEGER FK | 标签 ID |
| PK: (post_id, tag_id) | | |

### 4.5 posts_fts 表（SQLite FTS5 全文索引）

```sql
CREATE VIRTUAL TABLE posts_fts USING fts5(
    title,
    content_md,
    content='posts',
    content_rowid='id'
);
```

## 5. API 设计

### 5.1 公开接口

```
GET /api/posts?page=1&per_page=10&tag=xxx
  → { posts: [...], total: int, page: int, per_page: int }

GET /api/posts/:slug
  → { post: { title, slug, content_html, toc, tags, created_at, updated_at } }

GET /api/tags
  → { tags: [{ name: string, count: int }] }

GET /api/search?q=keyword&page=1
  → { posts: [...], total: int, page: int }   // FTS5 查询

GET /api/timeline
  → { years: [{ year: 2026, months: [{ month: 7, posts: [...] }] }] }
```

### 5.2 认证接口

```
POST /api/auth/login
  请求: { username: string, password: string }
  返回: { token: string, expires_at: string }
  错误: 401
```

### 5.3 管理员接口（需要 Authorization: Bearer <token>）

```
POST /api/admin/upload
  Content-Type: multipart/form-data
  字段: file (markdown file)
  返回: { post: {...} }

DELETE /api/admin/posts/:id
  → 204

PUT /api/admin/posts/:id
  请求: { title?, content_md?, tags? }
  → { post: {...} }
```

## 6. 前端路由

```
/              → 首页（文章列表 + 分页 + 标签筛选）
/post/:slug    → 文章详情（正文 + TOC 侧栏）
/timeline      → 时间线
/tags          → 标签云
/search?q=xxx  → 全文搜索结果
/login         → 管理员登录
/admin         → 管理页（拖拽上传 + 文章管理列表）
```

## 7. 前端组件树（主要组件）

```
App.vue
├── layouts/
│   ├── DefaultLayout.vue     // 公开页布局（Header + 内容 + Footer）
│   └── AdminLayout.vue       // 管理页布局（侧栏 + 内容）
├── components/
│   ├── PostCard.vue          // 文章列表卡片
│   ├── PostList.vue          // 文章列表（含分页）
│   ├── TagCloud.vue          // 标签云
│   ├── Timeline.vue          // 时间线条目
│   ├── SearchBar.vue         // 搜索框
│   ├── TOCSidebar.vue        // 文章 TOC 侧栏
│   ├── MarkdownRenderer.vue  // HTML 渲染器
│   ├── UploadZone.vue        // 拖拽上传区域
│   ├── ThemeToggle.vue       // 暗色模式切换
│   └── LoadingSpinner.vue    // 加载状态
├── views/
│   ├── HomeView.vue
│   ├── PostDetailView.vue
│   ├── TimelineView.vue
│   ├── TagsView.vue
│   ├── SearchView.vue
│   ├── LoginView.vue
│   └── AdminView.vue
├── stores/
│   ├── auth.ts               // JWT 管理（登录/登出/Token 刷新）
│   └── posts.ts               // 文章数据缓存
├── router/
│   └── index.ts              // 路由配置 + 导航守卫
└── utils/
    ├── api.ts                // Axios 封装（自动带 Token）
    └── theme.ts              // 暗色模式管理
```

## 8. Go 后端包结构

```
backend/
├── main.go                   // 入口，启动服务
├── config/
│   └── config.go             // 配置管理（端口、数据库路径、JWT 密钥）
├── internal/
│   ├── auth/
│   │   ├── jwt.go            // JWT 签发与验证
│   │   └── middleware.go     // Gin 中间件：JWT 校验
│   ├── handler/
│   │   ├── auth.go           // 登录
│   │   ├── post.go           // 文章 CRUD + 搜索
│   │   ├── tag.go            // 标签
│   │   ├── upload.go         // MD 文件上传 + 解析
│   │   └── timeline.go       // 时间线
│   ├── model/
│   │   ├── user.go
│   │   ├── post.go
│   │   └── tag.go
│   ├── parser/
│   │   └── markdown.go       // goldmark 解析 + TOC 提取
│   ├── service/
│   │   ├── post.go           // 文章业务逻辑
│   │   └── search.go          // FTS5 搜索
│   └── database/
│       ├── sqlite.go         // 数据库初始化 + 迁移
│       └── fts.go            // FTS5 索引管理
├── go.mod
└── go.sum
```

## 9. Markdown 解析流程

```
用户拖拽上传 .md 文件
        │
        ▼
后端 gin handler 接收 multipart
        │
        ▼
提取 frontmatter:
  ---
  title: xxx
  date: 2026-07-04
  tags: [Go, Vue]
  ---
        │
        ▼
goldmark 解析正文 → HTML
goldmark-meta 解析 frontmatter → title/date/tags
goldmark-toc 提取标题层级 → 生成 TOC JSON
        │
        ▼
shiki 对代码块做语法高亮（返回带样式的 HTML）
        │
        ▼
存入 SQLite:
  - content_md (原始)
  - content_html (渲染后)
  - toc_json
  - tags 表关联
        │
        ▼
返回结果给前端
```

## 10. 部署架构

```
博客服务器
├── blog-server (Go 二进制)
│   ├── 内嵌前端静态文件 (embed.FS)
│   ├── API 路由 (/api/*)
│   └── 静态文件托管 (/* → SPA index.html)
└── blog.db (SQLite 数据文件)

启动命令:  ./blog-server --port 8080 --db ./blog.db
```

---

## Appendix: 目录结构总览

```
blog2/
├── frontend/                # Vue 3 + Vite
│   ├── src/
│   │   ├── components/
│   │   ├── views/
│   │   ├── stores/
│   │   ├── router/
│   │   ├── utils/
│   │   ├── App.vue
│   │   └── main.ts
│   ├── index.html
│   ├── package.json
│   ├── vite.config.ts
│   └── tailwind.config.js
├── backend/                 # Go + Gin
│   ├── main.go
│   ├── go.mod
│   └── internal/
├── docs/
│   └── superpowers/
│       └── specs/
│           └── 2026-07-04-vue-go-blog-design.md  ← 本文档
└── .claude/
    ├── commands/            # frontend-design skills
    └── skills/              # superpowers skills
```
