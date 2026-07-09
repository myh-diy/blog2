import { ref, watch } from 'vue'
import { generateCSSVars, hexToRgb } from '../utils/color'

export interface ThemePreset {
  name: string
  brand: string
  accent: string
}

export const presets: ThemePreset[] = [
  { name: 'Spring', brand: '#22c55e', accent: '#14b8a6' },
  { name: 'Sakura', brand: '#fb7299', accent: '#f472b6' },
  { name: 'Sky',    brand: '#3b82f6', accent: '#06b6d4' },
  { name: 'Lavender', brand: '#8b5cf6', accent: '#a855f7' },
]

const STORAGE_KEY = 'theme-preset'

function findPreset(name: string): ThemePreset {
  return presets.find(p => p.name === name) || presets[0]
}

function loadPresetName(): string {
  try {
    return localStorage.getItem(STORAGE_KEY) || 'Spring'
  } catch {}
  return 'Spring'
}

const activePreset = ref<string>(loadPresetName())

export function useCustomTheme() {
  function apply() {
    const preset = findPreset(activePreset.value)
    const vars = generateCSSVars(preset.brand, preset.accent)
    const root = document.documentElement
    for (const [key, value] of Object.entries(vars)) {
      root.style.setProperty(key, value)
    }
    // Also set a hex reference for places that need a solid color
    const brandRgb = hexToRgb(preset.brand)
    const accentRgb = hexToRgb(preset.accent)
    root.style.setProperty('--brand-hex', preset.brand)
    root.style.setProperty('--accent-hex', preset.accent)
    root.style.setProperty('--brand-rgb', `${brandRgb.r} ${brandRgb.g} ${brandRgb.b}`)
    root.style.setProperty('--accent-rgb', `${accentRgb.r} ${accentRgb.g} ${accentRgb.b}`)
  }

  function setPreset(name: string) {
    activePreset.value = name
  }

  watch(activePreset, () => {
    apply()
    localStorage.setItem(STORAGE_KEY, activePreset.value)
  })

  return {
    presets,
    activePreset,
    apply,
    setPreset,
  }
}
