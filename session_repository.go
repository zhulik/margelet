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
	key := prefix + "margelet_sessions_"
	SessionRepo = &SessionRepository{key, redis}
}

func (session *SessionRepository) Create(chatId int, userId int, command, userAnswer string) {
	key := session.keyFor(chatId, userId)
	session.redis.HSet(key, "command", command)
}

func (session *SessionRepository) Add(chatId int, userId int, userAnswer string) {
	key := session.keyFor(chatId, userId)

	session.redis.RPush(key, userAnswer)
}

func (session *SessionRepository) Remove(chatId int, userId int) {
	key := session.keyFor(chatId, userId)
	session.redis.Del(key)
}

func (session *SessionRepository) Find(chatId int, userId int) []string {
	keys, _ := session.redis.Keys(fmt.Sprint("%s_%d_%d*", session.key, chatId, userId)).Result()
	if len(keys) == 0 {
		return []string{}
	}
	values, _ := session.redis.LRange(keys[0], 0, -1).Result()
	return values
}

func (session *SessionRepository) keyFor(chatId int, userId int) string {
	return fmt.Sprint("%s_%d_%d", session.key, chatId, userId)
}
