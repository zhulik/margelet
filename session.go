package margelet

import (
	"gopkg.in/telegram-bot-api.v4"
)

// MargeletSession incapsulates user's session
type margeletSession struct {
	*message
	bot         MargeletAPI
	lastMessage *tgbotapi.Message
	responses   []tgbotapi.Message
	finished    bool
}

func newMargetletSession(bot MargeletAPI, msg *tgbotapi.Message, responses []tgbotapi.Message) *margeletSession {
	return &margeletSession{
		message:   newMessage(bot, msg),
		responses: responses,
		finished:  false,
	}
}

// Responses returns all user's responses in session
func (s *margeletSession) Responses() []tgbotapi.Message {
	return s.responses
}

func (s *margeletSession) Finish() {
	s.finished = true
}
