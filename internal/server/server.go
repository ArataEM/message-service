package server

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	addr string
	db   *sql.DB
}

func NewServer(addr string, db *sql.DB) *Server {
	return &Server{
		addr: addr,
		db:   db,
	}
}

func (s *Server) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	messageHandler := NewHandler()
	messageHandler.RegisterRoutes(subrouter)

	slog.Info(fmt.Sprintf("Listening on %s", s.addr))

	return http.ListenAndServe(s.addr, router)
}
