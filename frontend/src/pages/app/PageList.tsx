import { useEffect, useState } from 'react';
import { Card, Table, Button, Drawer, Form, Input, Switch, message, Tag, Popconfirm, Space, Typography } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined, EyeOutlined } from '@ant-design/icons';
import { api } from '../../api/client';
import type { LandingPage } from '../../api/client';

export default function PageList() {
  const [pages, setPages] = useState<LandingPage[]>([]);
  const [loading, setLoading] = useState(true);
  const [open, setOpen] = useState(false);
  const [editing, setEditing] = useState<LandingPage | null>(null);
  const [previewHtml, setPreviewHtml] = useState('');
  const [previewOpen, setPreviewOpen] = useState(false);
  const [form] = Form.useForm();

  const load = () => { setLoading(true); api.get<LandingPage[]>('/pages').then(setPages).finally(() => setLoading(false)); };
  useEffect(load, []);

  const onSubmit = async (values: Record<string, unknown>) => {
    try {
      if (editing) { await api.put(`/pages/${editing.id}`, values); message.success('已更新'); }
      else { await api.post('/pages', values); message.success('已建立'); }
      setOpen(false); setEditing(null); form.resetFields(); load();
    } catch { message.error('操作失敗'); }
  };

  const openEdit = (p: LandingPage) => { setEditing(p); form.setFieldsValue(p); setOpen(true); };
  const openCreate = () => { setEditing(null); form.resetFields(); setOpen(true); };
  const preview = (html: string) => { setPreviewHtml(html); setPreviewOpen(true); };
  const del = async (id: string) => { try { await api.del(`/pages/${id}`); message.success('已刪除'); load(); } catch { message.error('刪除失敗'); } };

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16 }}>
        <Typography.Title level={3} style={{ margin: 0 }}>Landing Page 管理</Typography.Title>
        <Button type="primary" icon={<PlusOutlined />} onClick={openCreate}>新增頁面</Button>
      </div>

      <Card>
        <Table dataSource={pages} rowKey="id" loading={loading} columns={[
          { title: '名稱', dataIndex: 'name' },
          { title: '擷取憑證', dataIndex: 'capture_credentials', width: 100, render: (v: boolean) => v ? <Tag color="red">是</Tag> : <Tag>否</Tag> },
          { title: '來源', width: 100, render: (_: unknown, r: LandingPage) => r.tenant_id ? <Tag color="orange">自建</Tag> : <Tag color="purple">平台預建</Tag> },
          { title: '操作', width: 160, render: (_: unknown, r: LandingPage) => (
            <Space size="small">
              <Button size="small" icon={<EyeOutlined />} onClick={() => preview(r.html)}>預覽</Button>
              {r.tenant_id && <>
                <Button size="small" icon={<EditOutlined />} onClick={() => openEdit(r)} />
                <Popconfirm title="確定刪除？" onConfirm={() => del(r.id)}><Button size="small" danger icon={<DeleteOutlined />} /></Popconfirm>
              </>}
            </Space>
          )},
        ]} />
      </Card>

      <Drawer title={editing ? '編輯 Landing Page' : '新增 Landing Page'} open={open} onClose={() => { setOpen(false); setEditing(null); }} width={640}
        extra={<Button type="primary" onClick={() => form.submit()}>儲存</Button>}>
        <Form form={form} layout="vertical" onFinish={onSubmit} initialValues={{ capture_credentials: true }}>
          <Form.Item name="name" label="頁面名稱" rules={[{ required: true }]}><Input placeholder="例：仿公司登入頁" /></Form.Item>
          <Form.Item name="capture_credentials" label="擷取表單資料" valuePropName="checked"><Switch checkedChildren="是" unCheckedChildren="否" /></Form.Item>
          <Form.Item name="capture_fields" label="擷取欄位（JSON 陣列）"><Input placeholder='["email","password"]' /></Form.Item>
          <Form.Item name="redirect_url" label="提交後導向網址（選填）"><Input placeholder="https://www.google.com" /></Form.Item>
          <Form.Item name="html" label="頁面 HTML" rules={[{ required: true }]}>
            <Input.TextArea rows={16} placeholder={'<!DOCTYPE html>\n<html>\n<body>\n  <form action="{{.SubmitURL}}" method="POST">\n    <input name="email" placeholder="Email" />\n    <input name="password" type="password" placeholder="密碼" />\n    <button type="submit">登入</button>\n  </form>\n</body>\n</html>'} />
          </Form.Item>
          <Typography.Text type="secondary">
            提示：表單的 action 請使用 <Tag color="blue" style={{ fontFamily: 'monospace' }}>{'{{.SubmitURL}}'}</Tag>，系統會自動替換為追蹤網址。
          </Typography.Text>
        </Form>
      </Drawer>

      <Drawer title="頁面預覽" open={previewOpen} onClose={() => setPreviewOpen(false)} width={500}>
        <div dangerouslySetInnerHTML={{ __html: previewHtml }} style={{ border: '1px solid #d9d9d9', borderRadius: 8, minHeight: 400 }} />
      </Drawer>
    </div>
  );
}
