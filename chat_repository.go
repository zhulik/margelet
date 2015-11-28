package margelet

import (
	"gopkg.in/redis.v3"
	"strconv"
)

type chatRepository struct {
	key   string
	redis *redis.Client
}

func newChatRepository(prefix string, redis *redis.Client) *chatRepository {
	key := prefix + "margelet_chats"
	return &chatRepository{key, redis}
}

func (chat *chatRepository) Add(id int) {
	chat.redis.SAdd(chat.key, strconv.Itoa(id))
}

func (chat *chatRepository) Remove(id int) {
	chat.redis.SRem(chat.key, strconv.Itoa(id))
}

func (chat *chatRepository) All() []int {
	var result []int
	strings, _ := chat.redis.SMembers(chat.key).Result()

	for _, str := range strings {
		if c, err := strconv.Atoi(str); err == nil {
			result = append(result, c)
		}
	}
	return result
}
