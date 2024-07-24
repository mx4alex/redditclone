package usecase

import (
	"context"
	"redditclone/internal/model"
	"redditclone/internal/repository"
	"redditclone/internal/session"
)

type CommentService struct {
	r repository.Comment
	p repository.Post
}

func NewCommentService(r repository.Comment, p repository.Post) *CommentService {
	return &CommentService{
		r: r,
		p: p,
	}
}

func (c *CommentService) Get(commentID string, postID string) (*model.Comment, error) {
	return c.r.Get(commentID, postID)
}

func (c *CommentService) Add(text string, postID string, ctx context.Context) (model.Post, error) {
	currSession, err := session.SessionFromContext(ctx)
	if err != nil {
		return model.Post{}, err
	}

	author := &model.User{}
	author.ID = currSession.UserID
	author.Login = currSession.UserLogin

	err = c.r.Create(text, author, postID)
	if err != nil {
		return model.Post{}, err
	}

	post, err := c.p.Get(postID)
	if err != nil {
		return model.Post{}, err
	}

	comments, err := c.r.GetPostComments(postID)
	if err != nil {
		return model.Post{}, err
	}

	post.Comments = comments

	return post, nil
}

func (c *CommentService) Delete(commentID string, postID string, ctx context.Context) (model.Post, error) {
	currSession, err := session.SessionFromContext(ctx)
	if err != nil {
		return model.Post{}, err
	}

	currUser := &model.User{}
	currUser.ID = currSession.UserID
	currUser.Login = currSession.UserLogin

	comment, err := c.Get(commentID, postID)
	if err != nil {
		return model.Post{}, err
	}

	if currUser.ID != comment.Author.ID {
		return model.Post{}, model.ErrUnauthorized
	}

	err = c.r.Delete(commentID, postID)
	if err != nil {
		return model.Post{}, err
	}

	post, err := c.p.GetPost(postID)
	if err != nil {
		return model.Post{}, err
	}

	comments, err := c.r.GetPostComments(postID)
	if err != nil {
		return model.Post{}, err
	}

	post.Comments = comments

	return post, nil
}

func (c *CommentService) DeleteAll(postID string) {
	c.r.DeleteAll(postID)
}
