package main

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const jwtSecret = "secret"

type repositoryInterface interface {
	createUser(username, password string) error
	getUserPassword(username string) (string, error)
	saveMessage(text string) error
	getMessages() ([]string, error)
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
	stored, err := s.repo.getUserPassword(username)
	if err != nil {
		return "", errors.New("user not found")
	}

	if stored != password {
		return "", errors.New("wrong password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})

	return token.SignedString([]byte(jwtSecret))
}

func (s *service) saveMessage(text string) error {
	return s.repo.saveMessage(text)
}

func (s *service) getMessages() ([]string, error) {
	return s.repo.getMessages()
}
