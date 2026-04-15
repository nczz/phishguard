import { createContext, useContext, useEffect, useState } from 'react';
import type { ReactNode } from 'react';
import { useNavigate, Navigate } from 'react-router-dom';
import { Spin } from 'antd';
import { api } from '../api/client';
import type { User, LoginResponse } from '../api/client';

interface AuthContextValue {
  user: User | null;
  token: string | null;
  login: (email: string, password: string) => Promise<void>;
  logout: () => void;
  loading: boolean;
  impersonating: boolean;
  impersonate: (token: string) => void;
  exitImpersonation: () => void;
}

const AuthContext = createContext<AuthContextValue | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState<string | null>(() => localStorage.getItem('token'));
  const [loading, setLoading] = useState(!!token);
  const [impersonating, setImpersonating] = useState(() => !!localStorage.getItem('original_token'));
  const navigate = useNavigate();

  useEffect(() => {
    if (!token) return;
    api.get<User>('/auth/me')
      .then(setUser)
      .catch(() => {
        localStorage.removeItem('token');
        setToken(null);
      })
      .finally(() => setLoading(false));
  }, [token]);

  const login = async (email: string, password: string) => {
    const res = await api.post<LoginResponse>('/auth/login', { email, password });
    localStorage.setItem('token', res.token);
    setToken(res.token);
    setUser(res.user);
  };

  const logout = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('original_token');
    setToken(null);
    setUser(null);
    setImpersonating(false);
    navigate('/login');
  };

  const impersonate = (newToken: string) => {
    localStorage.setItem('original_token', token!);
    localStorage.setItem('token', newToken);
    setToken(newToken);
    setImpersonating(true);
  };

  const exitImpersonation = () => {
    const original = localStorage.getItem('original_token');
    if (!original) return;
    localStorage.setItem('token', original);
    localStorage.removeItem('original_token');
    setToken(original);
    setImpersonating(false);
    navigate('/admin/dashboard');
  };

  return (
    <AuthContext.Provider value={{ user, token, login, logout, loading, impersonating, impersonate, exitImpersonation }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) return { user: null, token: null, loading: true, impersonating: false, login: async () => {}, logout: () => {}, impersonate: () => {}, exitImpersonation: () => {} } as AuthContextValue;
  return ctx;
}

export function ProtectedRoute({ children, role }: { children: ReactNode; role?: string }) {
  const { user, loading } = useAuth();
  if (loading) return <Spin style={{ display: 'block', margin: '20vh auto' }} size="large" />;
  if (!user) return <Navigate to="/login" replace />;
  if (role && user.role !== role) return <Navigate to="/login" replace />;
  return <>{children}</>;
}
