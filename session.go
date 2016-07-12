package margelet

import (
	"gopkg.in/telegram-bot-api.v4"
)

type margeletSession struct {
	bot           MargeletAPI
	chatID        int64
	userID        int
	lastMessageID int
	responses     []tgbotapi.Message
}

func newMargetletSession(bot MargeletAPI, message *tgbotapi.Message, responses []tgbotapi.Message) Session {
	return &margeletSession{
		bot:           bot,
		chatID:        message.Chat.ID,
		userID:        message.From.ID,
		lastMessageID: message.MessageID,
		responses:     responses,
	}
}

func (s *margeletSession) Responses() []tgbotapi.Message {
	return s.responses
}

func (s *margeletSession) QuickSend(text string) (tgbotapi.Message, error) {
	return s.bot.QuickSend(s.chatID, text)
}

func (s *margeletSession) QuckReply(text string) (tgbotapi.Message, error) {
	return s.bot.QuickReply(s.chatID, s.lastMessageID, text)
}
