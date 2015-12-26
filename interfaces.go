package margelet

import (
	"github.com/Syfaro/telegram-bot-api"
	"gopkg.in/redis.v3"
)

// Responder - interface for message responders
type Responder interface {
	Response(bot MargeletAPI, message tgbotapi.Message) error
}

// CommandHandler - interface for command handlers
type CommandHandler interface {
	Responder
	HelpMessage() string
}

// SessionHandler - interface for session handlers
type SessionHandler interface {
	HandleResponse(bot MargeletAPI, message tgbotapi.Message, responses []tgbotapi.Message) (bool, error)
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
	GetSessionRepository() *SessionRepository
	GetRedis() *redis.Client
	HandleSession(message tgbotapi.Message, handler SessionHandler)
}

// TGBotAPI - interface, thar describe telegram-bot-api API
type TGBotAPI interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
	GetFileDirectURL(fileID string) (string, error)
	IsMessageToMe(message tgbotapi.Message) bool
	GetUpdatesChan(config tgbotapi.UpdateConfig) (<-chan tgbotapi.Update, error)
}
