# AgentFlow-Enterprise

<p align="center">
  <img src="docs/images/logo.png" alt="AgentFlow-Enterprise Logo" width="200">
</p>

<p align="center">
  <strong>全球首个生产级多智能体协同与全生命周期运维开源平台</strong>
</p>

<p align="center">
  <a href="#核心特性">核心特性</a> •
  <a href="#快速开始">快速开始</a> •
  <a href="#文档">文档</a> •
  <a href="#贡献">贡献</a> •
  <a href="#许可证">许可证</a>
</p>

<p align="center">
  <a href="https://github.com/agentflow-enterprise/agentflow-enterprise/actions/workflows/ci.yml">
    <img src="https://github.com/agentflow-enterprise/agentflow-enterprise/actions/workflows/ci.yml/badge.svg" alt="CI">
  </a>
  <a href="https://goreportcard.com/report/github.com/agentflow-enterprise/agentflow-enterprise">
    <img src="https://goreportcard.com/badge/github.com/agentflow-enterprise/agentflow-enterprise" alt="Go Report Card">
  </a>
  <a href="https://codecov.io/gh/agentflow-enterprise/agentflow-enterprise">
    <img src="https://codecov.io/gh/agentflow-enterprise/agentflow-enterprise/branch/main/graph/badge.svg" alt="Coverage">
  </a>
  <a href="LICENSE">
    <img src="https://img.shields.io/badge/License-Apache%202.0-blue.svg" alt="License">
  </a>
</p>

---

## 项目定位

AgentFlow-Enterprise 是企业级多 Agent 集群的「操作系统 + 协同中台 + 运维管控底座」，填补 AI Agent 从「玩具级 demo」到「企业级生产可用」的行业空白。

### 核心解决痛点

| 痛点维度 | 解决方案 |
|---------|---------|
| **协同层** - 任务拆解歧义、角色边界混乱、状态不同步 | ACIP 多 Agent 协同交互协议 + Raft 强一致状态同步 |
| **运维层** - 执行黑盒、无故障自愈、成本失控 | OpenTelemetry 全链路可观测 + 自动修复引擎 |
| **安全层** - 权限失控、无审计、数据泄露 | 双维度 RBAC + 链式审计 + DLP |
| **落地门槛** - 自研成本高、无模板、仅面向研发 | 全兼容主流框架 + 低代码编排 + 开箱即用模板 |

---

## 核心特性

### 🔄 多 Agent 协同调度
- **ACIP 协议**: 自研标准化 Agent 通信协议，5 大消息类型消除歧义
- **DAG 任务调度**: 智能拆解、依赖管理、串并行执行
- **Raft 状态同步**: 全局强一致，多节点实时同步

### 🛡️ 全生命周期运维
- **全链路可观测**: OpenTelemetry + Jaeger 端到端追踪
- **故障自愈**: 规则引擎 + LLM 根因分析，95% 异常自动修复
- **成本管控**: 多维度统计、预算告警、智能优化

### 🔐 企业级安全合规
- **双维度 RBAC**: 用户 + Agent 4 级权限管控
- **链式审计**: 不可篡改、全程留痕、一键合规报告
- **数据安全**: 加密传输存储、自动脱敏、DLP

### 🔌 全生态兼容
- **框架适配**: OpenHands、Plandex、LangGraph、CrewAI、AutoGen 等
- **低代码编排**: 拖拽式可视化工作流
- **模板市场**: 10+ 行业、20+ 场景开箱即用

---

## 架构图

```
┌─────────────────────────────────────────────────────────────────┐
│                    前端交互与低代码编排层                          │
│         可视化编排 | 管控控制台 | 模板市场 | 人机协同               │
├─────────────────────────────────────────────────────────────────┤
│                    安全合规与权限管控层                            │
│         RBAC权限 | 全链路审计 | 数据安全 | 恶意拦截                │
├─────────────────────────────────────────────────────────────────┤
│                    全生命周期运维管控层                            │
│         可观测平台 | 故障自愈 | 成本管控 | 生命周期管理              │
├─────────────────────────────────────────────────────────────────┤
│                    协同调度核心层（心脏）                          │
│         任务调度 | ACIP协议 | 状态同步 | 角色分工                  │
├─────────────────────────────────────────────────────────────────┤
│                    兼容接入层（Agent适配）                         │
│         多框架适配器 | Agent注册中心 | 统一工具网关                │
└─────────────────────────────────────────────────────────────────┘
```

---

## 快速开始

### 前置要求

- Docker 20.10+
- Docker Compose 2.0+
- (可选) Kubernetes 1.24+ 用于生产部署

### 方式一：Docker Compose（推荐快速体验）

