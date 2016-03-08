package margelet_test

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/zhulik/margelet"
	"gopkg.in/telegram-bot-api.v2"
	"testing"
	"time"
)

func TestMargelet(t *testing.T) {
	Convey("When given margelet", t, func() {
		m := getMargelet()

		Convey("When adding new message handler", func() {
			m.AddMessageHandler(EchoHandler{})

			Convey("It should be aded to message handlers array", func() {
				So(m.MessageHandlers, ShouldNotBeEmpty)
			})
		})

		Convey("When adding new command handler", func() {
			m.AddCommandHandler("/test", margelet.HelpHandler{})

			Convey("It should be aded to command handler hash", func() {
				So(m.CommandHandlers, ShouldNotBeEmpty)
			})

		})

		Convey("When adding new session handler", func() {
			m.AddSessionHandler("/test", SumSession{})

			Convey("It should be aded to command handler array", func() {
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

		Convey("When quick replying message", func() {
			_, err := m.QuickReply(0, 100500, "TEST")

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
			m.AddMessageHandler(EchoHandler{})
			m.AddMessageHandler(PanicHandler{})
			m.AddSessionHandler("/sum", SumSession{})
			m.InlineHandler = &InlineImage{}

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

			Convey("When running panic should not crash bot", func() {
				go m.Run()

				botMock.Updates <- tgbotapi.Update{Message: tgbotapi.Message{Text: "/panic"}}

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

			Convey("When running should handle inline query without panic", func() {
				go m.Run()

				botMock.Updates <-  tgbotapi.Update{InlineQuery: tgbotapi.InlineQuery{ID: "test_id", Query: "test"}}

				time.Sleep(100 * time.Millisecond)
				m.Stop()
			})
		})
	})
}
