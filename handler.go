package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type serviceInterface interface {
	ping() string
	register(username, password string) error
	login(username, password string) (string, error)
	createPost(userID int, content string) error
	getPosts(userID int) ([]string, error)
}

type handler struct {
	svc serviceInterface
}

func newHandler(svc serviceInterface) *handler {
	return &handler{svc: svc}
}

func (h *handler) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/test", h.test)
	mux.HandleFunc("/register", h.register)
	mux.HandleFunc("/login", h.login)
	mux.HandleFunc("/posts", authMiddleware(h.posts))
	return mux
}

func (h *handler) test(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(h.svc.ping()))
}

func (h *handler) register(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if err := h.svc.register(body.Username, body.Password); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("user created"))
}

func (h *handler) login(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	token, err := h.svc.login(body.Username, body.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Write([]byte(token))
}

func (h *handler) posts(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey).(int)

	switch r.Method {
	case http.MethodPost:
		var body struct {
			Content string `json:"content"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		if err := h.svc.createPost(userID, body.Content); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte("post created"))

	case http.MethodGet:
		posts, err := h.svc.getPosts(userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte(strings.Join(posts, "\n")))
	}
}
