package margelet

import (
	"github.com/zhulik/telegram-bot-api"
)

type Responder interface {
	Response(message tgbotapi.Message) (tgbotapi.Chattable, error)
}
