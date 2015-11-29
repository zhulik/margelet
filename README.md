# margelet
Telegram Bot Framework for Go based on [telegram-bot-api](https://github.com/Syfaro/telegram-bot-api)

It uses Redis for storing it's states, configs and so on. 

Any low-level interactions with Telegram Bot API(downloading files, keyboards and so on) should be performed through 
[telegram-bot-api](https://github.com/Syfaro/telegram-bot-api). 

Margelet it just thin layer, that allows you to solve
base bot tasks quickly and easy.

**Note: margelet in early beta now. Any advices and suggestions is required**

## Installation
`go get https://github.com/zhulik/margelet`

## Simple usage
```go
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

Out of box, margelet support only `/help` command, it respond some like this

`/help - Show bot help`

## Concept
Margelet based on some concepts:
* Message responders
* Command handlers
* Session handlers
* Chat configs

### Message responders
Message responder is struct that implements Responder interface. It receives all chat messages dependant on bot's
[Privacy mode](https://core.telegram.org/bots#privacy-mode). It don't receive commands.

Simple example:
```go
// EchoResponder is simple responder example
type EchoResponder struct {
}

// Response send message back to author
func (responder EchoResponder) Response(bot MargeletAPI, message tgbotapi.Message) error {
	_, err := bot.Send(tgbotapi.NewMessage(message.Chat.ID, message.Text))
	return err
}
```

This responder will repeat any user's message back to chat.

Message responders can be added to margelet with `AddMessageResponder` function:
```go
bot, err := margelet.NewMargelet("<your awesome bot name>", "<redis addr>", "<redis password>", 0, "your bot token", false)
bot.AddMessageResponder(EchoResponder{})
bot.Run()
```

### Command handlers
Command handler is struct that implements CommandHandler interface. CommandHandler can be subscribed on any command you need
and will receive all message messages with this command, only if there is no active session with this user in this chat

Simple example:
```go
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
```

Command handlers can be added to margelet with `AddSessionHandler` function:
```go
bot, err := margelet.NewMargelet("<your awesome bot name>", "<redis addr>", "<redis password>", 0, "your bot token", false)
bot.AddSessionHandler("/help", HelpResponder{bot})
bot.Run()
```

### Session handlers
Session here is interactive dialog with user, like [@BotFather](https://telegram.me/botfather) do. User runs session
with a command and then response to bot's questions until bot collects all needed information. It can be used for bot
configuration, for example.

**Session handlers API is still developing**

Session handler is struct that implements SessionHandler interface. Simple example:
```go
// SumSession - simple example session, that can sum numbers
type SumSession struct {
}

// HandleResponse - Handlers user response
func (session SumSession) HandleResponse(bot MargeletAPI, message tgbotapi.Message, responses []string) (bool, error) {
	var msg tgbotapi.MessageConfig
	switch len(responses) {
	case 0:
		msg = tgbotapi.MessageConfig{Text: "Hello, please, write one number per message, after some iterations write 'end'."}
		msg.ReplyMarkup = tgbotapi.ForceReply{true, true}
	default:
		if message.Text == "end" {
			var sum int
			for _, a := range responses {
				n, _ := strconv.Atoi(a)
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

