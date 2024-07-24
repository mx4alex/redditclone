package repository

import (
	"redditclone/internal/model"
)

type User interface {
	Authorize(login, pass string) (model.User, error)
	AddUser(login, pass string) (model.User, error)
}

type Post interface {
	Get(postID string) (model.Post, error)
	GetPost(postID string) (model.Post, error)
	GetCategory(category string) ([]model.Post, error)
	GetAllPosts() ([]model.Post, error)
	GetUserPosts(userLogin string) ([]model.Post, error)
	UpdatePost(post model.Post) error
	Create(post model.Post) (model.Post, error)
	Delete(postID string) error
}

type Comment interface {
	Get(commentID string, postID string) (*model.Comment, error)
	GetPostComments(postID string) ([]*model.Comment, error)
	Create(text string, author *model.User, postID string) error
	Delete(commentID string, postID string) error
	DeleteAll(postID string) error
}

type Repository struct {
	User
	Post
	Comment
}

func NewRepository() *Repository {
	return &Repository{
		User:    NewUserMemory(),
		Post:    NewPostMemory(),
		Comment: NewCommentMemory(),
	}
}
