import { useState, useEffect } from 'react'
import { Card, Form, Input, Button, Switch, Select, message, Divider, Space, Tag } from 'antd'
import { UserOutlined, LockOutlined, SaveOutlined, InfoCircleOutlined } from '@ant-design/icons'
import { useAuthStore } from '../stores/authStore'
import api from '../services/api'

export default function SettingsPage() {
  const { user } = useAuthStore()
  const [profileForm] = Form.useForm()
  const [passwordForm] = Form.useForm()
  const [saving, setSaving] = useState(false)
  const [changingPwd, setChangingPwd] = useState(false)

  useEffect(() => {
    if (user) {
      profileForm.setFieldsValue({
        username: user.username,
        email: user.email,
      })
    }
  }, [user])

  const handleSaveProfile = async () => {
    try {
      const values = await profileForm.validateFields()
      setSaving(true)
      await api.put('/user/profile', values)
      message.success('保存成功')
    } catch (err: any) {
      if (err.errorFields) return
      message.error(err.message || '保存失败')
    } finally {
      setSaving(false)
    }
  }

  const handleChangePassword = async () => {
    try {
      const values = await passwordForm.validateFields()
      setChangingPwd(true)
      await api.put('/user/password', values)
      message.success('密码修改成功')
      passwordForm.resetFields()
    } catch (err: any) {
      if (err.errorFields) return
      message.error(err.message || '修改失败')
    } finally {
      setChangingPwd(false)
    }
  }

  return (
    <div style={{ maxWidth: 600, margin: '0 auto' }}>
      <Card style={{ borderRadius: 16, border: 'none', boxShadow: '0 2px 12px rgba(0,0,0,0.08)', marginBottom: 24 }} styles={{ body: { padding: 0 } }}>
        <div style={{ padding: '20px 24px 0', borderBottom: '1px solid #f0f0f0' }}>
          <h3 style={{ margin: 0, fontSize: 16, fontWeight: 600 }}>个人信息</h3>
        </div>
        <div style={{ padding: 24 }}>
          <Form form={profileForm} layout="vertical">
            <Form.Item name="username" label="用户名">
              <Input prefix={<UserOutlined />} disabled />
            </Form.Item>
            <Form.Item name="email" label="邮箱">
              <Input placeholder="输入邮箱" />
            </Form.Item>
            <Form.Item>
              <Button type="primary" icon={<SaveOutlined />} loading={saving} onClick={handleSaveProfile} style={{ borderRadius: 8 }}>
                保存
              </Button>
            </Form.Item>
          </Form>
        </div>
      </Card>

      <Card style={{ borderRadius: 16, border: 'none', boxShadow: '0 2px 12px rgba(0,0,0,0.08)', marginBottom: 24 }} styles={{ body: { padding: 0 } }}>
        <div style={{ padding: '20px 24px 0', borderBottom: '1px solid #f0f0f0' }}>
          <h3 style={{ margin: 0, fontSize: 16, fontWeight: 600 }}>修改密码</h3>
        </div>
        <div style={{ padding: 24 }}>
          <Form form={passwordForm} layout="vertical">
            <Form.Item name="old_password" label="当前密码" rules={[{ required: true, message: '请输入当前密码' }]}>
              <Input.Password prefix={<LockOutlined />} placeholder="当前密码" />
            </Form.Item>
            <Form.Item name="new_password" label="新密码" rules={[{ required: true, min: 6, message: '密码至少 6 位' }]}>
              <Input.Password prefix={<LockOutlined />} placeholder="新密码" />
            </Form.Item>
            <Form.Item>
              <Button type="primary" icon={<SaveOutlined />} loading={changingPwd} onClick={handleChangePassword} style={{ borderRadius: 8 }}>
                修改密码
              </Button>
            </Form.Item>
          </Form>
        </div>
      </Card>

      <Card style={{ borderRadius: 16, border: 'none', boxShadow: '0 2px 12px rgba(0,0,0,0.08)' }}>
        <div style={{ padding: '20px 24px 0', borderBottom: '1px solid #f0f0f0' }}>
          <h3 style={{ margin: 0, fontSize: 16, fontWeight: 600 }}>关于</h3>
        </div>
        <div style={{ padding: 24 }}>
          <div style={{ display: 'flex', flexDirection: 'column', gap: 12, fontSize: 14 }}>
            <div style={{ display: 'flex', justifyContent: 'space-between' }}>
              <span style={{ color: '#86868b' }}>产品名称</span>
              <span>API Workbench</span>
            </div>
            <div style={{ display: 'flex', justifyContent: 'space-between' }}>
              <span style={{ color: '#86868b' }}>版本</span>
              <Tag>v0.1.0</Tag>
            </div>
            <div style={{ display: 'flex', justifyContent: 'space-between' }}>
              <span style={{ color: '#86868b' }}>技术栈</span>
              <span>Go + React + Ant Design</span>
            </div>
            <div style={{ display: 'flex', justifyContent: 'space-between' }}>
              <span style={{ color: '#86868b' }}>描述</span>
              <span>接口管理与自动化测试平台</span>
            </div>
          </div>
        </div>
      </Card>
    </div>
  )
}
