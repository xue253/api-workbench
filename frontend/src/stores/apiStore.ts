import { create } from 'zustand'
import api from '../services/api'

export interface Collection {
  id: number
  project_id: number
  parent_id: number | null
  name: string
  description: string
  sort_order: number
  children?: Collection[]
  apis?: APIItem[]
}

export interface APIItem {
  id: number
  collection_id: number
  name: string
  description: string
  protocol: string
  method: string
  url: string
  headers: string
  path_params: string
  query_params: string
  body_type: string
  body: string
  expected_status: number
  timeout_ms: number
}

interface APIState {
  collections: Collection[]
  currentAPI: APIItem | null
  loading: boolean
  fetchCollections: (projectId: number) => Promise<void>
  createCollection: (projectId: number, data: Partial<Collection>) => Promise<Collection>
  updateCollection: (id: number, data: Partial<Collection>) => Promise<void>
  deleteCollection: (id: number) => Promise<void>
  moveCollection: (id: number, parentId: number | null) => Promise<void>
  createAPI: (collectionId: number, data: Partial<APIItem>) => Promise<APIItem>
  updateAPI: (id: number, data: Partial<APIItem>) => Promise<void>
  deleteAPI: (id: number) => Promise<void>
  getAPI: (id: number) => Promise<APIItem>
  setCurrentAPI: (api: APIItem | null) => void
}

export const useAPIStore = create<APIState>((set, get) => ({
  collections: [],
  currentAPI: null,
  loading: false,

  fetchCollections: async (projectId) => {
    set({ loading: true })
    const res: any = await api.get(`/projects/${projectId}/collections`)
    const collections = res.data || []
    
    for (const col of collections) {
      const apiRes: any = await api.get(`/collections/${col.id}/apis`)
      col.apis = apiRes.data || []
    }
    
    set({ collections, loading: false })
  },

  createCollection: async (projectId, data) => {
    const res: any = await api.post(`/projects/${projectId}/collections`, data)
    set({ collections: [...get().collections, { ...res.data, apis: [] }] })
    return res.data
  },

  updateCollection: async (id, data) => {
    const res: any = await api.put(`/collections/${id}`, data)
    set({ collections: get().collections.map(c => c.id === id ? { ...c, ...res.data } : c) })
  },

  deleteCollection: async (id) => {
    await api.delete(`/collections/${id}`)
    set({ collections: get().collections.filter(c => c.id !== id) })
  },

  moveCollection: async (id, parentId) => {
    await api.post(`/collections/${id}/move`, { parent_id: parentId })
  },

  createAPI: async (collectionId, data) => {
    const res: any = await api.post(`/collections/${collectionId}/apis`, data)
    const collections = get().collections.map(c => {
      if (c.id === collectionId) {
        return { ...c, apis: [...(c.apis || []), res.data] }
      }
      return c
    })
    set({ collections })
    return res.data
  },

  updateAPI: async (id, data) => {
    const res: any = await api.put(`/apis/${id}`, data)
    const collections = get().collections.map(c => ({
      ...c,
      apis: (c.apis || []).map(a => a.id === id ? res.data : a)
    }))
    set({ collections })
  },

  deleteAPI: async (id) => {
    await api.delete(`/apis/${id}`)
    const collections = get().collections.map(c => ({
      ...c,
      apis: (c.apis || []).filter(a => a.id !== id)
    }))
    set({ collections })
  },

  getAPI: async (id) => {
    const res: any = await api.get(`/apis/${id}`)
    return res.data
  },

  setCurrentAPI: (api) => set({ currentAPI: api }),
}))
