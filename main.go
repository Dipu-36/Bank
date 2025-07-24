package main

import (
	"net/http"
)

type APIServer struct {
	listenAddr string
}

func NewAPIServer(port string) *APIServer {
	return &APIServer{
		listenAddr: port,
	}
}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}
