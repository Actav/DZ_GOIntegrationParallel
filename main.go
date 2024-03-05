package main

import (
	"dz_go/handlers"
	"dz_go/storage"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {
	s := storage.NewInMemoryStorage()
	h := handlers.Handlers{Storage: s}
	r := chi.NewRouter()

	r.Post("/create", h.CreateUser)
	r.Post("/make_friends", h.MakeFriends)
	r.Get("/friends/{userID}", h.ListFriends)
	r.Get("/user/{userID}", h.GetUser)
	r.Put("/user/{userID}", h.UpdateUser)
	r.Delete("/user", h.DeleteUser)

	log.Println("Start server")
	http.ListenAndServe(":8080", r)
}
