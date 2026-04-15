import { useEffect, useState } from 'react';
import { Card, Row, Col, Tag, Badge, Empty, Spin, Typography } from 'antd';
import { api } from '../../api/client';
import type { Scenario } from '../../api/client';

const categoryEmoji: Record<string, string> = {
  credential_harvest: '🔑',
  malware: '🦠',
  social_engineering: '🎭',
  spear_phishing: '🎯',
  whaling: '🐋',
  smishing: '📱',
  vishing: '📞',
};

export default function ScenarioList() {
  const [scenarios, setScenarios] = useState<Scenario[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    api.get<Scenario[]>('/scenarios').then(setScenarios).finally(() => setLoading(false));
  }, []);

  if (loading) return <Spin style={{ display: 'block', margin: '20vh auto' }} size="large" />;

  return (
    <div>
      <Typography.Title level={3}>情境庫</Typography.Title>
      {scenarios.length === 0 ? (
        <Empty description="尚無情境" />
      ) : (
        <Row gutter={[16, 16]}>
          {scenarios.map((s) => (
            <Col key={s.id} xs={24} sm={12} lg={8}>
              <Badge.Ribbon text={s.is_active ? '啟用' : '停用'} color={s.is_active ? 'green' : 'red'}>
                <Card hoverable>
                  <div style={{ fontSize: 36, textAlign: 'center', marginBottom: 8 }}>
                    {categoryEmoji[s.category] ?? '📧'}
                  </div>
                  <Typography.Title level={5} style={{ textAlign: 'center', margin: 0 }}>{s.name}</Typography.Title>
                  <div style={{ textAlign: 'center', margin: '8px 0' }}>
                    {'⭐'.repeat(Number(s.difficulty) || 1)}
                  </div>
                  <div style={{ textAlign: 'center' }}>
                    <Tag>{s.language}</Tag>
                    <Tag color="blue">{s.category}</Tag>
                    {s.tenant_id ? <Tag color="orange">自建</Tag> : <Tag color="purple">平台預建</Tag>}
                  </div>
                </Card>
              </Badge.Ribbon>
            </Col>
          ))}
        </Row>
      )}
    </div>
  );
}
