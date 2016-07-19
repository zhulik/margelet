package margelet

import (
	"gopkg.in/redis.v3"
	"gopkg.in/telegram-bot-api.v4"
)

// MessageHandler - interface for message handlers
type MessageHandler interface {
	HandleMessage(message Message) error
}

// InlineHandler - interface for message handlers
type InlineHandler interface {
	HandleInline(bot MargeletAPI, query *tgbotapi.InlineQuery) error
}

// CallbackHandler - interface for message handlers
type CallbackHandler interface {
	HandleCallback(query CallbackQuery) error
}

// CommandHandler - interface for command handlers
type CommandHandler interface {
	HandleCommand(msg Message) error
	HelpMessage() string
}

// SessionHandler - interface for session handlers
type SessionHandler interface {
	HandleSession(session Session) error
	CancelSession(session Session)
	HelpMessage() string
}

type Store interface {
	GetConfigRepository() *ChatConfigRepository
	GetSessionRepository() SessionRepository
	GetChatRepository() *ChatRepository
	GetRedis() *redis.Client
}

// MargeletAPI - interface, that describes margelet API
type MargeletAPI interface {
	Store
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
	AnswerInlineQuery(config tgbotapi.InlineConfig) (tgbotapi.APIResponse, error)
	AnswerCallbackQuery(config tgbotapi.CallbackConfig) (tgbotapi.APIResponse, error)
	QuickSend(chatID int64, message string) (tgbotapi.Message, error)
	QuickReply(chatID int64, messageID int, message string) (tgbotapi.Message, error)
	QuickForceReply(chatID int64, messageID int, message string) (tgbotapi.Message, error)

	GetFileDirectURL(fileID string) (string, error)
	IsMessageToMe(message tgbotapi.Message) bool
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

// Message - interface, that describes incapsulated info aboud user's message with some helper methods
type Message interface {
	Store
	Message() *tgbotapi.Message
	GetFileDirectURL(fileID string) (string, error)
	QuickSend(text string) (tgbotapi.Message, error)
	QuickReply(text string) (tgbotapi.Message, error)
	QuickForceReply(text string) (tgbotapi.Message, error)
	SendImageByURL(url string, caption string, replyMarkup interface{}) (tgbotapi.Message, error)
	SendTypingAction() error
	SendUploadPhotoAction() error
	SendRecordVideoAction() error
	SendUploadVideoAction() error
	SendRecordAudioAction() error
	SendUploadAudioAction() error
	SendUploadDocumentAction() error
	SendFindLocationAction() error
	SendHideKeyboard(message string) error
	GetCurrentUserpic() (string, error)
	Bot() MargeletAPI
	StartSession(command string)
}

// Session - interface, that describes incapsulated info aboud user's session with bot
type Session interface {
	Message
	Responses() []tgbotapi.Message
	Finish()
}

type CallbackQuery interface {
	Message
	Query() *tgbotapi.CallbackQuery
	Data() string
}
