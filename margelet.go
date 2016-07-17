package margelet

import (
	"fmt"
	"gopkg.in/redis.v3"
	"gopkg.in/telegram-bot-api.v4"
	"net/http"
	"path/filepath"
)

type policies []AuthorizationPolicy

func (p policies) Allow(message *tgbotapi.Message) error {
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
	CallbackHandler CallbackHandler

	running              bool
	verbose              bool
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

	return NewMargeletFromBot(botName, redisAddr, redisPassword, redisDB, bot, verbose)
}

// NewMargeletFromBot creates new Margelet instance from existing TGBotAPI(tgbotapi.BotAPI)
func NewMargeletFromBot(botName string, redisAddr string, redisPassword string, redisDB int64, bot TGBotAPI, verbose bool) (*Margelet, error) {
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
		CallbackHandler:      nil,
		running:              true,
		verbose:              verbose,
		Redis:                redis,
		ChatRepository:       chatRepository,
		SessionRepository:    sessionRepository,
		ChatConfigRepository: chatConfigRepository,
	}

	margelet.AddCommandHandler("help", HelpHandler{&margelet})

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

// AnswerCallbackQuery  - send answer to CallbackQuery
func (margelet *Margelet) AnswerCallbackQuery(config tgbotapi.CallbackConfig) (tgbotapi.APIResponse, error) {
	return margelet.bot.AnswerCallbackQuery(config)
}

// QuickSend - quick send text message to chatID
func (margelet *Margelet) QuickSend(chatID int64, message string) (tgbotapi.Message, error) {
	return margelet.bot.Send(tgbotapi.NewMessage(chatID, message))
}

// QuickReply - quick send text reply to message
func (margelet *Margelet) QuickReply(chatID int64, messageID int, message string) (tgbotapi.Message, error) {
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
func (margelet *Margelet) HandleSession(message *tgbotapi.Message, command string) {
	if authHandler, ok := margelet.SessionHandlers[command]; ok {
		handleSession(margelet, message, authHandler)
		return
	}
}

// StartSession - start new session with given command, adds message to dialog
func (margelet *Margelet) StartSession(message *tgbotapi.Message, command string) {
	if authHandler, ok := margelet.SessionHandlers[command]; ok {
		margelet.GetSessionRepository().Create(message.Chat.ID, message.From.ID, command)
		margelet.GetSessionRepository().Add(message.Chat.ID, message.From.ID, message)
		handleSession(margelet, message, authHandler)
		return
	}
}

// SendImageByURL - sends given by url image to chatID
func (margelet *Margelet) SendImageByURL(chatID int64, url string, caption string, replyMarkup interface{}) (tgbotapi.Message, error) {
	resp, err := http.Get(url)

	if err != nil {
		return tgbotapi.Message{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return tgbotapi.Message{}, fmt.Errorf("Error obtaining image by url %s", url)
	}
	reader := tgbotapi.FileReader{Name: fmt.Sprintf(filepath.Base(url)), Reader: resp.Body, Size: resp.ContentLength}

	cfg := tgbotapi.NewPhotoUpload(chatID, reader)
	cfg.Caption = caption
	cfg.ReplyMarkup = replyMarkup
	return margelet.Send(cfg)
}

// SendTypingAction - sends typing chat action to chatID
func (margelet *Margelet) SendTypingAction(chatID int64) error {
	_, err := margelet.bot.Send(tgbotapi.NewChatAction(chatID, "typing"))
	return err
}

// SendUploadPhotoAction - sends upload_photo chat action to chatID
func (margelet *Margelet) SendUploadPhotoAction(chatID int64) error {
	_, err := margelet.bot.Send(tgbotapi.NewChatAction(chatID, "upload_photo"))
	return err
}

// SendRecordVideoAction - sends record_video chat action to chatID
func (margelet *Margelet) SendRecordVideoAction(chatID int64) error {
	_, err := margelet.bot.Send(tgbotapi.NewChatAction(chatID, "record_video"))
	return err
}

// SendUploadVideoAction - sends upload_video chat action to chatID
func (margelet *Margelet) SendUploadVideoAction(chatID int64) error {
	_, err := margelet.bot.Send(tgbotapi.NewChatAction(chatID, "upload_video"))
	return err
}

// SendRecordAudioAction - sends record_audio chat action to chatID
func (margelet *Margelet) SendRecordAudioAction(chatID int64) error {
	_, err := margelet.bot.Send(tgbotapi.NewChatAction(chatID, "record_audio"))
	return err
}

// SendUploadAudioAction - sends upload_audio chat action to chatID
func (margelet *Margelet) SendUploadAudioAction(chatID int64) error {
	_, err := margelet.bot.Send(tgbotapi.NewChatAction(chatID, "upload_audio"))
	return err
}

// SendUploadDocumentAction - sends upload_document chat action to chatID
func (margelet *Margelet) SendUploadDocumentAction(chatID int64) error {
	_, err := margelet.bot.Send(tgbotapi.NewChatAction(chatID, "upload_document"))
	return err
}

// SendFindLocationAction - sends find_location chat action to chatID
func (margelet *Margelet) SendFindLocationAction(chatID int64) error {
	_, err := margelet.bot.Send(tgbotapi.NewChatAction(chatID, "find_location"))
	return err
}

// SendHideKeyboard - hides keyboard in chatID
func (margelet *Margelet) SendHideKeyboard(chatID int64, message string) error {
	msg := tgbotapi.NewMessage(chatID, message)
	msg.ReplyMarkup = tgbotapi.NewHideKeyboard(true)
	_, err := margelet.Send(msg)
	return err
}

// RawBot - returns low-level tgbotapi.Bot
func (margelet *Margelet) RawBot() *tgbotapi.BotAPI {
	return margelet.bot.(*tgbotapi.BotAPI)
}

// GetUserProfilePhotos - returns user profile photos by config
func (margelet *Margelet) GetUserProfilePhotos(config tgbotapi.UserProfilePhotosConfig) (tgbotapi.UserProfilePhotos, error) {
	return margelet.RawBot().GetUserProfilePhotos(config)
}

// GetCurrentUserpic - returs current userpic for userID or error
func (margelet *Margelet) GetCurrentUserpic(userID int) (string, error) {
	photos, err := margelet.GetUserProfilePhotos(tgbotapi.UserProfilePhotosConfig{UserID: userID, Offset: 0, Limit: 1})
	if err != nil {
		return "", err
	}

	if len(photos.Photos) > 0 {
		p := photos.Photos[len(photos.Photos)-1]
		return margelet.GetFileDirectURL(p[len(p)-1].FileID)
	}
	return "", fmt.Errorf("No userpic found")
}
