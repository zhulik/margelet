* [English](README.md)
 
[![Build Status](https://travis-ci.org/zhulik/margelet.svg?branch=master)](https://travis-ci.org/zhulik/margelet)

[![Build Status](https://travis-ci.org/zhulik/margelet.svg?branch=master)](https://travis-ci.org/zhulik/margelet)
<img src="https://img.shields.io/badge/last%20modified-today-brightgreen">
<img src="https://img.shields.io/badge/platform-linux--64%20%7C%20win--32%20%7C%20osx--64%20%7C%20win--64-lightgrey">
<img src="https://img.shields.io/badge/issues-1%20open-yellow">
<img src="https://img.shields.io/badge/license-MIT%20License-green">
# Margelet
Telegram Bot Framework для Go основан на Telegram-bot-API[telegram-bot-api](https://gopkg.in/telegram-bot-api.v4).

Он использует Redis для хранения своих состояний, конфигураций и так далее.

Любые низкоуровневые взаимодействия с Telegram Bot API  (загрузка файлов, клавиатур и т. д.) должны выполняться через telegram-bot-api (https://gopkg.in/telegram-bot-api.v4).

Margelet — это всего лишь тонкий слой, который позволяет быстро и легко решать основные задачи бота.

## Установка
`go get github.com/zhulik/margelet`

## Простое использование
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

	err = bot.Run()
	if err != nil {
		panic(err)
	}
}
```
Из коробки Margelet поддерживает только `/help` команду, отвечает примерно так

`/help - Show bot help`

## Концепция
Margelet основан на некоторых концепциях:

* Обработчики сообщений
* Обработчики команд
* Обработчики сеансов
* Конфигурации чата
* Встроенные обработчики
* Обработчики сообщений


### Обработчик сообщений — это структура, реализующая интерфейс Handler. Он получает все сообщения чата, зависящие от режима конфиденциальности бота [Privacy mode](https://core.telegram.org/bots#privacy-mode). Он не получает команды.

Простой пример:
```
package margelet_test

import (
	"../margelet"
)

// EchoHandler is simple handler example
type EchoHandler struct {
}

// Response send message back to author
func (handler EchoHandler) HandleMessage(m margelet.Message) error {
	_, err := m.QuickSend(m.Message().Text)
	return err
}

```

Этот обработчик будет повторять любое сообщение пользователя обратно в чат.

Помощники сообщений могут быть добавлены в margelet с помощью `AddMessageHandler` функции:
```
bot, err := margelet.NewMargelet("<your awesome bot name>", "<redis addr>", "<redis password>", 0, "your bot token", false)
bot.AddMessageHandler(EchoHandler{})
bot.Run()
```
### Обработчики команд
Обработчик команд — это структура, реализующая интерфейс CommandHandler. CommandHandler может быть подписан на любую нужную вам команду и будет получать все сообщения сообщения с этой командой, только если в этом чате нет активной сессии с этим пользователем

Простой пример:

```
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
```
Обработчики команд можно добавить в margelet с помощью `AddCommandHandler` функции:

``` go
bot, err := margelet.NewMargelet("<your awesome bot name>", "<redis addr>", "<redis password>", 0, "your bot token", false)
bot.AddCommandHandler("help", HelpHandler{bot})
bot.Run()
```

### Обработчики сеансов
Сессия здесь представляет собой интерактивный диалог с пользователем, как это делает [@BotFather](https://telegram.me/botfather). Пользователь запускает сеанс с помощью команды, а затем отвечает на вопросы бота, пока бот не соберет всю необходимую информацию. Его можно использовать, например, для настройки бота.

### API обработчиков сеансов все еще развивается**

Обработчик сеанса — это структура, реализующая интерфейс SessionHandler. Простой пример:

``` go 

package margelet_test

import (
	"fmt"
	"strconv"

	"../margelet"
	"gopkg.in/telegram-bot-api.v4"
)

// SumSession - simple example session, that can sum numbers
type SumSession struct {
}

// HandleResponse - Handlers user response
func (s SumSession) HandleSession(session margelet.Session) error {
	switch len(session.Responses()) {
	case 0:
		session.QuickReply("Hello, please, write one number per message, after some iterations write 'end'.")
	default:
		if session.Message().Text == "end" {
			var sum int
			for _, m := range session.Responses() {
				n, _ := strconv.Atoi(m.Text)
				sum += n
			}
			session.QuickReply(fmt.Sprintf("Your sum: %d", sum))
			session.Finish()
			return nil
		}

		_, err := strconv.Atoi(session.Message().Text)
		if err != nil {
			session.QuickReply("Sorry, not a number")
			return err
		}
	}

	return nil
}

// CancelResponse - Chance to clean up everything
func (s SumSession) CancelSession(session margelet.Session) {
	//Clean up all variables only used in the session

}

func (session SumSession) response(bot margelet.MargeletAPI, message *tgbotapi.Message, msg tgbotapi.MessageConfig) {
	msg.ChatID = message.Chat.ID
	msg.ReplyToMessageID = message.MessageID
	msg.ReplyMarkup = tgbotapi.ForceReply{ForceReply: true, Selective: true}
	bot.Send(msg)
}

// HelpMessage return help string for SumSession
func (session SumSession) HelpMessage() string {
	return "Sum your numbers and print result"
}
```

Обработчики сеансов могут быть добавлены в margelet с помощью `AddSessionHandler` функции:

``` go 
bot, err := margelet.NewMargelet("<your awesome bot name>", "<redis addr>", "<redis password>", 0, "your bot token", false)
bot.AddSessionHandler("help", SumSession{})
bot.Run()
```

При каждом ответе пользователя он получает все предыдущие ответы пользователя, поэтому вы можете восстановить состояние сеанса. Возвращаемые значения HandleResponse важны:

* first(bool) означает, что margelet должен завершить сеанс, поэтому верните true, если вы получили всю необходимую информацию от пользователя, иначе false
* second(err) означает, что бот не может обработать сообщение пользователя. Это сообщение не будет добавлено в историю диалогов сеанса. Возвращает любую ошибку, если вы можете обработать сообщение пользователя, и возвращает nil, если сообщение принято.

### Встроенные обработчики
Встроенный обработчик — это структура, реализующая интерфейс InlineHandler. InlineHandler может быть подписан на любые встроенные запросы.

Простой пример:

``` go 
package margelet_test

import (
	"github.com/zhulik/margelet"
	"gopkg.in/telegram-bot-api.v4"
)

type InlineImage struct {
}

func (handler InlineImage) HandleInline(bot margelet.MargeletAPI, query *tgbotapi.InlineQuery) error {
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

Встроенный обработчик может быть добавлен в margelet путем `InlineHandler` назначения:

``` go 

bot, err := margelet.NewMargelet("<your awesome bot name>", "<redis addr>", "<redis password>", 0, "your bot token", false)
m.InlineHandler = &InlineImage{}
bot.Run()
```

### Обработчики обратного вызова

Обработчик обратного вызова — это структура, реализующая интерфейс CallbackHandler. CallbackHandler может быть подписан на любые запросы обратного вызова.

Простой пример:

``` go 
package margelet_test

import (
	"../margelet"
	"gopkg.in/telegram-bot-api.v4"
)

type CallbackMessage struct {
}

func (handler CallbackMessage) HandleCallback(query margelet.CallbackQuery) error {
	config := tgbotapi.CallbackConfig{
		CallbackQueryID: query.Query().ID,
		Text:            "Done!",
		ShowAlert:       false,
	}

	query.Bot().AnswerCallbackQuery(config)
	return nil
}
```

Обработчик обратного вызова можно добавить в margelet по `CallbackHandler` назначению:

``` go 
bot, err := margelet.NewMargelet("<your awesome bot name>", "<redis addr>", "<redis password>", 0, "your bot token", false)
m.CallbackHandler = &CallbackMessage{}
bot.Run()
``` 

### Конфигурации чата
Боты могут хранить любую строку конфигурации (вы можете использовать сериализованный JSON) для любого чата. Его можно использовать для хранения пользовательских конфигураций и другой информации, связанной с пользователем. Простой пример:

``` go 
bot, err := margelet.NewMargelet("<your awesome bot name>", "<redis addr>", "<redis password>", 0, "your bot token", false)
...
bot.GetConfigRepository().Set(<chatID>, "<info>")
...
info := bot.GetConfigRepository().Get(<chatID>)
```
OR
``` go
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

Доступ к репозиторию конфигурации чата можно получить из обработчиков сеансов.

## Пример проекта
Простой и чистый пример проекта можно найти здесь -> [here](https://github.com/zhulik/cat_bot) . Он обеспечивает обработку команд и настройку сеанса.
