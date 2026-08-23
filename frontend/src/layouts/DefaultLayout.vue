<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { BookOpen, LogIn, PenLine, Search } from '@lucide/vue'
import ThemeToggle from '../components/ThemeToggle.vue'
import DisplayModeToggle from '../components/DisplayModeToggle.vue'
import ClassicThemeSelector from '../components/ClassicThemeSelector.vue'
import PageDecorations from '../components/PageDecorations.vue'
import { useBackgroundImage } from '../composables/useBackgroundImage'
import { useDisplayMode } from '../composables/useDisplayMode'
import { useClassicTheme } from '../composables/useClassicTheme'
import { useSiteTitle } from '../composables/useSiteTitle'
import { useAuthStore } from '../stores/auth'
import { generateCSSVars } from '../utils/color'
import api from '../utils/api'

const route = useRoute()
const { backgroundImage, bgOpacity } = useBackgroundImage()
const { siteTitle } = useSiteTitle()
const { displayMode } = useDisplayMode()
const { currentTheme } = useClassicTheme()
const isMinimal = computed(() => displayMode.value === 'minimal')
const classicBackground = computed(() => currentTheme.value.background || backgroundImage.value)
const classicOverlayOpacity = computed(() => currentTheme.value.background ? currentTheme.value.overlay : bgOpacity.value)
const layoutThemeStyle = computed(() => {
  if (isMinimal.value || !currentTheme.value.brand || !currentTheme.value.accent) return undefined
  return generateCSSVars(currentTheme.value.brand, currentTheme.value.accent)
})
const auth = useAuthStore()
const tags = ref<{ name: string; count: number }[]>([])

const publicNavLinks = [
  { to: '/', label: '首页' },
  { to: '/posts', label: '文章' },
  { to: '/timeline', label: '归档' },
  { to: '/tags', label: '标签' },
  { to: '/monitor', label: '监控' },
]
const navLinks = computed(() => auth.isAuthenticated
  ? [...publicNavLinks, { to: '/services', label: '私人服务' }]
  : publicNavLinks)

onMounted(async () => {
  try {
    const response = await api.get('/tags')
    tags.value = (response.data.tags || []).slice(0, 8)
  } catch {}
})
</script>

<template>
  <div class="relative min-h-screen flex flex-col" :style="layoutThemeStyle">
    <!-- Global background image layer -->
    <div
      class="fixed inset-0 z-0 bg-cover bg-center bg-fixed bg-no-repeat transition-all duration-700"
      :style="{
        backgroundImage: !isMinimal && classicBackground ? `url('${classicBackground}')` : 'none',
        backgroundPosition: currentTheme.position || 'center center',
      }"
      aria-hidden="true"
    >
      <div
        class="absolute inset-0 transition-colors"
        :class="!isMinimal && classicBackground ? 'bg-white dark:bg-slate-950' : 'bg-gray-50 dark:bg-slate-950'"
        :style="!isMinimal && classicBackground ? { opacity: classicOverlayOpacity } : undefined"
      ></div>
    </div>
    <!-- Navbar -->
    <header class="sticky top-0 z-50 border-b border-slate-200 bg-white/95 backdrop-blur dark:border-white/10 dark:bg-[#171b22]/95">
      <nav class="mx-auto flex h-14 max-w-[1200px] items-center gap-3 px-4">
        <router-link to="/" class="flex min-w-0 shrink-0 items-center gap-2 text-slate-900 dark:text-white">
          <span class="flex h-8 w-8 items-center justify-center bg-brand-500 text-white"><BookOpen :size="18" /></span>
          <span class="hidden max-w-32 truncate text-lg font-bold sm:block">{{ siteTitle }}</span>
        </router-link>

        <div class="hidden h-full items-center md:flex">
          <router-link v-for="link in navLinks" :key="link.to" :to="link.to"
            class="flex h-full items-center border-b-2 px-3 text-sm transition-colors"
            :class="route.path === link.to ? 'border-brand-500 text-brand-600' : 'border-transparent text-slate-500 hover:text-slate-900 dark:text-slate-400 dark:hover:text-white'"
          >
            {{ link.label }}
          </router-link>
        </div>

        <router-link to="/search" class="ml-auto hidden h-9 w-52 items-center justify-between border border-slate-200 bg-slate-50 px-3 text-sm text-slate-400 transition-colors hover:border-brand-400 lg:flex dark:border-white/10 dark:bg-[#11151b]">
          <span>搜索文章</span><Search :size="16" />
        </router-link>

        <div class="ml-auto flex items-center gap-1 lg:ml-0">
          <DisplayModeToggle />
          <ThemeToggle />
          <router-link v-if="auth.isAuthenticated" to="/admin" class="btn-primary h-9 px-3" title="进入管理后台">
            <PenLine :size="16" /><span class="hidden sm:inline">创作中心</span>
          </router-link>
          <router-link v-else to="/login" class="btn-secondary h-9 px-3" title="管理员登录">
            <LogIn :size="16" /><span class="hidden sm:inline">登录</span>
          </router-link>
        </div>
      </nav>
    </header>

    <div class="relative z-10 border-b border-slate-200 bg-white dark:border-white/10 dark:bg-[#171b22]">
      <div class="nav-scroll mx-auto flex h-11 max-w-[1200px] items-center gap-1 overflow-x-auto px-4">
        <router-link v-for="link in navLinks" :key="`mobile-${link.to}`" :to="link.to" class="shrink-0 px-3 py-2 text-sm md:hidden"
          :class="route.path === link.to ? 'text-brand-600' : 'text-slate-500 dark:text-slate-400'">{{ link.label }}</router-link>
        <span class="hidden h-4 w-px bg-slate-200 md:block dark:bg-white/10"></span>
        <router-link v-for="tag in tags" :key="tag.name" :to="`/posts?tag=${encodeURIComponent(tag.name)}`"
          class="shrink-0 px-3 py-2 text-sm text-slate-500 transition-colors hover:text-brand-600 dark:text-slate-400 dark:hover:text-brand-300">
          {{ tag.name }}
        </router-link>
        <router-link to="/search" class="ml-auto shrink-0 px-3 py-2 text-sm text-slate-500 lg:hidden dark:text-slate-400">
          搜索
        </router-link>
      </div>
    </div>

    <ClassicThemeSelector v-if="!isMinimal" />

    <main class="relative z-10 mx-auto w-full max-w-[1200px] flex-1 px-3 py-4 sm:px-4 sm:py-5">
      <PageDecorations v-if="!isMinimal" />
      <slot />
    </main>

    <footer class="relative z-10 mt-8 border-t border-slate-200 bg-white dark:border-white/10 dark:bg-[#171b22]">
      <div class="mx-auto flex max-w-[1200px] flex-col items-center justify-between gap-2 px-4 py-6 text-xs text-slate-400 sm:flex-row dark:text-slate-500">
        <p>© {{ new Date().getFullYear() }} {{ siteTitle }} · 记录技术与思考</p>
        <p>Powered by Go + Vue</p>
      </div>
    </footer>
  </div>
</template>
