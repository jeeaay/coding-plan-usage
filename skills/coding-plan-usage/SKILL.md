---
name: "coding-plan-usage"
description: "Queries the remaining hours of Alibaba Cloud Coding Plan using a command-line tool. Invoke when user asks for Coding Plan usage."
---

# Coding Plan Usage Helper

用于查询阿里云 Coding Plan 余量的命令行工具。

## 依赖项目

基于agent-browser，如果还没有安装，请先安装它。

- [agent-browser](https://github.com/vercel-labs/agent-browser)

## 何时使用

在以下场景主动调用：
- 用户希望“查询阿里云 Coding Plan余量”

## Release 地址

- https://github.com/jeeaay/coding-plan-usage/releases

## 平台与产物映射

- macOS Intel: `coding-plan-usage-darwin-amd64.tar.gz`
- macOS Apple Silicon: `coding-plan-usage-darwin-arm64.tar.gz`
- Linux x86_64: `coding-plan-usage-linux-amd64.tar.gz`
- Linux ARM64: `coding-plan-usage-linux-arm64.tar.gz`
- Windows x86_64: `coding-plan-usage-windows-amd64.zip`
- Windows ARM64: `coding-plan-usage-windows-arm64.zip`

## 标准执行流程

1. 识别用户平台与架构（darwin/linux/windows + amd64/arm64）
2. 选择对应 release 产物
3. 下载并解压
4. 复制 `.env.example` 为 `.env`（如需）
5. 运行二进制并解释输出

## macOS / Linux 示例

```bash
# 1) 下载（示例：macOS arm64）
curl -fL -o coding-plan-usage-darwin-arm64.tar.gz \
  https://github.com/jeeaay/coding-plan-usage/releases/latest/download/coding-plan-usage-darwin-arm64.tar.gz

# 2) 解压
tar -xzf coding-plan-usage-darwin-arm64.tar.gz
cd coding-plan-usage-darwin-arm64-bundle

# 3) 初始化配置（可选）
cp .env.example .env

# 4) 运行
chmod +x coding-plan-usage
./coding-plan-usage
```

## Windows 示例（PowerShell）

```powershell
# 1) 下载（示例：Windows amd64）
Invoke-WebRequest `
  -Uri "https://github.com/jeeaay/coding-plan-usage/releases/latest/download/coding-plan-usage-windows-amd64.zip" `
  -OutFile "coding-plan-usage-windows-amd64.zip"

# 2) 解压
Expand-Archive .\coding-plan-usage-windows-amd64.zip -DestinationPath .
Set-Location .\coding-plan-usage-windows-amd64-bundle

# 3) 初始化配置（可选）
Copy-Item .env.example .env

# 4) 运行
.\coding-plan-usage.exe
```

## 输出解释规则


- **未登录**：会自动打开阿里云首页并进入登录页，保存截图到当前目录 `aliyu-login.png`，终端提示你扫码；扫码后再次运行即可。如果频道允许发送图片 你可以直接发给用户，否则可以帮用户打开图片。

示例输出：

```text
Already logged in: false
Entered login page: true
请使用阿里云 App 扫码完成登录后，再次执行此程序以查询用量。
Login screenshot: /opt/coding-plan-usage/aliyu-login.png
Scan completed: false
Command finished successfully
```

- **已登录**：会自动进入 Coding Plan 页面并输出余量 JSON。

示例输出：

```json
{
  "hours5": {
    "usage": "0%",
    "resetTime": "2026-03-14 18:27:45"
  },
  "week": {
    "usage": "27%",
    "resetTime": "2026-03-16 00:00:00"
  },
  "month": {
    "usage": "15%",
    "resetTime": "2026-04-09 00:00:00"
  }
}
```

成功读取到用量后，程序会自动关闭浏览器会话。
