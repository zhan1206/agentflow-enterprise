export { Tasks } from './index';

import React from 'react';
import { Card, Table, Button, Space, Tag, Modal, Form, Input, Select, message, Progress } from 'antd';
import { PlusOutlined, PlayCircleOutlined, StopOutlined, EyeOutlined } from '@ant-design/icons';

const { Column } = Table;
const { Option } = Select;
const { TextArea } = Input;

interface Task {
  id: string;
  name: string;
  type: string;
  status: 'pending' | 'running' | 'completed' | 'failed';
  agent: string;
  progress: number;
  createdAt: string;
  duration?: string;
}

const mockTasks: Task[] = [
  { id: 'T001', name: '代码审查 - PR #234', type: 'serial', status: 'completed', agent: 'CodeReviewer-01', progress: 100, createdAt: '2024-01-15 10:30', duration: '2分30秒' },
  { id: 'T002', name: 'API 文档生成', type: 'parallel', status: 'running', agent: 'DocWriter-01', progress: 65, createdAt: '2024-01-15 11:00' },
  { id: 'T003', name: 'E2E 测试执行', type: 'dag', status: 'pending', agent: 'TestRunner-01', progress: 0, createdAt: '2024-01-15 11:15' },
  { id: 'T004', name: '数据分析 - 月报', type: 'serial', status: 'failed', agent: 'DataAnalyst-01', progress: 45, createdAt: '2024-01-15 10:00', duration: '1分20秒' },
];

export const Tasks: React.FC = () => {
  const [tasks, setTasks] = React.useState<Task[]>(mockTasks);
  const [loading, setLoading] = React.useState(false);
  const [modalVisible, setModalVisible] = React.useState(false);
  const [detailVisible, setDetailVisible] = React.useState(false);
  const [selectedTask, setSelectedTask] = React.useState<Task | null>(null);
  const [form] = Form.useForm();

  const handleRefresh = () => {
    setLoading(true);
    setTimeout(() => {
      setLoading(false);
      message.success('刷新成功');
    }, 500);
  };

  const handleExecute = (id: string) => {
    message.success(`任务 ${id} 已开始执行`);
  };

  const handleCancel = (id: string) => {
    message.warning(`任务 ${id} 已取消`);
  };

  const handleViewDetail = (task: Task) => {
    setSelectedTask(task);
    setDetailVisible(true);
  };

  const handleCreate = (values: any) => {
    console.log('Create task:', values);
    message.success('任务创建成功');
    setModalVisible(false);
    form.resetFields();
  };

  const statusColors: Record<string, string> = {
    pending: 'default',
    running: 'processing',
    completed: 'success',
    failed: 'error',
  };

  return (
    <div>
      <Card
        title="任务中心"
        extra={
          <Space>
            <Button onClick={handleRefresh}>刷新</Button>
            <Button type="primary" icon={<PlusOutlined />} onClick={() => setModalVisible(true)}>
              创建任务
            </Button>
          </Space>
        }
      >
        <Table dataSource={tasks} rowKey="id" loading={loading}>
          <Column title="任务ID" dataIndex="id" key="id" width={80} />
          <Column title="任务名称" dataIndex="name" key="name" />
          <Column
            title="类型"
            dataIndex="type"
            key="type"
            render={(type) => <Tag>{type.toUpperCase()}</Tag>}
          />
          <Column
            title="状态"
            dataIndex="status"
            key="status"
            render={(status) => <Tag color={statusColors[status]}>{status.toUpperCase()}</Tag>}
          />
          <Column title="执行Agent" dataIndex="agent" key="agent" />
          <Column
            title="进度"
            dataIndex="progress"
            key="progress"
            render={(progress: number) => <Progress percent={progress} size="small" />}
          />
          <Column title="创建时间" dataIndex="createdAt" key="createdAt" />
          <Column
            title="操作"
            key="action"
            width={200}
            render={(_, record: Task) => (
              <Space>
                <Button
                  type="link"
                  size="small"
                  icon={<PlayCircleOutlined />}
                  onClick={() => handleExecute(record.id)}
                  disabled={record.status === 'running' || record.status === 'completed'}
                >
                  执行
                </Button>
                <Button
                  type="link"
                  size="small"
                  danger
                  icon={<StopOutlined />}
                  onClick={() => handleCancel(record.id)}
                  disabled={record.status !== 'running'}
                >
                  取消
                </Button>
                <Button
                  type="link"
                  size="small"
                  icon={<EyeOutlined />}
                  onClick={() => handleViewDetail(record)}
                >
                  详情
                </Button>
              </Space>
            )}
          />
        </Table>
      </Card>

      <Modal
        title="创建任务"
        open={modalVisible}
        onCancel={() => setModalVisible(false)}
        onOk={() => form.submit()}
        width={600}
      >
        <Form form={form} layout="vertical" onFinish={handleCreate}>
          <Form.Item name="name" label="任务名称" rules={[{ required: true }]}>
            <Input placeholder="输入任务名称" />
          </Form.Item>
          <Form.Item name="type" label="执行类型" rules={[{ required: true }]}>
            <Select placeholder="选择执行类型">
              <Option value="serial">串行执行</Option>
              <Option value="parallel">并行执行</Option>
              <Option value="dag">DAG 工作流</Option>
            </Select>
          </Form.Item>
          <Form.Item name="agent" label="执行Agent" rules={[{ required: true }]}>
            <Select placeholder="选择 Agent">
              <Option value="CodeReviewer-01">CodeReviewer-01</Option>
              <Option value="DataAnalyst-01">DataAnalyst-01</Option>
              <Option value="DocWriter-01">DocWriter-01</Option>
            </Select>
          </Form.Item>
          <Form.Item name="description" label="任务描述">
            <TextArea rows={4} placeholder="描述任务内容和要求" />
          </Form.Item>
        </Form>
      </Modal>

      <Modal
        title="任务详情"
        open={detailVisible}
        onCancel={() => setDetailVisible(false)}
        footer={null}
      >
        {selectedTask && (
          <div>
            <p><strong>任务ID:</strong> {selectedTask.id}</p>
            <p><strong>任务名称:</strong> {selectedTask.name}</p>
            <p><strong>类型:</strong> {selectedTask.type}</p>
            <p><strong>状态:</strong> <Tag color={statusColors[selectedTask.status]}>{selectedTask.status.toUpperCase()}</Tag></p>
            <p><strong>执行Agent:</strong> {selectedTask.agent}</p>
            <p><strong>创建时间:</strong> {selectedTask.createdAt}</p>
            {selectedTask.duration && <p><strong>执行时长:</strong> {selectedTask.duration}</p>}
            <p><strong>进度:</strong></p>
            <Progress percent={selectedTask.progress} />
          </div>
        )}
      </Modal>
    </div>
  );
};
