import React from 'react';
import { Card, Row, Col, Typography, Statistic, Empty } from 'antd';
import { DashboardOutlined, AlertOutlined, ClockCircleOutlined } from '@ant-design/icons';

const { Title, Paragraph } = Typography;

export const Observability: React.FC = () => {
  return (
    <div>
      <Title level={4}>可观测性</Title>
      <Paragraph type="secondary">
        全链路追踪、指标监控、日志分析
      </Paragraph>

      <Row gutter={[16, 16]} style={{ marginTop: 16 }}>
        <Col xs={24} sm={8}>
          <Card>
            <Statistic
              title="追踪链路"
              value={1234}
              prefix={<DashboardOutlined />}
            />
          </Card>
        </Col>
        <Col xs={24} sm={8}>
          <Card>
            <Statistic
              title="告警数"
              value={5}
              prefix={<AlertOutlined />}
              valueStyle={{ color: '#ff4d4f' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={8}>
          <Card>
            <Statistic
              title="平均延迟"
              value={234}
              suffix="ms"
              prefix={<ClockCircleOutlined />}
            />
          </Card>
        </Col>
      </Row>

      <Card title="追踪面板" style={{ marginTop: 16 }}>
        <Empty description="Jaeger 集成即将上线" />
      </Card>

      <Card title="指标仪表盘" style={{ marginTop: 16 }}>
        <Empty description="Prometheus + Grafana 集成即将上线" />
      </Card>
    </div>
  );
};
