package cmd

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB


func InitializeDatabase() {
	var err error
	db, err = sql.Open("sqlite3", "messages.db")
	if err != nil {
		log.Fatal(err)
	}

	createTable := `
	CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT,
		message TEXT,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatal("Failed to create messages table:", err)
	}
}


func SaveMessage(username, message string) {
	_, err := db.Exec("INSERT INTO messages (username, message) VALUES (?, ?)", username, message)
	if err != nil {
		log.Println("Failed to save message:", err)
	}
}

func GetLastMessages(limit int) []string {
	rows, err := db.Query("SELECT username, message FROM messages ORDER BY timestamp DESC LIMIT ?", limit)
	if err != nil {
		log.Println("Failed to fetch messages:", err)
		return nil
	}
	defer rows.Close()

	var messages []string
	for rows.Next() {
		var username, message string
		if err := rows.Scan(&username, &message); err == nil {
			messages = append(messages, username+": "+message)
		}
	}
	return messages
}
