package main

import (
	"embed"
	"html/template"
	"io/fs"
	"log"
	"net/http"
)

// Embed templates and static assets
//
//go:embed templates/* assets/*
var resultFS embed.FS

func main() {
	// 1. Initialize Database
	InitDB("database.db")

	// 2. Setup Routing
	mux := http.NewServeMux()

	// Static Files (CSS, JS)
	assetsRoot, _ := fs.Sub(resultFS, "assets")
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.FS(assetsRoot))))

	// Page Routes (Protected)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	})
	mux.HandleFunc("/dashboard", authMiddleware(handleDashboard))
	mux.HandleFunc("/transactions", authMiddleware(handleTransactions))
	mux.HandleFunc("/categories", authMiddleware(handleCategories))

	// Public Routes
	mux.HandleFunc("/login", handleLogin)

	// API Routes (Auth)
	mux.HandleFunc("/api/login", apiLogin)
	mux.HandleFunc("/api/logout", apiLogout)

	// API Routes (Protected)
	mux.HandleFunc("/api/categories", authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getCategories(w, r)
		} else if r.Method == "POST" {
			createCategory(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.HandleFunc("/api/categories/", authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "DELETE" {
			deleteCategory(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.HandleFunc("/api/transactions", authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getTransactions(w, r)
		} else if r.Method == "POST" {
			createTransaction(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.HandleFunc("/api/transactions/", authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "DELETE" {
			deleteTransaction(w, r)
		} else if r.Method == "PUT" {
			updateTransaction(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	mux.HandleFunc("/api/summary", authMiddleware(getDashboardSummary))
	mux.HandleFunc("/api/stats/category", authMiddleware(getCategoryStats))

	// 3. Start Server
	port := ":8081"
	log.Printf("🚀 Server starting on http://localhost%s\n", port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatal(err)
	}
}

// --- Page Handlers ---

func handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		basicRender(w, "templates/login.html")
	}
}

func handleDashboard(w http.ResponseWriter, r *http.Request) {
	basicRender(w, "templates/dashboard.html")
}

func handleTransactions(w http.ResponseWriter, r *http.Request) {
	basicRender(w, "templates/transactions.html")
}

func handleCategories(w http.ResponseWriter, r *http.Request) {
	basicRender(w, "templates/categories.html")
}

// Helpers
func basicRender(w http.ResponseWriter, path string) {
	tmpl, err := template.ParseFS(resultFS, path)
	if err != nil {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}
	tmpl.Execute(w, nil)
}
