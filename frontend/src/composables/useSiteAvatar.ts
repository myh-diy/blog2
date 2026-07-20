import { ref, onMounted } from 'vue'
import api from '../utils/api'

const STORAGE_KEY = 'blog-site-avatar'

const DEFAULT_AVATAR = '/go-gopher.svg'

// Kept only to migrate installations that saved the previous placeholder.
const LEGACY_DEFAULT_AVATAR =
  'data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHZpZXdCb3g9IjAgMCAxMDAgMTAwIj48cmVjdCB3aWR0aD0iMTAwIiBoZWlnaHQ9IjEwMCIgZmlsbD0iI2UwZjdmYSIvPjxlbGxpcHNlIGN4PSI1MCIgY3k9IjY4IiByeD0iMzAiIHJ5PSIyMiIgZmlsbD0iIzgwZGVlYSIvPjxlbGxpcHNlIGN4PSI1MCIgY3k9IjQyIiByeD0iMjQiIHJ5PSIyMiIgZmlsbD0iIzgwZGVlYSIvPjxlbGxpcHNlIGN4PSIzMCIgY3k9IjI4IiByeD0iNyIgcnk9IjE2IiBmaWxsPSIjODBkZWVhIi8+PGVsbGlwc2UgY3g9IjcwIiBjeT0iMjgiIHJ4PSI3IiByeT0iMTYiIGZpbGw9IiM4MGRlZWEiLz48Y2lyY2xlIGN4PSIzNiIgY3k9IjQwIiByPSI4IiBmaWxsPSJ3aGl0ZSIvPjxjaXJjbGUgY3g9IjY0IiBjeT0iNDAiIHI9IjgiIGZpbGw9IndoaXRlIi8+PGNpcmNsZSBjeD0iMzciIGN5PSI0MSIgcj0iMy41IiBmaWxsPSIjMzMzIi8+PGNpcmNsZSBjeD0iNjMiIGN5PSI0MSIgcj0iMy41IiBmaWxsPSIjMzMzIi8+PGNpcmNsZSBjeD0iMzUiIGN5PSIzOCIgcj0iMS44IiBmaWxsPSJ3aGl0ZSIvPjxjaXJjbGUgY3g9IjYxIiBjeT0iMzgiIHI9IjEuOCIgZmlsbD0id2hpdGUiLz48cmVjdCB4PSI0MSIgeT0iNTYiIHdpZHRoPSI4IiBoZWlnaHQ9IjExIiByeD0iMS41IiBmaWxsPSJ3aGl0ZSIvPjxyZWN0IHg9IjUxIiB5PSI1NiIgd2lkdGg9IjgiIGhlaWdodD0iMTEiIHJ4PSIxLjUiIGZpbGw9IndoaXRlIi8+PC9zdmc+'

const siteAvatar = ref(DEFAULT_AVATAR)
const isDefaultAvatar = ref(true)

async function saveAvatar(url: string) {
  try {
    await api.put('/admin/settings', { site_avatar: url })
  } catch {
    // Offline or not admin; localStorage is already updated as cache
  }
}

export function useSiteAvatar() {
  onMounted(async () => {
    try {
      const r = await api.get('/settings')
      const s = r.data.settings as Record<string, string>
      if (s.site_avatar !== undefined) {
        const avatar = s.site_avatar === LEGACY_DEFAULT_AVATAR ? DEFAULT_AVATAR : s.site_avatar
        siteAvatar.value = avatar
        isDefaultAvatar.value = avatar === DEFAULT_AVATAR
        localStorage.setItem(STORAGE_KEY, avatar)
      } else {
        loadFromLocal()
      }
    } catch {
      loadFromLocal()
    }
  })

  function loadFromLocal() {
    const saved = localStorage.getItem(STORAGE_KEY)
    if (saved) {
      const avatar = saved === LEGACY_DEFAULT_AVATAR ? DEFAULT_AVATAR : saved
      siteAvatar.value = avatar
      isDefaultAvatar.value = avatar === DEFAULT_AVATAR
    }
  }

  function setSiteAvatar(url: string) {
    siteAvatar.value = url
    isDefaultAvatar.value = false
    localStorage.setItem(STORAGE_KEY, url)
    saveAvatar(url)
  }

  function resetSiteAvatar() {
    siteAvatar.value = DEFAULT_AVATAR
    isDefaultAvatar.value = true
    localStorage.removeItem(STORAGE_KEY)
    saveAvatar(DEFAULT_AVATAR)
  }

  return {
    siteAvatar,
    isDefaultAvatar,
    setSiteAvatar,
    resetSiteAvatar,
  }
}
