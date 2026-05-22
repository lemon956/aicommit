# 故障排查指南 / Troubleshooting Guide

## 🔒 API Key 泄露问题

### 问题：GitHub 阻止推送，提示包含 secrets

**错误信息**：
```
remote: error: GH013: Repository rule violations found
remote: - Push cannot contain secrets
```

**解决方案**：

1. **立即撤销泄露的 API Key**
   - OpenAI: https://platform.openai.com/api-keys
   - Claude: https://console.anthropic.com/settings/keys
   - DeepSeek: https://platform.deepseek.com/api_keys

2. **清理 Git 历史**：
   ```bash
   # 创建新的干净历史
   rm -rf .git
   git init
   git add .
   git commit -m "feat: initialize aicommit"
   git branch -M main
   git remote add origin git@github.com:username/aicommit.git
   git push -f origin main
   ```

3. **使用环境变量存储 API Key**：
   ```bash
   # 添加到 ~/.bashrc 或 ~/.zshrc
   export AICOMMIT_OPENAI_API_KEY="your-new-key"
   export AICOMMIT_CLAUDE_API_KEY="your-claude-key"
   ```

## 🔄 Gist 同步错误

### 问题：Failed to replace Gist: Error: No timestamp.json found in Gist

这个错误通常与 IDE 的设置同步功能相关。

**可能的原因**：
1. Cursor/VS Code 的设置同步功能
2. Settings Sync 插件
3. 损坏的 Gist 配置

**解决方案**：

1. **禁用 Cursor 同步功能**：
   - 打开 Cursor 设置
   - 搜索 "sync"
   - 禁用设置同步

2. **检查并删除损坏的 Gist**：
   - 访问 https://gist.github.com/
   - 找到相关的设置同步 Gist
   - 删除或重新创建

3. **清理本地同步数据**：
   ```bash
   # Linux/macOS
   rm -rf ~/.cursor-tutor
   rm -rf ~/.cursor/sync
   
   # Windows
   # 删除 %APPDATA%\.cursor\sync
   ```

4. **重置 Cursor 同步**：
   - Cursor → Settings → Settings Sync
   - Turn Off Settings Sync
   - 重启 Cursor
   - 重新启用（如果需要）

## 🧪 测试失败问题

### 问题：TestGit_GetDiff 失败

**错误信息**：
```
Error: no staged changes found
```

**原因**：测试运行时没有暂存的文件。

**解决方案**：

修改测试以处理空 diff 情况：

```go
func TestGit_GetDiff(t *testing.T) {
    g := New(".")
    diff, err := g.GetDiff()
    
    // 如果没有暂存的更改，这是正常的
    if err != nil && err.Error() == "no staged changes found" {
        t.Skip("No staged changes, skipping test")
        return
    }
    
    require.NoError(t, err)
    assert.NotEmpty(t, diff)
}
```

或者在运行测试前创建临时文件：

```bash
# 运行测试
echo "test" > test_file.txt
git add test_file.txt
make test
git reset HEAD test_file.txt
rm test_file.txt
```

## 🔨 构建问题

### 问题：undefined: SetDefaultPrompt

**错误信息**：
```
cmd/aicommit/main.go:68:4: undefined: SetDefaultPrompt
```

**原因**：Makefile 只编译了 main.go，没有包含其他源文件。

**解决方案**：

修改 Makefile：
```makefile
BINARY_PATH=./cmd/aicommit  # 而不是 ./cmd/aicommit/main.go
```

## 📦 依赖问题

### 问题：go.mod 文件过时

**解决方案**：
```bash
go mod tidy
go mod download
go mod verify
```

## 🌐 网络问题

### 问题：无法访问 AI API

**检查清单**：
1. ✅ 检查网络连接
2. ✅ 检查 API Key 是否有效
3. ✅ 检查 API 配额是否用尽
4. ✅ 检查是否需要代理
   ```bash
   export HTTP_PROXY=http://proxy:port
   export HTTPS_PROXY=http://proxy:port
   ```

## 🐛 常见错误

### 1. invalid commit message format

**确保格式**：
```
subject line

body (optional)
```

**常见要求**：
- 主题行（第 1 行）建议保持简洁
- 如果有正文，主题行与正文之间必须有一行空行
- 正文建议保持易读（必要时手动换行）
- 可选 trailers（如 `Co-authored-by: ...`、`Signed-off-by: ...`）

## 📝 日志调试

启用详细日志：
```bash
# 设置环境变量
export AICOMMIT_DEBUG=true

# 运行命令
aicommit --dry-run
```

## 🆘 获取帮助

如果以上方案都无法解决问题：

1. 🔍 搜索已有 issues: https://github.com/lemon956/aicommit/issues
2. 💬 在 Discussions 提问: https://github.com/lemon956/aicommit/discussions
3. 🐛 创建新 issue: https://github.com/lemon956/aicommit/issues/new

请提供：
- 操作系统和版本
- Go 版本
- aicommit 版本
- 完整的错误信息
- 复现步骤

---

**提示**：保持 aicommit 更新到最新版本可以避免很多已知问题。

```bash
# 更新到最新版本
go install github.com/lemon956/aicommit/cmd/aicommit@latest
```
