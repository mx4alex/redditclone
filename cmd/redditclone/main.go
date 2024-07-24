package main

import (
	"fmt"
	"net/http"
	"redditclone/internal/handlers"
	"redditclone/internal/middleware"
	"redditclone/internal/repository"
	"redditclone/internal/session"
	"redditclone/internal/usecase"

	"io/ioutil"

	"go.uber.org/zap"
)

func main() {
	sm := session.NewSessionsManager()
	zapLogger, err := zap.NewProduction()
	if err != nil {
		fmt.Println("Error in NewProduction: ", err)
		return
	}

	defer func(zapLogger *zap.Logger) {
		err = zapLogger.Sync()
		if err != nil {
			fmt.Println("Error in Sync: ", err)
		}
	}(zapLogger)

	logger := zapLogger.Sugar()

	repo := repository.NewRepository()
	service := usecase.NewService(repo, sm)
	handler := handlers.NewHandler(logger, service)

	mux := http.NewServeMux()

	mux.HandleFunc("/api/", handler.Handle)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	mux.Handle("/", http.FileServer(http.Dir("./static/html/")))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		var file []byte
		file, err = ioutil.ReadFile("./static/html/index.html")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			logger.Infow("Error in Read", err)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(file)
		if err != nil {
			logger.Infow("Error in Write", err)
			return
		}
	})

	h := middleware.Auth(sm, mux)
	h = middleware.AccessLog(logger, h)
	h = middleware.Panic(h)

	addr := ":8080"
	logger.Infow("starting server",
		"type", "START",
		"addr", addr,
	)

	err = http.ListenAndServe(addr, h)
	if err != nil {
		logger.Infow("Error in ListenAndServe", err)
		return
	}
}
