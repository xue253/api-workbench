import { useState, useEffect } from 'react'
import { Card, Form, Input, Select, Button, Tabs, Space, message } from 'antd'
import { SendOutlined } from '@ant-design/icons'
import { useApiStore } from '../stores/apiStore'
import api from '../services/api'

const { Option } = Select

export default function ApiEditor() {
  const { currentApi, updateAPI } = useApiStore()
  const [form] = Form.useForm()
  const [response, setResponse] = useState<any>(null)
  const [sending, setSending] = useState(false)

  useEffect(() => {
    if (currentApi) {
      form.setFieldsValue(currentApi)
    }
  }, [currentApi])

  const handleSend = async () => {
    if (!currentApi) return
    setSending(true)
    try {
      const res: any = await api.post(`/apis/${currentApi.id}/debug`)
      setResponse(res)
    } catch (err: any) {
      message.error(err.message)
    }
    setSending(false)
  }

  const handleSave = async () => {
    if (!currentApi) return
    const values = await form.validateFields()
    await updateAPI(currentApi.id, values)
    message.success('保存成功')
  }

  if (!currentApi) {
    return <Card><div style={{ textAlign: 'center', padding: 40, color: '#999' }}>选择接口开始编辑</div></Card>
  }

  return (
    <Card title={`编辑接口 - ${currentApi.name}`}>
      <Tabs items={[
        {
          key: 'config',
          label: '配置',
          children: (
            <Form form={form} layout="vertical">
              <Form.Item name="name" label="名称"><Input /></Form.Item>
              <Form.Item name="method" label="方法">
                <Select style={{ width: 120 }}>
                  {['GET', 'POST', 'PUT', 'DELETE', 'PATCH'].map(m => <Option key={m} value={m}>{m}</Option>)}
                </Select>
              </Form.Item>
              <Form.Item name="url" label="URL"><Input /></Form.Item>
              <Form.Item name="headers" label="请求头"><Input.TextArea rows={3} placeholder='JSON 格式' /></Form.Item>
              <Form.Item name="body" label="请求体"><Input.TextArea rows={5} /></Form.Item>
              <Form.Item name="body_type" label="Body 类型">
                <Select style={{ width: 120 }}>
                  {['json', 'form', 'raw', 'binary'].map(t => <Option key={t} value={t}>{t}</Option>)}
                </Select>
              </Form.Item>
              <Form.Item name="timeout_ms" label="超时(ms)"><Input type="number" /></Form.Item>
              <Space>
                <Button type="primary" onClick={handleSave}>保存</Button>
                <Button icon={<SendOutlined />} loading={sending} onClick={handleSend}>发送</Button>
              </Space>
            </Form>
          )
        },
        {
          key: 'response',
          label: '响应',
          children: response ? (
            <div>
              <div style={{ marginBottom: 8 }}>
                <strong>状态码:</strong> {response.status_code}
              </div>
              <pre style={{ background: '#1f1f1f', color: '#d4d4d4', padding: 16, borderRadius: 8, maxHeight: 400, overflow: 'auto' }}>
                {typeof response.body === 'string' ? response.body : JSON.stringify(response.body, null, 2)}
              </pre>
            </div>
          ) : (
            <div style={{ textAlign: 'center', padding: 40, color: '#999' }}>点击「发送」查看响应</div>
          )
        }
      ]} />
    </Card>
  )
}
