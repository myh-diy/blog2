import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../utils/api'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('token') || '')
  const isAuthenticated = ref(!!token.value)

  async function login(username: string, password: string) {
    const res = await api.post('/auth/login', { username, password })
    token.value = res.data.token
    localStorage.setItem('token', res.data.token)
    isAuthenticated.value = true
  }

  function logout() {
    token.value = ''
    localStorage.removeItem('token')
    isAuthenticated.value = false
  }

  return { token, isAuthenticated, login, logout }
})
