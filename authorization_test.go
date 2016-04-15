package margelet_test

import (
	"testing"
	"time"

	"../margelet"

	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/telegram-bot-api.v4"
)

func TestAuthorization(t *testing.T) {
	Convey("When given margelet", t, func() {
		m := getMargelet()

		Convey("with registered command with auth policy", func() {
			m.AddCommandHandler("/test", margelet.HelpHandler{}, margelet.UsernameAuthorizationPolicy{Usernames: []string{"test"}})

			Convey("sending message from allowed user", func() {
				from := tgbotapi.User{UserName: "test"}
				chat := tgbotapi.Chat{ID: 1}

				go m.Run()
				botMock.Updates <- tgbotapi.Update{Message: &tgbotapi.Message{From: &from, Text: "/test", Chat: &chat}}

				time.Sleep(100 * time.Millisecond)
				m.Stop()
			})

			Convey("sending message from disallowed user", func() {
				from := tgbotapi.User{UserName: "another_user"}
				chat := tgbotapi.Chat{ID: 1}

				go m.Run()
				botMock.Updates <- tgbotapi.Update{Message: &tgbotapi.Message{From: &from, Text: "/test", Chat: &chat}}

				time.Sleep(100 * time.Millisecond)
				m.Stop()
			})
		})
	})
}
