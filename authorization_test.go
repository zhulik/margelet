package margelet_test

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/zhulik/margelet"
	"testing"
	"time"
)

func TestAuthorization(t *testing.T) {
	Convey("When given margelet", t, func() {
		m := getMargelet()

		Convey("with registered command with auth policy", func() {
			m.AddCommandHandler("/test", margelet.HelpHandler{}, margelet.UsernameAuthorizationPolicy{Usernames: []string{"test"}})

			Convey("sending message from allowed user", func(){
				from := tgbotapi.User{UserName: "test"}

				go m.Run()
				botMock.Updates <- tgbotapi.Update{Message: tgbotapi.Message{From: from, Text: "/test"}}

				time.Sleep(100 * time.Millisecond)
				m.Stop()
			})

			Convey("sending message from disallowed user", func(){
				from := tgbotapi.User{UserName: "another_user"}

				go m.Run()
				botMock.Updates <- tgbotapi.Update{Message: tgbotapi.Message{From: from, Text: "/test"}}

				time.Sleep(100 * time.Millisecond)
				m.Stop()
			})
		})
	})
}
