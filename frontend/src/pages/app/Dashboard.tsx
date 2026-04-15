import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Row, Col, Card, Statistic, Tag, Button, Progress, Empty, Typography, Spin, message, Alert } from 'antd';
import { PlusOutlined, ExperimentOutlined } from '@ant-design/icons';
import { api } from '../../api/client';
import type { Campaign, Scenario } from '../../api/client';

interface DashboardStats { total_campaigns: number; avg_open_rate: number; avg_click_rate: number; avg_submit_rate: number; avg_report_rate: number; }

const statusColor: Record<string, string> = {
  draft: 'default', scheduled: 'orange', sending: 'processing', sent: 'blue', completed: 'green',
};

const rateColor = (v: number) => (v <= 20 ? '#52c41a' : v <= 50 ? '#faad14' : '#f5222d');

const emptyStats: DashboardStats = { total_campaigns: 0, avg_open_rate: 0, avg_click_rate: 0, avg_submit_rate: 0, avg_report_rate: 0 };

export default function Dashboard() {
  const navigate = useNavigate();
  const [campaigns, setCampaigns] = useState<Campaign[]>([]);
  const [stats, setStats] = useState<DashboardStats>(emptyStats);
  const [hasScenarios, setHasScenarios] = useState(true);
  const [seeding, setSeeding] = useState(false);
  const [loading, setLoading] = useState(true);

  const load = () => {
    setLoading(true);
    Promise.all([
      api.get<Campaign[]>('/campaigns'),
      api.get<Scenario[]>('/scenarios'),
      api.get<DashboardStats>('/reports/dashboard-stats'),
    ]).then(([c, s, d]) => { setCampaigns(c); setHasScenarios(s.length > 0); setStats(d); }).finally(() => setLoading(false));
  };

  useEffect(load, []);

  const seedData = async () => {
    setSeeding(true);
    try {
      await api.post('/seed-sample-data');
      message.success('範例資料已匯入！包含 5 個情境、5 個模板、2 個 Landing Page、5 位範例收件人');
      load();
    } catch { message.error('匯入失敗'); }
    setSeeding(false);
  };

  if (loading) return <Spin style={{ display: 'block', margin: '20vh auto' }} size="large" />;

  return (
    <div>
      <Typography.Title level={3}>歡迎回來</Typography.Title>

      {!hasScenarios && (
        <Alert type="info" showIcon style={{ marginBottom: 16 }}
          title="尚未設定測試資料"
          description="匯入範例資料可快速體驗完整功能，包含 5 個釣魚情境、信件模板、Landing Page 和範例收件人。"
          action={<Button type="primary" icon={<ExperimentOutlined />} loading={seeding} onClick={seedData}>匯入範例資料</Button>}
        />
      )}

      <Row gutter={[16, 16]}>
        <Col xs={12} sm={6}><Card><Statistic title="總測試數" value={stats.total_campaigns} /></Card></Col>
        <Col xs={12} sm={6}><Card><Statistic title="平均開信率" value={stats.avg_open_rate} suffix="%" styles={{ content: { color: rateColor(stats.avg_open_rate) } }} /></Card></Col>
        <Col xs={12} sm={6}><Card><Statistic title="平均點擊率" value={stats.avg_click_rate} suffix="%" styles={{ content: { color: rateColor(stats.avg_click_rate) } }} /></Card></Col>
        <Col xs={12} sm={6}><Card><Statistic title="平均提交率" value={stats.avg_submit_rate} suffix="%" styles={{ content: { color: rateColor(stats.avg_submit_rate) } }} /></Card></Col>
      </Row>

      <Card title="最近活動" style={{ marginTop: 24 }}>
        {campaigns.length === 0 ? (
          <Empty description="尚未建立任何測試" />
        ) : (
          <div>
            {campaigns.slice(0, 5).map((c) => (
              <div key={c.id} onClick={() => navigate(`/app/campaigns/${c.id}`)} style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', padding: '12px 0', borderBottom: '1px solid #f0f0f0', cursor: 'pointer' }}>
                <span style={{ fontWeight: 500 }}>{c.name}</span>
                <span><Tag color={statusColor[c.status]}>{c.status}</Tag><Progress percent={0} size="small" style={{ width: 120, marginLeft: 8 }} /></span>
              </div>
            ))}
          </div>
        )}
      </Card>

      <div style={{ textAlign: 'center', marginTop: 32 }}>
        <Button type="primary" size="large" icon={<PlusOutlined />} onClick={() => navigate('/app/campaigns/new')}>
          建立新測試
        </Button>
      </div>
    </div>
  );
}
