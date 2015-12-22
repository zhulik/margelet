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
	SessionRepository    *SessionRepository
	ChatConfigRepository *ChatConfigRepository
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

// QuickSend - quick send text message to chatID
func (margelet *Margelet) QuickSend(chatID int, message string) (tgbotapi.Message, error) {
	return margelet.bot.Send(tgbotapi.NewMessage(chatID, message))
}

// GetFileDirectURL - converts fileID to direct URL
func (margelet *Margelet) GetFileDirectURL(fileID string) (string, error) {
	return margelet.bot.GetFileDirectURL(fileID)
}

// IsMessageToMe - return true if message sent to this bot
func (margelet *Margelet) IsMessageToMe(message tgbotapi.Message) bool {
	return margelet.bot.IsMessageToMe(message)
}

// GetConfigRepository - returns chat config repository
func (margelet *Margelet) GetConfigRepository() *ChatConfigRepository {
	return margelet.ChatConfigRepository
}

// GetSessionRepository - returns session repository
func (margelet *Margelet) GetSessionRepository() *SessionRepository {
	return margelet.SessionRepository
}

// GetRedis - returns margelet's redis client
func (margelet *Margelet) GetRedis() *redis.Client {
	return margelet.Redis
}

// Run - starts message processing loop
func (margelet *Margelet) Run() error {
	updates, err := margelet.bot.GetUpdatesChan(tgbotapi.UpdateConfig{Timeout: 60})

	if err != nil {
		return err
	}

	for margelet.running {
		margelet.handleUpdate(<-updates)
	}
	return nil
}

// Stop - stops message processing loop
func (margelet *Margelet) Stop() {
	margelet.running = false
}

// HandleSession - handles any message as session message with handler
func (margelet *Margelet) HandleSession(message tgbotapi.Message, handler SessionHandler) {
	finish, err := handler.HandleResponse(margelet, message, margelet.SessionRepository.Dialog(message.Chat.ID, message.From.ID))
	if err == nil {
		margelet.SessionRepository.Add(message.Chat.ID, message.From.ID, message.Text)
	}

	if finish {
		margelet.SessionRepository.Remove(message.Chat.ID, message.From.ID)
	}
}


func (margelet *Margelet) handleUpdate(update tgbotapi.Update) {
	defer func() {
		if err := recover(); err != nil {
			margelet.QuickSend(update.Message.Chat.ID, "Panic occured!")
		}
	}()

	message := update.Message
	margelet.ChatRepository.Add(message.Chat.ID)

	// If we have active session in this chat with this user, handle it first
	if command := margelet.SessionRepository.Command(message.Chat.ID, message.From.ID); len(command) > 0 {
		// TODO: /cancel command should cancel any active session!
		if handler, ok := margelet.SessionHandlers[command]; ok {
			margelet.HandleSession(message, handler)
		}
	} else if message.IsCommand() {
		margelet.handleCommand(message)
	} else {
		margelet.handleMessage(message, margelet.MessageResponders)
	}
}

func (margelet *Margelet) handleCommand(message tgbotapi.Message) {
	if responder, ok := margelet.CommandResponders[message.Command()]; ok {
		margelet.handleMessage(message, []Responder{responder})
		return
	}

	if handler, ok := margelet.SessionHandlers[message.Command()]; ok {
		margelet.SessionRepository.Create(message.Chat.ID, message.From.ID, message.Command())
		margelet.HandleSession(message, handler)
		return
	}
}

func (margelet *Margelet) handleMessage(message tgbotapi.Message, responders []Responder) {
	for _, responder := range responders {
		err := responder.Response(margelet, message)

		if err != nil {
			margelet.QuickSend(message.Chat.ID, "Error occured: "+err.Error())
		}
	}
}
