package margelet

import (
	"fmt"
	"gopkg.in/redis.v3"
	"strings"
)

type chatConfigRepository struct {
	key   string
	redis *redis.Client
}

func newChatConfigRepository(prefix string, redis *redis.Client) *chatConfigRepository {
	key := strings.Join([]string{prefix, "margelet_chat_configs"}, "-")
	return &chatConfigRepository{key, redis}
}

func (chatConfig *chatConfigRepository) Set(chatID int, JSON string) {
	chatConfig.redis.Set(chatConfig.ketFor(chatID), JSON, 0)
}

func (chatConfig *chatConfigRepository) Remove(chatID int) {
	chatConfig.redis.Del(chatConfig.ketFor(chatID))
}

func (chatConfig *chatConfigRepository) Get(chatID int) string {
	json, _ := chatConfig.redis.Get(chatConfig.ketFor(chatID)).Result()
	return json
}

func (chatConfig *chatConfigRepository) ketFor(chatID int) string {
	return fmt.Sprintf("%s_%d", chatConfig.key, chatID)
}
