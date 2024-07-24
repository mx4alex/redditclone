package usecase

import (
	"crypto/md5"
	"errors"
	"fmt"
	"redditclone/internal/model"
	"redditclone/internal/repository"
	"redditclone/internal/session"
)

type UserService struct {
	r        repository.User
	sessions *session.SessionsManager
}

func NewUserService(r repository.User, sess *session.SessionsManager) *UserService {
	return &UserService{
		r:        r,
		sessions: sess,
	}
}

func (u *UserService) Authorize(login, pass string) (string, error) {
	passHash := fmt.Sprintf("%x", md5.Sum([]byte(pass)))

	user, err := u.r.Authorize(login, passHash)
	if err != nil {
		return "", err
	}

	sess, err := u.sessions.Create(user)
	if err != nil {
		return "", err
	}

	token, err := u.sessions.CreateToken(sess)
	if err != nil {
		return "", err
	}

	return token, err
}

func (u *UserService) AddUser(login, pass string) (string, error) {
	passHash := fmt.Sprintf("%x", md5.Sum([]byte(pass)))
	_, err := u.r.Authorize(login, passHash)

	if !errors.Is(err, model.ErrNoUser) {
		return "", model.ErrUserExist
	}

	user, err := u.r.AddUser(login, passHash)
	if err != nil {
		return "", err
	}

	sess, err := u.sessions.Create(user)
	if err != nil {
		return "", err
	}

	token, err := u.sessions.CreateToken(sess)
	if err != nil {
		return "", err
	}

	return token, nil
}
