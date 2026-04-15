import { useState } from 'react';
import { Layout, Menu, Button, Typography } from 'antd';
import {
  DashboardOutlined, SendOutlined, AppstoreOutlined,
  FileTextOutlined, TeamOutlined, SettingOutlined,
} from '@ant-design/icons';
import { Outlet, useNavigate, useLocation } from 'react-router-dom';
import { useAuth } from '../hooks/useAuth';

const { Header, Sider, Content } = Layout;

const menuItems = [
  { key: '/app/dashboard', icon: <DashboardOutlined />, label: 'Dashboard' },
  { key: '/app/campaigns', icon: <SendOutlined />, label: '釣魚測試' },
  { key: '/app/scenarios', icon: <AppstoreOutlined />, label: '情境庫' },
  { key: '/app/templates', icon: <FileTextOutlined />, label: '模板管理' },
  { key: '/app/recipients', icon: <TeamOutlined />, label: '收件人' },
  {
    key: 'settings',
    icon: <SettingOutlined />,
    label: '設定',
    children: [
      { key: '/app/settings/smtp', label: 'SMTP 設定' },
      { key: '/app/settings/audit', label: '稽核日誌' },
    ],
  },
];

export default function TenantLayout() {
  const { user, logout } = useAuth();
  const navigate = useNavigate();
  const { pathname } = useLocation();
  const [collapsed, setCollapsed] = useState(false);

  const openKey = pathname.startsWith('/app/settings') ? 'settings' : '';

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sider collapsible collapsed={collapsed} onCollapse={setCollapsed} breakpoint="lg">
        <div style={{ padding: 16, textAlign: 'center' }}>
          <Typography.Text strong style={{ color: '#fff' }}>
            {collapsed ? 'PG' : 'PhishGuard'}
          </Typography.Text>
        </div>
        <Menu
          theme="dark"
          mode="inline"
          selectedKeys={[pathname]}
          defaultOpenKeys={openKey ? [openKey] : []}
          items={menuItems}
          onClick={({ key }) => navigate(key)}
        />
      </Sider>
      <Layout>
        <Header style={{ display: 'flex', justifyContent: 'flex-end', alignItems: 'center', gap: 12, padding: '0 24px' }}>
          {user?.tenant_id && (
            <Typography.Text style={{ color: '#fff' }}>{user.tenant_id}</Typography.Text>
          )}
          <Typography.Text style={{ color: '#fff' }}>{user?.name}</Typography.Text>
          <Button size="small" onClick={logout}>Logout</Button>
        </Header>
        <Content style={{ padding: 24 }}>
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  );
}
