# 前端界面 B站少女粉紫风格重设计

## 背景

项目前端使用 Vue 3 + TypeScript + Vite + Tailwind CSS，当前视觉风格不统一：
- 主色调为暖橙/玫红，但 favicon 为紫色，图标风格不一致
- 首页使用固定定位破坏文档流
- 文章列表页混入大量原生 CSS，与 Tailwind 体系脱节
- 代码高亮未注册语言且主题不随 dark mode 切换
- 上传组件视觉突兀

本次设计目标：在不引入第三方 UI 组件库的前提下，将整体风格改造为**活泼、卡片化、粉紫渐变、带二次元装饰插画**的 B站少女感界面。

## 设计目标

1. 统一配色与品牌感：粉紫渐变主色贯穿全站
2. 卡片化布局：文章列表、标签、搜索等页面改为卡片流
3. 圆润可爱：大圆角、柔和阴影、hover 动效
4. 二次元点缀：在空状态、hero 区、分类处使用统一风格 SVG 插画
5. 修复既有视觉/交互 bug

## 配色系统

扩展 `tailwind.config.js`：

```js
colors: {
  brand: {
    50:  '#fff1f8',
    100: '#ffe4f3',
    200: '#ffc9e9',
    300: '#ff9dd4',
    400: '#ff5fb8',
    500: '#fb7299',   // B站粉
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
    400: '#a78bfa',   // 薰衣草紫
    500: '#8b5cf6',
    600: '#7c3aed',
    700: '#6d28d9',
  }
}
```

- 浅色背景：`bg-gray-50`
- 深色背景：`bg-slate-950`
- 卡片背景：`bg-white` / `bg-slate-900`
- 主按钮：`bg-gradient-to-r from-brand-400 to-accent-500`
- 标签：`bg-brand-100 text-brand-700`
- 链接 hover：`text-brand-500`

## 字体

继续使用现有 `Inter` + `Noto Sans SC`，标题使用适中字重，不引入新字体。

## 布局改造

### 全局布局 `DefaultLayout.vue`

- 导航栏：`bg-white/80 dark:bg-slate-900/80 backdrop-blur-sm`，底部 1px 边框
- 页面最大宽度：`max-w-6xl`（原 `max-w-4xl`）
- Logo：粉紫渐变圆角方块 + 字母 "B"
- Footer：简洁居中，与导航栏风格统一

### 首页 `/`

移除全屏固定定位，恢复正常文档流：
1. Hero 区：左侧欢迎文案 + 右侧二次元装饰插画，粉紫渐变背景
2. 精选文章卡片网格：桌面 3 列，平板 2 列，手机 1 列
3. 热门标签快捷入口
4. 页脚

### 文章列表 `/posts`

- 顶部：标题 + 搜索框 + 标签筛选
- 内容区：文章卡片网格（桌面 2 列，移动 1 列）
- 分页：圆润胶囊按钮

### 文章详情 `/post/:slug`

- 顶部：标题、标签、日期
- 主体：正文卡片 + 右侧 TOC 卡片
- 代码高亮主题随主题切换

### 时间线 `/timeline`

- 卡片时间轴布局
- 时间节点使用粉紫渐变圆点
-  yearly/monthly 分组清晰

### 标签页 `/tags`

- 彩色胶囊标签卡片网格
- hover 上浮 + 缩放

### 搜索 `/search`

- 搜索框居中放大
- 结果卡片化
- 无结果时显示 `EmptyState` 二次元插画

### 登录 `/login`

- 居中卡片，粉紫渐变边框与按钮

### 后台 `/admin`

- 上传区改用 Tailwind 重写
- 列表、表单区域卡片化

## 组件改造

| 组件 | 改造内容 |
|---|---|
| `DefaultLayout.vue` | 新导航栏、粉紫 logo、footer 统一 |
| `PostCard.vue` | 大圆角卡片、hover 上浮、标签胶囊、封面占位 |
| `MarkdownRenderer.vue` | 代码高亮注册语言、主题切换、引用块粉紫边框 |
| `TOCSidebar.vue` | 圆角卡片、当前项高亮 |
| `ThemeToggle.vue` | 统一太阳/月亮 SVG 图标 |
| `UploadZone.vue` | 完全用 Tailwind 重写，拖拽态粉紫边框 |
| `EmptyState.vue`（新增） | 空状态二次元插画组件 |
| `GradientButton.vue`（新增） | 粉紫渐变按钮 |
| `KawaiiIcon.vue`（新增） | 可复用二次元风格 SVG 插画集合 |

## 动画与交互

- 卡片 hover：`transition-all duration-300 hover:-translate-y-1 hover:shadow-lg`
- 按钮 hover：`hover:scale-105 transition-transform`
- 页面切换：淡入上移动效
- 主题切换：颜色平滑过渡
- 骨架屏：粉紫 shimmer（可选，视实现成本）

## 二次元装饰插画

使用内联 SVG，风格统一：圆润线条、粉紫配色、无复杂细节，避免版权风险。

- 首页 hero：挥手角色
- 搜索空状态：趴着的小人
- 错误/404：困惑小人
- 上传成功：开心举手小人
- 分类图标：星星、爱心、气泡等装饰

## Bug 修复

1. 标签链接从 `/?tag=xxx` 改为 `/posts?tag=xxx`
2. `MarkdownRenderer.vue` 注册常用语言：javascript, typescript, python, go, bash, html, css, vue, json, yaml, markdown
3. 代码高亮主题根据 dark mode 切换：light 用 `github`，dark 用 `github-dark`
4. `HomeView.vue` 移除固定定位
5. `PostsView.vue` 移除原生 `:root` CSS 变量，改用 Tailwind

## 非目标

- 不引入第三方 Vue UI 组件库
- 不添加后端接口改动
- 不添加新的业务功能（评论、点赞等）
- 不替换字体或引入外部插画资源

## 成功标准

1. 所有页面视觉风格统一为粉紫渐变少女风
2. 卡片化布局在桌面和移动端均正常显示
3. 代码高亮正常工作且主题随系统主题切换
4. 原有 bug 已修复
5. 构建通过，无新增 TypeScript 错误
