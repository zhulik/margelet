package margelet_test

import (
	"github.com/Syfaro/telegram-bot-api"
	"github.com/zhulik/margelet"
)

// PanicResponder - test responder that panics
type PanicResponder struct {
	Margelet *margelet.Margelet
}

// Response sends default help message
func (responder PanicResponder) Response(bot margelet.MargeletAPI, message tgbotapi.Message) error {
	panic("TEST")
	return nil
}

// HelpMessage return help string for HelpResponder
func (responder PanicResponder) HelpMessage() string {
	return "Panic!"
}
