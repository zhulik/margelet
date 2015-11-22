package margelet

import (
	"github.com/zhulik/telegram-bot-api"
)

type Margelet struct {
	bot               TGBotAPI
	MessageResponders []Responder
	CommandResponders map[string]CommandHandler
	running           bool
}

func NewMargelet(token string, verbose bool) (*Margelet, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	bot.Debug = verbose

	return &Margelet{bot, []Responder{}, map[string]CommandHandler{}, true}, nil
}

func NewMargeletFromBot(bot TGBotAPI) (*Margelet, error) {
	return &Margelet{bot, []Responder{}, map[string]CommandHandler{}, true}, nil
}

func (this *Margelet) AddMessageResponder(responder Responder) {
	this.MessageResponders = append(this.MessageResponders, responder)
}

func (this *Margelet) AddCommandResponder(command string, responder Responder) {
	this.CommandResponders[command] = responder
}

func (this *Margelet) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	return this.bot.Send(c)
}

func (this *Margelet) GetFileDirectURL(fileID string) (string, error) {
	return this.bot.GetFileDirectURL(fileID)
}

func (this *Margelet) IsMessageToMe(message tgbotapi.Message) bool {
	return this.bot.IsMessageToMe(message)
}

func (this *Margelet) Run() error {
	updates, err := this.bot.GetUpdatesChan(tgbotapi.UpdateConfig{Timeout: 60})

	if err != nil {
		return err
	}

	for this.running {
		select {
		case update := <-updates:
			message := update.Message
			if message.IsCommand() {
				if responder, ok := this.CommandResponders[message.Command()]; ok {
					this.handleMessage(message, []Responder{responder})
				}
			} else {
				this.handleMessage(message, this.MessageResponders)
			}
		}
	}
	return nil
}

func (this *Margelet) Stop() {
	this.running = false
}

func (this *Margelet) handleMessage(message tgbotapi.Message, responders []Responder) {
	for _, responder := range responders {
		err := responder.Response(this, message)

		if err != nil {
			msg := tgbotapi.NewMessage(message.Chat.ID, "Error occured: "+err.Error())
			this.Send(msg)
		}
	}
}
