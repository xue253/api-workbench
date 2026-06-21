import { useState, useEffect } from 'react'
import { Card, Tree, Button, Modal, Form, Input, Select, Space, Popconfirm, message, Tabs, Tag, Empty } from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined, FolderOutlined, FolderOpenOutlined, ApiOutlined, SendOutlined } from '@ant-design/icons'
import { useProjectStore } from '../stores/projectStore'
import { useAPIStore, Collection, APIItem } from '../stores/apiStore'
import api from '../services/api'

export default function ApiExplorerPage() {
  const { projects, fetch: fetchProjects } = useProjectStore()
  const { collections, loading, fetchCollections, createCollection, updateCollection, deleteCollection, createAPI, updateAPI, deleteAPI, currentAPI, setCurrentAPI } = useAPIStore()
  const [selectedProjectId, setSelectedProjectId] = useState<number | null>(null)
  const [colModalOpen, setColModalOpen] = useState(false)
  const [apiModalOpen, setApiModalOpen] = useState(false)
  const [editingCol, setEditingCol] = useState<Collection | null>(null)
  const [editingAPI, setEditingAPI] = useState<APIItem | null>(null)
  const [saving, setSaving] = useState(false)
  const [colForm] = Form.useForm()
  const [apiForm] = Form.useForm()

  useEffect(() => { fetchProjects() }, [])

  useEffect(() => {
    if (selectedProjectId) {
      fetchCollections(selectedProjectId)
      setCurrentAPI(null)
    }
  }, [selectedProjectId])

  const handleSaveCol = async () => {
    try {
      const values = await colForm.validateFields()
      setSaving(true)
      if (editingCol) {
        await updateCollection(editingCol.id, values)
      } else {
        await createCollection(selectedProjectId!, values)
      }
      message.success('保存成功')
      setColModalOpen(false)
      setEditingCol(null)
      colForm.resetFields()
    } catch (err: any) {
      if (err.errorFields) return
      message.error(err.message || '操作失败')
    } finally {
      setSaving(false)
    }
  }

  const handleSaveAPI = async () => {
    try {
      const values = await apiForm.validateFields()
      setSaving(true)
      if (editingAPI) {
        await updateAPI(editingAPI.id, values)
      } else {
        await createAPI(values.collection_id, values)
      }
      message.success('保存成功')
      setApiModalOpen(false)
      setEditingAPI(null)
      apiForm.resetFields()
    } catch (err: any) {
      if (err.errorFields) return
      message.error(err.message || '操作失败')
    } finally {
      setSaving(false)
    }
  }

  const handleEditAPI = async (apiId: number) => {
    try {
      const res: any = await api.get(`/apis/${apiId}`)
      setEditingAPI(res.data)
      apiForm.setFieldsValue(res.data)
      setApiModalOpen(true)
    } catch (err: any) {
      message.error(err.message || '加载失败')
    }
  }

  const buildTreeData = (cols: Collection[]): any[] => {
    return cols.map(col => ({
      key: `col-${col.id}`,
      title: (
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <span><FolderOutlined style={{ marginRight: 8, color: '#faad14' }} />{col.name}</span>
          <Space size={4}>
            <Button size="small" type="text" onClick={(e) => { e.stopPropagation(); setEditingCol(col); colForm.setFieldsValue(col); setColModalOpen(true) }}>
              <EditOutlined />
            </Button>
            <Popconfirm title="确认删除？" onConfirm={() => deleteCollection(col.id)}>
              <Button size="small" type="text" danger onClick={(e) => e.stopPropagation()}>
                <DeleteOutlined />
              </Button>
            </Popconfirm>
          </Space>
        </div>
      ),
      children: [
        ...(col.apis || []).map(apiItem => ({
          key: `api-${apiItem.id}`,
          title: (
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
              <span>
                <Tag color={apiItem.method === 'GET' ? 'green' : apiItem.method === 'POST' ? 'blue' : apiItem.method === 'PUT' ? 'orange' : 'red'} style={{ fontSize: 11, marginRight: 8 }}>
                  {apiItem.method}
                </Tag>
                {apiItem.name}
              </span>
              <Space size={4}>
                <Button size="small" type="text" onClick={(e) => { e.stopPropagation(); handleEditAPI(apiItem.id) }}>
                  <EditOutlined />
                </Button>
                <Popconfirm title="确认删除？" onConfirm={() => deleteAPI(apiItem.id)}>
                  <Button size="small" type="text" danger onClick={(e) => e.stopPropagation()}>
                    <DeleteOutlined />
                  </Button>
                </Popconfirm>
              </Space>
            </div>
          ),
          isLeaf: true,
        })),
      ],
    }))
  }

  return (
    <div style={{ display: 'flex', gap: 24, height: 'calc(100vh - 128px)' }}>
      <div style={{ width: 320, flexShrink: 0 }}>
        <Card style={{ borderRadius: 16, border: 'none', boxShadow: '0 2px 12px rgba(0,0,0,0.08)', height: '100%' }} styles={{ body: { padding: 0 } }}>
          <div style={{ padding: '20px 20px 12px', borderBottom: '1px solid #f0f0f0' }}>
            <h3 style={{ margin: '0 0 12px', fontSize: 15, fontWeight: 600, color: '#1d1d1f' }}>接口库</h3>
            <Select style={{ width: '100%' }} placeholder="选择项目" value={selectedProjectId} onChange={setSelectedProjectId} options={projects.map(p => ({ value: p.id, label: p.name }))} />
            {selectedProjectId && (
              <Space style={{ marginTop: 12 }}>
                <Button size="small" icon={<PlusOutlined />} onClick={() => { setEditingCol(null); colForm.resetFields(); setColModalOpen(true) }} style={{ borderRadius: 8 }}>新建集合</Button>
                <Button size="small" icon={<PlusOutlined />} onClick={() => { setEditingAPI(null); apiForm.resetFields(); setApiModalOpen(true) }} style={{ borderRadius: 8 }}>新建接口</Button>
              </Space>
            )}
          </div>
          <div style={{ padding: '8px 0', overflowY: 'auto', maxHeight: 'calc(100vh - 280px)' }}>
            {collections.length > 0 ? (
              <Tree
                showIcon
                defaultExpandAll
                treeData={buildTreeData(collections)}
                onSelect={(keys) => {
                  const key = keys[0] as string
                  if (key?.startsWith('api-')) {
                    const apiId = parseInt(key.replace('api-', ''))
                    const col = collections.find(c => c.apis?.some(a => a.id === apiId))
                    const api = col?.apis?.find(a => a.id === apiId)
                    if (api) setCurrentAPI(api)
                  }
                }}
                style={{ padding: '0 8px' }}
              />
            ) : (
              <div style={{ padding: 40, textAlign: 'center', color: '#86868b', fontSize: 13 }}>
                {loading ? '加载中...' : '暂无数据'}
              </div>
            )}
          </div>
        </Card>
      </div>

      <div style={{ flex: 1 }}>
        {currentAPI ? (
          <Card style={{ borderRadius: 16, border: 'none', boxShadow: '0 2px 12px rgba(0,0,0,0.08)' }} styles={{ body: { padding: 0 } }}>
            <div style={{ padding: '20px 24px 0', borderBottom: '1px solid #f0f0f0', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
              <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
                <Tag color={currentAPI.method === 'GET' ? 'green' : currentAPI.method === 'POST' ? 'blue' : 'orange'}>{currentAPI.method}</Tag>
                <span style={{ fontSize: 16, fontWeight: 600, color: '#1d1d1f' }}>{currentAPI.name}</span>
              </div>
              <Button type="primary" icon={<SendOutlined />} onClick={() => setCurrentAPI(null)} style={{ borderRadius: 8 }}>调试</Button>
            </div>
            <div style={{ padding: 24 }}>
              <Tabs items={[
                {
                  key: 'info',
                  label: '基本信息',
                  children: (
                    <div style={{ display: 'grid', gridTemplateColumns: '120px 1fr', gap: '12px 16px', fontSize: 14 }}>
                      <div style={{ color: '#86868b' }}>URL</div>
                      <div style={{ fontFamily: 'monospace', wordBreak: 'break-all' }}>{currentAPI.url}</div>
                      <div style={{ color: '#86868b' }}>协议</div>
                      <div>{currentAPI.protocol}</div>
                      <div style={{ color: '#86868b' }}>超时</div>
                      <div>{currentAPI.timeout_ms}ms</div>
                      <div style={{ color: '#86868b' }}>描述</div>
                      <div>{currentAPI.description || '-'}</div>
                    </div>
                  )
                },
                {
                  key: 'headers',
                  label: '请求头',
                  children: <pre style={{ background: '#f5f5f7', padding: 16, borderRadius: 8, fontSize: 12, fontFamily: 'monospace', margin: 0 }}>{currentAPI.headers || '{}'}</pre>
                },
                {
                  key: 'body',
                  label: '请求体',
                  children: (
                    <div>
                      <Tag style={{ marginBottom: 12 }}>{currentAPI.body_type || 'json'}</Tag>
                      <pre style={{ background: '#f5f5f7', padding: 16, borderRadius: 8, fontSize: 12, fontFamily: 'monospace', margin: 0, maxHeight: 300, overflow: 'auto' }}>{currentAPI.body || '空'}</pre>
                    </div>
                  )
                }
              ]} />
            </div>
          </Card>
        ) : (
          <Card style={{ borderRadius: 16, border: 'none', boxShadow: '0 2px 12px rgba(0,0,0,0.08)' }}>
            <div style={{ textAlign: 'center', padding: 80 }}>
              <div style={{ width: 64, height: 64, borderRadius: 16, background: '#f5f5f7', display: 'flex', alignItems: 'center', justifyContent: 'center', margin: '0 auto 24px' }}>
                <ApiOutlined style={{ fontSize: 28, color: '#86868b' }} />
              </div>
              <div style={{ fontSize: 20, fontWeight: 600, color: '#1d1d1f', marginBottom: 8 }}>选择接口查看详情</div>
              <div style={{ fontSize: 14, color: '#86868b' }}>从左侧列表选择一个接口</div>
            </div>
          </Card>
        )}
      </div>

      <Modal title={editingCol ? '编辑集合' : '新建集合'} open={colModalOpen} onOk={handleSaveCol} onCancel={() => setColModalOpen(false)} confirmLoading={saving} okButtonProps={{ style: { borderRadius: 8 } }} cancelButtonProps={{ style: { borderRadius: 8 } }}>
        <Form form={colForm} layout="vertical" style={{ marginTop: 16 }}>
          <Form.Item name="name" label="集合名称" rules={[{ required: true }]}><Input placeholder="输入集合名称" /></Form.Item>
          <Form.Item name="description" label="描述"><Input.TextArea rows={3} placeholder="描述（可选）" /></Form.Item>
        </Form>
      </Modal>

      <Modal title={editingAPI ? '编辑接口' : '新建接口'} open={apiModalOpen} onOk={handleSaveAPI} onCancel={() => setApiModalOpen(false)} confirmLoading={saving} width={640} okButtonProps={{ style: { borderRadius: 8 } }} cancelButtonProps={{ style: { borderRadius: 8 } }}>
        <Form form={apiForm} layout="vertical" style={{ marginTop: 16 }}>
          <Form.Item name="collection_id" label="所属集合" rules={[{ required: true }]}>
            <Select placeholder="选择集合" options={collections.map(c => ({ value: c.id, label: c.name }))} />
          </Form.Item>
          <Form.Item name="name" label="接口名称" rules={[{ required: true }]}><Input placeholder="输入接口名称" /></Form.Item>
          <Form.Item name="description" label="描述"><Input.TextArea rows={2} placeholder="描述（可选）" /></Form.Item>
          <div style={{ display: 'flex', gap: 12 }}>
            <Form.Item name="method" label="请求方法" rules={[{ required: true }]} style={{ width: 120 }}>
              <Select options={['GET', 'POST', 'PUT', 'DELETE', 'PATCH'].map(m => ({ value: m, label: m }))} />
            </Form.Item>
            <Form.Item name="url" label="URL" rules={[{ required: true }]} style={{ flex: 1 }}>
              <Input placeholder="https://api.example.com/path" />
            </Form.Item>
          </div>
          <Form.Item name="body_type" label="请求体类型">
            <Select options={[{ value: 'json', label: 'JSON' }, { value: 'form', label: 'Form' }, { value: 'text', label: 'Text' }, { value: 'xml', label: 'XML' }]} />
          </Form.Item>
          <Form.Item name="body" label="请求体">
            <Input.TextArea rows={4} placeholder='{"key": "value"}' style={{ fontFamily: 'monospace' }} />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}
