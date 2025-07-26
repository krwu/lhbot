# Lighthouse 实例监视器 (LHBot)

一个轻量级的 Go 机器人，用于监控腾讯云 Lighthouse 实例的可用性，并在实例有货时自动购买。

## 🎯 用途

该机器人持续监控特定腾讯云 Lighthouse 实例类型（尤其是“锐驰”系列）的可用性，并执行以下操作：
- 在实例可用时自动购买。
- 通过企业微信 Webhook 发送实时状态通知。
- 提供定期的监控报告。

## 🚀 功能

- **持续监控**：每 30 秒检查一次实例可用性。
- **自动购买**：在检测到库存时自动创建实例。
- **智能通知**：
  - 实例有货时立即提醒（并 @指定用户）。
  - 所有实例均售罄时，每小时发送一次心跳通知。
  - 购买成功后发送确认提醒。
- **重复购买防护**：防止在服务重启后重复购买相同类型的实例。
- **优雅关闭**：正确处理 SIGINT/SIGTERM 信号，以实现安全终止。
- **跨平台支持**：为 Linux AMD64 和 ARM64 提供预编译的二进制文件。

## 📋 先决条件

- 一个已激活 Lighthouse 服务的腾讯云账户。
- 腾讯云 API 凭证（SecretId 和 SecretKey）。
- 用于通知的企业微信 Webhook URL。
- Go 1.24.4 或更高版本（如果从源代码构建）。

## 🔧 配置

在 `~/.config/lhbot.env` 创建一个配置文件，并包含以下环境变量：

```bash
# 腾讯云凭证
CLIENT_ID="your_tencent_secret_id"
CLIENT_SECRET="your_tencent_secret_key"

# 企业微信通知
WEBHOOK="https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=your_webhook_key"
CHAT_ID="your_chat_id"

# 要监控的实例ID，多个用半角逗号隔开，bundle_rs_nmc_lin_med1_02=2核1G 40元套餐，bundle_rs_nmc_lin_med2_01=2核2G 55元套餐
BUNDLES="bundle_rs_nmc_lin_med1_02,bundle_rs_nmc_lin_med2_01"

# 可选：新实例的 root 密码（默认为 admin@2025）
ROOT_PASSWORD="your_secure_password"

# 可选：提及特定的用户 ID（默认为 @all 通知所有人）
MENTIONED_USERID="your_user_id"
# 是否启用自动购买，1=启用
ENABLE_PURCHASE="0"
```

## 🏗️ 构建

### 从源代码

```bash
# 克隆仓库
git clone <repository-url>
cd lhbot

# 为当前平台构建
go build -o lhbot .

# 为所有支持的平台（Linux AMD64/ARM64）交叉编译
./make.ps1 -Platform all
```

### 预编译的二进制文件

