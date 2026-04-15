import { Typography, Card, Row, Col, Divider, Tag } from 'antd';
import {
  SettingOutlined, TeamOutlined, AppstoreOutlined, SendOutlined,
  MailOutlined, GlobalOutlined, FormOutlined, BarChartOutlined,
  CheckCircleOutlined, ArrowDownOutlined,
} from '@ant-design/icons';

const { Title, Paragraph, Text } = Typography;

const stepStyle = { textAlign: 'center' as const, padding: '16px 12px' };
const iconStyle = (bg: string) => ({ fontSize: 32, color: '#fff', background: bg, borderRadius: '50%', padding: 16, display: 'inline-block', marginBottom: 8 });
const arrowDown = <div style={{ textAlign: 'center', padding: '4px 0' }}><ArrowDownOutlined style={{ fontSize: 20, color: '#bbb' }} /></div>;

export default function FlowDiagram() {
  return (
    <div style={{ maxWidth: 960, margin: '0 auto' }}>
      <Title level={3}>🗺️ 系統流程總覽</Title>
      <Paragraph type="secondary">從設定到報表，一張圖看懂 PhishGuard 的完整運作流程。</Paragraph>

      {/* ── Phase 1: Setup ── */}
      <Title level={4} style={{ marginTop: 32 }}>📋 前置設定（一次性）</Title>
      <Row gutter={16}>
        <Col xs={24} sm={8}>
          <Card hoverable style={stepStyle}>
            <div style={iconStyle('#1677ff')}><SettingOutlined /></div>
            <Title level={5}>① 設定 SMTP</Title>
            <Text type="secondary">選擇 SMTP / Mailgun / SES<br />設定寄件地址與認證<br />測試發信確認可用</Text>
          </Card>
        </Col>
        <Col xs={24} sm={8}>
          <Card hoverable style={stepStyle}>
            <div style={iconStyle('#52c41a')}><TeamOutlined /></div>
            <Title level={5}>② 匯入員工</Title>
            <Text type="secondary">下載 CSV 範本<br />填入 email、姓名、部門<br />上傳後自動解析分組</Text>
          </Card>
        </Col>
        <Col xs={24} sm={8}>
          <Card hoverable style={stepStyle}>
            <div style={iconStyle('#722ed1')}><AppstoreOutlined /></div>
            <Title level={5}>③ 確認情境</Title>
            <Text type="secondary">系統已預建 5 個情境<br />或自建信件模板 + Landing Page<br />+ 教育頁組合</Text>
          </Card>
        </Col>
      </Row>

      {arrowDown}

      {/* ── Phase 2: Campaign ── */}
      <Title level={4}>🚀 建立與發送測試</Title>
      <Card style={{ background: '#f6ffed', border: '1px solid #b7eb8f' }}>
        <Row gutter={16} align="middle">
          <Col xs={24} sm={6} style={{ textAlign: 'center' }}>
            <div style={iconStyle('#faad14')}><SendOutlined /></div>
            <Title level={5}>建立測試</Title>
          </Col>
          <Col xs={24} sm={18}>
            <Row gutter={[8, 8]}>
              <Col span={8}>
                <Card size="small" title="Step 1" style={{ height: '100%' }}>
                  <Text strong>選擇情境</Text><br />
                  <Text type="secondary">從情境庫選擇<br />或勾選「自動隨機」</Text>
                </Card>
              </Col>
              <Col span={8}>
                <Card size="small" title="Step 2" style={{ height: '100%' }}>
                  <Text strong>選擇對象</Text><br />
                  <Text type="secondary">全公司 / 指定部門<br />/ 隨機抽樣 N%</Text>
                </Card>
              </Col>
              <Col span={8}>
                <Card size="small" title="Step 3" style={{ height: '100%' }}>
                  <Text strong>確認發送</Text><br />
                  <Text type="secondary">預覽信件 → 發送<br />可選分散發送</Text>
                </Card>
              </Col>
            </Row>
          </Col>
        </Row>
      </Card>

      {arrowDown}

      {/* ── Phase 3: Tracking ── */}
      <Title level={4}>📡 追蹤階段（系統自動）</Title>
      <Paragraph type="secondary">信件發出後，系統自動追蹤收件人的每一步行為：</Paragraph>

      <div style={{ position: 'relative', padding: '0 24px' }}>
        {/* Funnel visualization */}
        {[
          { icon: <MailOutlined />, label: '寄達', desc: 'Worker 透過 SMTP/Mailgun/SES 發送信件', color: '#1677ff', width: '100%' },
          { icon: <MailOutlined />, label: '開信', desc: '收件人開啟信件 → 載入追蹤像素 → 記錄', color: '#13c2c2', width: '85%' },
          { icon: <GlobalOutlined />, label: '點擊連結', desc: '收件人點擊 {{.TrackURL}} → 記錄 → 導向 Landing Page', color: '#faad14', width: '65%' },
          { icon: <FormOutlined />, label: '提交表單', desc: '收件人在 Landing Page 填寫表單 → 記錄欄位名 → 顯示教育頁', color: '#ff4d4f', width: '45%' },
          { icon: <CheckCircleOutlined />, label: '舉報', desc: '收件人主動舉報可疑信件（正向指標 ✓）', color: '#52c41a', width: '30%' },
        ].map((step, i) => (
          <div key={i} style={{ marginBottom: 8 }}>
            <div style={{
              background: step.color, color: '#fff', borderRadius: 8, padding: '12px 16px',
              width: step.width, margin: '0 auto', display: 'flex', alignItems: 'center', gap: 12,
              transition: 'all 0.3s',
            }}>
              <span style={{ fontSize: 20 }}>{step.icon}</span>
              <div>
                <Text strong style={{ color: '#fff', fontSize: 15 }}>{step.label}</Text>
                <br />
                <Text style={{ color: 'rgba(255,255,255,0.85)', fontSize: 12 }}>{step.desc}</Text>
              </div>
            </div>
          </div>
        ))}
      </div>

      {arrowDown}

      {/* ── Phase 4: Results ── */}
      <Title level={4}>📊 報表與分析</Title>
      <Row gutter={16}>
        <Col xs={24} sm={8}>
          <Card hoverable style={stepStyle}>
            <div style={iconStyle('#1677ff')}><BarChartOutlined /></div>
            <Title level={5}>釣魚漏斗</Title>
            <Text type="secondary">寄達 → 開信 → 點擊<br />→ 提交 → 舉報<br />逐層轉換率圖表</Text>
          </Card>
        </Col>
        <Col xs={24} sm={8}>
          <Card hoverable style={stepStyle}>
            <div style={iconStyle('#faad14')}><TeamOutlined /></div>
            <Title level={5}>部門風險排名</Title>
            <Text type="secondary">各部門點擊率排序<br />識別高風險部門<br />針對性加強訓練</Text>
          </Card>
        </Col>
        <Col xs={24} sm={8}>
          <Card hoverable style={stepStyle}>
            <div style={iconStyle('#52c41a')}><CheckCircleOutlined /></div>
            <Title level={5}>統計摘要</Title>
            <Text type="secondary">開信率 / 點擊率 / 提交率<br />舉報率<br />匯出 PDF 報告</Text>
          </Card>
        </Col>
      </Row>

      <Divider />

      {/* ── Architecture ── */}
      <Title level={4}>⚙️ 系統架構</Title>
      <Paragraph type="secondary">PhishGuard 由四個核心元件組成：</Paragraph>
      <Row gutter={16}>
        {[
          { title: 'API Server', port: ':8080', desc: '處理所有管理操作：登入、建立測試、查看報表。前端 UI 透過此 API 操作。', color: '#1677ff' },
          { title: 'Track Server', port: ':8090', desc: '處理收件人行為追蹤：開信像素、連結點擊、表單提交。極輕量、低延遲。', color: '#13c2c2' },
          { title: 'Mail Worker', port: '背景', desc: '消費發信佇列，透過 SMTP/Mailgun/SES 發送信件。可水平擴展。', color: '#faad14' },
          { title: 'MySQL + Redis', port: '資料層', desc: 'MySQL 儲存所有業務資料（租戶隔離）。Redis 作為發信佇列。', color: '#722ed1' },
        ].map(c => (
          <Col xs={24} sm={12} lg={6} key={c.title}>
            <Card size="small" style={{ borderTop: `3px solid ${c.color}` }}>
              <Text strong>{c.title}</Text> <Tag>{c.port}</Tag>
              <Paragraph type="secondary" style={{ marginTop: 8, marginBottom: 0, fontSize: 12 }}>{c.desc}</Paragraph>
            </Card>
          </Col>
        ))}
      </Row>

      <Divider />

      {/* ── Data Flow ── */}
      <Title level={4}>🔄 資料流向</Title>
      <Card>
        <pre style={{ margin: 0, fontSize: 13, lineHeight: 1.8, overflow: 'auto' }}>{`
  客戶操作                    系統處理                      收件人行為
  ────────                    ────────                      ────────
  選情境 + 選對象              │                              │
       │                      ▼                              │
       └──→ 建立 Campaign ──→ 產生每人唯一 rid               │
                               │                              │
                               ▼                              │
                          渲染信件模板                         │
                          ├ 替換 {{.FirstName}} 等變數         │
                          ├ 嵌入追蹤像素 (1x1 GIF)            │
                          └ 改寫連結為追蹤 URL                 │
                               │                              │
                               ▼                              │
                          推入發信佇列 ──→ Worker 發信         │
                                                              │
                                              收件人開信 ──→ 載入像素 ──→ Track Server 記錄 opened
                                              收件人點連結 ──→ Track Server 記錄 clicked ──→ 導向 Landing Page
                                              收件人提交表單 ──→ Track Server 記錄 submitted ──→ 顯示教育頁
                                              收件人舉報 ──→ Track Server 記錄 reported ✓
                               │
                               ▼
                          報表即時更新
                          ├ 釣魚漏斗
                          ├ 部門風險排名
                          └ 統計摘要
`}</pre>
      </Card>
    </div>
  );
}
