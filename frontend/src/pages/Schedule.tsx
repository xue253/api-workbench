import { useState, useEffect } from 'react'
import { Card, Table, Button, Modal, Form, Input, InputNumber, Select, Switch, Space, Popconfirm, message, Tag } from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined, ClockCircleOutlined } from '@ant-design/icons'
import { useProjectStore } from '../stores/projectStore'
import api from '../services/api'

interface ScheduledTask {
  id: number
  project_id: number
  target_type: string
  target_id: number
  cron_expr: string
  enabled: boolean
  environment_id: number
  created_at: string
}

export default function SchedulePage() {
  const { projects, fetch: fetchProjects } = useProjectStore()
  const [selectedProjectId, setSelectedProjectId] = useState<number | null>(null)
  const [tasks, setTasks] = useState<ScheduledTask[]>([])
  const [loading, setLoading] = useState(false)
  const [modalOpen, setModalOpen] = useState(false)
  const [editing, setEditing] = useState<ScheduledTask | null>(null)
  const [saving, setSaving] = useState(false)
  const [form] = Form.useForm()

  useEffect(() => { fetchProjects() }, [])

  useEffect(() => {
    if (selectedProjectId) {
      loadTasks()
    }
  }, [selectedProjectId])

  const loadTasks = async () => {
    setLoading(true)
    try {
      const res: any = await api.get(`/projects/${selectedProjectId}/schedules`)
      setTasks(res.data || [])
    } catch (err) {
      console.error(err)
    }
    setLoading(false)
  }

  const handleSave = async () => {
    try {
      const values = await form.validateFields()
      setSaving(true)
      if (editing) {
        await api.put(`/schedules/${editing.id}`, values)
      } else {
        await api.post(`/projects/${selectedProjectId}/schedules`, values)
      }
      message.success('保存成功')
      setModalOpen(false)
      setEditing(null)
      form.resetFields()
      loadTasks()
    } catch (err: any) {
      if (err.errorFields) return
      message.error(err.message || '操作失败')
    } finally {
      setSaving(false)
    }
  }

  const handleDelete = async (id: number) => {
    try {
      await api.delete(`/schedules/${id}`)
      message.success('删除成功')
      loadTasks()
    } catch (err: any) {
      message.error(err.message || '删除失败')
    }
  }

  const handleToggle = async (record: ScheduledTask) => {
    try {
      await api.put(`/schedules/${record.id}`, { ...record, enabled: !record.enabled })
      message.success(record.enabled ? '已禁用' : '已启用')
      loadTasks()
    } catch (err: any) {
      message.error(err.message || '操作失败')
    }
  }

  const columns = [
    { title: '类型', dataIndex: 'target_type', key: 'target_type', width: 100, render: (v: string) => <Tag color={v === 'test_suite' ? 'blue' : 'green'}>{v === 'test_suite' ? '套件' : '用例'}</Tag> },
    { title: 'Cron 表达式', dataIndex: 'cron_expr', key: 'cron_expr', width: 150, render: (v: string) => <span style={{ fontFamily: 'monospace' }}>{v}</span> },
    { title: '状态', dataIndex: 'enabled', key: 'enabled', width: 80, render: (v: boolean) => <Tag color={v ? 'green' : 'default'}>{v ? '启用' : '禁用'}</Tag> },
    {
      title: '操作', key: 'action', width: 200,
      render: (_: any, record: ScheduledTask) => (
        <Space>
          <Switch size="small" checked={record.enabled} onChange={() => handleToggle(record)} />
          <Button size="small" type="text" style={{ color: '#0071e3' }} onClick={() => { setEditing(record); form.setFieldsValue(record); setModalOpen(true) }}>
            <EditOutlined /> 编辑
          </Button>
          <Popconfirm title="确认删除？" onConfirm={() => handleDelete(record.id)}>
            <Button size="small" type="text" danger><DeleteOutlined /> 删除</Button>
          </Popconfirm>
        </Space>
      )
    }
  ]

  const cronPresets = [
    { label: '每分钟', value: '* * * * *' },
    { label: '每小时', value: '0 * * * *' },
    { label: '每天凌晨', value: '0 0 * * *' },
    { label: '每天上午9点', value: '0 9 * * *' },
    { label: '每周一', value: '0 0 * * 1' },
    { label: '每月1号', value: '0 0 1 * *' },
  ]

  if (!selectedProjectId) {
    return (
      <Card style={{ borderRadius: 16, border: 'none', boxShadow: '0 2px 12px rgba(0,0,0,0.08)' }}>
        <div style={{ textAlign: 'center', padding: 80 }}>
          <div style={{ width: 64, height: 64, borderRadius: 16, background: '#f5f5f7', display: 'flex', alignItems: 'center', justifyContent: 'center', margin: '0 auto 24px' }}>
            <ClockCircleOutlined style={{ fontSize: 28, color: '#86868b' }} />
          </div>
          <div style={{ fontSize: 20, fontWeight: 600, color: '#1d1d1f', marginBottom: 8 }}>请先选择项目</div>
          <div style={{ fontSize: 14, color: '#86868b' }}>选择一个项目来管理定时调度</div>
        </div>
      </Card>
    )
  }

  return (
    <div>
      <div style={{ marginBottom: 24 }}>
        <Card style={{ borderRadius: 16, border: 'none', boxShadow: '0 2px 12px rgba(0,0,0,0.08)' }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: 16 }}>
            <span style={{ fontSize: 14, color: '#86868b', whiteSpace: 'nowrap' }}>选择项目：</span>
            <Select style={{ width: 300 }} placeholder="请选择项目" value={selectedProjectId} onChange={setSelectedProjectId} options={projects.map(p => ({ value: p.id, label: p.name }))} />
          </div>
        </Card>
      </div>

      <Card style={{ borderRadius: 16, border: 'none', boxShadow: '0 2px 12px rgba(0,0,0,0.08)' }} styles={{ body: { padding: 0 } }}>
        <div style={{ padding: '20px 24px 0', borderBottom: '1px solid #f0f0f0', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <h3 style={{ margin: 0, fontSize: 16, fontWeight: 600 }}>定时调度任务</h3>
          <Button type="primary" icon={<PlusOutlined />} onClick={() => { setEditing(null); form.resetFields(); setModalOpen(true) }} style={{ borderRadius: 8 }}>新建任务</Button>
        </div>
        <Table dataSource={tasks} columns={columns} rowKey="id" loading={loading} pagination={{ pageSize: 10 }} style={{ margin: '0 24px 24px' }} />
      </Card>

      <Modal
        title={editing ? '编辑调度任务' : '新建调度任务'}
        open={modalOpen}
        onOk={handleSave}
        onCancel={() => setModalOpen(false)}
        confirmLoading={saving}
        okButtonProps={{ style: { borderRadius: 8 } }}
        cancelButtonProps={{ style: { borderRadius: 8 } }}
      >
        <Form form={form} layout="vertical" style={{ marginTop: 16 }}>
          <Form.Item name="target_type" label="任务类型" rules={[{ required: true }]}>
            <Select options={[{ value: 'test_case', label: '测试用例' }, { value: 'test_suite', label: '测试套件' }]} />
          </Form.Item>
          <Form.Item name="target_id" label="任务 ID" rules={[{ required: true }]}>
            <InputNumber min={1} placeholder="输入测试用例或套件的 ID" style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item name="cron_expr" label="Cron 表达式" rules={[{ required: true }]}>
            <Input placeholder="* * * * *" />
          </Form.Item>
          <Form.Item label="快速选择">
            <Select placeholder="选择预设" onChange={(v) => form.setFieldValue('cron_expr', v)} options={cronPresets} allowClear />
          </Form.Item>
          <Form.Item name="enabled" label="启用" valuePropName="checked" initialValue={true}>
            <Switch />
          </Form.Item>
          <Form.Item name="environment_id" label="环境 ID">
            <InputNumber min={0} placeholder="可选" style={{ width: '100%' }} />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}
