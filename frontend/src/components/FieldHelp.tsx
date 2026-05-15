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
