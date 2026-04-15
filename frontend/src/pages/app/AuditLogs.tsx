import { useEffect, useState } from 'react';
import { Card, Table, Tag } from 'antd';
import dayjs from 'dayjs';
import { api } from '../../api/client';
import type { AuditLog } from '../../api/client';

const PAGE_SIZE = 20;

export default function AuditLogs() {
  const [data, setData] = useState<AuditLog[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);

  const load = (p: number) => {
    setLoading(true);
    api.get<{ logs: AuditLog[]; total: number }>(`/audit-logs?limit=${PAGE_SIZE}&offset=${(p - 1) * PAGE_SIZE}`)
      .then((res) => { setData(res.logs); setTotal(res.total); })
      .finally(() => setLoading(false));
  };

  useEffect(() => { load(page); }, [page]);

  const columns = [
    { title: '時間', dataIndex: 'created_at', width: 170, render: (v: string) => dayjs(v).format('YYYY-MM-DD HH:mm:ss') },
    { title: '使用者', dataIndex: 'user_email' },
    { title: '動作', dataIndex: 'action', width: 120, render: (v: string) => <Tag>{v}</Tag> },
    { title: '資源', dataIndex: 'resource', width: 120 },
    { title: 'IP', dataIndex: 'ip_address', width: 140 },
  ];

  return (
    <Card title="稽核日誌">
      <Table
        rowKey="id"
        loading={loading}
        columns={columns}
        dataSource={data}
        pagination={{ current: page, pageSize: PAGE_SIZE, total, onChange: setPage }}
        expandable={{ expandedRowRender: (r) => <pre style={{ margin: 0, whiteSpace: 'pre-wrap' }}>{r.detail}</pre> }}
      />
    </Card>
  );
}
