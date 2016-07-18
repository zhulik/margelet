package margelet

import (
	"fmt"
	"strings"
)

// HelpHandler Default handler for /help command. Margelet will add this automatically
type HelpHandler struct {
	Margelet *Margelet
}

// HandleCommand sends default help message
func (handler HelpHandler) HandleCommand(message Message) error {
	lines := []string{}
	for command, h := range handler.Margelet.CommandHandlers {
		lines = append(lines, fmt.Sprintf("/%s - %s", command, h.handler.HelpMessage()))
	}

	for command, h := range handler.Margelet.SessionHandlers {
		lines = append(lines, fmt.Sprintf("/%s - %s", command, h.handler.HelpMessage()))
	}

	_, err := message.QuickSend(strings.Join(lines, "\n"))
	return err
}

// HelpMessage return help string for HelpHandler
func (handler HelpHandler) HelpMessage() string {
	return "Show bot help"
}
