import { createContext } from 'react';
import type { User } from '../api/client';

export interface AuthContextValue {
  user: User | null;
  token: string | null;
  login: (email: string, password: string) => Promise<void>;
  logout: () => void;
  loading: boolean;
  impersonating: boolean;
  impersonate: (token: string) => void;
  exitImpersonation: () => void;
}

export const AuthContext = createContext<AuthContextValue | null>(null);
