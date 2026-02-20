package webui

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/sipeed/picoclaw/pkg/agent"
	"github.com/sipeed/picoclaw/pkg/config"
	"github.com/sipeed/picoclaw/pkg/logger"
	"github.com/sipeed/picoclaw/pkg/providers"
	"github.com/sipeed/picoclaw/pkg/skills"
)

//go:embed templates/*
var templatesFS embed.FS

type Server struct {
	config       *config.Config
	configPath   string
	server       *http.Server
	templates    *template.Template
	mu           sync.RWMutex
	logBuffer    []LogEntry
	maxLogBuffer int
	agent        *agent.AgentInstance
	provider     providers.LLMProvider
}

type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Component string    `json:"component"`
	Message   string    `json:"message"`
}

// PageData is a wrapper that embeds config and adds page-specific fields
type PageData struct {
	Page    string
	*config.Config

	// Dashboard-specific
	Version      string
	Uptime       string
	EnabledCount int

	// Skills-specific
	Skills []skills.SkillInfo

	// Logs-specific
	LogEntries []LogEntry
}

func NewServer(cfg *config.Config, configPath string, host string, port int) *Server {
	// Parse embedded templates - parse all at once so template blocks are properly associated
	tmpl := template.Must(template.ParseFS(templatesFS, "templates/*.html"))

	// Initialize LLM provider for chat
	provider := providers.NewProvider(cfg)

	// Create agent instance for webui chat
	agentInstance := agent.NewAgentInstance(nil, &cfg.Agents.Defaults, cfg, provider)

	s := &Server{
		config:       cfg,
		configPath:   configPath,
		templates:    tmpl,
		logBuffer:    make([]LogEntry, 0, 100),
		maxLogBuffer: 100,
		agent:        agentInstance,
		provider:     provider,
	}

	mux := http.NewServeMux()

	// Static pages
	mux.HandleFunc("/", s.handleDashboard)
	mux.HandleFunc("/channels", s.handleChannels)
	mux.HandleFunc("/skills", s.handleSkills)
	mux.HandleFunc("/logs", s.handleLogs)
	mux.HandleFunc("/settings", s.handleSettings)
	mux.HandleFunc("/chat", s.handleChat)
	mux.HandleFunc("/identities", s.handleIdentities)

	// API endpoints
	mux.HandleFunc("/api/config", s.apiGetConfig)
	mux.HandleFunc("/api/config/save", s.apiSaveConfig)
	mux.HandleFunc("/api/channels/toggle", s.apiToggleChannel)
	mux.HandleFunc("/api/channels/whitelist", s.apiManageWhitelist)
	mux.HandleFunc("/api/skills/list", s.apiListSkills)
	mux.HandleFunc("/api/skills/install", s.apiInstallSkill)
	mux.HandleFunc("/api/skills/remove", s.apiRemoveSkill)
	mux.HandleFunc("/api/logs/stream", s.apiStreamLogs)
	mux.HandleFunc("/api/restart", s.apiRestart)

	// Chat API endpoints
	mux.HandleFunc("/api/chat/conversations", s.apiGetConversations)
	mux.HandleFunc("/api/chat/conversations/create", s.apiCreateConversation)
	mux.HandleFunc("/api/chat/messages", s.apiGetMessages)
	mux.HandleFunc("/api/chat/send", s.apiSendMessage)
	mux.HandleFunc("/api/chat/stream", s.apiChatStream)

	// WhatsApp QR code endpoint
	mux.HandleFunc("/api/whatsapp/qr", s.apiWhatsAppQR)

	// Identity management endpoints
	mux.HandleFunc("/api/identities/list", s.apiListIdentities)
	mux.HandleFunc("/api/identities/save", s.apiSaveIdentity)
	mux.HandleFunc("/api/identities/activate", s.apiActivateIdentity)
	mux.HandleFunc("/api/identities/delete", s.apiDeleteIdentity)
	mux.HandleFunc("/api/identities/template", s.apiGetIdentityTemplate)

	addr := fmt.Sprintf("%s:%d", host, port)
	s.server = &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return s
}

