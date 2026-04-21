package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type repository struct {
	db *sql.DB
}

func newRepository(dsn string) (*repository, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS messages (
			id SERIAL PRIMARY KEY,
			text TEXT
		);
	`)
	if err != nil {
		return nil, err
	}

	return &repository{db: db}, nil
}

func (r *repository) createUser(username, password string) error {
	_, err := r.db.Exec(`INSERT INTO users (username, password) VALUES ($1, $2)`, username, password)
	return err
}

func (r *repository) getUserPassword(username string) (string, error) {
	var password string
	err := r.db.QueryRow(`SELECT password FROM users WHERE username = $1`, username).Scan(&password)
	return password, err
}

func (r *repository) saveMessage(text string) error {
	_, err := r.db.Exec(`INSERT INTO messages (text) VALUES ($1)`, text)
	return err
}

func (r *repository) getMessages() ([]string, error) {
	rows, err := r.db.Query(`SELECT text FROM messages`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []string
	for rows.Next() {
		var text string
		if err := rows.Scan(&text); err != nil {
			return nil, err
		}
		messages = append(messages, text)
	}
	return messages, nil
}
