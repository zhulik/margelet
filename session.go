package margelet

import (
	"gopkg.in/telegram-bot-api.v4"
)

// MargeletSession incapsulates user's session
type margeletSession struct {
	bot         MargeletAPI
	chatID      int64
	userID      int
	lastMessage *tgbotapi.Message
	responses   []tgbotapi.Message
}

func newMargetletSession(bot MargeletAPI, message *tgbotapi.Message, responses []tgbotapi.Message) Session {
	return &margeletSession{
		bot:         bot,
		chatID:      message.Chat.ID,
		userID:      message.From.ID,
		lastMessage: message,
		responses:   responses,
	}
}

// Responses returns all user's responses in session
func (s *margeletSession) Responses() []tgbotapi.Message {
	return s.responses
}

// QuickSend send test to session chat
func (s *margeletSession) QuickSend(text string) (tgbotapi.Message, error) {
	return s.bot.QuickSend(s.chatID, text)
}

// QuckReply send a reply to last session message
func (s *margeletSession) QuckReply(text string) (tgbotapi.Message, error) {
	return s.bot.QuickReply(s.chatID, s.lastMessage.MessageID, text)
}

// SendImageByURL send image by url to session chat
func (s *margeletSession) SendImageByURL(url string, caption string, replyMarkup interface{}) (tgbotapi.Message, error) {
	return s.bot.SendImageByURL(s.chatID, url, caption, replyMarkup)
}

// SendTypingAction send typing action to session chat
func (s *margeletSession) SendTypingAction() error {
	return s.bot.SendTypingAction(s.chatID)
}

// SendTypingAction send upload photo action to session chat
func (s *margeletSession) SendUploadPhotoAction() error {
	return s.bot.SendUploadPhotoAction(s.chatID)
}

// SendRecordVideoAction send record video action to session chat
func (s *margeletSession) SendRecordVideoAction() error {
	return s.bot.SendRecordVideoAction(s.chatID)
}

// SendUploadVideoAction send upload video action to session chat
func (s *margeletSession) SendUploadVideoAction() error {
	return s.bot.SendUploadVideoAction(s.chatID)
}

// SendRecordAudioAction send record audio action to session chat
func (s *margeletSession) SendRecordAudioAction() error {
	return s.bot.SendRecordAudioAction(s.chatID)
}

// SendUploadAudioAction send upload audio action to session chat
func (s *margeletSession) SendUploadAudioAction() error {
	return s.bot.SendUploadAudioAction(s.chatID)
}

// SendUploadDocumentAction send upload document action to session chat
func (s *margeletSession) SendUploadDocumentAction() error {
	return s.bot.SendUploadDocumentAction(s.chatID)
}

// SendFindLocationAction send find location action to session chat
func (s *margeletSession) SendFindLocationAction() error {
	return s.bot.SendFindLocationAction(s.chatID)
}

// SendHideKeyboard send message with hidding keyboard to session chat
func (s *margeletSession) SendHideKeyboard(chatID int64, message string) error {
	return s.bot.SendHideKeyboard(s.chatID, message)
}

func (s *margeletSession) Message() *tgbotapi.Message {
	return s.lastMessage
}
