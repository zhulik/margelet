package margelet

import (
	"fmt"
	"github.com/Syfaro/telegram-bot-api"
)

// PanicResponder - test responder that panics
type PanicResponder struct {
	Margelet *Margelet
}

// Response sends default help message
func (responder PanicResponder) Response(bot MargeletAPI, message tgbotapi.Message) error {
	panic("TEST")
	return nil
}

// HelpMessage return help string for HelpResponder
func (responder PanicResponder) HelpMessage() string {
	return "Panic!"
}
