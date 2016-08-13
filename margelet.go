package margelet

import (
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/redis.v3"
	"gopkg.in/telegram-bot-api.v4"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
)

type policies []AuthorizationPolicy

// RecoverCallback - callback wich will be called when margelet recovers from panic
type RecoverCallback func(margelet *Margelet, userID int, r interface{})
type ReceiveCallback func(from int, text string)
type SendCallback func(to int64, text string)

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

	RecoverCallback       RecoverCallback
	ReceiveCallback       ReceiveCallback
	SendCallback          SendCallback
	MessageHandlers       []MessageHandler
	CommandHandlers       map[string]authorizedCommandHandler
	UnknownCommandHandler *authorizedCommandHandler
	SessionHandlers       map[string]authorizedSessionHandler
	InlineHandler         InlineHandler
	CallbackHandler       CallbackHandler
	running               bool
	verbose               bool
	Redis                 *redis.Client
	ChatRepository        *ChatRepository
	SessionRepository     SessionRepository
	ChatConfigRepository  *ChatConfigRepository
	StatsRepository       StatsRepository
	token                 string
}

// NewMargelet creates new Margelet instance
func NewMargelet(botName string, redisAddr string, redisPassword string, redisDB int64, token string, verbose bool) (*Margelet, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	bot.Debug = verbose
	m, err := NewMargeletFromBot(botName, redisAddr, redisPassword, redisDB, bot, verbose)

	m.token = token
	return m, err
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
	statsRepository := newStatsRepository(botName, redis)

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
		StatsRepository:      statsRepository,
		ReceiveCallback:      func(f int, t string) {},
		SendCallback:         func(f int64, t string) {},
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

// SetUnknownCommandHandler - sets unknown command handler
func (margelet *Margelet) SetUnknownCommandHandler(handler CommandHandler, auth ...AuthorizationPolicy) {
	margelet.UnknownCommandHandler = &authorizedCommandHandler{auth, handler}
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
	v, err := strconv.Atoi(config.InlineQueryID)
	if err != nil {
		go margelet.SendCallback(int64(v), config.InlineQueryID)
	}
	return margelet.bot.AnswerInlineQuery(config)
}

// AnswerCallbackQuery  - send answer to CallbackQuery
func (margelet *Margelet) AnswerCallbackQuery(config tgbotapi.CallbackConfig) (tgbotapi.APIResponse, error) {
	v, err := strconv.Atoi(config.CallbackQueryID)
	if err != nil {
		go margelet.SendCallback(int64(v), config.Text)
	}
	return margelet.bot.AnswerCallbackQuery(config)
}

// QuickSend - quick send text message to chatID
func (margelet *Margelet) QuickSend(chatID int64, message string, replyMarkup interface{}) (tgbotapi.Message, error) {
	go margelet.SendCallback(chatID, message)
	msg := tgbotapi.NewMessage(chatID, message)
	msg.ReplyMarkup = replyMarkup
	return margelet.bot.Send(msg)
}

// QuickReply - quick send text reply to message
func (margelet *Margelet) QuickReply(chatID int64, messageID int, message string, replyMarkup interface{}) (tgbotapi.Message, error) {
	go margelet.SendCallback(chatID, message)
	msg := tgbotapi.NewMessage(chatID, message)
	msg.ReplyToMessageID = messageID
	msg.ReplyMarkup = replyMarkup
	return margelet.bot.Send(msg)
}

// QuickForceReply - quick send text force reply to message
func (margelet *Margelet) QuickForceReply(chatID int64, messageID int, message string) (tgbotapi.Message, error) {
	go margelet.SendCallback(chatID, message)
	msg := tgbotapi.NewMessage(chatID, message)
	msg.ReplyToMessageID = messageID
	msg.ReplyMarkup = tgbotapi.ForceReply{true, true}
	return margelet.bot.Send(msg)
}

// SendImage - sends image to chat
func (margelet *Margelet) SendImage(chatID int64, reader tgbotapi.FileReader, caption string, replyMarkup interface{}) (tgbotapi.Message, error) {
	go margelet.SendCallback(chatID, reader.Name)
	msg := tgbotapi.NewPhotoUpload(chatID, reader)
	msg.ReplyMarkup = replyMarkup
	msg.Caption = caption
	return margelet.Send(msg)
}

// SendDocument - sends document to chat
func (margelet *Margelet) SendDocument(chatID int64, reader tgbotapi.FileReader, replyMarkup interface{}) (tgbotapi.Message, error) {
	go margelet.SendCallback(chatID, reader.Name)
	msg := tgbotapi.NewDocumentUpload(chatID, reader)
	msg.ReplyMarkup = replyMarkup
	return margelet.Send(msg)
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

// GetStatsRepository - returns stats repository
func (margelet *Margelet) GetStatsRepository() StatsRepository {
	return margelet.StatsRepository
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
		msg := <-updates
		go handleUpdate(margelet, msg)
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
func (margelet *Margelet) SendImageByURL(chatID int64, imgURL string, caption string, replyMarkup interface{}) (tgbotapi.Message, error) {
	margelet.SendUploadPhotoAction(chatID)
	go margelet.SendCallback(chatID, imgURL)

	rmData, err := json.Marshal(replyMarkup)
	if err != nil {
		return tgbotapi.Message{}, err
	}

	params := url.Values{}
	params.Add("chat_id", strconv.FormatInt(chatID, 10))
	params.Add("photo", imgURL)
	params.Add("caption", caption)
	params.Add("reply_markup", string(rmData))
	method := fmt.Sprintf("https://api.telegram.org/beta/bot%s/sendPhoto", margelet.token)
	client := new(http.Client)
	resp, err := client.PostForm(method, params)
	if err != nil {
		return tgbotapi.Message{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden {
		return tgbotapi.Message{}, errors.New(tgbotapi.ErrAPIForbidden)
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return tgbotapi.Message{}, err
	}

	log.Println(method, string(bytes))

	var apiResp tgbotapi.APIResponse
	json.Unmarshal(bytes, &apiResp)

	if !apiResp.Ok {
		return tgbotapi.Message{}, errors.New(apiResp.Description)
	}

	return tgbotapi.Message{}, nil
	// resp, err := http.Get(url)
	//
	// if err != nil {
	// 	return tgbotapi.Message{}, err
	// }
	// defer resp.Body.Close()
	// if resp.StatusCode != 200 {
	// 	return tgbotapi.Message{}, fmt.Errorf("Error obtaining image by url %s", url)
	// }
	// reader := tgbotapi.FileReader{Name: fmt.Sprintf(filepath.Base(url)), Reader: resp.Body, Size: resp.ContentLength}
	//
	// cfg := tgbotapi.NewPhotoUpload(chatID, reader)
	// cfg.Caption = caption
	// cfg.ReplyMarkup = replyMarkup
	// return margelet.Send(cfg)
}

// SendDocumentByURL - sends given by url document to chatID
func (margelet *Margelet) SendDocumentByURL(chatID int64, url string, replyMarkup interface{}) (tgbotapi.Message, error) {
	margelet.SendUploadDocumentAction(chatID)
	go margelet.SendCallback(chatID, url)
	resp, err := http.Get(url)

	if err != nil {
		return tgbotapi.Message{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return tgbotapi.Message{}, fmt.Errorf("Error obtaining document by url %s", url)
	}
	reader := tgbotapi.FileReader{Name: fmt.Sprintf(filepath.Base(url)), Reader: resp.Body, Size: resp.ContentLength}

	cfg := tgbotapi.NewDocumentUpload(chatID, reader)
	cfg.ReplyMarkup = replyMarkup
	return margelet.Send(cfg)
}

// SendImageByID - sends given by FileID image to chatID
func (margelet *Margelet) SendImageByID(chatID int64, fileID string, caption string, replyMarkup interface{}) (tgbotapi.Message, error) {
	go margelet.SendCallback(chatID, fileID)
	cfg := tgbotapi.NewPhotoShare(chatID, fileID)
	cfg.Caption = caption
	cfg.ReplyMarkup = replyMarkup
	return margelet.Send(cfg)
}

// SendTypingAction - sends typing chat action to chatID
func (margelet *Margelet) SendTypingAction(chatID int64) error {
	go margelet.SendCallback(chatID, "typing")
	_, err := margelet.bot.Send(tgbotapi.NewChatAction(chatID, "typing"))
	return err
}

// SendUploadPhotoAction - sends upload_photo chat action to chatID
func (margelet *Margelet) SendUploadPhotoAction(chatID int64) error {
	go margelet.SendCallback(chatID, "upload_photo")
	_, err := margelet.bot.Send(tgbotapi.NewChatAction(chatID, "upload_photo"))
	return err
}

// SendRecordVideoAction - sends record_video chat action to chatID
func (margelet *Margelet) SendRecordVideoAction(chatID int64) error {
	go margelet.SendCallback(chatID, "record_video")
	_, err := margelet.bot.Send(tgbotapi.NewChatAction(chatID, "record_video"))
	return err
}

// SendUploadVideoAction - sends upload_video chat action to chatID
func (margelet *Margelet) SendUploadVideoAction(chatID int64) error {
	go margelet.SendCallback(chatID, "record_video")
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

// GetCurrentUserpic - returns current userpic URL for userID or error
func (margelet *Margelet) GetCurrentUserpic(userID int) (string, error) {
	fileID, err := margelet.GetCurrentUserpicID(userID)
	if err != nil {
		return "", err
	}

	return margelet.GetFileDirectURL(fileID)
}

// GetCurrentUserpicID - returns current userpic FileID for userID or error
func (margelet *Margelet) GetCurrentUserpicID(userID int) (string, error) {
	photos, err := margelet.GetUserProfilePhotos(tgbotapi.UserProfilePhotosConfig{UserID: userID, Offset: 0, Limit: 1})
	if err != nil {
		return "", err
	}

	if len(photos.Photos) > 0 {
		p := photos.Photos[len(photos.Photos)-1]
		return p[len(p)-1].FileID, nil
	}
	return "", fmt.Errorf("No userpic found")
}

// GetChatRepository - returns chats repository
func (margelet *Margelet) GetChatRepository() *ChatRepository {
	return margelet.ChatRepository
}
