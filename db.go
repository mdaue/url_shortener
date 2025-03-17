package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

type URL struct {
	Name          string    `json:"name"`
	CreatedAt     time.Time `json:"created_at"`
	Short         string    `json:"short"`
	RequestedFrom string    `json:"requested_from"`
	Clicks        int       `json:"clicks"`
}

func openDatabase() (db *sql.DB, err error) {
	db, err = sql.Open("sqlite", config.DatabasePath)
	if err == nil {
		fmt.Println("Opened Database")

		// Create tables if they don't exist
		_, err = db.Exec(`
			CREATE TABLE IF NOT EXISTS urls (
				name TEXT NOT NULL,
				created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
				short TEXT NOT NULL UNIQUE,
				requested_from TEXT NOT NULL,
				clicks INTEGER DEFAULT 0
			)
		`)
		if err == nil {
			fmt.Println("Database schema ready")
		}
	}
	return
}

func createURL(db *sql.DB, name string, short string, requestedFrom string) (URL, error) {
	_, err := db.Exec("INSERT INTO urls (name, short, requested_from) VALUES (?, ?, ?)", name, short, requestedFrom)
	if err != nil {
		return URL{}, err
	}
	var url URL
	row := db.QueryRow("SELECT name, created_at, short, requested_from FROM urls WHERE name = ? AND short = ?", name, short)
	err = row.Scan(&url.Name, &url.CreatedAt, &url.Short, &url.RequestedFrom)
	if err != nil {
		return URL{}, err
	}
	return url, nil
}

func addClicks(db *sql.DB, short string) error {
	_, err := db.Exec("UPDATE urls SET clicks = clicks + 1 WHERE short = ?", short)
	return err
}

func queryURLs(db *sql.DB) ([]URL, error) {
	rows, err := db.Query("SELECT * FROM urls")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var urls []URL
	for rows.Next() {
		var url URL
		err = rows.Scan(&url.Name, &url.CreatedAt, &url.Short, &url.RequestedFrom, &url.Clicks)
		if err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}
	return urls, nil
}

func queryURLsFromRequested(db *sql.DB, requestedFrom string) ([]URL, error) {
	rows, err := db.Query("SELECT * FROM urls WHERE requested_from = ?", requestedFrom)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var urls []URL
	for rows.Next() {
		var url URL
		err = rows.Scan(&url.Name, &url.CreatedAt, &url.Short, &url.RequestedFrom, &url.Clicks)
		if err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}
	return urls, nil
}

func queryRecentURLs(db *sql.DB, limit int) ([]URL, error) {
	rows, err := db.Query("SELECT * FROM urls ORDER BY created_at DESC LIMIT ?", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var urls []URL
	for rows.Next() {
		var url URL
		err = rows.Scan(&url.Name, &url.CreatedAt, &url.Short, &url.RequestedFrom, &url.Clicks)
		if err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}
	return urls, nil
}

func queryShortURL(db *sql.DB, short string) (URL, error) {
	var url URL
	row := db.QueryRow("SELECT * FROM urls WHERE short = ?", short)
	err := row.Scan(&url.Name, &url.CreatedAt, &url.Short, &url.RequestedFrom, &url.Clicks)
	if err != nil {
		return URL{}, err
	}
	return url, nil
}
