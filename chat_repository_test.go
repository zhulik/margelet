package margelet_test

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestChatRepository(t *testing.T) {
	Convey("Given chat repository", t, func() {
		m := getMargelet()

		Convey("When new chat added", func() {
			m.ChatRepository.Add(100500)
			Convey("New chat should be found in repo", func() {
				So(m.ChatRepository.All(), ShouldResemble, []int{100500})
			})
		})

		Convey("Given existing chats", func() {
			m.ChatRepository.Add(100500)
			m.ChatRepository.Add(100501)
			m.ChatRepository.Add(100502)
			m.ChatRepository.Add(100503)

			Convey("When removing chat", func() {
				m.ChatRepository.Remove(100500)

				Convey("Removed chat should not be found in repo", func() {
					So(m.ChatRepository.All(), ShouldResemble, []int{100501, 100502, 100503})
				})
			})

			Convey("Repo should return all chats", func() {
				So(m.ChatRepository.All(), ShouldResemble, []int{100500, 100501, 100502, 100503})
			})
		})
	})
}