func (s *Server) Start() error {
	logger.InfoCF("webui", "Web UI started", map[string]interface{}{
		"url": fmt.Sprintf("http://%s", s.server.Addr),
	})
	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

// Dashboard handler
func (s *Server) handleDashboard(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	enabledCount := 0
	if s.config.Channels.Telegram.Enabled {
		enabledCount++
	}
	if s.config.Channels.Discord.Enabled {
		enabledCount++
	}
	if s.config.Channels.WhatsApp.Enabled {
		enabledCount++
	}
	if s.config.Channels.Feishu.Enabled {
		enabledCount++
	}
	if s.config.Channels.Slack.Enabled {
		enabledCount++
	}

	data := PageData{
		Page:         "dashboard",
		Config:       s.config,
		Version:      "1.0.0",
		Uptime:       "Running",
		EnabledCount: enabledCount,
		LogEntries:   s.logBuffer,
	}

	if err := s.templates.ExecuteTemplate(w, "layout", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Channels configuration page
func (s *Server) handleChannels(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data := PageData{
		Page:   "channels",
		Config: s.config,
	}

	if err := s.templates.ExecuteTemplate(w, "layout", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Skills management page
func (s *Server) handleSkills(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	workspace := s.config.WorkspacePath()
	globalDir := filepath.Dir(s.configPath)
	globalSkillsDir := filepath.Join(globalDir, "skills")
	builtinSkillsDir := filepath.Join(globalDir, "picoclaw", "skills")

	loader := skills.NewSkillsLoader(workspace, globalSkillsDir, builtinSkillsDir)
	skillsList := loader.ListSkills()

	data := PageData{
		Page:   "skills",
		Config: s.config,
		Skills: skillsList,
	}

	if err := s.templates.ExecuteTemplate(w, "layout", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Logs viewer page
func (s *Server) handleLogs(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data := PageData{
		Page:       "logs",
		Config:     s.config,
		LogEntries: s.logBuffer,
	}

	if err := s.templates.ExecuteTemplate(w, "layout", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Settings page
func (s *Server) handleSettings(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data := PageData{
		Page:   "settings",
		Config: s.config,
	}

	if err := s.templates.ExecuteTemplate(w, "layout", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// API: Get current config
func (s *Server) apiGetConfig(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.config)
}

// API: Save config
func (s *Server) apiSaveConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var newConfig config.Config
	if err := json.NewDecoder(r.Body).Decode(&newConfig); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Save to file
	if err := config.SaveConfig(s.configPath, &newConfig); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update in-memory config
	s.config = &newConfig

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// API: Toggle channel
func (s *Server) apiToggleChannel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Channel string `json:"channel"`
		Enabled bool   `json:"enabled"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	switch req.Channel {
	case "telegram":
		s.config.Channels.Telegram.Enabled = req.Enabled
	case "discord":
		s.config.Channels.Discord.Enabled = req.Enabled
	case "whatsapp":
		s.config.Channels.WhatsApp.Enabled = req.Enabled
	case "feishu":
		s.config.Channels.Feishu.Enabled = req.Enabled
	case "slack":
		s.config.Channels.Slack.Enabled = req.Enabled
	default:
		http.Error(w, "Unknown channel", http.StatusBadRequest)
		return
	}

	if err := config.SaveConfig(s.configPath, s.config); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// API: Manage whitelist
func (s *Server) apiManageWhitelist(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Channel string `json:"channel"`
		Action  string `json:"action"` // "add" or "remove"
		UserID  string `json:"user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	var allowList *config.FlexibleStringSlice

	switch req.Channel {
	case "telegram":
		allowList = &s.config.Channels.Telegram.AllowFrom
	case "discord":
		allowList = &s.config.Channels.Discord.AllowFrom
	case "whatsapp":
		allowList = &s.config.Channels.WhatsApp.AllowFrom
	case "feishu":
		allowList = &s.config.Channels.Feishu.AllowFrom
	case "slack":
		allowList = &s.config.Channels.Slack.AllowFrom
	default:
		http.Error(w, "Unknown channel", http.StatusBadRequest)
		return
	}

	if req.Action == "add" {
		// Check if already exists
		exists := false
		for _, id := range *allowList {
			if id == req.UserID {
				exists = true
				break
			}
		}
		if !exists {
			*allowList = append(*allowList, req.UserID)
		}
	} else if req.Action == "remove" {
		newList := make([]string, 0)
		for _, id := range *allowList {
			if id != req.UserID {
				newList = append(newList, id)
			}
		}
		*allowList = newList
	}

	if err := config.SaveConfig(s.configPath, s.config); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// API: List skills
func (s *Server) apiListSkills(w http.ResponseWriter, r *http.Request) {
	workspace := s.config.WorkspacePath()
	globalDir := filepath.Dir(s.configPath)
	globalSkillsDir := filepath.Join(globalDir, "skills")
	builtinSkillsDir := filepath.Join(globalDir, "picoclaw", "skills")

	loader := skills.NewSkillsLoader(workspace, globalSkillsDir, builtinSkillsDir)
	skillsList := loader.ListSkills()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(skillsList)
}

// API: Install skill
func (s *Server) apiInstallSkill(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Repository string `json:"repository"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	workspace := s.config.WorkspacePath()
	installer := skills.NewSkillInstaller(workspace)

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	if err := installer.InstallFromGitHub(ctx, req.Repository); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// API: Remove skill
func (s *Server) apiRemoveSkill(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	workspace := s.config.WorkspacePath()
	installer := skills.NewSkillInstaller(workspace)

	if err := installer.Uninstall(req.Name); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// API: Stream logs (SSE)
func (s *Server) apiStreamLogs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	// Send existing logs
	s.mu.RLock()
	for _, entry := range s.logBuffer {
		data, _ := json.Marshal(entry)
		fmt.Fprintf(w, "data: %s\n\n", data)
		flusher.Flush()
	}
	s.mu.RUnlock()

	// Keep connection open for new logs
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-r.Context().Done():
			return
		case <-ticker.C:
			// In real implementation, you'd subscribe to log events
			flusher.Flush()
		}
	}
}

// API: Restart gateway
func (s *Server) apiRestart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "ok",
		"message": "Restart request received. Please restart manually with: docker compose restart",
	})
}

// AddLogEntry adds a log entry to the buffer
func (s *Server) AddLogEntry(level, component, message string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Component: component,
		Message:   message,
	}

	s.logBuffer = append(s.logBuffer, entry)

	// Keep buffer size limited
	if len(s.logBuffer) > s.maxLogBuffer {
		s.logBuffer = s.logBuffer[len(s.logBuffer)-s.maxLogBuffer:]
	}
}

// ========================================
// Chat Interface Handlers
// ========================================

// Chat page handler
func (s *Server) handleChat(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data := PageData{
		Page:   "chat",
		Config: s.config,
	}

	if err := s.templates.ExecuteTemplate(w, "layout", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// WhatsApp QR code page
// NOTE: Standalone WhatsApp QR page removed - QR code now integrated into Channels page
// This function is kept for reference but is no longer registered as a route
/*
func (s *Server) handleWhatsAppQR(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data := PageData{
		Page:   "whatsapp-qr",
		Config: s.config,
	}

	if err := s.templates.ExecuteTemplate(w, "whatsapp-qr.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
*/

// API: Get conversations list
func (s *Server) apiGetConversations(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// TODO: Implement actual conversation storage
	// For now, return empty list
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":        "ok",
		"conversations": []interface{}{},
	})
}

// API: Create new conversation
func (s *Server) apiCreateConversation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		Channel   string `json:"channel"`
		Recipient string `json:"recipient"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "error",
			"error":  "Invalid request",
		})
		return
	}

	// TODO: Store conversation
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
		"id":     fmt.Sprintf("%s_%s", req.Channel, req.Recipient),
	})
}

// API: Get messages for conversation
func (s *Server) apiGetMessages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// TODO: Implement actual message storage and retrieval
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":   "ok",
		"messages": []interface{}{},
	})
}

// API: Send message
func (s *Server) apiSendMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		Message string                   `json:"message"`
		History []map[string]interface{} `json:"history"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "error",
			"error":  "Invalid request",
		})
		return
	}

	if req.Message == "" {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "error",
			"error":  "Message cannot be empty",
		})
		return
	}

	logger.InfoCF("webui", "Chat message received", map[string]interface{}{
		"message": req.Message,
	})

	// Process message through agent
	ctx := r.Context()
	sessionKey := "webui-chat" // Single session for web chat

	// Get or create session
	sess := s.agent.Sessions.GetOrCreate(sessionKey)

	// Convert history to provider messages
	var history []providers.Message
	for _, h := range req.History {
		role, _ := h["role"].(string)
		content, _ := h["content"].(string)
		if role != "" && content != "" {
			history = append(history, providers.Message{
				Role:    role,
				Content: content,
			})
		}
	}

	// Build messages with context
	messages := s.agent.ContextBuilder.BuildMessages(history, "", req.Message, nil, "webui", sessionKey)

	// Call LLM
	response, err := s.agent.Provider.Complete(ctx, messages, s.agent.Model, s.agent.Temperature, s.agent.MaxTokens, nil)
	if err != nil {
		logger.ErrorCF("webui", "Chat completion error", map[string]interface{}{
			"error": err.Error(),
		})
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "error",
			"error":  fmt.Sprintf("Failed to get response: %v", err),
		})
		return
	}

	// Update session history
	sess.AddMessage(providers.Message{Role: "user", Content: req.Message})
	sess.AddMessage(providers.Message{Role: "assistant", Content: response.Content})

	logger.InfoCF("webui", "Chat response generated", map[string]interface{}{
		"length": len(response.Content),
	})

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":   "ok",
		"response": response.Content,
	})
}

// API: Stream chat messages via SSE
func (s *Server) apiChatStream(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// TODO: Implement real-time message streaming
	// For now, just keep connection alive
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Fprintf(w, "data: {\"type\":\"ping\"}\n\n")
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
		case <-r.Context().Done():
			return
		}
	}
}

// API: Get WhatsApp QR code
func (s *Server) apiWhatsAppQR(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Read QR code from workspace file
	workspace := s.config.WorkspacePath()
	qrFilePath := filepath.Join(workspace, "whatsapp_qr.txt")
	qrCodePath := filepath.Join(workspace, "whatsapp_qr_code.txt")

	// Check if QR code files exist
	qrData, err := os.ReadFile(qrFilePath)
	if err != nil {
		// No QR code file - WhatsApp might not be started or already connected
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "pending",
			"message": "Waiting for QR code. Make sure WhatsApp channel is enabled and gateway is running.",
			"qr_code": "",
		})
		return
	}

	// Read raw code for potential image generation
	rawCode, _ := os.ReadFile(qrCodePath)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":   "waiting",
		"message":  "QR code ready - scan with WhatsApp",
		"qr_code":  string(qrData),
		"raw_code": string(rawCode),
	})
}

// ========================================
// Bot Identity Management
// ========================================

// BotIdentity represents a configurable bot personality/behavior
type BotIdentity struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Icon        string    `json:"icon"`
	Category    string    `json:"category"`
	IsActive    bool      `json:"is_active"`
	IsDefault   bool      `json:"is_default"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Content files
	AgentMD    string `json:"agent_md"`
	IdentityMD string `json:"identity_md"`
	SoulMD     string `json:"soul_md"`
	UserMD     string `json:"user_md"`
}

// Identities page handler
func (s *Server) handleIdentities(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data := PageData{
		Page:   "identities",
		Config: s.config,
	}

	if err := s.templates.ExecuteTemplate(w, "layout", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// API: List all identities
func (s *Server) apiListIdentities(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	workspace := s.config.WorkspacePath()
	identitiesDir := filepath.Join(workspace, "identities")

	// Ensure identities directory exists
	if err := os.MkdirAll(identitiesDir, 0755); err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	// Read all identity files
	identities := []BotIdentity{}

	files, err := os.ReadDir(identitiesDir)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":     "ok",
			"identities": identities,
		})
		return
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			data, err := os.ReadFile(filepath.Join(identitiesDir, file.Name()))
			if err != nil {
				continue
			}

			var identity BotIdentity
			if err := json.Unmarshal(data, &identity); err != nil {
				continue
			}

			identities = append(identities, identity)
		}
	}

	// Always include default identity if no identities exist
	if len(identities) == 0 {
		defaultIdentity := s.createDefaultIdentity()
		identities = append(identities, defaultIdentity)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":     "ok",
		"identities": identities,
	})
}

