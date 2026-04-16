import { useEffect, useState } from 'react';
import { Card, Table, Tag, Button, Modal, Form, Input, InputNumber, Radio, Checkbox, Space, Popconfirm, message, Progress, Alert, Row, Col, Typography, Select } from 'antd';
import { PlusOutlined, SafetyCertificateOutlined, CheckCircleOutlined, WarningOutlined, CloseCircleOutlined } from '@ant-design/icons';
import FieldHelp, { tips } from '../../components/FieldHelp';
import { api } from '../../api/client';
import type { SMTPProfile } from '../../api/client';

const typeColor: Record<string, string> = { smtp: 'blue', mailgun: 'purple', ses: 'orange' };

export default function SMTPSettings() {
  const [data, setData] = useState<SMTPProfile[]>([]);
  const [loading, setLoading] = useState(true);
  const [open, setOpen] = useState(false);
  const [editing, setEditing] = useState<SMTPProfile | null>(null);
  const [testModal, setTestModal] = useState<string | null>(null);
  const [testEmail, setTestEmail] = useState('');
  const [form] = Form.useForm();
  const mailerType = Form.useWatch('mailer_type', form);

  // Compliance check state
  interface CompCheck { name: string; status: string; detail: string; fix?: string; }
  interface CompResult { domain: string; score: number; checks: CompCheck[]; }
  const [compResult, setCompResult] = useState<CompResult | null>(null);
  const [compLoading, setCompLoading] = useState(false);
  const [compProfileId, setCompProfileId] = useState<number | null>(null);

  const runCompliance = async () => {
    if (!compProfileId) { message.error('請選擇 SMTP 設定'); return; }
    setCompLoading(true);
    try {
      const res = await api.post<CompResult>('/smtp-profiles/check-compliance', { smtp_profile_id: compProfileId });
      setCompResult(res);
    } catch { message.error('檢測失敗'); }
    setCompLoading(false);
  };

  const statusIcon = (s: string) => s === 'pass' ? <CheckCircleOutlined style={{ color: '#52c41a' }} /> : s === 'warn' ? <WarningOutlined style={{ color: '#faad14' }} /> : <CloseCircleOutlined style={{ color: '#ff4d4f' }} />;
  const statusColor = (s: string) => s === 'pass' ? 'success' : s === 'warn' ? 'warning' : 'error';

  const load = () => {
    setLoading(true);
    api.get<SMTPProfile[]>('/smtp-profiles').then(setData).finally(() => setLoading(false));
  };

  useEffect(load, []);

  const onFinish = async (values: Record<string, unknown>) => {
    try {
      if (editing) {
        await api.put('/smtp-profiles/' + editing.id, values);
        message.success('已更新');
      } else {
        await api.post('/smtp-profiles', values);
        message.success('已新增');
      }
      setOpen(false); setEditing(null); form.resetFields(); load();
    } catch (e: unknown) {
      const msg = (e && typeof e === 'object' && 'displayMessage' in e) ? (e as { displayMessage: string }).displayMessage : '操作失敗';
      message.error(msg);
    }
  };

  const openEdit = (r: SMTPProfile) => { setEditing(r); form.setFieldsValue(r); setOpen(true); };
  const openCreate = () => { setEditing(null); form.resetFields(); setOpen(true); };

  const onDelete = async (id: string) => {
    await api.del('/smtp-profiles/' + id);
    message.success('已刪除');
    load();
  };

  const onTest = async () => {
    if (!testModal || !testEmail) return;
    try {
      await api.post('/smtp-profiles/' + testModal + '/test', { to: testEmail });
      message.success('測試信已發送');
      setTestModal(null);
      setTestEmail('');
    } catch (e: unknown) {
      const msg = (e && typeof e === 'object' && 'displayMessage' in e) ? (e as { displayMessage: string }).displayMessage : '發送失敗，請檢查 SMTP 設定是否正確';
      message.error(msg);
    }
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
          <Button size="small" onClick={() => openEdit(r)}>編輯</Button>
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
      extra={<Button type="primary" icon={<PlusOutlined />} onClick={openCreate}>新增設定</Button>}
    >
      <Table rowKey="id" loading={loading} columns={columns} dataSource={data} pagination={{ pageSize: 10 }} />

      <Modal title={editing ? '編輯 SMTP 設定' : '新增 SMTP 設定'} open={open} onOk={() => form.submit()} onCancel={() => setOpen(false)} width={520} destroyOnHidden>
        <Form form={form} layout="vertical" onFinish={onFinish} initialValues={{ mailer_type: 'smtp', port: 587, tls: true }}>
          <Form.Item name="name" label="名稱" rules={[{ required: true }]}><Input /></Form.Item>
          <Form.Item name="mailer_type" label={<FieldHelp label="發信方式" tip={tips.smtpType} guideAnchor="smtp" />}>
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

          <Form.Item name="from_address" label={<FieldHelp label="寄件地址" tip="收件人看到的寄件者 email 地址。Mailgun 用戶：域名必須與 Mailgun Domain 一致（如 noreply@yourdomain.com）" />} rules={[{ required: true, type: 'email' }]}><Input placeholder="noreply@yourdomain.com" /></Form.Item>
          <Form.Item name="from_name" label="寄件人名稱"><Input /></Form.Item>
        </Form>
      </Modal>

      <Modal title="發送測試信" open={!!testModal} onOk={onTest} onCancel={() => setTestModal(null)}>
        <Input placeholder="收件人 Email" value={testEmail} onChange={(e) => setTestEmail(e.target.value)} />
      </Modal>

      {/* Compliance Check Button */}
      <Card title={<><SafetyCertificateOutlined /> 發信合規檢測</>} style={{ marginTop: 24 }}>
        <Typography.Paragraph type="secondary">檢測寄件域名的 SPF、DKIM、DMARC 設定，確保信件不會進入垃圾郵件匣。</Typography.Paragraph>
        <Space>
          <Select placeholder="選擇 SMTP 設定" value={compProfileId} onChange={setCompProfileId} style={{ width: 320 }}>
            {data.map(p => <Select.Option key={p.id} value={p.id}>{p.name}（{p.from_address}）</Select.Option>)}
          </Select>
          <Button type="primary" icon={<SafetyCertificateOutlined />} loading={compLoading} onClick={runCompliance}>開始檢測</Button>
        </Space>

        {compResult && (
          <div style={{ marginTop: 24 }}>
            <Row gutter={16} align="middle" style={{ marginBottom: 16 }}>
              <Col>
                <Progress type="circle" percent={compResult.score} size={80}
                  strokeColor={compResult.score >= 80 ? '#52c41a' : compResult.score >= 50 ? '#faad14' : '#ff4d4f'} />
              </Col>
              <Col>
                <Typography.Title level={4} style={{ margin: 0 }}>
                  {compResult.domain} — {compResult.score >= 80 ? '良好' : compResult.score >= 50 ? '需改善' : '不合格'}
                </Typography.Title>
                <Typography.Text type="secondary">合規分數 {compResult.score}/100</Typography.Text>
              </Col>
            </Row>

            {compResult.checks.map((ck, i) => (
              <Alert key={i} type={statusColor(ck.status) as 'success' | 'warning' | 'error'} showIcon icon={statusIcon(ck.status)}
                title={ck.name} description={<>{ck.detail}{ck.fix && <div style={{ marginTop: 8, padding: 8, background: '#fafafa', borderRadius: 4, fontFamily: 'monospace', fontSize: 12 }}>💡 修正建議：{ck.fix}</div>}</>}
                style={{ marginBottom: 8 }} />
            ))}
          </div>
        )}
      </Card>
    </Card>
  );
}
