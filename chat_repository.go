package margelet

import (
	"gopkg.in/redis.v3"
	"strconv"
)

type ChatRepository struct {
	key   string
	redis *redis.Client
}

var ChatRepo *ChatRepository

func InitChatRepository(prefix string, redis *redis.Client) {
	key := prefix + "margelet_chats"
	ChatRepo = &ChatRepository{key, redis}
}

func (chat *ChatRepository) Add(id int) {
	chat.redis.SAdd(chat.key, strconv.Itoa(id))
}

func (chat *ChatRepository) Remove(id int) {
	chat.redis.SRem(chat.key, strconv.Itoa(id))
}

func (chat *ChatRepository) All() []int {
	var result []int
	strings, _ := chat.redis.SMembers(chat.key).Result()

	for _, str := range strings {
		if c, err := strconv.Atoi(str); err == nil {
			result = append(result, c)
		}
	}
	return result
}
