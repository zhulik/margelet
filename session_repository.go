package margelet

import (
	"fmt"
	"gopkg.in/redis.v3"
	"strings"
)

// SessionRepository - repository for sessions
type SessionRepository struct {
	key   string
	redis *redis.Client
}

func newSessionRepository(prefix string, redis *redis.Client) *SessionRepository {
	key := strings.Join([]string{prefix, "margelet_sessions"}, "-")
	return &SessionRepository{key, redis}
}

// Create - creates new session for chatID, userID and command
func (session *SessionRepository) Create(chatID int, userID int, command string) {
	key := session.keyFor(chatID, userID)
	session.redis.Set(key, command, 0)
}

// Add - adds user's answer to existing session
func (session *SessionRepository) Add(chatID int, userID int, userAnswer string) {
	key := session.dialogKeyFor(chatID, userID)

	session.redis.RPush(key, userAnswer)
}

// Remove - removes session
func (session *SessionRepository) Remove(chatID int, userID int) {
	key := session.keyFor(chatID, userID)
	session.redis.Del(key)

	key = session.dialogKeyFor(chatID, userID)
	session.redis.Del(key)
}

// Command - returns command for active session for chatID and userID, if exists
// otherwise returns empty string
func (session *SessionRepository) Command(chatID int, userID int) string {
	key := session.keyFor(chatID, userID)
	value, _ := session.redis.Get(key).Result()
	return value
}

// Dialog returns all user's answers history for chatID and userID
func (session *SessionRepository) Dialog(chatID int, userID int) []string {
	key := session.dialogKeyFor(chatID, userID)

	values, _ := session.redis.LRange(key, 0, -1).Result()
	return values
}

func (session *SessionRepository) keyFor(chatID int, userID int) string {
	return fmt.Sprintf("%s_%d_%d", session.key, chatID, userID)
}

func (session *SessionRepository) dialogKeyFor(chatID int, userID int) string {
	return fmt.Sprintf("%s_dialog", session.keyFor(chatID, userID))
}
