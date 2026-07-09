# 前端 B站少女粉紫风格重设计实施计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 在不引入第三方 UI 组件库的前提下，将博客前端整体视觉改造为粉紫渐变、卡片化、圆润可爱的 B站少女风格，并修复既有视觉/交互 bug。

**Architecture:** 以 Tailwind CSS 为核心，扩展 brand（粉）和 accent（紫）色阶，新增 3 个可复用组件（GradientButton、KawaiiIcon、EmptyState），逐个页面替换为卡片化布局；所有样式统一走 Tailwind 工具类，移除 PostsView/UploadZone 中的原生 CSS。

**Tech Stack:** Vue 3, TypeScript, Vite, Tailwind CSS 3.4, Pinia, highlight.js

## Global Constraints

- 不引入第三方 Vue UI 组件库
- 所有颜色必须来自 Tailwind 扩展色阶或合法的任意值语法
- 深色/浅色模式通过 `darkMode: 'class'` 和 `dark:` 前缀支持
- 图标使用内联 SVG，不引入图标字体库
- 二次元插画使用内联 SVG，风格统一为圆润线条、粉紫配色
- 保持现有路由、状态管理、API 调用不变
- 每次任务完成后运行 `npm run build` 或对应验证命令

---

## 文件变更总览

| 文件 | 操作 | 说明 |
|---|---|---|
| `frontend/tailwind.config.js` | 修改 | 扩展 brand/accent 色阶 |
| `frontend/src/style.css` | 修改 | 更新基础样式、selection、tag-pill、btn-primary |
| `frontend/public/favicon.svg` | 修改 | 改为粉紫渐变闪电/星星风格 |
| `frontend/src/components/GradientButton.vue` | 创建 | 粉紫渐变按钮 |
| `frontend/src/components/KawaiiIcon.vue` | 创建 | 可复用二次元 SVG 插画 |
| `frontend/src/components/EmptyState.vue` | 创建 | 空状态插画组件 |
| `frontend/src/layouts/DefaultLayout.vue` | 修改 | 新导航栏、footer、max-w-6xl |
| `frontend/src/views/HomeView.vue` | 修改 | 移除固定定位，卡片化 hero + 文章网格 |
| `frontend/src/components/PostCard.vue` | 修改 | 大圆角卡片、hover 上浮、修复标签链接 |
| `frontend/src/views/PostsView.vue` | 修改 | 卡片网格布局，移除原生 CSS |
| `frontend/src/views/PostDetailView.vue` | 修改 | 卡片化布局，修复标签链接 |
| `frontend/src/components/MarkdownRenderer.vue` | 修改 | 注册语言、主题随 dark mode 切换 |
| `frontend/src/components/TOCSidebar.vue` | 修改 | 圆角卡片、当前项高亮 |
| `frontend/src/views/TagsView.vue` | 修改 | 彩色胶囊标签卡片网格 |
| `frontend/src/views/TimelineView.vue` | 修改 | 卡片时间轴 |
| `frontend/src/views/SearchView.vue` | 修改 | 居中搜索 + EmptyState |
| `frontend/src/views/LoginView.vue` | 修改 | 居中卡片、粉紫渐变 |
| `frontend/src/views/AdminView.vue` | 修改 | 卡片化后台 |
| `frontend/src/components/UploadZone.vue` | 修改 | Tailwind 重写 |
| `frontend/src/components/ThemeToggle.vue` | 修改 | 统一太阳/月亮 SVG，粉紫强调色 |

---

## Task 1: 基础配置与全局样式

**Files:**
- Modify: `frontend/tailwind.config.js`
- Modify: `frontend/src/style.css`
- Modify: `frontend/public/favicon.svg`

**Interfaces:**
- Produces: `brand` 色阶（粉系）和 `accent` 色阶（紫系）可供所有组件使用
- Produces: `.tag-pill`、`.btn-primary` 更新为粉紫风格

