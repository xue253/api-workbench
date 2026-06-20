import { useState, useEffect } from 'react'
import { ConfigProvider, Card, List, Tag, Space, Typography, Spin, Alert, Button, theme } from 'antd'
import { ApiOutlined, CheckCircleOutlined, RocketOutlined } from '@ant-design/icons'
import { motion } from 'framer-motion'

const { Title, Text } = Typography

interface Item {
  id: number
  name: string
  value: number
}

function App() {
  const [items, setItems] = useState<Item[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    fetch('/api/items')
      .then(res => res.json())
      .then(data => {
        setItems(data)
        setLoading(false)
      })
      .catch(err => {
        setError(err.message)
        setLoading(false)
      })
  }, [])

  return (
    <ConfigProvider
      theme={{
        algorithm: theme.darkAlgorithm,
        token: {
          colorPrimary: '#1677ff',
          borderRadius: 8,
        },
      }}
    >
      <div style={{ 
        minHeight: '100vh', 
        background: '#141414',
        padding: '40px 20px'
      }}>
        <motion.div
          initial={{ opacity: 0, y: -20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5 }}
          style={{ maxWidth: 800, margin: '0 auto' }}
        >
          <Card 
            style={{ 
              background: '#1f1f1f',
              border: '1px solid #303030'
            }}
          >
            <Space direction="vertical" size="large" style={{ width: '100%' }}>
              <div style={{ textAlign: 'center' }}>
                <RocketOutlined style={{ fontSize: 48, color: '#1677ff' }} />
                <Title level={2} style={{ color: '#fff', margin: '16px 0 8px' }}>
                  API Workbench
                </Title>
                <Text type="secondary">React + Go 全栈项目实战</Text>
              </div>

              <div style={{ display: 'flex', gap: 16, justifyContent: 'center' }}>
                <Tag icon={<CheckCircleOutlined />} color="success">React 18</Tag>
                <Tag icon={<ApiOutlined />} color="processing">Go + Gin</Tag>
                <Tag color="warning">Ant Design 6</Tag>
              </div>

              {loading && (
                <div style={{ textAlign: 'center', padding: 40 }}>
                  <Spin size="large" tip="加载中..." />
                </div>
              )}

              {error && (
                <Alert
                  message="API 连接失败"
                  description="后端服务未启动，请先运行 Go 后端"
                  type="warning"
                  showIcon
                  action={
                    <Button size="small" href="https://go.dev/dl/" target="_blank">
                      安装 Go
                    </Button>
                  }
                />
              )}

              {!loading && !error && (
                <List
                  header={<Text strong style={{ color: '#fff' }}>API 数据</Text>}
                  bordered
                  dataSource={items}
                  renderItem={(item) => (
                    <List.Item>
                      <List.Item.Meta
                        title={<Text style={{ color: '#fff' }}>{item.name}</Text>}
                        description={<Text type="secondary">ID: {item.id}</Text>}
                      />
                      <Tag color="blue">{item.value}</Tag>
                    </List.Item>
                  )}
                />
              )}
            </Space>
          </Card>
        </motion.div>
      </div>
    </ConfigProvider>
  )
}

export default App
