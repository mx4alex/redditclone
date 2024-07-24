package repository

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"redditclone/internal/model"
)

type PostMemory struct {
	data []model.Post
	mu   sync.Mutex
}

func NewPostMemory() *PostMemory {
	return &PostMemory{
		data: make([]model.Post, 0),
		mu:   sync.Mutex{},
	}
}

func (p *PostMemory) Get(postID string) (model.Post, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, post := range p.data {
		if post.ID == postID {
			return post, nil
		}
	}

	return model.Post{}, model.ErrNoPost
}

func (p *PostMemory) GetPost(postID string) (model.Post, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for idx, post := range p.data {
		if post.ID == postID {
			post.Views += 1
			p.data[idx] = post
			return post, nil
		}
	}

	return model.Post{}, model.ErrNoPost
}

func (p *PostMemory) GetCategory(category string) ([]model.Post, error) {
	suitablePosts := make([]model.Post, 0)
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, post := range p.data {
		if post.Category == category {
			suitablePosts = append(suitablePosts, post)
		}
	}

	return suitablePosts, nil
}

func (p *PostMemory) GetAllPosts() ([]model.Post, error) {
	return p.data, nil
}

func (p *PostMemory) GetUserPosts(userLogin string) ([]model.Post, error) {
	suitablePosts := make([]model.Post, 0)
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, post := range p.data {
		if post.Author.Login == userLogin {
			suitablePosts = append(suitablePosts, post)
		}
	}

	return suitablePosts, nil
}

func (p *PostMemory) Create(post model.Post) (model.Post, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	post.Created = time.Now().Format(time.RFC3339)
	post.UpvotePercentage = 100
	post.Views = 0
	post.Score = 0
	post.Votes = make([]*model.Votes, 0)
	post.ID = uuid.New().String()
	p.data = append(p.data, post)

	return post, nil
}

func (p *PostMemory) UpdatePost(post model.Post) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	for idx, item := range p.data {
		if item.ID == post.ID {
			p.data[idx] = post
		}
	}

	return nil
}

func (p *PostMemory) Delete(postID string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	for i, post := range p.data {
		if post.ID == postID {
			p.data = append(p.data[:i], p.data[i+1:]...)
			return nil
		}
	}

	return model.ErrNoDel
}
