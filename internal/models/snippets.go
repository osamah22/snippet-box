package models

import (
	"database/sql"
	"errors"
	"time"

	_ "github.com/lib/pq"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModelInterface interface {
	Insert(title string, content string, expires int) (int, error)
	Get(id int) (*Snippet, error)
	Latest() ([]*Snippet, error)
}

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title, content string, expires int) (int, error) {
	command := `INSERT INTO snippets (title, content, created, expires) 
		VALUES ($1, $2, NOW(), NOW() + ($3 || ' days')::interval)
		RETURNING id`

	var id int
	err := m.DB.QueryRow(command, title, content, expires).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m *SnippetModel) Get(id int) (*Snippet, error) {
	query := `select id, title, content, created, expires FROM snippets
	WHERE expires > CURRENT_TIMESTAMP AND id = $1`

	snippet := Snippet{}
	err := m.DB.QueryRow(query, id).Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Created, &snippet.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}
	return &snippet, nil
}

func (m *SnippetModel) Latest() ([]*Snippet, error) {
	query := `select id, title, content, created, expires FROM snippets 
	where expires > CURRENT_TIMESTAMP LIMIT 10`

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	snippets := make([]*Snippet, 0)

	for rows.Next() {
		s := &Snippet{}

		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
