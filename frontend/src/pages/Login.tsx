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
      message.success('登录成功')
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
      message.success('注册成功')
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
      background: '#f5f5f7',
    }}>
      <div style={{
        width: 420,
        padding: 48,
        background: '#ffffff',
        borderRadius: 20,
        boxShadow: '0 4px 24px rgba(0,0,0,0.08)',
      }}>
        <div style={{ textAlign: 'center', marginBottom: 40 }}>
          <Title level={2} style={{ 
            color: '#1d1d1f', 
            margin: 0, 
            fontWeight: 600,
            letterSpacing: '-0.5px'
          }}>
            API Workbench
          </Title>
          <Text style={{ 
            color: '#86868b', 
            fontSize: 15,
            display: 'block',
            marginTop: 8
          }}>
            接口管理与自动化测试平台
          </Text>
        </div>
        <Tabs 
          centered 
          items={[
            {
              key: 'login',
              label: '登录',
              children: (
                <Form form={form} onFinish={handleLogin} layout="vertical" size="large">
                  <Form.Item name="username" rules={[{ required: true, message: '请输入用户名' }]}>
                    <Input 
                      prefix={<UserOutlined style={{ color: '#86868b' }} />} 
                      placeholder="用户名"
                      style={{ borderRadius: 8, height: 48 }}
                    />
                  </Form.Item>
                  <Form.Item name="password" rules={[{ required: true, message: '请输入密码' }]}>
                    <Input.Password 
                      prefix={<LockOutlined style={{ color: '#86868b' }} />} 
                      placeholder="密码"
                      style={{ borderRadius: 8, height: 48 }}
                    />
                  </Form.Item>
                  <Form.Item>
                    <Button 
                      type="primary" 
                      htmlType="submit" 
                      loading={loading} 
                      block 
                      style={{ 
                        borderRadius: 8, 
                        height: 48,
                        fontSize: 16,
                        fontWeight: 500,
                        background: '#0071e3',
                      }}
                    >
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
                    <Input 
                      prefix={<UserOutlined style={{ color: '#86868b' }} />} 
                      placeholder="用户名"
                      style={{ borderRadius: 8, height: 48 }}
                    />
                  </Form.Item>
                  <Form.Item name="password" rules={[{ required: true, min: 6, message: '密码至少 6 位' }]}>
                    <Input.Password 
                      prefix={<LockOutlined style={{ color: '#86868b' }} />} 
                      placeholder="密码"
                      style={{ borderRadius: 8, height: 48 }}
                    />
                  </Form.Item>
                  <Form.Item name="email">
                    <Input 
                      prefix={<MailOutlined style={{ color: '#86868b' }} />} 
                      placeholder="邮箱（选填）"
                      style={{ borderRadius: 8, height: 48 }}
                    />
                  </Form.Item>
                  <Form.Item>
                    <Button 
                      type="primary" 
                      htmlType="submit" 
                      loading={loading} 
                      block 
                      style={{ 
                        borderRadius: 8, 
                        height: 48,
                        fontSize: 16,
                        fontWeight: 500,
                        background: '#0071e3',
                      }}
                    >
                      注册
                    </Button>
                  </Form.Item>
                </Form>
              )
            }
          ]} 
        />
      </div>
    </div>
  )
}
