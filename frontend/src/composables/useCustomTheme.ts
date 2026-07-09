import { ref, watch } from 'vue'
import { generateCSSVars } from '../utils/color'

export interface ThemeColors {
  brand: string
  accent: string
}

const STORAGE_KEY = 'custom-theme'

const defaultColors: ThemeColors = {
  brand: '#22c55e',
  accent: '#14b8a6',
}

function loadColors(): ThemeColors {
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    if (raw) {
      const parsed = JSON.parse(raw)
      if (parsed.brand && parsed.accent) return parsed
    }
  } catch {}
  return { ...defaultColors }
}

const colors = ref<ThemeColors>(loadColors())

export function useCustomTheme() {
  function apply() {
    const vars = generateCSSVars(colors.value.brand, colors.value.accent)
    const root = document.documentElement
    for (const [key, value] of Object.entries(vars)) {
      root.style.setProperty(key, value)
    }
  }

  function setBrand(hex: string) {
    colors.value.brand = hex
  }

  function setAccent(hex: string) {
    colors.value.accent = hex
  }

  function reset() {
    colors.value = { ...defaultColors }
  }

  watch(colors, () => {
    apply()
    localStorage.setItem(STORAGE_KEY, JSON.stringify(colors.value))
  }, { deep: true })

  return {
    colors,
    apply,
    setBrand,
    setAccent,
    reset,
  }
}
