package margelet

import (
	"gopkg.in/redis.v3"
	"strconv"
	"strings"
)

// ChatRepository - repository for started chats
type ChatRepository struct {
	key   string
	redis *redis.Client
}

func newChatRepository(prefix string, redis *redis.Client) *ChatRepository {
	key := strings.Join([]string{prefix, "margelet_chats"}, "-")
	return &ChatRepository{key, redis}
}

func (chat *ChatRepository) Add(id int64) {
	chat.redis.SAdd(chat.key, strconv.FormatInt(id, 10))
}

func (chat *ChatRepository) Exist(id int64) (res bool) {
	res, _ = chat.redis.SIsMember(chat.key, strconv.FormatInt(id, 10)).Result()
	return
}

func (chat *ChatRepository) Remove(id int64) {
	chat.redis.SRem(chat.key, strconv.FormatInt(id, 10))
}

func (chat *ChatRepository) All() []int64 {
	var result []int64
	strings, _ := chat.redis.SMembers(chat.key).Result()

	for _, str := range strings {
		if c, err := strconv.ParseInt(str, 10, 64); err == nil {
			result = append(result, c)
		}
	}
	return result
}
