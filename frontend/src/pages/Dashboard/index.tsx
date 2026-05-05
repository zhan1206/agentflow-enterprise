import React from 'react';
import { Row, Col, Card, Statistic, Table, Tag, Progress, Space, Typography } from 'antd';
import {
  RobotOutlined,
  CheckCircleOutlined,
  ClockCircleOutlined,
  DollarOutlined,
  ArrowUpOutlined,
  ArrowDownOutlined,
} from '@ant-design/icons';
import {
  LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer,
  PieChart, Pie, Cell,
} from 'recharts';

const { Title } = Typography;

// Mock data
const taskStats = {
  total: 1247,
  running: 23,
  completed: 1189,
  failed: 35,
};

const agentStats = {
  total: 15,
  online: 12,
  offline: 3,
};

const costData = {
  today: 45.67,
  week: 312.45,
  month: 1234.56,
  trend: -12.5, // percentage
};

const chartData = [
  { name: '周一', tasks: 45, cost: 32 },
  { name: '周二', tasks: 52, cost: 45 },
  { name: '周三', tasks: 61, cost: 38 },
  { name: '周四', tasks: 48, cost: 52 },
  { name: '周五', tasks: 72, cost: 61 },
  { name: '周六', tasks: 35, cost: 28 },
  { name: '周日', tasks: 28, cost: 22 },
];

const pieData = [
  { name: '已完成', value: 1189, color: '#52c41a' },
  { name: '运行中', value: 23, color: '#1890ff' },
  { name: '失败', value: 35, color: '#ff4d4f' },
];

const recentTasks = [
  { id: 'T001', name: '代码审查 - PR #234', status: 'completed', agent: 'CodeReviewer', time: '2分钟前' },
  { id: 'T002', name: '文档生成 - API文档', status: 'running', agent: 'DocWriter', time: '正在执行' },
  { id: 'T003', name: '测试执行 - E2E测试', status: 'completed', agent: 'TestRunner', time: '15分钟前' },
  { id: 'T004', name: '数据分析 - 月报', status: 'failed', agent: 'DataAnalyst', time: '30分钟前' },
  { id: 'T005', name: '部署发布 - v1.2.0', status: 'completed', agent: 'DeployBot', time: '1小时前' },
];

const statusColors: Record<string, string> = {
  completed: 'success',
  running: 'processing',
  failed: 'error',
  pending: 'default',
};

export const Dashboard: React.FC = () => {
  const columns = [
    { title: '任务ID', dataIndex: 'id', key: 'id', width: 80 },
    { title: '任务名称', dataIndex: 'name', key: 'name' },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => <Tag color={statusColors[status]}>{status.toUpperCase()}</Tag>,
    },
    { title: '执行Agent', dataIndex: 'agent', key: 'agent' },
    { title: '时间', dataIndex: 'time', key: 'time' },
  ];

  return (
    <div className="dashboard">
      <Title level={4}>控制台概览</Title>

      {/* 统计卡片 */}
      <Row gutter={[16, 16]}>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="Agent 总数"
              value={agentStats.total}
              prefix={<RobotOutlined />}
              suffix={
                <Space style={{ fontSize: 14, marginLeft: 8 }}>
                  <Tag color="success">{agentStats.online} 在线</Tag>
                  <Tag>{agentStats.offline} 离线</Tag>
                </Space>
              }
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="运行中任务"
              value={taskStats.running}
              prefix={<ClockCircleOutlined />}
              valueStyle={{ color: '#1890ff' }}
            />
            <Progress
              percent={Math.round((taskStats.running / taskStats.total) * 100)}
              showInfo={false}
              strokeColor="#1890ff"
              size="small"
              style={{ marginTop: 8 }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="已完成任务"
              value={taskStats.completed}
              prefix={<CheckCircleOutlined />}
              valueStyle={{ color: '#52c41a' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card>
            <Statistic
              title="今日成本"
              value={costData.today}
              prefix={<DollarOutlined />}
              precision={2}
              suffix="USD"
              valueStyle={{ color: costData.trend < 0 ? '#52c41a' : '#ff4d4f' }}
            />
            <div style={{ marginTop: 8 }}>
              {costData.trend < 0 ? (
                <ArrowDownOutlined style={{ color: '#52c41a', marginRight: 4 }} />
              ) : (
                <ArrowUpOutlined style={{ color: '#ff4d4f', marginRight: 4 }} />
              )}
              <span style={{ color: costData.trend < 0 ? '#52c41a' : '#ff4d4f' }}>
                {Math.abs(costData.trend)}% 较昨日
              </span>
            </div>
          </Card>
        </Col>
      </Row>

      {/* 图表 */}
      <Row gutter={[16, 16]} style={{ marginTop: 16 }}>
        <Col xs={24} lg={16}>
          <Card title="任务执行趋势">
            <ResponsiveContainer width="100%" height={300}>
              <LineChart data={chartData}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="name" />
                <YAxis yAxisId="left" />
                <YAxis yAxisId="right" orientation="right" />
                <Tooltip />
                <Legend />
                <Line yAxisId="left" type="monotone" dataKey="tasks" stroke="#1890ff" name="任务数" />
                <Line yAxisId="right" type="monotone" dataKey="cost" stroke="#52c41a" name="成本($)" />
              </LineChart>
            </ResponsiveContainer>
          </Card>
        </Col>
        <Col xs={24} lg={8}>
          <Card title="任务状态分布">
            <ResponsiveContainer width="100%" height={300}>
              <PieChart>
                <Pie
                  data={pieData}
                  cx="50%"
                  cy="50%"
                  innerRadius={60}
                  outerRadius={100}
                  paddingAngle={5}
                  dataKey="value"
                  label={({ name, percent }) => `${name} ${(percent * 100).toFixed(0)}%`}
                >
                  {pieData.map((entry, index) => (
                    <Cell key={`cell-${index}`} fill={entry.color} />
                  ))}
                </Pie>
                <Tooltip />
              </PieChart>
            </ResponsiveContainer>
          </Card>
        </Col>
      </Row>

      {/* 最近任务 */}
      <Card title="最近任务" style={{ marginTop: 16 }}>
        <Table
          columns={columns}
          dataSource={recentTasks}
          rowKey="id"
          pagination={false}
          size="small"
        />
      </Card>
    </div>
  );
};
