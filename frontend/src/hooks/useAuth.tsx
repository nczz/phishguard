import { useContext } from 'react';
import { AuthContext } from './authContext';
import type { AuthContextValue } from './authContext';

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) return { user: null, token: null, loading: true, impersonating: false, login: async () => {}, logout: () => {}, impersonate: () => {}, exitImpersonation: () => {} } as AuthContextValue;
  return ctx;
}
