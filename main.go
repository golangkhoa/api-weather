package main

import (
	"api-weather/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", jsonHandler)
	r.Get("/location/{place}", jsonHandler)
	http.ListenAndServe(":3000", r)
}

func jsonHandler(w http.ResponseWriter, r *http.Request) {
	json := r.URL.Query().Get("json")
	q := chi.URLParam(r, "place")
	if json == "true" {
		handlers.JsonHandler(q, w, r)
	} else if json != "true" {
		handlers.PlaceHandler(q, w, r)
	}
}
