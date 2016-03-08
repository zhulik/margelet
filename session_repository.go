package margelet

import (
	"encoding/json"
	"fmt"
	"gopkg.in/redis.v3"
	"gopkg.in/telegram-bot-api.v2"
	"strings"
)

// SessionRepository - public interface for session repository
type SessionRepository interface {
	Create(chatID int, userID int, command string)
	Add(chatID int, userID int, message tgbotapi.Message)
	Remove(chatID int, userID int)
	Command(chatID int, userID int) string
	Dialog(chatID int, userID int) (messages []tgbotapi.Message)
}

type sessionRepository struct {
	key   string
	redis *redis.Client
}

func newSessionRepository(prefix string, redis *redis.Client) *sessionRepository {
	key := strings.Join([]string{prefix, "margelet_sessions"}, "-")
	return &sessionRepository{key, redis}
}

// Create - creates new session for chatID, userID and command
func (session *sessionRepository) Create(chatID int, userID int, command string) {
	key := session.keyFor(chatID, userID)
	session.redis.Set(key, command, 0)
}

// Add - adds user's answer to existing session
func (session *sessionRepository) Add(chatID int, userID int, message tgbotapi.Message) {
	key := session.dialogKeyFor(chatID, userID)

	json, _ := json.Marshal(message)

	session.redis.RPush(key, string(json))
}

// Remove - removes session
func (session *sessionRepository) Remove(chatID int, userID int) {
	key := session.keyFor(chatID, userID)
	session.redis.Del(key)

	key = session.dialogKeyFor(chatID, userID)
	session.redis.Del(key)
}

// Command - returns command for active session for chatID and userID, if exists
// otherwise returns empty string
func (session *sessionRepository) Command(chatID int, userID int) string {
	key := session.keyFor(chatID, userID)
	value, _ := session.redis.Get(key).Result()
	return value
}

// Dialog returns all user's answers history for chatID and userID
func (session *sessionRepository) Dialog(chatID int, userID int) (messages []tgbotapi.Message) {
	key := session.dialogKeyFor(chatID, userID)

	values := session.redis.LRange(key, 0, -1).Val()
	for _, value := range values {
		msg := tgbotapi.Message{}
		json.Unmarshal([]byte(value), &msg)
		messages = append(messages, msg)
	}
	return
}

func (session *sessionRepository) keyFor(chatID int, userID int) string {
	return fmt.Sprintf("%s_%d_%d", session.key, chatID, userID)
}

func (session *sessionRepository) dialogKeyFor(chatID int, userID int) string {
	return fmt.Sprintf("%s_dialog", session.keyFor(chatID, userID))
}
