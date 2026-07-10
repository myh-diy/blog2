import { ref, onMounted } from 'vue'
import api from '../utils/api'

const STORAGE_KEY = 'blog-site-title'

const siteTitle = ref('My Blog')

async function saveTitle(title: string) {
  try {
    await api.put('/admin/settings', { site_title: title })
  } catch {
    // Offline or not admin; localStorage is already updated as cache
  }
}

export function useSiteTitle() {
  onMounted(async () => {
    try {
      const r = await api.get('/settings')
      const s = r.data.settings as Record<string, string>
      if (s.site_title !== undefined) {
        siteTitle.value = s.site_title
        localStorage.setItem(STORAGE_KEY, s.site_title)
      } else {
        siteTitle.value = localStorage.getItem(STORAGE_KEY) || 'My Blog'
      }
    } catch {
      siteTitle.value = localStorage.getItem(STORAGE_KEY) || 'My Blog'
    }
    document.title = siteTitle.value
  })

  function setSiteTitle(title: string) {
    siteTitle.value = title.trim() || 'My Blog'
    localStorage.setItem(STORAGE_KEY, siteTitle.value)
    document.title = siteTitle.value
    saveTitle(siteTitle.value)
  }

  return {
    siteTitle,
    setSiteTitle,
  }
}
