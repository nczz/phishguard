import { useEffect, useState } from 'react';
import type { ReactNode } from 'react';
import { useNavigate } from 'react-router-dom';
import { api } from '../api/client';
import type { User, LoginResponse } from '../api/client';
import { AuthContext } from './authContext';

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
    window.location.href = '/admin/dashboard';
  };

  return (
    <AuthContext.Provider value={{ user, token, login, logout, loading, impersonating, impersonate, exitImpersonation }}>
      {children}
    </AuthContext.Provider>
  );
}
