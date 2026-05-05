# Contributing to AgentFlow-Enterprise

感谢你对 AgentFlow-Enterprise 的兴趣！我们欢迎所有形式的贡献。

## 🌟 贡献方式

- 提交 Bug 报告或功能建议
- 改进文档
- 提交代码修复或新功能
- 分享使用案例和最佳实践
- 帮助其他用户

## 🔧 开发环境设置

### 前置要求

- Go 1.21+
- Node.js 18+
- Docker & Docker Compose
- Make (可选)

### 后端开发

```bash
cd backend
go mod download
go run ./cmd/server --config ../deploy/configs/development.yaml
```

### 前端开发

```bash
cd frontend
npm install
npm run dev
```

### 运行测试

```bash
# 后端测试
cd backend && go test ./...

# 前端测试
cd frontend && npm test

# 集成测试
make test-integration
```

## 📝 代码规范

### Go 代码

- 遵循 [Effective Go](https://golang.org/doc/effective_go)
- 使用 `gofmt` 格式化代码
- 使用 `golint` 和 `go vet` 检查
- 为公共 API 编写文档注释

### TypeScript/React 代码

- 遵循 Airbnb Style Guide
- 使用 Prettier 格式化
- 使用 ESLint 检查
- 组件使用函数式组件 + Hooks

### 提交信息

使用 [Conventional Commits](https://www.conventionalcommits.org/):

```
feat: 添加多租户支持
fix: 修复状态同步竞态条件
docs: 更新部署文档
refactor: 重构调度器核心逻辑
test: 添加集成测试
chore: 更新依赖
```

## 🔄 Pull Request 流程

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 进行更改并确保测试通过
4. 提交符合规范的提交信息
5. 推送到你的 Fork (`git push origin feature/amazing-feature`)
6. 创建 Pull Request

### PR 检查清单

- [ ] 代码通过所有测试
- [ ] 新功能有对应的测试用例
- [ ] 更新了相关文档
- [ ] 提交信息符合规范
- [ ] 没有引入新的 warning

## 🏗️ 项目结构

详见 [README.md](README.md#项目结构)

## 🐛 报告 Bug

使用 [GitHub Issues](https://github.com/agentflow-enterprise/agentflow-enterprise/issues) 报告 Bug，请包含：

- 操作系统和版本
- AgentFlow 版本
- 复现步骤
- 期望行为
- 实际行为
- 日志/截图（如果适用）

## 💡 功能建议

使用 [GitHub Discussions](https://github.com/agentflow-enterprise/agentflow-enterprise/discussions) 提出功能建议。

## 📄 许可证

通过贡献代码，你同意你的代码将以 Apache 2.0 许可证发布。

---

再次感谢你的贡献！🎉
