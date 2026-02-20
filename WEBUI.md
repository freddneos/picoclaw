# 🌐 PicoClaw Web UI - Configuration Dashboard

A lightweight, pure-Go web interface for managing your PicoClaw instance.

## ✨ Features

- 📊 **Real-time Dashboard** - Monitor system status and activity
- 💬 **Channel Management** - Configure WhatsApp, Telegram, Discord, Slack, Feishu
- 🔐 **Whitelist Control** - Manage `allow_from` access lists visually
- 🎨 **Skills Manager** - Install, remove, and view skills
- ⚙️ **LLM Settings** - Configure models, API keys, and parameters
- 📝 **Live Logs** - Real-time log streaming
- 🪶 **Ultra-Lightweight** - Pure Go, no external dependencies, <1MB memory overhead

## 🚀 Quick Start

### Running Standalone

```bash
# Start the web UI
picoclaw webui

# Custom port
picoclaw webui --port 3000

# Bind to specific host
picoclaw webui --host 127.0.0.1 --port 8080
```

Then open your browser: http://localhost:8080

### Running with Docker

```bash
# Add to docker-compose.yml
services:
  picoclaw-webui:
    build: .
    command: webui --host 0.0.0.0 --port 8080
    ports:
      - "8080:8080"
    volumes:
      - ~/.picoclaw:/home/picoclaw/.picoclaw

# Start it
docker compose up picoclaw-webui
```

## 📱 Screenshots & Usage

### 1. Dashboard

The main dashboard shows:
- System status (uptime, version)
- Enabled channels
- Quick actions
- Recent activity logs

### 2. Channels Configuration

**WhatsApp Setup:**
1. Toggle WhatsApp channel ON
2. Restart gateway: `docker compose restart picoclaw-gateway`
3. Check logs for QR code: `docker compose logs picoclaw-gateway`
4. Scan QR code with WhatsApp app
5. Send a test message
6. Copy your WhatsApp ID from logs (format: `5511999999999@s.whatsapp.net`)
7. Add to whitelist in Web UI

**Whitelist Management:**
- Click "Add User" to add phone numbers
- Empty whitelist = open to everyone
- Click "Remove" to revoke access

**Other Channels:**
- Configure tokens/credentials directly in the UI
- Toggle channels on/off with a switch
- Changes saved automatically

### 3. Skills Manager

- **Install Skills:** Enter GitHub repo (e.g., `sipeed/picoclaw-skills/weather`)
- **View Installed:** See all active skills with descriptions
- **Remove Skills:** One-click removal

### 4. LLM Settings

Configure:
- Default model (Claude, GPT-4, Gemini, etc.)
- Max tokens
- Temperature
- Tool iteration limits
- API keys for all providers

### 5. Live Logs

- Real-time log streaming (Server-Sent Events)
- Color-coded by level and component
- Auto-scroll
- Clear button

## 🔧 Architecture

```
pkg/webui/
├── server.go              # HTTP server & routes
├── templates/
│   ├── layout.html        # Base layout with nav
│   ├── dashboard.html     # Main dashboard
│   ├── channels.html      # Channel configuration
│   ├── skills.html        # Skills management
│   ├── settings.html      # LLM settings
│   └── logs.html          # Live logs viewer
```

**Tech Stack:**
- Pure Go `net/http`
- `html/template` for rendering
- Server-Sent Events (SSE) for real-time logs
- No JavaScript frameworks (vanilla JS only)
- Embedded templates via `go:embed`

**Why Pure Go?**
- ✅ Maintains PicoClaw's ultra-lightweight philosophy
- ✅ No build toolchain required
- ✅ Single binary deployment
- ✅ <1MB memory overhead
- ✅ Fast cold start (<100ms)

## 🔐 Security

**Current State:**
- No authentication (local development only)

