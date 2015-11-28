package margelet_test

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestChatConfigRepository(t *testing.T) {
	Convey("Given ChatConfigRepository", t, func() {
		m := getMargelet()

		Convey("When adding new config", func() {
			m.ChatConfigRepository.Set(100500, "{\"a\": 100}")

			Convey("It can be found in repository", func() {
				So(m.ChatConfigRepository.Get(100500), ShouldEqual, "{\"a\": 100}")
			})
		})

		Convey("With existing config", func() {
			m.ChatConfigRepository.Set(100500, "{\"a\": 100}")

			Convey("When removing config", func() {
				m.ChatConfigRepository.Remove(100500)
				Convey("It cannot be found in repository", func() {
					So(m.ChatConfigRepository.Get(100500), ShouldBeEmpty)
				})
			})
		})

		Convey("When trying to get missing chatID", func() {
			Convey("Response should be empty", func(){
				So(m.ChatConfigRepository.Get(100500), ShouldBeEmpty)
			})
		})
	})
}
