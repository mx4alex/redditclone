package usecase

import (
	"context"
	"redditclone/internal/model"
	"redditclone/internal/repository"
	"redditclone/internal/session"
)

type User interface {
	Authorize(login, pass string) (string, error)
	AddUser(login, pass string) (string, error)
}

type Post interface {
	GetPost(postID string) (model.Post, error)
	GetCategory(category string) ([]model.Post, error)
	GetAllPosts() ([]model.Post, error)
	GetUserPosts(userLogin string) ([]model.Post, error)
	UpdateVote(vote int, postID string, ctx context.Context) (model.Post, error)
	Create(post model.Post, ctx context.Context) (model.Post, error)
	Delete(postID string, ctx context.Context) error
}

type Comment interface {
	Get(commentID string, postID string) (*model.Comment, error)
	Add(text string, postID string, ctx context.Context) (model.Post, error)
	Delete(commentID string, postID string, ctx context.Context) (model.Post, error)
	DeleteAll(postID string)
}

type Service struct {
	User
	Post
	Comment
}

func NewService(r *repository.Repository, sess *session.SessionsManager) *Service {
	return &Service{
		User:    NewUserService(r.User, sess),
		Post:    NewPostService(r.Post, r.Comment),
		Comment: NewCommentService(r.Comment, r.Post),
	}
}
