package margelet_test

import (
	"github.com/Syfaro/telegram-bot-api"
	"github.com/zhulik/margelet"
)

type BotMock struct {
	Updates chan tgbotapi.Update
}

func (bot BotMock) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	return tgbotapi.Message{}, nil
}

func (bot BotMock) GetFileDirectURL(fileID string) (string, error) {
	return "https://example.com/test.txt", nil
}

func (bot BotMock) IsMessageToMe(message tgbotapi.Message) bool {
	return false
}

func (bot BotMock) GetUpdatesChan(config tgbotapi.UpdateConfig) (<-chan tgbotapi.Update, error) {
	return this.Updates, nil
}

var (
	botMock = BotMock{}
)

func getMargelet() *margelet.Margelet {
	botMock.Updates = make(chan tgbotapi.Update, 10)
	m, _ := margelet.NewMargeletFromBot("test", "127.0.0.1:6379", "", 10, &botMock)

	m.Redis.FlushDb()
	return m
}

func ExampleUsage() {
	bot, err := margelet.NewMargelet("<your awesome bot name>", "<redis addr>", "<redis password>", 0, "your bot token", false)

	if err != nil {
		panic(err)
	}

	bot.Run()
}
