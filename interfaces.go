package margelet

import (
	"gopkg.in/redis.v3"
	"gopkg.in/telegram-bot-api.v2"
)

// MessageHandler - interface for message handlers
type MessageHandler interface {
	HandleMessage(bot MargeletAPI, message tgbotapi.Message) error
}

// InlineHandler - interface for message handlers
type InlineHandler interface {
	HandleInline(bot MargeletAPI, query tgbotapi.InlineQuery) error
}

// CommandHandler - interface for command handlers
type CommandHandler interface {
	HandleCommand(bot MargeletAPI, message tgbotapi.Message) error
	HelpMessage() string
}

// SessionHandler - interface for session handlers
type SessionHandler interface {
	HandleSession(bot MargeletAPI, message tgbotapi.Message, responses []tgbotapi.Message) (bool, error)
	CancelSession(bot MargeletAPI, message tgbotapi.Message, responses []tgbotapi.Message)
	HelpMessage() string
}

// MargeletAPI - interface, that describes margelet API
type MargeletAPI interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
	AnswerInlineQuery(config tgbotapi.InlineConfig) (tgbotapi.APIResponse, error)
	QuickSend(chatID int, message string) (tgbotapi.Message, error)
	QuickReply(chatID, messageID int, message string) (tgbotapi.Message, error)
	GetFileDirectURL(fileID string) (string, error)
	IsMessageToMe(message tgbotapi.Message) bool
	GetConfigRepository() *ChatConfigRepository
	GetSessionRepository() SessionRepository
	GetRedis() *redis.Client
	HandleSession(message tgbotapi.Message, command string)
}

// TGBotAPI - interface, that describe telegram-bot-api API
type TGBotAPI interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
	AnswerInlineQuery(config tgbotapi.InlineConfig) (tgbotapi.APIResponse, error)
	GetFileDirectURL(fileID string) (string, error)
	IsMessageToMe(message tgbotapi.Message) bool
	GetUpdatesChan(config tgbotapi.UpdateConfig) (<-chan tgbotapi.Update, error)
}

// AuthorizationPolicy - interface, that describes authorization policy for command or session
type AuthorizationPolicy interface {
	Allow(message tgbotapi.Message) error
}
