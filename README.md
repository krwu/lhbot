# Lighthouse Instance Monitor (LHBot)

A lightweight Go bot to monitor Tencent Cloud Lighthouse instance availability and automatically purchase them when they're back in stock.

## 🎯 Purpose

This bot continuously monitors the availability of specific Tencent Cloud Lighthouse instance types (especially the "锐驰" series) and performs the following actions:
- Automatically purchases an instance when it becomes available.
- Sends real-time status notifications via WeChat Work webhooks.
- Provides periodic monitoring reports.

## 🚀 Features

- **Continuous Monitoring**: Checks instance availability every 30 seconds.
- **Auto-Purchase**: Automatically creates an instance when stock is detected.
- **Intelligent Notifications**: 
  - Immediate alert when an instance is back in stock (with @user mention).
  - Hourly heartbeat notifications if all instances are sold out.
  - Purchase confirmation alerts.
- **Duplicate Purchase Prevention**: Prevents re-purchasing the same instance type after a service restart.
- **Graceful Shutdown**: Handles SIGINT/SIGTERM signals correctly for safe termination.
- **Cross-Platform Support**: Pre-built binaries for Linux AMD64 and ARM64.

## 📋 Prerequisites

- A Tencent Cloud account with the Lighthouse service activated.
- Tencent Cloud API Credentials (SecretId and SecretKey).
- A WeChat Work webhook URL for notifications.
- Go 1.24.4 or later (if building from source).

## 🔧 Configuration

Create a configuration file at `~/.config/lhbot.env` with the following environment variables:

```bash
# Tencent Cloud Credentials
CLIENT_ID="your_tencent_secret_id"
CLIENT_SECRET="your_tencent_secret_key"

# WeChat Work Notifications
WEBHOOK="https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=your_webhook_key"
CHAT_ID="your_chat_id"

# Optional: Root password for new instances (default: admin@2025)
ROOT_PASSWORD="your_secure_password"

# Optional: Mention a specific user ID (defaults to @all to notify everyone)
MENTIONED_USERID="your_user_id"
```

## 🏗️ Building

### From Source

```bash
# Clone the repository
git clone <repository-url>
cd lhbot

# Build for the current platform
go build -o lhbot .

# Cross-compile for all supported platforms (Linux AMD64/ARM64)
./make.ps1 -Platform all
```

### Pre-built Binaries

