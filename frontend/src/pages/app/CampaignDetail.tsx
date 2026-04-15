import { useEffect, useState, useCallback } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {
  Card, Tag, Statistic, List, Progress, Button, Row, Col, Spin, Breadcrumb, Tooltip, Typography,
} from 'antd';
import { ArrowLeftOutlined, FilePdfOutlined } from '@ant-design/icons';
import {
  BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip as RechartsTooltip,
  Cell, ResponsiveContainer,
} from 'recharts';
import dayjs from 'dayjs';
import { api, type Campaign, type CampaignReport, type DepartmentStat } from '../../api/client';

const STATUS_COLOR: Record<string, string> = {
  draft: 'default', scheduled: 'processing', sending: 'blue', completed: 'success', cancelled: 'error',
};

const FUNNEL_COLORS = ['#1677ff', '#13c2c2', '#faad14', '#ff4d4f', '#52c41a'];

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
  const [loading, setLoading] = useState(true);

  const fetchData = useCallback(async () => {
    if (!id) return;
    try {
      const [c, r] = await Promise.all([
        api.get<Campaign>('/campaigns/' + id),
        api.get<CampaignReport>('/campaigns/' + id + '/report'),
      ]);
      setCampaign(c);
      setReport(r);
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
            <List
              dataSource={deptsSorted}
              renderItem={(d: DepartmentStat) => {
                const rate = d.total ? Math.round((d.clicked / d.total) * 100) : 0;
                return (
                  <List.Item>
                    <List.Item.Meta title={d.department} />
                    <div style={{ width: 180 }}>
                      <Progress percent={rate} strokeColor={rateColor(rate)} size="small" />
                    </div>
                  </List.Item>
                );
              }}
            />
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

      {/* Section 3: Actions */}
      <div style={{ display: 'flex', gap: 12 }}>
        <Tooltip title="Phase 2 開發中">
          <Button icon={<FilePdfOutlined />} disabled>匯出 PDF</Button>
        </Tooltip>
        <Button icon={<ArrowLeftOutlined />} onClick={() => navigate('/app/campaigns')}>返回列表</Button>
      </div>
    </>
  );
}
