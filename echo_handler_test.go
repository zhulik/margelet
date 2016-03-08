package margelet_test

import (
	"github.com/zhulik/margelet"
	"gopkg.in/telegram-bot-api.v2"
)

// EchoHandler is simple handler example
type EchoHandler struct {
}

// Response send message back to author
func (handler EchoHandler) HandleMessage(bot margelet.MargeletAPI, message tgbotapi.Message) error {
	_, err := bot.Send(tgbotapi.NewMessage(message.Chat.ID, message.Text))
	return err
}
