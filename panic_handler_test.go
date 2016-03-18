package margelet_test

import (
	"../margelet"
	"gopkg.in/telegram-bot-api.v2"
)

// PanicHandler - test handler that panics
type PanicHandler struct {
	Margelet *margelet.Margelet
}

// Handle sends default help message
func (handler PanicHandler) HandleMessage(bot margelet.MargeletAPI, message tgbotapi.Message) error {
	panic("TEST")
	return nil
}

// HelpMessage return help string for PanicHandler
func (handler PanicHandler) HelpMessage() string {
	return "Panic!"
}
