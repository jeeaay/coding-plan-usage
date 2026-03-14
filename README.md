使用agent-browser操作浏览器获取数据，使用方法在

浏览器的使用需要保存会话以便后续使用 始终需要携带--session-name jeayapp 来保持会话
例如：
# 首次登录
agent-browser --session-name jeayapp open https://app.example.com/login
# ... 执行登录操作 ...
agent-browser close  # ← 自动保存！

# 下次使用
agent-browser --session-name jeayapp open https://app.example.com/dashboard
# ← 自动恢复登录状态！