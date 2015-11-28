package margelet_test

import (
	"github.com/zhulik/margelet"
	"reflect"
	"testing"
)

func TestAddChat(t *testing.T) {
	getMargelet()
	margelet.ChatRepo.Add(100500)

	if !reflect.DeepEqual(margelet.ChatRepo.All(), []int{100500}) {
		t.Fail()
	}
}

func TestRemoveChat(t *testing.T) {
	getMargelet()

	margelet.ChatRepo.Add(100500)
	margelet.ChatRepo.Add(100501)
	margelet.ChatRepo.Remove(100500)

	if !reflect.DeepEqual(margelet.ChatRepo.All(), []int{100501}) {
		t.Fail()
	}
}

func TestAll(t *testing.T) {
	getMargelet()

	margelet.ChatRepo.Add(100500)
	margelet.ChatRepo.Add(100501)
	margelet.ChatRepo.Add(100502)
	margelet.ChatRepo.Add(100503)

	if !reflect.DeepEqual(margelet.ChatRepo.All(), []int{100500, 100501, 100502, 100503}) {
		t.Fail()
	}
}
