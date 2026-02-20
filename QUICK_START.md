# PicoClaw Quick Start Guide

## 🚀 One-Command Setup

Start both the gateway and web UI with a single command:

```bash
docker compose up
```

This will start:
- **PicoClaw Gateway** - Handles WhatsApp and all messaging channels
- **PicoClaw Web UI** - Configuration dashboard and chat interface

## 📱 Access Points

Once running, access the web interface at:

### **http://localhost:8080**

## 🎯 Features

### 1. **Dashboard** (`/`)
- System status overview
- Enabled channels
- Recent activity
- Quick actions

### 2. **Chat Interface** (`/chat`)
- Send and receive messages through the web
- Real-time conversation view
- Multi-channel support (WhatsApp, Telegram, Discord, etc.)
- Test your bot without using your phone!

### 3. **WhatsApp QR Code** (`/whatsapp/qr`)
- Visual QR code display
- Auto-refresh every 60 seconds
- Connection status monitoring
- **No need to check Docker logs anymore!**

### 4. **Channels Configuration** (`/channels`)
- Enable/disable channels with toggle switches
- Manage whitelist (allow_from) visually
- Add/remove phone numbers with one click
- Configure API tokens and credentials

### 5. **Skills Management** (`/skills`)
- Install skills from GitHub
- View installed skills
- Remove skills

### 6. **Settings** (`/settings`)
- Configure LLM models
- Manage API keys
- Adjust agent parameters

### 7. **Live Logs** (`/logs`)
- Real-time log streaming
- Filter by component
- No need for `docker compose logs`!

## 🔧 Quick Setup Steps

1. **Start Everything**:
   ```bash
   docker compose up
   ```

2. **Connect WhatsApp**:
   - Open http://localhost:8080/whatsapp/qr
   - Scan the QR code with WhatsApp
   - Done! Your WhatsApp is connected

3. **Test the Chat**:
   - Go to http://localhost:8080/chat
   - Create a new conversation
   - Send test messages from the web interface
   - See responses in real-time

4. **Configure Channels**:
   - Go to http://localhost:8080/channels
   - Toggle channels on/off
   - Add phone numbers to whitelist
   - Save changes (auto-saves)

## 📋 Commands

### Start everything:
```bash
docker compose up
```

### Start in background:
```bash
docker compose up -d
```

### View logs:
```bash
docker compose logs -f
```

### Stop everything:
```bash
docker compose down
```

### Rebuild after code changes:
```bash
docker compose build
docker compose up
```

## 🌐 Web UI Navigation

```
┌─────────────────────────────────────────────┐
│  🦐 PicoClaw                                │
├─────────────────────────────────────────────┤
│  📊 Dashboard         - Overview            │
│  💬 Chat              - Web messaging       │
│  🔌 Channels          - Configuration       │
│  📱 WhatsApp QR       - Connect WhatsApp    │
│  🎨 Skills            - Manage skills       │
│  ⚙️  Settings         - LLM config          │
│  📝 Logs              - Live monitoring     │
└─────────────────────────────────────────────┘
```

## 💡 Pro Tips

1. **WhatsApp Connection**:
   - The QR code refreshes automatically
   - You can reconnect anytime at `/whatsapp/qr`
   - No need to restart containers

2. **Chat Interface**:
   - Perfect for testing without sending real WhatsApp messages
   - See AI responses in real-time
   - Works with all configured channels

3. **Whitelist Management**:
   - Add phone numbers in format: `5511999999999@s.whatsapp.net`
   - Empty whitelist = everyone can use the bot
   - Changes take effect immediately

4. **Configuration Changes**:
   - Most changes auto-save
   - Some require gateway restart:
     ```bash
     docker compose restart picoclaw-gateway
     ```

## 🎨 SaaS Ready

This setup is designed to be SaaS-ready:
- Each user gets their own Docker instance
- Isolated configuration per instance
- Web UI for easy onboarding
- No technical knowledge required for users

## 📚 More Documentation

- [Full Web UI Documentation](WEBUI.md)
- [WhatsApp Integration](pkg/channels/whatsapp.go)
- [Configuration Reference](config/config.example.json)

## 🐛 Troubleshooting

**Web UI not accessible?**
```bash
docker compose logs picoclaw-webui
```

**WhatsApp not connecting?**
- Check that WhatsApp channel is enabled in `/channels`
- Restart gateway: `docker compose restart picoclaw-gateway`
- Check logs at `/logs` page

**Gateway not starting?**
```bash
docker compose logs picoclaw-gateway
```

---

**That's it!** You now have a complete WhatsApp AI bot with a beautiful web interface running with just `docker compose up` 🎉
