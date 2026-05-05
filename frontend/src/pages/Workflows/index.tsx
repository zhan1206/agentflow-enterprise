export { Workflows } from './index';

import React from 'react';
import { Card, Button, Space, message, Empty, Row, Col, Typography } from 'antd';
import { PlusOutlined, FolderOpenOutlined } from '@ant-design/icons';

const { Title, Paragraph } = Typography;

export const Workflows: React.FC = () => {
  const handleCreate = () => {
    message.info('工作流编排器即将上线');
  };

  return (
    <div>
      <Card
        title="工作流编排"
        extra={
          <Space>
            <Button icon={<FolderOpenOutlined />}>模板库</Button>
            <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate}>
              新建工作流
            </Button>
          </Space>
        }
      >
        <Empty
          description={
            <div>
              <Title level={5}>工作流编排器即将上线</Title>
              <Paragraph type="secondary">
                支持可视化拖拽式工作流设计，DAG 依赖管理，多 Agent 协同编排
              </Paragraph>
            </div>
          }
        />
      </Card>
    </div>
  );
};
