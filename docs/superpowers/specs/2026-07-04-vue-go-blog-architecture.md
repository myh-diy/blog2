# Vue + Go 博客 — 架构文档

**日期**: 2026-07-04
**关联设计**: [2026-07-04-vue-go-blog-design.md](../specs/2026-07-04-vue-go-blog-design.md)
**关联计划**: [2026-07-04-vue-go-blog.md](../plans/2026-07-04-vue-go-blog.md)

---

## 1. 系统架构总览

```
┌──────────────────────────────────────────────────────┐
│                      浏览器                           │
│  ┌────────────────────────────────────────────────┐  │
│  │           Vue 3 SPA (Vite 构建)                 │  │
│  │  ┌──────────┐ ┌──────────┐ ┌───────────────┐  │  │
│  │  │  Views   │ │  Stores  │ │    Router      │  │  │
│  │  │ (7 pages)│ │ (Pinia)  │ │ (导航守卫+JWT)  │  │  │
│  │  └────┬─────┘ └────┬─────┘ └───────┬───────┘  │  │
│  │       │             │               │          │  │
│  │  ┌────┴─────────────┴───────────────┴───────┐  │  │
│  │  │           Axios (api.ts)                  │  │  │
│  │  │   • 自动附加 JWT Bearer Token              │  │  │
│  │  │   • 401 响应 → 跳转登录页                   │  │  │
│  │  └──────────────────┬───────────────────────┘  │  │
│  └─────────────────────┼──────────────────────────┘  │
└────────────────────────┼─────────────────────────────┘
                         │ HTTP/HTTPS
                         ▼
┌──────────────────────────────────────────────────────┐
│                  Go + Gin 后端 (:8080)                │
│                                                      │
│  ┌──────────┐  ┌───────────┐  ┌──────────────────┐  │
│  │  Router  │  │  Handler  │  │   Middleware      │  │
│  │ (Gin)    │──│ (7 files) │  │  • JWT Auth       │  │
│  │          │  │           │  │  • CORS (optional) │  │
│  └──────────┘  └─────┬─────┘  └──────────────────┘  │
│                      │                               │
│              ┌───────┼───────┐                       │
│              ▼       ▼       ▼                       │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐             │
│  │  Service │ │  Parser  │ │  Auth    │             │
│  │ (业务逻辑)│ │ (goldmark)│ │ (JWT)    │             │
│  └────┬─────┘ └──────────┘ └──────────┘             │
│       │                                              │
│  ┌────┴──────────────────────────┐                   │
│  │         GORM + SQLite         │                   │
│  │  ┌─────────────────────────┐ │                   │
│  │  │ posts │ tags │ post_tags │ │                   │
│  │  │ users │ posts_fts (FTS5) │ │                   │
│  │  └─────────────────────────┘ │                   │
│  └──────────────────────────────┘                   │
│                                                      │
│  ┌──────────────────────────────────────────────┐   │
│  │  embed.FS (frontend-dist/)                     │   │
│  │  • SPA index.html + JS/CSS 静态文件             │   │
│  └──────────────────────────────────────────────┘   │
└──────────────────────────────────────────────────────┘
```

## 2. 请求生命周期

### 2.1 公开请求（文章列表）

```
浏览器 GET /api/posts?page=1&tag=Go
  │
  ▼
Gin Router → handler.ListPosts()
  │
  ▼
c.Query("page"), c.Query("tag")
  │
  ▼
service.GetPosts(db, page, perPage, tag)
  │ → GORM Preload("Tags") 查询
  │ → COUNT 获取总数
  │ → ORDER BY created_at DESC
  │ → LIMIT/OFFSET 分页
  ▼
JSON: { posts, total, page, per_page }
  │
  ▼
Vue Axios → posts Store → PostCard 组件渲染
```

### 2.2 认证请求（上传 MD）

```
浏览器 POST /api/admin/upload (multipart/form-data)
  │ Authorization: Bearer <JWT>
  ▼
Gin Router → admin Group → auth.AuthMiddleware()
  │ → Parse JWT, validate signature + expiry
  │ → Extract userID, set in context
  │ → 失败则 401，中断链路
  ▼
handler.UploadPost()
  │ → c.FormFile("file")
  │ → io.ReadAll(file)
  │ → parser.ParseMarkdown(content)
  │   ├─ goldmark 解析 MD → HTML
  │   ├─ goldmark-meta 提取 frontmatter
  │   ├─ goldmark-toc 提取目录
  │   └─ slugify(title) 生成 URL
  │ → service.CreatePost(db, result, rawMD)
  │   ├─ 标签 FirstOrCreate
  │   ├─ INSERT INTO posts
  │   ├─ INSERT INTO post_tags
  │   └─ INSERT INTO posts_fts (FTS5 触发器自动执行)
  ▼
JSON: { post: { id, title, slug, ... } }
```

