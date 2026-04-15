import { lazy, Suspense } from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { ConfigProvider, Spin } from 'antd';
import zhTW from 'antd/locale/zh_TW';
import { AuthProvider, ProtectedRoute, useAuth } from './hooks/useAuth';
import AdminLayout from './layouts/AdminLayout';
import TenantLayout from './layouts/TenantLayout';

// ── Lazy pages ─────────────────────────────────────────

const Login = lazy(() => import('./pages/Login'));
const AdminDashboard = lazy(() => import('./pages/admin/AdminDashboard'));
const TenantList = lazy(() => import('./pages/admin/TenantList'));
const TenantCreate = lazy(() => import('./pages/admin/TenantCreate'));
const Dashboard = lazy(() => import('./pages/app/Dashboard'));
const CampaignList = lazy(() => import('./pages/app/CampaignList'));
const CampaignWizard = lazy(() => import('./pages/app/CampaignWizard'));
const CampaignDetail = lazy(() => import('./pages/app/CampaignDetail'));
const ScenarioList = lazy(() => import('./pages/app/ScenarioList'));
const TemplateList = lazy(() => import('./pages/app/TemplateList'));
const RecipientGroups = lazy(() => import('./pages/app/RecipientGroups'));
const SMTPSettings = lazy(() => import('./pages/app/SMTPSettings'));
const AuditLogs = lazy(() => import('./pages/app/AuditLogs'));

const fallback = <Spin style={{ display: 'block', margin: '20vh auto' }} size="large" />;

// ── Root redirect ──────────────────────────────────────

function RootRedirect() {
  const { user, loading } = useAuth();
  if (loading) return fallback;
  if (!user) return <Navigate to="/login" replace />;
  return <Navigate to={user.role === 'platform_admin' ? '/admin/dashboard' : '/app/dashboard'} replace />;
}

// ── App ────────────────────────────────────────────────

export default function App() {
  return (
    <ConfigProvider locale={zhTW}>
      <BrowserRouter>
        <AuthProvider>
          <Suspense fallback={fallback}>
            <Routes>
              <Route path="/login" element={<Login />} />

              <Route path="/admin" element={<ProtectedRoute role="platform_admin"><AdminLayout /></ProtectedRoute>}>
                <Route path="dashboard" element={<AdminDashboard />} />
                <Route path="tenants" element={<TenantList />} />
                <Route path="tenants/new" element={<TenantCreate />} />
              </Route>

              <Route path="/app" element={<ProtectedRoute><TenantLayout /></ProtectedRoute>}>
                <Route path="dashboard" element={<Dashboard />} />
                <Route path="campaigns" element={<CampaignList />} />
                <Route path="campaigns/new" element={<CampaignWizard />} />
                <Route path="campaigns/:id" element={<CampaignDetail />} />
                <Route path="scenarios" element={<ScenarioList />} />
                <Route path="templates" element={<TemplateList />} />
                <Route path="recipients" element={<RecipientGroups />} />
                <Route path="settings/smtp" element={<SMTPSettings />} />
                <Route path="settings/audit" element={<AuditLogs />} />
              </Route>

              <Route path="/" element={<RootRedirect />} />
            </Routes>
          </Suspense>
        </AuthProvider>
      </BrowserRouter>
    </ConfigProvider>
  );
}
