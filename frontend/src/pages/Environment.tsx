import { useState, useEffect } from 'react'
import { Card, Table, Button, Modal, Form, Input, Space, Popconfirm, message, Empty, Select } from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined, ApiOutlined } from '@ant-design/icons'
import { useProjectStore } from '../stores/projectStore'
import { useEnvStore, Environment, EnvVariable } from '../stores/envStore'

export default function EnvironmentPage() {
  const { projects, fetch: fetchProjects } = useProjectStore()
  const { environments, loading, fetchByProject, create, update, remove, current: currentEnv, setCurrent: setCurrentEnv, fetchVariables, saveVariables } = useEnvStore()
  const [selectedProjectId, setSelectedProjectId] = useState<number | null>(null)
  const [modalOpen, setModalOpen] = useState(false)
  const [editing, setEditing] = useState<Environment | null>(null)
  const [saving, setSaving] = useState(false)
  const [form] = Form.useForm()
  const [vars, setVars] = useState<EnvVariable[]>([])
  const [varsLoading, setVarsLoading] = useState(false)
  const [varsChanged, setVarsChanged] = useState(false)

  useEffect(() => { fetchProjects() }, []) // eslint-disable-line react-hooks/exhaustive-deps

  useEffect(() => {
    if (selectedProjectId) {
      fetchByProject(selectedProjectId)
      setCurrentEnv(null)
      setVars([])
    }
  }, [selectedProjectId, fetchByProject, setCurrentEnv])

  useEffect(() => {
    if (currentEnv) {
      loadVariables(currentEnv.id)
    }
  }, [currentEnv])

  const loadVariables = async (envId: number) => {
    setVarsLoading(true)
    const data = await fetchVariables(envId)
    setVars(data)
    setVarsChanged(false)
    setVarsLoading(false)
  }

  const handleSaveEnv = async () => {
    try {
      const values = await form.validateFields()
      setSaving(true)
      if (editing) {
        await update(editing.id, values)
        message.success('更新成功')
      } else {
        await create(selectedProjectId!, values)
        message.success('创建成功')
      }
      setModalOpen(false)
      setEditing(null)
      form.resetFields()
    } catch (err: any) {
      if (err.errorFields) return
      message.error(err.message || '操作失败')
    } finally {
      setSaving(false)
    }
  }

  const handleAddVar = () => {
    setVars([...vars, { id: 0, environment_id: currentEnv!.id, key: '', value: '' }])
    setVarsChanged(true)
  }

  const handleRemoveVar = (index: number) => {
    setVars(vars.filter((_, i) => i !== index))
    setVarsChanged(true)
  }

  const handleVarChange = (index: number, field: 'key' | 'value', val: string) => {
    const newVars = [...vars]
    newVars[index] = { ...newVars[index], [field]: val }
    setVars(newVars)
    setVarsChanged(true)
  }

  const handleSaveVars = async () => {
    if (!currentEnv) return
    try {
      await saveVariables(currentEnv.id, vars)
      message.success('变量保存成功')
      setVarsChanged(false)
    } catch (err: any) {
      message.error(err.message || '保存失败')
    }
  }

  const envColumns = [
    { title: '名称', dataIndex: 'name', key: 'name', render: (text: string) => <span style={{ fontWeight: 500 }}>{text}</span> },
    { title: '描述', dataIndex: 'description', key: 'description', ellipsis: true },
    {
      title: '操作', key: 'action', width: 150,
      render: (_: any, record: Environment) => (
        <Space>
          <Button size="small" type="text" style={{ color: '#0071e3' }} onClick={() => { setEditing(record); form.setFieldsValue(record); setModalOpen(true) }}>
            <EditOutlined /> 编辑
          </Button>
          <Popconfirm title="确认删除？" onConfirm={async () => { await remove(record.id); message.success('已删除'); if (currentEnv?.id === record.id) { setCurrentEnv(null); setVars([]) } }}>
            <Button size="small" type="text" danger><DeleteOutlined /> 删除</Button>
          </Popconfirm>
        </Space>
      )
    }
  ]

  return (
    <div>
      <div style={{ marginBottom: 24 }}>
        <Card style={{ borderRadius: 16, border: 'none', boxShadow: '0 2px 12px rgba(0,0,0,0.08)' }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: 16 }}>
            <span style={{ fontSize: 14, color: '#86868b', whiteSpace: 'nowrap' }}>选择项目：</span>
            <Select
              style={{ width: 300 }}
              placeholder="请选择项目"
              value={selectedProjectId}
              onChange={setSelectedProjectId}
              options={projects.map(p => ({ value: p.id, label: p.name }))}
            />
          </div>
        </Card>
      </div>

      {!selectedProjectId ? (
        <Card style={{ borderRadius: 16, border: 'none', boxShadow: '0 2px 12px rgba(0,0,0,0.08)' }}>
          <div style={{ textAlign: 'center', padding: 80 }}>
            <div style={{ width: 64, height: 64, borderRadius: 16, background: '#f5f5f7', display: 'flex', alignItems: 'center', justifyContent: 'center', margin: '0 auto 24px' }}>
              <ApiOutlined style={{ fontSize: 28, color: '#86868b' }} />
            </div>
            <div style={{ fontSize: 20, fontWeight: 600, color: '#1d1d1f', marginBottom: 8 }}>请先选择项目</div>
            <div style={{ fontSize: 14, color: '#86868b' }}>选择一个项目来管理其环境配置</div>
          </div>
        </Card>
      ) : (
        <div style={{ display: 'flex', gap: 24 }}>
          <div style={{ width: 360, flexShrink: 0 }}>
            <Card style={{ borderRadius: 16, border: 'none', boxShadow: '0 2px 12px rgba(0,0,0,0.08)' }} styles={{ body: { padding: 0 } }}>
              <div style={{ padding: '24px 24px 0', borderBottom: '1px solid #f0f0f0', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <h3 style={{ margin: 0, fontSize: 16, fontWeight: 600, color: '#1d1d1f' }}>环境列表</h3>
                <Button type="primary" size="small" icon={<PlusOutlined />} onClick={() => { setEditing(null); form.resetFields(); setModalOpen(true) }} style={{ borderRadius: 8 }}>
                  新建
                </Button>
              </div>
              <Table
                dataSource={environments}
                columns={envColumns}
                rowKey="id"
                loading={loading}
                pagination={false}
                size="small"
                onRow={(record) => ({
                  onClick: () => setCurrentEnv(record),
                  style: { cursor: 'pointer', background: currentEnv?.id === record.id ? '#f0f5ff' : undefined }
                })}
                style={{ margin: '0 16px 16px' }}
              />
            </Card>
          </div>

          <div style={{ flex: 1 }}>
            {currentEnv ? (
              <Card style={{ borderRadius: 16, border: 'none', boxShadow: '0 2px 12px rgba(0,0,0,0.08)' }} styles={{ body: { padding: 0 } }}>
                <div style={{ padding: '24px 24px 0', borderBottom: '1px solid #f0f0f0', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                  <div>
                    <h3 style={{ margin: 0, fontSize: 16, fontWeight: 600, color: '#1d1d1f' }}>{currentEnv.name} - 环境变量</h3>
                    <p style={{ margin: '4px 0 16px', fontSize: 12, color: '#86868b' }}>定义键值对，请求中使用 {'{{key}}'} 引用</p>
                  </div>
                  <Space>
                    <Button icon={<PlusOutlined />} onClick={handleAddVar} style={{ borderRadius: 8 }}>添加变量</Button>
                    {varsChanged && (
                      <Button type="primary" onClick={handleSaveVars} style={{ borderRadius: 8 }}>保存</Button>
                    )}
                  </Space>
                </div>
                <div style={{ padding: 24 }}>
                  {vars.length === 0 ? (
                    <Empty description="暂无变量，点击上方按钮添加" image={Empty.PRESENTED_IMAGE_SIMPLE} />
                  ) : (
                    <div style={{ display: 'flex', flexDirection: 'column', gap: 12 }}>
                      {vars.map((v, i) => (
                        <div key={i} style={{ display: 'flex', gap: 12, alignItems: 'center' }}>
                          <Input placeholder="变量名" value={v.key} onChange={(e) => handleVarChange(i, 'key', e.target.value)} style={{ width: 200, borderRadius: 8 }} />
                          <Input placeholder="值" value={v.value} onChange={(e) => handleVarChange(i, 'value', e.target.value)} style={{ flex: 1, borderRadius: 8 }} />
                          <Button type="text" danger icon={<DeleteOutlined />} onClick={() => handleRemoveVar(i)} />
                        </div>
                      ))}
                    </div>
                  )}
                </div>
              </Card>
            ) : (
              <Card style={{ borderRadius: 16, border: 'none', boxShadow: '0 2px 12px rgba(0,0,0,0.08)' }}>
                <div style={{ textAlign: 'center', padding: 60 }}>
                  <div style={{ fontSize: 16, color: '#86868b' }}>选择一个环境查看变量</div>
                </div>
              </Card>
            )}
          </div>
        </div>
      )}

      <Modal title={editing ? '编辑环境' : '新建环境'} open={modalOpen} onOk={handleSaveEnv} onCancel={() => setModalOpen(false)} confirmLoading={saving} okButtonProps={{ style: { borderRadius: 8 } }} cancelButtonProps={{ style: { borderRadius: 8 } }}>
        <Form form={form} layout="vertical" style={{ marginTop: 16 }}>
          <Form.Item name="name" label="环境名称" rules={[{ required: true }]}><Input placeholder="如：开发环境、测试环境" /></Form.Item>
          <Form.Item name="description" label="描述"><Input.TextArea rows={3} placeholder="环境描述（可选）" /></Form.Item>
        </Form>
      </Modal>
    </div>
  )
}
