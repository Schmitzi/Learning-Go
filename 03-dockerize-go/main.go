// @title Dockerize-Go 
// @version 1.1
// @description Simple Go REST API with Swagger Documentation

// @contact.name Michael Naysmith
// @contact.url https://schmitzi.nz
// @contact.email contact@schmitzi.nz

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
// @query.collection.format multi

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
	_ "bookstore-api/docs"
	"strconv"
	httpSwagger "github.com/swaggo/http-swagger"
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

// @Summary Return Welcome Message
// @ID welcomeMessage
// @Produce json
// @Success 200 {string} string "welcome message"
// @Failure 500 
// @Router / [get]
func (app *Application) welcomeMessage(w http.ResponseWriter, req *http.Request) {
	msg := "Welcome to the Bookstore"

  err := json.NewEncoder(w).Encode(msg)

	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

// @Summary Add Book to the Library
// @ID addBook
// @Produce json
// @Param book body Book true "Book object"
// @Success 201 {object} Book
// @Router /books [post]
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

// @Summary Get All Books
// @ID getBooks
// @Produce json
// @Success 200 {array} Book
// @Failure 404 {string} string "Library is empty"
// @Failure 500 {string} string "Internal server error"
// @Router /books [get]
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

// @Summary Get Book by ID
// @ID getBookByID
// @Produce json
// @Param id path int true "Book ID"
// @Success 200 {object} Book
// @Failure 404 {string} string "Book not found"
// @Failure 500 {string} string "Internal server error"
// @Router /books/{id} [get]
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

// @Summary Update Book by ID
// @ID updateBooks
// @Produce json
// @Param id path int true "Book ID"
// @Param book body Book true "Book object"
// @Success 200 {object} Book
// @Failure 404 {string} string "Library is empty"
// @Failure 500 {string} string "Internal server error"
// @Router /books/{id} [put]
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

// @Summary Delete Book by ID
// @ID deleteBooks
// @Produce json
// @Param id path int true "Book ID"
// @Success 204
// @Failure 404 {string} string "Library is empty"
// @Failure 500 {string} string "Internal server error"
// @Router /books/{id} [delete]
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

	// Get server serverPort
	port := os.Getenv("PORT")

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

	// Swagger Handler
	router.Get("/swagger/*", httpSwagger.WrapHandler)

	// Function Handlers
	router.Get("/", app.welcomeMessage)
	router.Get("/books", app.getAllBooks)
	router.Get("/books/{id}", app.getBooksById)
	router.Post("/books", app.addBook)
	router.Put("/books/{id}", app.updateBook)
	router.Delete("/books/{id}", app.deleteBook)

	fmt.Print("Server listening on localhost:" + port)


	// Start Server
	log.Fatal(http.ListenAndServeTLS(":"+port, "cert.pem", "key.pem", router))
}
