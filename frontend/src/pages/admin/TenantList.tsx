import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Card, Table, Tag, Badge, Button, Spin, Typography, Space } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import { api } from '../../api/client';
import type { Tenant } from '../../api/client';

export default function TenantList() {
  const nav = useNavigate();
  const [tenants, setTenants] = useState<Tenant[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    api.get<Tenant[]>('/admin/tenants').then(setTenants).finally(() => setLoading(false));
  }, []);

  if (loading) return <Spin style={{ display: 'block', margin: '20vh auto' }} size="large" />;

  const planColor: Record<string, string> = { free: 'default', pro: 'blue', enterprise: 'gold' };

  return (
    <div style={{ padding: 24 }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 16 }}>
        <Typography.Title level={3} style={{ margin: 0 }}>租戶管理</Typography.Title>
        <Button type="primary" icon={<PlusOutlined />} onClick={() => nav('/admin/tenants/new')}>新增租戶</Button>
      </div>
      <Card>
        <Table dataSource={tenants} rowKey="id" columns={[
          { title: '租戶名稱', dataIndex: 'name', render: (v: string, r: Tenant) => <a onClick={() => nav(`/admin/tenants/${r.id}`)}>{v}</a> },
          { title: '識別碼', dataIndex: 'slug' },
          { title: '方案', dataIndex: 'plan', render: (v: string) => <Tag color={planColor[v] || 'default'}>{v}</Tag> },
          { title: '狀態', dataIndex: 'is_active', render: (v: boolean) => <Badge status={v ? 'success' : 'error'} text={v ? '啟用' : '停用'} /> },
          { title: '收件人上限', dataIndex: 'max_recipients' },
          { title: '操作', render: (_: unknown, r: Tenant) => (
            <Space>
              <Button size="small" onClick={() => nav(`/admin/tenants/${r.id}`)}>管理</Button>
            </Space>
          )},
        ]} />
      </Card>
    </div>
  );
}
