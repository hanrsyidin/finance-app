package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// -- Models --

type Category struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	Type             string `json:"type"`
	Color            string `json:"color"`
	Icon             string `json:"icon"`
	TransactionCount int    `json:"transaction_count"`
}

type Transaction struct {
	ID            int     `json:"id"`
	Amount        float64 `json:"amount"`
	Description   string  `json:"note"`
	Date          string  `json:"date"`
	Type          string  `json:"type"`
	CategoryID    int     `json:"category_id"`
	CategoryName  string  `json:"category_name,omitempty"`
	CategoryColor string  `json:"category_color,omitempty"`
	CategoryIcon  string  `json:"category_icon,omitempty"`
}

// -- Handlers --

// Categories

func getCategories(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT c.id, c.name, c.type, c.color, c.icon, COUNT(t.id) as transaction_count
		FROM categories c
		LEFT JOIN transactions t ON c.id = t.category_id
		GROUP BY c.id
		ORDER BY c.name
	`
	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Type, &c.Color, &c.Icon, &c.TransactionCount); err != nil {
			continue
		}
		categories = append(categories, c)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

func createCategory(w http.ResponseWriter, r *http.Request) {
	var c Category
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := db.Exec("INSERT INTO categories (name, type, color, icon) VALUES (?, ?, ?, ?)", c.Name, c.Type, c.Color, c.Icon)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := res.LastInsertId()
	c.ID = int(id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(c)
}

func deleteCategory(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/api/categories/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("DELETE FROM categories WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Transactions

func getTransactions(w http.ResponseWriter, r *http.Request) {
	// Filters
	month := r.URL.Query().Get("month") // YYYY-MM
	limit := r.URL.Query().Get("limit")

	query := `
		SELECT t.id, t.amount, t.description, t.date, t.type, t.category_id,
		       c.name, c.color, c.icon
		FROM transactions t
		LEFT JOIN categories c ON t.category_id = c.id
		WHERE 1=1
	`
	var args []interface{}

	if month != "" {
		query += " AND strftime('%Y-%m', t.date) = ?"
		args = append(args, month)
	}

	query += " ORDER BY t.date DESC, t.id DESC"

	if limit != "" {
		query += " LIMIT ?"
		args = append(args, limit)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var t Transaction
		// Handle potential NULLs if needed, but schema says NOT NULL for most
		var catName, catColor, catIcon *string
		err := rows.Scan(&t.ID, &t.Amount, &t.Description, &t.Date, &t.Type, &t.CategoryID, &catName, &catColor, &catIcon)
		if err != nil {
			continue // or log
		}
		if catName != nil {
			t.CategoryName = *catName
			t.CategoryColor = *catColor
			t.CategoryIcon = *catIcon
		} else {
			t.CategoryName = "Uncategorized"
		}
		transactions = append(transactions, t)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}

func createTransaction(w http.ResponseWriter, r *http.Request) {
	var t Transaction
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := db.Exec("INSERT INTO transactions (amount, description, date, type, category_id) VALUES (?, ?, ?, ?, ?)",
		t.Amount, t.Description, t.Date, t.Type, t.CategoryID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := res.LastInsertId()
	t.ID = int(id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

func updateTransaction(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/api/transactions/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var t Transaction
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update query
	_, err = db.Exec("UPDATE transactions SET amount=?, description=?, date=?, type=?, category_id=? WHERE id=?",
		t.Amount, t.Description, t.Date, t.Type, t.CategoryID, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t.ID = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

func deleteTransaction(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/api/transactions/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("DELETE FROM transactions WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Dashboard Summary

func getDashboardSummary(w http.ResponseWriter, r *http.Request) {
	month := r.URL.Query().Get("month")
	if month == "" {
		// default to current month? or just return global?
		// for now require month or return 0
	}

	var income, expense float64

	// Calc Income
	db.QueryRow(`SELECT COALESCE(SUM(amount), 0) FROM transactions WHERE type='income' AND strftime('%Y-%m', date) = ?`, month).Scan(&income)

	// Calc Expense
	db.QueryRow(`SELECT COALESCE(SUM(amount), 0) FROM transactions WHERE type='expense' AND strftime('%Y-%m', date) = ?`, month).Scan(&expense)

	balance := income - expense

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"income":  income,
		"expense": expense,
		"balance": balance,
	})
}

func getCategoryStats(w http.ResponseWriter, r *http.Request) {
	month := r.URL.Query().Get("month")

	// 1. Get Total Expense for the month first
	var totalExpense float64
	err := db.QueryRow("SELECT COALESCE(SUM(amount), 0) FROM transactions WHERE type = 'expense' AND strftime('%Y-%m', date) = ?", month).Scan(&totalExpense)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 2. Get Expense by Category
	query := `
        SELECT c.name, c.color, c.icon, COALESCE(SUM(t.amount), 0) as total
        FROM transactions t
        JOIN categories c ON t.category_id = c.id
        WHERE t.type = 'expense' AND strftime('%Y-%m', t.date) = ?
        GROUP BY c.id
        ORDER BY total DESC
    `

	rows, err := db.Query(query, month)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type Stat struct {
		Name       string  `json:"name"`
		Color      string  `json:"color"`
		Icon       string  `json:"icon"`
		Amount     float64 `json:"amount"`
		Percentage float64 `json:"percentage"`
	}

	var stats []Stat
	for rows.Next() {
		var s Stat
		if err := rows.Scan(&s.Name, &s.Color, &s.Icon, &s.Amount); err != nil {
			continue
		}
		if totalExpense > 0 {
			s.Percentage = (s.Amount / totalExpense) * 100
		} else {
			s.Percentage = 0
		}
		stats = append(stats, s)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
