package margelet

import (
	"github.com/Syfaro/telegram-bot-api"
)

type Margelet struct {
	bot               TGBotAPI
	MessageResponders []Responder
	CommandResponders map[string]CommandHandler
	SessionHandlers   map[string]SessionHandler
	running           bool
}

func NewMargelet(token string, verbose bool) (*Margelet, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	bot.Debug = verbose

	return &Margelet{bot, []Responder{}, map[string]CommandHandler{}, map[string]SessionHandler{}, true}, nil
}

func NewMargeletFromBot(bot TGBotAPI) (*Margelet, error) {
	return &Margelet{bot, []Responder{}, map[string]CommandHandler{}, map[string]SessionHandler{}, true}, nil
}

func (this *Margelet) AddMessageResponder(responder Responder) {
	this.MessageResponders = append(this.MessageResponders, responder)
}

func (this *Margelet) AddCommandHandler(command string, responder CommandHandler) {
	this.CommandResponders[command] = responder
}

func (this *Margelet) AddSessionHandler(command string, responder SessionHandler) {
	this.SessionHandlers[command] = responder
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
			ChatRepo.Add(message.Chat.ID)
//			If we have active session with in this chat with this user, handle it first
			if len(SessionRepo.Find(message.Chat.ID, message.From.ID)) > 0 {
			} else if message.IsCommand() {
				this.handleCommand(message)
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

func (this *Margelet) handleCommand(message tgbotapi.Message) {
	if responder, ok := this.CommandResponders[message.Command()]; ok {
		this.handleMessage(message, []Responder{responder})
		return
	}

	if handler, ok := this.SessionHandlers[message.Command()]; ok {
		this.handleSession(message, handler)
		return
	}
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

func (this *Margelet) handleSession(message tgbotapi.Message, handler SessionHandler) {

	if err := handler.HandleResponse(this, message, SessionRepo.Find(message.Chat.ID, message.From.ID)); err == nil {
		//		SessionRepo.Add(message.Chat.ID, message.From.ID, message.Text, )
		//		SessionRepo.Remove(message.Chat.ID, message.From.ID)
	}
}