**For Production (Future):**
```go
// Option 1: Basic Auth
func basicAuth(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        user, pass, ok := r.BasicAuth()
        if !ok || user != "admin" || pass != "secret" {
            w.Header().Set("WWW-Authenticate", `Basic realm="picoclaw"`)
            http.Error(w, "Unauthorized", 401)
            return
        }
        next(w, r)
    }
}

// Option 2: JWT Token
// Option 3: OAuth2 (GitHub, Google)
```

**Recommended for SaaS:**
1. Add authentication middleware
2. Per-instance JWT tokens
3. Rate limiting
4. HTTPS enforcement
5. CORS configuration

## 🌍 SaaS Platform Architecture

### Multi-Tenant Setup

**Option 1: Container-per-User**
```yaml
# docker-compose.saas.yml
services:
  picoclaw-user-1:
    build: .
    environment:
      - USER_ID=user-1
      - INSTANCE_ID=abc123
    volumes:
      - ./data/user-1:/home/picoclaw/.picoclaw
    ports:
      - "8001:18790"
      - "8101:8080"  # Web UI

  picoclaw-user-2:
    build: .
    environment:
      - USER_ID=user-2
      - INSTANCE_ID=def456
    volumes:
      - ./data/user-2:/home/picoclaw/.picoclaw
    ports:
      - "8002:18790"
      - "8102:8080"
```

**Pros:**
- ✅ Complete isolation
- ✅ Easy per-user resource limits
- ✅ Simple to scale
- ✅ User-specific QR codes

**Cons:**
- ❌ More resource usage
- ❌ Container overhead per user

**Option 2: Shared Gateway with DB Isolation**
```
┌─────────────────────────────────────────────┐
│           Nginx Reverse Proxy               │
│   (Route by subdomain/path to instances)    │
└──────────┬──────────────────────────────────┘
           │
     ┌─────┴─────┬─────────┬─────────┐
     │           │         │         │
   User 1      User 2    User 3   User N
   Docker      Docker    Docker   Docker
```

**Dynamic Instance Creation:**
```go
func createUserInstance(userID string) error {
    // 1. Create user directory
    userDir := filepath.Join("/data", userID)
    os.MkdirAll(userDir, 0755)

    // 2. Generate config
    cfg := config.DefaultConfig()
    cfg.Channels.WhatsApp.Enabled = true
    config.SaveConfig(filepath.Join(userDir, "config.json"), cfg)

    // 3. Start Docker container
    cmd := exec.Command("docker", "run", "-d",
        "--name", "picoclaw-"+userID,
        "-v", userDir+":/home/picoclaw/.picoclaw",
        "-p", fmt.Sprintf("%d:18790", getAvailablePort()),
        "-p", fmt.Sprintf("%d:8080", getAvailablePort()+1000),
        "picoclaw:latest")

    return cmd.Run()
}
```

### SaaS Web UI Portal

```
Portal (SaaS Manager)
  ├── User Registration
  ├── Dashboard (Instance List)
  ├── Instance Manager
  │   ├── Create New Instance
  │   ├── Start/Stop/Restart
  │   └── View Logs
  └── Per-Instance Web UI (iframe/redirect)
```

**Example Portal Structure:**
```
saas-portal/
├── cmd/
│   └── portal/
│       └── main.go          # SaaS management server
├── pkg/
│   ├── instances/
│   │   ├── manager.go       # Docker instance manager
│   │   ├── provisioner.go   # Auto-provision new users
│   │   └── scaler.go        # Auto-scaling logic
│   ├── billing/             # Stripe integration
│   └── auth/                # User authentication
└── web/
    ├── dashboard.html       # User instance list
    └── instance_proxy.go    # Reverse proxy to user instances
```

## 📊 API Reference

### GET `/api/config`
Returns current configuration as JSON.

**Response:**
```json
{
  "channels": {
    "whatsapp": {
      "enabled": true,
      "allow_from": ["5511999999999@s.whatsapp.net"]
    }
  },
  "agents": {...},
  "providers": {...}
}
```

### POST `/api/config/save`
Save entire configuration.

