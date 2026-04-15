import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Card, Statistic, Table, Tag, Badge, Button, Row, Col, Spin } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import dayjs from 'dayjs';
import { api } from '../../api/client';
import type { PlatformStats, Tenant } from '../../api/client';

const planColor: Record<string, string> = { free: 'default', pro: 'blue', enterprise: 'gold' };

const columns = [
  { title: '公司名稱', dataIndex: 'name' },
  { title: 'Slug', dataIndex: 'slug' },
  { title: '方案', dataIndex: 'plan', render: (v: string) => <Tag color={planColor[v] ?? 'default'}>{v}</Tag> },
  { title: '狀態', dataIndex: 'is_active', render: (v: boolean) => <Badge status={v ? 'success' : 'error'} text={v ? '啟用' : '停用'} /> },
  { title: '建立時間', dataIndex: 'created_at', render: (v: string) => dayjs(v).format('YYYY-MM-DD') },
];

export default function AdminDashboard() {
  const navigate = useNavigate();
  const [stats, setStats] = useState<PlatformStats>();
  const [tenants, setTenants] = useState<Tenant[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    Promise.all([
      api.get<PlatformStats>('/admin/dashboard'),
      api.get<Tenant[]>('/admin/tenants'),
    ]).then(([s, t]) => { setStats(s); setTenants(t); }).finally(() => setLoading(false));
  }, []);

  if (loading) return <Spin style={{ display: 'block', margin: '20vh auto' }} size="large" />;

  return (
    <>
      <Row justify="space-between" align="middle" style={{ marginBottom: 16 }}>
        <Col><h2 style={{ margin: 0 }}>Platform Dashboard</h2></Col>
        <Col><Button type="primary" icon={<PlusOutlined />} onClick={() => navigate('/admin/tenants/new')}>新增租戶</Button></Col>
      </Row>

      <Row gutter={16} style={{ marginBottom: 24 }}>
        <Col span={12}>
          <Card><Statistic title="總租戶數" value={stats?.total_tenants} /></Card>
        </Col>
        <Col span={12}>
          <Card><Statistic title="啟用租戶" value={stats?.active_tenants} /></Card>
        </Col>
      </Row>

      <Card title="近期租戶">
        <Table
          rowKey="id"
          columns={columns}
          dataSource={tenants}
          pagination={{ pageSize: 10 }}
          onRow={(record) => ({ onClick: () => navigate(`/admin/tenants/${record.id}`), style: { cursor: 'pointer' } })}
        />
      </Card>
    </>
  );
}
