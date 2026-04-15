import { useEffect, useState } from 'react';
import { Card, Table, Tag, Button, Drawer, Form, Input, Select, Popconfirm, Space, message } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import { api } from '../../api/client';
import FieldHelp, { tips } from '../../components/FieldHelp';
import type { EmailTemplate } from '../../api/client';

const categories = ['credential', 'malware', 'social', 'compliance', 'custom'];
const languages = ['zh-TW', 'zh-CN', 'en', 'ja'];

export default function TemplateList() {
  const [data, setData] = useState<EmailTemplate[]>([]);
  const [loading, setLoading] = useState(true);
  const [open, setOpen] = useState(false);
  const [editing, setEditing] = useState<EmailTemplate | null>(null);
  const [form] = Form.useForm();

  const load = () => {
    setLoading(true);
    api.get<EmailTemplate[]>('/templates').then(setData).finally(() => setLoading(false));
  };

  useEffect(load, []);

  const openDrawer = (record?: EmailTemplate) => {
    setEditing(record ?? null);
    form.resetFields();
    if (record) form.setFieldsValue(record);
    setOpen(true);
  };

  const onFinish = async (values: Record<string, unknown>) => {
    if (editing) {
      await api.put('/templates/' + editing.id, values);
      message.success('已更新');
    } else {
      await api.post('/templates', values);
      message.success('已新增');
    }
    setOpen(false);
    load();
  };

  const onDelete = async (id: string) => {
    await api.del('/templates/' + id);
    message.success('已刪除');
    load();
  };

  const columns = [
    { title: '名稱', dataIndex: 'name' },
    { title: '主旨', dataIndex: 'subject' },
    { title: '分類', dataIndex: 'category', width: 100, render: (v: string) => <Tag>{v}</Tag> },
    { title: '語言', dataIndex: 'language', width: 80, render: (v: string) => <Tag>{v}</Tag> },
    {
      title: '來源', dataIndex: 'tenant_id', width: 80,
      render: (v: string | null) => v ? <Tag color="green">自建</Tag> : <Tag>平台</Tag>,
    },
    {
      title: '操作', width: 140,
      render: (_: unknown, r: EmailTemplate) => (
        <Space>
          {r.tenant_id && <Button size="small" onClick={() => openDrawer(r)}>編輯</Button>}
          <Popconfirm title="確定刪除？" onConfirm={() => onDelete(r.id)}>
            <Button size="small" danger>刪除</Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <Card
      title="模板管理"
      extra={<Button type="primary" icon={<PlusOutlined />} onClick={() => openDrawer()}>新增模板</Button>}
    >
      <Table rowKey="id" loading={loading} columns={columns} dataSource={data} pagination={{ pageSize: 10 }} />

      <Drawer
        title={editing ? '編輯模板' : '新增模板'}
        width={600}
        open={open}
        onClose={() => setOpen(false)}
        extra={<Button type="primary" onClick={() => form.submit()}>儲存</Button>}
      >
        <Form form={form} layout="vertical" onFinish={onFinish}>
          <Form.Item name="name" label="名稱" rules={[{ required: true }]}><Input /></Form.Item>
          <Form.Item name="subject" label={<FieldHelp label="主旨" tip="信件主旨，收件人在信箱中看到的標題" />} rules={[{ required: true }]}><Input /></Form.Item>
          <Form.Item name="category" label="分類" rules={[{ required: true }]}>
            <Select options={categories.map((c) => ({ label: c, value: c }))} />
          </Form.Item>
          <Form.Item name="language" label="語言" rules={[{ required: true }]}>
            <Select options={languages.map((l) => ({ label: l, value: l }))} />
          </Form.Item>
          <Form.Item name="html_body" label={<FieldHelp label="HTML 內容" tip={tips.templateVars} guideAnchor="variables" />}><Input.TextArea rows={8} /></Form.Item>
          <Form.Item name="text_body" label={<FieldHelp label="純文字內容" tip="當收件人的郵件客戶端不支援 HTML 時顯示的內容" />}><Input.TextArea rows={4} /></Form.Item>
        </Form>
      </Drawer>
    </Card>
  );
}