- [ ] **Step 1: 更新 tailwind.config.js**

将现有 `brand` 改为粉系，`pop` 改为 `accent` 紫系：

```js
colors: {
  brand: {
    50:  '#fff1f8',
    100: '#ffe4f3',
    200: '#ffc9e9',
    300: '#ff9dd4',
    400: '#ff5fb8',
    500: '#fb7299',
    600: '#e84a7a',
    700: '#c93062',
    800: '#a62b52',
    900: '#8b2848',
  },
  accent: {
    50:  '#f5f3ff',
    100: '#ede9fe',
    200: '#ddd6fe',
    300: '#c4b5fd',
    400: '#a78bfa',
    500: '#8b5cf6',
    600: '#7c3aed',
    700: '#6d28d9',
    800: '#5b21b6',
    900: '#4c1d95',
  },
}
```

- [ ] **Step 2: 更新 style.css**

```css
@layer base {
  html { scroll-behavior: smooth; }
  html, body { overflow-x: hidden; max-width: 100%; }
  body {
    @apply bg-gray-50 text-slate-700 antialiased;
    @apply dark:bg-slate-950 dark:text-slate-200;
  }
  ::selection {
    @apply bg-brand-200 dark:bg-brand-700/50;
  }
}

@layer components {
  .tag-pill {
    @apply inline-flex items-center text-xs px-2.5 py-1 rounded-full font-medium
           bg-brand-100 text-brand-700 dark:bg-brand-900/40 dark:text-brand-300
           hover:bg-brand-200 dark:hover:bg-brand-900/60 transition-colors;
  }
  .btn-primary {
    @apply inline-flex items-center justify-center px-4 py-2 rounded-xl font-semibold text-white
           bg-gradient-to-r from-brand-400 to-accent-500
           hover:from-brand-500 hover:to-accent-600
           active:scale-[0.98]
           transition-all shadow-md shadow-brand-500/20 hover:shadow-lg hover:shadow-brand-500/30;
  }
}
```

- [ ] **Step 3: 更新 favicon.svg**

改为粉紫渐变星星/闪电：

```svg
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100">
  <defs>
    <linearGradient id="g" x1="0" y1="0" x2="1" y2="1">
      <stop offset="0%" stop-color="#ff5fb8"/>
      <stop offset="100%" stop-color="#8b5cf6"/>
    </linearGradient>
  </defs>
  <rect width="100" height="100" rx="22" fill="url(#g)"/>
  <path d="M50 20 L58 42 L82 42 L63 56 L71 78 L50 64 L29 78 L37 56 L18 42 L42 42 Z" fill="white"/>
</svg>
```

- [ ] **Step 4: 验证构建**

Run: `cd frontend && npm run build`
Expected: PASS（仅改配置不应报错）

- [ ] **Step 5: 提交**

```bash
git add frontend/tailwind.config.js frontend/src/style.css frontend/public/favicon.svg
git commit -m "feat(ui): setup pink-purple color system and global styles"
```

---

## Task 2: 新增通用组件

**Files:**
- Create: `frontend/src/components/GradientButton.vue`
- Create: `frontend/src/components/KawaiiIcon.vue`
- Create: `frontend/src/components/EmptyState.vue`

**Interfaces:**
- Consumes: `brand-*` / `accent-*` Tailwind classes
- Produces: `<GradientButton>`、`<KawaiiIcon name="...">`、`<EmptyState title="..." description="...">`

- [ ] **Step 1: 创建 GradientButton.vue**

```vue
<script setup lang="ts">
defineProps<{ type?: 'button' | 'submit' }>()
</script>

<template>
  <button :type="type || 'button'"
    class="inline-flex items-center justify-center px-5 py-2.5 rounded-xl font-semibold text-white
           bg-gradient-to-r from-brand-400 to-accent-500
           hover:from-brand-500 hover:to-accent-600
           hover:-translate-y-0.5 hover:shadow-lg hover:shadow-brand-500/25
           active:scale-[0.98] active:translate-y-0
           transition-all duration-200 disabled:opacity-60 disabled:cursor-not-allowed">
    <slot />
  </button>
</template>
```

