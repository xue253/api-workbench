import { useEffect } from 'react'
import { Card, Table, Button, Modal, Form, Input, Space, Popconfirm, message } from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons'
import { useProjectStore, Project } from '../stores/projectStore'

export default function ProjectList() {
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

  const handleDelete = async (id: number) => {
    await remove(id)
    message.success('删除成功')
  }

  const columns = [
    { title: '名称', dataIndex: 'name', key: 'name' },
    { title: '描述', dataIndex: 'description', key: 'description', ellipsis: true },
    { title: '创建时间', dataIndex: 'created_at', key: 'created_at', width: 200 },
    {
      title: '操作', key: 'action', width: 150,
      render: (_: any, record: Project) => (
        <Space>
          <Button size="small" icon={<EditOutlined />} onClick={() => { setEditing(record); form.setFieldsValue(record); setModalOpen(true) }} />
          <Popconfirm title="确认删除？" onConfirm={() => handleDelete(record.id)}>
            <Button size="small" danger icon={<DeleteOutlined />} />
          </Popconfirm>
        </Space>
      )
    }
  ]

  return (
    <Card
      title="项目管理"
      extra={
        <Button type="primary" icon={<PlusOutlined />} onClick={() => { setEditing(null); form.resetFields(); setModalOpen(true) }}>
          新建项目
        </Button>
      }
    >
      <Table dataSource={projects} columns={columns} rowKey="id" loading={loading} />
      <Modal title={editing ? '编辑项目' : '新建项目'} open={modalOpen} onOk={handleSave} onCancel={() => setModalOpen(false)}>
        <Form form={form} layout="vertical">
          <Form.Item name="name" label="项目名称" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="description" label="描述">
            <Input.TextArea rows={3} />
          </Form.Item>
        </Form>
      </Modal>
    </Card>
  )
}

import { useState } from 'react'
