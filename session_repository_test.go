package margelet_test

import (
	"github.com/zhulik/margelet"
	"reflect"
	"testing"
)

func TestAddSession(t *testing.T) {
	margelet.Redis.FlushDb()
	margelet.SessionRepo.Add(100, 500, "Test")

	if !reflect.DeepEqual(margelet.SessionRepo.Find(100, 500), []string{"Test"}) {
		t.Fail()
	}
}

func TestFindSessionWithWrongIds(t *testing.T) {
	margelet.Redis.FlushDb()
	margelet.SessionRepo.Add(100, 500, "Test")

	if !reflect.DeepEqual(margelet.SessionRepo.Find(100, 501), []string{}) {
		t.Fail()
	}
}

func TestRemoveSession(t *testing.T) {
	margelet.Redis.FlushDb()
	margelet.SessionRepo.Add(100, 500, "Test")
	margelet.SessionRepo.Remove(100, 500)

	if !reflect.DeepEqual(margelet.SessionRepo.Find(100, 500), []string{}) {
		t.Fail()
	}
}
