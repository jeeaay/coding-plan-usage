# coding-plan-usage

阿里云 Coding Plan 余量查询命令行工具。

## 快速开始

macOS / Linux：

```bash
./coding-plan-usage
```

Windows：

```powershell
.\coding-plan-usage.exe
```

## 运行结果

- 未登录：程序会打开阿里云登录页，保存截图到当前目录 `aliyu-login.png`，按提示扫码后再次运行。

示例输出：

```text
Already logged in: false
Entered login page: true
请使用阿里云 App 扫码完成登录后，再次执行此程序以查询用量。
Login screenshot: /opt/coding-plan-usage/aliyu-login.png
Scan completed: false
Command finished successfully
```

- 已登录：程序会输出用量 JSON（近5小时、近一周、近一月）并自动关闭浏览器。

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

## 可选配置

将 `.env.example` 复制为 `.env` 并按需修改：

```bash
cp .env.example .env
```
