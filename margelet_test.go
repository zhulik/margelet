package margelet_test

import (
	"github.com/zhulik/margelet"
	"github.com/zhulik/telegram-bot-api"
	"testing"
	"time"
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
	m, _ := margelet.NewMargeletFromBot(&botMock)
	return m
}

func TestAddMessageResponder(t *testing.T) {
	m := getMargelet()
	m.AddMessageResponder(margelet.EchoResponder{})

	if len(m.MessageResponders) != 1 {
		t.Fail()
	}
}

func TestAddCommandResponder(t *testing.T) {
	m := getMargelet()
	m.AddCommandResponder("/test", margelet.EchoResponder{})

	if len(m.CommandResponders) != 1 {
		t.Fail()
	}
}

func TestSend(t *testing.T) {
	m := getMargelet()
	_, err := m.Send(tgbotapi.NewMessage(0, "TEST"))

	if err != nil {
		t.Fail()
	}
}

func TestGetFileDirectUrl(t *testing.T) {
	m := getMargelet()
	_, err := m.GetFileDirectUrl("test")

	if err != nil {
		t.Fail()
	}
}

func TestIsMessageToMe(t *testing.T) {
	m := getMargelet()
	m.IsMessageToMe(tgbotapi.Message{})

	if m.IsMessageToMe(tgbotapi.Message{}) != false {
		t.Fail()
	}
}

func TestRun(t *testing.T) {
	m := getMargelet()
	m.AddCommandResponder("/test", margelet.EchoResponder{})
	m.AddMessageResponder(margelet.EchoResponder{})
	go m.Run()
	botMock.Updates <- tgbotapi.Update{Message: tgbotapi.Message{Text: "/test"}}
	botMock.Updates <- tgbotapi.Update{Message: tgbotapi.Message{Text: "Test"}}
	time.Sleep(1 * time.Second)
	m.Stop()
}
