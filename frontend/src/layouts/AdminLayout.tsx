import { Layout, Menu, Button, Typography } from 'antd';
import { Outlet, useNavigate, useLocation } from 'react-router-dom';
import { useAuth } from '../hooks/useAuth';

const { Header, Content } = Layout;

const items = [
  { key: '/admin/dashboard', label: 'Dashboard' },
  { key: '/admin/tenants', label: 'Tenants' },
];

export default function AdminLayout() {
  const { user, logout } = useAuth();
  const navigate = useNavigate();
  const { pathname } = useLocation();

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Header style={{ display: 'flex', alignItems: 'center', gap: 16 }}>
        <Typography.Text strong style={{ color: '#fff', whiteSpace: 'nowrap' }}>
          PhishGuard Admin
        </Typography.Text>
        <Menu
          theme="dark"
          mode="horizontal"
          selectedKeys={[pathname]}
          items={items}
          onClick={({ key }) => navigate(key)}
          style={{ flex: 1 }}
        />
        <Typography.Text style={{ color: '#fff', whiteSpace: 'nowrap' }}>
          {user?.email}
        </Typography.Text>
        <Button size="small" onClick={logout}>Logout</Button>
      </Header>
      <Content style={{ padding: 24 }}>
        <Outlet />
      </Content>
    </Layout>
  );
}
