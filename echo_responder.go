package margelet

import (
	"github.com/zhulik/telegram-bot-api"
)

type EchoResponder struct {
}

func (this EchoResponder) Response(bot *Margelet, message tgbotapi.Message) (tgbotapi.Chattable, error) {
	return tgbotapi.NewMessage(message.Chat.ID, message.Text), nil
}
