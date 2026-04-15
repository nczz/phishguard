import { Typography, Card, Steps, Table, Tag, Alert, Row, Col, Divider } from 'antd';
import { SettingOutlined, TeamOutlined, AppstoreOutlined, SendOutlined, BarChartOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';

const { Title, Paragraph, Text } = Typography;

const variables = [
  { var: '{{.FirstName}}', desc: '收件人的名字', example: '小明' },
  { var: '{{.LastName}}', desc: '收件人的姓氏', example: '王' },
  { var: '{{.Email}}', desc: '收件人的 Email', example: 'wang@company.com' },
  { var: '{{.Department}}', desc: '收件人的部門', example: '業務部' },
  { var: '{{.Position}}', desc: '收件人的職稱', example: '業務經理' },
  { var: '{{.TrackURL}}', desc: '追蹤連結（必填，放在信件按鈕的 href）', example: 'https://t.phishguard.tw/t/c/abc123' },
  { var: '{{.SubmitURL}}', desc: 'Landing Page 表單提交網址（放在 form action）', example: 'https://t.phishguard.tw/t/s/abc123' },
  { var: '{{.RID}}', desc: '收件人追蹤 ID（系統自動產生）', example: 'abc123-def456' },
];

export default function Guide() {
  const nav = useNavigate();

  return (
    <div style={{ maxWidth: 900, margin: '0 auto' }}>
      <Title level={3}>📖 使用指南</Title>
      <Paragraph type="secondary">依照以下步驟，5 分鐘內完成第一次釣魚模擬測試。</Paragraph>

      <Steps direction="vertical" current={-1} items={[
        {
          title: <Text strong>Step 1：設定發信機制</Text>,
          icon: <SettingOutlined />,
          description: (
            <Card size="small" style={{ marginTop: 8, marginBottom: 16 }}>
              <Paragraph>前往 <a onClick={() => nav('/app/settings/smtp')}>SMTP 設定</a> 新增發信設定。支援三種方式：</Paragraph>
              <Row gutter={16}>
                <Col span={8}><Card size="small" title="SMTP"><Text type="secondary">直連公司或第三方 SMTP 伺服器（如 Gmail SMTP）</Text></Card></Col>
                <Col span={8}><Card size="small" title="Mailgun"><Text type="secondary">填入 Mailgun Domain 和 API Key</Text></Card></Col>
                <Col span={8}><Card size="small" title="AWS SES"><Text type="secondary">填入 Region、Access Key、Secret Key</Text></Card></Col>
              </Row>
              <Alert type="info" title="建議" description="設定完成後請使用「測試發信」功能，確認信件能正常寄達。" showIcon style={{ marginTop: 12 }} />
            </Card>
          ),
        },
        {
          title: <Text strong>Step 2：匯入員工名單</Text>,
          icon: <TeamOutlined />,
          description: (
            <Card size="small" style={{ marginTop: 8, marginBottom: 16 }}>
              <Paragraph>前往 <a onClick={() => nav('/app/recipients')}>收件人管理</a>，下載 CSV 範本，填入員工資料後上傳。</Paragraph>
              <Paragraph>CSV 欄位：<Tag>email</Tag><Tag>first_name</Tag><Tag>last_name</Tag><Tag>department</Tag><Tag>gender</Tag><Tag>position</Tag></Paragraph>
              <Paragraph type="secondary">重複匯入同一 email 會自動更新資料，不會重複建立。</Paragraph>
            </Card>
          ),
        },
        {
          title: <Text strong>Step 3：瀏覽情境庫</Text>,
          icon: <AppstoreOutlined />,
          description: (
            <Card size="small" style={{ marginTop: 8, marginBottom: 16 }}>
              <Paragraph>前往 <a onClick={() => nav('/app/scenarios')}>情境庫</a> 查看可用的釣魚情境。每個情境包含：</Paragraph>
              <ul>
                <li><strong>信件模板</strong> — 釣魚信的內容（支援變數替換）</li>
                <li><strong>Landing Page</strong> — 收件人點擊連結後看到的頁面</li>
                <li><strong>教育頁</strong> — 員工「中招」後顯示的資安教育內容</li>
              </ul>
              <Paragraph type="secondary">系統已預建 5 個常見情境，也可以自行新增。</Paragraph>
            </Card>
          ),
        },
        {
          title: <Text strong>Step 4：建立釣魚測試</Text>,
          icon: <SendOutlined />,
          description: (
            <Card size="small" style={{ marginTop: 8, marginBottom: 16 }}>
              <Paragraph>點擊 <a onClick={() => nav('/app/campaigns/new')}>建立新測試</a>，三步完成：</Paragraph>
              <ol>
                <li><strong>選擇情境</strong> — 從情境庫選擇，或勾選「自動隨機」</li>
                <li><strong>選擇對象</strong> — 全公司 / 指定部門 / 隨機抽樣 N%</li>
                <li><strong>確認發送</strong> — 預覽信件內容，確認後發送</li>
              </ol>
              <Alert type="warning" title="提醒" description="建議勾選「分散發送」，避免同一時間大量寄信觸發郵件伺服器警報。" showIcon />
            </Card>
          ),
        },
        {
          title: <Text strong>Step 5：查看報表</Text>,
          icon: <BarChartOutlined />,
          description: (
            <Card size="small" style={{ marginTop: 8, marginBottom: 16 }}>
              <Paragraph>測試發送後，前往活動詳情頁查看即時結果：</Paragraph>
              <ul>
                <li><strong>釣魚漏斗</strong> — 寄達 → 開信 → 點擊 → 提交 → 舉報</li>
                <li><strong>部門風險排名</strong> — 哪個部門最容易中招</li>
                <li><strong>統計摘要</strong> — 各項比率一目瞭然</li>
              </ul>
              <Paragraph type="secondary">發送中的活動每 30 秒自動更新數據。</Paragraph>
            </Card>
          ),
        },
      ]} />

      <Divider />

      <Title level={4}>📝 模板變數參考</Title>
      <Paragraph type="secondary">在信件模板和 Landing Page 中可使用以下變數，系統會自動替換為收件人的實際資料：</Paragraph>
      <Table dataSource={variables} rowKey="var" size="small" pagination={false} columns={[
        { title: '變數', dataIndex: 'var', width: 180, render: (v: string) => <Tag color="blue" style={{ fontFamily: 'monospace' }}>{v}</Tag> },
        { title: '說明', dataIndex: 'desc' },
        { title: '範例值', dataIndex: 'example', render: (v: string) => <Text type="secondary">{v}</Text> },
      ]} />

      <Alert type="info" title="重要" style={{ marginTop: 16 }} showIcon
        description={<>信件模板中的按鈕連結必須使用 <Tag color="blue" style={{ fontFamily: 'monospace' }}>{'{{.TrackURL}}'}</Tag> 變數，系統才能追蹤點擊行為。Landing Page 的表單 action 必須使用 <Tag color="blue" style={{ fontFamily: 'monospace' }}>{'{{.SubmitURL}}'}</Tag>。</>}
      />

      <Divider />

      <Title level={4}>📊 追蹤指標說明</Title>
      <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
        {[
          { label: '寄達', color: '#1677ff', desc: '信件成功送達收件人信箱' },
          { label: '開信', color: '#13c2c2', desc: '收件人開啟信件（透過追蹤像素偵測，精度約 60-80%）' },
          { label: '點擊', color: '#faad14', desc: '收件人點擊信件中的連結（高可靠度）' },
          { label: '提交', color: '#ff4d4f', desc: '收件人在 Landing Page 填寫並提交表單（高可靠度）' },
          { label: '舉報', color: '#52c41a', desc: '收件人主動舉報可疑信件（正向指標）' },
        ].map(m => (
          <Col xs={24} sm={12} lg={8} key={m.label}>
            <Card size="small"><Tag color={m.color}>{m.label}</Tag><Paragraph type="secondary" style={{ marginTop: 8, marginBottom: 0 }}>{m.desc}</Paragraph></Card>
          </Col>
        ))}
      </Row>
    </div>
  );
}
