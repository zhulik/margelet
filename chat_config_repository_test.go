package margelet_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
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
			Convey("Response should be empty", func() {
				So(m.ChatConfigRepository.Get(100500), ShouldBeEmpty)
			})
		})

		type testConfig struct {
			FavColor string
		}

		Convey("When adding new config using Struct", func() {
			testStruct := testConfig{FavColor: "Green"}
			m.ChatConfigRepository.SetWithStruct(100500, testStruct)

			Convey("It can be found in repository", func() {
				var testStruct2 testConfig
				m.ChatConfigRepository.GetWithStruct(100500, &testStruct2)
				So(testStruct2.FavColor, ShouldEqual, testStruct.FavColor)
			})
		})
	})
}
