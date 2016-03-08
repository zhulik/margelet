package margelet_test

import (
	"fmt"
	"github.com/zhulik/margelet"
	"gopkg.in/telegram-bot-api.v2"
	"strconv"
)

// SumSession - simple example session, that can sum numbers
type SumSession struct {
}

// HandleResponse - Handlers user response
func (session SumSession) HandleSession(bot margelet.MargeletAPI, message tgbotapi.Message, responses []tgbotapi.Message) (bool, error) {
	var msg tgbotapi.MessageConfig
	switch len(responses) {
	case 0:
		msg = tgbotapi.MessageConfig{Text: "Hello, please, write one number per message, after some iterations write 'end'."}
		msg.ReplyMarkup = tgbotapi.ForceReply{ForceReply: true, Selective: true}
	default:
		if message.Text == "end" {
			var sum int
			for _, m := range responses {
				n, _ := strconv.Atoi(m.Text)
				sum += n
			}
			msg = tgbotapi.MessageConfig{Text: fmt.Sprintf("Your sum: %d", sum)}
			session.response(bot, message, msg)
			msg.ReplyMarkup = tgbotapi.ForceReply{ForceReply: true, Selective: true}
			return true, nil
		}

		_, err := strconv.Atoi(message.Text)
		if err != nil {
			msg = tgbotapi.MessageConfig{Text: "Sorry, not a number"}
			session.response(bot, message, msg)
			msg.ReplyMarkup = tgbotapi.ForceReply{ForceReply: true, Selective: true}
			return false, err
		}
	}

	session.response(bot, message, msg)
	return false, nil
}

func (session SumSession) response(bot margelet.MargeletAPI, message tgbotapi.Message, msg tgbotapi.MessageConfig) {
	msg.ChatID = message.Chat.ID
	msg.ReplyToMessageID = message.MessageID
	bot.Send(msg)
}

// HelpMessage return help string for SumSession
func (session SumSession) HelpMessage() string {
	return "Sum your numbers and print result"
}
