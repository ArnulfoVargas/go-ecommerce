package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *Application) Routes() http.Handler {
  mux := chi.NewRouter()

  return mux
}
