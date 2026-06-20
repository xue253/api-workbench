import { useState, useEffect } from 'react'
import { Card, Tree, Button, Modal, Form, Input, message, Empty } from 'antd'
import { PlusOutlined, FolderOutlined, ApiOutlined } from '@ant-design/icons'
import { useApiStore, Collection } from '../stores/apiStore'
import { useProjectStore } from '../stores/projectStore'
import type { DataNode } from 'antd/es/tree'

export default function ApiExplorer() {
  const { current } = useProjectStore()
  const { collections, apis, fetchCollections, createCollection, fetchAPIs, setCurrentApi } = useApiStore()
  const [colModalOpen, setColModalOpen] = useState(false)
  const [form] = Form.useForm()

  useEffect(() => {
    if (current) fetchCollections(current.id)
  }, [current])

  const buildTree = (): DataNode[] => {
    const map = new Map<number | null, Collection[]>()
    collections.forEach(c => {
      const key = c.parent_id ?? null
      if (!map.has(key)) map.set(key, [])
      map.get(key)!.push(c)
    })

    const buildNode = (pid: number | null): DataNode[] => {
      return (map.get(pid) || []).map(c => ({
        key: `col-${c.id}`,
        title: c.name,
        icon: <FolderOutlined />,
        children: apis.filter(a => a.collection_id === c.id).map(a => ({
          key: `api-${a.id}`,
          title: a.name,
          icon: <ApiOutlined />,
          isLeaf: true,
        })),
      }))
    }

    return buildNode(null)
  }

  const handleSelect = async (keys: any[]) => {
    const key = keys[0] as string
    if (key?.startsWith('api-')) {
      const id = parseInt(key.replace('api-', ''))
      await setCurrentApi(apis.find(a => a.id === id) || null)
    } else if (key?.startsWith('col-')) {
      const id = parseInt(key.replace('col-', ''))
      await fetchAPIs(id)
    }
  }

  const handleCreateCol = async () => {
    if (!current) return
    const values = await form.validateFields()
    await createCollection(current.id, values)
    message.success('创建成功')
    setColModalOpen(false)
    form.resetFields()
  }

  if (!current) {
    return <Card><Empty description="请先选择项目" /></Card>
  }

  return (
    <Card
      title={`${current.name} - 接口库`}
      extra={<Button icon={<PlusOutlined />} onClick={() => setColModalOpen(true)}>新建集合</Button>}
    >
      {collections.length === 0 ? (
        <Empty description="暂无集合" />
      ) : (
        <Tree
          showIcon
          defaultExpandAll
          treeData={buildTree()}
          onSelect={handleSelect}
        />
      )}
      <Modal title="新建集合" open={colModalOpen} onOk={handleCreateCol} onCancel={() => setColModalOpen(false)}>
        <Form form={form} layout="vertical">
          <Form.Item name="name" label="集合名称" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="description" label="描述">
            <Input.TextArea rows={2} />
          </Form.Item>
        </Form>
      </Modal>
    </Card>
  )
}
