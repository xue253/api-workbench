import { useState, useEffect } from 'react'
import { ConfigProvider, Layout, Menu, Card, Table, Button, Modal, Form, Input, Space, Popconfirm, message, theme } from 'antd'
import { ProjectOutlined, ApiOutlined, SettingOutlined } from '@ant-design/icons'
import { useProjectStore, Project } from './stores/projectStore'

const { Sider, Content } = Layout

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
    { title: '名称', dataIndex: 'name', key: 'name' },
    { title: '描述', dataIndex: 'description', key: 'description', ellipsis: true },
    { title: '创建时间', dataIndex: 'created_at', key: 'created_at', width: 200 },
    {
      title: '操作', key: 'action', width: 150,
      render: (_: any, record: Project) => (
        <Space>
          <Button size="small" onClick={() => { setEditing(record); form.setFieldsValue(record); setModalOpen(true) }}>编辑</Button>
          <Popconfirm title="确认删除？" onConfirm={async () => { await remove(record.id); message.success('已删除') }}>
            <Button size="small" danger>删除</Button>
          </Popconfirm>
        </Space>
      )
    }
  ]

  return (
    <Card
      title="项目管理"
      extra={<Button type="primary" onClick={() => { setEditing(null); form.resetFields(); setModalOpen(true) }}>新建项目</Button>}
    >
      <Table dataSource={projects} columns={columns} rowKey="id" loading={loading} />
      <Modal title={editing ? '编辑项目' : '新建项目'} open={modalOpen} onOk={handleSave} onCancel={() => setModalOpen(false)}>
        <Form form={form} layout="vertical">
          <Form.Item name="name" label="项目名称" rules={[{ required: true }]}><Input /></Form.Item>
          <Form.Item name="description" label="描述"><Input.TextArea rows={3} /></Form.Item>
        </Form>
      </Modal>
    </Card>
  )
}

function Placeholder({ title }: { title: string }) {
  return <Card><div style={{ textAlign: 'center', padding: 60, color: '#666' }}>{title}</div></Card>
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
      case 'apis': return <Placeholder title="接口库 - 开发中" />
      case 'settings': return <Placeholder title="设置 - 开发中" />
      default: return null
    }
  }

  return (
    <ConfigProvider
      theme={{
        algorithm: theme.darkAlgorithm,
        token: { colorPrimary: '#1677ff', borderRadius: 8 },
      }}
    >
      <Layout style={{ minHeight: '100vh' }}>
        <Sider width={200} style={{ background: '#141414' }}>
          <div style={{ padding: '16px 20px', color: '#fff', fontSize: 16, fontWeight: 'bold' }}>
            API Workbench
          </div>
          <Menu
            mode="inline"
            selectedKeys={[activeKey]}
            items={menuItems}
            onClick={({ key }) => setActiveKey(key)}
            style={{ background: 'transparent', borderRight: 0 }}
          />
        </Sider>
        <Layout>
          <Content style={{ padding: 24, background: '#1a1a1a' }}>
            {renderContent()}
          </Content>
        </Layout>
      </Layout>
    </ConfigProvider>
  )
}

export default App
