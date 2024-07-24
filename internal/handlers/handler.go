package handlers

import (
	"go.uber.org/zap"
	"net/http"
	"redditclone/internal/usecase"
	"strings"
)

type Handler struct {
	Logger  *zap.SugaredLogger
	service *usecase.Service
}

func NewHandler(logger *zap.SugaredLogger, service *usecase.Service) *Handler {
	return &Handler{
		Logger:  logger,
		service: service,
	}
}

func (h *Handler) GetHelper(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) == 4 {
		h.GetPostByID(w, r)
	} else if len(parts) == 5 {
		h.PostRating(w, r)
	}
}

func (h *Handler) DeleteHelper(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) == 4 {
		h.DeletePost(w, r)
	} else if len(parts) == 5 {
		h.DeleteComment(w, r)
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost && r.URL.Path == "/api/register":
		h.Register(w, r)
	case r.Method == http.MethodPost && r.URL.Path == "/api/login":
		h.Login(w, r)
	case r.Method == http.MethodPost && r.URL.Path == "/api/posts":
		h.AddPost(w, r)
	case r.Method == http.MethodPost && strings.HasPrefix(r.URL.Path, "/api/post/"):
		h.AddComment(w, r)
	case r.Method == http.MethodGet && r.URL.Path == "/api/posts/":
		h.GetPosts(w, r)
	case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/api/posts/"):
		h.GetPostsCategory(w, r)
	case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/api/post/"):
		h.GetHelper(w, r)
	case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/api/user/"):
		h.GetUserPosts(w, r)
	case r.Method == http.MethodDelete && strings.HasPrefix(r.URL.Path, "/api/post/"):
		h.DeleteHelper(w, r)
	default:
		http.Error(w, "bad request", http.StatusBadRequest)
	}
}
