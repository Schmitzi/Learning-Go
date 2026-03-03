package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"strconv"
)

const serverPort = "8080"

type Application struct {
	db *sql.DB
}

type Book struct {
	ID       int     `json:"id"`
	Title    string  `json:"title"`
	Author   string  `json:"author"`
	Price    float64 `json:"price"`
	ImageURL string  `json:"imageurl"`
}

func (app *Application) addBook(w http.ResponseWriter, req *http.Request) {
	// Set Content-Type Header
	w.Header().Set("Content-Type", "application/json")

	var b Book

	err := json.NewDecoder(req.Body).Decode(&b)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = app.db.QueryRow(
		"INSERT INTO books (title, author, price, imageurl) VALUES ($1, $2, $3, $4) RETURNING id",
		b.Title, b.Author, b.Price, b.ImageURL,
	).Scan(&b.ID)

	// Handle Errors
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Set status code to 201
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(b)
}

func (app *Application) getAllBooks(w http.ResponseWriter, r *http.Request) {
	// Set Content-Type Header
	w.Header().Set("Content-Type", "application/json")

	// Save all entires to rows
	rows, err := app.db.Query("SELECT id, title, author, price, imageurl FROM books")

	// Handle Errors
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	// Create response object
	var books []Book
	for rows.Next() {
		var b Book
		err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.Price, &b.ImageURL)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		books = append(books, b)
	}

	// Create Encoder
	err = json.NewEncoder(w).Encode(books)

	// If error, return Internal Server Error
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	// If no error, 200 is sent
}

func (app *Application) getBooksById(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")

	// Set Content-Type Header
	w.Header().Set("Content-Type", "application/json")

	// Convert id to int
	var idx, err = strconv.Atoi(id)

	// Handle Errors
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Save all entires to rows
	item := app.db.QueryRow("SELECT * FROM books WHERE id = $1", idx)

	// Create response object
	var b Book
	err = item.Scan(&b.ID, &b.Title, &b.Author, &b.Price, &b.ImageURL)

	// If error
	if err == sql.ErrNoRows {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Create Encoder
	err = json.NewEncoder(w).Encode(b)

	// If no error, 200 is sent
}

func (app *Application) updateBook(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")

	// Set Content-Type Header
	w.Header().Set("Content-Type", "application/json")

	var b Book

	err := json.NewDecoder(req.Body).Decode(&b)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Convert id to int
	idx, err := strconv.Atoi(id)

	// Handle Errors
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = app.db.QueryRow(
		"UPDATE books SET title=$1, author=$2, price=$3, imageurl=$4 WHERE id=$5 RETURNING id",
		b.Title, b.Author, b.Price, b.ImageURL, idx,
	).Scan(&b.ID)

	// Handle Errors
	if err == sql.ErrNoRows {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(b)
}

func (app *Application) deleteBook(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")

	// Convert id to int
	var idx, err = strconv.Atoi(id)

	// Handle Errors
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Delete item
	item, err := app.db.Exec("DELETE FROM books WHERE id = $1", idx)

	// Handle Errors
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	rows, err := item.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if rows == 0 {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	// Set status code to 204
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	// Welcome Message
	fmt.Println("Welcome to the BookStore")

	// Load env
	godotenv.Load()

	// Create Database
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))

	// Error Handling
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	// Create Router
	router := chi.NewRouter()

	// Create App
	app := &Application{db: db}

	// Function Handlers
	router.Get("/books", app.getAllBooks)
	router.Get("/books/{id}", app.getBooksById)
	router.Post("/books", app.addBook)
	router.Put("/books/{id}", app.updateBook)
	router.Delete("/books/{id}", app.deleteBook)

	fmt.Print("Server listening on localhost:" + serverPort)

	// Start Server
	log.Fatal(http.ListenAndServe(":"+serverPort, router))
}
