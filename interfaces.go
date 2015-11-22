package margelet

import (
	"github.com/zhulik/telegram-bot-api"
)

type Responder interface {
	Response(bot *Margelet, message tgbotapi.Message) error
}

type CommandHandler Responder

type TGBotAPI interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
	GetFileDirectURL(fileID string) (string, error)
	IsMessageToMe(message tgbotapi.Message) bool
	GetUpdatesChan(config tgbotapi.UpdateConfig) (<-chan tgbotapi.Update, error)
}
