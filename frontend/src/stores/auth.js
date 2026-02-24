import { defineStore } from 'pinia'
import axios from 'axios'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null
  }),
  actions: {
    async login(username, password) {
      const res = await axios.post('/login', { username, password })
      localStorage.setItem('token', res.data.token)
      await this.fetchProfile()
    },
    async register(data) {
      await axios.post('/register', data)
    },
    async fetchProfile() {
      const res = await axios.get('/profile')
      this.user = res.data
    },
    logout() {
      localStorage.removeItem('token')
      this.user = null
      window.location = '/'
    }
  }
})