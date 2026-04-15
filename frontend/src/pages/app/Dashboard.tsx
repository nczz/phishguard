import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Row, Col, Card, Statistic, List, Tag, Button, Progress, Empty, Typography, Spin } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import { api } from '../../api/client';
import type { Campaign } from '../../api/client';

const statusColor: Record<string, string> = {
  draft: 'default', scheduled: 'orange', sending: 'processing', sent: 'blue', completed: 'green',
};

const rateColor = (v: number) => (v <= 20 ? '#52c41a' : v <= 50 ? '#faad14' : '#f5222d');

export default function Dashboard() {
  const navigate = useNavigate();
  const [campaigns, setCampaigns] = useState<Campaign[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    api.get<Campaign[]>('/campaigns').then(setCampaigns).finally(() => setLoading(false));
  }, []);

  const total = campaigns.length;
  // placeholder rates — real stats require per-campaign report fetching
  const avgOpen = 0;
  const avgClick = 0;
  const avgSubmit = 0;

  if (loading) return <Spin style={{ display: 'block', margin: '20vh auto' }} size="large" />;

  return (
    <div>
      <Typography.Title level={3}>歡迎回來</Typography.Title>

      <Row gutter={[16, 16]}>
        <Col xs={12} sm={6}><Card><Statistic title="總測試數" value={total} /></Card></Col>
        <Col xs={12} sm={6}><Card><Statistic title="平均開信率" value={avgOpen} suffix="%" valueStyle={{ color: rateColor(avgOpen) }} /></Card></Col>
        <Col xs={12} sm={6}><Card><Statistic title="平均點擊率" value={avgClick} suffix="%" valueStyle={{ color: rateColor(avgClick) }} /></Card></Col>
        <Col xs={12} sm={6}><Card><Statistic title="平均提交率" value={avgSubmit} suffix="%" valueStyle={{ color: rateColor(avgSubmit) }} /></Card></Col>
      </Row>

      <Card title="最近活動" style={{ marginTop: 24 }}>
        {campaigns.length === 0 ? (
          <Empty description="尚未建立任何測試" />
        ) : (
          <List
            dataSource={campaigns.slice(0, 5)}
            renderItem={(c) => (
              <List.Item
                style={{ cursor: 'pointer' }}
                onClick={() => navigate(`/app/campaigns/${c.id}`)}
              >
                <List.Item.Meta title={c.name} />
                <Tag color={statusColor[c.status]}>{c.status}</Tag>
                <Progress percent={0} size="small" style={{ width: 120 }} />
              </List.Item>
            )}
          />
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