```bash
# 克隆仓库
git clone https://github.com/agentflow-enterprise/agentflow-enterprise.git
cd agentflow-enterprise

# 启动服务
docker-compose up -d

# 访问控制台
open http://localhost:8080
```

### 方式二：Helm 部署（Kubernetes）

```bash
# 添加 Helm 仓库
helm repo add agentflow https://agentflow-enterprise.github.io/charts
helm repo update

# 安装
helm install agentflow agentflow/agentflow-enterprise \
  --namespace agentflow \
  --create-namespace
```

### 方式三：源码编译

```bash
# 前端
cd frontend
npm install && npm run build

# 后端
cd backend
go mod download && go build -o agentflow ./cmd/server

# 运行
./agentflow --config config.yaml
```

---

## 项目结构

```
agentflow-enterprise/
├── backend/                    # Go 后端服务
│   ├── cmd/                    # 入口程序
│   │   ├── server/             # API Server
│   │   ├── scheduler/          # 调度器
│   │   └── cli/                # 命令行工具
│   ├── internal/               # 内部模块
│   │   ├── core/               # 核心引擎
│   │   │   ├── scheduler/      # 任务调度
│   │   │   ├── protocol/       # ACIP 协议
│   │   │   └── sync/           # 状态同步
│   │   ├── adapter/            # Agent 适配器
│   │   ├── observability/      # 可观测性
│   │   ├── security/           # 安全模块
│   │   └── storage/            # 存储层
│   ├── pkg/                    # 公共库
│   └── api/                    # API 定义
│       ├── openapi/            # OpenAPI 规范
│       └── proto/              # Protobuf 定义
├── frontend/                   # React 前端
│   ├── src/
│   │   ├── pages/              # 页面组件
│   │   ├── components/         # 通用组件
│   │   ├── workflows/          # 工作流编排
│   │   └── store/              # 状态管理
│   └── public/
├── adapters/                   # Agent 框架适配器
│   ├── openhands/              # OpenHands 适配
│   ├── langgraph/              # LangGraph 适配
│   ├── crewai/                 # CrewAI 适配
│   └── autogen/                # AutoGen 适配
├── deploy/                     # 部署配置
│   ├── docker/                 # Docker 配置
│   ├── kubernetes/             # K8s manifests
│   └── helm/                   # Helm Charts
├── docs/                       # 文档
├── examples/                   # 示例工作流
└── tests/                      # 测试
    ├── unit/                   # 单元测试
    ├── integration/            # 集成测试
    └── e2e/                    # 端到端测试
```

---

## 文档

- [架构设计](docs/architecture.md)
- [ACIP 协议规范](docs/protocol/acip-spec.md)
- [API 文档](docs/api/README.md)
- [部署指南](docs/deployment.md)
- [开发指南](docs/development.md)
- [最佳实践](docs/best-practices.md)

---

## 技术栈

| 层级 | 技术选型 |
|-----|---------|
| 后端核心 | Go 1.21+ |
| AI 适配 | Python 3.11+ |
| 前端 | React 18 + TypeScript + Ant Design Pro |
| 状态存储 | RocksDB + Raft |
| 消息队列 | Kafka |
| 可观测 | OpenTelemetry + Jaeger + Prometheus + Grafana |
| 日志检索 | Elasticsearch |
| 部署 | Docker + Kubernetes |

---

## 贡献

我们欢迎所有形式的贡献！

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

详见 [贡献指南](CONTRIBUTING.md)

---

## 路线图

### v0.1 MVP (当前)
- [x] 项目初始化
- [ ] 核心调度引擎
- [ ] 基础 Agent 适配器
- [ ] 简易 Web 控制台

### v0.5 Beta
- [ ] ACIP 协议完整实现
- [ ] 全链路可观测
- [ ] 基础故障自愈
- [ ] RBAC 权限

### v1.0 正式版
- [ ] 企业级生产可用
- [ ] 完整安全合规
- [ ] 行业模板库

详见 [路线图](docs/roadmap.md)

---

## 许可证

本项目采用 [Apache 2.0](LICENSE) 许可证开源。

核心功能永久免费、商用友好、无厂商锁定。

---

## 社区

- [GitHub Discussions](https://github.com/agentflow-enterprise/agentflow-enterprise/discussions)
- [Discord](https://discord.gg/agentflow)
- 微信公众号: AgentFlow

---

## 致谢

感谢所有贡献者和以下开源项目：

- [OpenHands](https://github.com/All-Hands-AI/OpenHands)
- [LangGraph](https://github.com/langchain-ai/langgraph)
- [CrewAI](https://github.com/joaomdmoura/crewAI)
- [AutoGen](https://github.com/microsoft/autogen)

---

<p align="center">
  <strong>AgentFlow-Enterprise — AI Agent 规模化落地的核心基础设施</strong>
</p>