### 2.3 搜索请求

```
浏览器 GET /api/search?q=goroutine&page=1
  │
  ▼
handler.SearchPosts()
  │
  ▼
SELECT rowid FROM posts_fts WHERE posts_fts MATCH 'goroutine'
  ORDER BY rank LIMIT 10 OFFSET 0
  │
  ▼
SELECT * FROM posts WHERE id IN (result rows)
  Preload("Tags")
  ORDER BY created_at DESC
  │
  ▼
JSON: { posts, total, page, per_page }
```

## 3. 数据流图

```
┌─────────────────────────────────────────────────────────┐
│                     Markdown 上传流程                     │
│                                                         │
│  .md 文件拖拽                                             │
│      │                                                  │
│      ▼                                                  │
│  ┌──────────────┐                                       │
│  │ UploadZone   │  Vue 组件，接收 File 对象               │
│  │ (drag & drop)│                                       │
│  └──────┬───────┘                                       │
│         │ FormData { file }                             │
│         ▼                                               │
│  ┌──────────────┐                                       │
│  │ Axios POST   │  Authorization: Bearer <token>        │
│  │ /api/admin/  │                                       │
│  │   upload     │                                       │
│  └──────┬───────┘                                       │
│         │                                               │
│  ╔══════╪══════════════════════════════════════════╗    │
│  ║  Go  │  Handler                                  ║    │
│  ║      ▼                                           ║    │
│  ║ ┌──────────────┐  ┌─────────────────────────┐   ║    │
│  ║ │ parser.Parse │  │ goldmark 解析流程:        │   ║    │
│  ║ │ Markdown()   │─▶│ 1. YAML frontmatter     │   ║    │
│  ║ └──────┬───────┘  │ 2. GFM/Table 扩展       │   ║    │
│  ║        │          │ 3. TOC 提取              │   ║    │
│  ║        ▼          │ 4. Heading ID 自动生成    │   ║    │
│  ║ ┌──────────────┐  └─────────────────────────┘   ║    │
│  ║ │ ParseResult  │                                 ║    │
│  ║ │ • Title      │                                 ║    │
│  ║ │ • Date       │                                 ║    │
│  ║ │ • Tags[]     │                                 ║    │
│  ║ │ • HTML       │                                 ║    │
│  ║ │ • TOCJSON    │                                 ║    │
│  ║ │ • Slug       │                                 ║    │
│  ║ └──────┬───────┘                                 ║    │
│  ║        ▼                                         ║    │
│  ║ ┌──────────────┐                                 ║    │
│  ║ │ service.     │  GORM 写入:                      ║    │
│  ║ │ CreatePost() │  • Posts (title/slug/html/toc)  ║    │
│  ║ └──────┬───────┘  • Tags (FirstOrCreate)          ║    │
│  ║        │          • PostTags (关联表)              ║    │
│  ╚════════╪══════════════════════════════════════╝    │
│           │                                            │
│           ▼                                            │
│  ┌────────────────┐                                    │
│  │   SQLite       │                                    │
│  │   posts 表      │                                    │
│  │   tags 表       │                                    │
│  │   post_tags 表  │                                    │
│  │   posts_fts 表  │  ← FTS5 触发器自动同步              │
│  └────────────────┘                                    │
│           │                                            │
│           ▼                                            │
│  ┌────────────────┐                                    │
│  │ 201 Created    │  JSON { post: {...} }              │
│  └────────────────┘                                    │
│           │                                            │
│           ▼                                            │
│  ┌────────────────┐                                    │
│  │  Vue 响应       │  emit('uploaded')                  │
│  │  → loadPosts() │  → 刷新管理页列表                   │
│  └────────────────┘                                    │
└─────────────────────────────────────────────────────────┘
```

## 4. 前端架构

### 4.1 组件层级

