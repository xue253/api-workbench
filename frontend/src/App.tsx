import { useState, useEffect } from 'react'
import { ConfigProvider, Layout, Menu, Card, Table, Button, Modal, Form, Input, Space, Popconfirm, message, theme } from 'antd'
import { ProjectOutlined, ApiOutlined, SettingOutlined } from '@ant-design/icons'
import { useProjectStore, Project } from './stores/projectStore'

const { Sider, Content } = Layout

const pinkTokens = {
  colorPrimary: '#f759ab',
  colorBgContainer: '#fff0f6',
  colorBgLayout: '#fff5f7',
  colorBorder: '#ffadd2',
  borderRadius: 12,
  colorSuccess: '#52c41a',
  colorWarning: '#faad14',
  colorError: '#ff4d4f',
  fontFamily: "'Noto Sans SC', 'PingFang SC', 'Microsoft YaHei', sans-serif",
}

function ProjectManager() {
  const { projects, loading, fetch: fetchProjects, create, update, remove } = useProjectStore()
  const [modalOpen, setModalOpen] = useState(false)
  const [editing, setEditing] = useState<Project | null>(null)
  const [form] = Form.useForm()

  useEffect(() => { fetchProjects() }, [])

  const handleSave = async () => {
    const values = await form.validateFields()
    if (editing) {
      await update(editing.id, values)
      message.success('更新成功 ✿')
    } else {
      await create(values)
      message.success('创建成功 ✿')
    }
    setModalOpen(false)
    setEditing(null)
    form.resetFields()
  }

  const columns = [
    { title: '名称', dataIndex: 'name', key: 'name' },
    { title: '描述', dataIndex: 'description', key: 'description', ellipsis: true },
    { title: '创建时间', dataIndex: 'created_at', key: 'created_at', width: 200 },
    {
      title: '操作', key: 'action', width: 150,
      render: (_: any, record: Project) => (
        <Space>
          <Button size="small" style={{ color: '#f759ab', borderColor: '#f759ab' }} onClick={() => { setEditing(record); form.setFieldsValue(record); setModalOpen(true) }}>编辑</Button>
          <Popconfirm title="确认删除？" onConfirm={async () => { await remove(record.id); message.success('已删除') }}>
            <Button size="small" danger>删除</Button>
          </Popconfirm>
        </Space>
      )
    }
  ]

  return (
    <Card
      title="📋 项目管理"
      style={{ borderRadius: 16, border: '2px solid #ffadd2' }}
      extra={<Button type="primary" style={{ borderRadius: 8 }} onClick={() => { setEditing(null); form.resetFields(); setModalOpen(true) }}>+ 新建项目</Button>}
    >
      <Table dataSource={projects} columns={columns} rowKey="id" loading={loading} />
      <Modal title={editing ? '✏️ 编辑项目' : '🌟 新建项目'} open={modalOpen} onOk={handleSave} onCancel={() => setModalOpen(false)}>
        <Form form={form} layout="vertical">
          <Form.Item name="name" label="项目名称" rules={[{ required: true }]}><Input placeholder="给项目起个名字吧~" /></Form.Item>
          <Form.Item name="description" label="描述"><Input.TextArea rows={3} placeholder="简单描述一下项目内容~" /></Form.Item>
        </Form>
      </Modal>
    </Card>
  )
}

function Placeholder({ title }: { title: string }) {
  return (
    <Card style={{ borderRadius: 16, border: '2px solid #ffadd2' }}>
      <div style={{ textAlign: 'center', padding: 60 }}>
        <div style={{ fontSize: 48, marginBottom: 16 }}>🌸</div>
        <div style={{ fontSize: 16, color: '#f759ab' }}>{title}</div>
      </div>
    </Card>
  )
}

function App() {
  const [activeKey, setActiveKey] = useState('projects')

  const menuItems = [
    { key: 'projects', icon: <ProjectOutlined />, label: '项目管理' },
    { key: 'apis', icon: <ApiOutlined />, label: '接口库' },
    { key: 'settings', icon: <SettingOutlined />, label: '设置' },
  ]

  const renderContent = () => {
    switch (activeKey) {
      case 'projects': return <ProjectManager />
      case 'apis': return <Placeholder title="接口库 — 开发中" />
      case 'settings': return <Placeholder title="设置 — 开发中" />
      default: return null
    }
  }

  return (
    <ConfigProvider
      theme={{
        algorithm: theme.defaultAlgorithm,
        token: pinkTokens,
      }}
    >
      <Layout style={{ minHeight: '100vh', background: '#fff5f7' }}>
        <Sider
          width={220}
          style={{
            background: 'linear-gradient(180deg, #fff0f6 0%, #ffd6e7 100%)',
            borderRight: '2px solid #ffadd2',
          }}
        >
          <div style={{
            padding: '20px 24px',
            fontSize: 20,
            fontWeight: 'bold',
            color: '#f759ab',
            borderBottom: '1px solid #ffadd2',
          }}>
            ✿ API Workbench
          </div>
          <Menu
            mode="inline"
            selectedKeys={[activeKey]}
            items={menuItems}
            onClick={({ key }) => setActiveKey(key)}
            style={{
              background: 'transparent',
              borderRight: 0,
              marginTop: 8,
            }}
          />
        </Sider>
        <Layout>
          <Content style={{ padding: 24 }}>
            {renderContent()}
          </Content>
        </Layout>
      </Layout>
    </ConfigProvider>
  )
}

export default App
