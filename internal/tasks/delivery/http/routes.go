package http

import "github.com/go-chi/chi/v5"

func MapTasksRoutes(router chi.Router, handlers *TaskHandlers) {
	router.Post("/task", handlers.Create())
	router.Get("/task/{id}", handlers.Get())
}