```
App.vue
├── DefaultLayout.vue
│   ├── Header (nav)
│   │   ├── <router-link> (Home/Timeline/Tags/Search/Admin)
│   │   └── ThemeToggle.vue  ← useTheme() composable
│   └── <router-view>
│       ├── HomeView.vue
│       │   └── PostCard.vue (v-for)
│       ├── PostDetailView.vue
│       │   ├── MarkdownRenderer.vue  ← highlight.js
│       │   └── TOCSidebar.vue        ← sticky position
│       ├── TimelineView.vue
│       ├── TagsView.vue
│       ├── SearchView.vue
│       │   └── PostCard.vue
│       ├── LoginView.vue
│       └── AdminView.vue
│           ├── UploadZone.vue   ← drag & drop
│           └── Post Table       ← CRUD
```

### 4.2 状态管理

```
Pinia Stores
│
├── auth Store
│   State:  token (string), isAuthenticated (boolean)
│   Actions: login(), logout()
│   Persist: localStorage.getItem('token')
│
└── posts Store
    State:  posts[], total, loading
    Actions: fetchPosts(page, tag), fetchPost(slug)
    Cache:   无持久化，每次页面访问重新请求
```

### 4.3 路由守卫流程

```
用户访问 /admin
  │
  ▼
router.beforeEach(to, from, next)
  │
  ├─ to.meta.requiresAuth === true?
  │   ├─ Yes → localStorage 有 token?
  │   │   ├─ Yes → next() ✓
  │   │   └─ No  → next('/login')
  │   └─ No  → next() ✓
  │
  ▼
组件渲染
  │
  ▼
Axios interceptor:
  • request: 自动附加 Authorization header
  • response: 401 → 清除 token, 跳转 /login
```

## 5. 后端架构

### 5.1 分层设计

```
main.go (入口 + 路由注册 + embed 配置)
  │
  ├─ config/     → 环境变量 / 默认值
  ├─ internal/
  │   ├─ auth/
  │   │   ├─ jwt.go        → GenerateToken / ValidateToken
  │   │   └─ middleware.go  → AuthMiddleware (Gin HandlerFunc)
  │   ├─ handler/           → HTTP 层（请求解析 → 调用 service → 响应序列化）
  │   │   ├─ auth.go        → POST /api/auth/login
  │   │   ├─ post.go        → GET /api/posts, GET /api/posts/:slug
  │   │   ├─ upload.go      → POST /api/admin/upload
  │   │   ├─ tag.go         → GET /api/tags
  │   │   ├─ search.go      → GET /api/search
  │   │   └─ timeline.go    → GET /api/timeline
  │   ├─ service/           → 业务逻辑层（与数据库交互）
  │   │   ├─ post.go        → CRUD + 搜索 + 时间线
  │   │   ├─ tag.go         → 标签聚合查询
  │   │   └─ user.go        → 认证 + 默认用户种子
  │   ├─ model/             → GORM 模型定义
  │   │   ├─ user.go        → User
  │   │   ├─ post.go        → Post, TimelineEntry
  │   │   └─ tag.go         → Tag
  │   ├─ parser/
  │   │   └─ markdown.go    → goldmark 解析 + frontmatter + TOC
  │   └─ database/
  │       ├─ sqlite.go      → GORM 初始化 + AutoMigrate
  │       └─ fts.go         → FTS5 虚拟表 + 触发器
```

### 5.2 中间件链

```
请求进入
  │
  ▼
[1] gin.Default()
    • Logger (请求日志)
    • Recovery (panic → 500)
  │
  ▼
[2] auth.AuthMiddleware()  ← 仅 /api/admin/* 路由组
    • 检查 Authorization header
    • Parse + Validate JWT
    • 提取 userID → c.Set("userID", claims.UserID)
    • 失败 → 401 JSON
  │
  ▼
[3] handler 处理函数
```

## 6. 安全架构

```
┌──────────────────────────────────────┐
│              安全层次                  │
├──────────────────────────────────────┤
│ 传输层: HTTPS (TLS)                   │
│   • 生产环境建议 Caddy 自动证书        │
├──────────────────────────────────────┤
│ 认证层: JWT Bearer Token              │
│   • HS256 签名，7 天过期               │
│   • 密钥通过环境变量 JWT_SECRET 配置    │
│   • 默认值仅用于开发                    │
├──────────────────────────────────────┤
│ 密码层: bcrypt ($2a$ cost=10)         │
│   • 不可逆哈希                          │
│   • 默认密码 "admin" 仅首次种子         │
├──────────────────────────────────────┤
│ 授权层: 路由级中间件                   │
│   • /api/* 公开读写                    │
│   • /api/admin/* 需 JWT               │
│   • 管理员接口: 上传/更新/删除          │
├──────────────────────────────────────┤
│ 输入层: Gin binding + 校验            │
│   • JSON body → ShouldBindJSON        │
│   • FormData → FormFile + 文件类型校验 │
│   • URL param → strconv 类型转换       │
│   • SQL → GORM 参数化查询 (防注入)     │
│   • HTML 输出 → goldmark 安全渲染      │
├──────────────────────────────────────┤
│ 跨站: 个人博客，无需 CSRF              │
│   • JWT 存 localStorage (非 Cookie)   │
│   • CORS 默认关闭，仅同源访问           │
└──────────────────────────────────────┘
```

