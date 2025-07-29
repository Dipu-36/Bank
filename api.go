package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// APIServer represents our HTTP server instance.
// Contains configuration and handles routing.
type APIServer struct {
	listenAddr string
	store      Storage
}

// NewAPIServer creates a new APIServer instance.
// This constructor pattern is idiomatic Go for initialization.
// port: The port number to listen on (e.g., ":8080")
func NewAPIServer(port string, db Storage) *APIServer {
	return &APIServer{
		listenAddr: port,
		store:      db,
	}
}

// Run starts the HTTP server and configures routes.
// This is the main entry point to launch the server.
// Uses gorilla/mux for routing.
func (s *APIServer) Run() {
	// Create a new router instance
	router := mux.NewRouter()

	//Registers woutes with our handler decorater
	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleAccount))
	router.HandleFunc("/account{id}", makeHTTPHandleFunc(s.handleGetAccount))
	log.Println("JSON API server running on port: ", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAccount(w, r)
	}
	if r.Method == "POST" {
		return s.handleCreateAccount(w, r)
	}
	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	//account := NewAccount("Diptendu", "Pal")
	id := mux.Vars(r)["id"]
	fmt.Println(id)

	return WriteJSON(w, http.StatusOK, id)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleTranasfer(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// WriteJSON writes JSON response with the given status code , it sets the Content-Type as application-json Content-Type
// It also encodes the provides value to v & returns an erroor if JSON encoding fails
func WriteJSON(w http.ResponseWriter, status int, v any) error {

	//sets the content type to json
	w.Header().Set("Content-Type", "application/json")

	//sets the status code before writing the headers
	w.WriteHeader(status)

	//Encode the response body as JSON
	//http.ResponseWriter implements io.Writer, which NewEncoder requires
	return json.NewEncoder(w).Encode(v)
}

// apiFunc defines the signature for the API handler functions that returns errors
// This custom type allows handlers to return errors that will be automatically converted to JSON error responses
type apiFunc func(http.ResponseWriter, *http.Request) error

// ApiError represents a standerdized error response structure
// This ensures all API errors allow the same JSON format
type ApiError struct {
	Error string
}

// makeHTTPHandleFunc converts an apiFunc to a http.HandlerFunc.
// This decorator wraps our custom handlers to provide:
// - Automatic error handling
// - Consistent JSON error responses
// - Conversion between handler signatures
// Usage:
// router.HandleFunc("/path", makeHTTPHandleFunc(handler))
func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			//handle the error here
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}

}
