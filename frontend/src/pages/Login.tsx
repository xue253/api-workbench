import { useState } from 'react'
import { Card, Form, Input, Button, Tabs, message, Typography } from 'antd'
import { UserOutlined, LockOutlined, MailOutlined } from '@ant-design/icons'
import { useAuthStore } from '../stores/authStore'

const { Title, Text } = Typography

export default function LoginPage({ onLogin }: { onLogin: () => void }) {
  const { login, register } = useAuthStore()
  const [loading, setLoading] = useState(false)
  const [form] = Form.useForm()

  const handleLogin = async () => {
    const values = await form.validateFields()
    setLoading(true)
    try {
      await login(values.username, values.password)
      message.success('登录成功 ✿')
      onLogin()
    } catch (err: any) {
      message.error(err.message || '登录失败')
    }
    setLoading(false)
  }

  const handleRegister = async () => {
    const values = await form.validateFields()
    setLoading(true)
    try {
      await register(values.username, values.password, values.email)
      message.success('注册成功 ✿')
      onLogin()
    } catch (err: any) {
      message.error(err.message || '注册失败')
    }
    setLoading(false)
  }

  return (
    <div style={{
      minHeight: '100vh',
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'center',
      background: 'linear-gradient(135deg, #fff0f6 0%, #ffd6e7 50%, #fff5f7 100%)',
    }}>
      <Card style={{
        width: 420,
        borderRadius: 16,
        border: '2px solid #ffadd2',
        boxShadow: '0 8px 32px rgba(247,89,171,0.15)',
      }}>
        <div style={{ textAlign: 'center', marginBottom: 24 }}>
          <div style={{ fontSize: 48, marginBottom: 8 }}>🌸</div>
          <Title level={3} style={{ color: '#f759ab', margin: 0 }}>API Workbench</Title>
          <Text type="secondary">接口管理与自动化测试平台</Text>
        </div>
        <Tabs centered items={[
          {
            key: 'login',
            label: '登录',
            children: (
              <Form form={form} onFinish={handleLogin} layout="vertical" size="large">
                <Form.Item name="username" rules={[{ required: true, message: '请输入用户名' }]}>
                  <Input prefix={<UserOutlined />} placeholder="用户名" />
                </Form.Item>
                <Form.Item name="password" rules={[{ required: true, message: '请输入密码' }]}>
                  <Input.Password prefix={<LockOutlined />} placeholder="密码" />
                </Form.Item>
                <Form.Item>
                  <Button type="primary" htmlType="submit" loading={loading} block style={{ borderRadius: 8 }}>
                    登录
                  </Button>
                </Form.Item>
              </Form>
            )
          },
          {
            key: 'register',
            label: '注册',
            children: (
              <Form form={form} onFinish={handleRegister} layout="vertical" size="large">
                <Form.Item name="username" rules={[{ required: true, min: 3, max: 50, message: '用户名 3-50 个字符' }]}>
                  <Input prefix={<UserOutlined />} placeholder="用户名" />
                </Form.Item>
                <Form.Item name="password" rules={[{ required: true, min: 6, message: '密码至少 6 位' }]}>
                  <Input.Password prefix={<LockOutlined />} placeholder="密码" />
                </Form.Item>
                <Form.Item name="email">
                  <Input prefix={<MailOutlined />} placeholder="邮箱（选填）" />
                </Form.Item>
                <Form.Item>
                  <Button type="primary" htmlType="submit" loading={loading} block style={{ borderRadius: 8 }}>
                    注册
                  </Button>
                </Form.Item>
              </Form>
            )
          }
        ]} />
      </Card>
    </div>
  )
}
