package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Simple in-memory session store
var (
	sessions = make(map[string]int) // token -> user_id
	mu       sync.Mutex
)

// Login Request DTO
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func apiLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Fetch user from DB
	var id int
	var hash string
	err := db.QueryRow("SELECT id, password_hash FROM users WHERE username = ?", req.Username).Scan(&id, &hash)
	if err != nil {
		// User not found
		time.Sleep(100 * time.Millisecond) // mitigate timing attack
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate session token
	token := generateToken()
	mu.Lock()
	sessions[token] = id
	mu.Unlock()

	// Set cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func apiLogout(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	if err == nil {
		mu.Lock()
		delete(sessions, c.Value)
		mu.Unlock()
	}

	// Expire cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "logged out"})
}

// Middleware to protect routes
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Skip auth for login page and api
		if r.URL.Path == "/login" || r.URL.Path == "/api/login" || r.URL.Path == "/assets/" {
			next(w, r)
			return
		}

		c, err := r.Cookie("session_token")
		if err != nil {
			handleUnauth(w, r)
			return
		}

		mu.Lock()
		userId, ok := sessions[c.Value]
		mu.Unlock()

		if !ok {
			handleUnauth(w, r)
			return
		}

		// Add user_id to context
		ctx := context.WithValue(r.Context(), "user_id", userId)
		next(w, r.WithContext(ctx))
	}
}

func handleUnauth(w http.ResponseWriter, r *http.Request) {
	if isAPI(r.URL.Path) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	} else {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func isAPI(path string) bool {
	return len(path) >= 4 && path[0:4] == "/api"
}

func generateToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// Helper to get user ID from context
func getUserID(r *http.Request) int {
	id, ok := r.Context().Value("user_id").(int)
	if !ok {
		log.Println("Error: User ID not found in context")
		return 0
	}
	return id
}
