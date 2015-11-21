package margelet

import (
	"github.com/zhulik/telegram-bot-api"
)

type Responder interface {
	Response(bot *Margelet, message tgbotapi.Message) (tgbotapi.Chattable, error)
}

type TGBotAPI interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
	GetFileDirectUrl(fileID string) (string, error)
	IsMessageToMe(message tgbotapi.Message) bool
	GetUpdatesChan(config tgbotapi.UpdateConfig) (<-chan tgbotapi.Update, error)
}
