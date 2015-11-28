package margelet_test

import (
	"github.com/Syfaro/telegram-bot-api"
	"github.com/zhulik/margelet"
	"testing"
	"time"
)

func TestAddMessageResponder(t *testing.T) {
	m := getMargelet()
	m.AddMessageResponder(margelet.EchoResponder{})

	if len(m.MessageResponders) != 1 {
		t.Fail()
	}
}

func TestAddCommandResponder(t *testing.T) {
	m := getMargelet()
	m.AddCommandHandler("/test", margelet.EchoResponder{})

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
	_, err := m.GetFileDirectURL("test")

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
	m.AddCommandHandler("/test", margelet.EchoResponder{})
	m.AddMessageResponder(margelet.EchoResponder{})
	go m.Run()
	botMock.Updates <- tgbotapi.Update{Message: tgbotapi.Message{Text: "/test"}}
	botMock.Updates <- tgbotapi.Update{Message: tgbotapi.Message{Text: "Test"}}
	time.Sleep(1 * time.Second)
	m.Stop()
}
