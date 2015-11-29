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
Margelet uses some base concepts:
* Message responders
* Command responders
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