- [ ] **Step 2: 创建 KawaiiIcon.vue**

```vue
<script setup lang="ts">
const props = defineProps<{ name: 'wave' | 'search' | 'sad' | 'happy' | 'star' | 'heart' }>()
</script>

<template>
  <svg v-if="props.name === 'wave'" viewBox="0 0 120 120" class="w-full h-full">
    <circle cx="60" cy="60" r="56" fill="#fff1f8"/>
    <circle cx="42" cy="50" r="6" fill="#8b5cf6"/>
    <circle cx="78" cy="50" r="6" fill="#8b5cf6"/>
    <path d="M45 72 Q60 85 75 72" stroke="#fb7299" stroke-width="4" fill="none" stroke-linecap="round"/>
    <circle cx="20" cy="35" r="8" fill="#ffc9e9"/>
    <circle cx="100" cy="30" r="6" fill="#ddd6fe"/>
    <circle cx="105" cy="80" r="5" fill="#ff9dd4"/>
  </svg>
  <!-- search, sad, happy, star, heart 类似，保持圆润粉紫风格 -->
</template>
```

（实际实现时为每个 name 补全 SVG，风格统一。）

- [ ] **Step 3: 创建 EmptyState.vue**

```vue
<script setup lang="ts">
import KawaiiIcon from './KawaiiIcon.vue'
defineProps<{ icon: 'search' | 'sad' | 'happy'; title: string; description?: string }>()
</script>

<template>
  <div class="flex flex-col items-center justify-center py-16 text-center">
    <div class="w-28 h-28 mb-4">
      <KawaiiIcon :name="icon" />
    </div>
    <h3 class="text-lg font-bold text-slate-700 dark:text-slate-200">{{ title }}</h3>
    <p v-if="description" class="mt-1 text-sm text-slate-400 dark:text-slate-500 max-w-xs">{{ description }}</p>
    <slot />
  </div>
</template>
```

- [ ] **Step 4: 验证构建**

Run: `cd frontend && npm run build`
Expected: PASS

- [ ] **Step 5: 提交**

```bash
git add frontend/src/components/GradientButton.vue frontend/src/components/KawaiiIcon.vue frontend/src/components/EmptyState.vue
git commit -m "feat(ui): add GradientButton, KawaiiIcon, EmptyState components"
```

---

## Task 3: 全局布局 DefaultLayout.vue

**Files:**
- Modify: `frontend/src/layouts/DefaultLayout.vue`

**Interfaces:**
- Consumes: `brand-*` / `accent-*` classes
- Produces: 新导航栏、页面宽度 `max-w-6xl`、统一 footer

- [ ] **Step 1: 重写 DefaultLayout.vue**

关键改动：
- 外层背景改为 `bg-gray-50 dark:bg-slate-950`
- header 改为 `bg-white/80 dark:bg-slate-900/80 backdrop-blur-sm border-b border-gray-200 dark:border-white/5`
- nav 容器 `max-w-6xl`
- Logo 方块改为 `from-brand-400 to-accent-500`
- 导航链接高亮改为 `text-brand-600 dark:text-brand-400 bg-brand-50 dark:bg-brand-900/20`
- main 容器 `max-w-6xl mx-auto px-4 py-10`
- footer 简洁居中

- [ ] **Step 2: 验证**

Run: `cd frontend && npm run build`
Expected: PASS

- [ ] **Step 3: 提交**

```bash
git add frontend/src/layouts/DefaultLayout.vue
git commit -m "feat(ui): redesign default layout with pink-purple navbar and footer"
```

---

## Task 4: 首页 HomeView.vue

**Files:**
- Modify: `frontend/src/views/HomeView.vue`

