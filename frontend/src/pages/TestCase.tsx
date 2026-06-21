import { useState, useEffect } from 'react'
import { Card, Table, Button, Modal, Form, Input, Space, Popconfirm, message, Select, Tag, Empty } from 'antd'
import { PlusOutlined, PlayCircleOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons'
import { useProjectStore } from '../stores/projectStore'
import api from '../services/api'

interface TestCase {
  id: number
  project_id: number
  name: string
  description: string
  created_at: string
}

interface TestSuite {
  id: number
  project_id: number
  name: string
  description: string
  run_mode: string
  max_concurrency: number
  created_at: string
}

interface TestRun {
  id: number
  target_type: string
  target_id: number
  status: string
  total: number
  passed: number
  failed: number
  duration_ms: number
  started_at: string
}

export default function TestCasePage() {
  const { projects, fetch: fetchProjects } = useProjectStore()
  const [selectedProjectId, setSelectedProjectId] = useState<number | null>(null)
  const [testCases, setTestCases] = useState<TestCase[]>([])
  const [testSuites, setTestSuites] = useState<TestSuite[]>([])
  const [testRuns, setTestRuns] = useState<TestRun[]>([])
  const [loading, setLoading] = useState(false)
  const [modalOpen, setModalOpen] = useState(false)
  const [modalType, setModalType] = useState<'case' | 'suite'>('case')
  const [editing, setEditing] = useState<any>(null)
  const [saving, setSaving] = useState(false)
  const [form] = Form.useForm()

  useEffect(() => { fetchProjects() }, [])

  useEffect(() => {
    if (selectedProjectId) {
      loadData()
    }
  }, [selectedProjectId])

  const loadData = async () => {
    setLoading(true)
    try {
      const [casesRes, suitesRes, runsRes]: any[] = await Promise.all([
        api.get(`/projects/${selectedProjectId}/test-cases`),
        api.get(`/projects/${selectedProjectId}/test-suites`),
        api.get('/test-runs'),
      ])
      setTestCases(casesRes.data || [])
      setTestSuites(suitesRes.data || [])
      setTestRuns(runsRes.data || [])
    } catch (err) {
      console.error(err)
    }
    setLoading(false)
  }

  const handleSave = async () => {
    try {
      const values = await form.validateFields()
      setSaving(true)
      if (modalType === 'case') {
        if (editing) {
          await api.put(`/test-cases/${editing.id}`, values)
        } else {
          await api.post(`/projects/${selectedProjectId}/test-cases`, values)
        }
      } else {
        if (editing) {
          await api.put(`/test-suites/${editing.id}`, values)
        } else {
          await api.post(`/projects/${selectedProjectId}/test-suites`, values)
        }
      }
      message.success('保存成功')
      setModalOpen(false)
      setEditing(null)
      form.resetFields()
      loadData()
    } catch (err: any) {
      if (err.errorFields) return
      message.error(err.message || '操作失败')
    } finally {
      setSaving(false)
    }
  }

  const handleDelete = async (type: string, id: number) => {
    try {
      if (type === 'case') {
        await api.delete(`/test-cases/${id}`)
      } else {
        await api.delete(`/test-suites/${id}`)
      }
      message.success('删除成功')
      loadData()
    } catch (err: any) {
      message.error(err.message || '删除失败')
    }
  }

  const handleRun = async (type: string, id: number) => {
    try {
      let res: any
      if (type === 'case') {
        res = await api.post(`/test-cases/${id}/run`, {})
      } else {
        res = await api.post(`/test-suites/${id}/run`, {})
      }
      message.success('测试已启动')
      loadData()
    } catch (err: any) {
      message.error(err.message || '启动失败')
    }
  }

  const caseColumns = [
    { title: '名称', dataIndex: 'name', key: 'name', render: (text: string) => <span style={{ fontWeight: 500 }}>{text}</span> },
    { title: '描述', dataIndex: 'description', key: 'description', ellipsis: true },
    { title: '创建时间', dataIndex: 'created_at', key: 'created_at', width: 180 },
    {
      title: '操作', key: 'action', width: 200,
      render: (_: any, record: TestCase) => (
        <Space>
          <Button size="small" type="text" style={{ color: '#52c41a' }} onClick={() => handleRun('case', record.id)}>
            <PlayCircleOutlined /> 运行
          </Button>
          <Button size="small" type="text" style={{ color: '#0071e3' }} onClick={() => { setModalType('case'); setEditing(record); form.setFieldsValue(record); setModalOpen(true) }}>
            <EditOutlined /> 编辑
          </Button>
          <Popconfirm title="确认删除？" onConfirm={() => handleDelete('case', record.id)}>
            <Button size="small" type="text" danger><DeleteOutlined /> 删除</Button>
          </Popconfirm>
        </Space>
      )
    }
  ]

  const suiteColumns = [
    { title: '名称', dataIndex: 'name', key: 'name', render: (text: string) => <span style={{ fontWeight: 500 }}>{text}</span> },
    { title: '描述', dataIndex: 'description', key: 'description', ellipsis: true },
    { title: '执行模式', dataIndex: 'run_mode', key: 'run_mode', width: 100, render: (v: string) => <Tag color={v === 'parallel' ? 'blue' : 'green'}>{v === 'parallel' ? '并行' : '顺序'}</Tag> },
    {
      title: '操作', key: 'action', width: 200,
      render: (_: any, record: TestSuite) => (
        <Space>
          <Button size="small" type="text" style={{ color: '#52c41a' }} onClick={() => handleRun('suite', record.id)}>
            <PlayCircleOutlined /> 运行
          </Button>
          <Button size="small" type="text" style={{ color: '#0071e3' }} onClick={() => { setModalType('suite'); setEditing(record); form.setFieldsValue(record); setModalOpen(true) }}>
            <EditOutlined /> 编辑
          </Button>
          <Popconfirm title="确认删除？" onConfirm={() => handleDelete('suite', record.id)}>
            <Button size="small" type="text" danger><DeleteOutlined /> 删除</Button>
          </Popconfirm>
        </Space>
      )
    }
  ]

  const runColumns = [
    { title: '类型', dataIndex: 'target_type', key: 'target_type', width: 100, render: (v: string) => <Tag color={v === 'test_suite' ? 'blue' : 'green'}>{v === 'test_suite' ? '套件' : '用例'}</Tag> },
    { title: '状态', dataIndex: 'status', key: 'status', width: 100, render: (v: string) => <Tag color={v === 'done' ? 'green' : v === 'failed' ? 'red' : 'blue'}>{v === 'done' ? '通过' : v === 'failed' ? '失败' : '运行中'}</Tag> },
    { title: '通过/失败', key: 'result', width: 120, render: (_: any, r: TestRun) => <span><span style={{ color: '#52c41a' }}>{r.passed}</span> / <span style={{ color: '#ff4d4f' }}>{r.failed}</span></span> },
    { title: '耗时', dataIndex: 'duration_ms', key: 'duration_ms', width: 100, render: (v: number) => `${v}ms` },
    { title: '时间', dataIndex: 'started_at', key: 'started_at', width: 180 },
  ]

  if (!selectedProjectId) {
    return (
      <Card style={{ borderRadius: 16, border: 'none', boxShadow: '0 2px 12px rgba(0,0,0,0.08)' }}>
        <div style={{ textAlign: 'center', padding: 80 }}>
          <div style={{ width: 64, height: 64, borderRadius: 16, background: '#f5f5f7', display: 'flex', alignItems: 'center', justifyContent: 'center', margin: '0 auto 24px' }}>
            <PlayCircleOutlined style={{ fontSize: 28, color: '#86868b' }} />
          </div>
          <div style={{ fontSize: 20, fontWeight: 600, color: '#1d1d1f', marginBottom: 8 }}>请先选择项目</div>
          <div style={{ fontSize: 14, color: '#86868b' }}>选择一个项目来管理测试用例</div>
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

      <div style={{ display: 'flex', gap: 24, marginBottom: 24 }}>
        <Card style={{ flex: 1, borderRadius: 16, border: 'none', boxShadow: '0 2px 12px rgba(0,0,0,0.08)' }} styles={{ body: { padding: 0 } }}>
          <div style={{ padding: '20px 24px 0', borderBottom: '1px solid #f0f0f0', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <h3 style={{ margin: 0, fontSize: 16, fontWeight: 600 }}>测试用例</h3>
            <Button type="primary" size="small" icon={<PlusOutlined />} onClick={() => { setModalType('case'); setEditing(null); form.resetFields(); setModalOpen(true) }} style={{ borderRadius: 8 }}>新建</Button>
          </div>
          <Table dataSource={testCases} columns={caseColumns} rowKey="id" loading={loading} pagination={false} size="small" style={{ margin: '0 16px 16px' }} />
        </Card>

        <Card style={{ flex: 1, borderRadius: 16, border: 'none', boxShadow: '0 2px 12px rgba(0,0,0,0.08)' }} styles={{ body: { padding: 0 } }}>
          <div style={{ padding: '20px 24px 0', borderBottom: '1px solid #f0f0f0', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <h3 style={{ margin: 0, fontSize: 16, fontWeight: 600 }}>测试套件</h3>
            <Button type="primary" size="small" icon={<PlusOutlined />} onClick={() => { setModalType('suite'); setEditing(null); form.resetFields(); setModalOpen(true) }} style={{ borderRadius: 8 }}>新建</Button>
          </div>
          <Table dataSource={testSuites} columns={suiteColumns} rowKey="id" loading={loading} pagination={false} size="small" style={{ margin: '0 16px 16px' }} />
        </Card>
      </div>

      <Card style={{ borderRadius: 16, border: 'none', boxShadow: '0 2px 12px rgba(0,0,0,0.08)' }} styles={{ body: { padding: 0 } }}>
        <div style={{ padding: '20px 24px 0', borderBottom: '1px solid #f0f0f0' }}>
          <h3 style={{ margin: 0, fontSize: 16, fontWeight: 600 }}>执行记录</h3>
        </div>
        <Table dataSource={testRuns} columns={runColumns} rowKey="id" loading={loading} pagination={{ pageSize: 10 }} size="small" style={{ margin: '0 16px 16px' }} />
      </Card>

      <Modal
        title={modalType === 'case' ? (editing ? '编辑测试用例' : '新建测试用例') : (editing ? '编辑测试套件' : '新建测试套件')}
        open={modalOpen}
        onOk={handleSave}
        onCancel={() => setModalOpen(false)}
        confirmLoading={saving}
        okButtonProps={{ style: { borderRadius: 8 } }}
        cancelButtonProps={{ style: { borderRadius: 8 } }}
      >
        <Form form={form} layout="vertical" style={{ marginTop: 16 }}>
          <Form.Item name="name" label="名称" rules={[{ required: true }]}><Input placeholder="输入名称" /></Form.Item>
          <Form.Item name="description" label="描述"><Input.TextArea rows={3} placeholder="描述（可选）" /></Form.Item>
          {modalType === 'suite' && (
            <>
              <Form.Item name="run_mode" label="执行模式" initialValue="sequential">
                <Select options={[{ value: 'sequential', label: '顺序执行' }, { value: 'parallel', label: '并行执行' }]} />
              </Form.Item>
              <Form.Item name="max_concurrency" label="最大并发数" initialValue={5}>
                <Input type="number" placeholder="5" />
              </Form.Item>
            </>
          )}
        </Form>
      </Modal>
    </div>
  )
}
