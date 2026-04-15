import { useEffect, useState } from 'react';
import { Card, Table, Tag, Typography, Empty, Spin } from 'antd';
import { api } from '../../api/client';

interface OffenderHistory { campaign_id: number; campaign_name: string; status: string; }
interface Offender { email: string; first_name: string; last_name: string; department: string; history: OffenderHistory[]; click_count: number; submit_count: number; }

const riskTag = (o: Offender) =>
  o.submit_count >= 2 ? <Tag color="red">高風險</Tag> :
  o.click_count >= 2 ? <Tag color="orange">中風險</Tag> :
  <Tag>低風險</Tag>;

const historyColumns = [
  { title: '活動名稱', dataIndex: 'campaign_name', key: 'campaign_name' },
  { title: '狀態', dataIndex: 'status', key: 'status', render: (s: string) => <Tag>{s}</Tag> },
];

const columns = [
  { title: 'Email', dataIndex: 'email', key: 'email' },
  { title: '姓名', key: 'name', render: (_: unknown, r: Offender) => `${r.last_name}${r.first_name}` },
  { title: '部門', dataIndex: 'department', key: 'department', render: (d: string) => <Tag>{d}</Tag> },
  { title: '點擊次數', dataIndex: 'click_count', key: 'click_count', sorter: (a: Offender, b: Offender) => a.click_count - b.click_count, render: (v: number) => <span style={v >= 3 ? { color: '#ff4d4f', fontWeight: 'bold' } : undefined}>{v}</span> },
  { title: '提交次數', dataIndex: 'submit_count', key: 'submit_count', defaultSortOrder: 'descend' as const, sorter: (a: Offender, b: Offender) => a.submit_count - b.submit_count, render: (v: number) => <span style={v >= 2 ? { color: '#ff4d4f', fontWeight: 'bold' } : undefined}>{v}</span> },
  { title: '風險等級', key: 'risk', render: (_: unknown, r: Offender) => riskTag(r) },
];

export default function RepeatOffenders() {
  const [data, setData] = useState<Offender[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    api.get<Offender[]>('/reports/offenders').then(setData).finally(() => setLoading(false));
  }, []);

  if (loading) return <Spin style={{ display: 'block', margin: '20vh auto' }} size="large" />;

  return (
    <Card>
      <Typography.Title level={3}>累犯追蹤</Typography.Title>
      <Typography.Text type="secondary">曾在釣魚測試中點擊或提交資料的員工</Typography.Text>
      {data.length === 0 ? (
        <Empty description="目前沒有累犯記錄，表示員工資安意識良好！" style={{ marginTop: 48 }} />
      ) : (
        <Table
          rowKey="email"
          columns={columns}
          dataSource={data}
          style={{ marginTop: 24 }}
          expandable={{
            expandedRowRender: (r) => (
              <Table rowKey="campaign_id" columns={historyColumns} dataSource={r.history} pagination={false} size="small" />
            ),
          }}
        />
      )}
    </Card>
  );
}
