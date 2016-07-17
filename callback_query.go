package margelet

import (
	"gopkg.in/telegram-bot-api.v4"
)

type callbackQuery struct {
	*message
	query *tgbotapi.CallbackQuery
}

func newCallbackQuery(bot MargeletAPI, query *tgbotapi.CallbackQuery) *callbackQuery {
	return &callbackQuery{
		message: newMessage(bot, query.Message),
		query:   query,
	}
}

// Query returns tgbotapi query
func (s *callbackQuery) Query() *tgbotapi.CallbackQuery {
	return s.query
}

func (s *callbackQuery) Data() string {
	return s.query.Data
}
