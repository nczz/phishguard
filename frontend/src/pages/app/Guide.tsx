import { Typography, Card, Table, Tag, Alert, Row, Col, Divider } from 'antd';
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
  { var: '{{.ReportURL}}', desc: '舉報連結（系統自動注入信件底部，也可手動放置）', example: 'https://t.phishguard.tw/t/r/abc123' },
  { var: '{{.RID}}', desc: '收件人追蹤 ID（系統自動產生）', example: 'abc123-def456' },
];

export default function Guide() {
  const nav = useNavigate();

  return (
    <div style={{ maxWidth: 900, margin: '0 auto' }}>
      <Title level={3}>📖 使用指南</Title>
      <Paragraph type="secondary">依照以下步驟，5 分鐘內完成第一次釣魚模擬測試。</Paragraph>

      {/* ── Step 1: SMTP ── */}
      <div id="smtp" style={{ scrollMarginTop: 80 }}>
        <Title level={4}>Step 1：設定發信機制</Title>
        <Card size="small" style={{ marginBottom: 24 }}>
          <Paragraph>前往 <a onClick={() => nav('/app/settings/smtp')}>SMTP 設定</a> 新增發信設定。支援三種方式：</Paragraph>
          <Row gutter={16}>
            <Col span={8}><Card size="small" title="SMTP"><Text type="secondary">直連公司或第三方 SMTP 伺服器。需填入 host、port、帳號密碼。適合已有郵件伺服器的企業。</Text></Card></Col>
            <Col span={8}><Card size="small" title="Mailgun"><Text type="secondary">填入 Mailgun Domain 和 API Key。適合沒有自建 SMTP 的團隊，到達率高。</Text></Card></Col>
            <Col span={8}><Card size="small" title="AWS SES"><Text type="secondary">填入 Region、Access Key、Secret Key。成本最低（約 NT$0.003/封），適合大量發送。</Text></Card></Col>
          </Row>
          <Alert type="info" title="建議" description="設定完成後請使用「測試發信」功能，確認信件能正常寄達再開始測試。" showIcon style={{ marginTop: 12 }} />
        </Card>
      </div>

      {/* ── Step 2: Recipients ── */}
      <div id="recipients" style={{ scrollMarginTop: 80 }}>
        <Title level={4}>Step 2：匯入員工名單</Title>
        <Card size="small" style={{ marginBottom: 24 }}>
          <Paragraph>前往 <a onClick={() => nav('/app/recipients')}>收件人管理</a>，下載 CSV 範本，填入員工資料後上傳。</Paragraph>
          <Paragraph><strong>CSV 欄位：</strong></Paragraph>
          <Table size="small" pagination={false} dataSource={[
            { field: 'email', required: '必填', desc: '員工 Email 地址（作為唯一識別）' },
            { field: 'first_name', required: '選填', desc: '名字' },
            { field: 'last_name', required: '選填', desc: '姓氏' },
            { field: 'department', required: '建議填', desc: '部門名稱（用於部門篩選和報表分析）' },
            { field: 'gender', required: '選填', desc: '性別：男 / 女 / 不指定' },
            { field: 'position', required: '選填', desc: '職稱' },
          ]} rowKey="field" columns={[
            { title: '欄位', dataIndex: 'field', width: 120, render: (v: string) => <Tag style={{ fontFamily: 'monospace' }}>{v}</Tag> },
            { title: '必填', dataIndex: 'required', width: 80 },
            { title: '說明', dataIndex: 'desc' },
          ]} />
          <Alert type="success" title="去重機制" description="重複匯入同一 email 會自動更新該員工的資料（姓名、部門等），不會重複建立。" showIcon style={{ marginTop: 12 }} />
        </Card>
      </div>

      {/* ── Step 3: Scenarios ── */}
      <div id="scenarios" style={{ scrollMarginTop: 80 }}>
        <Title level={4}>Step 3：瀏覽情境庫</Title>
        <Card size="small" style={{ marginBottom: 24 }}>
          <Paragraph>前往 <a onClick={() => nav('/app/scenarios')}>情境庫</a> 查看可用的釣魚情境。</Paragraph>
          <Paragraph><strong>每個情境包含三個部分：</strong></Paragraph>
          <Row gutter={16}>
            <Col span={8}><Card size="small" title="✉️ 信件模板"><Text type="secondary">釣魚信的 HTML 內容，支援變數替換（見下方變數表）。按鈕連結必須使用 {'{{.TrackURL}}'} 變數。</Text></Card></Col>
            <Col span={8}><Card size="small" title="🖥️ Landing Page"><Text type="secondary">收件人點擊連結後看到的頁面。表單 action 必須使用 {'{{.SubmitURL}}'} 變數。</Text></Card></Col>
            <Col span={8}><Card size="small" title="📚 教育頁"><Text type="secondary">員工提交表單後顯示的資安教育內容，說明這是測試並教導辨識技巧。</Text></Card></Col>
          </Row>
          <Paragraph type="secondary" style={{ marginTop: 12 }}>
            <strong>難度分級：</strong>⭐ 簡單（明顯可疑）/ ⭐⭐ 中等（需仔細辨識）/ ⭐⭐⭐ 困難（高度擬真）。建議首次測試用簡單難度建立 baseline。
          </Paragraph>
        </Card>
      </div>

      {/* ── Step 4: Campaign ── */}
      <div id="campaign" style={{ scrollMarginTop: 80 }}>
        <Title level={4}>Step 4：建立釣魚測試</Title>
        <Card size="small" style={{ marginBottom: 24 }}>
          <Paragraph>點擊 <a onClick={() => nav('/app/campaigns/new')}>建立新測試</a>，三步完成：</Paragraph>
          <Row gutter={16}>
            <Col span={8}>
              <Card size="small" title="① 選擇情境">
                <Text type="secondary">從情境庫選擇一個情境，或勾選「自動隨機」讓系統隨機選擇。</Text>
              </Card>
            </Col>
            <Col span={8}>
              <Card size="small" title="② 選擇對象">
                <ul style={{ paddingLeft: 16, margin: 0 }}>
                  <li><strong>全公司</strong> — 所有員工</li>
                  <li><strong>指定部門</strong> — 勾選特定部門</li>
                  <li><strong>隨機抽樣</strong> — 隨機抽取 N%</li>
                </ul>
              </Card>
            </Col>
            <Col span={8}>
              <Card size="small" title="③ 確認發送">
                <Text type="secondary">確認摘要、預覽信件內容，按下發送。建議勾選「分散發送」。</Text>
              </Card>
            </Col>
          </Row>
          <Alert type="warning" title="分散發送" style={{ marginTop: 12 }} showIcon
            description="勾選後信件會在數小時內隨機分散發送，避免同一時間大量寄信觸發郵件伺服器警報，也防止員工互相通風報信。" />
        </Card>
      </div>

      {/* ── Step 5: Reports ── */}
      <div id="reports" style={{ scrollMarginTop: 80 }}>
        <Title level={4}>Step 5：查看報表</Title>
        <Card size="small" style={{ marginBottom: 24 }}>
          <Paragraph>測試發送後，前往活動詳情頁查看即時結果：</Paragraph>
          <ul>
            <li><strong>釣魚漏斗</strong> — 寄達 → 開信 → 點擊 → 提交 → 舉報，逐層轉換率</li>
            <li><strong>部門風險排名</strong> — 哪個部門點擊率最高，需要加強訓練</li>
            <li><strong>統計摘要</strong> — 各項比率數字一目瞭然</li>
          </ul>
          <Paragraph type="secondary">發送中的活動每 30 秒自動更新數據。</Paragraph>
        </Card>
      </div>

      <Divider />

      {/* ── Variables ── */}
      <div id="variables" style={{ scrollMarginTop: 80 }}>
        <Title level={4}>📝 模板變數參考</Title>
        <Paragraph type="secondary">在信件模板和 Landing Page 中可使用以下變數，系統會自動替換為收件人的實際資料：</Paragraph>
        <Table dataSource={variables} rowKey="var" size="small" pagination={false} columns={[
          { title: '變數', dataIndex: 'var', width: 180, render: (v: string) => <Tag color="blue" style={{ fontFamily: 'monospace' }}>{v}</Tag> },
          { title: '說明', dataIndex: 'desc' },
          { title: '範例值', dataIndex: 'example', render: (v: string) => <Text type="secondary">{v}</Text> },
        ]} />
        <Alert type="info" title="重要" style={{ marginTop: 16 }} showIcon
          description={<>信件模板中的按鈕連結必須使用 <Tag color="blue" style={{ fontFamily: 'monospace' }}>{'{{.TrackURL}}'}</Tag>，系統才能追蹤點擊行為。Landing Page 的表單 action 必須使用 <Tag color="blue" style={{ fontFamily: 'monospace' }}>{'{{.SubmitURL}}'}</Tag>。</>}
        />
      </div>

      <Divider />

      {/* ── Metrics ── */}
      <div id="metrics" style={{ scrollMarginTop: 80 }}>
        <Title level={4}>📊 追蹤指標說明</Title>
        <Paragraph type="secondary">系統追蹤以下五個指標，構成完整的釣魚漏斗：</Paragraph>
        <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
          {[
            { label: '寄達', color: '#1677ff', desc: '信件成功送達收件人信箱。失敗原因可能是 email 不存在或被 mail gateway 擋下。' },
            { label: '開信', color: '#13c2c2', desc: '收件人開啟信件。透過 1x1 追蹤像素偵測，精度約 60-80%（部分郵件客戶端會預載或封鎖圖片）。' },
            { label: '點擊', color: '#faad14', desc: '收件人點擊信件中的連結。高可靠度指標。注意：部分 mail gateway（如 Microsoft Safe Links）會預先點擊連結，系統會過濾 bot 請求。' },
            { label: '提交', color: '#ff4d4f', desc: '收件人在 Landing Page 填寫並提交表單。高可靠度指標。系統僅記錄欄位名稱，不記錄實際輸入值（如密碼）。' },
            { label: '舉報', color: '#52c41a', desc: '收件人主動舉報可疑信件。這是正向指標，舉報率越高代表員工資安意識越好。' },
          ].map(m => (
            <Col xs={24} sm={12} key={m.label}>
              <Card size="small"><Tag color={m.color} style={{ fontSize: 14 }}>{m.label}</Tag><Paragraph type="secondary" style={{ marginTop: 8, marginBottom: 0 }}>{m.desc}</Paragraph></Card>
            </Col>
          ))}
        </Row>
      </div>

      <Divider />

      {/* ── Auto Test ── */}
      <div id="autotest" style={{ scrollMarginTop: 80 }}>
        <Title level={4}>🔄 自動定期測試</Title>
        <Card size="small" style={{ marginBottom: 24 }}>
          <Paragraph>前往 <a onClick={() => nav('/app/settings/auto-test')}>設定 → 自動測試</a> 開啟自動排程。</Paragraph>
          <ul>
            <li><strong>頻率</strong>：每月 / 每季 / 每半年</li>
            <li><strong>對象</strong>：全公司或隨機抽樣 N%</li>
            <li><strong>情境</strong>：系統自動從情境庫隨機選擇</li>
            <li><strong>冷卻期</strong>：同一人 30 天內不會被重複測試</li>
            <li><strong>報表</strong>：測試完成後自動寄送報表給租戶管理員</li>
          </ul>
        </Card>
      </div>

      <Divider />

      {/* ── Reports ── */}
      <div id="reports-detail" style={{ scrollMarginTop: 80 }}>
        <Title level={4}>📈 報表功能</Title>
        <Card size="small" style={{ marginBottom: 24 }}>
          <Row gutter={16}>
            <Col span={8}><Card size="small" title="匯出 PDF"><Text type="secondary">活動詳情頁 → 匯出 PDF，包含漏斗統計和部門排名，可直接交給管理層。</Text></Card></Col>
            <Col span={8}><Card size="small" title="匯出 CSV"><Text type="secondary">活動詳情頁 → 匯出 CSV，包含每位收件人的完整狀態和時間戳。</Text></Card></Col>
            <Col span={8}><Card size="small" title="寄送報表"><Text type="secondary">活動完成後，點擊「寄送報表」或系統自動寄送給租戶管理員。</Text></Card></Col>
          </Row>
          <Paragraph style={{ marginTop: 12 }}>
            <strong>進階報表：</strong>側邊欄 → 報表 → <a onClick={() => nav('/app/reports/offenders')}>累犯追蹤</a>（跨活動個人歷史）和 <a onClick={() => nav('/app/reports/trend')}>趨勢分析</a>（折線圖對比各次活動指標變化）。
          </Paragraph>
        </Card>
      </div>

      <Divider />

      {/* ── FAQ ── */}
      <div id="faq" style={{ scrollMarginTop: 80 }}>
        <Title level={4}>❓ 常見問題</Title>
        <Card size="small">
          <Paragraph><strong>Q：信件寄出去都進垃圾郵件怎麼辦？</strong></Paragraph>
          <Paragraph type="secondary">請確認 SMTP 設定的寄件地址有正確的 SPF/DKIM/DMARC 記錄。如果使用公司內部測試，建議請 IT 部門將發信 IP 或域名加入 mail gateway 白名單。</Paragraph>
          <Divider style={{ margin: '12px 0' }} />
          <Paragraph><strong>Q：開信率為什麼不準確？</strong></Paragraph>
          <Paragraph type="secondary">開信追蹤依賴載入追蹤像素。Apple Mail Privacy Protection、Gmail Image Proxy 等會影響準確度。建議將開信率視為「下限估計」。</Paragraph>
          <Divider style={{ margin: '12px 0' }} />
          <Paragraph><strong>Q：可以對同一批人重複測試嗎？</strong></Paragraph>
          <Paragraph type="secondary">可以。每次建立新測試都會產生新的追蹤 ID，不會與之前的測試混淆。</Paragraph>
          <Divider style={{ margin: '12px 0' }} />
          <Paragraph><strong>Q：Landing Page 的 {'{{.SubmitURL}}'} 是什麼？</strong></Paragraph>
          <Paragraph type="secondary">這是系統自動產生的表單提交網址。將它放在 {'<form action="{{.SubmitURL}}">'}，系統就能記錄收件人的提交行為並顯示教育頁。</Paragraph>
        </Card>
      </div>
    </div>
  );
}
