# GitHub Workflows 使用指南

## 📋 Workflows 概览

本项目包含以下 GitHub Actions workflows：

### 1. CI Workflow (`.github/workflows/ci.yml`)

**触发条件：**
- 推送到 `master`、`main` 或 `develop` 分支
- 创建针对这些分支的 Pull Request

**功能：**
- ✅ 在多个 Go 版本 (1.21, 1.22, 1.23) 上运行测试
- ✅ 运行 golangci-lint 代码检查
- ✅ 生成测试覆盖率报告
- ✅ 构建二进制文件
- ✅ 上传构建产物

**查看状态：**
```
https://github.com/lemon956/aicommit/actions?query=workflow%3ACI
```

### 2. Release Workflow (`.github/workflows/release.yml`)

**触发条件：**
- 推送以 `v` 开头的 tag（如 `v1.0.0`）

**功能：**
- 🚀 为多个平台构建二进制文件：
  - Linux (AMD64, ARM64)
  - macOS (AMD64, ARM64)
  - Windows (AMD64, ARM64)
- 📦 创建压缩包和校验和
- 📝 自动生成 release notes
- 🎉 创建 GitHub Release

**如何发布新版本：**

```bash
# 1. 更新版本号和 CHANGELOG
vim CHANGELOG.md

# 2. 提交更改
git add .
git commit -m "chore: prepare for v1.0.0 release"

# 3. 创建并推送 tag
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0

# 4. GitHub Actions 会自动：
#    - 构建多平台二进制文件
#    - 创建 GitHub Release
#    - 上传所有构建产物
```

## 🔧 本地测试

### 测试 CI workflow

```bash
# 运行测试
make test

# 运行 linter
make lint

# 构建
make build
```

### 测试多平台构建

```bash
# Linux AMD64
GOOS=linux GOARCH=amd64 go build -o build/aicommit-linux-amd64 ./cmd/aicommit

# macOS ARM64
GOOS=darwin GOARCH=arm64 go build -o build/aicommit-darwin-arm64 ./cmd/aicommit

# Windows AMD64
GOOS=windows GOARCH=amd64 go build -o build/aicommit-windows-amd64.exe ./cmd/aicommit
```

## 📊 状态徽章

在 README 中添加以下徽章：

```markdown
[![CI](https://github.com/lemon956/aicommit/workflows/CI/badge.svg)](https://github.com/lemon956/aicommit/actions)
[![Release](https://img.shields.io/github/v/release/lemon956/aicommit)](https://github.com/lemon956/aicommit/releases)
```

## 🤖 Dependabot

Dependabot 配置文件位于 `.github/dependabot.yml`，会自动：

- 每周检查 Go 模块更新
- 每周检查 GitHub Actions 更新
- 自动创建 PR 进行依赖更新

## 🔐 Secrets 配置

如果需要额外的 secrets（如发布到其他平台），在 GitHub 仓库设置中添加：

```
Settings → Secrets and variables → Actions → New repository secret
```

常用 secrets：
- `GITHUB_TOKEN` - 自动提供，用于创建 releases
- `CODECOV_TOKEN` - Codecov 集成（可选）

## 📝 自定义 Workflows

### 修改 Go 版本

编辑 `.github/workflows/ci.yml`：

```yaml
strategy:
  matrix:
    go-version: ['1.21', '1.22', '1.23']  # 修改这里
```

### 添加新的构建平台

编辑 `.github/workflows/release.yml`，添加：

```yaml
# FreeBSD AMD64
GOOS=freebsd GOARCH=amd64 go build -o build/aicommit-freebsd-amd64 ./cmd/aicommit
```

### 修改触发分支

编辑 workflow 文件的 `on` 部分：

```yaml
on:
  push:
    branches: [ master, main, develop, staging ]  # 添加更多分支
```

## 🐛 故障排查

### CI 失败

1. **测试失败**：检查测试日志，修复失败的测试
2. **Lint 失败**：运行 `make lint` 查看具体错误
3. **构建失败**：检查依赖是否正确

### Release 失败

1. **权限问题**：确保 GITHUB_TOKEN 有足够权限
2. **Tag 格式错误**：确保 tag 以 `v` 开头
3. **构建错误**：在本地测试多平台构建

## 📚 相关资源

- [GitHub Actions 文档](https://docs.github.com/en/actions)
- [Go Actions](https://github.com/actions/setup-go)
- [golangci-lint Action](https://github.com/golangci/golangci-lint-action)
- [Dependabot 文档](https://docs.github.com/en/code-security/dependabot)

## 💡 最佳实践

1. ✅ 在本地测试后再推送
2. ✅ 保持 workflows 简洁高效
3. ✅ 定期更新 GitHub Actions 版本
4. ✅ 使用缓存加速构建
5. ✅ 为重要步骤添加清晰的注释
6. ✅ 监控 workflow 运行时间和成本

---

如有问题，请在 [Discussions](https://github.com/lemon956/aicommit/discussions) 中提问。

