import { useState, useEffect } from 'react'
import { Card, Table, Button, Select, Tag, Space, Statistic, Row, Col, Empty, Spin, message } from 'antd'
import { DownloadOutlined, ReloadOutlined, FileTextOutlined } from '@ant-design/icons'
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
  skipped: number
  duration_ms: number
  started_at: string
  finished_at: string
}

interface RunDetail {
  id: number
  api_id: number
  test_case_id: number
  status: string
  status_code: number
  duration_ms: number
  error_message: string
  response_body: string
}

export default function ReportPage() {
  const { projects, fetch: fetchProjects } = useProjectStore()
  const [selectedProjectId, setSelectedProjectId] = useState<number | null>(null)
  const [testRuns, setTestRuns] = useState<TestRun[]>([])
  const [loading, setLoading] = useState(false)
  const [selectedRun, setSelectedRun] = useState<TestRun | null>(null)
  const [runDetails, setRunDetails] = useState<RunDetail[]>([])
  const [detailsLoading, setDetailsLoading] = useState(false)

  useEffect(() => { fetchProjects() }, []) // eslint-disable-line react-hooks/exhaustive-deps

  useEffect(() => {
    if (selectedProjectId) {
      loadTestRuns()
    }
  }, [selectedProjectId])

  const loadTestRuns = async () => {
    setLoading(true)
    try {
      const res: any = await api.get('/test-runs')
      setTestRuns(res.data || [])
    } catch (err) {
      console.error(err)
    }
    setLoading(false)
  }

  const loadRunDetails = async (run: TestRun) => {
    setSelectedRun(run)
    setDetailsLoading(true)
    try {
      const res: any = await api.get(`/test-runs/${run.id}/report`)
      setRunDetails(res.data?.details || [])
    } catch (err) {
      console.error(err)
    }
    setDetailsLoading(false)
  }

  const handleExport = async (runId: number, format: string) => {
    try {
      const res = await api.get(`/test-runs/${runId}/export?format=${format}`, { responseType: 'blob' })
      const blob = new Blob([res as any])
      const url = window.URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = `report-${runId}.${format}`
      a.click()
      window.URL.revokeObjectURL(url)
      message.success('导出成功')
    } catch (err: any) {
      message.error(err.message || '导出失败')
    }
  }

  const runColumns = [
    { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
    { title: '类型', dataIndex: 'target_type', key: 'target_type', width: 100, render: (v: string) => <Tag color={v === 'test_suite' ? 'blue' : 'green'}>{v === 'test_suite' ? '套件' : '用例'}</Tag> },
    { title: '状态', dataIndex: 'status', key: 'status', width: 80, render: (v: string) => <Tag color={v === 'done' ? 'green' : v === 'failed' ? 'red' : 'blue'}>{v === 'done' ? '通过' : v === 'failed' ? '失败' : '运行中'}</Tag> },
    { title: '通过/失败', key: 'result', width: 100, render: (_: any, r: TestRun) => <span><span style={{ color: '#52c41a' }}>{r.passed}</span> / <span style={{ color: '#ff4d4f' }}>{r.failed}</span></span> },
    { title: '耗时', dataIndex: 'duration_ms', key: 'duration_ms', width: 100, render: (v: number) => `${v}ms` },
    { title: '时间', dataIndex: 'started_at', key: 'started_at', width: 180 },
    {
      title: '操作', key: 'action', width: 200,
      render: (_: any, record: TestRun) => (
        <Space>
          <Button size="small" type="text" style={{ color: '#0071e3' }} onClick={() => loadRunDetails(record)}>
            查看详情
          </Button>
          <Button size="small" type="text" icon={<DownloadOutlined />} onClick={() => handleExport(record.id, 'md')}>
            MD
          </Button>
          <Button size="small" type="text" icon={<DownloadOutlined />} onClick={() => handleExport(record.id, 'html')}>
            HTML
          </Button>
        </Space>
      )
    }
  ]

  const detailColumns = [
    { title: '#', key: 'index', width: 50, render: (_: any, __: any, i: number) => i + 1 },
    { title: '状态', dataIndex: 'status', key: 'status', width: 80, render: (v: string) => <Tag color={v === 'passed' ? 'green' : 'red'}>{v === 'passed' ? '通过' : '失败'}</Tag> },
    { title: '状态码', dataIndex: 'status_code', key: 'status_code', width: 80 },
    { title: '耗时', dataIndex: 'duration_ms', key: 'duration_ms', width: 100, render: (v: number) => `${v}ms` },
    { title: '错误信息', dataIndex: 'error_message', key: 'error_message', ellipsis: true },
  ]

  const passRate = selectedRun && selectedRun.total > 0
    ? ((selectedRun.passed / selectedRun.total) * 100).toFixed(1)
    : '0'

  return (
    <div>
      <div style={{ marginBottom: 24 }}>
        <Card style={{ borderRadius: 16, border: 'none', boxShadow: '0 2px 12px rgba(0,0,0,0.08)' }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: 16 }}>
            <span style={{ fontSize: 14, color: '#86868b', whiteSpace: 'nowrap' }}>选择项目：</span>
            <Select style={{ width: 300 }} placeholder="请选择项目" value={selectedProjectId} onChange={setSelectedProjectId} options={projects.map(p => ({ value: p.id, label: p.name }))} />
            <Button icon={<ReloadOutlined />} onClick={loadTestRuns} style={{ borderRadius: 8 }}>刷新</Button>
          </div>
        </Card>
      </div>

      {selectedRun && (
        <div style={{ marginBottom: 24 }}>
          <Card style={{ borderRadius: 16, border: 'none', boxShadow: '0 2px 12px rgba(0,0,0,0.08)' }}>
            <Row gutter={24}>
              <Col span={6}>
                <Statistic title="总数" value={selectedRun.total} />
              </Col>
              <Col span={6}>
                <Statistic title="通过" value={selectedRun.passed} valueStyle={{ color: '#52c41a' }} />
              </Col>
              <Col span={6}>
                <Statistic title="失败" value={selectedRun.failed} valueStyle={{ color: '#ff4d4f' }} />
              </Col>
              <Col span={6}>
                <Statistic title="通过率" value={`${passRate}%`} suffix="" />
              </Col>
            </Row>
          </Card>
        </div>
      )}

      <div style={{ display: 'flex', gap: 24 }}>
        <div style={{ width: selectedRun ? 400 : '100%', flexShrink: 0 }}>
          <Card style={{ borderRadius: 16, border: 'none', boxShadow: '0 2px 12px rgba(0,0,0,0.08)' }} styles={{ body: { padding: 0 } }}>
            <div style={{ padding: '20px 24px 0', borderBottom: '1px solid #f0f0f0' }}>
              <h3 style={{ margin: 0, fontSize: 16, fontWeight: 600 }}>执行记录</h3>
            </div>
            <Table
              dataSource={testRuns}
              columns={runColumns}
              rowKey="id"
              loading={loading}
              pagination={{ pageSize: 10 }}
              size="small"
              onRow={(record) => ({
                onClick: () => loadRunDetails(record),
                style: { cursor: 'pointer', background: selectedRun?.id === record.id ? '#f0f5ff' : undefined }
              })}
              style={{ margin: '0 16px 16px' }}
            />
          </Card>
        </div>

        {selectedRun && (
          <div style={{ flex: 1 }}>
            <Card style={{ borderRadius: 16, border: 'none', boxShadow: '0 2px 12px rgba(0,0,0,0.08)' }} styles={{ body: { padding: 0 } }}>
              <div style={{ padding: '20px 24px 0', borderBottom: '1px solid #f0f0f0', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <h3 style={{ margin: 0, fontSize: 16, fontWeight: 600 }}>执行详情</h3>
                <Space>
                  <Button size="small" icon={<DownloadOutlined />} onClick={() => handleExport(selectedRun.id, 'md')} style={{ borderRadius: 8 }}>导出 MD</Button>
                  <Button size="small" icon={<DownloadOutlined />} onClick={() => handleExport(selectedRun.id, 'html')} style={{ borderRadius: 8 }}>导出 HTML</Button>
                </Space>
              </div>
              <Table
                dataSource={runDetails}
                columns={detailColumns}
                rowKey="id"
                loading={detailsLoading}
                pagination={false}
                size="small"
                style={{ margin: '0 16px 16px' }}
              />
            </Card>
          </div>
        )}
      </div>
    </div>
  )
}
