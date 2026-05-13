import { useEffect, useState } from 'react';
import { Card, Table, Tag, Tooltip } from 'antd';
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

  const formatDetail = (v: string) => {
    if (!v) return '-';
    try {
      const obj = JSON.parse(v);
      const keys = Object.keys(obj).filter(k => obj[k] !== '***').slice(0, 3);
      return keys.map(k => `${k}: ${typeof obj[k] === 'string' ? obj[k] : JSON.stringify(obj[k])}`).join(', ');
    } catch { return v; }
  };

  const columns = [
    { title: '時間', dataIndex: 'created_at', width: 170, render: (v: string) => dayjs(v).format('YYYY-MM-DD HH:mm:ss') },
    { title: '使用者', dataIndex: 'user_email', width: 180, ellipsis: { showTitle: false }, render: (v: string) => <Tooltip title={v}>{v}</Tooltip> },
    { title: '動作', dataIndex: 'action', width: 150, render: (v: string) => <Tag>{v}</Tag> },
    { title: '資源', dataIndex: 'resource', width: 80 },
    { title: 'ID', dataIndex: 'resource_id', width: 60, render: (v: number | null) => v ?? '-' },
    { title: '摘要', dataIndex: 'detail', width: 250, ellipsis: { showTitle: false }, render: (v: string) => {
      const text = formatDetail(v);
      return text === '-' ? <span style={{ color: '#999' }}>-</span> : <Tooltip title={text}>{text}</Tooltip>;
    }},
    { title: 'IP', dataIndex: 'ip_address', width: 120 },
  ];

  return (
    <Card title="稽核日誌">
      <Table
        rowKey="id"
        loading={loading}
        columns={columns}
        dataSource={data}
        scroll={{ x: 1010 }}
        pagination={{ current: page, pageSize: PAGE_SIZE, total, onChange: setPage }}
        expandable={{ expandedRowRender: (r) => {
          if (!r.detail) return <span style={{ color: '#999' }}>無詳細資料</span>;
          try {
            const obj = JSON.parse(r.detail);
            return <pre style={{ margin: 0, whiteSpace: 'pre-wrap', fontSize: 12 }}>{JSON.stringify(obj, null, 2)}</pre>;
          } catch { return <pre style={{ margin: 0, whiteSpace: 'pre-wrap' }}>{r.detail}</pre>; }
        }}}
      />
    </Card>
  );
}
