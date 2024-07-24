package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"redditclone/internal/model"
)

func (h *Handler) GetPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := h.service.Post.GetAllPosts()
	if err != nil {
		h.Logger.Infow("Error in GetAllPosts", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(posts)
	if err != nil {
		h.Logger.Infow("Error marshalings", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(resp)
	if err != nil {
		h.Logger.Infow("Error write response", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetPostsCategory(w http.ResponseWriter, r *http.Request) {
	args := strings.Split(r.URL.Path, "/")
	if len(args) != 4 {
		h.Logger.Infow("Not correct parameters")
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	category := args[3]
	posts, err := h.service.Post.GetCategory(category)
	if err != nil {
		h.Logger.Infow("Error in GetCategory", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(posts)
	if err != nil {
		h.Logger.Infow("Error marshalings", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(resp)
	if err != nil {
		h.Logger.Infow("Error write response", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetPostByID(w http.ResponseWriter, r *http.Request) {
	args := strings.Split(r.URL.Path, "/")
	if len(args) != 4 {
		h.Logger.Infow("Not correct parameters")
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	postID := args[3]
	post, err := h.service.Post.GetPost(postID)
	if err != nil {
		h.Logger.Infow("Error in GetPost", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(post)
	if err != nil {
		h.Logger.Infow("Error marshaling", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(resp)
	if err != nil {
		h.Logger.Infow("Error write response", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) PostRating(w http.ResponseWriter, r *http.Request) {
	args := strings.Split(r.URL.Path, "/")
	n := len(args)
	postID := args[n-2]

	voteType := args[n-1]
	var vote int
	switch voteType {
	case "upvote":
		vote = 1
	case "downvote":
		vote = -1
	default:
		vote = 0
	}

	elem, err := h.service.Post.UpdateVote(vote, postID, r.Context())
	if err != nil {
		h.Logger.Infow("Error in UpdateVote", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, errMarshal := json.Marshal(elem)
	if errMarshal != nil {
		h.Logger.Infow("Error marshaling", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(resp)
	if err != nil {
		h.Logger.Infow("Error write response", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetUserPosts(w http.ResponseWriter, r *http.Request) {
	args := strings.Split(r.URL.Path, "/")
	if len(args) != 4 {
		h.Logger.Infow("Not correct parameters")
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	userLogin := args[3]
	posts, err := h.service.Post.GetUserPosts(userLogin)
	if err != nil {
		h.Logger.Infow("Error in GetUserPosts", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(posts)
	if err != nil {
		h.Logger.Infow("Error marshaling", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(resp)
	if err != nil {
		h.Logger.Infow("Error write response", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) AddPost(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.Logger.Infow("Error reading body", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	post := model.Post{}
	err = json.Unmarshal(body, &post)
	if err != nil {
		h.Logger.Infow("Error unmarshaling", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	post, err = h.service.Post.Create(post, r.Context())
	if err != nil {
		h.Logger.Infow("Error in CreatePost", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(post)
	if err != nil {
		h.Logger.Infow("Error marshaling", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(resp)
	if err != nil {
		h.Logger.Infow("Error write response", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusCreated)
}

func (h *Handler) DeletePost(w http.ResponseWriter, r *http.Request) {
	args := strings.Split(r.URL.Path, "/")
	if len(args) != 4 {
		h.Logger.Infow("Not correct parameters")
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	postID := args[3]

	err := h.service.Post.Delete(postID, r.Context())
	if err != nil {
		if errors.Is(err, model.ErrUnauthorized) {
			h.Logger.Infow("Unauthorized", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		h.Logger.Infow("Error in Delete post", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(map[string]interface{}{
		"message": "success",
	})

	if err != nil {
		h.Logger.Infow("Error marshaling", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(resp)
	if err != nil {
		h.Logger.Infow("Error write response", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
