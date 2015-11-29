package margelet_test

import (
	"github.com/Syfaro/telegram-bot-api"
	"github.com/zhulik/margelet"
)

type BotMock struct {
	Updates chan tgbotapi.Update
}

func (this BotMock) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	return tgbotapi.Message{}, nil
}

func (this BotMock) GetFileDirectURL(fileID string) (string, error) {
	return "https://example.com/test.txt", nil
}

func (this BotMock) IsMessageToMe(message tgbotapi.Message) bool {
	return false
}

func (this BotMock) GetUpdatesChan(config tgbotapi.UpdateConfig) (<-chan tgbotapi.Update, error) {
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