import { computed, ref } from 'vue'

export interface ClassicTheme {
  id: string
  name: string
  description: string
  brand?: string
  accent?: string
  background?: string
  position?: string
  overlay?: number
}

export const classicThemes: ClassicTheme[] = [
  {
    id: 'aemeath',
    name: '爱弥斯',
    description: '星火与霓光',
    brand: '#e95f9f',
    accent: '#22b8d6',
    background: '/themes/aemeath.png',
    position: 'center center',
    overlay: 0.74,
  },
  {
    id: 'shorekeeper',
    name: '守岸人',
    description: '星海与蓝蝶',
    brand: '#3478e5',
    accent: '#35cce8',
    background: '/themes/shorekeeper.jpg',
    position: 'center center',
    overlay: 0.7,
  },
  {
    id: 'qingxiao',
    name: '清霄',
    description: '云剑与清辉',
    brand: '#356fb6',
    accent: '#39bcae',
    background: '/themes/qingxiao.jpg',
    position: 'center center',
    overlay: 0.74,
  },
  {
    id: 'custom',
    name: '自定义',
    description: '使用后台背景',
  },
]

const STORAGE_KEY = 'classic-theme'
const DEFAULT_THEME = 'aemeath'
const storedTheme = localStorage.getItem(STORAGE_KEY)
const activeThemeId = ref(classicThemes.some(theme => theme.id === storedTheme) ? storedTheme! : DEFAULT_THEME)
const currentTheme = computed(() => classicThemes.find(theme => theme.id === activeThemeId.value) ?? classicThemes[0])

export function useClassicTheme() {
  function setClassicTheme(id: string) {
    if (!classicThemes.some(theme => theme.id === id)) return
    activeThemeId.value = id
    localStorage.setItem(STORAGE_KEY, id)
  }

  return { classicThemes, activeThemeId, currentTheme, setClassicTheme }
}
