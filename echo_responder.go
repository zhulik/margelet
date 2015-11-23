package margelet

import (
	"github.com/Syfaro/telegram-bot-api"
)

type EchoResponder struct {
}

func (this EchoResponder) Response(bot MargeletAPI, message tgbotapi.Message) error {
	_, err := bot.Send(tgbotapi.NewMessage(message.Chat.ID, message.Text))
	return err
}
