import { ref, onMounted } from 'vue'

const STORAGE_KEY = 'blog-site-title'

const siteTitle = ref('My Blog')

export function useSiteTitle() {
  onMounted(() => {
    const saved = localStorage.getItem(STORAGE_KEY)
    if (saved) siteTitle.value = saved
  })

  function setSiteTitle(title: string) {
    siteTitle.value = title.trim() || 'My Blog'
    localStorage.setItem(STORAGE_KEY, siteTitle.value)
    document.title = siteTitle.value
  }

  return {
    siteTitle,
    setSiteTitle,
  }
}
