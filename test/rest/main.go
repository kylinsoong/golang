package main

import (
	"encoding/json"
	"fmt"
        "strconv"
	"net/http"
	"github.com/gorilla/mux"
)

// User represents the data structure for a user.
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// In-memory storage for users (replace with a database in a real-world scenario).
var users []User

// Handler for the POST /users endpoint.
func createUserHandler(w http.ResponseWriter, r *http.Request) {
	// Decode JSON from request body into User struct
	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Assign a unique ID (for simplicity in this example)
	newUser.ID = len(users) + 1

	// Add the new user to the in-memory storage
	users = append(users, newUser)

	// Respond with the created user as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}

// Handler for the GET /users/{id} endpoint.
func getUserHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the "id" parameter from the request URL
	vars := mux.Vars(r)
	userID := vars["id"]

	// Convert the user ID to an integer
	id, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Find the user by ID in the in-memory storage
	var foundUser User
	for _, user := range users {
		if user.ID == id {
			foundUser = user
			break
		}
	}

	// Check if the user was found
	if foundUser.ID == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Respond with the found user as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(foundUser)
}

func main() {
	// Create a new router using mux
	r := mux.NewRouter()

	// Register the handler functions for the /users endpoints
	r.HandleFunc("/users", createUserHandler).Methods("POST")
	r.HandleFunc("/users/{id}", getUserHandler).Methods("GET")

	// Start the server on port 8080
	port := 8080
	fmt.Printf("Server listening on :%d...\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), r)
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}

