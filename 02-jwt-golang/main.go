package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
)

// Struct defines
type statusMessage struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type loginForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type tokenMessage struct {
	Token string `json:"token"`
}

// Consts
const serverPort = "4000"

// Should be an environmental variable, fine for learning though
var jwtSecret = []byte("supersecret")

// Same here. Not good practice
var users = map[string]string{
	"schmitzi": "password",
}

func welcome(w http.ResponseWriter, r *http.Request) {
	// Set Content-Type
	w.Header().Set("Content-Type", "application/json")

	// Define JSON object
	data := statusMessage{
		Status:  "Success",
		Message: "Welcome to Golang with JWT authentication",
	}

	// Encode object
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Return and send 200
}

func secure(w http.ResponseWriter, r *http.Request) {
	// Extract, JWT token from Header
	token := r.Header.Get("Authorization")

	// Check if "Bearer " exists in the header
	ok := strings.HasPrefix(token, "Bearer ")
	if !ok {
		// Set 401 status code
		w.WriteHeader(http.StatusUnauthorized)

		// Form response object
		ret := statusMessage{
			Status:  "Failure",
			Message: "You are not authorized to view this page",
		}

		// If success, then encode the result
		err := json.NewEncoder(w).Encode(ret)

		// Handle error
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		// Return with 401
		return
	}

	// Remove "Bearer " prefix
	token = strings.TrimPrefix(token, "Bearer ")

	// Check that the token was signed with HS256
	parsed, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return jwtSecret, nil
	})

	// Handle error
	if err != nil {
		http.Error(w, err.Error(), 401)
		return
	}

	// Check parsed tokens validity
	if parsed.Valid {
		// Grab username from JWT claim to use in response message
		claims, ok := parsed.Claims.(jwt.MapClaims)
		if ok != true {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		// Form success message
		mes := "Congrats " + claims["sub"].(string) + " and Welcome to the Secure page!"
		ret := statusMessage{
			Status:  "Success",
			Message: mes,
		}

		// Encode message
		err = json.NewEncoder(w).Encode(ret)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		// Return with 200 status code
		return
	} else {
		// Set status code
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	// Set Content-Type Header
	w.Header().Set("Content-Type", "application/json")

	var l loginForm

	// Decode body from request
	err := json.NewDecoder(r.Body).Decode(&l)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Find User: Password from map of users
	expectedPassword, ok := users[l.Username]

	// Handle errors
	if !ok {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	// Compare password given with saved password
	if expectedPassword != l.Password {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Create JWT claims
	var jwtClaims = jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub": l.Username,
			"exp": time.Now().Add(24 * time.Hour).Unix(),
			"iat": time.Now().Unix(),
		})

	// Sign the string
	s, err := jwtClaims.SignedString(jwtSecret)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Prepare the token to return
	ret := tokenMessage{
		Token: s,
	}

	// Encode signed string
	err = json.NewEncoder(w).Encode(ret)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Return with 200 OK
}

func main() {

	// Create Router
	router := chi.NewRouter()

	log.Printf("Server listening on localhost:%s", serverPort)

	// Function Handlers
	router.Get("/", welcome)
	router.Get("/secure", secure)
	router.Post("/login", login)

	// Start Server
	log.Fatal(http.ListenAndServe(":"+serverPort, router))
}
