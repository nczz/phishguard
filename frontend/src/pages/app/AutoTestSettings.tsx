import { useEffect, useState } from 'react';
import { Card, Form, Switch, Radio, Slider, Button, Alert, Typography, Spin, Descriptions, message } from 'antd';
import dayjs from 'dayjs';
import { api } from '../../api/client';

interface AutoTestConfig {
  id: number;
  tenant_id: number;
  is_enabled: boolean;
  frequency: string;
  target_mode: string;
  sample_percent: number;
  next_run_at: string | null;
}

export default function AutoTestSettings() {
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [form] = Form.useForm();
  const isEnabled = Form.useWatch('is_enabled', form);
  const targetMode = Form.useWatch('target_mode', form);

  useEffect(() => {
    api.get<AutoTestConfig>('/auto-test').then((cfg) => {
      form.setFieldsValue(cfg);
    }).finally(() => setLoading(false));
  }, [form]);

  const onFinish = async (values: Record<string, unknown>) => {
    setSaving(true);
    try {
      const { is_enabled, frequency, target_mode, sample_percent } = values;
      await api.put('/auto-test', { is_enabled, frequency, target_mode, sample_percent });
      message.success('已儲存');
    } catch {
      message.error('儲存失敗');
    }
    setSaving(false);
  };

  if (loading) return <Spin style={{ display: 'block', margin: '20vh auto' }} size="large" />;

  return (
    <>
      <Typography.Title level={3}>自動定期測試</Typography.Title>
      <Typography.Paragraph type="secondary">
        開啟後系統會自動按排程執行釣魚測試，自動選人、自動選情境、自動產生報表。
      </Typography.Paragraph>

      <Card>
        <Form form={form} layout="vertical" onFinish={onFinish} initialValues={{ frequency: 'quarterly', target_mode: 'random', sample_percent: 30 }}>
          <Form.Item name="is_enabled" label="啟用自動測試" valuePropName="checked">
            <Switch checkedChildren="開" unCheckedChildren="關" />
          </Form.Item>

          {isEnabled && (
            <>
              <Form.Item name="frequency" label="測試頻率">
                <Radio.Group>
                  <Radio value="monthly">每月</Radio>
                  <Radio value="quarterly">每季</Radio>
                  <Radio value="biannual">每半年</Radio>
                </Radio.Group>
              </Form.Item>

              <Form.Item name="target_mode" label="測試對象">
                <Radio.Group>
                  <Radio value="all">全公司</Radio>
                  <Radio value="random">隨機抽樣</Radio>
                </Radio.Group>
              </Form.Item>

              {targetMode === 'random' && (
                <Form.Item name="sample_percent" label="抽樣比例">
                  <Slider min={10} max={100} step={5} marks={{ 10: '10%', 50: '50%', 100: '100%' }} tooltip={{ formatter: (v) => `${v}%` }} />
                </Form.Item>
              )}

              <Descriptions column={1} style={{ marginBottom: 24 }}>
                <Descriptions.Item label="下次執行時間">
                  {form.getFieldValue('next_run_at')
                    ? dayjs(form.getFieldValue('next_run_at')).format('YYYY-MM-DD HH:mm')
                    : '尚未排程'}
                </Descriptions.Item>
              </Descriptions>
            </>
          )}

          <Form.Item>
            <Button type="primary" htmlType="submit" loading={saving}>儲存設定</Button>
          </Form.Item>
        </Form>
      </Card>

      <Alert
        type="info"
        showIcon
        style={{ marginTop: 24 }}
        title="自動測試說明"
        description={
          <ul style={{ margin: 0, paddingLeft: 20 }}>
            <li>系統會自動從情境庫隨機選擇測試情境</li>
            <li>同一人 30 天內不會被重複測試（冷卻期）</li>
            <li>測試完成後會自動寄送報表給租戶管理員</li>
          </ul>
        }
      />
    </>
  );
}
