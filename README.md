# coding-plan-usage

一个用 Go + agent-browser 实现的阿里云 Coding Plan 余量查询工具。

核心目标：**一条命令快速拿到余量 JSON**。

## 快速使用

### 0）交给AI吧

本项目已经内置skills，复制`skills`下的`coding-plan-usage`到你的agent的技能目录下，然后告诉agent你需要查询 Coding Plan 余量。

### 1）直接运行

```bash
go run .
```

### 2）看结果

- **未登录**：会自动打开阿里云首页并进入登录页，保存截图到当前目录 `aliyu-login.png`，终端提示你扫码；扫码后再次运行即可。

示例输出：

```text
Already logged in: false
Entered login page: true
请使用阿里云 App 扫码完成登录后，再次执行此程序以查询用量。
Login screenshot: /opt/coding-plan-usage/aliyu-login.png
Scan completed: false
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

## 工作流程说明

每次执行程序时：
- 打开阿里云首页
- 判断是否已登录
- 已登录：跳转到 Coding Plan 详情页并解析“近5小时/近一周/近一月”用量与重置时间
- 未登录：点击登录入口，截图保存登录页，等待你扫码后再次执行

## GitHub 多平台编译

GitHub Actions 工作流：
- [build-multi-platform.yml](https://github.com/jeeaay/coding-plan-usage/blob/main/.github/workflows/build-multi-platform.yml)

触发方式：
- 手动触发：`Actions -> build-multi-platform -> Run workflow`
- 打标签触发：推送 `v*` 标签（如 `v1.0.0`）

打标签：

```bash
git tag v1.0.0
git push origin v1.0.0
```