**Request Body:**
```json
{
  "channels": {...},
  "agents": {...},
  "providers": {...}
}
```

### POST `/api/channels/toggle`
Enable/disable a channel.

**Request:**
```json
{
  "channel": "whatsapp",
  "enabled": true
}
```

### POST `/api/channels/whitelist`
Add/remove user from whitelist.

**Request:**
```json
{
  "channel": "whatsapp",
  "action": "add",
  "user_id": "5511999999999@s.whatsapp.net"
}
```

### GET `/api/skills/list`
List all installed skills.

### POST `/api/skills/install`
Install skill from GitHub.

**Request:**
```json
{
  "repository": "sipeed/picoclaw-skills/weather"
}
```

### POST `/api/skills/remove`
Remove installed skill.

**Request:**
```json
{
  "name": "weather"
}
```

### GET `/api/logs/stream` (SSE)
Stream logs in real-time using Server-Sent Events.

## 🚦 Production Deployment

### Nginx Reverse Proxy

```nginx
server {
    listen 80;
    server_name picoclaw.example.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }

    # SSE for logs
    location /api/logs/stream {
        proxy_pass http://localhost:8080;
        proxy_buffering off;
        proxy_cache off;
        proxy_set_header Connection '';
        chunked_transfer_encoding off;
    }
}
```

### Docker Production Setup

```yaml
version: '3.8'

services:
  picoclaw-gateway:
    image: picoclaw:latest
    restart: unless-stopped
    volumes:
      - ./data/config:/home/picoclaw/.picoclaw
    ports:
      - "18790:18790"

  picoclaw-webui:
    image: picoclaw:latest
    command: webui --host 0.0.0.0 --port 8080
    restart: unless-stopped
    volumes:
      - ./data/config:/home/picoclaw/.picoclaw
    ports:
      - "8080:8080"
    depends_on:
      - picoclaw-gateway

  nginx:
    image: nginx:alpine
    restart: unless-stopped
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - ./ssl:/etc/nginx/ssl
    depends_on:
      - picoclaw-webui
```

## 🎯 Future Enhancements

- [ ] User authentication (OAuth2, JWT)
- [ ] Multi-instance management dashboard
- [ ] QR code display in Web UI (convert string to image)
- [ ] Real-time message monitoring
- [ ] Statistics & analytics dashboard
- [ ] API key encryption at rest
- [ ] Two-factor authentication
- [ ] Role-based access control (RBAC)
- [ ] Audit logs
- [ ] Backup/restore functionality
- [ ] Dark mode toggle
- [ ] Mobile-responsive design
- [ ] WebSocket for bidirectional communication
- [ ] Prometheus metrics endpoint

## 📚 Development

### Adding a New Page

1. Create template in `pkg/webui/templates/newpage.html`:
```html
{{template "layout" .}}

{{define "content"}}
<div class="header">
    <h1>New Page</h1>
</div>

<div class="card">
    <!-- Your content -->
</div>
{{end}}
```

2. Add handler in `pkg/webui/server.go`:
```go
func (s *Server) handleNewPage(w http.ResponseWriter, r *http.Request) {
    data := struct {
        Config *config.Config
    }{
        Config: s.config,
    }
    s.templates.ExecuteTemplate(w, "newpage.html", data)
}
```

3. Register route:
```go
mux.HandleFunc("/newpage", s.handleNewPage)
```

### Adding a New API Endpoint

```go
func (s *Server) apiNewEndpoint(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var req struct {
        Field string `json:"field"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Process request...

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
```

## 🤝 Contributing

This Web UI is designed to be:
- **Lightweight** - No heavy frameworks
- **Fast** - Minimal JavaScript, server-side rendering
- **Simple** - Easy to understand and extend
- **Embeddable** - Templates embedded in binary

When contributing, please maintain these principles!

## 📝 License

Same as PicoClaw - MIT License

---

**Built with 🦐 by the PicoClaw Team**

Ready for SaaS? Let's scale! 🚀
