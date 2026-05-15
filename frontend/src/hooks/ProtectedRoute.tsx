import type { ReactNode } from 'react';
import { Navigate } from 'react-router-dom';
import { Spin } from 'antd';
import { useAuth } from './useAuth';

export function ProtectedRoute({ children, role }: { children: ReactNode; role?: string }) {
  const { user, loading } = useAuth();
  if (loading) return <Spin style={{ display: 'block', margin: '20vh auto' }} size="large" />;
  if (!user) return <Navigate to="/login" replace />;
  if (role && user.role !== role) return <Navigate to="/login" replace />;
  return <>{children}</>;
}
