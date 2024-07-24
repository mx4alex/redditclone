package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"redditclone/internal/model"
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.Logger.Infow("Error reading body", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	logForm := &LoginForm{}
	err = json.Unmarshal(body, logForm)
	if err != nil {
		h.Logger.Infow("Error unmarshaling", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token, err := h.service.User.Authorize(logForm.Login, logForm.Password)
	if err != nil {
		h.Logger.Infow("Unauthorized", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	resp, err := json.Marshal(map[string]interface{}{
		"token": token,
	})

	if err != nil {
		h.Logger.Infow("Error marshaling", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	_, errWrite := w.Write(resp)
	if errWrite != nil {
		h.Logger.Infow("Error write response", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.Logger.Infow("Error reading body", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()

	logForm := &LoginForm{}
	err = json.Unmarshal(body, logForm)
	if err != nil {
		h.Logger.Infow("Error unmarshaling", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var respErr []byte

	token, err := h.service.User.AddUser(logForm.Login, logForm.Password)
	if errors.Is(err, model.ErrUserExist) {
		respErr, err = json.Marshal(map[string][]ErrForm{
			"errors": {
				{
					Location: "body",
					Msg:      "already exists",
					Param:    "username",
					Value:    logForm.Login,
				},
			}})

		if err != nil {
			h.Logger.Infow("Error marshaling", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		h.Logger.Infow("Unable to process the instructions", err)
		http.Error(w, "", http.StatusUnprocessableEntity)

		_, errWrite := w.Write(respErr)
		if errWrite != nil {
			h.Logger.Infow("Error write response", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		return
	}

	if err != nil {
		h.Logger.Infow("Error in AddUser", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(map[string]interface{}{
		"token": token,
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