You can also download pre-built binaries from the [Releases](https://github.com/krwu/lhbot/releases) page:

- `lhbot_linux-amd64`
- `lhbot_linux-arm64`

## 🔧 Installation & Deployment

### 1. Install the Binary

```bash
# Create a local bin directory
mkdir -p ~/.local/bin

# Copy the binary to the local bin directory
cp bin/lhbot_linux-amd64 ~/.local/bin/lhbot
chmod +x ~/.local/bin/lhbot

# Create the configuration directory
mkdir -p ~/.config
```

### 2. Run as a User-Level systemd Service

#### Service Setup

```bash
# Copy the systemd service file to the user's config directory
cp systemd/lhbot.service ~/.config/systemd/user/

# Enable linger to ensure the service continues running after you log out
loginctl enable-linger $USER

# Reload the systemd user daemon
systemctl --user daemon-reload

# Enable the service to start on boot
systemctl --user enable lhbot.service

# Start the service now
systemctl --user start lhbot.service

# Check the service status
systemctl --user status lhbot.service
```

#### View Logs

```bash
# View real-time logs
journalctl --user -f -u lhbot.service

# View logs from the past hour
journalctl --user -u lhbot.service --since "1 hour ago"
```

#### Service Management

```bash
# Stop the service
systemctl --user stop lhbot.service

# Restart the service
systemctl --user restart lhbot.service

# Disable the service
systemctl --user disable lhbot.service
```

### 3. Manual Execution

```bash
# Load environment variables from your config file
source ~/.config/lhbot.env

# Run the bot directly
~/.local/bin/lhbot
```

## 📊 Monitoring & Notifications

The bot sends notifications in Chinese via WeChat Work.

### Notification Types

1.  **Stock Available** (Immediate)
    > ⚠️ **发现可用套餐**
    > - **锐驰-2C1G**: AVAILABLE
    > - **锐驰-2C2G**: SOLD_OUT
    >
    > **通知时间**：2025-07-25 21:15:00
    >
    > *Mentions the specified user (e.g., @kairee)*

2.  **Heartbeat** (Hourly, when all instances are sold out)
    > ⚙️ **监控服务运行中**
    > - **锐驰-2C1G**: SOLD_OUT
    > - **锐驰-2C2G**: SOLD_OUT
    >
    > **通知时间**：2025-07-25 21:00:00

3.  **Purchase Confirmation**
    > ✅ **锐驰自动购买成功**
    > - **型号**: 锐驰-2C2G
    >
    > **通知时间**：2025-07-25 21:16:30

### Duplicate Purchase Prevention

The bot creates a lock file to prevent duplicate purchases, even if the service restarts.

- **Lock File**: `~/lhbot-bought.lock`
- **Purpose**: To persist the purchase status across restarts.
- **Behavior**: If `~/lhbot-bought.lock` exists, the auto-purchase function will be disabled.

## ⚙️ systemd Configuration

The user-level service is defined in `~/.config/systemd/user/lhbot.service`.

Key directives:
- **EnvironmentFile**: Loads configuration from `~/.config/lhbot.env`.
- **ExecStart**: Specifies the command to run: `~/.local/bin/lhbot`.
- **Restart**: Automatically restarts the service on failure.
- **Logging**: Standard output and error are captured by the systemd journal, viewable with `journalctl`.

## 🔍 Troubleshooting

### Common Issues

**Service stops after user logs out:**
Enable linger for your user to allow services to run in the background.
```bash
loginctl enable-linger $USER
```

**Service fails to start:**
```bash
# Check the service status for errors
systemctl --user status lhbot.service

# View detailed logs
journalctl --user -u lhbot.service --no-pager

# Verify the binary has execute permissions
ls -la ~/.local/bin/lhbot
```

**API Authentication Errors:**
- Ensure `CLIENT_ID` and `CLIENT_SECRET` in `~/.config/lhbot.env` are correct.
- Confirm the API key has the necessary permissions in the Tencent Cloud console.

**Webhook Notification Failures:**
- Verify the `WEBHOOK` URL is correct and accessible.
- Check the `CHAT_ID` configuration.

### Debug Mode

To run the bot with verbose logging for debugging:
```bash
# Load environment variables
source ~/.config/lhbot.env

# Run the bot and pipe output to a log file
~/.local/bin/lhbot 2>&1 | tee debug.log
```

## 🛡️ Security Best Practices

- **Never commit secrets**: Do not commit API credentials or other secrets to version control.
- **Use environment variables**: Store sensitive data in environment variables, not in the code.
- **Rotate keys**: Regularly rotate your Tencent Cloud API keys.
- **Monitor logs**: Check logs for suspicious activity: `journalctl --user -u lhbot.service`.

## 📄 Recent Updates

### v1.1.0 (2025-07-25)
- **Enhanced Notifications**: Added immediate alerts for stock availability with user mentions.
- **Duplicate Purchase Prevention**: Implemented a file-based lock to avoid re-purchasing after a restart.
- **Improved Error Handling**: Added more detailed logging for `CreateInstances` API failures.
- **Systemd Service Fix**: Documented the `loginctl enable-linger` requirement for persistent user services.

## 🤝 Contributing

1.  Fork the repository.
2.  Create your feature branch.
3.  Commit your changes.
4.  Ensure your changes are well-tested.
5.  Submit a pull request.

## 📞 Support

If you encounter issues:
- Review the troubleshooting section.
- Check the system logs: `journalctl --user -u lhbot.service`.
- Open an issue on GitHub.
