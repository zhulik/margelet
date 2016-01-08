package margelet

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"gopkg.in/redis.v3"
)

// Handler - interface for message handlers
type MessageHandler interface {
	HandleMessage(bot MargeletAPI, message tgbotapi.Message) error
}

// CommandHandler - interface for command handlers
type CommandHandler interface {
	HandleCommand(bot MargeletAPI, message tgbotapi.Message) error
	HelpMessage() string
}

// SessionHandler - interface for session handlers
type SessionHandler interface {
	HandleSession(bot MargeletAPI, message tgbotapi.Message, responses []tgbotapi.Message) (bool, error)
	HelpMessage() string
}

// MargeletAPI - interface, that describes margelet API
type MargeletAPI interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
	QuickSend(chatID int, message string) (tgbotapi.Message, error)
	QuickReply(chatID, messageID int, message string) (tgbotapi.Message, error)
	GetFileDirectURL(fileID string) (string, error)
	IsMessageToMe(message tgbotapi.Message) bool
	GetConfigRepository() *ChatConfigRepository
	GetSessionRepository() SessionRepository
	GetRedis() *redis.Client
	HandleSession(message tgbotapi.Message, handler SessionHandler)
}

// TGBotAPI - interface, that describe telegram-bot-api API
type TGBotAPI interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
	GetFileDirectURL(fileID string) (string, error)
	IsMessageToMe(message tgbotapi.Message) bool
	GetUpdatesChan(config tgbotapi.UpdateConfig) (<-chan tgbotapi.Update, error)
}

// AuthorizationPolicy - interface, that describes authorization policy for command or session
type AuthorizationPolicy interface {
	Allow(message tgbotapi.Message) error
}
