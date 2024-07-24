package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

func (h *Handler) AddComment(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.Logger.Infow("Error reading body", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	args := strings.Split(r.URL.Path, "/")
	if len(args) != 4 {
		h.Logger.Infow("Not correct parameters")
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	postID := args[3]

	commentForm := &CommentForm{}
	err = json.Unmarshal(body, commentForm)
	if err != nil {
		h.Logger.Infow("Error unmarshaling", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	post, err := h.service.Comment.Add(commentForm.Comment, postID, r.Context())
	if err != nil {
		h.Logger.Infow("Error in Create comment", err)
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

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	args := strings.Split(r.URL.Path, "/")
	n := len(args)

	postID := args[n-2]
	commentID := args[n-1]

	post, err := h.service.Comment.Delete(commentID, postID, r.Context())
	if err != nil {
		h.Logger.Infow("Error in Delete comment", err)
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
