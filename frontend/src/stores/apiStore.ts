import { create } from 'zustand'
import api from '../services/api'

export interface Collection {
  id: number
  project_id: number
  parent_id: number | null
  name: string
  description: string
  sort_order: number
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
  proto_service: string
  proto_method: string
  expected_status: number
  timeout_ms: number
}

interface ApiStoreState {
  collections: Collection[]
  apis: APIItem[]
  currentApi: APIItem | null
  loading: boolean
  fetchCollections: (pid: number) => Promise<void>
  createCollection: (pid: number, data: Partial<Collection>) => Promise<Collection>
  updateCollection: (id: number, data: Partial<Collection>) => Promise<void>
  deleteCollection: (id: number) => Promise<void>
  fetchAPIs: (cid: number) => Promise<void>
  getApi: (id: number) => Promise<void>
  createAPI: (cid: number, data: Partial<APIItem>) => Promise<APIItem>
  updateAPI: (id: number, data: Partial<APIItem>) => Promise<void>
  deleteAPI: (id: number) => Promise<void>
  setCurrentApi: (a: APIItem | null) => void
}

export const useApiStore = create<ApiStoreState>((set, get) => ({
  collections: [],
  apis: [],
  currentApi: null,
  loading: false,

  fetchCollections: async (pid) => {
    set({ loading: true })
    const res: any = await api.get(`/projects/${pid}/collections`)
    set({ collections: res.data || [], loading: false })
  },
  createCollection: async (pid, data) => {
    const res: any = await api.post(`/projects/${pid}/collections`, data)
    set({ collections: [...get().collections, res.data] })
    return res.data
  },
  updateCollection: async (id, data) => {
    const res: any = await api.put(`/collections/${id}`, data)
    set({ collections: get().collections.map(c => c.id === id ? res.data : c) })
  },
  deleteCollection: async (id) => {
    await api.delete(`/collections/${id}`)
    set({ collections: get().collections.filter(c => c.id !== id) })
  },

  fetchAPIs: async (cid) => {
    set({ loading: true })
    const res: any = await api.get(`/collections/${cid}/apis`)
    set({ apis: res.data || [], loading: false })
  },
  getApi: async (id) => {
    const res: any = await api.get(`/apis/${id}`)
    set({ currentApi: res.data })
  },
  createAPI: async (cid, data) => {
    const res: any = await api.post(`/collections/${cid}/apis`, data)
    set({ apis: [...get().apis, res.data] })
    return res.data
  },
  updateAPI: async (id, data) => {
    const res: any = await api.put(`/apis/${id}`, data)
    set({ apis: get().apis.map(a => a.id === id ? res.data : a) })
  },
  deleteAPI: async (id) => {
    await api.delete(`/apis/${id}`)
    set({ apis: get().apis.filter(a => a.id !== id) })
  },
  setCurrentApi: (a) => set({ currentApi: a }),
}))
