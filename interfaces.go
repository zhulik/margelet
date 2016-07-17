package margelet

import (
	"gopkg.in/redis.v3"
	"gopkg.in/telegram-bot-api.v4"
)

// MessageHandler - interface for message handlers
type MessageHandler interface {
	HandleMessage(bot MargeletAPI, message *tgbotapi.Message) error
}

// InlineHandler - interface for message handlers
type InlineHandler interface {
	HandleInline(bot MargeletAPI, query *tgbotapi.InlineQuery) error
}

// CallbackHandler - interface for message handlers
type CallbackHandler interface {
	HandleCallback(bot MargeletAPI, query *tgbotapi.CallbackQuery) error
}

// CommandHandler - interface for command handlers
type CommandHandler interface {
	HandleCommand(bot MargeletAPI, message *tgbotapi.Message) error
	HelpMessage() string
}

// SessionHandler - interface for session handlers
type SessionHandler interface {
	HandleSession(bot MargeletAPI, session Session) error
	CancelSession(bot MargeletAPI, session Session)
	HelpMessage() string
}

// MargeletAPI - interface, that describes margelet API
type MargeletAPI interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
	AnswerInlineQuery(config tgbotapi.InlineConfig) (tgbotapi.APIResponse, error)
	AnswerCallbackQuery(config tgbotapi.CallbackConfig) (tgbotapi.APIResponse, error)
	QuickSend(chatID int64, message string) (tgbotapi.Message, error)
	QuickReply(chatID int64, messageID int, message string) (tgbotapi.Message, error)
	QuickForceReply(chatID int64, messageID int, message string) (tgbotapi.Message, error)

	GetFileDirectURL(fileID string) (string, error)
	IsMessageToMe(message tgbotapi.Message) bool
	GetConfigRepository() *ChatConfigRepository
	GetSessionRepository() SessionRepository
	GetRedis() *redis.Client
	HandleSession(message *tgbotapi.Message, command string)
	StartSession(message *tgbotapi.Message, command string)
	SendImageByURL(chatID int64, url string, caption string, replyMarkup interface{}) (tgbotapi.Message, error)

	SendTypingAction(chatID int64) error
	SendUploadPhotoAction(chatID int64) error
	SendRecordVideoAction(chatID int64) error
	SendUploadVideoAction(chatID int64) error
	SendRecordAudioAction(chatID int64) error
	SendUploadAudioAction(chatID int64) error
	SendUploadDocumentAction(chatID int64) error
	SendFindLocationAction(chatID int64) error

	SendHideKeyboard(chatID int64, message string) error

	RawBot() *tgbotapi.BotAPI
	GetUserProfilePhotos(config tgbotapi.UserProfilePhotosConfig) (tgbotapi.UserProfilePhotos, error)
	GetCurrentUserpic(userID int) (string, error)
}

// TGBotAPI - interface, that describe telegram-bot-api API
type TGBotAPI interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
	AnswerInlineQuery(config tgbotapi.InlineConfig) (tgbotapi.APIResponse, error)
	AnswerCallbackQuery(config tgbotapi.CallbackConfig) (tgbotapi.APIResponse, error)
	GetFileDirectURL(fileID string) (string, error)
	IsMessageToMe(message tgbotapi.Message) bool
	GetUpdatesChan(config tgbotapi.UpdateConfig) (<-chan tgbotapi.Update, error)
}

// AuthorizationPolicy - interface, that describes authorization policy for command or session
type AuthorizationPolicy interface {
	Allow(message *tgbotapi.Message) error
}

// Session - interface, that describes incapsulated info aboud user's session with bot
type Session interface {
	Responses() []tgbotapi.Message
	Message() *tgbotapi.Message
	Finish()

	QuickSend(text string) (tgbotapi.Message, error)
	QuckReply(text string) (tgbotapi.Message, error)
	QuickForceReply(text string) (tgbotapi.Message, error)
	// SendImageByURL send image by url to session chat
	SendImageByURL(url string, caption string, replyMarkup interface{}) (tgbotapi.Message, error)

	// SendTypingAction send typing action to session chat
	SendTypingAction() error

	// SendTypingAction send upload photo action to session chat
	SendUploadPhotoAction() error

	// SendRecordVideoAction send record video action to session chat
	SendRecordVideoAction() error

	// SendUploadVideoAction send upload video action to session chat
	SendUploadVideoAction() error

	// SendRecordAudioAction send record audio action to session chat
	SendRecordAudioAction() error

	// SendUploadAudioAction send upload audio action to session chat
	SendUploadAudioAction() error

	// SendUploadDocumentAction send upload document action to session chat
	SendUploadDocumentAction() error

	// SendFindLocationAction send find location action to session chat
	SendFindLocationAction() error

	// SendHideKeyboard send message with hidding keyboard to session chat
	SendHideKeyboard(chatID int64, message string) error
}
