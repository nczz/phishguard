import { useEffect, useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { Card, Table, Tag, Button, Empty } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import dayjs from 'dayjs';
import { api } from '../../api/client';
import type { Campaign } from '../../api/client';

const statusMap: Record<string, { color: string; label: string }> = {
  draft: { color: 'default', label: '草稿' },
  scheduled: { color: 'orange', label: '排程中' },
  sending: { color: 'processing', label: '發送中' },
  paused: { color: 'warning', label: '已暫停' },
  sent: { color: 'blue', label: '已發送' },
  stopped: { color: 'red', label: '已終止' },
  completed: { color: 'green', label: '已完成' },
};

const columns = [
  {
    title: '活動名稱',
    dataIndex: 'name',
    render: (name: string, r: Campaign) => <Link to={`/app/campaigns/${r.id}`}>{name}</Link>,
  },
  {
    title: '狀態',
    dataIndex: 'status',
    width: 100,
    render: (s: string) => {
      const m = statusMap[s] ?? { color: 'default', label: s };
      return <Tag color={m.color}>{m.label}</Tag>;
    },
  },
  {
    title: '發送日期',
    dataIndex: 'launched_at',
    width: 160,
    render: (v: string) => (v ? dayjs(v).format('YYYY-MM-DD HH:mm') : '-'),
  },
  {
    title: '建立日期',
    dataIndex: 'created_at',
    width: 160,
    render: (v: string) => dayjs(v).format('YYYY-MM-DD HH:mm'),
  },
];

export default function CampaignList() {
  const navigate = useNavigate();
  const [campaigns, setCampaigns] = useState<Campaign[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    api.get<Campaign[]>('/campaigns').then(setCampaigns).finally(() => setLoading(false));
  }, []);

  return (
    <Card
      title="釣魚測試活動"
      extra={
        <Button type="primary" icon={<PlusOutlined />} onClick={() => navigate('/app/campaigns/new')}>
          建立新測試
        </Button>
      }
    >
      {!loading && campaigns.length === 0 ? (
        <Empty description="尚未建立任何測試">
          <Button type="primary" onClick={() => navigate('/app/campaigns/new')}>建立新測試</Button>
        </Empty>
      ) : (
        <Table rowKey="id" loading={loading} columns={columns} dataSource={campaigns} pagination={{ pageSize: 10 }} />
      )}
    </Card>
  );
}
