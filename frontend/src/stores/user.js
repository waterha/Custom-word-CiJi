import { defineStore } from 'pinia'
import apiClient from '../api/axios'

export const useUserStore = defineStore('user', {
  state: () => ({
    token: localStorage.getItem('token') || '',
    user: JSON.parse(localStorage.getItem('user')) || null,
    currentWord: null,
    showMeaning: false,
    progress: 0,
    completed: 0,
    total: 0
  }),
  getters: {
    isAuthenticated: (state) => !!state.token,
    isAdmin: (state) => state.user && state.user.role === 'admin'
  },
  actions: {
    setToken(token) {
      this.token = token
      localStorage.setItem('token', token)
    },
    setUser(user) {
      this.user = user
      localStorage.setItem('user', JSON.stringify(user))
    },
    logout() {
      this.token = ''
      this.user = null
      localStorage.removeItem('token')
      localStorage.removeItem('user')
    },
    async login(username, password) {
      const response = await apiClient.post('/login', {
        username,
        password
      })
      this.setToken(response.data.token)
      this.setUser(response.data.user)
      return response.data
    },
    async register(userData) {
      const response = await apiClient.post('/register', userData)
      return response.data
    },
    async getNextWord(level = 'all') {
      try {
        if (level === 'custom') {
          const response = await apiClient.get('/learn/custom/next')
          this.currentWord = response.data
          this.showMeaning = false
          return response.data
        }
        const response = await apiClient.get(`/learn/next?level=${level}`)
        this.currentWord = response.data
        this.showMeaning = false
        return response.data
      } catch (error) {
        if (error.response && error.response.status === 404) {
          // 本轮学习完成
          this.currentWord = null
          this.showMeaning = false
          return null
        }
        throw error
      }
    },
    async submitAnswer(wordID, known) {
      const response = await apiClient.post('/learn/answer', {
        word_id: wordID,
        known
      })
      return response.data
    },
    async submitCustomAnswer(wordID, known) {
      const response = await apiClient.post('/learn/custom/answer', {
        word_id: wordID,
        known
      })
      return response.data
    },
    async getProgress(level = 'all') {
      if (level === 'custom') {
        const response = await apiClient.get('/learn/custom/progress')
        this.progress = response.data.progress
        this.completed = response.data.completed
        this.total = response.data.total
        return response.data
      }
      const response = await apiClient.get(`/learn/progress?level=${level}`)
      this.progress = response.data.progress
      this.completed = response.data.completed
      this.total = response.data.total
      return response.data
    }
  }
})
