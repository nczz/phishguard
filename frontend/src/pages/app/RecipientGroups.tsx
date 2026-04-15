import { useEffect, useState } from 'react';
import { Card, List, Table, Modal, Form, Input, Button, Tabs, Space, Empty, Typography, message, Spin } from 'antd';
import { PlusOutlined, UploadOutlined } from '@ant-design/icons';
import { api } from '../../api/client';
import type { RecipientGroup, Recipient } from '../../api/client';

const recipientColumns = [
  { title: 'Email', dataIndex: 'email' },
  { title: '名', dataIndex: 'first_name' },
  { title: '姓', dataIndex: 'last_name' },
  { title: '部門', dataIndex: 'department' },
  { title: '職稱', dataIndex: 'position' },
];

export default function RecipientGroups() {
  const [groups, setGroups] = useState<RecipientGroup[]>([]);
  const [loading, setLoading] = useState(true);
  const [selected, setSelected] = useState<RecipientGroup | null>(null);
  const [newOpen, setNewOpen] = useState(false);
  const [importOpen, setImportOpen] = useState(false);
  const [newName, setNewName] = useState('');
  const [jsonText, setJsonText] = useState('');
  const [importTab, setImportTab] = useState('json');
  const [form] = Form.useForm();

  const fetchGroups = () => {
    setLoading(true);
    api.get<RecipientGroup[]>('/recipient-groups').then(setGroups).finally(() => setLoading(false));
  };

  useEffect(fetchGroups, []);

  const selectGroup = (g: RecipientGroup) => {
    api.get<RecipientGroup[]>('/recipient-groups').then((all) => {
      const found = all.find((x) => x.id === g.id);
      setSelected(found ?? g);
    });
  };

  const handleCreate = async () => {
    if (!newName.trim()) return;
    await api.post('/recipient-groups', { name: newName.trim() });
    message.success('群組已建立');
    setNewOpen(false);
    setNewName('');
    fetchGroups();
  };

  const handleImport = async () => {
    let recipients: Omit<Recipient, 'id'>[];
    if (importTab === 'json') {
      try {
        recipients = JSON.parse(jsonText);
        if (!Array.isArray(recipients)) throw new Error();
      } catch {
        message.error('JSON 格式錯誤');
        return;
      }
    } else {
      const values = await form.validateFields();
      recipients = [values];
    }
    const res = await api.post<{ count: number }>('/recipient-groups/import', {
      group_id: selected!.id,
      recipients,
    });
    message.success(`成功匯入 ${res.count ?? recipients.length} 筆收件人`);
    setImportOpen(false);
    setJsonText('');
    form.resetFields();
    selectGroup(selected!);
    fetchGroups();
  };

  if (loading) return <Spin style={{ display: 'block', margin: '20vh auto' }} size="large" />;

  return (
    <div>
      <Typography.Title level={3}>收件人管理</Typography.Title>
      <div style={{ display: 'flex', gap: 16 }}>
        {/* Left: group list */}
        <Card style={{ width: 320, flexShrink: 0 }} title="群組" extra={<Button icon={<PlusOutlined />} size="small" onClick={() => setNewOpen(true)}>新增群組</Button>}>
          {groups.length === 0 ? (
            <Empty description="尚無群組" />
          ) : (
            <List
              dataSource={groups}
              renderItem={(g) => (
                <List.Item
                  style={{ cursor: 'pointer', background: selected?.id === g.id ? '#e6f4ff' : undefined, padding: '8px 12px' }}
                  onClick={() => selectGroup(g)}
                >
                  <List.Item.Meta title={g.name} description={`${g.recipients?.length ?? 0} 位收件人`} />
                </List.Item>
              )}
            />
          )}
        </Card>

        {/* Right: detail */}
        <Card style={{ flex: 1 }}>
          {selected ? (
            <>
              <Space style={{ marginBottom: 16, justifyContent: 'space-between', width: '100%' }}>
                <Typography.Title level={4} style={{ margin: 0 }}>{selected.name}</Typography.Title>
                <Button icon={<UploadOutlined />} onClick={() => setImportOpen(true)}>匯入收件人</Button>
              </Space>
              <Table
                rowKey="id"
                columns={recipientColumns}
                dataSource={selected.recipients ?? []}
                pagination={{ pageSize: 10 }}
              />
            </>
          ) : (
            <Empty description="請選擇一個群組" />
          )}
        </Card>
      </div>

      {/* New group modal */}
      <Modal title="新增群組" open={newOpen} onOk={handleCreate} onCancel={() => setNewOpen(false)} okText="建立">
        <Input placeholder="群組名稱" value={newName} onChange={(e) => setNewName(e.target.value)} />
      </Modal>

      {/* Import modal */}
      <Modal title="匯入收件人" open={importOpen} onOk={handleImport} onCancel={() => setImportOpen(false)} okText="匯入" width={600}>
        <Tabs activeKey={importTab} onChange={setImportTab} items={[
          {
            key: 'json',
            label: 'JSON 匯入',
            children: (
              <Input.TextArea
                rows={8}
                placeholder='[{"email":"a@b.com","first_name":"","last_name":"","department":"","position":""}]'
                value={jsonText}
                onChange={(e) => setJsonText(e.target.value)}
              />
            ),
          },
          {
            key: 'form',
            label: '表單新增',
            children: (
              <Form form={form} layout="vertical">
                <Form.Item name="email" label="Email" rules={[{ required: true, type: 'email' }]}>
                  <Input />
                </Form.Item>
                <Form.Item name="first_name" label="名"><Input /></Form.Item>
                <Form.Item name="last_name" label="姓"><Input /></Form.Item>
                <Form.Item name="department" label="部門"><Input /></Form.Item>
                <Form.Item name="position" label="職稱"><Input /></Form.Item>
              </Form>
            ),
          },
        ]} />
      </Modal>
    </div>
  );
}
