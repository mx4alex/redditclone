package session

import (
	jwt "github.com/dgrijalva/jwt-go"
	"net/http"
	"redditclone/internal/model"
	"strings"
	"sync"
	"time"
)

type SessionsManager struct {
	data map[string]*Session
	mu   *sync.RWMutex
}

func NewSessionsManager() *SessionsManager {
	return &SessionsManager{
		data: make(map[string]*Session),
		mu:   &sync.RWMutex{},
	}
}

func (s *SessionsManager) Check(w http.ResponseWriter, r *http.Request) (*Session, error) {
	hashSecretGetter := func(token *jwt.Token) (interface{}, error) {
		method, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok || method.Alg() != "HS256" {
			return nil, model.ErrSignMethod
		}
		return []byte("salt"), nil
	}

	inToken := r.Header.Get("authorization")
	if inToken == "" {
		return nil, model.ErrNoAuth
	}

	inToken = strings.Split(inToken, " ")[1]
	token, errJwt := jwt.Parse(inToken, hashSecretGetter)
	if errJwt != nil {
		return nil, errJwt
	}
	payload, ok := token.Claims.(jwt.MapClaims)

	if ok {
		sessClaims := payload["user"].(map[string]interface{})

		sess := &Session{}
		sess.UserID = sessClaims["id"].(string)
		sess.UserLogin = sessClaims["username"].(string)

		return sess, nil
	}

	return nil, model.ErrNoAuth
}

func (s *SessionsManager) CreateToken(sess *Session) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": map[string]interface{}{
			"username": sess.UserLogin,
			"id":       sess.UserID,
		},
		"iat": time.Now().Unix(),
		"exp": time.Now().Unix() + 1200,
	})
	tokenString, err := token.SignedString([]byte("salt"))

	return tokenString, err
}

func (s *SessionsManager) Create(user model.User) (*Session, error) {
	sess, err := NewSession(user)
	if err != nil {
		return nil, err
	}

	s.mu.Lock()
	s.data[sess.ID] = sess
	s.mu.Unlock()

	return sess, nil
}
