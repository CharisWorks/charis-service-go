package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

func main() {
	r := chi.NewRouter()
	r.Get("/register", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("register"))
	})
	log.Print("server started")
	http.ListenAndServe(":8080", r)
}
