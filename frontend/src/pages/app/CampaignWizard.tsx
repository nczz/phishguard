import { useEffect, useMemo, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Steps, Card, Row, Col, Button, Checkbox, Radio, Slider, Input,
  Tag, Modal, Spin, Typography, message, Space, Alert, Descriptions, Tooltip,
} from 'antd';
import { RocketOutlined, EyeOutlined, QuestionCircleOutlined } from '@ant-design/icons';
import { tips } from '../../components/FieldHelp';
import type { Scenario, RecipientGroup, SMTPProfile, Campaign } from '../../api/client';
import { api } from '../../api/client';

const { Title, Text } = Typography;

type SelectionMode = 'all' | 'department' | 'sample';

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

function difficultyStars(d: string) {
  const n = parseInt(d, 10) || 1;
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
  const [spreadSend, setSpreadSend] = useState(false);
  const [campaignName, setCampaignName] = useState(defaultCampaignName);

  useEffect(() => {
    Promise.all([
      api.get<Scenario[]>('/scenarios'),
      api.get<RecipientGroup[]>('/recipient-groups'),
      api.get<SMTPProfile[]>('/smtp-profiles'),
    ]).then(([s, g, p]) => {
      setScenarios(s);
      setGroups(g);
      setSmtpProfiles(p);
    }).finally(() => setLoading(false));
  }, []);

  // Derived
  const allRecipients = useMemo(
    () => groups.filter(g => selectedGroups.includes(g.id)).flatMap(g => g.recipients ?? []),
    [groups, selectedGroups],
  );

  const allDepartments = useMemo(
    () => [...new Set(allRecipients.map(r => r.department).filter(Boolean))].sort(),
    [allRecipients],
  );

  const estimatedCount = useMemo(() => {
    let pool = allRecipients;
    if (selectionMode === 'department') {
      pool = pool.filter(r => departments.includes(r.department));
    }
    if (selectionMode === 'sample') {
      return Math.max(1, Math.round(pool.length * samplePercent / 100));
    }
    return pool.length;
  }, [allRecipients, selectionMode, departments, samplePercent]);

  const scenarioObj = scenarios.find(s => s.id === selectedScenario);
  const smtpProfile = smtpProfiles[0] ?? null;

  const canNext = (s: number) => {
    if (s === 0) return autoRandom || !!selectedScenario;
    if (s === 1) return selectedGroups.length > 0 && (selectionMode !== 'department' || departments.length > 0);
    return true;
  };

  async function handleSubmit() {
    setSubmitting(true);
    try {
      const campaign = await api.post<Campaign>('/campaigns', {
        name: campaignName,
        scenario_id: autoRandom ? null : selectedScenario,
        smtp_profile_id: smtpProfile?.id,
        group_ids: selectedGroups,
        phish_url: window.location.origin + '/phish',
        selection_mode: selectionMode,
        sample_percent: selectionMode === 'sample' ? samplePercent : undefined,
        departments: selectionMode === 'department' ? departments : undefined,
        spread_send: spreadSend,
      });
      await api.post('/campaigns/' + campaign.id + '/launch');
      message.success('測試已發送！');
      navigate('/app/campaigns/' + campaign.id);
    } catch {
      message.error('發送失敗，請稍後再試');
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
        </>
      )}

      {/* ── Step 2: 選擇對象 ── */}
      {step === 1 && (
        <>
          <Title level={4}>選擇測試對象</Title>

          <Text strong style={{ display: 'block', marginBottom: 8 }}>選擇收件人群組</Text>
          <Checkbox.Group
            value={selectedGroups}
            onChange={v => setSelectedGroups(v as string[])}
            style={{ marginBottom: 16 }}
          >
            <Space orientation="vertical">
              {groups.map(g => (
                <Checkbox key={g.id} value={g.id}>
                  {g.name}（{g.recipients?.length ?? 0} 人）
                </Checkbox>
              ))}
            </Space>
          </Checkbox.Group>

          <Text strong style={{ display: 'block', marginBottom: 8 }}>發送範圍</Text>
          <Radio.Group
            value={selectionMode}
            onChange={e => setSelectionMode(e.target.value)}
            style={{ marginBottom: 16 }}
          >
            <Space orientation="vertical">
              <Radio value="all">全公司（共 {allRecipients.length} 人） <Tooltip title={tips.selectionAll}><QuestionCircleOutlined style={{color:'#999'}} /></Tooltip></Radio>
              <Radio value="department">指定部門 <Tooltip title={tips.selectionDept}><QuestionCircleOutlined style={{color:'#999'}} /></Tooltip></Radio>
              <Radio value="sample">隨機抽樣 <Tooltip title={tips.selectionRandom}><QuestionCircleOutlined style={{color:'#999'}} /></Tooltip></Radio>
            </Space>
          </Radio.Group>

          {selectionMode === 'department' && (
            <Checkbox.Group
              value={departments}
              onChange={v => setDepartments(v as string[])}
              style={{ display: 'block', marginBottom: 16, paddingLeft: 24 }}
            >
              <Space orientation="vertical">
                {allDepartments.map(d => (
                  <Checkbox key={d} value={d}>{d}</Checkbox>
                ))}
              </Space>
            </Checkbox.Group>
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

          <Checkbox checked={spreadSend} onChange={e => setSpreadSend(e.target.checked)}>
            分散發送（避免員工互相通風報信） <Tooltip title={tips.spreadSend}><QuestionCircleOutlined style={{color:'#999'}} /></Tooltip>
          </Checkbox>
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
              </Descriptions.Item>
              <Descriptions.Item label="發送方式">
                {spreadSend ? '分散發送' : '立即發送'}
              </Descriptions.Item>
              <Descriptions.Item label="SMTP">
                {smtpProfile
                  ? `${smtpProfile.name}（${smtpProfile.from_address}）`
                  : <Alert type="warning" title="尚未設定 SMTP" showIcon banner />}
              </Descriptions.Item>
            </Descriptions>
          </Card>

          <Input
            addonBefore="測試名稱"
            value={campaignName}
            onChange={e => setCampaignName(e.target.value)}
            style={{ marginBottom: 16, maxWidth: 500 }}
          />

          {scenarioObj?.template?.html_body && (
            <>
              <Button icon={<EyeOutlined />} onClick={() => setPreviewOpen(true)} style={{ marginBottom: 24 }}>
                Preview email
              </Button>
              <Modal
                title="郵件預覽"
                open={previewOpen}
                onCancel={() => setPreviewOpen(false)}
                footer={null}
                width={640}
              >
                <div dangerouslySetInnerHTML={{ __html: scenarioObj.template.html_body }} />
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
              disabled={!smtpProfile}
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
