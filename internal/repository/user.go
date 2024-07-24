package repository

import (
	"sync"

	"github.com/google/uuid"
	"redditclone/internal/model"
)

type UserMemory struct {
	data map[string]model.User
	mu   sync.Mutex
}

func NewUserMemory() *UserMemory {
	return &UserMemory{
		data: make(map[string]model.User),
		mu:   sync.Mutex{},
	}
}

func (repo *UserMemory) Authorize(login, pass string) (model.User, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	user, ok := repo.data[login]
	if !ok {
		return model.User{}, model.ErrNoUser
	}

	if user.Password != pass {
		return user, model.ErrBadPass
	}

	return user, nil
}

func (repo *UserMemory) AddUser(login, pass string) (model.User, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	repo.data[login] = model.User{
		ID:       uuid.New().String(),
		Login:    login,
		Password: pass,
	}

	return repo.data[login], nil
}
