package margelet_test

import (
	"github.com/Syfaro/telegram-bot-api"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/zhulik/margelet"
	"testing"
	"time"
)

func TestMargelet(t *testing.T) {
	Convey("When given margelet", t, func() {
		m := getMargelet()

		Convey("When adding new message responder", func() {
			m.AddMessageResponder(margelet.EchoResponder{})

			Convey("It should be aded to message responders array", func() {
				So(m.MessageResponders, ShouldNotBeEmpty)
			})
		})

		Convey("When adding new command responder", func() {
			m.AddCommandHandler("/test", margelet.HelpResponder{})

			Convey("It should be aded to command responders hash", func() {
				So(m.CommandResponders, ShouldNotBeEmpty)
			})

		})

		Convey("When adding new session handler", func() {
			m.AddSessionHandler("/test", margelet.SumSession{})

			Convey("It should be aded to command responders array", func() {
				So(m.SessionHandlers, ShouldNotBeEmpty)
			})

		})

		Convey("When sending message", func() {
			_, err := m.Send(tgbotapi.NewMessage(0, "TEST"))

			Convey("Error should be nil", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("When quick sending message", func() {
			_, err := m.QuickSend(0, "TEST")

			Convey("Error should be nil", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("When getting file direct URL", func() {
			_, err := m.GetFileDirectURL("test")

			Convey("Error should be nil", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("When asking is message sent to me", func() {
			msg := tgbotapi.Message{}
			res := m.IsMessageToMe(msg)
			Convey("It should return false", func() {
				So(res, ShouldBeFalse)
			})
		})

		Convey("When trying to get config repository", func() {
			repo := m.GetConfigRepository()
			Convey("It should not return nil", func() {
				So(repo, ShouldNotBeNil)
			})
		})

		Convey("Given configured margelet", func() {
			m.AddMessageResponder(margelet.EchoResponder{})
			m.AddSessionHandler("/sum", margelet.SumSession{})

			Convey("When running should handle message without panic", func() {
				go m.Run()

				botMock.Updates <- tgbotapi.Update{Message: tgbotapi.Message{Text: "Test"}}
				time.Sleep(10 * time.Millisecond)
				m.Stop()
			})

			Convey("When running should handle command without panic", func() {
				go m.Run()

				botMock.Updates <- tgbotapi.Update{Message: tgbotapi.Message{Text: "/help"}}

				time.Sleep(10 * time.Millisecond)
				m.Stop()
			})

			Convey("When running should handle session without panic", func() {
				go m.Run()

				botMock.Updates <- tgbotapi.Update{Message: tgbotapi.Message{Text: "/sum"}}
				botMock.Updates <- tgbotapi.Update{Message: tgbotapi.Message{Text: "10"}}
				botMock.Updates <- tgbotapi.Update{Message: tgbotapi.Message{Text: "20"}}
				botMock.Updates <- tgbotapi.Update{Message: tgbotapi.Message{Text: "test"}}
				botMock.Updates <- tgbotapi.Update{Message: tgbotapi.Message{Text: "end"}}

				time.Sleep(100 * time.Millisecond)
				m.Stop()
			})
		})
	})
}
