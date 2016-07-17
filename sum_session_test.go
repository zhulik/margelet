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
func (s SumSession) HandleSession(session margelet.Session) error {
	switch len(session.Responses()) {
	case 0:
		session.QuickReply("Hello, please, write one number per message, after some iterations write 'end'.")
	default:
		if session.Message().Text == "end" {
			var sum int
			for _, m := range session.Responses() {
				n, _ := strconv.Atoi(m.Text)
				sum += n
			}
			session.QuickReply(fmt.Sprintf("Your sum: %d", sum))
			session.Finish()
			return nil
		}

		_, err := strconv.Atoi(session.Message().Text)
		if err != nil {
			session.QuickReply("Sorry, not a number")
			return err
		}
	}

	return nil
}

// CancelResponse - Chance to clean up everything
func (s SumSession) CancelSession(session margelet.Session) {
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
