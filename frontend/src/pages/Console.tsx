import { useState, useEffect, useRef } from 'react'
import { Card, Select, Button, Table, Tag, Progress, Space, message, Empty } from 'antd'
import { PlayCircleOutlined, ReloadOutlined, CheckCircleOutlined, CloseCircleOutlined, LoadingOutlined } from '@ant-design/icons'
import { useProjectStore } from '../stores/projectStore'
import api from '../services/api'

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

interface RunDetail {
  index: number
  total: number
  passed: number
  failed: number
  api_name: string
  status: string
}

export default function ConsolePage() {
  const { projects, fetch: fetchProjects } = useProjectStore()
  const [selectedProjectId, setSelectedProjectId] = useState<number | null>(null)
  const [testCases, setTestCases] = useState<any[]>([])
  const [testSuites, setTestSuites] = useState<any[]>([])
  const [currentRun, setCurrentRun] = useState<TestRun | null>(null)
  const [runDetails, setRunDetails] = useState<RunDetail[]>([])
  const [isRunning, setIsRunning] = useState(false)
  const wsRef = useRef<WebSocket | null>(null)

  useEffect(() => { fetchProjects() }, []) // eslint-disable-line react-hooks/exhaustive-deps

  useEffect(() => {
    if (selectedProjectId) {
      loadData()
    }
  }, [selectedProjectId])

  useEffect(() => {
    return () => {
      if (wsRef.current) {
        wsRef.current.close()
      }
    }
  }, [])

  const loadData = async () => {
    try {
      const [casesRes, suitesRes]: any[] = await Promise.all([
        api.get(`/projects/${selectedProjectId}/test-cases`),
        api.get(`/projects/${selectedProjectId}/test-suites`),
      ])
      setTestCases(casesRes.data || [])
      setTestSuites(suitesRes.data || [])
    } catch (err) {
      console.error(err)
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
      const run = res.data
      setCurrentRun({ ...run, passed: 0, failed: 0, total: 0 })
      setRunDetails([])
      setIsRunning(true)
      connectWebSocket(run.id)
    } catch (err: any) {
      message.error(err.message || '启动失败')
    }
  }

  const connectWebSocket = (runId: number) => {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const ws = new WebSocket(`${protocol}//${window.location.host}/api/v1/ws/test-runs/${runId}`)
    wsRef.current = ws

    ws.onmessage = (event) => {
      const msg = JSON.parse(event.data)
      
      if (msg.type === 'progress') {
        setCurrentRun(prev => prev ? { ...prev, ...msg.data } : null)
      } else if (msg.type === 'detail') {
        setRunDetails(prev => [...prev, msg.data])
        setCurrentRun(prev => prev ? { ...prev, passed: msg.data.passed, failed: msg.data.failed } : null)
      } else if (msg.type === 'complete') {
        setCurrentRun(prev => prev ? { ...prev, status: msg.data.status, duration_ms: msg.data.duration_ms } : null)
        setIsRunning(false)
        ws.close()
        message.success('测试执行完成')
      }
    }

    ws.onerror = () => {
      setIsRunning(false)
      message.error('WebSocket 连接错误')
    }

    ws.onclose = () => {
      setIsRunning(false)
    }
  }

  const columns = [
    { title: '#', dataIndex: 'index', key: 'index', width: 60 },
    { title: '接口', dataIndex: 'api_name', key: 'api_name', ellipsis: true },
    { title: '状态', dataIndex: 'status', key: 'status', width: 80, render: (v: string) => <Tag color={v === 'passed' ? 'green' : 'red'}>{v === 'passed' ? '通过' : '失败'}</Tag> },
    { title: '进度', key: 'progress', width: 200, render: (_: any, r: RunDetail) => <Progress percent={Math.round((r.index / r.total) * 100)} size="small" /> },
  ]

  const passRate = currentRun && currentRun.total > 0
    ? Math.round((currentRun.passed / currentRun.total) * 100)
    : 0

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
          <div style={{ padding: '20px 24px 0', borderBottom: '1px solid #f0f0f0' }}>
            <h3 style={{ margin: 0, fontSize: 16, fontWeight: 600 }}>测试用例</h3>
          </div>
          <div style={{ padding: 16, maxHeight: 200, overflowY: 'auto' }}>
            {testCases.length > 0 ? testCases.map(tc => (
              <div key={tc.id} style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', padding: '8px 0', borderBottom: '1px solid #f0f0f0' }}>
                <span>{tc.name}</span>
                <Button size="small" type="primary" icon={<PlayCircleOutlined />} onClick={() => handleRun('case', tc.id)} disabled={isRunning} style={{ borderRadius: 8 }}>运行</Button>
              </div>
            )) : <Empty description="暂无用例" image={Empty.PRESENTED_IMAGE_SIMPLE} />}
          </div>
        </Card>

        <Card style={{ flex: 1, borderRadius: 16, border: 'none', boxShadow: '0 2px 12px rgba(0,0,0,0.08)' }} styles={{ body: { padding: 0 } }}>
          <div style={{ padding: '20px 24px 0', borderBottom: '1px solid #f0f0f0' }}>
            <h3 style={{ margin: 0, fontSize: 16, fontWeight: 600 }}>测试套件</h3>
          </div>
          <div style={{ padding: 16, maxHeight: 200, overflowY: 'auto' }}>
            {testSuites.length > 0 ? testSuites.map(ts => (
              <div key={ts.id} style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', padding: '8px 0', borderBottom: '1px solid #f0f0f0' }}>
                <span>{ts.name}</span>
                <Button size="small" type="primary" icon={<PlayCircleOutlined />} onClick={() => handleRun('suite', ts.id)} disabled={isRunning} style={{ borderRadius: 8 }}>运行</Button>
              </div>
            )) : <Empty description="暂无套件" image={Empty.PRESENTED_IMAGE_SIMPLE} />}
          </div>
        </Card>
      </div>

      {currentRun && (
        <Card style={{ borderRadius: 16, border: 'none', boxShadow: '0 2px 12px rgba(0,0,0,0.08)' }} styles={{ body: { padding: 0 } }}>
          <div style={{ padding: '20px 24px 0', borderBottom: '1px solid #f0f0f0', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <div style={{ display: 'flex', alignItems: 'center', gap: 16 }}>
              <h3 style={{ margin: 0, fontSize: 16, fontWeight: 600 }}>执行进度</h3>
              {isRunning ? (
                <Tag icon={<LoadingOutlined />} color="processing">运行中</Tag>
              ) : (
                <Tag icon={currentRun.status === 'done' ? <CheckCircleOutlined /> : <CloseCircleOutlined />} color={currentRun.status === 'done' ? 'success' : 'error'}>
                  {currentRun.status === 'done' ? '完成' : '失败'}
                </Tag>
              )}
            </div>
            <Space>
              <span>通过: <span style={{ color: '#52c41a', fontWeight: 600 }}>{currentRun.passed}</span></span>
              <span>失败: <span style={{ color: '#ff4d4f', fontWeight: 600 }}>{currentRun.failed}</span></span>
              <span>通过率: <span style={{ fontWeight: 600 }}>{passRate}%</span></span>
            </Space>
          </div>
          <Table
            dataSource={runDetails}
            columns={columns}
            rowKey="index"
            pagination={false}
            size="small"
            style={{ margin: '0 24px 24px' }}
          />
        </Card>
      )}

      {!currentRun && (
        <Card style={{ borderRadius: 16, border: 'none', boxShadow: '0 2px 12px rgba(0,0,0,0.08)' }}>
          <div style={{ textAlign: 'center', padding: 60 }}>
            <div style={{ fontSize: 16, color: '#86868b' }}>选择一个用例或套件开始执行</div>
          </div>
        </Card>
      )}
    </div>
  )
}
