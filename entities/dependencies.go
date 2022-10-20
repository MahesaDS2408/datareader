package entities

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Dependencies struct {
	Config *YodelConfig
	Router *chi.Mux
}

func (d Dependencies) Run() {
	http.ListenAndServe(":8080", d.Router)
}
