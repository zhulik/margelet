package margelet

import (
	"gopkg.in/redis.v3"
	"gopkg.in/telegram-bot-api.v2"
)

type policies []AuthorizationPolicy

func (p policies) Allow(message tgbotapi.Message) error {
	if len(p) == 0 {
		return nil
	}

	for _, policy := range p {
		if err := policy.Allow(message); err != nil {
			return err
		}
	}

	return nil
}

type authorizedCommandHandler struct {
	policies
	handler CommandHandler
}

type authorizedSessionHandler struct {
	policies
	handler SessionHandler
}

// Margelet - main struct in package, handles all interactions
type Margelet struct {
	bot TGBotAPI

	MessageHandlers []MessageHandler
	CommandHandlers map[string]authorizedCommandHandler
	SessionHandlers map[string]authorizedSessionHandler
	InlineHandler   InlineHandler

	running              bool
	Redis                *redis.Client
	ChatRepository       *chatRepository
	SessionRepository    SessionRepository
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

	margelet := Margelet{
		bot:                  bot,
		MessageHandlers:      []MessageHandler{},
		CommandHandlers:      map[string]authorizedCommandHandler{},
		SessionHandlers:      map[string]authorizedSessionHandler{},
		InlineHandler:        nil,
		running:              true,
		Redis:                redis,
		ChatRepository:       chatRepository,
		SessionRepository:    sessionRepository,
		ChatConfigRepository: chatConfigRepository,
	}

	margelet.AddCommandHandler("/help", HelpHandler{&margelet})

	return &margelet, nil
}

// AddMessageHandler - adds new MessageHandler to Margelet
func (margelet *Margelet) AddMessageHandler(handler MessageHandler) {
	margelet.MessageHandlers = append(margelet.MessageHandlers, handler)
}

// AddCommandHandler - adds new CommandHandler to Margelet
func (margelet *Margelet) AddCommandHandler(command string, handler CommandHandler, auth ...AuthorizationPolicy) {
	margelet.CommandHandlers[command] = authorizedCommandHandler{auth, handler}
}

// AddSessionHandler - adds new SessionHandler to Margelet
func (margelet *Margelet) AddSessionHandler(command string, handler SessionHandler, auth ...AuthorizationPolicy) {
	margelet.SessionHandlers[command] = authorizedSessionHandler{auth, handler}
}

// Send - send message to Telegram
func (margelet *Margelet) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	return margelet.bot.Send(c)
}

// AnswerInlineQuery  - send answer to InlineQuery
func (margelet *Margelet) AnswerInlineQuery(config tgbotapi.InlineConfig) (tgbotapi.APIResponse, error) {
	return margelet.bot.AnswerInlineQuery(config)
}

// QuickSend - quick send text message to chatID
func (margelet *Margelet) QuickSend(chatID int, message string) (tgbotapi.Message, error) {
	return margelet.bot.Send(tgbotapi.NewMessage(chatID, message))
}

// QuickReply - quick send text reply to message
func (margelet *Margelet) QuickReply(chatID, messageID int, message string) (tgbotapi.Message, error) {
	msg := tgbotapi.NewMessage(chatID, message)
	msg.ReplyToMessageID = messageID
	return margelet.bot.Send(msg)
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
func (margelet *Margelet) GetSessionRepository() SessionRepository {
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
		handleUpdate(margelet, <-updates)
	}
	return nil
}

// Stop - stops message processing loop
func (margelet *Margelet) Stop() {
	margelet.running = false
}

// HandleSession - handles any message as session message with handler
func (margelet *Margelet) HandleSession(message tgbotapi.Message, command string) {
	if authHandler, ok := margelet.SessionHandlers[command]; ok {
		handleSession(margelet, message, authHandler)
		return
	}
}
