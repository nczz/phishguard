import { useState } from 'react';
import { Layout, Menu, Button, Typography, Alert } from 'antd';
import {
  DashboardOutlined, SendOutlined, AppstoreOutlined,
  FileTextOutlined, TeamOutlined, SettingOutlined, BookOutlined, LayoutOutlined, NodeIndexOutlined,
  BarChartOutlined,
} from '@ant-design/icons';
import { Outlet, useNavigate, useLocation } from 'react-router-dom';
import { useAuth } from '../hooks/useAuth';

const { Header, Sider, Content } = Layout;

const menuItems = [
  { key: '/app/dashboard', icon: <DashboardOutlined />, label: 'Dashboard' },
  { key: '/app/campaigns', icon: <SendOutlined />, label: '釣魚測試' },
  { key: '/app/scenarios', icon: <AppstoreOutlined />, label: '情境庫' },
  { key: '/app/templates', icon: <FileTextOutlined />, label: '模板管理' },
  { key: '/app/pages', icon: <LayoutOutlined />, label: 'Landing Page' },
  { key: '/app/recipients', icon: <TeamOutlined />, label: '收件人' },
  {
    key: 'reports',
    icon: <BarChartOutlined />,
    label: '報表',
    children: [
      { key: '/app/reports/offenders', label: '累犯追蹤' },
      { key: '/app/reports/trend', label: '趨勢分析' },
    ],
  },
  {
    key: 'settings',
    icon: <SettingOutlined />,
    label: '設定',
    children: [
      { key: '/app/settings/smtp', label: 'SMTP 設定' },
      { key: '/app/settings/audit', label: '稽核日誌' },
      { key: '/app/settings/auto-test', label: '自動測試' },
    ],
  },
  { key: '/app/guide', icon: <BookOutlined />, label: '使用指南' },
  { key: '/app/flow', icon: <NodeIndexOutlined />, label: '流程總覽' },
];

export default function TenantLayout() {
  const { user, logout, impersonating, exitImpersonation } = useAuth();
  const navigate = useNavigate();
  const { pathname } = useLocation();
  const [collapsed, setCollapsed] = useState(false);

  const openKey = pathname.startsWith('/app/settings') ? 'settings' : pathname.startsWith('/app/reports') ? 'reports' : '';

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
          {impersonating && (
            <Alert
              type="warning"
              banner
              message="您正在以租戶身份操作"
              action={<Button size="small" onClick={exitImpersonation}>返回管理員</Button>}
              style={{ marginBottom: 16 }}
            />
          )}
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  );
}
