import { useState, useEffect } from 'react'
import { Card, Select, Button, Input, Tabs, Space, Tag, Typography, Empty, Spin, message } from 'antd'
import { SendOutlined, ApiOutlined, EnvironmentOutlined } from '@ant-design/icons'
import { useProjectStore } from '../stores/projectStore'
import { useEnvStore } from '../stores/envStore'
import api from '../services/api'

const { Text, Title } = Typography
const { TextArea } = Input

interface DebugResult {
  status_code: number
  headers: Record<string, string>
  body: string
  duration_ms: number
  content_length: number
  error?: string
}

interface APIItem {
  id: number
  name: string
  method: string
  url: string
  protocol: string
}

interface CollectionItem {
  id: number
  name: string
  apis?: APIItem[]
}

export default function DebugPage() {
  const { projects, fetch: fetchProjects } = useProjectStore()
  const { environments, fetchByProject } = useEnvStore()
  const [selectedProjectId, setSelectedProjectId] = useState<number | null>(null)
  const [selectedEnvId, setSelectedEnvId] = useState<number | null>(null)
  const [collections, setCollections] = useState<CollectionItem[]>([])
  const [selectedApi, setSelectedApi] = useState<APIItem | null>(null)

  const [method, setMethod] = useState('GET')
  const [url, setUrl] = useState('')
  const [headers, setHeaders] = useState('')
  const [queryParams, setQueryParams] = useState('')
  const [bodyType, setBodyType] = useState('json')
  const [body, setBody] = useState('')
  const [loading, setLoading] = useState(false)
  const [result, setResult] = useState<DebugResult | null>(null)

  useEffect(() => { fetchProjects() }, [])

  useEffect(() => {
    if (selectedProjectId) {
      fetchByProject(selectedProjectId)
      loadCollections(selectedProjectId)
    }
  }, [selectedProjectId])

  const loadCollections = async (projectId: number) => {
    try {
      const res: any = await api.get(`/projects/${projectId}/collections`)
      const cols = res.data || []
      for (const col of cols) {
        const apiRes: any = await api.get(`/collections/${col.id}/apis`)
        col.apis = apiRes.data || []
      }
      setCollections(cols)
    } catch (err) {
      console.error(err)
    }
  }

  const handleSelectApi = async (apiId: number) => {
    try {
      const res: any = await api.get(`/apis/${apiId}`)
      const apiDef = res.data
      setSelectedApi(apiDef)
      setMethod(apiDef.method || 'GET')
      setUrl(apiDef.url || '')
      setHeaders(apiDef.headers || '')
      setQueryParams(apiDef.query_params || '')
      setBodyType(apiDef.body_type || 'json')
      setBody(apiDef.body || '')
    } catch (err: any) {
      message.error(err.message || '加载接口失败')
    }
  }

  const handleSend = async () => {
    if (!url) {
      message.warning('请输入请求 URL')
      return
    }
    if (!selectedApi?.id) {
      message.warning('请先选择一个接口')
      return
    }
    setLoading(true)
    setResult(null)
    try {
      let parsedHeaders: Record<string, string> = {}
      if (headers) {
        try { parsedHeaders = JSON.parse(headers) } catch { parsedHeaders = {} }
      }
      let parsedQuery: Record<string, string> = {}
      if (queryParams) {
        try { parsedQuery = JSON.parse(queryParams) } catch { parsedQuery = {} }
      }
      const res: any = await api.post(`/apis/${selectedApi?.id || 0}/debug`, {
        method,
        url,
        headers: parsedHeaders,
        query_params: parsedQuery,
        body_type: bodyType,
        body,
        env_id: selectedEnvId,
        timeout_ms: 30000,
      })
      setResult(res.data)
    } catch (err: any) {
      setResult({
        status_code: 0,
        headers: {},
        body: '',
        duration_ms: 0,
        content_length: 0,
        error: err.message || '请求失败',
      })
    } finally {
      setLoading(false)
    }
  }

  const allApis = collections.flatMap(c => (c.apis || []).map(a => ({ ...a, collectionName: c.name })))

  return (
    <div style={{ display: 'flex', gap: 24, height: 'calc(100vh - 128px)' }}>
      {/* Left: API List */}
      <div style={{ width: 280, flexShrink: 0 }}>
        <Card style={{ borderRadius: 16, border: 'none', boxShadow: '0 2px 12px rgba(0,0,0,0.08)', height: '100%' }} styles={{ body: { padding: 0 } }}>
          <div style={{ padding: '20px 20px 12px', borderBottom: '1px solid #f0f0f0' }}>
            <h3 style={{ margin: '0 0 12px', fontSize: 15, fontWeight: 600, color: '#1d1d1f' }}>接口列表</h3>
            <Select style={{ width: '100%' }} placeholder="选择项目" value={selectedProjectId} onChange={setSelectedProjectId} options={projects.map(p => ({ value: p.id, label: p.name }))} />
          </div>
          <div style={{ padding: '8px 0', overflowY: 'auto', maxHeight: 'calc(100vh - 280px)' }}>
            {collections.map(col => (
              <div key={col.id}>
                <div style={{ padding: '8px 20px', fontSize: 12, color: '#86868b', fontWeight: 500 }}>{col.name}</div>
                {(col.apis || []).map(apiItem => (
                  <div
                    key={apiItem.id}
                    onClick={() => handleSelectApi(apiItem.id)}
                    style={{
                      padding: '8px 20px 8px 32px',
                      cursor: 'pointer',
                      background: selectedApi?.id === apiItem.id ? '#f0f5ff' : undefined,
                      fontSize: 13,
                      color: '#1d1d1f',
                    }}
                  >
                    <Tag color={apiItem.method === 'GET' ? 'green' : apiItem.method === 'POST' ? 'blue' : apiItem.method === 'PUT' ? 'orange' : 'red'} style={{ fontSize: 11, marginRight: 8 }}>
                      {apiItem.method}
                    </Tag>
                    {apiItem.name}
                  </div>
                ))}
              </div>
            ))}
            {collections.length === 0 && (
              <div style={{ padding: 40, textAlign: 'center', color: '#86868b', fontSize: 13 }}>暂无接口</div>
            )}
          </div>
        </Card>
      </div>

      {/* Center: Request Editor */}
      <div style={{ flex: 1, display: 'flex', flexDirection: 'column', gap: 16 }}>
        <Card style={{ borderRadius: 16, border: 'none', boxShadow: '0 2px 12px rgba(0,0,0,0.08)' }} styles={{ body: { padding: 16 } }}>
          <div style={{ display: 'flex', gap: 8, marginBottom: 12 }}>
            <Select style={{ width: 120 }} value={method} onChange={setMethod} options={['GET', 'POST', 'PUT', 'DELETE', 'PATCH', 'HEAD', 'OPTIONS'].map(m => ({ value: m, label: m }))} />
            <Input value={url} onChange={e => setUrl(e.target.value)} placeholder="请输入请求 URL" style={{ flex: 1, borderRadius: 8 }} />
            <Select style={{ width: 160 }} placeholder="选择环境" value={selectedEnvId} onChange={setSelectedEnvId} allowClear options={environments.map(e => ({ value: e.id, label: e.name }))} />
            <Button type="primary" icon={<SendOutlined />} loading={loading} onClick={handleSend} style={{ borderRadius: 8 }}>
              发送
            </Button>
          </div>

          <Tabs
            size="small"
            items={[
              {
                key: 'headers',
                label: '请求头',
                children: <TextArea value={headers} onChange={e => setHeaders(e.target.value)} placeholder='{"Content-Type": "application/json"}' rows={4} style={{ borderRadius: 8, fontFamily: 'monospace' }} />
              },
              {
                key: 'params',
                label: 'Query 参数',
                children: <TextArea value={queryParams} onChange={e => setQueryParams(e.target.value)} placeholder='{"key": "value"}' rows={4} style={{ borderRadius: 8, fontFamily: 'monospace' }} />
              },
              {
                key: 'body',
                label: '请求体',
                children: (
                  <div>
                    <Select size="small" value={bodyType} onChange={setBodyType} style={{ width: 120, marginBottom: 8 }} options={[{ value: 'json', label: 'JSON' }, { value: 'text', label: 'Text' }, { value: 'form', label: 'Form' }]} />
                    <TextArea value={body} onChange={e => setBody(e.target.value)} placeholder='{"key": "value"}' rows={6} style={{ borderRadius: 8, fontFamily: 'monospace' }} />
                  </div>
                )
              }
            ]}
          />
        </Card>

        {/* Response */}
        <Card style={{ borderRadius: 16, border: 'none', boxShadow: '0 2px 12px rgba(0,0,0,0.08)', flex: 1 }} styles={{ body: { padding: 0 } }}>
          <div style={{ padding: '12px 20px', borderBottom: '1px solid #f0f0f0', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <h3 style={{ margin: 0, fontSize: 14, fontWeight: 600, color: '#1d1d1f' }}>响应</h3>
            {result && (
              <Space>
                <Tag color={result.error ? 'red' : result.status_code < 400 ? 'green' : 'orange'}>
                  {result.error ? '错误' : result.status_code}
                </Tag>
                {!result.error && <Text type="secondary" style={{ fontSize: 12 }}>{result.duration_ms}ms</Text>}
              </Space>
            )}
          </div>
          <div style={{ padding: 16, overflow: 'auto', maxHeight: 400 }}>
            {loading ? (
              <div style={{ textAlign: 'center', padding: 40 }}><Spin /></div>
            ) : result ? (
              result.error ? (
                <div style={{ color: '#ff4d4f', padding: 20 }}>{result.error}</div>
              ) : (
                <div>
                  <div style={{ marginBottom: 12 }}>
                    <Text type="secondary" style={{ fontSize: 12 }}>响应头</Text>
                    <div style={{ marginTop: 4, fontSize: 12, fontFamily: 'monospace', background: '#f5f5f7', padding: 12, borderRadius: 8 }}>
                      {Object.entries(result.headers).map(([k, v]) => (
                        <div key={k}><span style={{ color: '#0071e3' }}>{k}</span>: {v}</div>
                      ))}
                    </div>
                  </div>
                  <div>
                    <Text type="secondary" style={{ fontSize: 12 }}>响应体</Text>
                    <pre style={{ marginTop: 4, fontSize: 12, fontFamily: 'monospace', background: '#f5f5f7', padding: 12, borderRadius: 8, overflow: 'auto', maxHeight: 300, whiteSpace: 'pre-wrap', wordBreak: 'break-all' }}>
                      {result.body}
                    </pre>
                  </div>
                </div>
              )
            ) : (
              <Empty description="点击发送按钮查看响应" image={Empty.PRESENTED_IMAGE_SIMPLE} />
            )}
          </div>
        </Card>
      </div>
    </div>
  )
}
