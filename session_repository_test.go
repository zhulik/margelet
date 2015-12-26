package margelet_test

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"github.com/Syfaro/telegram-bot-api"
)

func TestSessionRepository(t *testing.T) {
	Convey("Given session repository", t, func() {
		m := getMargelet()

		msg := tgbotapi.Message{Text: "Test"}

		Convey("When creating new session", func() {
			m.SessionRepository.Create(100, 500, "/test")

			Convey("New session should be found in repo", func() {
				So(m.SessionRepository.Command(100, 500), ShouldEqual, "/test")
			})
		})

		Convey("When adding new dialog", func() {
			m.SessionRepository.Add(100, 500, msg)

			Convey("It shound be found in repo", func() {
				So(m.SessionRepository.Dialog(100, 500)[0].Text, ShouldEqual, "Test")
			})
		})

		Convey("When dialogs exists", func() {
			m.SessionRepository.Add(100, 500, msg)

			Convey("Trying to get dialog for non-exists session shound return empty array", func() {
				m.SessionRepository.Add(100, 500, msg)
				So(m.SessionRepository.Dialog(100, 501), ShouldBeEmpty)

			})

			Convey("Removed session shound not be found in repo", func() {
				m.SessionRepository.Remove(100, 500)
				So(m.SessionRepository.Dialog(100, 500), ShouldBeEmpty)
			})
		})
	})
}
