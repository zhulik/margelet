package margelet_test

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/zhulik/margelet"
	"testing"
)

func TestSessionRepository(t *testing.T) {
	Convey("Given session repository", t, func() {
		getMargelet()

		Convey("When creating new session", func() {
			margelet.SessionRepo.Create(100, 500, "/test")

			Convey("New session should be found in repo", func() {
				So(margelet.SessionRepo.Command(100, 500), ShouldEqual, "/test")
			})
		})

		Convey("When adding new dialog", func() {
			margelet.SessionRepo.Add(100, 500, "Test")

			Convey("It shound be found in repo", func() {
				So(margelet.SessionRepo.Dialog(100, 500), ShouldResemble, []string{"Test"})
			})
		})

		Convey("When dialogs exists", func() {
			margelet.SessionRepo.Add(100, 500, "Test")

			Convey("Trying to get dialog for non-exists session shound return empty array", func() {
				margelet.SessionRepo.Add(100, 500, "Test")
				So(margelet.SessionRepo.Dialog(100, 501), ShouldResemble, []string{})

			})

			Convey("Removed session shound not be found in repo", func() {
				margelet.SessionRepo.Remove(100, 500)
				So(margelet.SessionRepo.Dialog(100, 500), ShouldResemble, []string{})
			})
		})
	})
}
