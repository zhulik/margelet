package margelet

import (
	"fmt"
	"gopkg.in/telegram-bot-api.v4"
	"strings"
)

func handleUpdate(margelet *Margelet, update tgbotapi.Update) {
	defer func() {
		if err := recover(); err != nil {

			var panicMessage string

			if margelet.verbose {
				panicMessage = fmt.Sprintf("Panic occured: %v", err)
			} else {
				panicMessage = "Panic occured!"
			}

			margelet.QuickSend(update.Message.Chat.ID, panicMessage)
		}
	}()

	switch {
	case update.Message != nil:
		message := update.Message
		margelet.ChatRepository.Add(message.Chat.ID)

		// If we have active session in this chat with this user, handle it first
		if command := margelet.SessionRepository.Command(message.Chat.ID, message.From.ID); len(command) > 0 {
			margelet.HandleSession(message, command)
		} else {
			if message.IsCommand() {
				handleCommand(margelet, message)
			} else {
				handleMessage(margelet, message)
			}
		}
	case update.InlineQuery != nil:
		handleInline(margelet, update.InlineQuery)
	case update.CallbackQuery != nil:
		handleCallback(margelet, update.CallbackQuery)
	}
}

func handleInline(margelet *Margelet, query *tgbotapi.InlineQuery) {
	handler := margelet.InlineHandler

	if handler != nil {
		handler.HandleInline(margelet, query)
	}
}

func handleCallback(margelet *Margelet, query *tgbotapi.CallbackQuery) {
	handler := margelet.CallbackHandler

	if handler != nil {
		handler.HandleCallback(margelet, query)
	}
}

func handleCommand(margelet *Margelet, message *tgbotapi.Message) {
	if authHandler, ok := margelet.CommandHandlers[strings.TrimSpace(message.Command())]; ok {
		if err := authHandler.Allow(message); err != nil {
			margelet.QuickSend(message.Chat.ID, "Authorization error: "+err.Error())
			return
		}
		err := authHandler.handler.HandleCommand(margelet, message)

		if err != nil {
			margelet.QuickSend(message.Chat.ID, "Error occured: "+err.Error())
		}
		return
	}

	if authHandler, ok := margelet.SessionHandlers[strings.TrimSpace(message.Command())]; ok {
		margelet.SessionRepository.Create(message.Chat.ID, message.From.ID, strings.TrimSpace(message.Command()))
		handleSession(margelet, message, authHandler)
		return
	}
}

func handleMessage(margelet *Margelet, message *tgbotapi.Message) {
	for _, handler := range margelet.MessageHandlers {
		err := handler.HandleMessage(margelet, message)

		if err != nil {
			margelet.QuickSend(message.Chat.ID, "Error occured: "+err.Error())
		}
	}
}

func handleSession(margelet *Margelet, message *tgbotapi.Message, authHandler authorizedSessionHandler) {
	if err := authHandler.Allow(message); err != nil {
		margelet.QuickSend(message.Chat.ID, "Authorization error: "+err.Error())
		return
	}
	if strings.TrimSpace(message.Command()) == "/cancel" {
		authHandler.handler.CancelSession(margelet, message, margelet.SessionRepository.Dialog(message.Chat.ID, message.From.ID))
		margelet.SessionRepository.Remove(message.Chat.ID, message.From.ID)
		return
	}

	finish, err := authHandler.handler.HandleSession(margelet, message, margelet.SessionRepository.Dialog(message.Chat.ID, message.From.ID))

	if finish {
		margelet.SessionRepository.Remove(message.Chat.ID, message.From.ID)
		return
	}

	if err == nil {
		margelet.SessionRepository.Add(message.Chat.ID, message.From.ID, message)
	}

}