**Interfaces:**
- Consumes: `/api/quotes`, `PostCard`, `GradientButton`, `KawaiiIcon`
- Produces: 可滚动首页，包含 hero、精选文章网格、标签入口

- [ ] **Step 1: 移除固定定位**

删除所有 `<style>` 中的 `.home-*` 固定定位样式，恢复文档流。

- [ ] **Step 2: 新增文章获取逻辑**

```ts
import { usePostsStore } from '../stores/posts'
const postStore = usePostsStore()
onMounted(() => { postStore.fetchPosts(1, '') })
```

- [ ] **Step 3: 重写模板**

```vue
<template>
  <div class="space-y-16">
    <!-- Hero -->
    <section class="relative overflow-hidden rounded-3xl bg-gradient-to-br from-brand-100 via-white to-accent-100 dark:from-brand-900/30 dark:via-slate-900 dark:to-accent-900/30 p-8 md:p-14">
      <div class="relative z-10 flex flex-col md:flex-row items-center gap-10">
        <div class="flex-1 text-center md:text-left">
          <h1 class="text-4xl md:text-6xl font-black text-slate-800 dark:text-slate-100 leading-tight">
            My <span class="text-transparent bg-clip-text bg-gradient-to-r from-brand-400 to-accent-500">Blog</span>
          </h1>
          <p class="mt-4 text-lg text-slate-500 dark:text-slate-400">Learning in public, one post at a time.</p>
          <div class="mt-8 flex flex-wrap justify-center md:justify-start gap-3">
            <router-link to="/posts">
              <GradientButton>Explore posts</GradientButton>
            </router-link>
            <router-link to="/search"
              class="px-5 py-2.5 rounded-xl font-semibold text-slate-600 dark:text-slate-300
                     bg-white dark:bg-slate-800 border border-gray-200 dark:border-white/10
                     hover:border-brand-300 dark:hover:border-brand-700 transition-all">
              Search
            </router-link>
          </div>
        </div>
        <div class="w-40 h-40 md:w-56 md:h-56 shrink-0">
          <KawaiiIcon name="wave" />
        </div>
      </div>
    </section>

    <!-- Featured posts -->
    <section>
      <div class="flex items-center justify-between mb-6">
        <h2 class="text-2xl font-bold text-slate-800 dark:text-slate-100">Latest Posts</h2>
        <router-link to="/posts" class="text-sm font-medium text-brand-600 dark:text-brand-400 hover:underline">View all →</router-link>
      </div>
      <div v-if="postStore.loading" class="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div v-for="i in 6" :key="i" class="h-40 rounded-2xl bg-gray-200 dark:bg-white/5 animate-pulse"></div>
      </div>
      <div v-else-if="postStore.posts.length" class="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
        <PostCard v-for="post in postStore.posts.slice(0, 6)" :key="post.id" :post="post" />
      </div>
      <EmptyState v-else icon="sad" title="No posts yet" description="Upload your first markdown post in admin." />
    </section>
  </div>
</template>
```

- [ ] **Step 4: 验证构建**

Run: `cd frontend && npm run build`
Expected: PASS

- [ ] **Step 5: 提交**

```bash
git add frontend/src/views/HomeView.vue
git commit -m "feat(ui): redesign homepage with hero, card grid, and kawaii illustration"
```

---

## Task 5: 文章卡片 PostCard.vue

**Files:**
- Modify: `frontend/src/components/PostCard.vue`

**Interfaces:**
- Consumes: `Post` type, `.tag-pill`
- Produces: 新卡片样式，标签链接修复为 `/posts?tag=xxx`

- [ ] **Step 1: 重写 PostCard.vue**

