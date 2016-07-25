package margelet

import (
	"fmt"
	"gopkg.in/redis.v3"
	"strconv"
	"strings"
)

// StatsRepository - public interface for session repository
type StatsRepository interface {
	Incr(chatID int64, userID int, name string)
	Get(chatID int64, userID int, name string) int
}

type statsRepository struct {
	key   string
	redis *redis.Client
}

func newStatsRepository(prefix string, redis *redis.Client) *statsRepository {
	key := strings.Join([]string{prefix, "margelet_sessions"}, "-")
	return &statsRepository{key, redis}
}

// Inc - adds user's answer to existing session
func (stats *statsRepository) Incr(chatID int64, userID int, name string) {
	key := stats.keyFor(chatID, userID, name)
	stats.redis.Incr(key)
}

// Get - adds user's answer to existing session
func (stats *statsRepository) Get(chatID int64, userID int, name string) int {
	key := stats.keyFor(chatID, userID, name)
	value, _ := stats.redis.Get(key).Result()
	v, _ := strconv.Atoi(value)
	return v
}

func (stats *statsRepository) keyFor(chatID int64, userID int, name string) string {
	return fmt.Sprintf("%s_%d_%d_%d", stats.key, chatID, userID, name)
}
