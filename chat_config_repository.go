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
func (chatConfig *ChatConfigRepository) Set(chatID int, JSON string) {
	chatConfig.redis.Set(chatConfig.ketFor(chatID), JSON, 0)
}

// SetWithStruct - stores any config for chatID using a struct
func (chatConfig *ChatConfigRepository) SetWithStruct(chatID int, obj interface{}) {
	valueBytes, _ := json.Marshal(obj)
	valueString := string(valueBytes)
	chatConfig.Set(chatID, valueString)
}

// Remove - removes config for chatID
func (chatConfig *ChatConfigRepository) Remove(chatID int) {
	chatConfig.redis.Del(chatConfig.ketFor(chatID))
}

// Get - returns config for chatID
func (chatConfig *ChatConfigRepository) Get(chatID int) string {
	json, _ := chatConfig.redis.Get(chatConfig.ketFor(chatID)).Result()
	return json
}

// GetWithStruct - returns config for chatID using a struct
func (chatConfig *ChatConfigRepository) GetWithStruct(chatID int, obj interface{}) {
	valueString := chatConfig.Get(chatID)
	json.Unmarshal([]byte(valueString), obj)
}

func (chatConfig *ChatConfigRepository) ketFor(chatID int) string {
	return fmt.Sprintf("%s_%d", chatConfig.key, chatID)
}