```vue
<script setup lang="ts">
import type { Post } from '../stores/posts'
defineProps<{ post: Post }>()
</script>

<template>
  <article class="group flex flex-col h-full bg-white dark:bg-slate-900 rounded-2xl border border-gray-100 dark:border-white/5
                   shadow-sm hover:shadow-xl hover:shadow-brand-500/10
                   hover:-translate-y-1 transition-all duration-300 overflow-hidden">
    <!-- 封面占位 -->
    <div class="h-32 bg-gradient-to-br from-brand-100 to-accent-100 dark:from-brand-900/30 dark:to-accent-900/30"></div>
    <div class="flex flex-col flex-1 p-5">
      <div class="flex items-center gap-2 text-xs text-slate-400 dark:text-slate-500 mb-3">
        <time>{{ new Date(post.created_at).toLocaleDateString('zh-CN') }}</time>
        <span>·</span>
        <span>{{ post.tags.length }} tags</span>
      </div>
      <router-link :to="`/post/${post.slug}`" class="block flex-1">
        <h2 class="text-lg font-bold text-slate-800 dark:text-slate-100 group-hover:text-brand-600 dark:group-hover:text-brand-400 transition-colors line-clamp-2">
          {{ post.title }}
        </h2>
      </router-link>
      <div v-if="post.tags.length" class="flex gap-2 flex-wrap mt-4">
        <router-link v-for="tag in post.tags" :key="tag.id" :to="`/posts?tag=${tag.name}`" class="tag-pill">
          #{{ tag.name }}
        </router-link>
      </div>
    </div>
  </article>
</template>
```

- [ ] **Step 2: 验证构建**

Run: `cd frontend && npm run build`
Expected: PASS

- [ ] **Step 3: 提交**

```bash
git add frontend/src/components/PostCard.vue
git commit -m "feat(ui): redesign PostCard with card style and fix tag links"
```

---

## Task 6: 文章列表 PostsView.vue

**Files:**
- Modify: `frontend/src/views/PostsView.vue`

**Interfaces:**
- Consumes: `PostCard`, `allTags`, `store.posts`
- Produces: 卡片网格 + 标签筛选侧边栏 + 分页

- [ ] **Step 1: 移除 `<style>` 块和所有 CSS 变量**

全部用 Tailwind 替代。

- [ ] **Step 2: 重写模板**

```vue
<template>
  <div>
    <div class="flex flex-col lg:flex-row gap-10">
      <!-- Posts grid -->
      <div class="flex-1">
        <h1 class="text-3xl font-black text-slate-800 dark:text-slate-100 mb-6">
          {{ route.query.tag ? `#${route.query.tag}` : 'All Posts' }}
        </h1>

        <div v-if="store.loading" class="grid md:grid-cols-2 gap-6">
          <div v-for="i in 4" :key="i" class="h-48 rounded-2xl bg-gray-200 dark:bg-white/5 animate-pulse"></div>
        </div>

        <EmptyState v-else-if="!store.posts.length" icon="sad" title="No posts found" />

        <template v-else>
          <div class="grid md:grid-cols-2 gap-6">
            <PostCard v-for="post in store.posts" :key="post.id" :post="post" />
          </div>

          <div v-if="totalPages > 1" class="flex justify-center gap-2 mt-10">
            <button v-for="p in totalPages" :key="p" @click="$router.push({query:{...route.query,page:p}})"
              class="w-10 h-10 rounded-full text-sm font-semibold transition-all"
              :class="p === currentPage
                ? 'bg-gradient-to-r from-brand-400 to-accent-500 text-white shadow-md'
                : 'bg-white dark:bg-slate-900 text-slate-600 dark:text-slate-300 border border-gray-200 dark:border-white/10 hover:border-brand-300 dark:hover:border-brand-700'">
              {{ p }}
            </button>
          </div>
        </template>
      </div>

      <!-- Sidebar -->
      <aside class="lg:w-64 shrink-0">
        <div class="sticky top-24 bg-white dark:bg-slate-900 rounded-2xl border border-gray-100 dark:border-white/5 p-5 shadow-sm">
          <h3 class="text-xs font-bold text-slate-400 dark:text-slate-500 uppercase tracking-wider mb-4">Tags</h3>
          <div class="space-y-1">
            <button @click="$router.push('/posts')"
              class="w-full flex items-center justify-between px-3 py-2 rounded-xl text-sm transition-all"
              :class="!route.query.tag
                ? 'bg-brand-50 dark:bg-brand-900/20 text-brand-700 dark:text-brand-300 font-medium'
                : 'text-slate-600 dark:text-slate-300 hover:bg-gray-50 dark:hover:bg-white/5'">
              <span>All posts</span>
            </button>
            <button v-for="t in allTags" :key="t.name" @click="$router.push({path:'/posts',query:{tag:t.name}})"
              class="w-full flex items-center justify-between px-3 py-2 rounded-xl text-sm transition-all"
              :class="route.query.tag === t.name
                ? 'bg-brand-50 dark:bg-brand-900/20 text-brand-700 dark:text-brand-300 font-medium'
                : 'text-slate-600 dark:text-slate-300 hover:bg-gray-50 dark:hover:bg-white/5'">
              <span>#{{ t.name }}</span>
              <span class="text-xs px-2 py-0.5 rounded-full bg-gray-100 dark:bg-white/10 text-slate-400">{{ t.count }}</span>
            </button>
          </div>
        </div>
      </aside>
    </div>
  </div>
