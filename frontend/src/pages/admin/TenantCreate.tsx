import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Card, Form, Input, Select, Button, Breadcrumb, message } from 'antd';
import { api } from '../../api/client';

export default function TenantCreate() {
  const navigate = useNavigate();
  const [form] = Form.useForm();
  const [submitting, setSubmitting] = useState(false);

  const onFinish = async (values: Record<string, string>) => {
    setSubmitting(true);
    try {
      await api.post('/admin/tenants', values);
      message.success('租戶建立成功');
      navigate('/admin/dashboard');
    } catch {
      message.error('建立失敗');
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <>
      <Breadcrumb style={{ marginBottom: 16 }} items={[
        { title: 'Admin' },
        { title: '租戶管理' },
        { title: '新增租戶' },
      ]} />

      <Card title="新增租戶" style={{ maxWidth: 600 }}>
        <Form form={form} layout="vertical" onFinish={onFinish} initialValues={{ plan: 'free' }}>
          <Form.Item label="公司名稱" name="name" rules={[{ required: true, message: '請輸入公司名稱' }]}>
            <Input onChange={(e) => form.setFieldValue('slug', e.target.value.toLowerCase().replace(/\s+/g, '-'))} />
          </Form.Item>
          <Form.Item label="識別碼" name="slug" rules={[{ required: true, message: '請輸入識別碼' }]}>
            <Input />
          </Form.Item>
          <Form.Item label="方案" name="plan">
            <Select options={[
              { value: 'free', label: 'Free' },
              { value: 'pro', label: 'Pro' },
              { value: 'enterprise', label: 'Enterprise' },
            ]} />
          </Form.Item>
          <Form.Item label="管理員 Email" name="admin_email" rules={[{ required: true, type: 'email', message: '請輸入有效 Email' }]}>
            <Input />
          </Form.Item>
          <Form.Item label="管理員密碼" name="admin_password" rules={[{ required: true, min: 8, message: '密碼至少 8 字元' }]}>
            <Input.Password />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={submitting}>建立租戶</Button>
          </Form.Item>
        </Form>
      </Card>
    </>
  );
}