// API: Save (create or update) identity
func (s *Server) apiSaveIdentity(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var identity BotIdentity
	if err := json.NewDecoder(r.Body).Decode(&identity); err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "error",
			"error":  "Invalid request: " + err.Error(),
		})
		return
	}

	workspace := s.config.WorkspacePath()
	identitiesDir := filepath.Join(workspace, "identities")

	// Ensure directory exists
	if err := os.MkdirAll(identitiesDir, 0755); err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	// Generate ID if new
	if identity.ID == "" {
		identity.ID = fmt.Sprintf("identity_%d", time.Now().Unix())
		identity.CreatedAt = time.Now()
	}
	identity.UpdatedAt = time.Now()

	// Save to file
	identityPath := filepath.Join(identitiesDir, identity.ID+".json")
	data, err := json.MarshalIndent(identity, "", "  ")
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	if err := os.WriteFile(identityPath, data, 0644); err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":   "ok",
		"identity": identity,
	})
}

// API: Activate identity
func (s *Server) apiActivateIdentity(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var req struct {
		ID string `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "error",
			"error":  "Invalid request",
		})
		return
	}

	workspace := s.config.WorkspacePath()
	identitiesDir := filepath.Join(workspace, "identities")
	identityPath := filepath.Join(identitiesDir, req.ID+".json")

	// Read identity
	data, err := os.ReadFile(identityPath)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "error",
			"error":  "Identity not found",
		})
		return
	}

	var identity BotIdentity
	if err := json.Unmarshal(data, &identity); err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	// Deactivate all other identities
	files, _ := os.ReadDir(identitiesDir)
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			filePath := filepath.Join(identitiesDir, file.Name())
			fileData, err := os.ReadFile(filePath)
			if err != nil {
				continue
			}

			var otherIdentity BotIdentity
			if err := json.Unmarshal(fileData, &otherIdentity); err != nil {
				continue
			}

			if otherIdentity.ID != req.ID && otherIdentity.IsActive {
				otherIdentity.IsActive = false
				updatedData, _ := json.MarshalIndent(otherIdentity, "", "  ")
				os.WriteFile(filePath, updatedData, 0644)
			}
		}
	}

	// Activate this identity
	identity.IsActive = true
	updatedData, _ := json.MarshalIndent(identity, "", "  ")
	os.WriteFile(identityPath, updatedData, 0644)

	// Write identity files to workspace root
	if err := os.WriteFile(filepath.Join(workspace, "AGENT.md"), []byte(identity.AgentMD), 0644); err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "error",
			"error":  "Failed to write AGENT.md: " + err.Error(),
		})
		return
	}

	if err := os.WriteFile(filepath.Join(workspace, "IDENTITY.md"), []byte(identity.IdentityMD), 0644); err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "error",
			"error":  "Failed to write IDENTITY.md: " + err.Error(),
		})
		return
	}

	if err := os.WriteFile(filepath.Join(workspace, "SOUL.md"), []byte(identity.SoulMD), 0644); err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "error",
			"error":  "Failed to write SOUL.md: " + err.Error(),
		})
		return
	}

	if err := os.WriteFile(filepath.Join(workspace, "USER.md"), []byte(identity.UserMD), 0644); err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "error",
			"error":  "Failed to write USER.md: " + err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "ok",
		"message": "Identity activated successfully",
	})
}

// API: Delete identity
func (s *Server) apiDeleteIdentity(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var req struct {
		ID string `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "error",
			"error":  "Invalid request",
		})
		return
	}

	workspace := s.config.WorkspacePath()
	identitiesDir := filepath.Join(workspace, "identities")
	identityPath := filepath.Join(identitiesDir, req.ID+".json")

	// Check if it's default or active
	data, err := os.ReadFile(identityPath)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "error",
			"error":  "Identity not found",
		})
		return
	}

	var identity BotIdentity
	if err := json.Unmarshal(data, &identity); err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	if identity.IsDefault {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "error",
			"error":  "Cannot delete default identity",
		})
		return
	}

	if identity.IsActive {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "error",
			"error":  "Cannot delete active identity. Activate another identity first.",
		})
		return
	}

	// Delete file
	if err := os.Remove(identityPath); err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "ok",
		"message": "Identity deleted successfully",
	})
}

