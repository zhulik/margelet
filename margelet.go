package margelet

import (
	"github.com/zhulik/telegram-bot-api"
)

type Margelet struct {
	bot        *tgbotapi.BotAPI
	responders []Responder
}

func NewMargelet(token string) (*Margelet, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	bot.Debug = false

	var ucfg tgbotapi.UpdateConfig = tgbotapi.NewUpdate(0)
	ucfg.Timeout = 60
	err = bot.UpdatesChan(ucfg)

	if err != nil {
		return nil, err
	}

	return &Margelet{bot, []Responder{}}, nil
}

func (this *Margelet) AddResponder(responder Responder) {
	this.responders = append(this.responders, responder)
}

func (this *Margelet) Send(c tgbotapi.Chattable) {
	this.bot.Send(c)
}

func (this *Margelet) Run() {
	for {
		select {
		case update := <-this.bot.Updates:
			for _, responder := range this.responders {
				msg, err := responder.Response(update.Message)

				if err != nil {
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Error occured: "+err.Error())
				}

				this.Send(msg)
			}
		}
	}
}
