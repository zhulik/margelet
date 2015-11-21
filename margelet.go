package margelet

import (
	"github.com/zhulik/telegram-bot-api"
)

type Margelet struct {
	bot        *tgbotapi.BotAPI
	messageResponders []Responder
	commandResponders []Responder
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

	return &Margelet{bot, []Responder{}, []Responder{}}, nil
}

func (this *Margelet) AddMessageResponder(responder Responder) {
	this.messageResponders = append(this.messageResponders, responder)
}

func (this *Margelet) AddCommandResponder(responder Responder) {
	this.commandResponders = append(this.commandResponders, responder)
}

func (this *Margelet) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	return this.bot.Send(c)
}

func (this *Margelet) GetFileDirectUrl(fileID string) (string, error) {
	return this.bot.GetFileDirectUrl(fileID)
}

func (this *Margelet) IsMessageToMe(message tgbotapi.Message) bool {
	return this.bot.IsMessageToMe(message)
}

func (this *Margelet) Run() {
	for {
		select {
		case update := <-this.bot.Updates:
			message := update.Message
			if message.IsCommand() {
				this.handleMessage(message, this.commandResponders)
			} else {
				this.handleMessage(message, this.messageResponders)
			}
		}
	}
}

func (this *Margelet) handleMessage(message tgbotapi.Message, responders []Responder) {
	for _, responder := range responders {
		msg, err := responder.Response(this, message)

		if err != nil {
			msg = tgbotapi.NewMessage(message.Chat.ID, "Error occured: " + err.Error())
		}

		_, err = this.Send(msg)
		if err != nil {
			msg = tgbotapi.NewMessage(message.Chat.ID, "Error occured: " + err.Error())
			this.Send(msg)
		}
	}
}