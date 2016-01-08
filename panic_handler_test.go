package margelet_test

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/zhulik/margelet"
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
