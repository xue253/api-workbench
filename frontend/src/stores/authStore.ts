import { create } from 'zustand'
import api from '../services/api'

export interface User {
  id: number
  username: string
  email: string
  avatar: string
  created_at: string
}

interface AuthState {
  user: User | null
  token: string | null
  loading: boolean
  login: (username: string, password: string) => Promise<void>
  register: (username: string, password: string, email?: string) => Promise<void>
  logout: () => void
  fetchProfile: () => Promise<void>
  isLoggedIn: () => boolean
}

export const useAuthStore = create<AuthState>((set, get) => ({
  user: null,
  token: localStorage.getItem('token'),
  loading: false,

  login: async (username, password) => {
    const res: any = await api.post('/auth/login', { username, password })
    const { token, user } = res.data
    localStorage.setItem('token', token)
    api.defaults.headers.common['Authorization'] = `Bearer ${token}`
    set({ user, token })
  },

  register: async (username, password, email) => {
    const res: any = await api.post('/auth/register', { username, password, email })
    const { token, user } = res.data
    localStorage.setItem('token', token)
    api.defaults.headers.common['Authorization'] = `Bearer ${token}`
    set({ user, token })
  },

  logout: () => {
    localStorage.removeItem('token')
    delete api.defaults.headers.common['Authorization']
    set({ user: null, token: null })
  },

  fetchProfile: async () => {
    const token = get().token
    if (!token) return
    api.defaults.headers.common['Authorization'] = `Bearer ${token}`
    try {
      const res: any = await api.get('/user/profile')
      set({ user: res.data })
    } catch {
      get().logout()
    }
  },

  isLoggedIn: () => !!get().token,
}))
