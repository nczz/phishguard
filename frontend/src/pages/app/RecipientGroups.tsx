import { useEffect, useState, useMemo } from 'react';
import { Card, Table, Button, Modal, Upload, Tag, message, Row, Col, Statistic, Select, Typography, Empty, Space } from 'antd';
import { UploadOutlined, DownloadOutlined, TeamOutlined, InboxOutlined } from '@ant-design/icons';
import { api } from '../../api/client';
import type { RecipientGroup, Recipient } from '../../api/client';

const CSV_TEMPLATE = `email,first_name,last_name,department,gender,position
wang@example.com,小明,王,業務部,男,業務經理
chen@example.com,小華,陳,財務部,女,會計
lin@example.com,小美,林,研發部,不指定,工程師
`;

interface ParsedRow { email: string; first_name: string; last_name: string; department: string; gender: string; position: string }

function parseCSV(text: string): ParsedRow[] {
  const lines = text.trim().split(/\r?\n/);
  if (lines.length < 2) return [];
  const header = lines[0].toLowerCase().split(',').map(h => h.trim());
  const ei = header.indexOf('email');
  const fi = header.indexOf('first_name');
  const li = header.indexOf('last_name');
  const di = header.indexOf('department');
  const gi = header.indexOf('gender');
  const pi = header.indexOf('position');
  if (ei < 0) { message.error('CSV 必須包含 email 欄位'); return []; }
  return lines.slice(1).filter(l => l.trim()).map(line => {
    const cols = line.split(',').map(c => c.trim());
    return { email: cols[ei] || '', first_name: cols[fi] ?? '', last_name: cols[li] ?? '', department: cols[di] ?? '', gender: cols[gi] ?? '不指定', position: cols[pi] ?? '' };
  }).filter(r => r.email.includes('@'));
}

