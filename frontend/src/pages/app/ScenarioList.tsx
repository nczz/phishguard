import { useEffect, useState } from 'react';
import { Card, Row, Col, Tag, Badge, Empty, Spin, Typography, Button, Drawer, Form, Input, Select, message, Popconfirm, Space } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons';
import { api } from '../../api/client';
import FieldHelp, { tips } from '../../components/FieldHelp';
import type { Scenario, EmailTemplate, LandingPage } from '../../api/client';

const categoryEmoji: Record<string, string> = {
  password_reset: '🔐', invoice: '💳', package: '📦', hr_notice: '📋',
  it_alert: '🖥️', social_engineering: '🎭', credential_harvest: '🔑',
};

const categories = [
  { value: 'password_reset', label: '🔐 密碼重設' },
  { value: 'invoice', label: '💳 發票/報價單' },
  { value: 'package', label: '📦 包裹通知' },
  { value: 'hr_notice', label: '📋 HR 公告' },
  { value: 'it_alert', label: '🖥️ IT 警告' },
  { value: 'social_engineering', label: '🎭 社交工程' },
  { value: 'credential_harvest', label: '🔑 憑證竊取' },
];

export default function ScenarioList() {
  const [scenarios, setScenarios] = useState<Scenario[]>([]);
  const [templates, setTemplates] = useState<EmailTemplate[]>([]);
  const [pages, setPages] = useState<LandingPage[]>([]);
  const [loading, setLoading] = useState(true);
  const [open, setOpen] = useState(false);
  const [editing, setEditing] = useState<Scenario | null>(null);
  const [previewHtml, setPreviewHtml] = useState('');
  const [previewOpen, setPreviewOpen] = useState(false);
  const [form] = Form.useForm();

  const load = () => {
    setLoading(true);
    Promise.all([
      api.get<Scenario[]>('/scenarios'),
      api.get<EmailTemplate[]>('/templates'),
      api.get<LandingPage[]>('/pages'),
    ]).then(([s, t, p]) => { setScenarios(s); setTemplates(t); setPages(p); }).finally(() => setLoading(false));
  };

  useEffect(load, []);

  const onSubmit = async (values: Record<string, unknown>) => {
    try {
      if (editing) {
        await api.put(`/scenarios/${editing.id}`, values);
        message.success('情境已更新');
      } else {
        await api.post('/scenarios', values);
        message.success('情境建立成功');
      }
      setOpen(false);
      setEditing(null);
      form.resetFields();
      load();
    } catch { message.error('操作失敗'); }
  };

  const openEdit = (s: Scenario) => {
    setEditing(s);
    form.setFieldsValue({ name: s.name, category: s.category, difficulty: s.difficulty, language: s.language, template_id: s.template_id, page_id: s.page_id, education_html: s.education_html });
    setOpen(true);
  };

  const openCreate = () => { setEditing(null); form.resetFields(); setOpen(true); };

  const deleteScenario = async (id: number | string) => {
    try { await api.del(`/scenarios/${id}`); message.success('已刪除'); load(); }
    catch { message.error('刪除失敗'); }
  };

  if (loading) return <Spin style={{ display: 'block', margin: '20vh auto' }} size="large" />;

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16 }}>
        <Typography.Title level={3} style={{ margin: 0 }}>情境庫</Typography.Title>
        <Button type="primary" icon={<PlusOutlined />} onClick={openCreate}>新增情境</Button>
      </div>

      {scenarios.length === 0 ? (
        <Empty description="尚無情境">
          <Button type="primary" onClick={openCreate}>建立第一個情境</Button>
        </Empty>
      ) : (
        <Row gutter={[16, 16]}>
          {scenarios.map((s) => (
            <Col key={s.id} xs={24} sm={12} lg={8}>
              <Badge.Ribbon text={s.is_active ? '啟用' : '停用'} color={s.is_active ? 'green' : 'red'}>
                <Card hoverable>
                  <div style={{ fontSize: 36, textAlign: 'center', marginBottom: 8 }}>{categoryEmoji[s.category] ?? '📧'}</div>
                  <Typography.Title level={5} style={{ textAlign: 'center', margin: 0 }}>{s.name}</Typography.Title>
                  <div style={{ textAlign: 'center', margin: '8px 0' }}>{'⭐'.repeat(Number(s.difficulty) || 1)}</div>
                  <div style={{ textAlign: 'center' }}>
                    <Tag>{s.language}</Tag>
                    <Tag color="blue">{s.category}</Tag>
                    {s.tenant_id ? <Tag color="orange">自建</Tag> : <Tag color="purple">平台預建</Tag>}
                  </div>
                  {s.tenant_id && (
                    <div style={{ textAlign: 'center', marginTop: 12, borderTop: '1px solid #f0f0f0', paddingTop: 8 }}>
                      <Space>
                        <Button size="small" icon={<EditOutlined />} onClick={() => openEdit(s)}>編輯</Button>
                        <Popconfirm title="確定刪除此情境？" onConfirm={() => deleteScenario(s.id)}>
                          <Button size="small" danger icon={<DeleteOutlined />}>刪除</Button>
                        </Popconfirm>
                      </Space>
                    </div>
                  )}
                </Card>
              </Badge.Ribbon>
            </Col>
          ))}
        </Row>
      )}

      <Drawer title={editing ? '編輯情境' : '新增情境'} open={open} onClose={() => { setOpen(false); setEditing(null); }} width={520} extra={<Button type="primary" onClick={() => form.submit()}>儲存</Button>}>
        <Form form={form} layout="vertical" onFinish={onSubmit} initialValues={{ difficulty: 2, language: 'zh-TW', is_active: true }}>
          <Form.Item name="name" label={<FieldHelp label="情境名稱" tip={tips.scenario} guideAnchor="scenarios" />} rules={[{ required: true }]}><Input placeholder="例：密碼到期通知" /></Form.Item>
          <Form.Item name="category" label="分類" rules={[{ required: true }]}><Select options={categories} placeholder="選擇分類" /></Form.Item>
          <Form.Item name="difficulty" label={<FieldHelp label="難度" tip={tips.difficulty} />} rules={[{ required: true }]}>
            <Select options={[{ value: 1, label: '⭐ 簡單' }, { value: 2, label: '⭐⭐ 中等' }, { value: 3, label: '⭐⭐⭐ 困難' }]} />
          </Form.Item>
          <Form.Item name="language" label="語言"><Select options={[{ value: 'zh-TW', label: '繁體中文' }, { value: 'en', label: 'English' }]} /></Form.Item>
          <Form.Item name="template_id" label={<FieldHelp label="信件模板" tip={tips.templateVars} guideAnchor="variables" />} rules={[{ required: true, message: '請先建立模板再選擇' }]}>
            <Select placeholder={templates.length ? '選擇模板' : '請先到模板管理建立模板'} options={templates.map(t => ({ value: t.id, label: `${t.name} — ${t.subject}` }))} />
          </Form.Item>
          <Form.Item name="page_id" label={<FieldHelp label="Landing Page" tip={tips.submitURL} guideAnchor="variables" />} rules={[{ required: true, message: '請先建立 Landing Page' }]}>
            <Select placeholder={pages.length ? '選擇 Landing Page' : '請先到模板管理建立頁面'} options={pages.map(p => ({ value: p.id, label: p.name }))} />
          </Form.Item>
          <Form.Item name="education_html" label={<FieldHelp label="教育頁內容 (HTML)" tip={tips.educationHTML} guideAnchor="metrics" />} rules={[{ required: true }]}>
            <Input.TextArea rows={6} placeholder="<h1>這是一封釣魚測試信</h1><p>以下是辨識方法...</p>" />
          </Form.Item>
          <Button type="dashed" block style={{ marginBottom: 16 }} onClick={() => { setPreviewHtml(form.getFieldValue('education_html') || ''); setPreviewOpen(true); }}>👁 預覽教育頁</Button>
        </Form>
      </Drawer>

      <Drawer title="教育頁預覽" open={previewOpen} onClose={() => setPreviewOpen(false)} width={520}>
        <div dangerouslySetInnerHTML={{ __html: previewHtml }} style={{ border: '1px solid #d9d9d9', borderRadius: 8, padding: 16, minHeight: 400 }} />
      </Drawer>
    </div>
  );
}
