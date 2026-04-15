import { useEffect, useState } from 'react';
import { Navigate } from 'react-router-dom';
import { Card, Form, Input, Button, Typography, Space, Checkbox, message } from 'antd';
import { MailOutlined, LockOutlined, SafetyOutlined } from '@ant-design/icons';
import { useAuth } from '../hooks/useAuth';

const { Title, Text } = Typography;

export default function Login() {
  const { user, login, loading } = useAuth();
  const [submitting, setSubmitting] = useState(false);

  // prevent message context warning
  useEffect(() => () => message.destroy(), []);

  if (!loading && user) {
    return <Navigate to={user.role === 'platform_admin' ? '/admin/dashboard' : '/app/dashboard'} replace />;
  }

  const onFinish = async (values: { email: string; password: string }) => {
    setSubmitting(true);
    try {
      await login(values.email, values.password);
    } catch {
      message.error('登入失敗，請確認帳號密碼');
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div style={{ minHeight: '100vh', display: 'flex', alignItems: 'center', justifyContent: 'center', background: '#f0f2f5' }}>
      <Card style={{ width: 400 }}>
        <Space orientation="vertical" align="center" style={{ width: '100%', marginBottom: 24 }}>
          <SafetyOutlined style={{ fontSize: 48, color: '#1677ff' }} />
          <Title level={2} style={{ margin: 0 }}>PhishGuard</Title>
          <Text type="secondary">企業釣魚模擬測試平台</Text>
        </Space>

        <Form layout="vertical" onFinish={onFinish}>
          <Form.Item name="email" rules={[{ required: true, message: '請輸入 Email' }]}>
            <Input prefix={<MailOutlined />} placeholder="Email" size="large" />
          </Form.Item>
          <Form.Item name="password" rules={[{ required: true, message: '請輸入密碼' }]}>
            <Input.Password prefix={<LockOutlined />} placeholder="密碼" size="large" />
          </Form.Item>
          <Form.Item name="remember" valuePropName="checked">
            <Checkbox>Remember me</Checkbox>
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" block size="large" loading={submitting}>
              登入
            </Button>
          </Form.Item>
        </Form>
      </Card>
    </div>
  );
}
