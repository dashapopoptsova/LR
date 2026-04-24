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
    CREATE TABLE IF NOT EXISTS posts (
      id SERIAL PRIMARY KEY,
      user_id INT NOT NULL REFERENCES users(id),
      content TEXT NOT NULL
    );
    CREATE TABLE IF NOT EXISTS dbtest_logs (
      id SERIAL PRIMARY KEY,
      value TEXT NOT NULL
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

func (r *repository) getUser(username string) (id int, password string, err error) {
	err = r.db.QueryRow(`SELECT id, password FROM users WHERE username = $1`, username).Scan(&id, &password)
	return
}

func (r *repository) createPost(userID int, content string) error {
	_, err := r.db.Exec(`INSERT INTO posts (user_id, content) VALUES ($1, $2)`, userID, content)
	return err
}

func (r *repository) getPosts(userID int) ([]string, error) {
	rows, err := r.db.Query(`SELECT content FROM posts WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []string
	for rows.Next() {
		var content string
		if err := rows.Scan(&content); err != nil {
			return nil, err
		}
		posts = append(posts, content)
	}
	return posts, nil
}

func (r *repository) saveDBTest(value string) error {
	_, err := r.db.Exec(`INSERT INTO dbtest_logs (value) VALUES ($1)`, value)
	return err
}
