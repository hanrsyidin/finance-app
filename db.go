package main

import (
	"database/sql"
	"log"

	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

var db *sql.DB

func InitDB(dataSourceName string) {
	var err error
	db, err = sql.Open("sqlite", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	createTables()
	seedDefaultData()
}

func createTables() {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		type TEXT NOT NULL, -- 'income' or 'expense'
		color TEXT, -- hex code or tailwind class
		icon TEXT
	);

	CREATE TABLE IF NOT EXISTS transactions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		amount REAL NOT NULL,
		description TEXT,
		date TEXT NOT NULL, -- YYYY-MM-DD
		type TEXT NOT NULL, -- 'income' or 'expense'
		category_id INTEGER,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (category_id) REFERENCES categories(id)
	);
	`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("Error creating tables: ", err)
	}
}

func seedDefaultData() {
	// Check if admin user exists
	var count int
	row := db.QueryRow("SELECT COUNT(*) FROM users WHERE username = 'admin'")
	err := row.Scan(&count)
	if err != nil {
		log.Println("Error checking users:", err)
		return
	}

	if count == 0 {
		// Create default admin user (password: admin123)
		hash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		_, err = db.Exec("INSERT INTO users (username, password_hash) VALUES (?, ?)", "admin", string(hash))
		if err != nil {
			log.Println("Error seeding admin user:", err)
		} else {
			log.Println("✅ Default user created: admin / admin123")
		}
	}

	// Check if categories exist
	row = db.QueryRow("SELECT COUNT(*) FROM categories")
	row.Scan(&count)
	if count == 0 {
		// Seed some default categories
		cats := []struct {
			Name  string
			Type  string
			Color string
			Icon  string
		}{
			{"Salary", "income", "bg-emerald-500", "💰"},
			{"Freelance", "income", "bg-blue-500", "💻"},
			{"Food", "expense", "bg-orange-500", "🍔"},
			{"Transport", "expense", "bg-indigo-500", "🚌"},
			{"Utilities", "expense", "bg-yellow-500", "⚡"},
			{"Entertainment", "expense", "bg-pink-500", "🎬"},
		}

		stmt, _ := db.Prepare("INSERT INTO categories (name, type, color, icon) VALUES (?, ?, ?, ?)")
		for _, c := range cats {
			stmt.Exec(c.Name, c.Type, c.Color, c.Icon)
		}
		log.Println("✅ Default categories seeded")
	}
}
