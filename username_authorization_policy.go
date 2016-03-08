package margelet

import (
	"fmt"
	"gopkg.in/telegram-bot-api.v2"
)

// UsernameAuthorizationPolicy - simple authorization policy, that checks sender's username
type UsernameAuthorizationPolicy struct {
	Usernames []string
}

// Allow check message author's username and returns nil if it in Usernames
// otherwise, returns an authorization error message
func (p UsernameAuthorizationPolicy) Allow(message tgbotapi.Message) error {
	for _, username := range p.Usernames {
		if message.From.UserName == username {
			return nil
		}
	}

	return fmt.Errorf("user %s is not allowed to do it", message.From.UserName)
}
