# 贡献指南 / Contributing Guide

感谢你考虑为 aicommit 做出贡献！

Thank you for considering contributing to aicommit!

## 📋 目录 / Table of Contents

- [行为准则 / Code of Conduct](#行为准则--code-of-conduct)
- [如何贡献 / How to Contribute](#如何贡献--how-to-contribute)
- [开发环境设置 / Development Setup](#开发环境设置--development-setup)
- [提交代码 / Submitting Changes](#提交代码--submitting-changes)
- [代码规范 / Coding Standards](#代码规范--coding-standards)
- [测试 / Testing](#测试--testing)

## 行为准则 / Code of Conduct

我们致力于维护一个开放和友好的社区。请尊重所有贡献者。

We are committed to maintaining an open and welcoming community. Please respect all contributors.

## 如何贡献 / How to Contribute

### 报告 Bug / Reporting Bugs

使用 [Bug Report 模板](https://github.com/lemon956/aicommit/issues/new?template=bug_report.md) 创建 issue。

Use the [Bug Report template](https://github.com/lemon956/aicommit/issues/new?template=bug_report.md) to create an issue.

### 建议功能 / Suggesting Features

使用 [Feature Request 模板](https://github.com/lemon956/aicommit/issues/new?template=feature_request.md) 创建 issue。

Use the [Feature Request template](https://github.com/lemon956/aicommit/issues/new?template=feature_request.md) to create an issue.

### 提交 PR / Submitting Pull Requests

1. Fork 本仓库 / Fork the repository
2. 创建功能分支 / Create a feature branch
   ```bash
   git checkout -b feat/amazing-feature
   ```
3. 提交更改 / Commit your changes
   ```bash
   git commit -m "feat: add amazing feature"
   ```
4. 推送到分支 / Push to the branch
   ```bash
   git push origin feat/amazing-feature
   ```
5. 创建 Pull Request / Open a Pull Request

## 开发环境设置 / Development Setup

### 前置要求 / Prerequisites

- Go 1.21 或更高版本 / Go 1.21 or higher
- Git
- Make (可选 / optional)

### 安装 / Installation

```bash
# 克隆仓库 / Clone the repository
git clone https://github.com/lemon956/aicommit.git
cd aicommit

# 安装依赖 / Install dependencies
go mod download

# 构建 / Build
make build

# 运行测试 / Run tests
make test
```

### 项目结构 / Project Structure

```
aicommit/
├── cmd/aicommit/          # 主程序入口 / Main application entry
├── internal/              # 内部包 / Internal packages
│   ├── config/           # 配置管理 / Configuration management
│   ├── git/              # Git 操作 / Git operations
│   └── model/            # AI 模型集成 / AI model integration
├── pkg/                   # 公共包 / Public packages
│   ├── prompt/           # 提示词管理 / Prompt management
│   └── validator/        # 验证器 / Validators
└── .github/              # GitHub 配置 / GitHub configuration
```

## 提交代码 / Submitting Changes

### Commit 消息格式 / Commit Message Format

我们遵循 [Conventional Commits](https://www.conventionalcommits.org/) 规范：

We follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

**类型 / Types:**
- `feat`: 新功能 / New feature
- `fix`: Bug 修复 / Bug fix
- `docs`: 文档更新 / Documentation
- `style`: 代码格式 / Code style
- `refactor`: 重构 / Refactoring
- `perf`: 性能优化 / Performance
- `test`: 测试 / Testing
- `chore`: 构建/工具 / Build/tooling
- `ci`: CI/CD 配置 / CI/CD configuration

**示例 / Examples:**
```
feat(model): add support for GPT-4
fix(git): handle empty diff correctly
docs: update installation guide
refactor(prompt): simplify template structure
```

## 代码规范 / Coding Standards

### Go 代码风格 / Go Code Style

- 遵循 [Effective Go](https://golang.org/doc/effective_go.html) 指南
- 使用 `gofmt` 格式化代码
- 使用 `golangci-lint` 进行代码检查

```bash
# 格式化代码 / Format code
make fmt

# 运行 linter / Run linter
make lint
```

### 代码质量 / Code Quality

- 为新功能编写测试 / Write tests for new features
- 保持函数简短和专注 / Keep functions short and focused
- 添加必要的注释 / Add meaningful comments
- 处理所有错误 / Handle all errors
- 避免不必要的复杂性 / Avoid unnecessary complexity

## 测试 / Testing

### 运行测试 / Running Tests

```bash
# 运行所有测试 / Run all tests
make test

# 运行带覆盖率的测试 / Run tests with coverage
make test-coverage

# 查看覆盖率报告 / View coverage report
go tool cover -html=coverage.html
```

### 编写测试 / Writing Tests

- 为所有公共函数编写单元测试 / Write unit tests for all public functions
- 使用表驱动测试 / Use table-driven tests
- Mock 外部依赖 / Mock external dependencies
- 测试错误情况 / Test error cases

**示例 / Example:**

```go
func TestGenerateCommitMessage(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {
            name:    "valid diff",
            input:   "diff content",
            want:    "feat: add feature",
            wantErr: false,
        },
        // More test cases...
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

## 发布流程 / Release Process

维护者会定期创建新版本。发布流程：

Maintainers will create new releases regularly. Release process:

1. 更新版本号 / Update version number
2. 更新 CHANGELOG / Update CHANGELOG
3. 创建 Git tag / Create Git tag
4. 推送 tag 触发自动发布 / Push tag to trigger automatic release

```bash
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

## 获取帮助 / Getting Help

如果你在贡献过程中遇到问题：

If you need help during the contribution process:

- 💬 在 [Discussions](https://github.com/lemon956/aicommit/discussions) 提问
- 📧 通过 issue 联系维护者
- 📖 查看现有的 issues 和 PRs

## 许可证 / License

通过提交 PR，你同意你的贡献将在与本项目相同的许可证下发布。

By submitting a PR, you agree that your contributions will be licensed under the same license as this project.

---

再次感谢你的贡献！🎉

Thank you again for your contributions! 🎉

