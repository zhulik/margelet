package margelet

import (
	"gopkg.in/redis.v3"
	"gopkg.in/telegram-bot-api.v4"
)

type message struct {
	bot     MargeletAPI
	message *tgbotapi.Message
}

func newMessage(bot MargeletAPI, msg *tgbotapi.Message) *message {
	return &message{
		bot:     bot,
		message: msg,
	}
}

func (s *message) Bot() MargeletAPI {
	return s.bot
}

// QuickSend send test to session chat
func (s *message) QuickSend(text string, replyMarkup ...interface{}) (tgbotapi.Message, error) {
	return s.bot.QuickSend(s.message.Chat.ID, text, replyMarkup)
}

// QuckReply send a reply to last session message
func (s *message) QuickReply(text string, replyMarkup ...interface{}) (tgbotapi.Message, error) {
	return s.bot.QuickReply(s.message.Chat.ID, s.message.MessageID, text, replyMarkup)
}

func (s *message) QuickForceReply(text string) (tgbotapi.Message, error) {
	return s.bot.QuickForceReply(s.message.Chat.ID, s.message.MessageID, text)
}

// SendImageByURL send image by url to session chat
func (s *message) SendImageByURL(url string, caption string, replyMarkup interface{}) (tgbotapi.Message, error) {
	return s.bot.SendImageByURL(s.message.Chat.ID, url, caption, replyMarkup)
}

// SendTypingAction send typing action to session chat
func (s *message) SendTypingAction() error {
	return s.bot.SendTypingAction(s.message.Chat.ID)
}

// SendTypingAction send upload photo action to session chat
func (s *message) SendUploadPhotoAction() error {
	return s.bot.SendUploadPhotoAction(s.message.Chat.ID)
}

// SendRecordVideoAction send record video action to session chat
func (s *message) SendRecordVideoAction() error {
	return s.bot.SendRecordVideoAction(s.message.Chat.ID)
}

// SendUploadVideoAction send upload video action to session chat
func (s *message) SendUploadVideoAction() error {
	return s.bot.SendUploadVideoAction(s.message.Chat.ID)
}

// SendRecordAudioAction send record audio action to session chat
func (s *message) SendRecordAudioAction() error {
	return s.bot.SendRecordAudioAction(s.message.Chat.ID)
}

// SendUploadAudioAction send upload audio action to session chat
func (s *message) SendUploadAudioAction() error {
	return s.bot.SendUploadAudioAction(s.message.Chat.ID)
}

// SendUploadDocumentAction send upload document action to session chat
func (s *message) SendUploadDocumentAction() error {
	return s.bot.SendUploadDocumentAction(s.message.Chat.ID)
}

// SendFindLocationAction send find location action to session chat
func (s *message) SendFindLocationAction() error {
	return s.bot.SendFindLocationAction(s.message.Chat.ID)
}

// SendHideKeyboard send message with hidding keyboard to session chat
func (s *message) SendHideKeyboard(message string) error {
	return s.bot.SendHideKeyboard(s.message.Chat.ID, message)
}

// Message returns user's message
func (s *message) Message() *tgbotapi.Message {
	return s.message
}

func (s *message) GetFileDirectURL(fileID string) (string, error) {
	return s.bot.GetFileDirectURL(fileID)
}

func (s *message) GetCurrentUserpic() (string, error) {
	return s.bot.GetCurrentUserpic(s.message.From.ID)
}

// GetConfigRepository - returns chat config repository
func (s *message) GetConfigRepository() *ChatConfigRepository {
	return s.bot.GetConfigRepository()
}

// GetSessionRepository - returns session repository
func (s *message) GetSessionRepository() SessionRepository {
	return s.bot.GetSessionRepository()
}

// GetRedis - returns margelet's redis client
func (s *message) GetRedis() *redis.Client {
	return s.bot.GetRedis()
}

func (s *message) StartSession(command string) {
	s.bot.StartSession(s.message, command)
}

func (s *message) GetChatRepository() *ChatRepository {
	return s.bot.GetChatRepository()
}
