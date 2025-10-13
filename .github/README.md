# GitHub 配置文件说明

本目录包含 aicommit 项目的所有 GitHub 配置文件。

## 📁 目录结构

```
.github/
├── workflows/              # GitHub Actions workflows
│   ├── ci.yml             # 持续集成
│   └── release.yml        # 自动发布
├── ISSUE_TEMPLATE/        # Issue 模板
│   ├── bug_report.md      # Bug 报告模板
│   ├── feature_request.md # 功能请求模板
│   └── config.yml         # Issue 配置
├── CODEOWNERS             # 代码所有者
├── CONTRIBUTING.md        # 贡献指南
├── PULL_REQUEST_TEMPLATE.md # PR 模板
├── WORKFLOWS.md           # Workflows 使用指南
├── dependabot.yml         # Dependabot 配置
└── README.md             # 本文件
```

## 🚀 Workflows

### CI Workflow
- **文件**: `workflows/ci.yml`
- **功能**: 自动测试、代码检查、构建
- **触发**: 推送到主分支或创建 PR

### Release Workflow
- **文件**: `workflows/release.yml`
- **功能**: 多平台构建和发布
- **触发**: 推送版本 tag (v*)

详情请查看 [WORKFLOWS.md](WORKFLOWS.md)

## 📝 模板

### Issue 模板
- **Bug Report**: 报告 bug
- **Feature Request**: 请求新功能

### PR 模板
包含完整的检查清单和说明

## 🤝 贡献

请查看 [CONTRIBUTING.md](CONTRIBUTING.md) 了解：
- 如何设置开发环境
- 代码规范
- 提交流程
- 测试要求

## 🔧 配置文件

### CODEOWNERS
定义代码审查责任人

### dependabot.yml
自动依赖更新配置：
- Go 模块：每周检查
- GitHub Actions：每周检查

## 📚 相关文件

项目根目录中的相关文件：
- `/.golangci.yml` - Linter 配置
- `/CHANGELOG.md` - 变更日志
- `/SECURITY.md` - 安全政策

## 🎯 使用说明

1. **首次使用**：
   - 将所有 `lemon956` 替换为实际的 GitHub 用户名
   - 更新 CODEOWNERS 中的维护者信息
   - 配置必要的 GitHub Secrets

2. **发布新版本**：
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```

3. **查看 CI 状态**：
   - 访问 Actions 标签页
   - 查看 workflow 运行历史

## 🔗 快速链接

- [Actions](https://github.com/lemon956/aicommit/actions)
- [Releases](https://github.com/lemon956/aicommit/releases)
- [Issues](https://github.com/lemon956/aicommit/issues)
- [Pull Requests](https://github.com/lemon956/aicommit/pulls)
- [Discussions](https://github.com/lemon956/aicommit/discussions)

---

💡 **提示**: 记得将所有 URL 中的 `lemon956` 替换为实际的 GitHub 用户名！

