package server

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/message", h.handleCreateMessage).Methods(http.MethodPost)
	router.HandleFunc("/message", h.handleGetMessage).Methods(http.MethodGet)
}

func (h *Handler) handleCreateMessage(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) handleGetMessage(w http.ResponseWriter, r *http.Request) {
	slog.Info("Get message-service invoked")
}
