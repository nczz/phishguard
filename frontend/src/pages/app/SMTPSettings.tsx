import { useEffect, useState } from 'react';
import { Card, Table, Tag, Button, Modal, Form, Input, InputNumber, Radio, Checkbox, Space, Popconfirm, message } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import { api } from '../../api/client';
import type { SMTPProfile } from '../../api/client';

const typeColor: Record<string, string> = { smtp: 'blue', mailgun: 'purple', ses: 'orange' };

export default function SMTPSettings() {
  const [data, setData] = useState<SMTPProfile[]>([]);
  const [loading, setLoading] = useState(true);
  const [open, setOpen] = useState(false);
  const [testModal, setTestModal] = useState<string | null>(null);
  const [testEmail, setTestEmail] = useState('');
  const [form] = Form.useForm();
  const mailerType = Form.useWatch('mailer_type', form);

  const load = () => {
    setLoading(true);
    api.get<SMTPProfile[]>('/smtp-profiles').then(setData).finally(() => setLoading(false));
  };

  useEffect(load, []);

  const onFinish = async (values: Record<string, unknown>) => {
    await api.post('/smtp-profiles', values);
    message.success('已新增');
    setOpen(false);
    load();
  };

  const onDelete = async (id: string) => {
    await api.del('/smtp-profiles/' + id);
    message.success('已刪除');
    load();
  };

  const onTest = async () => {
    if (!testModal || !testEmail) return;
    await api.post('/smtp-profiles/' + testModal + '/test', { to: testEmail });
    message.success('測試信已發送');
    setTestModal(null);
    setTestEmail('');
  };

  const columns = [
    { title: '名稱', dataIndex: 'name' },
    {
      title: '類型', dataIndex: 'mailer_type', width: 100,
      render: (v: string) => <Tag color={typeColor[v]}>{v.toUpperCase()}</Tag>,
    },
    { title: '寄件地址', dataIndex: 'from_address' },
    { title: '寄件人', dataIndex: 'from_name' },
    {
      title: '操作', width: 160,
      render: (_: unknown, r: SMTPProfile) => (
        <Space>
          <Button size="small" onClick={() => { setTestModal(r.id); setTestEmail(''); }}>測試</Button>
          <Popconfirm title="確定刪除？" onConfirm={() => onDelete(r.id)}>
            <Button size="small" danger>刪除</Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <Card
      title="SMTP 設定"
      extra={<Button type="primary" icon={<PlusOutlined />} onClick={() => { form.resetFields(); form.setFieldValue('mailer_type', 'smtp'); setOpen(true); }}>新增設定</Button>}
    >
      <Table rowKey="id" loading={loading} columns={columns} dataSource={data} pagination={{ pageSize: 10 }} />

      <Modal title="新增 SMTP 設定" open={open} onOk={() => form.submit()} onCancel={() => setOpen(false)} width={520} destroyOnClose>
        <Form form={form} layout="vertical" onFinish={onFinish} initialValues={{ mailer_type: 'smtp', port: 587, tls: true }}>
          <Form.Item name="name" label="名稱" rules={[{ required: true }]}><Input /></Form.Item>
          <Form.Item name="mailer_type" label="類型">
            <Radio.Group>
              <Radio.Button value="smtp">SMTP</Radio.Button>
              <Radio.Button value="mailgun">Mailgun</Radio.Button>
              <Radio.Button value="ses">SES</Radio.Button>
            </Radio.Group>
          </Form.Item>

          {mailerType === 'smtp' && (
            <>
              <Form.Item name="host" label="Host" rules={[{ required: true }]}><Input /></Form.Item>
              <Form.Item name="port" label="Port" rules={[{ required: true }]}><InputNumber style={{ width: '100%' }} /></Form.Item>
              <Form.Item name="username" label="Username"><Input /></Form.Item>
              <Form.Item name="password" label="Password"><Input.Password /></Form.Item>
              <Form.Item name="tls" valuePropName="checked"><Checkbox>啟用 TLS</Checkbox></Form.Item>
            </>
          )}
          {mailerType === 'mailgun' && (
            <>
              <Form.Item name="mailgun_domain" label="Mailgun Domain" rules={[{ required: true }]}><Input /></Form.Item>
              <Form.Item name="mailgun_api_key" label="Mailgun API Key" rules={[{ required: true }]}><Input.Password /></Form.Item>
            </>
          )}
          {mailerType === 'ses' && (
            <>
              <Form.Item name="ses_region" label="SES Region" rules={[{ required: true }]}><Input /></Form.Item>
              <Form.Item name="ses_access_key" label="Access Key" rules={[{ required: true }]}><Input /></Form.Item>
              <Form.Item name="ses_secret_key" label="Secret Key" rules={[{ required: true }]}><Input.Password /></Form.Item>
            </>
          )}

          <Form.Item name="from_address" label="寄件地址" rules={[{ required: true }]}><Input /></Form.Item>
          <Form.Item name="from_name" label="寄件人名稱"><Input /></Form.Item>
        </Form>
      </Modal>

      <Modal title="發送測試信" open={!!testModal} onOk={onTest} onCancel={() => setTestModal(null)}>
        <Input placeholder="收件人 Email" value={testEmail} onChange={(e) => setTestEmail(e.target.value)} />
      </Modal>
    </Card>
  );
}
