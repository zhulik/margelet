package margelet_test

import (
	"../margelet"
)

// PanicHandler - test handler that panics
type PanicHandler struct {
	Margelet *margelet.Margelet
}

// Handle sends default help message
func (handler PanicHandler) HandleMessage(m margelet.Message) error {
	panic("TEST")
	return nil
}

// HelpMessage return help string for PanicHandler
func (handler PanicHandler) HelpMessage() string {
	return "Panic!"
}
