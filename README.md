# aicommit

[![CI](https://github.com/lemon956/aicommit/workflows/CI/badge.svg)](https://github.com/lemon956/aicommit/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/lemon956/aicommit)](https://goreportcard.com/report/github.com/lemon956/aicommit)
[![License](https://img.shields.io/github/license/lemon956/aicommit)](LICENSE)
[![Release](https://img.shields.io/github/v/release/lemon956/aicommit)](https://github.com/lemon956/aicommit/releases)
[![Go Version](https://img.shields.io/github/go-mod/go-version/lemon956/aicommit)](go.mod)

AI-powered git commit message generator that uses various AI models to create meaningful commit messages based on your staged changes.

## Features

- 🤖 **Multiple AI Model Support**: Claude, OpenAI, and DeepSeek
- ⚙️ **Configurable**: Easy configuration via YAML file or environment variables
- 🎯 **Git Standards Compliant**: Generates commit messages following conventional commits format
- 🔒 **Secure**: API keys can be stored in environment variables
- 🧪 **Well Tested**: Comprehensive test coverage
- 🚀 **Simple & Fast**: No bloat, just works

## Installation

```bash
go install github.com/aicommit/aicommit/cmd/aicommit@latest
```

Or build from source:

```bash
git clone https://github.com/aicommit/aicommit.git
cd aicommit
go build -o aicommit cmd/aicommit/main.go
```

## Quick Start

1. **Initialize configuration**:
```bash
aicommit config init
```

2. **Add your API key** (choose one):
   - Edit `~/.config/aicommit/aicommit.yaml`
   - Set environment variable: `export AICOMMIT_CLAUDE_API_KEY=your-key`

3. **Stage your changes**:
```bash
git add .
```

4. **Generate commit message**:
```bash
aicommit
```

## Configuration

### Configuration File

Create a configuration file at `~/.config/aicommit/aicommit.yaml`:

```yaml
# AI model to use
model: claude-3-sonnet-20240229

# Provider: claude, openai, or deepseek
provider: claude

# API keys (alternatively use environment variables)
api_keys:
  claude: "your-claude-api-key"
  openai: "your-openai-api-key"
  deepseek: "your-deepseek-api-key"
```

### Environment Variables

You can also configure aicommit using environment variables:

```bash
export AICOMMIT_PROVIDER=claude
export AICOMMIT_MODEL=claude-3-sonnet-20240229
export AICOMMIT_CLAUDE_API_KEY=your-key
export AICOMMIT_OPENAI_API_KEY=your-key
export AICOMMIT_DEEPSEEK_API_KEY=your-key
```

### Supported Models

#### Claude Models
- `claude-3-sonnet-20240229` (recommended)
- `claude-3-opus-20240229`
- `claude-3-haiku-20240307`

#### OpenAI Models
- `gpt-4`
- `gpt-3.5-turbo`

#### DeepSeek Models
- `deepseek-chat`

## Usage

### Basic Usage

```bash
# Stage your changes
git add .

# Generate and commit with AI
aicommit

# Preview the commit message without committing
aicommit --dry-run
```

### Advanced Usage

```bash
# Use specific config file
aicommit --config /path/to/config.yaml

# Use environment variables (overrides config file)
export AICOMMIT_PROVIDER=openai
export AICOMMIT_OPENAI_API_KEY=your-key
aicommit
```

## Commit Message Format

aicommit generates commit messages following the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>: <description>

Examples:
feat: add user authentication
fix: resolve login issue
docs: update README with installation instructions
style: format code with gofmt
refactor: simplify error handling logic
test: add unit tests for user service
chore: update dependencies
```

Valid types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`, `perf`, `build`, `ci`

## API Key Setup

### Claude (Anthropic)
1. Sign up at [Anthropic](https://www.anthropic.com/)
2. Get your API key from the dashboard
3. Set it in config or use `AICOMMIT_CLAUDE_API_KEY`

### OpenAI
1. Sign up at [OpenAI](https://openai.com/)
2. Get your API key from the dashboard
3. Set it in config or use `AICOMMIT_OPENAI_API_KEY`

### DeepSeek
1. Sign up at [DeepSeek](https://deepseek.com/)
2. Get your API key from the dashboard
3. Set it in config or use `AICOMMIT_DEEPSEEK_API_KEY`

## Development

### Project Structure

```
aicommit/
├── cmd/aicommit/        # CLI application
├── internal/
│   ├── config/          # Configuration management
│   ├── git/             # Git operations
│   └── model/           # AI model providers
├── pkg/
│   ├── prompt/          # Commit message utilities
│   └── validator/       # Validation utilities
└── go.mod
```

### Running Tests

```bash
go test ./...
```

### Building

```bash
go build -o aicommit cmd/aicommit/main.go
```

## Contributing

我们欢迎各种形式的贡献！请查看 [贡献指南](.github/CONTRIBUTING.md) 了解如何开始。

We welcome contributions of all kinds! Please see our [Contributing Guide](.github/CONTRIBUTING.md) to get started.

**快速开始 / Quick Start:**

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## CI/CD

本项目使用 GitHub Actions 进行持续集成和部署：

This project uses GitHub Actions for continuous integration and deployment:

- **CI Workflow**: 自动运行测试、代码检查和构建 / Automatically runs tests, linting, and builds
- **Release Workflow**: 创建 tag 时自动构建多平台二进制文件并发布 / Automatically builds multi-platform binaries and creates releases when tags are pushed
- **Dependabot**: 自动更新依赖 / Automatically updates dependencies

查看 [.github/workflows](.github/workflows) 目录了解详情。

See the [.github/workflows](.github/workflows) directory for details.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Viper](https://github.com/spf13/viper) - Configuration management
- [Conventional Commits](https://www.conventionalcommits.org/) - Commit message specification