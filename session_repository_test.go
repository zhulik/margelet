package margelet_test

import (
	"github.com/zhulik/margelet"
	"reflect"
	"testing"
)

func TestCreateSession(t *testing.T) {
	getMargelet()

	margelet.SessionRepo.Create(100, 500, "/test")

	if margelet.SessionRepo.Command(100, 500) != "/test" {
		t.Fail()
	}
}

func TestAddDialog(t *testing.T) {
	getMargelet()
	margelet.SessionRepo.Add(100, 500, "Test")

	if !reflect.DeepEqual(margelet.SessionRepo.Dialog(100, 500), []string{"Test"}) {
		t.Fail()
	}
}

func TestAddDialogWithWrongIds(t *testing.T) {
	getMargelet()
	margelet.SessionRepo.Add(100, 500, "Test")

	if !reflect.DeepEqual(margelet.SessionRepo.Dialog(100, 501), []string{}) {
		t.Fail()
	}
}

func TestRemoveSession(t *testing.T) {
	getMargelet()
	margelet.SessionRepo.Add(100, 500, "Test")
	margelet.SessionRepo.Remove(100, 500)

	if !reflect.DeepEqual(margelet.SessionRepo.Dialog(100, 500), []string{}) {
		t.Fail()
	}
}