export default function RecipientGroups() {
  const [groups, setGroups] = useState<RecipientGroup[]>([]);
  const [allRecipients, setAllRecipients] = useState<Recipient[]>([]);
  const [importOpen, setImportOpen] = useState(false);
  const [parsed, setParsed] = useState<ParsedRow[]>([]);
  const [importing, setImporting] = useState(false);
  const [deptFilter, setDeptFilter] = useState<string | undefined>();

  const load = () => {
    api.get<RecipientGroup[]>('/recipient-groups').then(gs => {
      setGroups(gs);
      const all = gs.flatMap(g => g.recipients ?? []);
      setAllRecipients(all);
    });
  };

  useEffect(load, []);

  // Department stats
  const deptStats = useMemo(() => {
    const map = new Map<string, number>();
    allRecipients.forEach(r => { const d = r.department || '未分類'; map.set(d, (map.get(d) || 0) + 1); });
    return Array.from(map.entries()).sort((a, b) => b[1] - a[1]);
  }, [allRecipients]);

  const departments = useMemo(() => deptStats.map(([d]) => d), [deptStats]);

  const filteredRecipients = useMemo(() => {
    if (!deptFilter) return allRecipients;
    return allRecipients.filter(r => (r.department || '未分類') === deptFilter);
  }, [allRecipients, deptFilter]);

  // CSV template download
  const downloadTemplate = () => {
    const blob = new Blob([CSV_TEMPLATE], { type: 'text/csv;charset=utf-8;' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a'); a.href = url; a.download = 'phishguard_員工名單範本.csv'; a.click();
    URL.revokeObjectURL(url);
  };

  // File upload handler
  const handleFile = (file: File) => {
    const reader = new FileReader();
    reader.onload = (e) => {
      const text = e.target?.result as string;
      const rows = parseCSV(text);
      if (rows.length === 0) { message.error('無法解析 CSV 或沒有有效資料'); return; }
      setParsed(rows);
      setImportOpen(true);
    };
    reader.readAsText(file);
    return false; // prevent auto upload
  };

  // Import confirmed
  const doImport = async () => {
    if (parsed.length === 0) return;
    setImporting(true);
    try {
      // Create or find the "全公司" group
      let group = groups.find(g => g.name === '全公司');
      if (!group) {
        group = await api.post<RecipientGroup>('/recipient-groups', { name: '全公司' });
      }
      const res = await api.post<{ created: number; updated: number; total: number }>('/recipient-groups/import', { group_id: group.id, recipients: parsed });
      message.success(`匯入完成：新增 ${res.created} 人，更新 ${res.updated} 人`);
      setImportOpen(false);
      setParsed([]);
      load();
    } catch { message.error('匯入失敗'); }
    setImporting(false);
  };

  // Preview table columns
  const previewDepts = useMemo(() => {
    const map = new Map<string, number>();
    parsed.forEach(r => { const d = r.department || '未分類'; map.set(d, (map.get(d) || 0) + 1); });
    return Array.from(map.entries()).sort((a, b) => b[1] - a[1]);
  }, [parsed]);

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16 }}>
        <Typography.Title level={3} style={{ margin: 0 }}>收件人管理</Typography.Title>
        <Space>
          <Button icon={<DownloadOutlined />} onClick={downloadTemplate}>下載 CSV 範本</Button>
          <Upload accept=".csv,.txt" showUploadList={false} beforeUpload={handleFile}>
            <Button type="primary" icon={<UploadOutlined />}>匯入員工名單</Button>
          </Upload>
        </Space>
      </div>

      {/* Department overview */}
      {deptStats.length > 0 && (
        <Row gutter={[12, 12]} style={{ marginBottom: 24 }}>
          <Col xs={12} sm={6}>
            <Card size="small"><Statistic title="總人數" value={allRecipients.length} prefix={<TeamOutlined />} /></Card>
          </Col>
          <Col xs={12} sm={6}>
            <Card size="small"><Statistic title="部門數" value={deptStats.length} /></Card>
          </Col>
          {deptStats.slice(0, 4).map(([dept, count]) => (
            <Col xs={12} sm={6} key={dept}>
              <Card size="small" hoverable onClick={() => setDeptFilter(deptFilter === dept ? undefined : dept)}
                style={{ borderColor: deptFilter === dept ? '#1677ff' : undefined }}>
                <Statistic title={dept} value={count} suffix="人" />
              </Card>
            </Col>
          ))}
        </Row>
      )}

      {/* Employee table */}
      <Card title={
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <span>員工名單 {deptFilter && <Tag closable onClose={() => setDeptFilter(undefined)}>{deptFilter}</Tag>}</span>
          {departments.length > 0 && (
            <Select allowClear placeholder="篩選部門" style={{ width: 160 }} value={deptFilter} onChange={setDeptFilter}
              options={departments.map(d => ({ value: d, label: d }))} />
          )}
        </div>
      }>
        {allRecipients.length === 0 ? (
          <Empty description="尚未匯入員工名單" image={<InboxOutlined style={{ fontSize: 48, color: '#ccc' }} />}>
            <Space>
              <Button onClick={downloadTemplate}>下載 CSV 範本</Button>
              <Upload accept=".csv,.txt" showUploadList={false} beforeUpload={handleFile}>
                <Button type="primary">匯入第一份名單</Button>
              </Upload>
            </Space>
          </Empty>
        ) : (
          <Table dataSource={filteredRecipients} rowKey="id" size="small" pagination={{ pageSize: 20 }}
            columns={[
              { title: 'Email', dataIndex: 'email', sorter: (a: Recipient, b: Recipient) => a.email.localeCompare(b.email) },
              { title: '姓', dataIndex: 'last_name', width: 80 },
              { title: '名', dataIndex: 'first_name', width: 80 },
              { title: '部門', dataIndex: 'department', width: 120, render: (d: string) => <Tag>{d || '未分類'}</Tag>,
                filters: departments.map(d => ({ text: d, value: d })), onFilter: (v: unknown, r: Recipient) => r.department === v },
              { title: '性別', dataIndex: 'gender', width: 80, render: (g: string) => g || '不指定',
                filters: [{ text: '男', value: '男' }, { text: '女', value: '女' }, { text: '不指定', value: '不指定' }], onFilter: (v: unknown, r: Recipient) => (r.gender || '不指定') === v },
              { title: '職稱', dataIndex: 'position', width: 120 },
            ]}
          />
        )}
      </Card>

      {/* Import preview modal */}
      <Modal title="確認匯入" open={importOpen} onCancel={() => setImportOpen(false)} width={800}
        onOk={doImport} okText={`確認匯入 ${parsed.length} 人`} confirmLoading={importing}>
        <div style={{ marginBottom: 16 }}>
          <Typography.Text>解析到 <strong>{parsed.length}</strong> 位員工，分佈在 <strong>{previewDepts.length}</strong> 個部門：</Typography.Text>
          <div style={{ marginTop: 8 }}>
            {previewDepts.map(([dept, count]) => (
              <Tag key={dept} color="blue" style={{ marginBottom: 4 }}>{dept}: {count} 人</Tag>
            ))}
          </div>
        </div>
        <Table dataSource={parsed} rowKey="email" size="small" pagination={{ pageSize: 10 }}
          columns={[
            { title: 'Email', dataIndex: 'email' },
            { title: '姓', dataIndex: 'last_name', width: 80 },
            { title: '名', dataIndex: 'first_name', width: 80 },
            { title: '部門', dataIndex: 'department', width: 120 },
            { title: '性別', dataIndex: 'gender', width: 80 },
            { title: '職稱', dataIndex: 'position', width: 120 },
          ]}
        />
      </Modal>
    </div>
  );
}
