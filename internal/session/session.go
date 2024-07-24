package session

import (
	"context"
	"crypto/rand"
	"fmt"
	"redditclone/internal/model"
)

type sessKey string

var SessionKey sessKey = "sessionKey"

type Session struct {
	ID        string
	UserID    string
	UserLogin string
}

func NewSession(user model.User) (*Session, error) {
	randID := make([]byte, 16)
	_, err := rand.Read(randID)
	if err != nil {
		return nil, err
	}

	sess := &Session{
		ID:        fmt.Sprintf("%x", randID),
		UserID:    user.ID,
		UserLogin: user.Login,
	}

	return sess, nil
}

func SessionFromContext(ctx context.Context) (*Session, error) {
	sess, ok := ctx.Value(SessionKey).(*Session)
	if !ok || sess == nil {
		return nil, model.ErrNoAuth
	}
	return sess, nil
}

func ContextWithSession(ctx context.Context, sess *Session) context.Context {
	return context.WithValue(ctx, SessionKey, sess)
}
