# margelet
Telegram Bot Framework for Go based on [telegram-bot-api](https://github.com/Syfaro/telegram-bot-api)

It uses Redis for storing it's states, configs and so on

**Note: margelet in early beta now. Any advices and 
suggestions is required**

## Installation
`go get https://github.com/zhulik/margelet`

## Simple usage
```go
import (
    "github.com/zhulik/margelet"
)

func main() {
    margelet, err := margelet.NewMargelet("<your awesome bot name>", "<redis addr>", "<redis password>", 0, "your bot token", false)
    
    if err != nil {
        panic(err)
    }
    
    margelet.Run()
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