## 7. 部署架构

```
┌─────────────────────────────────────────┐
│           部署方式: 单一二进制             │
│                                         │
│  blog-server (Go binary)                │
│  ├─ 内嵌 frontend-dist/ (embed.FS)      │
│  │   ├─ index.html                      │
│  │   ├─ assets/index-xxx.js             │
│  │   └─ assets/index-xxx.css            │
│  ├─ API 路由 (/api/*)                   │
│  └─ SPA fallback (/* → index.html)     │
│                                         │
│  blog.db (SQLite, 同目录)               │
│                                         │
│  启动: ./blog-server                     │
│  环境变量:                               │
│    PORT=8080                            │
│    DB_PATH=./blog.db                    │
│    JWT_SECRET=<随机密钥>                 │
└─────────────────────────────────────────┘

反向代理 (可选):
  Caddy / Nginx → 静态文件直出 + /api/* 代理到 blog-server
```

## 8. 构建流水线

```
┌──────────┐    ┌──────────────┐    ┌──────────────────┐
│ frontend │    │   copy dist  │    │     backend       │
│ npm run  │───▶│ frontend/dist│───▶│ go build          │
│  build   │    │     →        │    │  -o blog-server . │
│          │    │ backend/     │    │                   │
│          │    │ frontend-dist│    │ embed 自动包含     │
└──────────┘    └──────────────┘    └──────────────────┘
                                            │
                                            ▼
                                     ┌──────────────┐
                                     │ blog-server  │
                                     │  (+ blog.db) │
                                     └──────────────┘

单命令构建脚本: ./build.sh
```

## 9. 错误处理策略

```
Go 后端:
┌───────────────────────────────────────┐
│ handler 层:                           │
│   • 参数错误 → 400 + JSON error       │
│   • 认证失败 → 401 + JSON error       │
│   • 资源不存在 → 404 + JSON error     │
│   • 服务端错误 → 500 + JSON error     │
│   • panic → Recovery 中间件 → 500     │
│                                       │
│ service 层:                           │
│   • 数据库错误 → error 向上传播        │
│   • 无记录 → gorm.ErrRecordNotFound   │
│                                       │
│ parser 层:                            │
│   • 解析失败 → error (含上下文信息)    │
│   • frontmatter 缺失 → 使用默认值     │
└───────────────────────────────────────┘

Vue 前端:
┌───────────────────────────────────────┐
│ Axios interceptor:                    │
│   • 网络错误 → 控制台 + 用户提示       │
│   • 401 → 清除 token, 跳转 /login     │
│   • 其他 HTTP 错误 → 组件级处理        │
│                                       │
│ Router:                               │
│   • 404 路由 → 默认显示空状态          │
│   • 懒加载失败 → Vite 自动重试         │
│                                       │
│ 组件:                                 │
│   • loading 状态 → spinner / skeleton │
│   • empty 状态 → "No posts yet"       │
│   • error 状态 → 错误消息 + 重试按钮   │
└───────────────────────────────────────┘
```

## 10. 关键设计决策

| 决策 | 选择 | 原因 |
|------|------|------|
| 代码高亮位置 | 前端 (highlight.js) | 避免 Node.js sidecar，减少部署复杂度 |
| FTS 同步方式 | SQLite 触发器 | 零代码同步，无需 service 层手动维护 |
| 前端渲染模式 | SPA | 个人日记不需要 SEO，上传即见 |
| 数据库 | SQLite | 单文件零运维，适合个人博客 |
| TOC 存储 | JSON 预计算存储 | 避免每次请求重新解析 |
| Go embed | 前端 dist 目录 | 单一二进制，部署简单 |
| 居中条件 | 默认 max-w-4xl (896px) | 博客阅读友好宽度 |
