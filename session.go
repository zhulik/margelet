package margelet

import (
	"gopkg.in/telegram-bot-api.v4"
)

type margeletSession struct {
	bot           MargeletAPI
	chatID        int64
	userID        int
	lastMessageID int
	responses     []tgbotapi.Message
}

func newMargetletSession(bot MargeletAPI, message *tgbotapi.Message, responses []tgbotapi.Message) Session {
	return &margeletSession{
		bot:           bot,
		chatID:        message.Chat.ID,
		userID:        message.From.ID,
		lastMessageID: message.MessageID,
		responses:     responses,
	}
}

func (s *margeletSession) Responses() []tgbotapi.Message {
	return s.responses
}

func (s *margeletSession) QuickSend(text string) (tgbotapi.Message, error) {
	return s.bot.QuickSend(s.chatID, text)
}

func (s *margeletSession) QuckReply(text string) (tgbotapi.Message, error) {
	return s.bot.QuickReply(s.chatID, s.lastMessageID, text)
}

func (s *margeletSession) SendImageByURL(url string, caption string, replyMarkup interface{}) (tgbotapi.Message, error) {
	return s.bot.SendImageByURL(s.chatID, url, caption, replyMarkup)
}

func (s *margeletSession) SendTypingAction() error {
	return s.bot.SendTypingAction(s.chatID)
}

func (s *margeletSession) SendUploadPhotoAction() error {
	return s.bot.SendUploadPhotoAction(s.chatID)
}

func (s *margeletSession) SendRecordVideoAction() error {
	return s.bot.SendRecordVideoAction(s.chatID)
}

func (s *margeletSession) SendUploadVideoAction() error {
	return s.bot.SendUploadVideoAction(s.chatID)
}

func (s *margeletSession) SendRecordAudioAction() error {
	return s.bot.SendRecordAudioAction(s.chatID)
}

func (s *margeletSession) SendUploadAudioAction() error {
	return s.bot.SendUploadAudioAction(s.chatID)
}

func (s *margeletSession) SendUploadDocumentAction() error {
	return s.bot.SendUploadDocumentAction(s.chatID)
}

func (s *margeletSession) SendFindLocationAction() error {
	return s.bot.SendFindLocationAction(s.chatID)
}

func (s *margeletSession) SendHideKeyboard(chatID int64, message string) error {
	return s.bot.SendHideKeyboard(s.chatID, message)
}
