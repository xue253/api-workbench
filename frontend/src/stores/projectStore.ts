import { create } from 'zustand'
import api from '../services/api'

export interface Project {
  id: number
  name: string
  description: string
  created_at: string
}

interface ProjectState {
  projects: Project[]
  current: Project | null
  loading: boolean
  fetch: () => Promise<void>
  create: (data: Partial<Project>) => Promise<Project>
  update: (id: number, data: Partial<Project>) => Promise<void>
  remove: (id: number) => Promise<void>
  setCurrent: (p: Project | null) => void
}

export const useProjectStore = create<ProjectState>((set, get) => ({
  projects: [],
  current: null,
  loading: false,
  fetch: async () => {
    set({ loading: true })
    const res: any = await api.get('/projects')
    set({ projects: res.data || [], loading: false })
  },
  create: async (data) => {
    const res: any = await api.post('/projects', data)
    set({ projects: [...get().projects, res.data] })
    return res.data
  },
  update: async (id, data) => {
    const res: any = await api.put(`/projects/${id}`, data)
    set({ projects: get().projects.map(p => p.id === id ? res.data : p) })
  },
  remove: async (id) => {
    await api.delete(`/projects/${id}`)
    set({ projects: get().projects.filter(p => p.id !== id) })
  },
  setCurrent: (p) => set({ current: p }),
}))
