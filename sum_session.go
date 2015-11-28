package margelet

import (
	"fmt"
	"github.com/Syfaro/telegram-bot-api"
	"strconv"
)

// SumSession - simple example session, that cat sum numbers
type SumSession struct {
}

// HandleResponse - Handlers user response
func (session SumSession) HandleResponse(bot MargeletAPI, message tgbotapi.Message, responses []string) (bool, error) {
	var msg tgbotapi.MessageConfig
	switch len(responses) {
	case 0:
		msg = tgbotapi.MessageConfig{Text: "Hello, please, write one number per message, after some iterations write 'end'."}
		msg.ReplyMarkup = tgbotapi.ForceReply{true, true}
	default:
		if message.Text == "end" {
			var sum int
			for _, a := range responses {
				n, _ := strconv.Atoi(a)
				sum += n
			}
			msg = tgbotapi.MessageConfig{Text: fmt.Sprintf("Your sum: %d", sum)}
			session.response(bot, message, msg)
			msg.ReplyMarkup = tgbotapi.ForceReply{false, true}
			return true, nil
		}

		_, err := strconv.Atoi(message.Text)
		if err != nil {
			msg = tgbotapi.MessageConfig{Text: "Sorry, not a number"}
			session.response(bot, message, msg)
			msg.ReplyMarkup = tgbotapi.ForceReply{true, true}
			return false, err
		}
	}

	session.response(bot, message, msg)
	return false, nil
}

func (session SumSession) response(bot MargeletAPI, message tgbotapi.Message, msg tgbotapi.MessageConfig) {
	msg.ChatID = message.Chat.ID
	msg.ReplyToMessageID = message.MessageID
	bot.Send(msg)
}
