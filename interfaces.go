package margelet

import (
	"github.com/Syfaro/telegram-bot-api"
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
	HandleResponse(bot MargeletAPI, message tgbotapi.Message, responses []string) (bool, error)
	HelpMessage() string
}

// MargeletAPI - interface, that describes margelet API
type MargeletAPI interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
	GetFileDirectURL(fileID string) (string, error)
	IsMessageToMe(message tgbotapi.Message) bool
}

// TGBotAPI - interface, thar describe telegram-bot-api API
type TGBotAPI interface {
	MargeletAPI
	GetUpdatesChan(config tgbotapi.UpdateConfig) (<-chan tgbotapi.Update, error)
}