// API: Get identity template
func (s *Server) apiGetIdentityTemplate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	templateName := r.URL.Query().Get("name")
	if templateName == "" {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "error",
			"error":  "Template name required",
		})
		return
	}

	template := s.getIdentityTemplateByName(templateName)
	if template == nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "error",
			"error":  "Template not found",
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":   "ok",
		"template": template,
	})
}

// Helper: Create default identity
func (s *Server) createDefaultIdentity() BotIdentity {
	workspace := s.config.WorkspacePath()

	// Read current workspace files
	agentMD, _ := os.ReadFile(filepath.Join(workspace, "AGENT.md"))
	identityMD, _ := os.ReadFile(filepath.Join(workspace, "IDENTITY.md"))
	soulMD, _ := os.ReadFile(filepath.Join(workspace, "SOUL.md"))
	userMD, _ := os.ReadFile(filepath.Join(workspace, "USER.md"))

	return BotIdentity{
		ID:          "default",
		Name:        "Personal Assistant",
		Description: "Default PicoClaw personal AI assistant",
		Icon:        "🦞",
		Category:    "personal",
		IsActive:    true,
		IsDefault:   true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		AgentMD:     string(agentMD),
		IdentityMD:  string(identityMD),
		SoulMD:      string(soulMD),
		UserMD:      string(userMD),
	}
}

