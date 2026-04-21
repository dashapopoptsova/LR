package main

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const jwtSecret = "secret"

type repositoryInterface interface {
	createUser(username, password string) error
	getUser(username string) (int, string, error)
	createPost(userID int, content string) error
	getPosts(userID int) ([]string, error)
}

type service struct {
	repo repositoryInterface
}

func newService(repo repositoryInterface) *service {
	return &service{repo: repo}
}

func (s *service) ping() string {
	return "Hello!"
}

func (s *service) register(username, password string) error {
	return s.repo.createUser(username, password)
}

func (s *service) login(username, password string) (string, error) {
	id, stored, err := s.repo.getUser(username)
	if err != nil {
		return "", errors.New("user not found")
	}

	if stored != password {
		return "", errors.New("wrong password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": id,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	return token.SignedString([]byte(jwtSecret))
}

func (s *service) createPost(userID int, content string) error {
	return s.repo.createPost(userID, content)
}

func (s *service) getPosts(userID int) ([]string, error) {
	return s.repo.getPosts(userID)
}
