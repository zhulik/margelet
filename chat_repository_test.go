package margelet_test

import (
	"github.com/zhulik/margelet"
	"testing"
	"reflect"
)

func TestMain(m *testing.M) {
	margelet.InitRedis("127.0.0.1:6379", "password", 10)
	margelet.Redis.FlushDb()

	margelet.InitChatRepository("", margelet.Redis)
	m.Run()
}

func TestAddChat(t *testing.T) {
	margelet.Chat.AddChat(100500)

	if !reflect.DeepEqual(margelet.Chat.All(), []int{100500}) {
		t.Fail()
	}
}

func TestRemoveChat(t *testing.T) {
	margelet.Redis.FlushDb()

	margelet.Chat.AddChat(100500)
	margelet.Chat.AddChat(100501)
	margelet.Chat.RemoveChat(100500)

	if !reflect.DeepEqual(margelet.Chat.All(), []int{100501}) {
		t.Fail()
	}
}

func TestAll(t *testing.T) {
	margelet.Redis.FlushDb()

	margelet.Chat.AddChat(100500)
	margelet.Chat.AddChat(100501)
	margelet.Chat.AddChat(100502)
	margelet.Chat.AddChat(100503)

	t.Log(margelet.Chat.All())

	if !reflect.DeepEqual(margelet.Chat.All(), []int{100500, 100501, 100502, 100503}) {
		t.Fail()
	}
}
