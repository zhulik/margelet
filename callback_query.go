package margelet

import (
	"gopkg.in/redis.v3"
	"gopkg.in/telegram-bot-api.v4"
)

type callbackQuery struct {
	bot   MargeletAPI
	query *tgbotapi.CallbackQuery
}

func newCallbackQuery(bot MargeletAPI, query *tgbotapi.CallbackQuery) *callbackQuery {
	return &callbackQuery{
		bot:   bot,
		query: query,
	}
}

// Query returns tgbotapi query
func (s *callbackQuery) Query() *tgbotapi.CallbackQuery {
	return s.query
}

func (s *callbackQuery) Bot() MargeletAPI {
	return s.bot
}

// Message returns user's message
func (s *callbackQuery) Message() *tgbotapi.Message {
	return s.query.Message
}

// QuickSend send test to session chat
func (s *callbackQuery) QuickSend(text string) (tgbotapi.Message, error) {
	return s.bot.QuickSend(s.query.Message.Chat.ID, text)
}

// QuckReply send a reply to last session message
func (s *callbackQuery) QuickReply(text string) (tgbotapi.Message, error) {
	return s.bot.QuickReply(s.query.Message.Chat.ID, s.query.Message.MessageID, text)
}

func (s *callbackQuery) QuickForceReply(text string) (tgbotapi.Message, error) {
	return s.bot.QuickForceReply(s.query.Message.Chat.ID, s.query.Message.MessageID, text)
}

// SendImageByURL send image by url to session chat
func (s *callbackQuery) SendImageByURL(url string, caption string, replyMarkup interface{}) (tgbotapi.Message, error) {
	return s.bot.SendImageByURL(s.query.Message.Chat.ID, url, caption, replyMarkup)
}

// SendTypingAction send typing action to session chat
func (s *callbackQuery) SendTypingAction() error {
	return s.bot.SendTypingAction(s.query.Message.Chat.ID)
}

// SendTypingAction send upload photo action to session chat
func (s *callbackQuery) SendUploadPhotoAction() error {
	return s.bot.SendUploadPhotoAction(s.query.Message.Chat.ID)
}

// SendRecordVideoAction send record video action to session chat
func (s *callbackQuery) SendRecordVideoAction() error {
	return s.bot.SendRecordVideoAction(s.query.Message.Chat.ID)
}

// SendUploadVideoAction send upload video action to session chat
func (s *callbackQuery) SendUploadVideoAction() error {
	return s.bot.SendUploadVideoAction(s.query.Message.Chat.ID)
}

// SendRecordAudioAction send record audio action to session chat
func (s *callbackQuery) SendRecordAudioAction() error {
	return s.bot.SendRecordAudioAction(s.query.Message.Chat.ID)
}

// SendUploadAudioAction send upload audio action to session chat
func (s *callbackQuery) SendUploadAudioAction() error {
	return s.bot.SendUploadAudioAction(s.query.Message.Chat.ID)
}

// SendUploadDocumentAction send upload document action to session chat
func (s *callbackQuery) SendUploadDocumentAction() error {
	return s.bot.SendUploadDocumentAction(s.query.Message.Chat.ID)
}

// SendFindLocationAction send find location action to session chat
func (s *callbackQuery) SendFindLocationAction() error {
	return s.bot.SendFindLocationAction(s.query.Message.Chat.ID)
}

// SendHideKeyboard send message with hidding keyboard to session chat
func (s *callbackQuery) SendHideKeyboard(message string) error {
	return s.bot.SendHideKeyboard(s.query.Message.Chat.ID, message)
}

func (s *callbackQuery) GetFileDirectURL(fileID string) (string, error) {
	return s.bot.GetFileDirectURL(fileID)
}

func (s *callbackQuery) Data() string {
	return s.query.Data
}

// GetConfigRepository - returns chat config repository
func (s *callbackQuery) GetConfigRepository() *ChatConfigRepository {
	return s.bot.GetConfigRepository()
}

// GetSessionRepository - returns session repository
func (s *callbackQuery) GetSessionRepository() SessionRepository {
	return s.bot.GetSessionRepository()
}

// GetRedis - returns margelet's redis client
func (s *callbackQuery) GetRedis() *redis.Client {
	return s.bot.GetRedis()
}

func (s *callbackQuery) GetCurrentUserpic() (string, error) {
	return s.bot.GetCurrentUserpic(s.query.Message.From.ID)
}
