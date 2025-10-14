# 合并 Dependabot PRs 指南

## 📋 待合并的 PR 列表

根据 https://github.com/lemon956/aicommit/pulls，有以下 8 个 Dependabot PR：

| PR # | 标题 | 类型 |
|------|------|------|
| #8 | bump golangci/golangci-lint-action from 4 to 8 | CI 工具 |
| #7 | bump actions/checkout from 4 to 5 | CI 工具 |
| #6 | bump actions/setup-go from 5 to 6 | CI 工具 |
| #5 | bump codecov/codecov-action from 4 to 5 | CI 工具 |
| #4 | bump github.com/spf13/cobra from 1.8.0 to 1.10.1 | Go 依赖 |
| #3 | bump github.com/spf13/viper from 1.18.2 to 1.21.0 | Go 依赖 |
| #2 | bump github.com/stretchr/testify from 1.8.4 to 1.11.1 | Go 依赖 |
| #1 | bump softprops/action-gh-release from 1 to 2 | CI 工具 |

## 🚀 方法 1: 通过 GitHub 网页合并（推荐）

### 步骤：

1. 访问 https://github.com/lemon956/aicommit/pulls

2. 对每个 PR：
   - 点击 PR 标题
   - 等待 CI 检查通过（绿色 ✓）
   - 点击 "Squash and merge" 按钮
   - 确认合并

3. **建议合并顺序**（从简单到复杂）：

   **第一批 - GitHub Actions 更新**（无代码影响）：
   - #8 golangci-lint-action (4 → 8)
   - #7 actions/checkout (4 → 5)
   - #6 actions/setup-go (5 → 6)
   - #5 codecov-action (4 → 5)
   - #1 softprops/action-gh-release (1 → 2)

   **第二批 - Go 依赖更新**（可能影响代码）：
   - #2 testify (1.8.4 → 1.11.1) - 测试库
   - #4 cobra (1.8.0 → 1.10.1) - CLI 框架
   - #3 viper (1.18.2 → 1.21.0) - 配置管理

## 🛠️ 方法 2: 使用 GitHub CLI（如果已安装）

```bash
# 安装 GitHub CLI
# Fedora/RHEL: sudo dnf install gh
# Ubuntu/Debian: sudo apt install gh
# macOS: brew install gh

# 登录
gh auth login

# 批量启用自动合并（CI 通过后自动合并）
for i in {1..8}; do
  gh pr merge $i --auto --squash --delete-branch
  sleep 1
done
```

## 🔄 方法 3: 手动使用 Git（不推荐）

```bash
# 对每个 PR 执行以下操作
git fetch origin pull/1/head:dependabot-pr-1
git checkout main
git merge --squash dependabot-pr-1
git commit -m "chore(ci)(deps): bump softprops/action-gh-release from 1 to 2"
git push origin main
git branch -D dependabot-pr-1
```

## ⚡ 快速批量合并脚本

如果 CI 已经全部通过，可以使用：

```bash
# 安装 GitHub CLI 后运行
gh pr list --state open --json number --jq '.[].number' | while read pr; do
  echo "Merging PR #$pr"
  gh pr merge "$pr" --squash --delete-branch --admin || echo "Failed to merge PR #$pr"
  sleep 2
done
```

## 📊 合并后验证

合并所有 PR 后，验证项目状态：

```bash
# 拉取最新代码
git pull origin main

# 更新依赖
go mod download
go mod tidy

# 运行测试
make test

# 运行 lint
make lint

# 构建项目
make build

# 验证二进制文件
./build/aicommit --version
```

## 🎯 预期结果

合并后，您的项目将：
- ✅ 使用最新版本的 GitHub Actions
- ✅ 使用最新版本的 Go 依赖
- ✅ 提高安全性和性能
- ✅ 修复已知的 bug

## ⚠️ 注意事项

1. **CI 必须通过**：确保每个 PR 的 CI 检查都通过后再合并
2. **逐个合并**：建议逐个合并，而不是一次性全部合并
3. **测试**：合并后运行完整测试套件
4. **Breaking Changes**：虽然这些都是小版本更新，但仍需注意 viper 从 1.18 到 1.21 的变化

## 🐛 如果遇到问题

### CI 失败

如果某个 PR 的 CI 失败：
1. 查看失败日志
2. 本地修复问题
3. 推送到 main 分支
4. Dependabot 会自动 rebase 它的 PR

### 合并冲突

如果有合并冲突：
1. 关闭冲突的 PR
2. 手动更新依赖：
   ```bash
   go get github.com/spf13/cobra@latest
   go mod tidy
   ```
3. 提交并推送

## 🔗 相关链接

- [GitHub PR 页面](https://github.com/lemon956/aicommit/pulls)
- [GitHub Actions 文档](https://docs.github.com/en/actions)
- [Dependabot 文档](https://docs.github.com/en/code-security/dependabot)

---

**提示**：所有 lint 错误已在 commit `930f9b1` 中修复，现在 CI 应该能通过了！