</template>
```

- [ ] **Step 3: 验证构建**

Run: `cd frontend && npm run build`
Expected: PASS

- [ ] **Step 4: 提交**

```bash
git add frontend/src/views/PostsView.vue
git commit -m "feat(ui): card grid layout for posts page, remove custom CSS"
```

---

## Task 7: 文章详情 + MarkdownRenderer + TOCSidebar

**Files:**
- Modify: `frontend/src/views/PostDetailView.vue`
- Modify: `frontend/src/components/MarkdownRenderer.vue`
- Modify: `frontend/src/components/TOCSidebar.vue`

**Interfaces:**
- Consumes: `post` data, `html`, `tocJson`
- Produces: 卡片化详情页、正常代码高亮、主题切换

- [ ] **Step 1: 修复 MarkdownRenderer.vue**

```ts
import { onMounted, watch, ref } from 'vue'
import hljs from 'highlight.js/lib/core'
import javascript from 'highlight.js/lib/languages/javascript'
import typescript from 'highlight.js/lib/languages/typescript'
import python from 'highlight.js/lib/languages/python'
import go from 'highlight.js/lib/languages/go'
import bash from 'highlight.js/lib/languages/bash'
import xml from 'highlight.js/lib/languages/xml'
import css from 'highlight.js/lib/languages/css'

hljs.registerLanguage('javascript', javascript)
hljs.registerLanguage('typescript', typescript)
hljs.registerLanguage('python', python)
hljs.registerLanguage('go', go)
hljs.registerLanguage('bash', bash)
hljs.registerLanguage('shell', bash)
hljs.registerLanguage('html', xml)
hljs.registerLanguage('css', css)

// 动态导入主题
const lightTheme = () => import('highlight.js/styles/github.css')
const darkTheme = () => import('highlight.js/styles/github-dark.css')

async function loadTheme() {
  const isDark = document.documentElement.classList.contains('dark')
  if (isDark) await darkTheme(); else await lightTheme()
}

