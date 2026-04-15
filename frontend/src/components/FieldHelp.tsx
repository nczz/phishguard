import { Tooltip } from 'antd';
import { QuestionCircleOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';



interface Props {
  label: string;
  tip: string;
  guideAnchor?: string; // links to /app/guide#anchor
}

export default function FieldHelp({ label, tip, guideAnchor }: Props) {
  const nav = useNavigate();
  return (
    <span>
      {label}{' '}
      <Tooltip title={<span>{tip}{guideAnchor && <><br /><a onClick={() => nav(`/app/guide#${guideAnchor}`)} style={{ color: '#69c0ff' }}>查看使用指南 →</a></>}</span>}>
        <QuestionCircleOutlined style={{ color: '#999', cursor: 'help' }} />
      </Tooltip>
    </span>
  );
}

// Predefined tips for reuse
export const tips = {
  trackURL: '信件中按鈕的連結必須使用 {{.TrackURL}}，系統才能追蹤點擊行為',
  submitURL: '表單的 action 必須使用 {{.SubmitURL}}，系統才能記錄提交行為',
  templateVars: '可用變數：{{.FirstName}} {{.LastName}} {{.Email}} {{.Department}} {{.Position}} {{.TrackURL}}',
  educationHTML: '員工點擊釣魚連結並提交表單後看到的教育頁面，建議包含釣魚辨識技巧',
  captureFields: 'JSON 陣列格式，指定要記錄哪些表單欄位名稱（僅記錄欄位名，不記錄實際值）',
  smtpType: 'SMTP：直連郵件伺服器 / Mailgun：Mailgun API / SES：AWS SES 服務',
  selectionAll: '對所有收件人群組中的員工發送測試信',
  selectionDept: '僅對指定部門的員工發送',
  selectionRandom: '從收件人中隨機抽取指定比例進行測試',
  spreadSend: '將信件分散在數小時內發送，避免同一時間大量寄信被偵測或員工互相通風報信',
  csvFormat: '必填欄位：email。選填：first_name, last_name, department, gender, position',
  difficulty: '⭐ 簡單：明顯可疑 / ⭐⭐ 中等：需仔細辨識 / ⭐⭐⭐ 困難：高度擬真',
  scenario: '情境 = 信件模板 + Landing Page + 教育頁的完整組合',
  phishURL: '追蹤伺服器的網址，系統用來產生追蹤連結和 Landing Page 網址',
} as const;
