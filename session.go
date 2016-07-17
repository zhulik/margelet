package margelet

import (
	"gopkg.in/telegram-bot-api.v4"
)

// MargeletSession incapsulates user's session
type MargeletSession struct {
	bot           MargeletAPI
	chatID        int64
	userID        int
	lastMessageID int
	responses     []tgbotapi.Message
}

func newMargetletSession(bot MargeletAPI, message *tgbotapi.Message, responses []tgbotapi.Message) Session {
	return &MargeletSession{
		bot:           bot,
		chatID:        message.Chat.ID,
		userID:        message.From.ID,
		lastMessageID: message.MessageID,
		responses:     responses,
	}
}

// Responses returns all user's responses in session
func (s *MargeletSession) Responses() []tgbotapi.Message {
	return s.responses
}

// QuickSend send test to session chat
func (s *MargeletSession) QuickSend(text string) (tgbotapi.Message, error) {
	return s.bot.QuickSend(s.chatID, text)
}

// QuckReply send a reply to last session message
func (s *MargeletSession) QuckReply(text string) (tgbotapi.Message, error) {
	return s.bot.QuickReply(s.chatID, s.lastMessageID, text)
}

// SendImageByURL send image by url to session chat
func (s *MargeletSession) SendImageByURL(url string, caption string, replyMarkup interface{}) (tgbotapi.Message, error) {
	return s.bot.SendImageByURL(s.chatID, url, caption, replyMarkup)
}

// SendTypingAction send typing action to session chat
func (s *MargeletSession) SendTypingAction() error {
	return s.bot.SendTypingAction(s.chatID)
}

// SendTypingAction send upload photo action to session chat
func (s *MargeletSession) SendUploadPhotoAction() error {
	return s.bot.SendUploadPhotoAction(s.chatID)
}

// SendRecordVideoAction send record video action to session chat
func (s *MargeletSession) SendRecordVideoAction() error {
	return s.bot.SendRecordVideoAction(s.chatID)
}

// SendUploadVideoAction send upload video action to session chat
func (s *MargeletSession) SendUploadVideoAction() error {
	return s.bot.SendUploadVideoAction(s.chatID)
}

// SendRecordAudioAction send record audio action to session chat
func (s *MargeletSession) SendRecordAudioAction() error {
	return s.bot.SendRecordAudioAction(s.chatID)
}

// SendUploadAudioAction send upload audio action to session chat
func (s *MargeletSession) SendUploadAudioAction() error {
	return s.bot.SendUploadAudioAction(s.chatID)
}

// SendUploadDocumentAction send upload document action to session chat
func (s *MargeletSession) SendUploadDocumentAction() error {
	return s.bot.SendUploadDocumentAction(s.chatID)
}

// SendFindLocationAction send find location action to session chat
func (s *MargeletSession) SendFindLocationAction() error {
	return s.bot.SendFindLocationAction(s.chatID)
}

// SendHideKeyboard send message with hidding keyboard to session chat
func (s *MargeletSession) SendHideKeyboard(chatID int64, message string) error {
	return s.bot.SendHideKeyboard(s.chatID, message)
}
