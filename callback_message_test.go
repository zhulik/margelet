package margelet_test

import (
	"../margelet"
	"gopkg.in/telegram-bot-api.v4"
)

type CallbackMessage struct {
}

func (handler CallbackMessage) HandleCallback(bot margelet.MargeletAPI, query *tgbotapi.CallbackQuery) error {

	config := tgbotapi.CallbackConfig{
		CallbackQueryID: query.ID,
		Text:            "Done!",
		ShowAlert:       false,
	}

	bot.AnswerCallbackQuery(config)
	return nil
}
