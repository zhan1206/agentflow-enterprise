import React from 'react';
import { Card, Form, Input, Switch, Button, message, Divider, Typography } from 'antd';

const { Title } = Typography;

export const Settings: React.FC = () => {
  const [form] = Form.useForm();

  const handleSave = () => {
    message.success('设置已保存');
  };

  return (
    <div>
      <Title level={4}>系统设置</Title>

      <Card title="基础配置" style={{ marginTop: 16 }}>
        <Form form={form} layout="vertical" onFinish={handleSave}>
          <Form.Item label="系统名称" name="systemName" initialValue="AgentFlow Enterprise">
            <Input />
          </Form.Item>
          <Form.Item label="API 端点" name="apiEndpoint" initialValue="http://localhost:8080">
            <Input />
          </Form.Item>
          <Form.Item label="启用调试模式" name="debugMode" valuePropName="checked" initialValue={false}>
            <Switch />
          </Form.Item>
          <Button type="primary" htmlType="submit">
            保存
          </Button>
        </Form>
      </Card>

      <Card title="安全配置" style={{ marginTop: 16 }}>
        <Form layout="vertical">
          <Form.Item label="JWT 密钥" name="jwtSecret">
            <Input.Password placeholder="输入 JWT 密钥" />
          </Form.Item>
          <Form.Item label="启用审计日志" name="auditLog" valuePropName="checked" initialValue={true}>
            <Switch />
          </Form.Item>
          <Button type="primary">保存</Button>
        </Form>
      </Card>

      <Card title="通知配置" style={{ marginTop: 16 }}>
        <Form layout="vertical">
          <Form.Item label="邮件服务器" name="smtpServer">
            <Input placeholder="smtp.example.com" />
          </Form.Item>
          <Form.Item label="Webhook URL" name="webhookUrl">
            <Input placeholder="https://webhook.example.com" />
          </Form.Item>
          <Button type="primary">保存</Button>
        </Form>
      </Card>
    </div>
  );
};
