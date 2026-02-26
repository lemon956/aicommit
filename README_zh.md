# aicommit

[English](README.md) | [中文](README_zh.md)

[![CI](https://github.com/lemon956/aicommit/workflows/CI/badge.svg)](https://github.com/lemon956/aicommit/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/lemon956/aicommit)](https://goreportcard.com/report/github.com/lemon956/aicommit)
[![License](https://img.shields.io/github/license/lemon956/aicommit)](LICENSE)
[![Release](https://img.shields.io/github/v/release/lemon956/aicommit)](https://github.com/lemon956/aicommit/releases)
[![Go Version](https://img.shields.io/github/go-mod/go-version/lemon956/aicommit)](go.mod)

AI 驱动的 git commit 消息生成器，使用多种 AI 模型根据暂存的更改生成有意义的提交消息。

## 特性

- 🤖 **多 AI 模型支持**：Claude、OpenAI 和 DeepSeek
- ⚙️ **可配置**：通过 YAML 文件或环境变量轻松配置
- 🎯 **符合 Git 规范**：生成符合 `gitcommit(5)` 建议的提交消息
- 🔒 **安全**：API 密钥可存储在环境变量中
- 🧪 **完善的测试**：全面的测试覆盖
- 🚀 **简单快速**：无冗余，开箱即用

## 安装

```bash
go install github.com/aicommit/aicommit/cmd/aicommit@latest
```

或从源码构建：

```bash
git clone https://github.com/aicommit/aicommit.git
cd aicommit
go build -o aicommit cmd/aicommit/main.go
```

## 快速开始

1. **初始化配置**：
```bash
aicommit config init
```

2. **添加 API 密钥**（选择一种方式）：
   - 编辑 `~/.config/aicommit/aicommit.yaml`
   - 设置环境变量：`export AICOMMIT_CLAUDE_API_KEY=your-key`

3. **暂存更改**：
```bash
git add .
```

4. **生成提交消息**：
```bash
aicommit
```

## 配置

### 配置文件

在 `~/.config/aicommit/aicommit.yaml` 创建配置文件：

```yaml
# 使用的 AI 模型
model: claude-3-sonnet-20240229

# 提供商：claude、openai 或 deepseek
provider: claude

# API 密钥（也可使用环境变量）
api_keys:
  claude: "your-claude-api-key"
  openai: "your-openai-api-key"
  deepseek: "your-deepseek-api-key"
```

### 环境变量

也可以使用环境变量配置 aicommit：

```bash
export AICOMMIT_PROVIDER=claude
export AICOMMIT_MODEL=claude-3-sonnet-20240229
export AICOMMIT_CLAUDE_API_KEY=your-key
export AICOMMIT_OPENAI_API_KEY=your-key
export AICOMMIT_DEEPSEEK_API_KEY=your-key
```

### 支持的模型

#### Claude 模型
- `claude-3-sonnet-20240229`（推荐）
- `claude-3-opus-20240229`
- `claude-3-haiku-20240307`

#### OpenAI 模型
- `gpt-4`
- `gpt-3.5-turbo`

#### DeepSeek 模型
- `deepseek-chat`

## 使用方法

### 基本用法

```bash
# 暂存更改
git add .

# 使用 AI 生成并提交
aicommit

# 预览提交消息而不实际提交
aicommit --dry-run
```

### 高级用法

```bash
# 使用环境变量（优先于配置文件）
export AICOMMIT_PROVIDER=openai
export AICOMMIT_OPENAI_API_KEY=your-key
aicommit
```

## 提交消息格式

aicommit 生成的提交消息遵循 Git 官方 `gitcommit(5)` 的建议：

```
<主题行>

<正文>（可选）

示例：
Add JWT auth to CLI login

Add editor support for reviewing commit message

This lets users edit the generated message before committing and reduces
incorrect commits caused by prompt misunderstandings.
```

## API 密钥设置

### Claude (Anthropic)
1. 在 [Anthropic](https://www.anthropic.com/) 注册
2. 从控制台获取 API 密钥
3. 在配置中设置或使用 `AICOMMIT_CLAUDE_API_KEY`

### OpenAI
1. 在 [OpenAI](https://openai.com/) 注册
2. 从控制台获取 API 密钥
3. 在配置中设置或使用 `AICOMMIT_OPENAI_API_KEY`

### DeepSeek
1. 在 [DeepSeek](https://deepseek.com/) 注册
2. 从控制台获取 API 密钥
3. 在配置中设置或使用 `AICOMMIT_DEEPSEEK_API_KEY`

## 开发

### 项目结构

```
aicommit/
├── cmd/aicommit/        # CLI 应用程序
├── internal/
│   ├── config/          # 配置管理
│   ├── git/             # Git 操作
│   └── model/           # AI 模型提供商
├── pkg/
│   ├── prompt/          # 提交消息工具
│   └── validator/       # 验证工具
└── go.mod
```

### 运行测试

```bash
go test ./...
```

### 构建

```bash
go build -o aicommit cmd/aicommit/main.go
```

## 贡献

我们欢迎各种形式的贡献！请查看 [贡献指南](.github/CONTRIBUTING.md) 了解如何开始。

**快速开始：**

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'feat: add amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

## CI/CD

本项目使用 GitHub Actions 进行持续集成和部署：

- **CI 工作流**：自动运行测试、代码检查和构建
- **Release 工作流**：创建 tag 时自动构建多平台二进制文件并发布
- **Dependabot**：自动更新依赖

查看 [.github/workflows](.github/workflows) 目录了解详情。

## 故障排查

遇到问题？查看我们的 [故障排查指南](TROUBLESHOOTING.md) 获取常见问题的解决方案。

常见问题：
- 🔒 API Key 泄露和清理
- 🔄 Gist 同步错误
- 🧪 测试失败
- 🔨 构建问题

## 许可证

本项目采用 MIT 许可证 - 详情请参阅 [LICENSE](LICENSE) 文件。

## 致谢

- [Cobra](https://github.com/spf13/cobra) - CLI 框架
- [Viper](https://github.com/spf13/viper) - 配置管理
- Git 官方文档（`gitcommit(5)`）- 提交消息建议