onMounted(() => { loadTheme().then(highlight) })
watch(() => props.html, () => { setTimeout(() => { loadTheme().then(highlight) }, 0) })
```

- [ ] **Step 2: 重写 PostDetailView.vue**

将页面改为卡片化，标签链接改为 `/posts?tag=xxx`。

- [ ] **Step 3: 重写 TOCSidebar.vue**

改为圆角卡片，使用 brand 色高亮当前项（通过 IntersectionObserver 或保持现有样式）。

- [ ] **Step 4: 验证构建**

Run: `cd frontend && npm run build`
Expected: PASS

- [ ] **Step 5: 提交**

```bash
git add frontend/src/views/PostDetailView.vue frontend/src/components/MarkdownRenderer.vue frontend/src/components/TOCSidebar.vue
git commit -m "feat(ui): card-style post detail, fix syntax highlighting and theme"
```

---

## Task 8: 标签页 TagsView.vue

**Files:**
- Modify: `frontend/src/views/TagsView.vue`

**Interfaces:**
- Consumes: `tags` list
- Produces: 彩色胶囊标签卡片网格

- [ ] **Step 1: 重写为卡片网格**

```vue
<template>
  <div>
    <h1 class="text-3xl font-black text-slate-800 dark:text-slate-100 mb-1">Tags</h1>
    <p class="text-slate-400 dark:text-slate-500 mb-8 text-sm">{{ tags.length }} topics to explore</p>

    <div v-if="tags.length" class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
      <router-link v-for="t in tags" :key="t.name" :to="`/posts?tag=${t.name}`"
        class="group flex flex-col items-center justify-center p-6 rounded-2xl
               bg-white dark:bg-slate-900 border border-gray-100 dark:border-white/5
               hover:border-brand-300 dark:hover:border-brand-700
               hover:shadow-lg hover:shadow-brand-500/10 hover:-translate-y-1
               transition-all duration-300">
        <span class="text-2xl mb-2">🏷️</span>
        <span class="font-bold text-slate-700 dark:text-slate-200 group-hover:text-brand-600 dark:group-hover:text-brand-400 transition-colors">{{ t.name }}</span>
        <span class="text-xs text-slate-400 mt-1">{{ t.count }} posts</span>
      </router-link>
    </div>

    <EmptyState v-else icon="sad" title="No tags yet" />
  </div>