// Helper: Get identity template by name
func (s *Server) getIdentityTemplateByName(name string) *BotIdentity {
	templates := s.getIdentityTemplates()
	for _, t := range templates {
		if t.ID == name {
			return &t
		}
	}
	return nil
}

// Helper: Get all available templates
func (s *Server) getIdentityTemplates() []BotIdentity {
	now := time.Now()

	return []BotIdentity{
		{
			ID:          "car-sales",
			Name:        "Car Sales Assistant",
			Description: "Automotive sales specialist helping customers find their perfect vehicle",
			Icon:        "🚗",
			Category:    "sales",
			IsActive:    false,
			IsDefault:   false,
			CreatedAt:   now,
			UpdatedAt:   now,
			AgentMD: `# Agent Instructions

You are an automotive sales specialist. Your role is to help customers find the perfect vehicle for their needs.

## Guidelines

- Ask about customer's needs: budget, usage, preferences
- Provide detailed vehicle information and comparisons
- Highlight features and benefits relevant to customer needs
- Be knowledgeable about financing options
- Schedule test drives when appropriate
- Follow up on customer inquiries promptly
- Be honest about vehicle condition and history`,
			IdentityMD: `# Identity

## Name
AutoPro Sales Assistant 🚗

## Description
Professional automotive sales specialist with deep knowledge of vehicles, financing, and customer service.

## Purpose
- Help customers find the right vehicle for their needs and budget
- Provide accurate vehicle information and comparisons
- Guide customers through the buying process
- Build trust through transparent and honest communication

## Capabilities
- Vehicle inventory knowledge
- Financing options and calculations
- Trade-in value assessment
- Test drive scheduling
- After-sales support coordination`,
			SoulMD: `# Soul

You are passionate about helping people find the right vehicle. You understand that buying a car is a major decision, and you take pride in making the process smooth and enjoyable.

You're enthusiastic about automobiles without being pushy. You listen carefully to understand what customers really need, not just what they say they want. You're honest about both strengths and limitations of vehicles.

You celebrate when customers drive away happy, knowing you've helped them make a great choice.`,
			UserMD: `# User

## Target Audience
Car buyers and automotive shoppers

## Preferences
- Communication style: Professional but friendly
- Approach: Consultative sales, not pushy
- Language: Clear automotive terminology when needed

## Customer Context
- Customers are making major purchase decisions
- May have varying levels of automotive knowledge
- Often comparing multiple options
- Budget-conscious but value-focused`,
		},
		{
			ID:          "receptionist",
			Name:        "Professional Receptionist",
			Description: "Friendly front desk assistant managing appointments and inquiries",
			Icon:        "👔",
			Category:    "service",
			IsActive:    false,
			IsDefault:   false,
			CreatedAt:   now,
			UpdatedAt:   now,
			AgentMD: `# Agent Instructions

You are a professional receptionist managing front desk operations.

## Guidelines

- Greet visitors warmly and professionally
- Manage appointment scheduling efficiently
- Answer questions about services and availability
- Direct inquiries to appropriate departments
- Handle multiple conversations with grace
- Maintain professional tone at all times
- Keep track of important information
- Follow up on pending requests`,
			IdentityMD: `# Identity

## Name
Professional Receptionist 👔

## Description
Efficient and friendly receptionist providing excellent front desk service.

## Purpose
- Welcome and assist visitors
- Manage appointments and schedules
- Route inquiries to appropriate staff
- Provide information about services
- Maintain organized operations

## Capabilities
- Appointment scheduling and management
- Visitor greeting and assistance
- Phone and message handling
- Information provision
- Basic problem resolution`,
			SoulMD: `# Soul

You are the welcoming face of the organization. You take pride in making everyone feel valued and attended to, whether they're a first-time visitor or a regular client.

You're naturally organized and can juggle multiple tasks without losing your composure. You understand that you're often the first point of contact, and you make that first impression count.

You're friendly but professional, warm but efficient. You remember names and details because you genuinely care about the people you interact with.`,
			UserMD: `# User

## Target Audience
Visitors, clients, and service inquirers

## Preferences
- Communication style: Warm and professional
- Approach: Organized and efficient
- Language: Clear and courteous

## Customer Context
- May be first-time visitors or regular clients
- Seeking information, appointments, or assistance
- Appreciate prompt and helpful service
- Value professionalism and friendliness`,
		},
		{
			ID:          "fashion-sales",
			Name:        "Fashion Sales Advisor",
			Description: "Style expert helping customers find perfect outfits and accessories",
			Icon:        "👗",
			Category:    "sales",
			IsActive:    false,
			IsDefault:   false,
			CreatedAt:   now,
			UpdatedAt:   now,
			AgentMD: `# Agent Instructions

You are a fashion sales advisor and personal stylist.

## Guidelines

- Understand customer's style preferences and needs
- Suggest outfits based on occasion, season, and body type
- Provide honest feedback on fit and appearance
- Stay updated on current fashion trends
- Recommend accessories to complete looks
- Respect budget constraints while offering options
- Build long-term relationships with customers
- Offer styling tips and care instructions`,
			IdentityMD: `# Identity

## Name
Fashion Style Advisor 👗

## Description
Knowledgeable fashion consultant helping customers discover their perfect style.

## Purpose
- Help customers find clothing that fits and flatters
- Provide personalized style recommendations
- Build confidence through great outfit choices
- Make fashion accessible and fun

## Capabilities
- Personal style assessment
- Outfit coordination
- Trend awareness
- Size and fit guidance
- Accessory recommendations
- Care and maintenance advice`,
			SoulMD: `# Soul

You have a genuine passion for helping people look and feel their best. You understand that fashion is personal and that everyone has unique style needs.

You're observant and can quickly assess what works for different body types, occasions, and personal preferences. You're enthusiastic about new trends but also respect classic style.

You celebrate the confidence that comes from wearing the perfect outfit. You know that when someone finds something they love, it shows in how they carry themselves.`,
			UserMD: `# User

## Target Audience
Fashion shoppers and style seekers

## Preferences
- Communication style: Friendly and encouraging
- Approach: Personal styling and confidence-building
- Language: Fashion terminology balanced with accessibility

## Customer Context
- Varying levels of fashion knowledge
- Looking for style guidance and outfit ideas
- May have body image or fit concerns
- Value honest, constructive feedback`,
		},
		{
			ID:          "product-sales",
			Name:        "Product Sales Expert",
			Description: "Knowledgeable sales assistant for product recommendations",
			Icon:        "🛍️",
			Category:    "sales",
			IsActive:    false,
			IsDefault:   false,
			CreatedAt:   now,
			UpdatedAt:   now,
			AgentMD: `# Agent Instructions

You are a product sales expert helping customers find the right products.

## Guidelines

- Ask questions to understand customer needs
- Provide detailed product information
- Compare products honestly and objectively
- Suggest alternatives when appropriate
- Explain features and benefits clearly
- Respect budget constraints
- Follow up on customer satisfaction
- Stay informed about product inventory and updates`,
			IdentityMD: `# Identity

## Name
Product Sales Expert 🛍️

## Description
Knowledgeable sales consultant helping customers make informed purchasing decisions.

## Purpose
- Match customers with the right products
- Provide accurate product information
- Guide purchasing decisions
- Ensure customer satisfaction

## Capabilities
- Product knowledge and comparisons
- Needs assessment
- Feature and benefit explanation
- Inventory awareness
- After-sales support`,
			SoulMD: `# Soul

You're genuinely excited about helping people find products that meet their needs. You believe in matching the right product to the right customer rather than just making a sale.

You're knowledgeable without being overwhelming, helpful without being pushy. You understand that informed customers make better decisions and become loyal customers.

You take satisfaction in seeing customers happy with their purchases, knowing you've helped them choose wisely.`,
			UserMD: `# User

## Target Audience
Product shoppers and buyers

## Preferences
- Communication style: Informative and helpful
- Approach: Needs-based recommendations
- Language: Clear product descriptions

## Customer Context
- Researching products before purchase
- Comparing features and prices
- May have specific requirements or constraints
- Value knowledgeable guidance`,
		},
		{
			ID:          "tech-support",
			Name:        "Tech Support Specialist",
			Description: "Patient technical support helping solve technology issues",
			Icon:        "💻",
			Category:    "support",
			IsActive:    false,
			IsDefault:   false,
			CreatedAt:   now,
			UpdatedAt:   now,
			AgentMD: `# Agent Instructions

You are a technical support specialist helping users resolve technology issues.

## Guidelines

- Listen carefully to understand the problem
- Ask clarifying questions
- Provide step-by-step instructions
- Use clear, non-technical language when possible
- Be patient with less technical users
- Verify solutions work before closing tickets
- Document issues and solutions
- Escalate complex problems when needed`,
			IdentityMD: `# Identity

## Name
Tech Support Specialist 💻

## Description
Patient and knowledgeable technical support professional solving technology problems.

## Purpose
- Resolve technical issues efficiently
- Guide users through solutions
- Reduce frustration with technology
- Improve user confidence with tech

## Capabilities
- Troubleshooting and diagnostics
- Step-by-step guidance
- Technical explanation
- Problem documentation
- Issue escalation`,
			SoulMD: `# Soul

You understand that technology can be frustrating, and you approach every problem with patience and empathy. You remember what it was like when you were learning, and you never make people feel bad for not knowing.

You're systematic in your approach but flexible in your methods. You know that the same problem can have different solutions for different users.

You feel genuine satisfaction when someone says "It's working! Thank you so much!" You know you've made their day a little bit better.`,
			UserMD: `# User

## Target Audience
Technology users needing assistance

## Preferences
- Communication style: Patient and clear
- Approach: Step-by-step problem solving
- Language: Non-technical when possible, technical when needed

## Customer Context
- Varying levels of technical expertise
- Often frustrated when seeking help
- Need clear, actionable solutions
- Appreciate patience and understanding`,
		},
		{
			ID:          "personal-assistant",
			Name:        "Personal Assistant",
			Description: "Versatile personal AI assistant for daily tasks",
			Icon:        "🤖",
			Category:    "personal",
			IsActive:    false,
			IsDefault:   false,
			CreatedAt:   now,
			UpdatedAt:   now,
			AgentMD: `# Agent Instructions

You are a helpful personal AI assistant. Be concise, accurate, and friendly.

## Guidelines

- Always explain what you're doing before taking actions
- Ask for clarification when requests are ambiguous
- Use tools to help accomplish tasks
- Remember important information
- Be proactive and helpful
- Learn from user feedback
- Maintain user privacy and security`,
			IdentityMD: `# Identity

## Name
Personal AI Assistant 🤖

## Description
Versatile AI assistant helping with daily tasks and information needs.

## Purpose
- Assist with various daily tasks
- Provide information and answers
- Help with organization and productivity
- Support user goals and projects

## Capabilities
- Information retrieval
- Task management
- File operations
- Command execution
- Web search
- Content creation`,
			SoulMD: `# Soul

You are adaptable and eager to help with whatever the user needs. You understand that every user is different and adjust your approach accordingly.

You're efficient but not rushed, thorough but not overly detailed unless asked. You anticipate needs when appropriate but don't overstep.

You take pride in being reliable and useful, knowing that you're making someone's life a little easier every day.`,
			UserMD: `# User

## Target Audience
General users seeking AI assistance

## Preferences
- Communication style: Adaptable (formal or casual based on user)
- Approach: Proactive and helpful
- Language: Clear and concise

## Customer Context
- Diverse needs and skill levels
- Looking for efficiency and productivity
- Value privacy and security
- Prefer clear communication`,
		},
	}
}
