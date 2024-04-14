package session

import (
	"context"
	"fmt"

	"github.com/spf13/viper"
)

type Session struct {
	// токен
	ParticipantToken string
	// является ли пользователь админом
	IsAdmin bool
}

type sessKey string

var SessionKey sessKey = sessKey(viper.GetString("sessKey"))

// SessionFromContext получает сессию из контекста
func SessionFromContext(ctx context.Context) (*Session, error) {
	sess, ok := ctx.Value(SessionKey).(*Session)
	if !ok || sess == nil {
		return nil, fmt.Errorf("error: no session found")
	}
	return sess, nil
}

// ContextWithSession помещает сессию в контекст
func ContextWithSession(ctx context.Context, sess *Session) context.Context {
	return context.WithValue(ctx, SessionKey, sess)
}
