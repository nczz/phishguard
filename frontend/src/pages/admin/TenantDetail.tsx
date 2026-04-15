import { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Card, Descriptions, Tag, Badge, Button, Spin, message } from 'antd';
import { api } from '../../api/client';
import type { Tenant } from '../../api/client';
import dayjs from 'dayjs';

export default function TenantDetail() {
  const { id } = useParams();
  const nav = useNavigate();
  const [tenant, setTenant] = useState<Tenant | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    api.get<Tenant>(`/admin/tenants/${id}`).then(t => setTenant(t)).catch(() => message.error('租戶不存在')).finally(() => setLoading(false));
  }, [id]);

  if (loading) return <Spin style={{ display: 'block', margin: '20vh auto' }} size="large" />;
  if (!tenant) return null;

  const planColor = { free: 'default', pro: 'blue', enterprise: 'gold' }[tenant.plan] || 'default';

  return (
    <div style={{ padding: 24 }}>
      <Button onClick={() => nav('/admin/tenants')} style={{ marginBottom: 16 }}>← 返回列表</Button>
      <Card title={tenant.name}>
        <Descriptions column={2} bordered>
          <Descriptions.Item label="識別碼">{tenant.slug}</Descriptions.Item>
          <Descriptions.Item label="方案"><Tag color={planColor}>{tenant.plan}</Tag></Descriptions.Item>
          <Descriptions.Item label="狀態"><Badge status={tenant.is_active ? 'success' : 'error'} text={tenant.is_active ? '啟用' : '停用'} /></Descriptions.Item>
          <Descriptions.Item label="收件人上限">{tenant.max_recipients}</Descriptions.Item>
          <Descriptions.Item label="建立時間">{dayjs(tenant.created_at).format('YYYY-MM-DD HH:mm')}</Descriptions.Item>
        </Descriptions>
      </Card>
    </div>
  );
}
