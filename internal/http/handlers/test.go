package handlers

import (
	"github.com/go-chi/render"
	"net/http"
)

// HelloWorld godoc
// @Summary HelloWorld summary
// @Description HelloWorld description
// @Tags hw tag
// @Accept json
// @Produce json
// @Success 200 {string} string	"ok"
// @Router / [get]
func HelloWorld() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, "Hello World")
	}
}
