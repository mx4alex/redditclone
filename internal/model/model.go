package model

type Post struct {
	Author           User       `json:"author"`
	Category         string     `json:"category"`
	Comments         []*Comment `json:"comments"`
	Created          string     `json:"created"`
	ID               string     `json:"id"`
	Score            int        `json:"score"`
	Text             string     `json:"text,omitempty"`
	URL              string     `json:"url,omitempty"`
	Title            string     `json:"title"`
	Type             string     `json:"type"`
	UpvotePercentage int        `json:"upvotePercentage"`
	Views            int        `json:"views"`
	Votes            []*Votes   `json:"votes"`
}

type User struct {
	ID       string `json:"id"`
	Login    string `json:"username"`
	Password string
}

type Comment struct {
	ID      string `json:"id"`
	Author  User   `json:"author"`
	Body    string `json:"body"`
	Created string `json:"created"`
}

type Votes struct {
	User string `json:"user"`
	Vote int    `json:"vote"`
}
