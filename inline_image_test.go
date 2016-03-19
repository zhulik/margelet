package margelet_test

import (
	"../margelet"
	"gopkg.in/telegram-bot-api.v2"
)

type InlineImage struct {
}

func (handler InlineImage) HandleInline(bot margelet.MargeletAPI, query tgbotapi.InlineQuery) error {
	testPhotoQuery := tgbotapi.NewInlineQueryResultPhoto(query.ID, "https://telegram.org/img/t_logo.png")
	testPhotoQuery.ThumbURL = "https://telegram.org/img/t_logo.png"

	config := tgbotapi.InlineConfig{
		InlineQueryID: query.ID,
		CacheTime:     2,
		IsPersonal:    false,
		Results:       []interface{}{testPhotoQuery},
		NextOffset:    "",
	}

	bot.AnswerInlineQuery(config)
	return nil
}
