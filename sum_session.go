package margelet

import (
	"fmt"
	"github.com/Syfaro/telegram-bot-api"
	"strconv"
)

type SumSession struct {
}

func (this SumSession) HandleResponse(bot MargeletAPI, message tgbotapi.Message, responses []string) (bool, error) {
	var msg tgbotapi.MessageConfig
	switch len(responses) {
	case 0:
		msg = tgbotapi.MessageConfig{Text: "Hello, please, write one number per message, after some iterations write 'end'."}
	default:
		if message.Text == "end" {
			var sum int
			for _, a := range responses {
				n, _ := strconv.Atoi(a)
				sum += n
			}
			msg = tgbotapi.MessageConfig{Text: fmt.Sprintf("Your sum: %d", sum)}
			this.response(bot, message, msg)
			return true, nil
		} else {
			_, err := strconv.Atoi(message.Text)
			if err != nil {
				msg = tgbotapi.MessageConfig{Text: "Sorry, not a number"}
				this.response(bot, message, msg)
				return false, err
			}
		}
	}

	this.response(bot, message, msg)
	return false, nil
}

func (this SumSession) response(bot MargeletAPI, message tgbotapi.Message, msg tgbotapi.MessageConfig) {
	msg.ChatID = message.Chat.ID
	msg.ReplyToMessageID = message.MessageID
	bot.Send(msg)
}
