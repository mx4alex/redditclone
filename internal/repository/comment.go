package repository

import (
	"sync"
	"time"

	"redditclone/internal/model"

	"github.com/google/uuid"
)

type CommentMemory struct {
	data map[string][]*model.Comment
	mu   sync.Mutex
}

func NewCommentMemory() *CommentMemory {
	return &CommentMemory{
		data: make(map[string][]*model.Comment),
		mu:   sync.Mutex{},
	}
}

func (c *CommentMemory) Get(commentID string, postID string) (*model.Comment, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, comment := range c.data[postID] {
		if comment.ID == commentID {
			return comment, nil
		}
	}

	return nil, model.ErrNoComment
}

func (c *CommentMemory) GetPostComments(postID string) ([]*model.Comment, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	comments := make([]*model.Comment, 0)
	for _, comment := range c.data[postID] {
		comments = append(comments, comment)
	}

	return comments, nil
}

func (c *CommentMemory) Create(text string, author *model.User, postID string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	comment := new(model.Comment)
	comment.ID = uuid.New().String()
	comment.Author = *author
	comment.Created = time.Now().Format(time.RFC3339)
	comment.Body = text

	if _, ok := c.data[postID]; ok {
		c.data[postID] = append(c.data[postID], comment)
	} else {
		c.data[postID] = make([]*model.Comment, 1, 10)
		c.data[postID][0] = comment
	}

	return nil
}

func (c *CommentMemory) Delete(commentID string, postID string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	for i, comment := range c.data[postID] {
		if comment.ID == commentID {
			c.data[postID] = append(c.data[postID][:i], c.data[postID][i+1:]...)
			return nil
		}
	}

	return model.ErrNoDel
}

func (c *CommentMemory) DeleteAll(postID string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, postID)

	return nil
}
