package margelet_test

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/zhulik/margelet"
	"testing"
)

func TestChatRepository(t *testing.T) {
	Convey("Given chat repository", t, func() {
		getMargelet()

		Convey("When new chat added", func() {
			margelet.ChatRepo.Add(100500)
			Convey("New chat should be found in repo", func() {
				So(margelet.ChatRepo.All(), ShouldResemble, []int{100500})
			})
		})

		Convey("Given existing chats", func() {
			margelet.ChatRepo.Add(100500)
			margelet.ChatRepo.Add(100501)
			margelet.ChatRepo.Add(100502)
			margelet.ChatRepo.Add(100503)

			Convey("When removing chat", func() {
				margelet.ChatRepo.Remove(100500)

				Convey("Removed chat should not be found in repo", func() {
					So(margelet.ChatRepo.All(), ShouldResemble, []int{100501, 100502, 100503})
				})
			})

			Convey("Repo should return all chats", func() {
				So(margelet.ChatRepo.All(), ShouldResemble, []int{100500, 100501, 100502, 100503})
			})
		})
	})
}
