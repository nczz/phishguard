import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Card, Statistic, Table, Tag, Badge, Button, Row, Col, Alert, Spin, Typography,
} from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import { api } from '../../api/client';

interface TenantStat {
  tenant_id: number;
  tenant_name: string;
  slug: string;
  plan: string;
  is_active: boolean;
  recipient_count: number;
  campaign_count: number;
  emails_sent: number;
}

interface AlertItem {
  type: string;
  message: string;
}

interface DashboardData {
  total_tenants: number;
  active_tenants: number;
  total_recipients: number;
  total_campaigns: number;
  total_emails: number;
  tenants: TenantStat[];
  alerts: AlertItem[];
}

const planColor: Record<string, string> = { free: 'default', pro: 'blue', enterprise: 'gold' };

export default function AdminDashboard() {
  const navigate = useNavigate();
  const [data, setData] = useState<DashboardData>();
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    api.get<DashboardData>('/admin/dashboard')
      .then(setData)
      .finally(() => setLoading(false));
  }, []);

  if (loading) return <Spin style={{ display: 'block', margin: '20vh auto' }} size="large" />;

  const columns: ColumnsType<TenantStat> = [
    {
      title: '租戶名稱',
      dataIndex: 'tenant_name',
      render: (v, r) => <a onClick={() => navigate(`/admin/tenants/${r.tenant_id}`)}>{v}</a>,
    },
    {
      title: '方案',
      dataIndex: 'plan',
      render: (v: string) => <Tag color={planColor[v] ?? 'default'}>{v}</Tag>,
    },
    {
      title: '狀態',
      dataIndex: 'is_active',
      render: (v: boolean) => <Badge status={v ? 'success' : 'error'} text={v ? '啟用' : '停用'} />,
    },
    { title: '收件人數', dataIndex: 'recipient_count' },
    { title: '活動數', dataIndex: 'campaign_count' },
    { title: '發信量', dataIndex: 'emails_sent' },
    {
      title: '操作',
      render: (_, r) => <Button size="small" onClick={() => navigate(`/admin/tenants/${r.tenant_id}`)}>進入</Button>,
    },
  ];

  const stats = [
    { title: '總租戶數', value: data?.total_tenants },
    { title: '活躍租戶', value: data?.active_tenants },
    { title: '總收件人', value: data?.total_recipients },
    { title: '總活動數', value: data?.total_campaigns },
    { title: '總發信量', value: data?.total_emails },
  ];

  return (
    <>
      <Row justify="space-between" align="middle" style={{ marginBottom: 16 }}>
        <Col><Typography.Title level={3} style={{ margin: 0 }}>Platform Dashboard</Typography.Title></Col>
        <Col>
          <Button type="primary" icon={<PlusOutlined />} onClick={() => navigate('/admin/tenants/new')}>
            + 新增租戶
          </Button>
        </Col>
      </Row>

      <Row gutter={16} style={{ marginBottom: 24 }}>
        {stats.map((s) => (
          <Col key={s.title} xs={24} sm={12} md={8} lg={4} xl={4} style={{ marginBottom: 8 }}>
            <Card><Statistic title={s.title} value={s.value} /></Card>
          </Col>
        ))}
      </Row>

      {data?.alerts && data.alerts.length > 0 && (
        <div style={{ marginBottom: 24 }}>
          {data.alerts.map((a, i) => (
            <Alert
              key={i}
              type={a.type as 'warning' | 'error' | 'info'}
              title={a.message}
              showIcon
              style={{ marginBottom: 8 }}
            />
          ))}
        </div>
      )}

      <Card title="租戶列表">
        <Table<TenantStat>
          rowKey="tenant_id"
          columns={columns}
          dataSource={data?.tenants}
          pagination={{ pageSize: 10 }}
        />
      </Card>
    </>
  );
}
