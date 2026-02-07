package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"auth-go-skd/auth"
	"auth-go-skd/config"
	"auth-go-skd/provider/google"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	service := auth.New(auth.Opts{
		Secret: "super-secret-key-change-me",
		URL:    "http://localhost:" + cfg.HTTP.Port,
	})

	service.Add(google.New(
		cfg.OAuth.Google.ClientID,
		cfg.OAuth.Google.ClientSecret,
		cfg.OAuth.Google.RedirectURL,
	))

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	authHandler, avatarHandler := service.Handlers()
	r.Mount("/auth", authHandler)
	r.Mount("/avatar", avatarHandler)

	r.Group(func(r chi.Router) {
		r.Use(service.Middleware().Auth)
		r.Get("/private", func(w http.ResponseWriter, r *http.Request) {
			user := auth.User(r)
			w.Write([]byte("Hello " + user.Name))
		})
	})

	log.Printf("Server listening on port %s", cfg.HTTP.Port)
	http.ListenAndServe(":"+cfg.HTTP.Port, r)
}
