import axios from 'axios';

// ── Types ──────────────────────────────────────────────

export interface Tenant {
  id: string;
  name: string;
  slug: string;
  plan: string;
  max_recipients: number;
  is_active: boolean;
  created_at: string;
}

export interface User {
  id: string;
  tenant_id: string;
  email: string;
  name: string;
  role: string;
  is_active: boolean;
  last_login: string;
}

export interface EmailTemplate {
  id: string;
  tenant_id: string;
  name: string;
  subject: string;
  html_body: string;
  text_body: string;
  category: string;
  language: string;
}

export interface LandingPage {
  id: string;
  tenant_id: string;
  name: string;
  html: string;
  capture_credentials: boolean;
  redirect_url: string;
}

export interface Scenario {
  id: string;
  tenant_id: string;
  name: string;
  category: string;
  difficulty: string;
  language: string;
  template_id: string;
  page_id: string;
  education_html: string;
  is_active: boolean;
  template?: EmailTemplate;
  page?: LandingPage;
}

export interface Recipient {
  id: string;
  email: string;
  first_name: string;
  last_name: string;
  department: string;
  gender: string;
  position: string;
  is_active: boolean;
}

export interface RecipientGroup {
  id: string;
  tenant_id: string;
  name: string;
  recipients?: Recipient[];
}

export interface SMTPProfile {
  id: string;
  tenant_id: string;
  name: string;
  mailer_type: string;
  host: string;
  port: number;
  from_address: string;
  from_name: string;
}

export interface Campaign {
  id: string;
  tenant_id: string;
  name: string;
  status: string;
  scenario_id: string;
  template_id: string;
  page_id: string;
  smtp_profile_id: string;
  phish_url: string;
  selection_mode: string;
  sample_percent: number;
  launched_at: string;
  send_by: string;
  schedule_start: string;
  working_hours_only: boolean;
  skip_weekends: boolean;
  completed_at: string;
  created_at: string;
}

export interface FunnelStats {
  total: number;
  sent: number;
  opened: number;
  clicked: number;
  downloaded: number;
  submitted: number;
  reported: number;
}

export interface DepartmentStat {
  department: string;
  total: number;
  clicked: number;
}

export interface CampaignReport {
  funnel: FunnelStats;
  departments: DepartmentStat[];
}

export interface RecipientResult {
  email: string;
  first_name: string;
  last_name: string;
  department: string;
  status: string;
  sent_at: string | null;
  opened_at: string | null;
  clicked_at: string | null;
  submitted_at: string | null;
  reported_at: string | null;
}

export interface AuditLog {
  id: string;
  tenant_id: string;
  user_id: string;
  user_email: string;
  role: string;
  action: string;
  resource: string;
  resource_id: string;
  detail: string;
  ip_address: string;
  created_at: string;
}

export interface PlatformStats {
  total_tenants: number;
  active_tenants: number;
}

export interface LoginResponse {
  token: string;
  user: User;
}

// ── Axios instance ─────────────────────────────────────

const instance = axios.create({ baseURL: '/api' });

instance.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) config.headers.Authorization = `Bearer ${token}`;
  return config;
});

instance.interceptors.response.use(
  (res) => res,
  (err) => {
    if (err.response?.status === 401) {
      localStorage.removeItem('token');
      window.location.href = '/login';
    }
    // Normalize error message from various API response formats
    const data = err.response?.data;
    const errorMsg = data?.error || data?.message || err.message || '未知錯誤';
    err.displayMessage = errorMsg;
    return Promise.reject(err);
  },
);

// Extract user-friendly error message from API errors
export function getErrorMessage(err: unknown): string {
  if (err && typeof err === 'object' && 'displayMessage' in err) return (err as { displayMessage: string }).displayMessage;
  if (err instanceof Error) return err.message;
  return '未知錯誤';
}

// ── Typed helpers ──────────────────────────────────────

export const api = {
  get: <T>(url: string, params?: Record<string, unknown>) =>
    instance.get<T>(url, { params }).then((r) => r.data),
  post: <T>(url: string, data?: unknown) =>
    instance.post<T>(url, data).then((r) => r.data),
  put: <T>(url: string, data?: unknown) =>
    instance.put<T>(url, data).then((r) => r.data),
  patch: <T>(url: string, data?: unknown) =>
    instance.patch<T>(url, data).then((r) => r.data),
  del: <T>(url: string) =>
    instance.delete<T>(url).then((r) => r.data),
};
