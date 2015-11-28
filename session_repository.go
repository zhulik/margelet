package margelet

import (
	"fmt"
	"gopkg.in/redis.v3"
	"strings"
)

type sessionRepository struct {
	key   string
	redis *redis.Client
}

func newSessionRepository(prefix string, redis *redis.Client) *sessionRepository {
	key := strings.Join([]string{prefix, "margelet_sessions"}, "-")
	return &sessionRepository{key, redis}
}

func (session *sessionRepository) Create(chatID int, userID int, command string) {
	key := session.keyFor(chatID, userID)
	session.redis.Set(key, command, 0)
}

func (session *sessionRepository) Add(chatID int, userID int, userAnswer string) {
	key := session.dialogKeyFor(chatID, userID)

	session.redis.RPush(key, userAnswer)
}

func (session *sessionRepository) Remove(chatID int, userID int) {
	key := session.keyFor(chatID, userID)
	session.redis.Del(key)

	key = session.dialogKeyFor(chatID, userID)
	session.redis.Del(key)
}

func (session *sessionRepository) Command(chatID int, userID int) string {
	key := session.keyFor(chatID, userID)
	value, _ := session.redis.Get(key).Result()
	return value
}

func (session *sessionRepository) Dialog(chatID int, userID int) []string {
	key := session.dialogKeyFor(chatID, userID)

	values, _ := session.redis.LRange(key, 0, -1).Result()
	return values
}

func (session *sessionRepository) keyFor(chatID int, userID int) string {
	return fmt.Sprintf("%s_%d_%d", session.key, chatID, userID)
}

func (session *sessionRepository) dialogKeyFor(chatID int, userID int) string {
	return fmt.Sprintf("%s_dialog", session.keyFor(chatID, userID))
}
