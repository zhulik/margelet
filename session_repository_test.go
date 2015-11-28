package margelet_test

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestSessionRepository(t *testing.T) {
	Convey("Given session repository", t, func() {
		m := getMargelet()

		Convey("When creating new session", func() {
			m.SessionRepository.Create(100, 500, "/test")

			Convey("New session should be found in repo", func() {
				So(m.SessionRepository.Command(100, 500), ShouldEqual, "/test")
			})
		})

		Convey("When adding new dialog", func() {
			m.SessionRepository.Add(100, 500, "Test")

			Convey("It shound be found in repo", func() {
				So(m.SessionRepository.Dialog(100, 500), ShouldResemble, []string{"Test"})
			})
		})

		Convey("When dialogs exists", func() {
			m.SessionRepository.Add(100, 500, "Test")

			Convey("Trying to get dialog for non-exists session shound return empty array", func() {
				m.SessionRepository.Add(100, 500, "Test")
				So(m.SessionRepository.Dialog(100, 501), ShouldResemble, []string{})

			})

			Convey("Removed session shound not be found in repo", func() {
				m.SessionRepository.Remove(100, 500)
				So(m.SessionRepository.Dialog(100, 500), ShouldResemble, []string{})
			})
		})
	})
}
