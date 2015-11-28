package margelet

import (
	"github.com/Syfaro/telegram-bot-api"
)

type Responder interface {
	Response(bot MargeletAPI, message tgbotapi.Message) error
}

type CommandHandler Responder

type SessionHandler interface {
	HandleResponse(bot MargeletAPI, message tgbotapi.Message, responses []string) (bool, error)
}

type MargeletAPI interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
	GetFileDirectURL(fileID string) (string, error)
	IsMessageToMe(message tgbotapi.Message) bool
}

type TGBotAPI interface {
	MargeletAPI
	GetUpdatesChan(config tgbotapi.UpdateConfig) (<-chan tgbotapi.Update, error)
}
