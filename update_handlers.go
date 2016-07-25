package margelet

import (
	"fmt"
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"runtime/debug"
	"strings"
)

func handleUpdate(margelet *Margelet, update tgbotapi.Update) {
	defer func() {
		if err := recover(); err != nil {
			if margelet.verbose {
				log.Println(string(debug.Stack()))
			}

			var panicMessage string

			if margelet.verbose {
				panicMessage = fmt.Sprintf("Panic occured: %v", err)
			} else {
				panicMessage = "Panic occured!"
			}

			margelet.QuickSend(update.Message.Chat.ID, panicMessage, nil)
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
		handler.HandleCallback(newCallbackQuery(margelet, query))
	}
}

func handleCommand(margelet *Margelet, message *tgbotapi.Message) {
	if authHandler, ok := margelet.CommandHandlers[strings.TrimSpace(message.Command())]; ok {
		if err := authHandler.Allow(message); err != nil {
			margelet.QuickSend(message.Chat.ID, "Authorization error: "+err.Error(), nil)
			return
		}
		err := authHandler.handler.HandleCommand(newMessage(margelet, message))

		if err != nil {
			margelet.QuickSend(message.Chat.ID, "Error occured: "+err.Error(), nil)
		}
		return
	}

	if authHandler, ok := margelet.SessionHandlers[strings.TrimSpace(message.Command())]; ok {
		margelet.SessionRepository.Create(message.Chat.ID, message.From.ID, strings.TrimSpace(message.Command()))
		handleSession(margelet, message, authHandler)
		return
	}

	if margelet.UnknownCommandHandler != nil {
		if err := margelet.UnknownCommandHandler.Allow(message); err != nil {
			margelet.QuickSend(message.Chat.ID, "Authorization error: "+err.Error(), nil)
			return
		}
		err := margelet.UnknownCommandHandler.handler.HandleCommand(newMessage(margelet, message))

		if err != nil {
			margelet.QuickSend(message.Chat.ID, "Error occured: "+err.Error(), nil)
		}
	}
}

func handleMessage(margelet *Margelet, msg *tgbotapi.Message) {
	for _, handler := range margelet.MessageHandlers {
		m := newMessage(margelet, msg)
		err := handler.HandleMessage(m)

		if err != nil {
			m.QuickSend("Error occured: "+err.Error(), nil)
		}
	}
}

func handleSession(margelet *Margelet, message *tgbotapi.Message, authHandler authorizedSessionHandler) {
	if err := authHandler.Allow(message); err != nil {
		margelet.QuickSend(message.Chat.ID, "Authorization error: "+err.Error(), nil)
		return
	}
	dialog := margelet.SessionRepository.Dialog(message.Chat.ID, message.From.ID)
	session := newMargetletSession(margelet, message, dialog)
	if strings.TrimSpace(message.Command()) == "cancel" {
		authHandler.handler.CancelSession(session)
		margelet.SessionRepository.Remove(message.Chat.ID, message.From.ID)
		return
	}

	err := authHandler.handler.HandleSession(session)

	if session.finished {
		margelet.SessionRepository.Remove(message.Chat.ID, message.From.ID)
		return
	}

	if err != nil {
		log.Printf("Margelet handling session error %s", err.Error())
		return
	}
	margelet.SessionRepository.Add(message.Chat.ID, message.From.ID, message)
}
