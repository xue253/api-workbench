import { useState, useEffect } from 'react'
import { ConfigProvider, Layout, Menu, Card, Table, Button, Modal, Form, Input, Space, Popconfirm, message, theme } from 'antd'
import { ProjectOutlined, ApiOutlined, SettingOutlined, LogoutOutlined, UserOutlined, PlusOutlined } from '@ant-design/icons'
import { useProjectStore, Project } from './stores/projectStore'
import { useAuthStore } from './stores/authStore'
import LoginPage from './pages/Login'

const { Sider, Content } = Layout

const appleTokens = {
  colorPrimary: '#0071e3',
  colorBgContainer: '#ffffff',
  colorBgLayout: '#f5f5f7',
  colorBorder: '#d2d2d7',
  borderRadius: 12,
  fontFamily: "'Inter', -apple-system, BlinkMacSystemFont, 'SF Pro Display', 'SF Pro Text', 'Helvetica Neue', 'Helvetica', 'Arial', sans-serif",
  colorText: '#1d1d1f',
  colorTextSecondary: '#86868b',
  fontSize: 14,
  lineHeight: 1.5,
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
      message.success('更新成功')
    } else {
      await create(values)
      message.success('创建成功')
    }
    setModalOpen(false)
    setEditing(null)
    form.resetFields()
  }

  const columns = [
    { title: '名称', dataIndex: 'name', key: 'name', render: (text: string) => <span style={{ fontWeight: 500 }}>{text}</span> },
    { title: '描述', dataIndex: 'description', key: 'description', ellipsis: true },
    { title: '创建时间', dataIndex: 'created_at', key: 'created_at', width: 200 },
    {
      title: '操作', key: 'action', width: 150,
      render: (_: any, record: Project) => (
        <Space>
          <Button 
            size="small" 
            type="text"
            style={{ color: '#0071e3' }} 
            onClick={() => { setEditing(record); form.setFieldsValue(record); setModalOpen(true) }}
          >
            编辑
          </Button>
          <Popconfirm title="确认删除？" onConfirm={async () => { await remove(record.id); message.success('已删除') }}>
            <Button size="small" type="text" danger>删除</Button>
          </Popconfirm>
        </Space>
      )
    }
  ]

  return (
    <Card 
      style={{ 
        borderRadius: 16, 
        border: 'none',
        boxShadow: '0 2px 12px rgba(0,0,0,0.08)',
      }}
      styles={{ body: { padding: 0 } }}
    >
      <div style={{ 
        padding: '24px 24px 0', 
        borderBottom: '1px solid #f0f0f0',
        display: 'flex',
        justifyContent: 'space-between',
        alignItems: 'center'
      }}>
        <div>
          <h2 style={{ margin: 0, fontSize: 20, fontWeight: 600, color: '#1d1d1f' }}>项目管理</h2>
          <p style={{ margin: '4px 0 16px', fontSize: 14, color: '#86868b' }}>管理你的 API 测试项目</p>
        </div>
        <Button 
          type="primary" 
          icon={<PlusOutlined />}
          onClick={() => { setEditing(null); form.resetFields(); setModalOpen(true) }}
          style={{ borderRadius: 8 }}
        >
          新建项目
        </Button>
      </div>
      <Table 
        dataSource={projects} 
        columns={columns} 
        rowKey="id" 
        loading={loading}
        pagination={{ pageSize: 10 }}
        style={{ margin: '0 24px 24px' }}
      />
      <Modal 
        title={editing ? '编辑项目' : '新建项目'} 
        open={modalOpen} 
        onOk={handleSave} 
        onCancel={() => setModalOpen(false)}
        okButtonProps={{ style: { borderRadius: 8 } }}
        cancelButtonProps={{ style: { borderRadius: 8 } }}
      >
        <Form form={form} layout="vertical" style={{ marginTop: 16 }}>
          <Form.Item name="name" label="项目名称" rules={[{ required: true }]}>
            <Input placeholder="输入项目名称" />
          </Form.Item>
          <Form.Item name="description" label="描述">
            <Input.TextArea rows={3} placeholder="项目描述（可选）" />
          </Form.Item>
        </Form>
      </Modal>
    </Card>
  )
}

