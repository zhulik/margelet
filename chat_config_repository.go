package margelet

import (
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/redis.v3"
)

// ChatConfigRepository - repository for chat configs
type ChatConfigRepository struct {
	key   string
	redis *redis.Client
}

func newChatConfigRepository(prefix string, redis *redis.Client) *ChatConfigRepository {
	key := strings.Join([]string{prefix, "margelet_chat_configs"}, "-")
	return &ChatConfigRepository{key, redis}
}

// Set - stores any config for chatID
func (chatConfig *ChatConfigRepository) Set(chatID int64, JSON string) {
	chatConfig.redis.Set(chatConfig.keyFor(chatID), JSON, 0)
}

// SetWithStruct - stores any config for chatID using a struct
func (chatConfig *ChatConfigRepository) SetWithStruct(chatID int64, obj interface{}) {
	valueBytes, _ := json.Marshal(obj)
	valueString := string(valueBytes)
	chatConfig.Set(chatID, valueString)
}

// Remove - removes config for chatID
func (chatConfig *ChatConfigRepository) Remove(chatID int64) {
	chatConfig.redis.Del(chatConfig.keyFor(chatID))
}

// Get - returns config for chatID
func (chatConfig *ChatConfigRepository) Get(chatID int64) string {
	json, _ := chatConfig.redis.Get(chatConfig.keyFor(chatID)).Result()
	return json
}

// GetWithStruct - returns config for chatID using a struct
func (chatConfig *ChatConfigRepository) GetWithStruct(chatID int64, obj interface{}) {
	valueString := chatConfig.Get(chatID)
	json.Unmarshal([]byte(valueString), obj)
}

func (chatConfig *ChatConfigRepository) keyFor(chatID int64) string {
	return fmt.Sprintf("%s_%d", chatConfig.key, chatID)
}