</template>
```

- [ ] **Step 2: 验证构建**

Run: `cd frontend && npm run build`
Expected: PASS

- [ ] **Step 3: 提交**

```bash
git add frontend/src/views/TagsView.vue
git commit -m "feat(ui): redesign tags page as card grid"
```

---

## Task 9: 时间线 TimelineView.vue

**Files:**
- Modify: `frontend/src/views/TimelineView.vue`

**Interfaces:**
- Consumes: `timeline` list
- Produces: 卡片时间轴布局

- [ ] **Step 1: 重写时间轴**

将 yearly node 改为粉紫渐变圆点，月度分组使用卡片，整体更圆润。

- [ ] **Step 2: 验证构建**

Run: `cd frontend && npm run build`
Expected: PASS

- [ ] **Step 3: 提交**

```bash
git add frontend/src/views/TimelineView.vue
git commit -m "feat(ui): redesign timeline with card style and gradient nodes"
```

---

## Task 10: 搜索页 SearchView.vue

**Files:**
- Modify: `frontend/src/views/SearchView.vue`

**Interfaces:**
- Consumes: `PostCard`, `EmptyState`
- Produces: 居中搜索 + 结果卡片

- [ ] **Step 1: 重写搜索页**

- 搜索框居中、更大
- 结果使用 `PostCard` 网格
- 无结果使用 `EmptyState icon="search"`

- [ ] **Step 2: 验证构建**

Run: `cd frontend && npm run build`
Expected: PASS

- [ ] **Step 3: 提交**

```bash
git add frontend/src/views/SearchView.vue
git commit -m "feat(ui): redesign search page with centered input and card results"
```

---

## Task 11: 登录页 LoginView.vue

**Files:**
- Modify: `frontend/src/views/LoginView.vue`

**Interfaces:**
- Consumes: `useAuthStore`
- Produces: 粉紫渐变居中卡片

- [ ] **Step 1: 更新 Logo 和按钮颜色**

- Logo 方块改为 `from-brand-400 to-accent-500`
- 标题颜色改为 slate 体系
- 表单边框改为 `border-gray-100 dark:border-white/5`
- 登录按钮改为 `GradientButton` 或 `.btn-primary`
- 输入框 focus ring 改为 brand

- [ ] **Step 2: 验证构建**

Run: `cd frontend && npm run build`
Expected: PASS

- [ ] **Step 3: 提交**

```bash
git add frontend/src/views/LoginView.vue
git commit -m "feat(ui): redesign login page with pink-purple card"
```

---

## Task 12: 后台 AdminView.vue + UploadZone.vue

**Files:**
- Modify: `frontend/src/views/AdminView.vue`
- Modify: `frontend/src/components/UploadZone.vue`

**Interfaces:**
- Consumes: `GradientButton`
- Produces: 卡片化后台、Tailwind 上传区

- [ ] **Step 1: 重写 UploadZone.vue**

移除所有 `<style>`，使用 Tailwind：
- 外层 `border-2 border-dashed border-gray-300 dark:border-white/10 rounded-2xl p-10 text-center bg-white dark:bg-slate-900 transition-all`
- drag 态 `border-brand-400 bg-brand-50 dark:bg-brand-900/20 scale-[1.01]`
- 图标用 KawaiiIcon 或统一 SVG
- spinner 用 `animate-spin border-brand-500`

- [ ] **Step 2: 更新 AdminView.vue**

- 各区域改为圆角卡片
- 按钮统一为 `.btn-primary` 或 `GradientButton`
- 表格保持，但容器卡片化

- [ ] **Step 3: 验证构建**

Run: `cd frontend && npm run build`
Expected: PASS

- [ ] **Step 4: 提交**

```bash
git add frontend/src/views/AdminView.vue frontend/src/components/UploadZone.vue
git commit -m "feat(ui): redesign admin and upload zone with Tailwind card style"
```

---

## Task 13: ThemeToggle.vue 统一

**Files:**
- Modify: `frontend/src/components/ThemeToggle.vue`

**Interfaces:**
- Produces: 统一太阳/月亮 SVG，hover 背景融入新主题

- [ ] **Step 1: 更新样式**

- hover 背景改为 `hover:bg-brand-100 dark:hover:bg-white/10`
- 图标颜色 light 模式下也使用 brand/accent

- [ ] **Step 2: 验证构建**

Run: `cd frontend && npm run build`
Expected: PASS

- [ ] **Step 3: 提交**

```bash
git add frontend/src/components/ThemeToggle.vue
git commit -m "feat(ui): update theme toggle colors"
```

---

## Task 14: 最终构建验证

**Files:**
- All changed files

- [ ] **Step 1: 完整构建**

Run: `cd frontend && npm run build`
Expected: PASS with no TS errors

- [ ] **Step 2: 类型检查**

Run: `cd frontend && npx vue-tsc --noEmit`
Expected: PASS

- [ ] **Step 3: 提交最终变更（如果未提交）**

```bash
git add -A
git commit -m "feat(ui): complete pink-purple kawaii redesign"
```

---

## Spec Coverage Check

| 设计文档要求 | 对应 Task |
|---|---|
| 粉紫渐变配色 | Task 1 |
| 全局布局 max-w-6xl | Task 3 |
| 首页卡片化 hero | Task 4 |
| 文章列表卡片网格 | Task 5, 6 |
| 文章详情卡片化 + TOC | Task 7 |
| Markdown 代码高亮修复 | Task 7 |
| 标签页卡片化 | Task 8 |
| 时间线卡片化 | Task 9 |
| 搜索页卡片化 + EmptyState | Task 10 |
| 登录页粉紫卡片 | Task 11 |
| 后台卡片化 + 上传区重写 | Task 12 |
| 二次元装饰插画 | Task 2, 4, 10 |
| 标签链接修复 | Task 5, 7 |

## Placeholder Scan

- 无 TBD/TODO
- KawaiiIcon.vue 中需为每个 `name` 补全 SVG（执行时按风格统一实现）
- 所有验证命令已给出
