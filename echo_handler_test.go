package margelet_test

import (
	"../margelet"
	"gopkg.in/telegram-bot-api.v3"
)

// EchoHandler is simple handler example
type EchoHandler struct {
}

// Response send message back to author
func (handler EchoHandler) HandleMessage(bot margelet.MargeletAPI, message tgbotapi.Message) error {
	_, err := bot.Send(tgbotapi.NewMessage(message.Chat.ID, message.Text))
	return err
}
