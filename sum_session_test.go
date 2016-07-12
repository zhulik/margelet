package margelet_test

import (
	"fmt"
	"strconv"

	"../margelet"
	"gopkg.in/telegram-bot-api.v4"
)

// SumSession - simple example session, that can sum numbers
type SumSession struct {
}

// HandleResponse - Handlers user response
func (s SumSession) HandleSession(bot margelet.MargeletAPI, message *tgbotapi.Message, session margelet.Session) (bool, error) {
	var msg tgbotapi.MessageConfig
	switch len(session.Responses()) {
	case 0:
		msg = tgbotapi.MessageConfig{Text: "Hello, please, write one number per message, after some iterations write 'end'."}
		s.response(bot, message, msg)
	default:
		if message.Text == "end" {
			var sum int
			for _, m := range session.Responses() {
				n, _ := strconv.Atoi(m.Text)
				sum += n
			}
			msg = tgbotapi.MessageConfig{Text: fmt.Sprintf("Your sum: %d", sum)}
			s.response(bot, message, msg)
			return true, nil
		}

		_, err := strconv.Atoi(message.Text)
		if err != nil {
			msg = tgbotapi.MessageConfig{Text: "Sorry, not a number"}
			s.response(bot, message, msg)
			return false, err
		}
	}

	return false, nil
}

// CancelResponse - Chance to clean up everything
func (s SumSession) CancelSession(bot margelet.MargeletAPI, message *tgbotapi.Message, session margelet.Session) {
	//Clean up all variables only used in the session

}

func (session SumSession) response(bot margelet.MargeletAPI, message *tgbotapi.Message, msg tgbotapi.MessageConfig) {
	msg.ChatID = message.Chat.ID
	msg.ReplyToMessageID = message.MessageID
	msg.ReplyMarkup = tgbotapi.ForceReply{ForceReply: true, Selective: true}
	bot.Send(msg)
}

// HelpMessage return help string for SumSession
func (session SumSession) HelpMessage() string {
	return "Sum your numbers and print result"
}
