import React from 'react';
import { Card, Table, Button, Space, Tag, Modal, Form, Input, Select, message } from 'antd';
import { PlusOutlined, ReloadOutlined, PlayCircleOutlined, StopOutlined } from '@ant-design/icons';

const { Column } = Table;
const { Option } = Select;

interface Agent {
  id: string;
  name: string;
  type: string;
  status: 'online' | 'offline' | 'busy';
  capabilities: string[];
  endpoint: string;
}

const mockAgents: Agent[] = [
  { id: '1', name: 'CodeReviewer-01', type: 'openhands', status: 'online', capabilities: ['code-review', 'debugging'], endpoint: 'http://localhost:5001' },
  { id: '2', name: 'DataAnalyst-01', type: 'langgraph', status: 'busy', capabilities: ['data-analysis', 'visualization'], endpoint: 'http://localhost:5002' },
  { id: '3', name: 'DocWriter-01', type: 'crewai', status: 'online', capabilities: ['documentation', 'translation'], endpoint: 'http://localhost:5003' },
];

export const Agents: React.FC = () => {
  const [agents, setAgents] = React.useState<Agent[]>(mockAgents);
  const [loading, setLoading] = React.useState(false);
  const [modalVisible, setModalVisible] = React.useState(false);
  const [form] = Form.useForm();

  const handleRefresh = () => {
    setLoading(true);
    setTimeout(() => {
      setLoading(false);
      message.success('刷新成功');
    }, 500);
  };

  const handleAdd = () => {
    setModalVisible(true);
  };

  const handleCreate = (values: any) => {
    console.log('Create agent:', values);
    message.success('Agent 创建成功');
    setModalVisible(false);
    form.resetFields();
  };

  const handleStart = (id: string) => {
    message.success(`Agent ${id} 已启动`);
  };

  const handleStop = (id: string) => {
    message.warning(`Agent ${id} 已停止`);
  };

  const statusColors: Record<string, string> = {
    online: 'success',
    offline: 'default',
    busy: 'processing',
  };

  return (
    <div>
      <Card
        title="Agent 管理"
        extra={
          <Space>
            <Button icon={<ReloadOutlined />} onClick={handleRefresh}>
              刷新
            </Button>
            <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>
              添加 Agent
            </Button>
          </Space>
        }
      >
        <Table dataSource={agents} rowKey="id" loading={loading}>
          <Column title="ID" dataIndex="id" key="id" width={80} />
          <Column title="名称" dataIndex="name" key="name" />
          <Column
            title="类型"
            dataIndex="type"
            key="type"
            render={(type) => <Tag color="blue">{type.toUpperCase()}</Tag>}
          />
          <Column
            title="状态"
            dataIndex="status"
            key="status"
            render={(status) => <Tag color={statusColors[status]}>{status.toUpperCase()}</Tag>}
          />
          <Column
            title="能力"
            dataIndex="capabilities"
            key="capabilities"
            render={(caps: string[]) => caps.map((c) => <Tag key={c}>{c}</Tag>)}
          />
          <Column title="端点" dataIndex="endpoint" key="endpoint" ellipsis />
          <Column
            title="操作"
            key="action"
            width={150}
            render={(_, record: Agent) => (
              <Space>
                <Button
                  type="link"
                  size="small"
                  icon={<PlayCircleOutlined />}
                  onClick={() => handleStart(record.id)}
                  disabled={record.status === 'busy'}
                >
                  启动
                </Button>
                <Button
                  type="link"
                  size="small"
                  danger
                  icon={<StopOutlined />}
                  onClick={() => handleStop(record.id)}
                  disabled={record.status === 'offline'}
                >
                  停止
                </Button>
              </Space>
            )}
          />
        </Table>
      </Card>

      <Modal
        title="添加 Agent"
        open={modalVisible}
        onCancel={() => setModalVisible(false)}
        onOk={() => form.submit()}
      >
        <Form form={form} layout="vertical" onFinish={handleCreate}>
          <Form.Item name="name" label="名称" rules={[{ required: true }]}>
            <Input placeholder="输入 Agent 名称" />
          </Form.Item>
          <Form.Item name="type" label="类型" rules={[{ required: true }]}>
            <Select placeholder="选择 Agent 类型">
              <Option value="openhands">OpenHands</Option>
              <Option value="langgraph">LangGraph</Option>
              <Option value="crewai">CrewAI</Option>
              <Option value="autogen">AutoGen</Option>
            </Select>
          </Form.Item>
          <Form.Item name="endpoint" label="端点" rules={[{ required: true }]}>
            <Input placeholder="http://localhost:5000" />
          </Form.Item>
          <Form.Item name="capabilities" label="能力">
            <Select mode="multiple" placeholder="选择能力">
              <Option value="code-generation">代码生成</Option>
              <Option value="code-review">代码审查</Option>
              <Option value="debugging">调试</Option>
              <Option value="testing">测试</Option>
              <Option value="documentation">文档</Option>
              <Option value="data-analysis">数据分析</Option>
            </Select>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};
