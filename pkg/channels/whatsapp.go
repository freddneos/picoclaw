package channels

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"

	waProto "go.mau.fi/whatsmeow/binary/proto"
	qrcode "github.com/skip2/go-qrcode"

	"github.com/sipeed/picoclaw/pkg/bus"
	"github.com/sipeed/picoclaw/pkg/config"
	"github.com/sipeed/picoclaw/pkg/logger"
	"github.com/sipeed/picoclaw/pkg/utils"

	_ "github.com/mattn/go-sqlite3"
)

type WhatsAppChannel struct {
	*BaseChannel
	client    *whatsmeow.Client
	config    config.WhatsAppConfig
	container *sqlstore.Container
	dbPath    string
}

func NewWhatsAppChannel(cfg *config.Config, bus *bus.MessageBus) (*WhatsAppChannel, error) {
	whatsappCfg := cfg.Channels.WhatsApp
	base := NewBaseChannel("whatsapp", whatsappCfg, bus, whatsappCfg.AllowFrom)

	// Create database path
	dbPath := filepath.Join(cfg.WorkspacePath(), "whatsapp.db")

	return &WhatsAppChannel{
		BaseChannel: base,
		config:      whatsappCfg,
		dbPath:      dbPath,
	}, nil
}

func (c *WhatsAppChannel) Start(ctx context.Context) error {
	logger.InfoC("whatsapp", "Starting WhatsApp channel...")

	// Ensure database directory exists
	dbDir := filepath.Dir(c.dbPath)
	if err := os.MkdirAll(dbDir, 0700); err != nil {
		return fmt.Errorf("failed to create database directory: %w", err)
	}

	// Create database container
	dbLog := waLog.Stdout("Database", "INFO", true)
	container, err := sqlstore.New(ctx, "sqlite3", "file:"+c.dbPath+"?_foreign_keys=on", dbLog)
	if err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}
	c.container = container

	// Get first device or create new
	deviceStore, err := container.GetFirstDevice(ctx)
	if err != nil {
		return fmt.Errorf("failed to get device: %w", err)
	}

	clientLog := waLog.Stdout("Client", "INFO", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)
	c.client = client

	// Add event handler
	client.AddEventHandler(c.handleEvent)

	// Connect
	if client.Store.ID == nil {
		// No ID stored, need to pair
		qrChan, err := client.GetQRChannel(ctx)
		if err != nil {
			return fmt.Errorf("failed to get QR channel: %w", err)
		}

		err = client.Connect()
		if err != nil {
			return fmt.Errorf("failed to connect: %w", err)
		}

		logger.InfoC("whatsapp", "")
		logger.InfoC("whatsapp", "═══════════════════════════════════════════════════════════════")
		logger.InfoC("whatsapp", "  Scan the QR code below with WhatsApp to link your device:")
		logger.InfoC("whatsapp", "  1. Open WhatsApp on your phone")
		logger.InfoC("whatsapp", "  2. Tap Settings → Linked Devices → Link a Device")
		logger.InfoC("whatsapp", "  3. Scan the QR code below")
		logger.InfoC("whatsapp", "═══════════════════════════════════════════════════════════════")
		logger.InfoC("whatsapp", "")

		for evt := range qrChan {
			if evt.Event == "code" {
				// Generate terminal QR code
				qr, err := qrcode.New(evt.Code, qrcode.Medium)
				if err != nil {
					logger.ErrorCF("whatsapp", "Failed to generate QR code", map[string]interface{}{"error": err})
					logger.InfoC("whatsapp", "QR Code string: "+evt.Code)
					continue
				}

				// Print QR code as ASCII art
				qrString := qr.ToSmallString(false)
				logger.InfoC("whatsapp", "\n"+qrString)
				logger.InfoC("whatsapp", "")

				// Write QR code to file for Web UI
				workspace := filepath.Dir(c.dbPath)
				qrFilePath := filepath.Join(workspace, "whatsapp_qr.txt")
				os.WriteFile(qrFilePath, []byte(qrString), 0644)

				// Also write the raw code for image generation
				qrCodePath := filepath.Join(workspace, "whatsapp_qr_code.txt")
				os.WriteFile(qrCodePath, []byte(evt.Code), 0644)
			} else {
				logger.InfoCF("whatsapp", "QR channel result", map[string]interface{}{"event": evt.Event})
			}
		}
	} else {
		// Already paired
		err = client.Connect()
		if err != nil {
			return fmt.Errorf("failed to connect: %w", err)
		}
	}

	c.setRunning(true)
	logger.InfoC("whatsapp", "WhatsApp connected successfully")

	return nil
}

func (c *WhatsAppChannel) Stop(ctx context.Context) error {
	logger.InfoC("whatsapp", "Stopping WhatsApp channel...")

	if c.client != nil {
		c.client.Disconnect()
	}

	if c.container != nil {
		if err := c.container.Close(); err != nil {
			logger.ErrorCF("whatsapp", "Failed to close database", map[string]interface{}{"error": err})
		}
	}

	c.setRunning(false)
	return nil
}

func (c *WhatsAppChannel) Send(ctx context.Context, msg bus.OutboundMessage) error {
	if c.client == nil {
		return fmt.Errorf("whatsapp client not connected")
	}

	// Parse JID from chat ID
	jid, err := types.ParseJID(msg.ChatID)
	if err != nil {
		return fmt.Errorf("invalid chat ID: %w", err)
	}

	// Send message
	_, err = c.client.SendMessage(ctx, jid, &waProto.Message{
		Conversation: proto.String(msg.Content),
	})

	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	logger.InfoCF("whatsapp", "Sent message", map[string]interface{}{"chat_id": msg.ChatID})
	return nil
}

func (c *WhatsAppChannel) handleEvent(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		c.handleMessage(v)
	case *events.Connected:
		logger.InfoC("whatsapp", "Connected to WhatsApp")
	case *events.Disconnected:
		logger.InfoC("whatsapp", "Disconnected from WhatsApp")
	case *events.LoggedOut:
		logger.InfoC("whatsapp", "Logged out from WhatsApp")
	}
}

func (c *WhatsAppChannel) handleMessage(evt *events.Message) {
	// Ignore messages from self
	if evt.Info.IsFromMe {
		return
	}

	// Get message text
	var content string
	if evt.Message.Conversation != nil {
		content = *evt.Message.Conversation
	} else if evt.Message.ExtendedTextMessage != nil && evt.Message.ExtendedTextMessage.Text != nil {
		content = *evt.Message.ExtendedTextMessage.Text
	}

	// Ignore empty messages
	if content == "" {
		return
	}

	senderID := evt.Info.Sender.String()
	chatID := evt.Info.Chat.String()

	logger.InfoCF("whatsapp", "Message received", map[string]interface{}{
		"sender":  senderID,
		"content": utils.Truncate(content, 50),
	})

	// Handle message through base channel
	metadata := map[string]string{
		"message_id": evt.Info.ID,
		"timestamp":  fmt.Sprintf("%d", evt.Info.Timestamp.Unix()),
	}

	if evt.Info.PushName != "" {
		metadata["user_name"] = evt.Info.PushName
	}

	c.HandleMessage(senderID, chatID, content, nil, metadata)
}
