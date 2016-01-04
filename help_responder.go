package margelet

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"strings"
)

// HelpResponder Default responder for /help command. Margelet will add this automatically
type HelpResponder struct {
	Margelet *Margelet
}

// Response sends default help message
func (responder HelpResponder) Response(bot MargeletAPI, message tgbotapi.Message) error {
	lines := []string{}
	for command, responder := range responder.Margelet.CommandResponders {
		lines = append(lines, fmt.Sprintf("%s - %s", command, responder.HelpMessage()))
	}

	for command, responder := range responder.Margelet.SessionHandlers {
		lines = append(lines, fmt.Sprintf("%s - %s", command, responder.HelpMessage()))
	}

	_, err := bot.Send(tgbotapi.NewMessage(message.Chat.ID, strings.Join(lines, "\n")))
	return err
}

// HelpMessage return help string for HelpResponder
func (responder HelpResponder) HelpMessage() string {
	return "Show bot help"
}
