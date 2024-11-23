package application

import (
	"net/http"

	"github.com/ArataEM/message-service/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (a *App) loadRoutes() {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello Artem!"))
	})

	router.Route("/messages", a.loadMessageRoutes)

	a.router = router
}

func (a *App) loadMessageRoutes(router chi.Router) {
	messageHandler := &handler.Message{
		Repo: a.rdb,
	}

	router.Get("/", messageHandler.List)
	router.Post("/", messageHandler.Create)
	router.Get("/{id}", messageHandler.GetById)
	router.Put("/{id}", messageHandler.UpdateById)
	router.Delete("/{id}", messageHandler.DeleteById)
}
