import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

const backendUrl = process.env.VITE_BACKEND_URL || 'http://localhost:8080'

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue()],
  server: {
    proxy: {
      '/api': {
        target: backendUrl,
        changeOrigin: true,
      },
      '/uploads': {
        target: backendUrl,
        changeOrigin: true,
      },
    },
  },
})
