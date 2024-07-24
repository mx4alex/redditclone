package model

import "errors"

var (
	ErrNoPost       = errors.New("no post found")
	ErrNoDel        = errors.New("there is no post being deleted")
	ErrNoDelComm    = errors.New("there is no comment being deleted")
	ErrNoUser       = errors.New("no user found")
	ErrBadPass      = errors.New("invalid password")
	ErrNoComment    = errors.New("no comment found")
	ErrBadVote      = errors.New("bad vote number for vote")
	ErrNoAuth       = errors.New("no session found")
	ErrUnauthorized = errors.New("unauthorized")
	ErrUserExist    = errors.New("user already exists")
	ErrSignMethod   = errors.New("bad sign method")
)
