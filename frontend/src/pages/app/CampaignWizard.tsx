import { useEffect, useMemo, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Steps, Card, Row, Col, Button, Checkbox, Radio, Slider, Input, DatePicker,
  Tag, Modal, Spin, Typography, message, Space, Alert, Descriptions, Tooltip, Table,
} from 'antd';
import { RocketOutlined, EyeOutlined, QuestionCircleOutlined } from '@ant-design/icons';
import dayjs from 'dayjs';
import { tips } from '../../components/FieldHelp';
import type { Scenario, RecipientGroup, SMTPProfile, Campaign } from '../../api/client';
import { api, getErrorMessage } from '../../api/client';

const { Title, Text } = Typography;

type SelectionMode = 'all' | 'department' | 'sample' | 'individual';

const CATEGORY_ICON: Record<string, string> = {
  password: '🔐',
  credential: '🔑',
  package: '📦',
  invoice: '💳',
  email: '📧',
};

function categoryIcon(cat: string) {
  return CATEGORY_ICON[cat] ?? '📧';
}

function difficultyStars(d: number | string) {
  const n = parseInt(String(d), 10) || 1;
  return '⭐'.repeat(n);
}

function defaultCampaignName() {
  const d = new Date();
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')} 釣魚測試`;
}

export default function CampaignWizard() {
  const navigate = useNavigate();
  const [step, setStep] = useState(0);
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [previewOpen, setPreviewOpen] = useState(false);

  // Data
  const [scenarios, setScenarios] = useState<Scenario[]>([]);
  const [groups, setGroups] = useState<RecipientGroup[]>([]);
  const [smtpProfiles, setSmtpProfiles] = useState<SMTPProfile[]>([]);

  // Wizard state
  const [autoRandom, setAutoRandom] = useState(false);
  const [selectedScenario, setSelectedScenario] = useState<string | null>(null);
  const [selectionMode, setSelectionMode] = useState<SelectionMode>('all');
  const [samplePercent, setSamplePercent] = useState(30);
  const [selectedGroups, setSelectedGroups] = useState<string[]>([]);
  const [departments, setDepartments] = useState<string[]>([]);
  const [sendMode, setSendMode] = useState<'immediate' | 'scheduled'>('immediate');
  const [scheduleStart, setScheduleStart] = useState<string>('');
  const [scheduleEnd, setScheduleEnd] = useState<string>('');
  const [workingHoursOnly, setWorkingHoursOnly] = useState(true);
  const [skipWeekends, setSkipWeekends] = useState(true);
  const [campaignName, setCampaignName] = useState(defaultCampaignName);
  const [selectedRecipientIds, setSelectedRecipientIds] = useState<string[]>([]);
  const [recipientSearch, setRecipientSearch] = useState('');
  const [skipCooldown, setSkipCooldown] = useState(false);

  useEffect(() => {
    Promise.all([
      api.get<Scenario[]>('/scenarios'),
      api.get<RecipientGroup[]>('/recipient-groups'),
      api.get<SMTPProfile[]>('/smtp-profiles'),
    ]).then(([s, g, p]) => {
      setScenarios(s);
      setGroups(g);
      setSmtpProfiles(p);
      // Auto-select all groups
      setSelectedGroups(g.map((grp: RecipientGroup) => String(grp.id)));
    }).finally(() => setLoading(false));
  }, []);

  // Derived
  const allRecipients = useMemo(
    () => groups.filter(g => selectedGroups.includes(String(g.id))).flatMap(g => g.recipients ?? []),
    [groups, selectedGroups],
  );

  const allDepartments = useMemo(
    () => [...new Set(allRecipients.map(r => r.department).filter(Boolean))].sort(),
    [allRecipients],
  );

  const estimatedCount = useMemo(() => {
    if (selectionMode === 'individual') return selectedRecipientIds.length;
    let pool = allRecipients;
    if (selectionMode === 'department') {
      pool = pool.filter(r => departments.includes(r.department));
    }
    if (selectionMode === 'sample') {
      return Math.max(1, Math.round(pool.length * samplePercent / 100));
    }
    return pool.length;
  }, [allRecipients, selectionMode, departments, samplePercent, selectedRecipientIds]);

  const scenarioObj = scenarios.find(s => s.id === selectedScenario);
  const smtpProfile = smtpProfiles[0] ?? null;
  const trimmedCampaignName = campaignName.trim();
  const scheduleInvalid = sendMode === 'scheduled' && (
    !scheduleStart || !scheduleEnd || !dayjs(scheduleEnd).isAfter(dayjs(scheduleStart))
  );

  const canNext = (s: number) => {
    if (s === 0) return (autoRandom && scenarios.length > 0) || !!selectedScenario;
    if (s === 1) {
      return selectedGroups.length > 0
        && estimatedCount > 0
        && !scheduleInvalid
        && (selectionMode === 'individual' ? selectedRecipientIds.length > 0 : selectionMode !== 'department' || departments.length > 0);
    }
    return true;
  };

  const canSubmit = !!smtpProfile && trimmedCampaignName.length > 0 && estimatedCount > 0 && !scheduleInvalid;

  async function handleSubmit() {
    if (!canSubmit) return;
    setSubmitting(true);
    try {
      const campaign = await api.post<Campaign>('/campaigns', {
        name: trimmedCampaignName,
        scenario_id: autoRandom ? null : selectedScenario,
        smtp_profile_id: smtpProfile?.id,
        group_ids: selectedGroups.map(Number),
        phish_url: window.location.origin + '/phish',
        selection_mode: selectionMode,
        sample_percent: selectionMode === 'sample' ? samplePercent : undefined,
        departments: selectionMode === 'department' ? departments : undefined,
        send_mode: sendMode,
        schedule_start: sendMode === 'scheduled' ? scheduleStart : undefined,
        schedule_end: sendMode === 'scheduled' ? scheduleEnd : undefined,
        working_hours_only: workingHoursOnly,
        skip_weekends: skipWeekends,
        timezone: Intl.DateTimeFormat().resolvedOptions().timeZone,
      });
      await api.post('/campaigns/' + campaign.id + '/launch', {
        skip_cooldown: skipCooldown,
        recipient_ids: selectionMode === 'individual' ? selectedRecipientIds.map(Number) : undefined,
      });
      message.success('測試已發送！');
      navigate('/app/campaigns/' + campaign.id);
    } catch (err) {
      message.error(getErrorMessage(err));
    } finally {
      setSubmitting(false);
    }
  }

  if (loading) return <Spin style={{ display: 'block', margin: '20vh auto' }} size="large" />;

  return (
    <div style={{ maxWidth: 900, margin: '0 auto' }}>
      <Steps current={step} style={{ marginBottom: 32 }} items={[
        { title: '選擇情境' },
        { title: '選擇對象' },
        { title: '確認發送' },
      ]} />

      {/* ── Step 1: 選擇情境 ── */}
      {step === 0 && (
        <>
          <Title level={4}>選擇測試情境</Title>
          <Checkbox
            checked={autoRandom}
            onChange={e => { setAutoRandom(e.target.checked); setSelectedScenario(null); }}
            style={{ marginBottom: 16 }}
          >
            ☑ 自動隨機（系統從情境庫隨機選擇）
          </Checkbox>

          {!autoRandom && (
            <Row gutter={[16, 16]}>
              {scenarios.map(s => (
                <Col xs={24} sm={12} md={8} key={s.id}>
                  <Card
                    hoverable
                    onClick={() => setSelectedScenario(s.id)}
                    style={{
                      borderColor: selectedScenario === s.id ? '#1677ff' : undefined,
                      borderWidth: selectedScenario === s.id ? 2 : 1,
                    }}
                  >
                    <div style={{ fontSize: 28, marginBottom: 8 }}>{categoryIcon(s.category)}</div>
                    <Text strong>{s.name}</Text>
                    <div style={{ marginTop: 8 }}>
                      <span>{difficultyStars(s.difficulty)}</span>
                      <Tag style={{ marginLeft: 8 }}>{s.language}</Tag>
                    </div>
                  </Card>
                </Col>
              ))}
            </Row>
          )}
          {autoRandom && scenarios.length === 0 && (
            <Alert type="warning" showIcon message="目前沒有可用情境，請先建立或匯入情境。" />
          )}
        </>
      )}

      {/* ── Step 2: 選擇對象 ── */}
      {step === 1 && (
        <>
          <Title level={4}>選擇測試對象</Title>

          <Text strong style={{ display: 'block', marginBottom: 8 }}>測試範圍（目前共 {allRecipients.length} 位收件人）</Text>
          <Radio.Group
            value={selectionMode}
            onChange={e => setSelectionMode(e.target.value)}
            style={{ marginBottom: 16 }}
          >
            <Space direction="vertical">
              <Radio value="all">全部發送（{allRecipients.length} 人）<Tooltip title={tips.selectionAll}><QuestionCircleOutlined style={{color:'#999'}} /></Tooltip></Radio>
              <Radio value="department">依部門篩選 <Tooltip title={tips.selectionDept}><QuestionCircleOutlined style={{color:'#999'}} /></Tooltip></Radio>
              <Radio value="sample">隨機抽樣 <Tooltip title={tips.selectionRandom}><QuestionCircleOutlined style={{color:'#999'}} /></Tooltip></Radio>
              <Radio value="individual">手動挑選個別人員</Radio>
            </Space>
          </Radio.Group>

          <Checkbox
            checked={skipCooldown}
            onChange={e => setSkipCooldown(e.target.checked)}
            style={{ marginBottom: 16, marginLeft: 24 }}
          >
            忽略 30 天冷卻期（強制發送給所有選中對象）
          </Checkbox>

          {selectionMode === 'department' && (
            allDepartments.length > 0 ? (
              <Checkbox.Group
                value={departments}
                onChange={v => setDepartments(v as string[])}
                style={{ display: 'block', marginBottom: 16, paddingLeft: 24 }}
              >
                <Space direction="vertical">
                  {allDepartments.map(d => (
                    <Checkbox key={d} value={d}>{d}</Checkbox>
                  ))}
                </Space>
              </Checkbox.Group>
            ) : (
              <Alert type="warning" showIcon message="所選群組沒有部門資料，請改用全部發送或手動挑選。" style={{ marginBottom: 16 }} />
            )
          )}

          {selectionMode === 'sample' && (
            <div style={{ paddingLeft: 24, marginBottom: 16, maxWidth: 400 }}>
              <Slider
                min={10} max={100} step={5}
                value={samplePercent}
                onChange={setSamplePercent}
                marks={{ 10: '10%', 50: '50%', 100: '100%' }}
              />
              <Text type="secondary">預估 {estimatedCount} 人</Text>
            </div>
          )}

          {selectionMode === 'department' && departments.length > 0 && (
            <Text type="secondary" style={{ display: 'block', marginBottom: 16 }}>
              預估 {estimatedCount} 人
            </Text>
          )}

          {selectionMode === 'individual' && (
            <div style={{ marginBottom: 16 }}>
              <Input
                placeholder="搜尋 Email、姓名或部門"
                allowClear
                value={recipientSearch}
                onChange={e => setRecipientSearch(e.target.value)}
                style={{ marginBottom: 8, maxWidth: 300 }}
              />
              <Table
                size="small"
                rowKey="id"
                dataSource={allRecipients.filter(r => {
                  if (!recipientSearch) return true;
                  const q = recipientSearch.toLowerCase();
                  return r.email.toLowerCase().includes(q) || (r.last_name + r.first_name).toLowerCase().includes(q) || r.department.toLowerCase().includes(q);
                })}
                scroll={{ x: 500 }}
                pagination={{ pageSize: 10, showSizeChanger: true, size: 'small' }}
                rowSelection={{
                  selectedRowKeys: selectedRecipientIds,
                  onChange: (keys) => setSelectedRecipientIds(keys as string[]),
                }}
                columns={[
                  { title: 'Email', dataIndex: 'email', width: 200, ellipsis: { showTitle: false }, render: (v: string) => <Tooltip title={v}>{v}</Tooltip> },
                  { title: '姓名', key: 'name', width: 100, render: (_: unknown, r: { last_name: string; first_name: string }) => r.last_name + r.first_name },
                  { title: '部門', dataIndex: 'department', width: 100 },
                ]}
              />
              <Text type="secondary">已選 {selectedRecipientIds.length} 人</Text>
            </div>
          )}

          <Card size="small" title="發送排程" style={{ marginTop: 16 }}>
            <Radio.Group value={sendMode} onChange={e => setSendMode(e.target.value)} style={{ marginBottom: 12 }}>
              <Radio value="immediate">立即發送（系統自動控制速率）</Radio>
              <Radio value="scheduled">排程發送（指定時間窗口）</Radio>
            </Radio.Group>

            {sendMode === 'scheduled' && (
              <div style={{ marginLeft: 24 }}>
                <div style={{ marginBottom: 8 }}>
                  <Text>開始時間：</Text>
                  <DatePicker showTime format="YYYY-MM-DD HH:mm" placeholder="選擇開始時間"
                    value={scheduleStart ? dayjs(scheduleStart) : null}
                    onChange={v => setScheduleStart(v ? v.toISOString() : '')}
                    style={{ marginLeft: 8 }} />
                </div>
                <div style={{ marginBottom: 8 }}>
                  <Text>結束時間：</Text>
                  <DatePicker showTime format="YYYY-MM-DD HH:mm" placeholder="選擇結束時間"
                    value={scheduleEnd ? dayjs(scheduleEnd) : null}
                    onChange={v => setScheduleEnd(v ? v.toISOString() : '')}
                    style={{ marginLeft: 8 }} />
                </div>
                {scheduleInvalid && (
                  <Alert
                    type="warning"
                    showIcon
                    message="請選擇有效的排程時間，結束時間必須晚於開始時間。"
                    style={{ marginTop: 8, marginBottom: 8 }}
                  />
                )}
              </div>
            )}

            <div style={{ marginTop: 8 }}>
              <Checkbox checked={workingHoursOnly} onChange={e => setWorkingHoursOnly(e.target.checked)}>
                僅工作時間發送（09:00-17:00）
              </Checkbox>
            </div>
            <div>
              <Checkbox checked={skipWeekends} onChange={e => setSkipWeekends(e.target.checked)}>
                避開週末（週六、週日不發送）
              </Checkbox>
            </div>
            <Text type="secondary" style={{ display: 'block', marginTop: 8, fontSize: 12 }}>
              💡 發信速率由系統根據發信服務商（SES/Mailgun/SMTP）的政策自動控制，無需手動設定。
            </Text>
          </Card>
        </>
      )}

      {/* ── Step 3: 確認發送 ── */}
      {step === 2 && (
        <>
          <Title level={4}>確認測試內容</Title>

          <Card style={{ marginBottom: 24 }}>
            <Descriptions column={1} bordered size="small">
              <Descriptions.Item label="情境">
                {autoRandom
                  ? '自動隨機'
                  : scenarioObj
                    ? `${scenarioObj.name} ${difficultyStars(scenarioObj.difficulty)}`
                    : '—'}
              </Descriptions.Item>
              <Descriptions.Item label="對象">
                {selectionMode === 'all' && `全公司 ${estimatedCount} 人`}
                {selectionMode === 'department' && `指定部門（${departments.join('、')}）${estimatedCount} 人`}
                {selectionMode === 'sample' && `隨機抽樣 ${samplePercent}%（約 ${estimatedCount} 人）`}
                {selectionMode === 'individual' && `手動挑選 ${estimatedCount} 人`}
              </Descriptions.Item>
              <Descriptions.Item label="發送方式">
                {sendMode === 'scheduled' ? `排程發送 (${scheduleStart ? new Date(scheduleStart).toLocaleString('zh-TW') : '?'} ~ ${scheduleEnd ? new Date(scheduleEnd).toLocaleString('zh-TW') : '?'})` : '立即發送'}
                {workingHoursOnly && ' · 僅工作時間'}
                {skipWeekends && ' · 避開週末'}
              </Descriptions.Item>
              <Descriptions.Item label="SMTP">
                {smtpProfile
                  ? `${smtpProfile.name}（${smtpProfile.from_address}）`
                  : <Alert type="warning" title="尚未設定 SMTP" showIcon banner />}
              </Descriptions.Item>
              <Descriptions.Item label="預估發送時間">
                {(() => {
                  const rate = smtpProfile?.mailer_type === 'ses' ? 12 : smtpProfile?.mailer_type === 'mailgun' ? 40 : 3;
                  const secs = Math.ceil(estimatedCount / rate);
                  if (secs < 60) return `約 ${secs} 秒`;
                  if (secs < 3600) return `約 ${Math.ceil(secs / 60)} 分鐘`;
                  return `約 ${(secs / 3600).toFixed(1)} 小時`;
                })()}
                <Text type="secondary" style={{ fontSize: 12 }}>（依 {smtpProfile?.mailer_type?.toUpperCase() || 'SMTP'} 速率限制自動控制）</Text>
              </Descriptions.Item>
            </Descriptions>
          </Card>

          <div style={{ marginBottom: 16, maxWidth: 500 }}>
            <Text strong style={{ marginRight: 8 }}>測試名稱</Text>
            <Input
              value={campaignName}
              onChange={e => setCampaignName(e.target.value)}
              status={trimmedCampaignName ? undefined : 'error'}
              style={{ width: 300 }}
            />
          </div>

          {!trimmedCampaignName && <Alert type="warning" showIcon message="請輸入測試名稱。" style={{ marginBottom: 16 }} />}

          {scenarioObj?.template?.html_body && (
            <>
              <Button icon={<EyeOutlined />} onClick={() => setPreviewOpen(true)} style={{ marginBottom: 24 }}>
                預覽信件
              </Button>
              <Modal
                title="郵件預覽"
                open={previewOpen}
                onCancel={() => setPreviewOpen(false)}
                footer={null}
                width={700}
              >
                {(() => {
                  // Get first recipient from selected groups for preview
                  const firstRecipient = groups
                    .filter((g: RecipientGroup) => selectedGroups.includes(String(g.id)))
                    .flatMap((g: RecipientGroup) => g.recipients ?? [])[0];
                  const vars: Record<string, string> = {
                    '{{.FirstName}}': firstRecipient?.first_name || '小明',
                    '{{.LastName}}': firstRecipient?.last_name || '王',
                    '{{.Email}}': firstRecipient?.email || 'user@example.com',
                    '{{.Department}}': firstRecipient?.department || '業務部',
                    '{{.Position}}': firstRecipient?.position || '員工',
                    '{{.TrackURL}}': '#',
                    '{{.ReportURL}}': '#',
                  };
                  let html = scenarioObj.template.html_body;
                  for (const [k, v] of Object.entries(vars)) {
                    html = html.split(k).join(v);
                  }
                  return (
                    <div>
                      <div style={{ marginBottom: 12, padding: '8px 12px', background: '#f5f5f5', borderRadius: 4, fontSize: 13 }}>
                        <strong>主旨：</strong>{scenarioObj.template.subject}<br />
                        <strong>收件人：</strong>{firstRecipient ? `${firstRecipient.last_name}${firstRecipient.first_name} <${firstRecipient.email}>` : '（請先選擇收件人群組）'}
                      </div>
                      <iframe
                        srcDoc={html}
                        style={{ width: '100%', height: 500, border: '1px solid #d9d9d9', borderRadius: 4 }}
                        sandbox=""
                        title="email-preview"
                      />
                    </div>
                  );
                })()}
              </Modal>
            </>
          )}
        </>
      )}

      {/* ── Navigation ── */}
      <div style={{ marginTop: 32, display: 'flex', justifyContent: 'space-between' }}>
        <div>
          {step > 0 && <Button onClick={() => setStep(s => s - 1)}>上一步</Button>}
        </div>
        <div>
          {step < 2 && (
            <Button type="primary" disabled={!canNext(step)} onClick={() => setStep(s => s + 1)}>
              下一步
            </Button>
          )}
          {step === 2 && (
            <Button
              type="primary"
              size="large"
              icon={<RocketOutlined />}
              loading={submitting}
              disabled={!canSubmit}
              onClick={handleSubmit}
            >
              🚀 確認發送
            </Button>
          )}
        </div>
      </div>
    </div>
  );
}
