package models

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"
)

var (
	ErrNoRecord = errors.New("models: no matching record found")
)

type URLModelInterface interface {
	Get(short string) (HURL, error)
	Insert(short, url, creator string) (HURL, error)
}

type URLModel struct {
	DB *sql.DB
}

type HURL struct {
	ID      int       `json:"id"`
	Short   string    `json:"short"`
	URL     string    `json:"url"`
	Creator string    `json:"creator"`
	Created time.Time `json:"created"`
	Expires time.Time `json:"expires"`
}

func (u *URLModel) Get(short string) (HURL, error) {
	query := `SELECT id, short, uri, creator, created, expires FROM urls WHERE expires > current_timestamp AND short = ?`

	row := u.DB.QueryRow(query, short)
	var hurl HURL
	if err := row.Scan(&hurl.ID, &hurl.Short, &hurl.URL, &hurl.Creator, &hurl.Created, &hurl.Expires); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return HURL{}, ErrNoRecord
		} else {
			return HURL{}, err
		}
	}
	return hurl, nil
}

func (u *URLModel) Insert(short, url, creator string) (HURL, error) {

	blob, err := json.Marshal(creator)
	if err != nil {
		return HURL{}, nil
	}

	query := `INSERT INTO urls (short, uri, creator) VALUES (?, ?, ?)`

	if _, err := u.DB.Exec(query, short, url, blob); err != nil {
		return HURL{}, err
	}

	return u.Get(short)
}
