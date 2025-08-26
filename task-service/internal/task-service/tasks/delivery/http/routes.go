package http

import "github.com/go-chi/chi/v5"

// MapTasksRoutes creates routes for task handlers.
func MapTasksRoutes(router chi.Router, handlers *TaskHandlers) {
	router.Post("/task", handlers.Create())
	router.Get("/task/{id}", handlers.Get())
}
