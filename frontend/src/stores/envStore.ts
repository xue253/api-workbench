import { create } from 'zustand'
import api from '../services/api'

export interface Environment {
  id: number
  project_id: number
  name: string
  description: string
  sort_order: number
  variables?: EnvVariable[]
}

export interface EnvVariable {
  id: number
  environment_id: number
  key: string
  value: string
}

interface EnvState {
  environments: Environment[]
  current: Environment | null
  loading: boolean
  fetchByProject: (projectId: number) => Promise<void>
  create: (projectId: number, data: Partial<Environment>) => Promise<Environment>
  update: (id: number, data: Partial<Environment>) => Promise<void>
  remove: (id: number) => Promise<void>
  setCurrent: (env: Environment | null) => void
  fetchVariables: (envId: number) => Promise<EnvVariable[]>
  saveVariables: (envId: number, vars: Partial<EnvVariable>[]) => Promise<void>
}

export const useEnvStore = create<EnvState>((set, get) => ({
  environments: [],
  current: null,
  loading: false,

  fetchByProject: async (projectId) => {
    set({ loading: true })
    const res: any = await api.get(`/projects/${projectId}/environments`)
    set({ environments: res.data || [], loading: false })
  },

  create: async (projectId, data) => {
    const res: any = await api.post(`/projects/${projectId}/environments`, data)
    set({ environments: [...get().environments, res.data] })
    return res.data
  },

  update: async (id, data) => {
    const res: any = await api.put(`/environments/${id}`, data)
    set({ environments: get().environments.map(e => e.id === id ? res.data : e) })
  },

  remove: async (id) => {
    await api.delete(`/environments/${id}`)
    set({ environments: get().environments.filter(e => e.id !== id) })
  },

  setCurrent: (env) => set({ current: env }),

  fetchVariables: async (envId) => {
    const res: any = await api.get(`/environments/${envId}/variables`)
    return res.data || []
  },

  saveVariables: async (envId, vars) => {
    await api.put(`/environments/${envId}/variables`, vars)
  },
}))
