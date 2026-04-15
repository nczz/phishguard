import { useEffect, useState, useCallback } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {
  Card, Tag, Statistic, Progress, Button, Row, Col, Spin, Breadcrumb, Table, Typography, message,
} from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { ArrowLeftOutlined, DownloadOutlined, FilePdfOutlined, MailOutlined } from '@ant-design/icons';
import {
  BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip as RechartsTooltip,
  Cell, ResponsiveContainer,
} from 'recharts';
import dayjs from 'dayjs';
import { api, type Campaign, type CampaignReport, type DepartmentStat, type RecipientResult } from '../../api/client';

const STATUS_COLOR: Record<string, string> = {
  draft: 'default', scheduled: 'processing', sending: 'blue', completed: 'success', cancelled: 'error',
};

const FUNNEL_COLORS = ['#1677ff', '#13c2c2', '#faad14', '#ff4d4f', '#52c41a'];

const RECIPIENT_STATUS_COLOR: Record<string, string> = {
  sent: 'blue', opened: 'cyan', clicked: 'orange', submitted: 'red', reported: 'green',
};

function rateColor(rate: number) {
  if (rate > 40) return '#ff4d4f';
  if (rate > 20) return '#faad14';
  return '#52c41a';
}

export default function CampaignDetail() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [campaign, setCampaign] = useState<Campaign>();
  const [report, setReport] = useState<CampaignReport>();
  const [recipients, setRecipients] = useState<RecipientResult[]>([]);
  const [loading, setLoading] = useState(true);

  const fetchData = useCallback(async () => {
    if (!id) return;
    try {
      const [c, r, rec] = await Promise.all([
        api.get<Campaign>('/campaigns/' + id),
        api.get<CampaignReport>('/campaigns/' + id + '/report'),
        api.get<RecipientResult[]>('/campaigns/' + id + '/recipients'),
      ]);
      setCampaign(c);
      setReport(r);
      setRecipients(rec);
    } finally {
      setLoading(false);
    }
  }, [id]);

  useEffect(() => { fetchData(); }, [fetchData]);

  // Auto-refresh every 30s while sending
  useEffect(() => {
    if (campaign?.status !== 'sending') return;
    const timer = setInterval(fetchData, 30_000);
    return () => clearInterval(timer);
  }, [campaign?.status, fetchData]);

  if (loading) return <Spin style={{ display: 'block', margin: '20vh auto' }} size="large" />;
  if (!campaign || !report) return null;

  const { funnel } = report;
  const funnelData = [
    { name: '寄達', value: funnel.sent },
    { name: '開信', value: funnel.opened },
    { name: '點擊', value: funnel.clicked },
    { name: '提交', value: funnel.submitted },
    { name: '舉報', value: funnel.reported },
  ];

  const pct = (n: number) => funnel.total ? ((n / funnel.total) * 100).toFixed(1) + '%' : '0%';

  const deptsSorted = [...report.departments].sort(
    (a: DepartmentStat, b: DepartmentStat) => (b.total ? b.clicked / b.total : 0) - (a.total ? a.clicked / a.total : 0),
  );

  const exportCSV = async () => {
    const token = localStorage.getItem('token');
    const res = await fetch('/api/campaigns/' + id + '/export/csv', {
      headers: { Authorization: 'Bearer ' + token },
    });
    const blob = await res.blob();
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'campaign_results.csv';
    a.click();
    URL.revokeObjectURL(url);
  };

  const exportPDF = async () => {
    const res = await fetch('/api/campaigns/' + id + '/report/pdf', { headers: { Authorization: 'Bearer ' + localStorage.getItem('token') } });
    const blob = await res.blob();
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a'); a.href = url; a.download = 'campaign_report.pdf'; a.click();
    URL.revokeObjectURL(url);
  };

  const sendReport = async () => {
    try {
      await api.post('/campaigns/' + id + '/send-report');
      message.success('報表已寄送給租戶管理員');
    } catch { message.error('寄送失敗'); }
  };

  const fmt = (v: string | null) => (v ? dayjs(v).format('MM-DD HH:mm') : '-');

  const recipientColumns: ColumnsType<RecipientResult> = [
    {
      title: 'Email', dataIndex: 'email', key: 'email',
      filterDropdown: ({ setSelectedKeys, selectedKeys, confirm, clearFilters }) => (
        <div style={{ padding: 8 }}>
          <input
            placeholder="搜尋 Email"
            value={selectedKeys[0] as string}
            onChange={(e) => setSelectedKeys(e.target.value ? [e.target.value] : [])}
            onKeyDown={(e) => e.key === 'Enter' && confirm()}
            style={{ width: 200, marginBottom: 8, display: 'block' }}
          />
          <Button size="small" type="primary" onClick={() => confirm()} style={{ width: 90, marginRight: 8 }}>搜尋</Button>
          <Button size="small" onClick={() => { clearFilters?.(); confirm(); }} style={{ width: 90 }}>重設</Button>
        </div>
      ),
      onFilter: (value, record) => record.email.toLowerCase().includes(String(value).toLowerCase()),
    },
    { title: '姓名', key: 'name', render: (_, r) => r.last_name + r.first_name },
    { title: '部門', dataIndex: 'department', key: 'department' },
    {
      title: '狀態', dataIndex: 'status', key: 'status',
      render: (s: string) => <Tag color={RECIPIENT_STATUS_COLOR[s] ?? 'default'}>{s}</Tag>,
    },
    { title: '寄達', key: 'sent_at', render: (_, r) => fmt(r.sent_at) },
    { title: '開信', key: 'opened_at', render: (_, r) => fmt(r.opened_at) },
    { title: '點擊', key: 'clicked_at', render: (_, r) => fmt(r.clicked_at) },
    { title: '提交', key: 'submitted_at', render: (_, r) => fmt(r.submitted_at) },
    { title: '舉報', key: 'reported_at', render: (_, r) => fmt(r.reported_at) },
  ];

  return (
    <>
      <Breadcrumb
        style={{ marginBottom: 16 }}
        items={[
          { title: <a onClick={() => navigate('/app/campaigns')}>釣魚測試</a> },
          { title: campaign.name },
        ]}
      />

      <Typography.Title level={3} style={{ marginBottom: 8 }}>
        {campaign.name}{' '}
        <Tag color={STATUS_COLOR[campaign.status] ?? 'default'}>{campaign.status}</Tag>
      </Typography.Title>
      {campaign.launched_at && (
        <Typography.Text type="secondary" style={{ display: 'block', marginBottom: 24 }}>
          發送時間：{dayjs(campaign.launched_at).format('YYYY-MM-DD HH:mm')}
        </Typography.Text>
      )}

      {/* Section 1: Funnel */}
      <Card title="釣魚漏斗" style={{ marginBottom: 24 }}>
        <ResponsiveContainer width="100%" height={300}>
          <BarChart data={funnelData} layout="vertical" margin={{ left: 20, right: 60 }}>
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis type="number" />
            <YAxis type="category" dataKey="name" width={60} />
            <RechartsTooltip formatter={(v) => [String(v), '人數']} />
            <Bar dataKey="value" radius={[0, 6, 6, 0]} label={{ position: 'right', formatter: (v) => `${v} (${pct(Number(v))})` }}>
              {funnelData.map((_, i) => (
                <Cell key={i} fill={FUNNEL_COLORS[i]} />
              ))}
            </Bar>
          </BarChart>
        </ResponsiveContainer>
      </Card>

      {/* Section 2: Department risk + Summary */}
      <Row gutter={24} style={{ marginBottom: 24 }}>
        <Col xs={24} lg={12}>
          <Card title="部門風險排名" style={{ height: '100%' }}>
            {deptsSorted.map((d: DepartmentStat) => {
              const rate = d.total ? Math.round((d.clicked / d.total) * 100) : 0;
              return (
                <div key={d.department} style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', padding: '8px 0', borderBottom: '1px solid #f0f0f0' }}>
                  <span>{d.department}</span>
                  <div style={{ width: 180 }}><Progress percent={rate} strokeColor={rateColor(rate)} size="small" /></div>
                </div>
              );
            })}
          </Card>
        </Col>
        <Col xs={24} lg={12}>
          <Card title="統計摘要" style={{ height: '100%' }}>
            <Row gutter={[16, 16]}>
              <Col span={12}><Statistic title="總收件人" value={funnel.total} /></Col>
              <Col span={12}><Statistic title="開信率" value={pct(funnel.opened)} /></Col>
              <Col span={12}><Statistic title="點擊率" value={pct(funnel.clicked)} /></Col>
              <Col span={12}><Statistic title="提交率" value={pct(funnel.submitted)} /></Col>
              <Col span={12}><Statistic title="舉報率" value={pct(funnel.reported)} /></Col>
            </Row>
          </Card>
        </Col>
      </Row>

      {/* Section 3: Recipient Detail */}
      <Card title="收件人明細" style={{ marginBottom: 24 }}>
        <Table<RecipientResult>
          columns={recipientColumns}
          dataSource={recipients}
          rowKey="email"
          size="small"
          pagination={{ pageSize: 20, showSizeChanger: true }}
        />
      </Card>

      {/* Section 4: Actions */}
      <div style={{ display: 'flex', gap: 12 }}>
        <Button icon={<DownloadOutlined />} onClick={exportCSV}>匯出 CSV</Button>
        <Button icon={<FilePdfOutlined />} onClick={exportPDF}>匯出 PDF</Button>
        {campaign.status === 'completed' && <Button icon={<MailOutlined />} onClick={sendReport}>寄送報表</Button>}
        <Button icon={<ArrowLeftOutlined />} onClick={() => navigate('/app/campaigns')}>返回列表</Button>
      </div>
    </>
  );
}
