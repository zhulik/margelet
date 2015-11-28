package margelet

import (
	"fmt"
	"gopkg.in/redis.v3"
)

type SessionRepository struct {
	key   string
	redis *redis.Client
}

var SessionRepo *SessionRepository

func InitSessionRepository(prefix string, redis *redis.Client) {
	key := prefix + "margelet_sessions"
	SessionRepo = &SessionRepository{key, redis}
}

func (session *SessionRepository) Create(chatId int, userId int, command string) {
	key := session.keyFor(chatId, userId)
	session.redis.Set(key, command, 0)
}

func (session *SessionRepository) Add(chatId int, userId int, userAnswer string) {
	key := session.dialogKeyFor(chatId, userId)

	session.redis.RPush(key, userAnswer)
}

func (session *SessionRepository) Remove(chatId int, userId int) {
	key := session.keyFor(chatId, userId)
	session.redis.Del(key)

	key = session.dialogKeyFor(chatId, userId)
	session.redis.Del(key)
}

func (session *SessionRepository) Command(chatId int, userId int) string {
	key := session.keyFor(chatId, userId)
	value, _ := session.redis.Get(key).Result()
	return value
}

func (session *SessionRepository) Dialog(chatId int, userId int) []string {
	key := session.dialogKeyFor(chatId, userId)

	values, _ := session.redis.LRange(key, 0, -1).Result()
	return values
}

func (session *SessionRepository) keyFor(chatId int, userId int) string {
	return fmt.Sprintf("%s_%d_%d", session.key, chatId, userId)
}

func (session *SessionRepository) dialogKeyFor(chatId int, userId int) string {
	return fmt.Sprintf("%s_dialog", session.keyFor(chatId, userId))
}
