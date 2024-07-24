package usecase

import (
	"context"
	"math"
	"redditclone/internal/model"
	"redditclone/internal/repository"
	"redditclone/internal/session"
	"sort"
)

type PostService struct {
	r repository.Post
	c repository.Comment
}

func NewPostService(r repository.Post, c repository.Comment) *PostService {
	return &PostService{
		r: r,
		c: c,
	}
}

type PostSort []model.Post

func (a PostSort) Len() int           { return len(a) }
func (a PostSort) Less(i, j int) bool { return a[i].Score > a[j].Score }
func (a PostSort) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func (p *PostService) GetPost(postID string) (model.Post, error) {
	post, err := p.r.GetPost(postID)
	if err != nil {
		return model.Post{}, err
	}

	comments, err := p.c.GetPostComments(postID)
	if err != nil {
		return model.Post{}, err
	}

	post.Comments = comments

	return post, nil
}

func (p *PostService) GetAllPosts() ([]model.Post, error) {
	posts, err := p.r.GetAllPosts()
	if err != nil {
		return nil, err
	}

	var comments []*model.Comment
	for i, post := range posts {
		comments, err = p.c.GetPostComments(post.ID)
		if err != nil {
			return nil, err
		}

		posts[i].Comments = comments
	}

	sort.Sort(PostSort(posts))
	return posts, err
}

func (p *PostService) GetCategory(category string) ([]model.Post, error) {
	posts, err := p.r.GetCategory(category)
	if err != nil {
		return nil, err
	}

	var comments []*model.Comment
	for i, post := range posts {
		comments, err = p.c.GetPostComments(post.ID)
		if err != nil {
			return nil, err
		}

		posts[i].Comments = comments
	}

	sort.Sort(PostSort(posts))
	return posts, err
}

func (p *PostService) GetUserPosts(userLogin string) ([]model.Post, error) {
	posts, err := p.r.GetUserPosts(userLogin)
	if err != nil {
		return nil, err
	}

	var comments []*model.Comment
	for i, post := range posts {
		comments, err = p.c.GetPostComments(post.ID)
		if err != nil {
			return nil, err
		}

		posts[i].Comments = comments
	}

	sort.Sort(PostSort(posts))
	return posts, err
}

func (p *PostService) UpdateVote(vote int, postID string, ctx context.Context) (model.Post, error) {
	currSession, err := session.SessionFromContext(ctx)
	if err != nil {
		return model.Post{}, err
	}

	currUser := &model.User{}
	currUser.ID = currSession.UserID
	currUser.Login = currSession.UserLogin

	post, err := p.r.Get(postID)
	if err != nil {
		return model.Post{}, err
	}

	newVote := &model.Votes{
		User: currUser.ID,
		Vote: vote,
	}

	delIDx := -1
	isNewVote := true
	for idx, item := range post.Votes {
		if item.User == newVote.User {
			if vote == 0 {
				delIDx = idx
			} else {
				post.Votes[idx] = newVote
				isNewVote = false
			}
			break
		}
	}

	if delIDx != -1 {
		post.Votes = append(post.Votes[:delIDx], post.Votes[delIDx+1:]...)
	} else if isNewVote {
		post.Votes = append(post.Votes, newVote)
	}

	var score, upvotes, countVotes int

	for _, item := range post.Votes {
		score += item.Vote
		if item.Vote == 1 {
			upvotes++
		}
		countVotes++
	}

	post.Score = score
	if countVotes == 0 {
		post.UpvotePercentage = 0
	} else {
		post.UpvotePercentage = int(math.Abs(float64(upvotes) / float64(countVotes) * 100))
	}

	err = p.r.UpdatePost(post)
	if err != nil {
		return model.Post{}, err
	}

	comments, err := p.c.GetPostComments(postID)
	if err != nil {
		return model.Post{}, err
	}

	post.Comments = comments

	return post, nil
}

func (p *PostService) Create(post model.Post, ctx context.Context) (model.Post, error) {
	currSession, err := session.SessionFromContext(ctx)
	if err != nil {
		return model.Post{}, err
	}

	currUser := &model.User{}
	currUser.ID = currSession.UserID
	currUser.Login = currSession.UserLogin

	post.Author = *currUser

	post, err = p.r.Create(post)
	if err != nil {
		return model.Post{}, err
	}

	return p.UpdateVote(1, post.ID, ctx)
}

func (p *PostService) Delete(postID string, ctx context.Context) error {
	currSession, err := session.SessionFromContext(ctx)
	if err != nil {
		return err
	}

	currUser := &model.User{}
	currUser.ID = currSession.UserID
	currUser.Login = currSession.UserLogin

	post, err := p.r.Get(postID)
	if err != nil {
		return err
	}

	comments, err := p.c.GetPostComments(postID)
	if err != nil {
		return err
	}

	post.Comments = comments

	if currUser.ID != post.Author.ID {
		return model.ErrUnauthorized
	}

	err = p.r.Delete(postID)
	if err != nil {
		return err
	}

	err = p.c.DeleteAll(postID)
	if err != nil {
		return err
	}

	return nil
}
