import { useState, useEffect } from 'react';
import { Table, Card, Tag, Typography } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import dayjs from 'dayjs';
import { api } from '../../api/client';

interface AuditLog {
  id: number;
  tenant_id: number | null;
  user_email: string;
  action: string;
  resource: string;
  resource_id: number | null;
  detail: string;
  ip_address: string;
  created_at: string;
}

const PAGE_SIZE = 30;

const columns: ColumnsType<AuditLog> = [
  { title: '時間', dataIndex: 'created_at', width: 180, render: (v: string) => dayjs(v).format('YYYY-MM-DD HH:mm:ss') },
  { title: '租戶', dataIndex: 'tenant_id', width: 80, render: (v: number | null) => v ?? '平台' },
  { title: '使用者', dataIndex: 'user_email', ellipsis: true },
  { title: '動作', dataIndex: 'action', width: 120, render: (v: string) => <Tag>{v}</Tag> },
  { title: '資源', dataIndex: 'resource', width: 120 },
  { title: 'IP', dataIndex: 'ip_address', width: 140 },
];

export default function AdminAuditLogs() {
  const [logs, setLogs] = useState<AuditLog[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    setLoading(true);
    const offset = (page - 1) * PAGE_SIZE;
    api.get<{ logs: AuditLog[]; total: number }>('/admin/audit-logs?limit=30&offset=' + offset)
      .then((d) => { setLogs(d.logs); setTotal(d.total); })
      .finally(() => setLoading(false));
  }, [page]);

  return (
    <Card>
      <Typography.Title level={4}>平台稽核日誌</Typography.Title>
      <Table<AuditLog>
        rowKey="id"
        columns={columns}
        dataSource={logs}
        loading={loading}
        pagination={{ current: page, pageSize: PAGE_SIZE, total, onChange: setPage, showSizeChanger: false }}
        expandable={{ expandedRowRender: (r) => <pre style={{ margin: 0, whiteSpace: 'pre-wrap' }}>{r.detail}</pre> }}
      />
    </Card>
  );
}
