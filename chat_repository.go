package margelet

import (
	"gopkg.in/redis.v3"
	"strconv"
	"strings"
)

type chatRepository struct {
	key   string
	redis *redis.Client
}

func newChatRepository(prefix string, redis *redis.Client) *chatRepository {
	key := strings.Join([]string{prefix, "margelet_chats"}, "-")
	return &chatRepository{key, redis}
}

func (chat *chatRepository) Add(id int64) {
	chat.redis.SAdd(chat.key, strconv.FormatInt(id, 10))
}

func (chat *chatRepository) Remove(id int64) {
	chat.redis.SRem(chat.key, strconv.FormatInt(id, 10))
}

func (chat *chatRepository) All() []int64 {
	var result []int64
	strings, _ := chat.redis.SMembers(chat.key).Result()

	for _, str := range strings {
		if c, err := strconv.ParseInt(str, 10, 64); err == nil {
			result = append(result, c)
		}
	}
	return result
}
