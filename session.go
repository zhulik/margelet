package margelet

import (
	"gopkg.in/redis.v3"
	"gopkg.in/telegram-bot-api.v4"
)

// MargeletSession incapsulates user's session
type margeletSession struct {
	bot         MargeletAPI
	lastMessage *tgbotapi.Message
	responses   []tgbotapi.Message
	finished    bool
}

func newMargetletSession(bot MargeletAPI, message *tgbotapi.Message, responses []tgbotapi.Message) *margeletSession {
	return &margeletSession{
		bot:         bot,
		lastMessage: message,
		responses:   responses,
		finished:    false,
	}
}

// Responses returns all user's responses in session
func (s *margeletSession) Responses() []tgbotapi.Message {
	return s.responses
}

func (s *margeletSession) Bot() MargeletAPI {
	return s.bot
}

// QuickSend send test to session chat
func (s *margeletSession) QuickSend(text string) (tgbotapi.Message, error) {
	return s.bot.QuickSend(s.lastMessage.Chat.ID, text)
}

// QuckReply send a reply to last session message
func (s *margeletSession) QuickReply(text string) (tgbotapi.Message, error) {
	return s.bot.QuickReply(s.lastMessage.Chat.ID, s.lastMessage.MessageID, text)
}

func (s *margeletSession) QuickForceReply(text string) (tgbotapi.Message, error) {
	return s.bot.QuickForceReply(s.lastMessage.Chat.ID, s.lastMessage.MessageID, text)
}

// SendImageByURL send image by url to session chat
func (s *margeletSession) SendImageByURL(url string, caption string, replyMarkup interface{}) (tgbotapi.Message, error) {
	return s.bot.SendImageByURL(s.lastMessage.Chat.ID, url, caption, replyMarkup)
}

// SendTypingAction send typing action to session chat
func (s *margeletSession) SendTypingAction() error {
	return s.bot.SendTypingAction(s.lastMessage.Chat.ID)
}

// SendTypingAction send upload photo action to session chat
func (s *margeletSession) SendUploadPhotoAction() error {
	return s.bot.SendUploadPhotoAction(s.lastMessage.Chat.ID)
}

// SendRecordVideoAction send record video action to session chat
func (s *margeletSession) SendRecordVideoAction() error {
	return s.bot.SendRecordVideoAction(s.lastMessage.Chat.ID)
}

// SendUploadVideoAction send upload video action to session chat
func (s *margeletSession) SendUploadVideoAction() error {
	return s.bot.SendUploadVideoAction(s.lastMessage.Chat.ID)
}

// SendRecordAudioAction send record audio action to session chat
func (s *margeletSession) SendRecordAudioAction() error {
	return s.bot.SendRecordAudioAction(s.lastMessage.Chat.ID)
}

// SendUploadAudioAction send upload audio action to session chat
func (s *margeletSession) SendUploadAudioAction() error {
	return s.bot.SendUploadAudioAction(s.lastMessage.Chat.ID)
}

// SendUploadDocumentAction send upload document action to session chat
func (s *margeletSession) SendUploadDocumentAction() error {
	return s.bot.SendUploadDocumentAction(s.lastMessage.Chat.ID)
}

// SendFindLocationAction send find location action to session chat
func (s *margeletSession) SendFindLocationAction() error {
	return s.bot.SendFindLocationAction(s.lastMessage.Chat.ID)
}

// SendHideKeyboard send message with hidding keyboard to session chat
func (s *margeletSession) SendHideKeyboard(message string) error {
	return s.bot.SendHideKeyboard(s.lastMessage.Chat.ID, message)
}

// Message returns user's message
func (s *margeletSession) Message() *tgbotapi.Message {
	return s.lastMessage
}

func (s *margeletSession) Finish() {
	s.finished = true
}

func (s *margeletSession) GetFileDirectURL(fileID string) (string, error) {
	return s.bot.GetFileDirectURL(fileID)
}

// GetConfigRepository - returns chat config repository
func (s *margeletSession) GetConfigRepository() *ChatConfigRepository {
	return s.bot.GetConfigRepository()
}

// GetSessionRepository - returns session repository
func (s *margeletSession) GetSessionRepository() SessionRepository {
	return s.bot.GetSessionRepository()
}

// GetRedis - returns margelet's redis client
func (s *margeletSession) GetRedis() *redis.Client {
	return s.bot.GetRedis()
}

func (s *margeletSession) GetCurrentUserpic() (string, error) {
	return s.bot.GetCurrentUserpic(s.lastMessage.From.ID)
}
