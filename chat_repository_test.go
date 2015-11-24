package margelet_test

import (
	"github.com/zhulik/margelet"
	"reflect"
	"testing"
)

func TestMain(m *testing.M) {
	margelet.InitRedis("127.0.0.1:6379", "", 0)
	margelet.Redis.FlushDb()

	margelet.InitChatRepository("", margelet.Redis)
	margelet.InitSessionRepository("", margelet.Redis)
	m.Run()
}

func TestAddChat(t *testing.T) {
	margelet.ChatRepo.Add(100500)

	if !reflect.DeepEqual(margelet.ChatRepo.All(), []int{100500}) {
		t.Fail()
	}
}

func TestRemoveChat(t *testing.T) {
	margelet.Redis.FlushDb()

	margelet.ChatRepo.Add(100500)
	margelet.ChatRepo.Add(100501)
	margelet.ChatRepo.Remove(100500)

	if !reflect.DeepEqual(margelet.ChatRepo.All(), []int{100501}) {
		t.Fail()
	}
}

func TestAll(t *testing.T) {
	margelet.Redis.FlushDb()

	margelet.ChatRepo.Add(100500)
	margelet.ChatRepo.Add(100501)
	margelet.ChatRepo.Add(100502)
	margelet.ChatRepo.Add(100503)

	t.Log(margelet.ChatRepo.All())

	if !reflect.DeepEqual(margelet.ChatRepo.All(), []int{100500, 100501, 100502, 100503}) {
		t.Fail()
	}
}
