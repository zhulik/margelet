package margelet_test

import (
	"github.com/Syfaro/telegram-bot-api"
	"github.com/zhulik/margelet"
)

// EchoResponder is simple responder example
type EchoResponder struct {
}

// Response send message back to author
func (responder EchoResponder) Response(bot margelet.MargeletAPI, message tgbotapi.Message) error {
	_, err := bot.Send(tgbotapi.NewMessage(message.Chat.ID, message.Text))
	return err
}
