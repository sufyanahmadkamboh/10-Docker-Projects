package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

type Item struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

var db *sql.DB

func main() {
	var err error

	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "goapp")
	dbSSLMode := getEnv("DB_SSLMODE", "disable")
	port := getEnv("PORT", "8080")

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode,
	)

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	log.Println("Connected to Postgres ✅")

	mux := http.NewServeMux()

	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/api/items", itemsHandler)     // GET, POST
	mux.HandleFunc("/api/items/", itemByIDHandler) // GET, PUT, DELETE

	address := ":" + port
	log.Printf("Server running on %s\n", address)
	if err := http.ListenAndServe(address, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func getEnv(key, def string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return def
}

// ---------- Handlers ----------

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, `
		<!DOCTYPE html>
		<html>
		<head><title>Go CRUD API</title></head>
		<body style="font-family: sans-serif;">
			<h1>🚀 Go CRUD API</h1>
			<p>API is running.</p>
			<p>Try:</p>
			<ul>
				<li><code>GET /api/items</code></li>
				<li><code>POST /api/items</code> with JSON body: {"name": "Item 1", "description": "Test"}</li>
				<li><code>GET /api/items/{id}</code></li>
				<li><code>PUT /api/items/{id}</code></li>
				<li><code>DELETE /api/items/{id}</code></li>
			</ul>
		</body>
		</html>
	`)
}

func itemsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getItems(w, r)
	case http.MethodPost:
		createItem(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func itemByIDHandler(w http.ResponseWriter, r *http.Request) {
	// Path: /api/items/{id}
	idStr := strings.TrimPrefix(r.URL.Path, "/api/items/")
	if idStr == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		getItem(w, r, id)
	case http.MethodPut:
		updateItem(w, r, id)
	case http.MethodDelete:
		deleteItem(w, r, id)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// ---------- CRUD Logic ----------

func getItems(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name, description FROM items ORDER BY id")
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		log.Println("getItems error:", err)
		return
	}
	defer rows.Close()

	var items []Item

	for rows.Next() {
		var it Item
		if err := rows.Scan(&it.ID, &it.Name, &it.Description); err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			log.Println("getItems scan error:", err)
			return
		}
		items = append(items, it)
	}

	writeJSON(w, items)
}

func getItem(w http.ResponseWriter, r *http.Request, id int) {
	var it Item
	err := db.QueryRow("SELECT id, name, description FROM items WHERE id = $1", id).
		Scan(&it.ID, &it.Name, &it.Description)

	if err == sql.ErrNoRows {
		http.Error(w, "item not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		log.Println("getItem error:", err)
		return
	}

	writeJSON(w, it)
}

func createItem(w http.ResponseWriter, r *http.Request) {
	var input Item
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if input.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	err := db.QueryRow(
		"INSERT INTO items (name, description) VALUES ($1, $2) RETURNING id",
		input.Name, input.Description,
	).Scan(&input.ID)

	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		log.Println("createItem error:", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	writeJSON(w, input)
}

func updateItem(w http.ResponseWriter, r *http.Request, id int) {
	var input Item
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if input.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	res, err := db.Exec(
		"UPDATE items SET name = $1, description = $2 WHERE id = $3",
		input.Name, input.Description, id,
	)
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		log.Println("updateItem error:", err)
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "item not found", http.StatusNotFound)
		return
	}

	input.ID = id
	writeJSON(w, input)
}

func deleteItem(w http.ResponseWriter, r *http.Request, id int) {
	res, err := db.Exec("DELETE FROM items WHERE id = $1", id)
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		log.Println("deleteItem error:", err)
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "item not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ---------- Helpers ----------

func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Println("writeJSON error:", err)
	}
}
