[![Build Status](https://travis-ci.org/zhulik/margelet.svg?branch=master)](https://travis-ci.org/zhulik/margelet)
# Margelet
Telegram Bot Framework for Go is based on [telegram-bot-api](https://gopkg.in/telegram-bot-api.v2)

It uses Redis to store it's states, configs and so on.

Any low-level interactions with Telegram Bot API(downloading files, keyboards and so on) should be performed through
[telegram-bot-api](https://"gopkg.in/telegram-bot-api.v2").

Margelet is just a thin layer, that allows you to solve
basic bot tasks quickly and easy.

**Note: margelet is in early beta now. Any advices and suggestions are welcome**

## Installation
`go get https://github.com/zhulik/margelet`

## Simple usage
```go
package main

import (
    "github.com/zhulik/margelet"
)

func main() {
    bot, err := margelet.NewMargelet("<your awesome bot name>", "<redis addr>", "<redis password>", 0, "your bot token", false)

    if err != nil {
        panic(err)
    }

    bot.Run()
}
```

Out of the box, margelet supports only `/help` command, it responds something like this

`/help - Show bot help`

## Concept
Margelet is based on some concepts:
* Message handlers
* Command handlers
* Session handlers
* Chat configs
* Inline handlers

### Message handlers
Message handler is a struct that implements Handler interface. It receives all chat messages dependant on bot's
[Privacy mode](https://core.telegram.org/bots#privacy-mode). It doesn't receive commands.

Simple example:
```go
// EchoHandler is simple handler example
type EchoHandler struct {
}

// Response send message back to author
func (handler EchoHandler) HandleMessage(bot margelet.MargeletAPI, message tgbotapi.Message) error {
	_, err := bot.Send(tgbotapi.NewMessage(message.Chat.ID, message.Text))
	return err
}

```

This handler will repeat any user's message back to chat.

Message helpers can be added to margelet with `AddMessageHandler` function:
```go
bot, err := margelet.NewMargelet("<your awesome bot name>", "<redis addr>", "<redis password>", 0, "your bot token", false)
bot.AddMessageHandler(EchoHandler{})
bot.Run()
```

### Command handlers
Command handler is struct that implements CommandHandler interface. CommandHandler can be subscribed on any command you need
and will receive all message messages with this command, only if there is no active session with this user in this chat

Simple example:
```go
// HelpHandler Default handler for /help command. Margelet will add this automatically
type HelpHandler struct {
	Margelet *Margelet
}

// Handle sends default help message
func (handler HelpHandler) HandleCommand(bot MargeletAPI, message tgbotapi.Message) error {
	lines := []string{}
	for command, h := range handler.Margelet.CommandHandlers {
		lines = append(lines, fmt.Sprintf("%s - %s", command, h.handler.HelpMessage()))
	}

	for command, h := range handler.Margelet.SessionHandlers {
		lines = append(lines, fmt.Sprintf("%s - %s", command, h.handler.HelpMessage()))
	}

	_, err := bot.Send(tgbotapi.NewMessage(message.Chat.ID, strings.Join(lines, "\n")))
	return err
}

// HelpMessage return help string for HelpHandler
func (handler HelpHandler) HelpMessage() string {
	return "Show bot help"
}
```

Command handlers can be added to margelet with `AddCommandHandler` function:
```go
bot, err := margelet.NewMargelet("<your awesome bot name>", "<redis addr>", "<redis password>", 0, "your bot token", false)
bot.AddCommandHandler("help", HelpHandler{bot})
bot.Run()
```

### Session handlers
Session here is an interactive dialog with user, like [@BotFather](https://telegram.me/botfather) does. User runs session
with a command and then response to bot's questions until bot collects all needed information. It can be used for bot
configuration, for example.

**Session handlers API is still developing**

Session handler is struct that implements SessionHandler interface. Simple example:
```go
// SumSession - simple example session, that can sum numbers
type SumSession struct {
}

// HandleResponse - Handlers user response
func (session SumSession) HandleResponse(bot MargeletAPI, message tgbotapi.Message, responses []tgbotapi.Message) (bool, error) {
	var msg tgbotapi.MessageConfig
	switch len(responses) {
	case 0:
		msg = tgbotapi.MessageConfig{Text: "Hello, please, write one number per message, after some iterations write 'end'."}
		msg.ReplyMarkup = tgbotapi.ForceReply{true, true}
	default:
		if message.Text == "end" {
			var sum int
			for _, m := range responses {
				n, _ := strconv.Atoi(m.Text)
				sum += n
			}
			msg = tgbotapi.MessageConfig{Text: fmt.Sprintf("Your sum: %d", sum)}
			session.response(bot, message, msg)
			msg.ReplyMarkup = tgbotapi.ForceReply{false, true}
			return true, nil
		}

		_, err := strconv.Atoi(message.Text)
		if err != nil {
			msg = tgbotapi.MessageConfig{Text: "Sorry, not a number"}
			session.response(bot, message, msg)
			msg.ReplyMarkup = tgbotapi.ForceReply{true, true}
			return false, err
		}
	}

	session.response(bot, message, msg)
	return false, nil
}

// CancelResponse - Chance to clean up everything
func (session SumSession) CancelSession(bot MargeletAPI, message tgbotapi.Message, responses []tgbotapi.Message){
  //Clean up all variables only used in the session

}

func (session SumSession) response(bot MargeletAPI, message tgbotapi.Message, msg tgbotapi.MessageConfig) {
	msg.ChatID = message.Chat.ID
	msg.ReplyToMessageID = message.MessageID
	bot.Send(msg)
}

// HelpMessage return help string for SumSession
func (session SumSession) HelpMessage() string {
	return "Sum your numbers and print result"
}
```
Command handlers can be added to margelet with `AddSessionHandler` function:
```go
bot, err := margelet.NewMargelet("<your awesome bot name>", "<redis addr>", "<redis password>", 0, "your bot token", false)
bot.AddSessionHandler("help", SumSession{})
bot.Run()
```

On each user response it receives all previous user responses, so you can restore session state. HandleResponse return values
it important:
* first(bool), means that margelet should finish session, so return true if you receive all needed info from user, false otherwise
* second(err), means that bot cannot handle user's message. This message will not be added to session dialog history.
Return any error if you can handle user's message and return nil if message is accepted.

### Inline handlers
Inline handler is struct that implements InlineHandler interface. InlineHandler can be subscribed on any inline queries.

Simple example:
```go

package margelet_test

import (
	"github.com/zhulik/margelet"
	"gopkg.in/telegram-bot-api.v2"
)

type InlineImage struct {
}

func (handler InlineImage) HandleInline(bot margelet.MargeletAPI, query tgbotapi.InlineQuery) error {
	testPhotoQuery := tgbotapi.NewInlineQueryResultPhoto(query.ID, "https://telegram.org/img/t_logo.png")
	testPhotoQuery.ThumbURL = "https://telegram.org/img/t_logo.png"

	config := tgbotapi.InlineConfig{
		InlineQueryID: query.ID,
		CacheTime:     2,
		IsPersonal:    false,
		Results:       []interface{}{testPhotoQuery},
		NextOffset:    "",
	}

	bot.AnswerInlineQuery(config)
	return nil
}


```

Inline handler can be added to margelet by `InlineHandler` assignment:
```go
bot, err := margelet.NewMargelet("<your awesome bot name>", "<redis addr>", "<redis password>", 0, "your bot token", false)
m.InlineHandler = &InlineImage{}
bot.Run()
```

### Chat configs
Bots can store any config string(you can use serialized JSON) for any chat. It can be used for storing user's
configurations and other user-related information. Simple example:
```go
bot, err := margelet.NewMargelet("<your awesome bot name>", "<redis addr>", "<redis password>", 0, "your bot token", false)
...
bot.GetConfigRepository().Set(<chatID>, "<info>")
...
info := bot.GetConfigRepository().Get(<chatID>)

OR

type userInfo struct{
  FavColor string // First character has to be Capital otherwise it wont be saved
}
...
user := userInfo{FavColor: "Green"}
bot.GetConfigRepository().SetWithStruct(<chatID>, user)
...
var user userInfo
bot.GetConfigRepository().GetWithStruct(<chatID>, &user)

```
Chat config repository can be accessed from session handlers.

## Example project
Simple and clean example project can be found [here](https://github.com/zhulik/cat_bot). It provides command handling
and session configuration.
