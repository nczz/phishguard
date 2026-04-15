import { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {
  Card, Descriptions, Table, Tabs, Tag, Badge, Button, Switch, Modal,
  Form, Input, Select, Popconfirm, Spin, Space, Typography, message,
} from 'antd';
import { api } from '../../api/client';
import type { Tenant, Campaign, User } from '../../api/client';
import dayjs from 'dayjs';

const { Title } = Typography;

const planColor: Record<string, string> = { free: 'default', pro: 'blue', enterprise: 'gold' };
const statusColor: Record<string, string> = { draft: 'default', running: 'processing', completed: 'green', paused: 'orange' };

export default function TenantDetail() {
  const { id } = useParams<{ id: string }>();
  const nav = useNavigate();
  const [tenant, setTenant] = useState<Tenant | null>(null);
  const [campaigns, setCampaigns] = useState<Campaign[]>([]);
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [toggling, setToggling] = useState(false);
  const [modalOpen, setModalOpen] = useState(false);
  const [form] = Form.useForm();
  const [submitting, setSubmitting] = useState(false);
  const [editUser, setEditUser] = useState<User | null>(null);
  const [editForm] = Form.useForm();

  const fetchAll = () => {
    setLoading(true);
    Promise.all([
      api.get<Tenant>('/admin/tenants/' + id),
      api.get<Campaign[]>('/admin/tenants/' + id + '/campaigns'),
      api.get<User[]>('/admin/tenants/' + id + '/users'),
    ])
      .then(([t, c, u]) => { setTenant(t); setCampaigns(c); setUsers(u); })
      .catch(() => message.error('載入失敗'))
      .finally(() => setLoading(false));
  };

  useEffect(fetchAll, [id]);

  const handleToggle = async (checked: boolean) => {
    setToggling(true);
    try {
      await api.patch('/admin/tenants/' + id + '/toggle', { is_active: checked });
      fetchAll();
    } catch { message.error('切換失敗'); }
    finally { setToggling(false); }
  };

  const handleImpersonate = async () => {
    try {
      const res = await api.post<{ token: string; tenant_id: string }>('/admin/tenants/' + id + '/impersonate');
      localStorage.setItem('impersonate_from', localStorage.getItem('token') || '');
      localStorage.setItem('impersonate_token', res.token);
      nav('/app/dashboard');
    } catch { message.error('模擬失敗'); }
  };

  const handleAddUser = async (values: Record<string, string>) => {
    setSubmitting(true);
    try {
      await api.post('/admin/tenants/' + id + '/users', values);
      message.success('使用者已新增');
      setModalOpen(false);
      form.resetFields();
      fetchAll();
    } catch { message.error('新增失敗'); }
    finally { setSubmitting(false); }
  };

  const handleDeleteUser = async (uid: string) => {
    try {
      await api.del('/admin/tenants/' + id + '/users/' + uid);
      message.success('已刪除');
      fetchAll();
    } catch { message.error('刪除失敗'); }
  };

  const openEditUser = (u: User) => {
    setEditUser(u);
    editForm.setFieldsValue({ name: u.name, role: u.role, is_active: u.is_active, password: '' });
  };

  const handleEditUser = async (values: { name: string; role: string; is_active: boolean; password: string }) => {
    if (!editUser) return;
    try {
      const payload: Record<string, unknown> = { name: values.name, role: values.role, is_active: values.is_active };
      if (values.password) payload.password = values.password;
      await api.put('/admin/tenants/' + id + '/users/' + editUser.id, payload);
      message.success('已更新');
      setEditUser(null);
      fetchAll();
    } catch { message.error('更新失敗'); }
  };

  if (loading) return <Spin style={{ display: 'block', margin: '20vh auto' }} size="large" />;
  if (!tenant) return null;

  return (
    <div style={{ padding: 24 }}>
      {/* Header */}
      <Space align="center" style={{ marginBottom: 16 }}>
        <Button onClick={() => nav('/admin/tenants')}>← 返回列表</Button>
        <Title level={4} style={{ margin: 0 }}>{tenant.name}</Title>
        <Badge status={tenant.is_active ? 'success' : 'error'} text={tenant.is_active ? '啟用' : '停用'} />
        <Tag color={planColor[tenant.plan]}>{tenant.plan}</Tag>
      </Space>

      {/* Action buttons */}
      <Space style={{ marginBottom: 16 }}>
        <Switch
          checked={tenant.is_active}
          loading={toggling}
          onChange={handleToggle}
          checkedChildren="啟用"
          unCheckedChildren="停用"
        />
        <Button type="primary" onClick={handleImpersonate}>以此租戶身份操作</Button>
      </Space>

      {/* Info card */}
      <Card style={{ marginBottom: 16 }}>
        <Descriptions column={2} bordered>
          <Descriptions.Item label="識別碼">{tenant.slug}</Descriptions.Item>
          <Descriptions.Item label="方案"><Tag color={planColor[tenant.plan]}>{tenant.plan}</Tag></Descriptions.Item>
          <Descriptions.Item label="收件人上限">{tenant.max_recipients}</Descriptions.Item>
          <Descriptions.Item label="建立時間">{dayjs(tenant.created_at).format('YYYY-MM-DD HH:mm')}</Descriptions.Item>
        </Descriptions>
      </Card>

      {/* Tabs */}
      <Tabs items={[
        {
          key: 'campaigns',
          label: '活動列表',
          children: (
            <Table<Campaign> rowKey="id" dataSource={campaigns} pagination={false} columns={[
              { title: '名稱', dataIndex: 'name' },
              { title: '狀態', dataIndex: 'status', render: (s: string) => <Tag color={statusColor[s]}>{s}</Tag> },
              { title: '發送時間', dataIndex: 'launched_at', render: (v: string) => v ? dayjs(v).format('YYYY-MM-DD HH:mm') : '-' },
              { title: '建立時間', dataIndex: 'created_at', render: (v: string) => dayjs(v).format('YYYY-MM-DD HH:mm') },
            ]} />
          ),
        },
        {
          key: 'users',
          label: '使用者管理',
          children: (
            <>
              <Button type="primary" style={{ marginBottom: 12 }} onClick={() => setModalOpen(true)}>+ 新增使用者</Button>
              <Table<User> rowKey="id" dataSource={users} pagination={false} columns={[
                { title: 'Email', dataIndex: 'email' },
                { title: '姓名', dataIndex: 'name' },
                { title: '角色', dataIndex: 'role', render: (r: string) => <Tag>{r}</Tag> },
                { title: '狀態', dataIndex: 'is_active', render: (v: boolean) => <Badge status={v ? 'success' : 'error'} text={v ? '啟用' : '停用'} /> },
                { title: '最後登入', dataIndex: 'last_login', render: (v: string) => v ? dayjs(v).format('YYYY-MM-DD HH:mm') : '-' },
                {
                  title: '操作', render: (_, record) => (
                    <Space>
                      <Button size="small" onClick={() => openEditUser(record)}>編輯</Button>
                      <Popconfirm title="確定刪除此使用者？" onConfirm={() => handleDeleteUser(record.id)}>
                        <Button danger size="small">刪除</Button>
                      </Popconfirm>
                    </Space>
                  ),
                },
              ]} />
            </>
          ),
        },
      ]} />

      {/* Add user modal */}
      <Modal
        title="新增使用者"
        open={modalOpen}
        onCancel={() => setModalOpen(false)}
        onOk={() => form.submit()}
        confirmLoading={submitting}
        destroyOnHidden
      >
        <Form form={form} layout="vertical" onFinish={handleAddUser}>
          <Form.Item label="Email" name="email" rules={[{ required: true, type: 'email', message: '請輸入有效 Email' }]}>
            <Input />
          </Form.Item>
          <Form.Item label="姓名" name="name" rules={[{ required: true, message: '請輸入姓名' }]}>
            <Input />
          </Form.Item>
          <Form.Item label="密碼" name="password" rules={[{ required: true, min: 8, message: '密碼至少 8 字元' }]}>
            <Input.Password />
          </Form.Item>
          <Form.Item label="角色" name="role" rules={[{ required: true, message: '請選擇角色' }]} initialValue="operator">
            <Select options={[
              { value: 'tenant_admin', label: '租戶管理員' },
              { value: 'operator', label: '操作員' },
              { value: 'viewer', label: '檢視者' },
            ]} />
          </Form.Item>
        </Form>
      </Modal>

      {/* Edit user modal */}
      <Modal title="編輯使用者" open={!!editUser} onCancel={() => setEditUser(null)} onOk={() => editForm.submit()} destroyOnHidden>
        <Form form={editForm} layout="vertical" onFinish={handleEditUser}>
          <Form.Item label="姓名" name="name" rules={[{ required: true }]}><Input /></Form.Item>
          <Form.Item label="角色" name="role" rules={[{ required: true }]}>
            <Select options={[
              { value: 'tenant_admin', label: '租戶管理員' },
              { value: 'operator', label: '操作員' },
              { value: 'viewer', label: '檢視者' },
            ]} />
          </Form.Item>
          <Form.Item label="啟用" name="is_active" valuePropName="checked">
            <Switch checkedChildren="啟用" unCheckedChildren="停用" />
          </Form.Item>
          <Form.Item label="重設密碼（留空不變更）" name="password">
            <Input.Password placeholder="輸入新密碼（至少 8 字元）" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
}
