package margelet_test

import (
	"../margelet"
	"gopkg.in/telegram-bot-api.v4"
)

type CallbackMessage struct {
}

func (handler CallbackMessage) HandleCallback(query margelet.CallbackQuery) error {
	config := tgbotapi.CallbackConfig{
		CallbackQueryID: query.Query().ID,
		Text:            "Done!",
		ShowAlert:       false,
	}

	query.Bot().AnswerCallbackQuery(config)
	return nil
}
