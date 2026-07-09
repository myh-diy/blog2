export interface RGB { r: number; g: number; b: number }
export interface HSL { h: number; s: number; l: number }

export function hexToRgb(hex: string): RGB {
  const clean = hex.replace('#', '')
  const bigint = parseInt(clean, 16)
  return {
    r: (bigint >> 16) & 255,
    g: (bigint >> 8) & 255,
    b: bigint & 255,
  }
}

export function rgbToHex({ r, g, b }: RGB): string {
  return '#' + [r, g, b].map(v => Math.round(v).toString(16).padStart(2, '0')).join('')
}

export function rgbToHsl({ r, g, b }: RGB): HSL {
  r /= 255; g /= 255; b /= 255
  const max = Math.max(r, g, b)
  const min = Math.min(r, g, b)
  let h = 0, s = 0, l = (max + min) / 2

  if (max !== min) {
    const d = max - min
    s = l > 0.5 ? d / (2 - max - min) : d / (max + min)
    switch (max) {
      case r: h = (g - b) / d + (g < b ? 6 : 0); break
      case g: h = (b - r) / d + 2; break
      case b: h = (r - g) / d + 4; break
    }
    h /= 6
  }

  return { h: h * 360, s: s * 100, l: l * 100 }
}

export function hslToRgb({ h, s, l }: HSL): RGB {
  s /= 100; l /= 100
  const c = (1 - Math.abs(2 * l - 1)) * s
  const x = c * (1 - Math.abs(((h / 60) % 2) - 1))
  const m = l - c / 2
  let r = 0, g = 0, b = 0

  if (h < 60) { r = c; g = x; b = 0 }
  else if (h < 120) { r = x; g = c; b = 0 }
  else if (h < 180) { r = 0; g = c; b = x }
  else if (h < 240) { r = 0; g = x; b = c }
  else if (h < 300) { r = x; g = 0; b = c }
  else { r = c; g = 0; b = x }

  return {
    r: Math.round((r + m) * 255),
    g: Math.round((g + m) * 255),
    b: Math.round((b + m) * 255),
  }
}

export function hexToHsl(hex: string): HSL {
  return rgbToHsl(hexToRgb(hex))
}

export function hslToHex(hsl: HSL): string {
  return rgbToHex(hslToRgb(hsl))
}

export function clamp(num: number, min: number, max: number): number {
  return Math.max(min, Math.min(max, num))
}

export function generateScale(baseHex: string): Record<number, string> {
  const base = hexToHsl(baseHex)
  const steps: Record<number, number> = {
    50: 48,
    100: 40,
    200: 30,
    300: 18,
    400: 8,
    500: 0,
    600: -10,
    700: -18,
    800: -26,
    900: -34,
  }

  const scale: Record<number, string> = {}
  for (const [level, delta] of Object.entries(steps)) {
    const l = clamp(base.l + delta, 5, 98)
    const s = level === '500' ? base.s : clamp(base.s + (delta < 0 ? 5 : -5), 10, 100)
    scale[Number(level)] = hslToHex({ h: base.h, s, l })
  }
  return scale
}

export function generateCSSVars(brandHex: string, accentHex: string): Record<string, string> {
  const brand = generateScale(brandHex)
  const accent = generateScale(accentHex)
  const vars: Record<string, string> = {}

  for (const [level, hex] of Object.entries(brand)) {
    const rgb = hexToRgb(hex)
    vars[`--brand-${level}`] = `${rgb.r} ${rgb.g} ${rgb.b}`
  }
  for (const [level, hex] of Object.entries(accent)) {
    const rgb = hexToRgb(hex)
    vars[`--accent-${level}`] = `${rgb.r} ${rgb.g} ${rgb.b}`
  }

  return vars
}