function Placeholder({ title }: { title: string }) {
  return (
    <Card style={{ 
      borderRadius: 16, 
      border: 'none',
      boxShadow: '0 2px 12px rgba(0,0,0,0.08)',
    }}>
      <div style={{ textAlign: 'center', padding: 80 }}>
        <div style={{ 
          width: 64, 
          height: 64, 
          borderRadius: 16, 
          background: '#f5f5f7',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          margin: '0 auto 24px'
        }}>
          <ApiOutlined style={{ fontSize: 28, color: '#86868b' }} />
        </div>
        <div style={{ fontSize: 20, fontWeight: 600, color: '#1d1d1f', marginBottom: 8 }}>{title}</div>
        <div style={{ fontSize: 14, color: '#86868b' }}>功能开发中</div>
      </div>
    </Card>
  )
}

function App() {
  const [activeKey, setActiveKey] = useState('projects')
  const { user, token, fetchProfile, logout } = useAuthStore()
  const [showLogin, setShowLogin] = useState(false)

  useEffect(() => {
    if (token) fetchProfile()
  }, [token])

  if (!token || showLogin) {
    return (
      <ConfigProvider theme={{ algorithm: theme.defaultAlgorithm, token: appleTokens }}>
        <LoginPage onLogin={() => setShowLogin(false)} />
      </ConfigProvider>
    )
  }

  const menuItems = [
    { key: 'projects', icon: <ProjectOutlined />, label: '项目' },
    { key: 'apis', icon: <ApiOutlined />, label: '接口库' },
    { key: 'settings', icon: <SettingOutlined />, label: '设置' },
  ]

  const renderContent = () => {
    switch (activeKey) {
      case 'projects': return <ProjectManager />
      case 'apis': return <Placeholder title="接口库" />
      case 'settings': return <Placeholder title="设置" />
      default: return null
    }
  }

  return (
    <ConfigProvider theme={{ algorithm: theme.defaultAlgorithm, token: appleTokens }}>
      <Layout style={{ minHeight: '100vh', background: '#f5f5f7' }}>
        <Sider
          width={240}
          style={{
            background: '#ffffff',
            borderRight: '1px solid #e5e5e5',
          }}
        >
          <div style={{
            padding: '24px 20px',
            borderBottom: '1px solid #e5e5e5',
          }}>
            <div style={{ 
              fontSize: 18, 
              fontWeight: 600, 
              color: '#1d1d1f',
              letterSpacing: '-0.3px'
            }}>
              API Workbench
            </div>
          </div>
          <Menu
            mode="inline"
            selectedKeys={[activeKey]}
            items={menuItems}
            onClick={({ key }) => setActiveKey(key)}
            style={{ 
              background: 'transparent', 
              borderRight: 0, 
              padding: '8px 0',
            }}
          />
          <div style={{
            position: 'absolute',
            bottom: 0,
            left: 0,
            right: 0,
            padding: '16px 20px',
            borderTop: '1px solid #e5e5e5',
            background: '#ffffff',
          }}>
            <div style={{ 
              display: 'flex', 
              alignItems: 'center', 
              gap: 12, 
              marginBottom: 12 
            }}>
              <div style={{
                width: 36,
                height: 36,
                borderRadius: '50%',
                background: '#f5f5f7',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
              }}>
                <UserOutlined style={{ color: '#86868b', fontSize: 16 }} />
              </div>
              <div>
                <div style={{ fontSize: 14, fontWeight: 500, color: '#1d1d1f' }}>{user?.username}</div>
                <div style={{ fontSize: 12, color: '#86868b' }}>在线</div>
              </div>
            </div>
            <Button
              size="small"
              icon={<LogoutOutlined />}
              onClick={() => { logout(); setShowLogin(true) }}
              style={{ 
                color: '#86868b', 
                borderColor: '#d2d2d7',
                borderRadius: 8,
              }}
              block
            >
              退出登录
            </Button>
          </div>
        </Sider>
        <Layout>
          <Content style={{ padding: 32 }}>
            {renderContent()}
          </Content>
        </Layout>
      </Layout>
    </ConfigProvider>
  )
}

export default App
