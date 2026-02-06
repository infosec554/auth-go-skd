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
	// 1. Load Config (Optional, you can just hardcode strings for simple apps)
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// 2. Initialize Service (One-line setup)
	service := auth.New(auth.Opts{
		Secret: "super-secret-key-change-me",
		URL:    "http://localhost:" + cfg.HTTP.Port,
	})

	// 3. Add Providers
	service.Add(google.New(
		cfg.OAuth.Google.ClientID,
		cfg.OAuth.Google.ClientSecret,
		cfg.OAuth.Google.RedirectURL,
	))

	// 4. Setup Router
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Mount Auth Handlers (One-line mount)
	authHandler, avatarHandler := service.Handlers()
	r.Mount("/auth", authHandler)
	r.Mount("/avatar", avatarHandler)

	// 5. Protected Routes
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
