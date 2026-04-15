import { useState } from 'react';
import { Card, Button, Typography, message, Alert, Row, Col } from 'antd';
import { ExperimentOutlined, CheckCircleOutlined } from '@ant-design/icons';
import { api } from '../../api/client';

export default function SeedData() {
  const [loading, setLoading] = useState(false);
  const [done, setDone] = useState(false);

  const seed = async () => {
    setLoading(true);
    try {
      await api.post('/seed-sample-data');
      message.success('範例資料已匯入！');
      setDone(true);
    } catch { message.error('匯入失敗，可能已存在相同資料'); }
    setLoading(false);
  };

  return (
    <div style={{ maxWidth: 700, margin: '0 auto' }}>
      <Typography.Title level={3}>匯入範例資料</Typography.Title>
      <Typography.Paragraph type="secondary">一鍵匯入完整的測試資料，讓您快速體驗所有功能。</Typography.Paragraph>

      <Card>
        <Typography.Paragraph><strong>匯入後將包含：</strong></Typography.Paragraph>
        <Row gutter={[16, 12]}>
          {[
            { icon: '✉️', title: '5 個信件模板', desc: '密碼到期、包裹通知、薪資單、資安警告、發票確認' },
            { icon: '🖥️', title: '2 個 Landing Page', desc: '仿登入頁面、確認資訊頁面（含表單擷取）' },
            { icon: '📋', title: '5 個釣魚情境', desc: '模板 + Landing Page + 教育頁完整打包' },
            { icon: '👥', title: '5 位範例收件人', desc: '分佈在 4 個部門，可直接用於測試' },
            { icon: '📚', title: '教育頁面', desc: '員工中招後看到的資安教育內容' },
          ].map((item, i) => (
            <Col span={12} key={i}>
              <Card size="small">
                <span style={{ fontSize: 20, marginRight: 8 }}>{item.icon}</span>
                <strong>{item.title}</strong>
                <Typography.Paragraph type="secondary" style={{ margin: '4px 0 0', fontSize: 12 }}>{item.desc}</Typography.Paragraph>
              </Card>
            </Col>
          ))}
        </Row>

        <Alert type="info" title="提示" description="匯入的資料會新增到現有資料中，不會覆蓋已有的模板或情境。可多次匯入（但會產生重複資料）。" showIcon style={{ margin: '16px 0' }} />

        {done ? (
          <Alert type="success" title="匯入完成" description="範例資料已成功匯入！前往情境庫或模板管理查看。" showIcon icon={<CheckCircleOutlined />} />
        ) : (
          <Button type="primary" size="large" icon={<ExperimentOutlined />} loading={loading} onClick={seed} block>
            匯入範例資料
          </Button>
        )}
      </Card>
    </div>
  );
}
