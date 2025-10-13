# 安全政策 / Security Policy

## 支持的版本 / Supported Versions

当前支持以下版本的安全更新：

Currently supported versions with security updates:

| Version | Supported          |
| ------- | ------------------ |
| 1.x.x   | :white_check_mark: |
| < 1.0   | :x:                |

## 报告漏洞 / Reporting a Vulnerability

我们非常重视安全问题。如果你发现了安全漏洞，请**不要**在公开 issue 中报告。

We take security issues seriously. If you discover a security vulnerability, please **DO NOT** report it in a public issue.

### 如何报告 / How to Report

请通过以下方式之一报告安全漏洞：

Please report security vulnerabilities through one of the following methods:

1. **GitHub Security Advisory** (推荐 / Recommended)
   - 访问 / Visit: https://github.com/lemon956/aicommit/security/advisories
   - 点击 "Report a vulnerability" / Click "Report a vulnerability"

2. **邮件 / Email**
   - 发送至 / Send to: 15230727732xlm@gmail.com
   - 主题 / Subject: `[SECURITY] aicommit vulnerability report`

### 报告应包含 / Report Should Include

请在报告中包含以下信息：

Please include the following information in your report:

- 漏洞类型 / Type of vulnerability
- 受影响的版本 / Affected versions
- 复现步骤 / Steps to reproduce
- 潜在影响 / Potential impact
- 建议的修复方案（如有）/ Suggested fix (if any)

### 响应时间 / Response Timeline

- **确认收到**: 24-48 小时内 / **Acknowledgment**: Within 24-48 hours
- **初步评估**: 7 天内 / **Initial Assessment**: Within 7 days
- **修复发布**: 根据严重程度，通常 30 天内 / **Fix Release**: Based on severity, typically within 30 days

## 安全最佳实践 / Security Best Practices

### API 密钥管理 / API Key Management

1. **永远不要**将 API 密钥提交到版本控制 / **Never** commit API keys to version control
2. 使用环境变量存储敏感信息 / Use environment variables for sensitive information
3. 定期轮换 API 密钥 / Rotate API keys regularly
4. 限制 API 密钥权限 / Limit API key permissions

### 配置文件 / Configuration Files

```bash
# 确保配置文件权限正确
# Ensure proper file permissions
chmod 600 ~/.config/aicommit/aicommit.yaml
```

### 依赖更新 / Dependency Updates

我们使用 Dependabot 自动更新依赖。请及时更新到最新版本。

We use Dependabot to automatically update dependencies. Please update to the latest version promptly.

## 安全功能 / Security Features

- ✅ API 密钥通过环境变量存储 / API keys stored via environment variables
- ✅ 不记录敏感信息 / No logging of sensitive information
- ✅ 最小权限原则 / Principle of least privilege
- ✅ 定期安全审计 / Regular security audits
- ✅ 依赖漏洞扫描 / Dependency vulnerability scanning

## 已知问题 / Known Issues

目前没有已知的安全问题。

No known security issues at this time.

## 漏洞披露政策 / Vulnerability Disclosure Policy

- 我们遵循**负责任的披露**原则 / We follow **responsible disclosure** principles
- 修复发布后，将公开披露漏洞详情 / Vulnerability details will be disclosed after fix is released
- 安全研究人员将获得适当的致谢 / Security researchers will receive appropriate credit

## 联系方式 / Contact

如有安全相关问题，请联系：

For security-related questions, please contact:

- Security Email: 15230727732xlm@gmail.com
- GitHub Security: https://github.com/lemon956/aicommit/security

---

感谢你帮助保护 aicommit 及其用户！

Thank you for helping to keep aicommit and its users safe!

