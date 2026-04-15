import { useEffect, useState } from 'react';
import { Card, Table, Typography, Empty, Spin } from 'antd';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';
import { api } from '../../api/client';

interface TrendPoint { campaign_id: number; campaign_name: string; launched_at: string; open_rate: number; click_rate: number; submit_rate: number; report_rate: number; }

const lines = [
  { key: 'open_rate', name: '開信率', color: '#13c2c2' },
  { key: 'click_rate', name: '點擊率', color: '#faad14' },
  { key: 'submit_rate', name: '提交率', color: '#ff4d4f' },
  { key: 'report_rate', name: '舉報率', color: '#52c41a' },
] as const;

const columns = [
  { title: '活動名稱', dataIndex: 'campaign_name', key: 'campaign_name' },
  { title: '發送時間', dataIndex: 'launched_at', key: 'launched_at' },
  ...lines.map(l => ({ title: l.name, dataIndex: l.key, key: l.key, render: (v: number) => `${v}%` })),
];

export default function TrendAnalysis() {
  const [data, setData] = useState<TrendPoint[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    api.get<TrendPoint[]>('/reports/trend').then(setData).finally(() => setLoading(false));
  }, []);

  if (loading) return <Spin style={{ display: 'block', margin: '20vh auto' }} size="large" />;

  if (data.length === 0) return <Card><Typography.Title level={3}>趨勢分析</Typography.Title><Empty description="目前沒有趨勢資料" /></Card>;

  return (
    <Card>
      <Typography.Title level={3}>趨勢分析</Typography.Title>
      <Typography.Text type="secondary">跨活動的指標變化趨勢</Typography.Text>
      <div style={{ marginTop: 24 }}>
        <ResponsiveContainer width="100%" height={400}>
          <LineChart data={data}>
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis dataKey="campaign_name" />
            <YAxis domain={[0, 100]} unit="%" />
            <Tooltip formatter={(v: unknown) => `${v}%`} />
            <Legend />
            {lines.map(l => <Line key={l.key} type="monotone" dataKey={l.key} name={l.name} stroke={l.color} />)}
          </LineChart>
        </ResponsiveContainer>
      </div>
      <Table rowKey="campaign_id" columns={columns} dataSource={data} pagination={false} style={{ marginTop: 24 }} />
    </Card>
  );
}
