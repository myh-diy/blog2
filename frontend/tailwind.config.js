import typography from '@tailwindcss/typography'

function cssVarColor(name) {
  return `rgb(var(${name}) / <alpha-value>)`
}

/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        brand: {
          50:  cssVarColor('--brand-50'),
          100: cssVarColor('--brand-100'),
          200: cssVarColor('--brand-200'),
          300: cssVarColor('--brand-300'),
          400: cssVarColor('--brand-400'),
          500: cssVarColor('--brand-500'),
          600: cssVarColor('--brand-600'),
          700: cssVarColor('--brand-700'),
          800: cssVarColor('--brand-800'),
          900: cssVarColor('--brand-900'),
        },
        accent: {
          50:  cssVarColor('--accent-50'),
          100: cssVarColor('--accent-100'),
          200: cssVarColor('--accent-200'),
          300: cssVarColor('--accent-300'),
          400: cssVarColor('--accent-400'),
          500: cssVarColor('--accent-500'),
          600: cssVarColor('--accent-600'),
          700: cssVarColor('--accent-700'),
          800: cssVarColor('--accent-800'),
          900: cssVarColor('--accent-900'),
        },
        // Warm neutrals (keep for compatibility)
        warm: {
          50:  '#fafaf9',
          100: '#f5f5f4',
          200: '#e7e5e4',
          300: '#d6d3d1',
          400: '#a8a29e',
          500: '#78716c',
          600: '#57534e',
          700: '#44403c',
          800: '#292524',
          900: '#1c1917',
        },
      },
      fontFamily: {
        sans: ['"Noto Sans SC"', '"Inter"', 'system-ui', '-apple-system', 'sans-serif'],
        mono: ['"JetBrains Mono"', '"Fira Code"', 'monospace'],
      },
    },
  },
  plugins: [
    typography,
  ],
}
