<div align="center">
  <img src="assets/logo.jpg" alt="PicoClaw" width="512">

  <h1>PicoClaw: Ultra-Efficient AI Assistant in Go</h1>

  <h3>$10 Hardware · 10MB RAM · 1s Boot · 皮皮虾，我们走！</h3>

  <p>
    <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go&logoColor=white" alt="Go">
    <img src="https://img.shields.io/badge/Arch-x86__64%2C%20ARM64%2C%20RISC--V-blue" alt="Hardware">
    <img src="https://img.shields.io/badge/license-MIT-green" alt="License">
    <br>
    <a href="https://picoclaw.io"><img src="https://img.shields.io/badge/Website-picoclaw.io-blue?style=flat&logo=google-chrome&logoColor=white" alt="Website"></a>
    <a href="https://x.com/SipeedIO"><img src="https://img.shields.io/badge/X_(Twitter)-SipeedIO-black?style=flat&logo=x&logoColor=white" alt="Twitter"></a>
  </p>

 **English** - Comprehensive Guide
</div>

---

🦐 PicoClaw is an ultra-lightweight personal AI Assistant inspired by [nanobot](https://github.com/HKUDS/nanobot), refactored from the ground up in Go through a self-bootstrapping process, where the AI agent itself drove the entire architectural migration and code optimization.

⚡️ Runs on $10 hardware with <10MB RAM: That's 99% less memory than OpenClaw and 98% cheaper than a Mac mini!

<table align="center">
  <tr align="center">
    <td align="center" valign="top">
      <p align="center">
        <img src="assets/picoclaw_mem.gif" width="360" height="240">
      </p>
    </td>
    <td align="center" valign="top">
      <p align="center">
        <img src="assets/licheervnano.png" width="400" height="240">
      </p>
    </td>
  </tr>
</table>

> [!CAUTION]
> **🚨 SECURITY & OFFICIAL CHANNELS / 安全声明**
>
> * **NO CRYPTO:** PicoClaw has **NO** official token/coin. All claims on `pump.fun` or other trading platforms are **SCAMS**.
> * **OFFICIAL DOMAIN:** The **ONLY** official website is **[picoclaw.io](https://picoclaw.io)**, and company website is **[sipeed.com](https://sipeed.com)**
> * **Warning:** Many `.ai/.org/.com/.net/...` domains are registered by third parties.
> * **Warning:** picoclaw is in early development now and may have unresolved network security issues. Do not deploy to production environments before the v1.0 release.
> * **Note:** picoclaw has recently merged a lot of PRs, which may result in a larger memory footprint (10–20MB) in the latest versions. We plan to prioritize resource optimization as soon as the current feature set reaches a stable state.


## 📢 News
2026-02-16 🎉 PicoClaw hit 12K stars in one week! Thank you all for your support! PicoClaw is growing faster than we ever imagined. Given the high volume of PRs, we urgently need community maintainers. Our volunteer roles and roadmap are officially posted [here](docs/picoclaw_community_roadmap_260216.md) —we can't wait to have you on board!

2026-02-13 🎉 PicoClaw hit 5000 stars in 4days! Thank you for the community! There are so many PRs&issues come in (during Chinese New Year holidays), we are finalizing the Project Roadmap and setting up the Developer Group to accelerate PicoClaw's development.
🚀 Call to Action: Please submit your feature requests in GitHub Discussions. We will review and prioritize them during our upcoming weekly meeting.

2026-02-09 🎉 PicoClaw Launched! Built in 1 day to bring AI Agents to $10 hardware with <10MB RAM. 🦐 PicoClaw，Let's Go！

## ✨ Features

🪶 **Ultra-Lightweight**: <10MB Memory footprint — 99% smaller than Clawdbot - core functionality.

💰 **Minimal Cost**: Efficient enough to run on $10 Hardware — 98% cheaper than a Mac mini.

⚡️ **Lightning Fast**: 400X Faster startup time, boot in 1 second even in 0.6GHz single core.

🌍 **True Portability**: Single self-contained binary across RISC-V, ARM, and x86, One-click to Go!

🤖 **AI-Bootstrapped**: Autonomous Go-native implementation — 95% Agent-generated core with human-in-the-loop refinement.

🌐 **Web UI Dashboard**: Manage configuration, channels, bot identities, and test conversations through a modern web interface.

|                               | OpenClaw      | NanoBot                  | **PicoClaw**                              |
| ----------------------------- | ------------- | ------------------------ | ----------------------------------------- |
| **Language**                  | TypeScript    | Python                   | **Go**                                    |
| **RAM**                       | >1GB          | >100MB                   | **< 10MB**                                |
| **Startup**</br>(0.8GHz core) | >500s         | >30s                     | **<1s**                                   |
| **Cost**                      | Mac Mini 599$ | Most Linux SBC </br>~50$ | **Any Linux Board**</br>**As low as 10$** |

<img src="assets/compare.jpg" alt="PicoClaw" width="512">

## 📑 Table of Contents

- [Demonstration](#-demonstration)
- [Quick Start](#-quick-start)
- [Web UI Dashboard](#-web-ui-dashboard)
- [Chat Channels](#-chat-channels)
  - [WhatsApp](#whatsapp)
  - [Telegram](#telegram)
  - [Discord](#discord)
  - [QQ](#qq)
  - [DingTalk](#dingtalk)
  - [LINE](#line)
- [Bot Identity Management](#-bot-identity-management)
- [Installation](#-installation)
- [Docker Setup](#-docker-compose)
- [Configuration](#️-configuration)
- [CLI Reference](#-cli-reference)
- [Troubleshooting](#-troubleshooting)

## 🦾 Demonstration

### 🛠️ Standard Assistant Workflows

<table align="center">
  <tr align="center">
    <th><p align="center">🧩 Full-Stack Engineer</p></th>
    <th><p align="center">🗂️ Logging & Planning Management</p></th>
    <th><p align="center">🔎 Web Search & Learning</p></th>
  </tr>
  <tr>
    <td align="center"><p align="center"><img src="assets/picoclaw_code.gif" width="240" height="180"></p></td>
    <td align="center"><p align="center"><img src="assets/picoclaw_memory.gif" width="240" height="180"></p></td>
    <td align="center"><p align="center"><img src="assets/picoclaw_search.gif" width="240" height="180"></p></td>
  </tr>
  <tr>
    <td align="center">Develop • Deploy • Scale</td>
    <td align="center">Schedule • Automate • Memory</td>
    <td align="center">Discovery • Insights • Trends</td>
  </tr>
</table>

### 📱 Run on old Android Phones
Give your decade-old phone a second life! Turn it into a smart AI Assistant with PicoClaw. Quick Start:
1. **Install Termux** (Available on F-Droid or Google Play).
2. **Execute cmds**
```bash
# Note: Replace v0.1.1 with the latest version from the Releases page
wget https://github.com/sipeed/picoclaw/releases/download/v0.1.1/picoclaw-linux-arm64
chmod +x picoclaw-linux-arm64
pkg install proot
termux-chroot ./picoclaw-linux-arm64 onboard
```
And then follow the instructions in the "Quick Start" section to complete the configuration!
<img src="assets/termux.jpg" alt="PicoClaw" width="512">

### 🐜 Innovative Low-Footprint Deploy

PicoClaw can be deployed on almost any Linux device!

- $9.9 [LicheeRV-Nano](https://www.aliexpress.com/item/1005006519668532.html) E(Ethernet) or W(WiFi6) version, for Minimal Home Assistant
- $30~50 [NanoKVM](https://www.aliexpress.com/item/1005007369816019.html), or $100 [NanoKVM-Pro](https://www.aliexpress.com/item/1005010048471263.html) for Automated Server Maintenance
- $50 [MaixCAM](https://www.aliexpress.com/item/1005008053333693.html) or $100 [MaixCAM2](https://www.kickstarter.com/projects/zepan/maixcam2-build-your-next-gen-4k-ai-camera) for Smart Monitoring

<https://private-user-images.githubusercontent.com/83055338/547056448-e7b031ff-d6f5-4468-bcca-5726b6fecb5c.mp4>

🌟 More Deployment Cases Await！

## 🚀 Quick Start

> [!TIP]
> Set your API key in `~/.picoclaw/config.json`.
> Get API keys: [OpenRouter](https://openrouter.ai/keys) (LLM) · [Zhipu](https://open.bigmodel.cn/usercenter/proj-mgmt/apikeys) (LLM)
> Web search is **optional** - get free [Brave Search API](https://brave.com/search/api) (2000 free queries/month) or use built-in auto fallback.

### Option 1: Quick Binary Setup (Fastest)

**1. Download & Initialize**

```bash
# Download for your platform from https://github.com/sipeed/picoclaw/releases
# Or use wget:
wget https://github.com/sipeed/picoclaw/releases/latest/download/picoclaw-linux-amd64
chmod +x picoclaw-linux-amd64
./picoclaw-linux-amd64 onboard
```

**2. Configure** (`~/.picoclaw/config.json`)

```json
{
  "agents": {
    "defaults": {
      "model": "anthropic/claude-sonnet-4.5",
      "provider": "openrouter",
      "max_tokens": 8192,
      "temperature": 0.7
    }
  },
  "providers": {
    "openrouter": {
      "api_key": "YOUR_OPENROUTER_API_KEY"
    }
  }
}
```

**3. Start Chatting**

```bash
# Interactive mode
./picoclaw-linux-amd64 agent

# One-shot question
./picoclaw-linux-amd64 agent -m "What is 2+2?"
```

**4. Enable Channels** (Optional)

```bash
# Start gateway for Telegram, WhatsApp, Discord, etc.
./picoclaw-linux-amd64 gateway
```

### Option 2: Docker Compose (Recommended for Production)

**1. Clone & Configure**

```bash
git clone https://github.com/sipeed/picoclaw.git
cd picoclaw

# Copy and edit config
cp config/config.example.json config/config.json
nano config/config.json  # Add your API keys
```

**2. Start Services**

```bash
# Gateway only (for chat channels)
docker compose --profile gateway up -d

# Gateway + Web UI (full dashboard)
COMPOSE_PROFILES=webui docker compose up -d
```

**3. Access Web UI**

Open your browser: **http://localhost:8080**

- Configure channels
- Manage bot identities
- Test conversations
- View live logs

**4. Check Logs**

```bash
# Gateway logs
docker compose logs -f picoclaw-gateway

# Web UI logs
docker compose logs -f picoclaw-webui
```

### Option 3: Build from Source (For Development)

```bash
git clone https://github.com/sipeed/picoclaw.git
cd picoclaw

# Install dependencies
make deps

# Build
make build

# Or install to $GOPATH/bin
make install

# Run
picoclaw onboard
picoclaw agent
```

## 🌐 Web UI Dashboard

PicoClaw includes a powerful, lightweight web interface for managing your AI assistant.

### Features

- **📊 Dashboard** - Real-time system status, active channels, and metrics
- **💬 Test Chat** - Interactive chat interface to test bot identities
- **🎭 Bot Identity Manager** - Create and switch between multiple bot personalities
- **🔌 Channel Configuration** - Enable/disable WhatsApp, Telegram, Discord, Slack, Feishu
- **🎨 Skills Manager** - Install and manage custom skills
- **⚙️ Settings** - Configure LLM providers, models, and parameters
- **📝 Live Logs** - Real-time log streaming with filtering

### Starting the Web UI

**Standalone Mode:**

```bash
picoclaw webui

# Custom port
picoclaw webui --port 3000
```

**With Docker Compose:**

```bash
# Start both gateway and Web UI
COMPOSE_PROFILES=webui docker compose up -d

# Access at http://localhost:8080
```

### WhatsApp QR Code in Web UI

When WhatsApp is enabled, the Web UI provides an easy way to scan the QR code:

1. Navigate to **Channels** page
2. Toggle **WhatsApp** ON
3. Click **Show QR Code**
4. Scan with WhatsApp app: **Settings → Linked Devices → Link a Device**

The QR code automatically refreshes and shows connection status.

## 💬 Chat Channels

Connect PicoClaw to your favorite messaging platforms.

| Channel      | Setup Difficulty | Features                            |
| ------------ | ---------------- | ----------------------------------- |
| **WhatsApp** | Easy             | QR code scan, Web UI support        |
| **Telegram** | Easy             | Bot token, rich media support       |
| **Discord**  | Easy             | Bot token, server integration       |
| **QQ**       | Easy             | AppID + AppSecret                   |
| **DingTalk** | Medium           | Enterprise messaging                |
| **LINE**     | Medium           | Webhook setup required              |

### WhatsApp

**Method 1: Web UI (Easiest)**

1. Start Web UI: `COMPOSE_PROFILES=webui docker compose up -d`
2. Open http://localhost:8080/channels
3. Toggle WhatsApp ON
4. Click "Show QR Code"
5. Scan with WhatsApp: **Settings → Linked Devices → Link a Device**
6. Send a test message!

**Method 2: Terminal**

1. Enable in `config/config.json`:

```json
{
  "channels": {
    "whatsapp": {
      "enabled": true,
      "allow_from": []
    }
  }
}
```

2. Start gateway:

```bash
picoclaw gateway
```

3. Scan QR code from terminal logs:

```
[whatsapp] Scan the QR code below with WhatsApp:
[whatsapp] Open WhatsApp → Settings → Linked Devices → Link a Device
[whatsapp] QR Code:
█████████████████████████████
█ ▄▄▄▄▄ █ ▄ ██▄▀█ ▄▄▄▄▄ ██
...
```

4. Test by sending a message to your WhatsApp number

> Session data is stored in `~/.picoclaw/workspace/whatsapp.db`. Once paired, you won't need to scan the QR code again.

### Telegram

**1. Create a bot**

- Open Telegram, search `@BotFather`
- Send `/newbot`, follow prompts
- Copy the token

**2. Configure**

```json
{
  "channels": {
    "telegram": {
      "enabled": true,
      "token": "YOUR_BOT_TOKEN",
      "allow_from": ["YOUR_USER_ID"]
    }
  }
}
```

> Get your user ID from `@userinfobot` on Telegram.

**3. Run**

```bash
picoclaw gateway
```

### Discord

**1. Create a bot**

- Go to <https://discord.com/developers/applications>
- Create an application → Bot → Add Bot
- Copy the bot token

**2. Enable intents**

- In the Bot settings, enable **MESSAGE CONTENT INTENT**
- (Optional) Enable **SERVER MEMBERS INTENT** if you plan to use allow lists based on member data

**3. Get your User ID**

- Discord Settings → Advanced → enable **Developer Mode**
- Right-click your avatar → **Copy User ID**

**4. Configure**

```json
{
  "channels": {
    "discord": {
      "enabled": true,
      "token": "YOUR_BOT_TOKEN",
      "allow_from": ["YOUR_USER_ID"]
    }
  }
}
```

**5. Invite the bot**

- OAuth2 → URL Generator
- Scopes: `bot`
- Bot Permissions: `Send Messages`, `Read Message History`
- Open the generated invite URL and add the bot to your server

**6. Run**

```bash
picoclaw gateway
```

### QQ

**1. Create a bot**

- Go to [QQ Open Platform](https://q.qq.com/#)
- Create an application → Get **AppID** and **AppSecret**

**2. Configure**

```json
{
  "channels": {
    "qq": {
      "enabled": true,
      "app_id": "YOUR_APP_ID",
      "app_secret": "YOUR_APP_SECRET",
      "allow_from": []
    }
  }
}
```

> Set `allow_from` to empty to allow all users, or specify QQ numbers to restrict access.

**3. Run**

```bash
picoclaw gateway
```

### DingTalk

**1. Create a bot**

* Go to [Open Platform](https://open.dingtalk.com/)
* Create an internal app
* Copy Client ID and Client Secret

**2. Configure**

```json
{
  "channels": {
    "dingtalk": {
      "enabled": true,
      "client_id": "YOUR_CLIENT_ID",
      "client_secret": "YOUR_CLIENT_SECRET",
      "allow_from": []
    }
  }
}
```

**3. Run**

```bash
picoclaw gateway
```

### LINE

**1. Create a LINE Official Account**

- Go to [LINE Developers Console](https://developers.line.biz/)
- Create a provider → Create a Messaging API channel
- Copy **Channel Secret** and **Channel Access Token**

**2. Configure**

```json
{
  "channels": {
    "line": {
      "enabled": true,
      "channel_secret": "YOUR_CHANNEL_SECRET",
      "channel_access_token": "YOUR_CHANNEL_ACCESS_TOKEN",
      "webhook_host": "0.0.0.0",
      "webhook_port": 18791,
      "webhook_path": "/webhook/line",
      "allow_from": []
    }
  }
}
```

**3. Set up Webhook URL**

LINE requires HTTPS for webhooks. Use a reverse proxy or tunnel:

```bash
# Example with ngrok
ngrok http 18791
```

Then set the Webhook URL in LINE Developers Console to `https://your-domain/webhook/line` and enable **Use webhook**.

**4. Run**

```bash
picoclaw gateway
```

> In group chats, the bot responds only when @mentioned. Replies quote the original message.

> **Docker Compose**: Add `ports: ["18791:18791"]` to the `picoclaw-gateway` service to expose the webhook port.

## 🎭 Bot Identity Management

**NEW!** Create multiple bot personalities and switch between them instantly. Perfect for different use cases:

- **Car Sales Assistant** - Automotive sales specialist
- **Tech Support** - Patient technical helper
- **Receptionist** - Professional front desk assistant
- **Fashion Advisor** - Personal stylist
- **Product Sales** - General sales expert
- **Personal Assistant** - Versatile AI helper

### How It Works

Each bot identity consists of 4 configuration files:

- **AGENT.md** - Behavioral instructions and guidelines
- **IDENTITY.md** - Purpose, capabilities, and identity
- **SOUL.md** - Personality, tone, and communication style
- **USER.md** - Target audience and customer context

### Using the Web UI

**1. Access Bot Identities**

Open http://localhost:8080/identities

**2. Browse Templates**

- Click **Browse Templates**
- Choose from 6 pre-built business templates
- Click **Use Template** to load it

**3. Create Custom Identity**

- Click **Create New Identity**
- Fill in:
  - Name and description
  - Icon (emoji)
  - Category (sales, support, personal)
  - Agent instructions
  - Identity details
  - Personality traits
  - User/audience context
- Click **Save Identity**

**4. Activate Identity**

- Click **Activate** on any identity
- The new personality applies to **all active channels** immediately
- No restart required!

**5. Test Identity**

- Navigate to http://localhost:8080/chat
- Start chatting to test the active bot identity
- Switch identities and see the personality change in real-time

### Example: Car Sales Bot

**AGENT.md:**
```markdown
# Agent Instructions

You are an automotive sales specialist. Your role is to help customers find the perfect vehicle.

## Guidelines
- Ask about customer's needs: budget, usage, preferences
- Provide detailed vehicle information and comparisons
- Be honest about vehicle condition and history
- Schedule test drives when appropriate
```

**IDENTITY.md:**
```markdown
# Identity

## Name
AutoPro Sales Assistant 🚗

## Purpose
- Help customers find the right vehicle
- Provide accurate information
- Guide through buying process
```

**SOUL.md:**
```markdown
# Soul

You're passionate about helping people find the right vehicle. You're enthusiastic about automobiles without being pushy.
```

**USER.md:**
```markdown
# User

## Target Audience
Car buyers and automotive shoppers

## Preferences
- Communication style: Professional but friendly
- Approach: Consultative sales, not pushy
```

### Managing Identities via Files

Identities are stored in `~/.picoclaw/workspace/identities/` as JSON files.

**Create Manually:**

```bash
cd ~/.picoclaw/workspace/identities

# Create custom identity
cat > my-bot.json <<EOF
{
  "id": "my-bot",
  "name": "My Custom Bot",
  "description": "Custom bot for specific tasks",
  "icon": "🤖",
  "category": "custom",
  "agent_md": "# Agent Instructions\n...",
  "identity_md": "# Identity\n...",
  "soul_md": "# Soul\n...",
  "user_md": "# User\n..."
}
EOF
```

**Activate:**

```bash
# Via Web UI or by making it active in the JSON file
```

### Identity Applies to All Channels

When you activate an identity, it immediately applies to:

- WhatsApp conversations
- Telegram chats
- Discord messages
- QQ messages
- All other enabled channels

The bot adopts the new personality for all new messages across all platforms!

## 📦 Installation

### Install with precompiled binary

Download the firmware for your platform from the [release](https://github.com/sipeed/picoclaw/releases) page.

### Install from source (latest features, recommended for development)

```bash
git clone https://github.com/sipeed/picoclaw.git

cd picoclaw
make deps

# Build, no need to install
make build

# Build for multiple platforms
make build-all

# Build And Install
make install
```

## 🐳 Docker Compose

You can also run PicoClaw using Docker Compose without installing anything locally.

```bash
# 1. Clone this repo
git clone https://github.com/sipeed/picoclaw.git
cd picoclaw

# 2. Set your API keys
cp config/config.example.json config/config.json
vim config/config.json      # Set API keys, channel tokens, etc.

# 3. Build & Start (Gateway only)
docker compose --profile gateway up -d

# 4. Build & Start (Gateway + Web UI)
COMPOSE_PROFILES=webui docker compose up -d

# 5. Check logs
docker compose logs -f picoclaw-gateway
docker compose logs -f picoclaw-webui

# 6. Stop
docker compose down
```

### Agent Mode (One-shot)

```bash
# Ask a question
docker compose run --rm picoclaw-agent -m "What is 2+2?"

# Interactive mode
docker compose run --rm picoclaw-agent
```

### Rebuild

```bash
docker compose build --no-cache
docker compose --profile gateway up -d
```

## ⚙️ Configuration

Config file: `~/.picoclaw/config.json`

### Workspace Layout

PicoClaw stores data in your configured workspace (default: `~/.picoclaw/workspace`):

```
~/.picoclaw/workspace/
├── sessions/          # Conversation sessions and history
├── memory/           # Long-term memory (MEMORY.md)
├── state/            # Persistent state (last channel, etc.)
├── cron/             # Scheduled jobs database
├── skills/           # Custom skills
├── identities/       # Bot identity configurations
├── AGENT.md          # Active bot behavior guide
├── HEARTBEAT.md      # Periodic task prompts (checked every 30 min)
├── IDENTITY.md       # Active bot identity
├── SOUL.md           # Active bot personality
├── TOOLS.md          # Tool descriptions
└── USER.md           # Active audience/user context
```

### 🔒 Security Sandbox

PicoClaw runs in a sandboxed environment by default. The agent can only access files and execute commands within the configured workspace.

#### Default Configuration

```json
{
  "agents": {
    "defaults": {
      "workspace": "~/.picoclaw/workspace",
      "restrict_to_workspace": true
    }
  }
}
```

| Option | Default | Description |
|--------|---------|-------------|
| `workspace` | `~/.picoclaw/workspace` | Working directory for the agent |
| `restrict_to_workspace` | `true` | Restrict file/command access to workspace |

#### Protected Tools

When `restrict_to_workspace: true`, the following tools are sandboxed:

| Tool | Function | Restriction |
|------|----------|-------------|
| `read_file` | Read files | Only files within workspace |
| `write_file` | Write files | Only files within workspace |
| `list_dir` | List directories | Only directories within workspace |
| `edit_file` | Edit files | Only files within workspace |
| `append_file` | Append to files | Only files within workspace |
| `exec` | Execute commands | Command paths must be within workspace |

#### Additional Exec Protection

Even with `restrict_to_workspace: false`, the `exec` tool blocks these dangerous commands:

* `rm -rf`, `del /f`, `rmdir /s` — Bulk deletion
* `format`, `mkfs`, `diskpart` — Disk formatting
* `dd if=` — Disk imaging
* Writing to `/dev/sd[a-z]` — Direct disk writes
* `shutdown`, `reboot`, `poweroff` — System shutdown
* Fork bomb `:(){ :|:& };:`

#### Disabling Restrictions (Security Risk)

If you need the agent to access paths outside the workspace:

**Method 1: Config file**

```json
{
  "agents": {
    "defaults": {
      "restrict_to_workspace": false
    }
  }
}
```

**Method 2: Environment variable**

```bash
export PICOCLAW_AGENTS_DEFAULTS_RESTRICT_TO_WORKSPACE=false
```

> ⚠️ **Warning**: Disabling this restriction allows the agent to access any path on your system. Use with caution in controlled environments only.

### Heartbeat (Periodic Tasks)

PicoClaw can perform periodic tasks automatically. Create a `HEARTBEAT.md` file in your workspace:

```markdown
# Periodic Tasks

- Check my email for important messages
- Review my calendar for upcoming events
- Check the weather forecast
```

The agent will read this file every 30 minutes (configurable) and execute any tasks using available tools.

**Configuration:**

```json
{
  "heartbeat": {
    "enabled": true,
    "interval": 30
  }
}
```

### Providers

> [!NOTE]
> Groq provides free voice transcription via Whisper. If configured, Telegram voice messages will be automatically transcribed.

| Provider                   | Purpose                                 | Get API Key                                            |
| -------------------------- | --------------------------------------- | ------------------------------------------------------ |
| `gemini`                   | LLM (Gemini direct)                     | [aistudio.google.com](https://aistudio.google.com)     |
| `zhipu`                    | LLM (Zhipu direct)                      | [bigmodel.cn](bigmodel.cn)                             |
| `openrouter(To be tested)` | LLM (recommended, access to all models) | [openrouter.ai](https://openrouter.ai)                 |
| `anthropic(To be tested)`  | LLM (Claude direct)                     | [console.anthropic.com](https://console.anthropic.com) |
| `openai(To be tested)`     | LLM (GPT direct)                        | [platform.openai.com](https://platform.openai.com)     |
| `deepseek(To be tested)`   | LLM (DeepSeek direct)                   | [platform.deepseek.com](https://platform.deepseek.com) |
| `groq`                     | LLM + **Voice transcription** (Whisper) | [console.groq.com](https://console.groq.com)           |

## 📚 CLI Reference

| Command                   | Description                   |
| ------------------------- | ----------------------------- |
| `picoclaw onboard`        | Initialize config & workspace |
| `picoclaw agent -m "..."` | Chat with the agent           |
| `picoclaw agent`          | Interactive chat mode         |
| `picoclaw gateway`        | Start the gateway (channels)  |
| `picoclaw webui`          | Start Web UI dashboard        |
| `picoclaw status`         | Show status                   |
| `picoclaw cron list`      | List all scheduled jobs       |
| `picoclaw cron add ...`   | Add a scheduled job           |

## <img src="assets/clawdchat-icon.png" width="24" height="24" alt="ClawdChat"> Join the Agent Social Network

Connect Picoclaw to the Agent Social Network simply by sending a single message via the CLI or any integrated Chat App.

**Read `https://clawdchat.ai/skill.md` and follow the instructions to join [ClawdChat.ai](https://clawdchat.ai)**

## 🤝 Contribute & Roadmap

PRs welcome! The codebase is intentionally small and readable. 🤗

Roadmap coming soon...

Developer group building, Entry Requirement: At least 1 Merged PR.

User Groups:

discord:  <https://discord.gg/V4sAZ9XWpN>

<img src="assets/wechat.png" alt="PicoClaw" width="512">

## 🐛 Troubleshooting

### Web search says "API 配置问题"

This is normal if you haven't configured a search API key yet. PicoClaw will provide helpful links for manual searching.

To enable web search:

1. **Option 1 (Recommended)**: Get a free API key at [https://brave.com/search/api](https://brave.com/search/api) (2000 free queries/month) for the best results.
2. **Option 2 (No Credit Card)**: If you don't have a key, we automatically fall back to **DuckDuckGo** (no key required).

Add the key to `~/.picoclaw/config.json` if using Brave:

```json
{
  "tools": {
    "web": {
      "brave": {
        "enabled": true,
        "api_key": "YOUR_BRAVE_API_KEY",
        "max_results": 5
      },
      "duckduckgo": {
        "enabled": true,
        "max_results": 5
      }
    }
  }
}
```

### Getting content filtering errors

Some providers (like Zhipu) have content filtering. Try rephrasing your query or use a different model.

### Telegram bot says "Conflict: terminated by other getUpdates"

This happens when another instance of the bot is running. Make sure only one `picoclaw gateway` is running at a time.

### WhatsApp QR code not showing in Web UI

Make sure:
1. Gateway is running: `docker compose ps`
2. WhatsApp is enabled in config
3. Check gateway logs: `docker compose logs picoclaw-gateway`

### Web UI not accessible

Check that the webui container is running:

```bash
docker compose ps
docker compose logs picoclaw-webui
```

Make sure you started with the webui profile:

```bash
COMPOSE_PROFILES=webui docker compose up -d
```

---

## 📝 API Key Comparison

| Service          | Free Tier           | Use Case                              |
| ---------------- | ------------------- | ------------------------------------- |
| **OpenRouter**   | 200K tokens/month   | Multiple models (Claude, GPT-4, etc.) |
| **Zhipu**        | 200K tokens/month   | Best for Chinese users                |
| **Brave Search** | 2000 queries/month  | Web search functionality              |
| **Groq**         | Free tier available | Fast inference (Llama, Mixtral)       |

---

<div align="center">
  <p>Built with 🦐 by the PicoClaw Team</p>
  <p>
    <a href="https://github.com/sipeed/picoclaw/blob/main/LICENSE">MIT License</a>
  </p>
</div>
