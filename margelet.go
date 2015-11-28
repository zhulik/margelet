package margelet

import (
	"github.com/Syfaro/telegram-bot-api"
	"gopkg.in/redis.v3"
)

// Margelet - main struct in package, handles all interactions
type Margelet struct {
	bot                  TGBotAPI
	MessageResponders    []Responder
	CommandResponders    map[string]CommandHandler
	SessionHandlers      map[string]SessionHandler
	running              bool
	Redis                *redis.Client
	ChatRepository       *chatRepository
	SessionRepository    *sessionRepository
	ChatConfigRepository *chatConfigRepository
}

// NewMargelet creates new Margelet instance
func NewMargelet(botName string, redisAddr string, redisPassword string, redisDB int64, token string, verbose bool) (*Margelet, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	bot.Debug = verbose

	return NewMargeletFromBot(botName, redisAddr, redisPassword, redisDB, bot)
}

// NewMargeletFromBot creates new Margelet instance from existing TGBotAPI(tgbotapi.BotAPI)
func NewMargeletFromBot(botName string, redisAddr string, redisPassword string, redisDB int64, bot TGBotAPI) (*Margelet, error) {
	redis := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	})

	if _, err := redis.Ping().Result(); err != nil {
		return &Margelet{}, err
	}

	chatRepository := newChatRepository(botName, redis)
	sessionRepository := newSessionRepository(botName, redis)
	chatConfigRepository := newChatConfigRepository(botName, redis)

	margelet := Margelet{bot, []Responder{}, map[string]CommandHandler{}, map[string]SessionHandler{}, true, redis, chatRepository, sessionRepository, chatConfigRepository}

	margelet.AddCommandHandler("/help", HelpResponder{&margelet})

	return &margelet, nil
}

// AddMessageResponder - adds new MessageResponder to Margelet
func (margelet *Margelet) AddMessageResponder(responder Responder) {
	margelet.MessageResponders = append(margelet.MessageResponders, responder)
}

// AddCommandHandler - adds new CommandHandler to Margelet
func (margelet *Margelet) AddCommandHandler(command string, responder CommandHandler) {
	margelet.CommandResponders[command] = responder
}

// AddSessionHandler - adds new SessionHandler to Margelet
func (margelet *Margelet) AddSessionHandler(command string, responder SessionHandler) {
	margelet.SessionHandlers[command] = responder
}

// Send - send message to Telegram
func (margelet *Margelet) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	return margelet.bot.Send(c)
}

// GetFileDirectURL - converts fileID to direct URL
func (margelet *Margelet) GetFileDirectURL(fileID string) (string, error) {
	return margelet.bot.GetFileDirectURL(fileID)
}

// IsMessageToMe - return true if message sent to this bot
func (margelet *Margelet) IsMessageToMe(message tgbotapi.Message) bool {
	return margelet.bot.IsMessageToMe(message)
}

// Run - starts message processing loop
func (margelet *Margelet) Run() error {
	updates, err := margelet.bot.GetUpdatesChan(tgbotapi.UpdateConfig{Timeout: 60})

	if err != nil {
		return err
	}

	for margelet.running {
		select {
		case update := <-updates:
			message := update.Message
			margelet.ChatRepository.Add(message.Chat.ID)

			// If we have active session in this chat with this user, handle it first
			if command := margelet.SessionRepository.Command(message.Chat.ID, message.From.ID); len(command) > 0 {
				if handler, ok := margelet.SessionHandlers[command]; ok {
					margelet.handleSession(message, handler)
				}
			} else if message.IsCommand() {
				margelet.handleCommand(message)
			} else {
				margelet.handleMessage(message, margelet.MessageResponders)
			}
		}
	}
	return nil
}

// Stop - stops message processing loop
func (margelet *Margelet) Stop() {
	margelet.running = false
}

func (margelet *Margelet) handleCommand(message tgbotapi.Message) {

	if responder, ok := margelet.CommandResponders[message.Command()]; ok {
		margelet.handleMessage(message, []Responder{responder})
		return
	}

	if handler, ok := margelet.SessionHandlers[message.Command()]; ok {
		margelet.SessionRepository.Create(message.Chat.ID, message.From.ID, message.Command())
		margelet.handleSession(message, handler)
		return
	}
}

func (margelet *Margelet) handleMessage(message tgbotapi.Message, responders []Responder) {
	for _, responder := range responders {
		err := responder.Response(margelet, message)

		if err != nil {
			msg := tgbotapi.NewMessage(message.Chat.ID, "Error occured: "+err.Error())
			margelet.Send(msg)
		}
	}
}

func (margelet *Margelet) handleSession(message tgbotapi.Message, handler SessionHandler) {

	finish, err := handler.HandleResponse(margelet, message, margelet.SessionRepository.Dialog(message.Chat.ID, message.From.ID))
	if err == nil {
		margelet.SessionRepository.Add(message.Chat.ID, message.From.ID, message.Text)
	}

	if finish {
		margelet.SessionRepository.Remove(message.Chat.ID, message.From.ID)
	}
}
