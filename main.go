package main

import (
	"fmt"
	"net/http"

	envy "github.com/codegangsta/envy/lib"
	"github.com/go-chi/chi/v5"
)

func main() {
	fmt.Println("Hello world")

	r := chi.NewRouter()

	r.Get("/sites/{id}", Edit)
	r.Put("/sites/{id}", Update)

	r.Get("/callback", AuthCallback)

	port := envy.MustGet("PORT")
	fmt.Println("Listening on port :" + port)
	http.ListenAndServe(":"+port, r)
}

func Edit(w http.ResponseWriter, r *http.Request) {
}

func Update(w http.ResponseWriter, r *http.Request) {
}