您也可以从 [Releases](https://github.com/krwu/lhbot/releases) 页面下载预编译的二进制文件：

- `lhbot_linux-amd64`
- `lhbot_linux-arm64`

## 🔧 安装与部署

### 1. 安装二进制文件

```bash
# 创建本地 bin 目录
mkdir -p ~/.local/bin

# 将二进制文件复制到本地 bin 目录
cp bin/lhbot_linux-amd64 ~/.local/bin/lhbot
chmod +x ~/.local/bin/lhbot

# 创建配置目录
mkdir -p ~/.config
```

### 2. 作为用户级 systemd 服务运行

#### 服务设置

```bash
# 将 systemd 服务文件复制到用户的配置目录
cp systemd/lhbot.service ~/.config/systemd/user/

# 启用 linger 以确保在您注销后服务能继续运行
loginctl enable-linger $USER

# 重新加载 systemd 用户守护进程
systemctl --user daemon-reload

# 启用服务以在启动时运行
systemctl --user enable lhbot.service

# 立即启动服务
systemctl --user start lhbot.service

# 检查服务状态
systemctl --user status lhbot.service
```

#### 查看日志

```bash
# 查看实时日志
tail -f /var/log/lhbot.log
```

#### 服务管理

```bash
# 停止服务
systemctl --user stop lhbot.service

# 重启服务
systemctl --user restart lhbot.service

# 禁用服务
systemctl --user disable lhbot.service
```

### 3. 手动执行

```bash
# 从配置文件加载环境变量
source ~/.config/lhbot.env

# 直接运行机器人
~/.local/bin/lhbot
```

## 📊 监控与通知

机器人通过企业微信发送中文通知。

### 通知类型

1.  **有货提醒** (立即)
    > ⚠️ **发现可用套餐**
    > - **锐驰-2C1G**: AVAILABLE
    > - **锐驰-2C2G**: SOLD_OUT
    >
    > **通知时间**：2025-07-25 21:15:00
    >
    > *提及指定用户 (例如 @kairee)*

2.  **心跳** (每小时，当所有实例售罄时)
    > ⚙️ **监控服务运行中**
    > - **锐驰-2C1G**: SOLD_OUT
    > - **锐驰-2C2G**: SOLD_OUT
    >
    > **通知时间**：2025-07-25 21:00:00

3.  **购买确认**
    > ✅ **锐驰自动购买成功**
    > - **型号**: 锐驰-2C2G
    >
    > **通知时间**：2025-07-25 21:16:30

### 重复购买防护

机器人创建一个锁文件以防止重复购买，即使服务重启。

- **锁文件**：`~/lhbot-bought.lock`
- **目的**：在重启后保持购买状态。
- **行为**：如果 `~/lhbot-bought.lock` 存在，自动购买功能将被禁用。

## ⚙️ systemd 配置

用户级服务在 `~/.config/systemd/user/lhbot.service` 中定义。

关键指令：
- **EnvironmentFile**：从 `~/.config/lhbot.env` 加载配置。
- **ExecStart**：指定要运行的命令：`~/.local/bin/lhbot`。
- **Restart**：在失败时自动重启服务。
- **Logging**：标准输出和错误由 systemd 日志捕获，可通过 `journalctl` 查看。

## 🔍 故障排除

### 常见问题

**用户注销后服务停止：**
为您的用户启用 linger，以允许服务在后台运行。
```bash
loginctl enable-linger $USER
```

**服务无法启动：**
```bash
# 检查服务状态以获取错误信息
systemctl --user status lhbot.service

# 查看详细日志
journalctl --user -u lhbot.service --no-pager

# 验证二进制文件是否具有执行权限
ls -la ~/.local/bin/lhbot
```

**API 认证错误：**
- 确保 `~/.config/lhbot.env` 中的 `CLIENT_ID` 和 `CLIENT_SECRET` 正确。
- 在腾讯云控制台中确认 API 密钥具有必要的权限。

**Webhook 通知失败：**
- 验证 `WEBHOOK` URL 是否正确且可访问。
- 检查 `CHAT_ID` 配置。

### 调试模式

要以详细日志记录模式运行机器人以进行调试：
```bash
# 加载环境变量
source ~/.config/lhbot.env

# 运行机器人并将输出通过管道传输到日志文件
~/.local/bin/lhbot 2>&1 | tee debug.log
```

## 🛡️ 安全最佳实践

- **切勿提交机密**：不要将 API 凭证或其他机密提交到版本控制。
- **使用环境变量**：将敏感数据存储在环境变量中，而不是代码中。
- **轮换密钥**：定期轮换您的腾讯云 API 密钥。
- **监控日志**：检查日志以发现可疑活动：`journalctl --user -u lhbot.service`。

## 📄 最近更新

### v1.1.0 (2025-07-25)
- **增强通知**：为库存可用性添加了带用户提及的即时提醒。
- **重复购买防护**：实现了基于文件的锁，以避免在重启后重复购买。
- **改进的错误处理**：为 `CreateInstances` API 失败添加了更详细的日志记录。
- **Systemd 服务修复**：记录了持久用户服务所需的 `loginctl enable-linger` 要求。

## 🤝 贡献

1.  Fork 仓库。
2.  创建您的功能分支。
3.  提交您的更改。
4.  确保您的更改经过充分测试。
5.  提交拉取请求。

## 📞 支持

如果您遇到问题：
- 查看故障排除部分。
- 检查系统日志：`journalctl --user -u lhbot.service`。
- 在 GitHub 上提出问题。
