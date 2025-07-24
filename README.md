# LightHouse Razor Bot (LHBot)

A lightweight Go bot that monitors Tencent Cloud Lighthouse instance availability and automatically purchases instances when stock becomes available.

## 🎯 Purpose

This bot continuously monitors specific Tencent Cloud Lighthouse bundle availability (particularly "锐驰" series) and:
- Automatically purchases instances when stock becomes available
- Sends notifications via WeChat Work webhooks about stock status
- Provides periodic monitoring reports

## 🚀 Features

- **Real-time monitoring**: Checks bundle availability every 30 seconds
- **Auto-purchase**: Automatically creates instances when stock is detected
- **Smart notifications**: Sends hourly status updates and purchase confirmations
- **Graceful shutdown**: Handles SIGINT/SIGTERM signals properly
- **Cross-platform**: Supports Linux AMD64 and ARM64 architectures

## 📋 Prerequisites

- **Tencent Cloud Account** with Lighthouse service enabled
- **API Credentials**: Tencent Cloud SecretId and SecretKey
- **WeChat Work** webhook URL and chat ID for notifications
- **Go 1.24.4** or later for building from source

## 🔧 Configuration

### Environment Variables

Create a configuration file at `~/.config/lhbot.env`:

```bash
# Tencent Cloud Credentials
CLIENT_ID=your_tencent_secret_id
CLIENT_SECRET=your_tencent_secret_key

# WeChat Work Notifications
WEBHOOK=https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=your_webhook_key
CHAT_ID=your_chat_id

# Optional: Root password for created instances (default: admin@2025)
ROOT_PASSWORD=your_secure_password
```

## 🏗️ Building

### From Source

```bash
# Clone the repository
git clone <repository-url>
cd lhbot

# Build for current platform
go build -o lhbot .

# Build for all supported platforms (Linux AMD64/ARM64)
./make.ps1 -Platform all
```

### Pre-built Binaries

Pre-compiled binaries are available in the `bin/` directory:
- `lhbot_linux-amd64` - Linux AMD64
- `lhbot_linux-arm64` - Linux ARM64

## 🔧 Installation & Deployment

### 1. Binary Installation

```bash
# Create local bin directory
mkdir -p ~/.local/bin

# Copy binary
cp bin/lhbot_linux-amd64 ~/.local/bin/lhbot
chmod +x ~/.local/bin/lhbot

# Create config directory
mkdir -p ~/.config
```

### 2. User-Level systemd Service

#### Setup Service

```bash
# Copy systemd service file
cp systemd/lhbot.service ~/.config/systemd/user/

# Reload systemd user daemon
systemctl --user daemon-reload

# Enable service to start at boot
systemctl --user enable lhbot.service

# Start the service
systemctl --user start lhbot.service

# Check service status
systemctl --user status lhbot.service
```

#### View Logs

```bash
# Real-time logs
journalctl --user -f -u lhbot.service

# Recent logs
journalctl --user -u lhbot.service --since "1 hour ago"
```

#### Service Management

```bash
# Stop service
systemctl --user stop lhbot.service

# Restart service
systemctl --user restart lhbot.service

# Disable service
systemctl --user disable lhbot.service
```

### 3. Manual Execution

```bash
# Load environment variables
source ~/.config/lhbot.env

# Run directly
~/.local/bin/lhbot
```

## 📊 Monitoring

The bot provides two types of notifications:

### Status Reports (Hourly)
```
⚙️ **监控服务运行中**
- **锐驰-2C4G**: SOLD_OUT
- **锐驰-4C8G**: AVAILABLE

**通知时间**：2025-07-24 10:30:00
```

### Purchase Confirmations
```
✅ **锐驰自动购买成功**
- **型号**: 锐驰-4C8G

**通知时间**：2025-07-24 10:31:15
```

## ⚙️ Configuration Files

### systemd Service File
Located at `~/.config/systemd/user/lhbot.service`

Key configurations:
- **EnvironmentFile**: `~/.config/lhbot.env`
- **ExecStart**: `~/.local/bin/lhbot`
- **Restart**: Auto-restart on failure (max 3 times in 30 seconds)
- **Logging**: Output to `/var/log/lhbot.log`

## 🔍 Troubleshooting

### Common Issues

**Service won't start:**
```bash
# Check service status
systemctl --user status lhbot.service

# View detailed logs
journalctl --user -u lhbot.service --no-pager

# Check binary permissions
ls -la ~/.local/bin/lhbot
```

**API authentication errors:**
- Verify `CLIENT_ID` and `CLIENT_SECRET` are correct
- Check Tencent Cloud API permissions

**Webhook notifications failing:**
- Verify `WEBHOOK` URL is accessible
- Check `CHAT_ID` configuration

### Debug Mode

```bash
# Run with debug output
source ~/.config/lhbot.env
~/.local/bin/lhbot 2>&1 | tee debug.log
```

## 🛡️ Security Notes

- Never commit API credentials to version control
- Use environment variables for sensitive configuration
- Regularly rotate Tencent Cloud API keys
- Monitor `/var/log/lhbot.log` for security events

## 📄 License

This project is open source. Please review the license file for details.

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## 📞 Support

For issues and questions:
- Check the troubleshooting section
- Review system logs: `journalctl --user -u lhbot.service`
- Open an issue in the repository