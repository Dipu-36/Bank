package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
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
	router.HandleFunc("/account/{id}", withJwtAuth(makeHTTPHandleFunc(s.handleGetAccountbyID), s.store))

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
	accounts, err := s.store.GetAccounts()
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, accounts)
}

func (s *APIServer) handleGetAccountbyID(w http.ResponseWriter, r *http.Request) error {
	//account := NewAccount("Diptendu", "Pal")

	// We are fetching the ID from the request or the param therefore it is going to be in tokenString
	// therefore we need to convert it to Int using the strconv.Atoi

	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fmt.Errorf("Invalid id given %s", idStr)
	}

	account, err := s.store.GetAccountbyID(id)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {

	//Allocating memory for the struct CreateAccountRequest & returning a pointer to that struct
	createAccountReq := new(CreateAccountRequest)

	//Reading the raw JSON body from the HTTP request & decoding it back to struct, returns an err
	if err := json.NewDecoder(r.Body).Decode(createAccountReq); err != nil {
		return err
	}

	//Calling the constructor function to create a new Account struct and pass the values to in
	//account variable stores the struct fields, in struct form, decoded from JSON body
	account := NewAccount(createAccountReq.FirstName, createAccountReq.LastName)

	//Using store field of the APIServer instance s to get the interface methods to store the Account info in DB
	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	//Creating a JWT token for the user when he creates his acccount
	tokenString, err := createJWT(account)
	if err != nil {
		return err
	}

	fmt.Println("JWT token: ", tokenString)

	return WriteJSON(w, http.StatusOK, account)
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
func createJWT(account *Account) (string, error) {
	claims := &jwt.MapClaims{
		"expiresAt":     15000,
		"accountNumber": account.Number,
	}

	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50TnVtYmVyIjoyOTczMywiZXhwaXJlc0F0IjoxNTAwMH0.8qv6KZ-5Ps6Gj2lkFpQwz2GVJ0rKyrI_xJS6jLpzIgU
func permissionDenied(w http.ResponseWriter) {
	WriteJSON(w, http.StatusForbidden, ApiError{Error: "Permission Denied"})
}

func withJwtAuth(handlerFunc http.HandlerFunc, s Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("calling JWT auth middleware")

		tokenString := r.Header.Get("x-jwt-token")

		token, err := validateJWT(tokenString)
		if err != nil {
			permissionDenied(w)
			return
		}

		//claims := token.Claims.(jwt.MapClaims)
		if !token.Valid {
			permissionDenied(w)
			return
		}

		userID, err := getID(r)
		if err != nil {
			permissionDenied(w)
		}

		account, err := s.GetAccountbyID(userID)
		if err != nil {
			permissionDenied(w)
			return
		}
		claims := token.Claims.(jwt.MapClaims)

		if account.Number != int64(claims["account_number"].(float64)) {
			permissionDenied(w)
			return
		}

		handlerFunc(w, r)
	}
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
}

// apiFunc defines the signature for the API handler functions that returns errors
// This custom type allows handlers to return errors that will be automatically converted to JSON error responses
type apiFunc func(http.ResponseWriter, *http.Request) error

// ApiError represents a standerdized error response structure
// This ensures all API errors allow the same JSON format
type ApiError struct {
	Error string `json:"error"`
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
