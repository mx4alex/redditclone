package handlers

type PostForm struct {
	Category string `json:"category"`
	Text     string `json:"text"`
	Title    string `json:"title"`
	Type     string `json:"type,omitempty"`
	URL      string `json:"url,omitempty"`
}

type CommentForm struct {
	Comment string `json:"comment"`
}

type LoginForm struct {
	Login    string `json:"username"`
	Password string `json:"password"`
}

type ErrForm struct {
	Location string `json:"location"`
	Param    string `json:"param"`
	Msg      string `json:"msg"`
	Value    string `json:"value,omitempty"`
}
