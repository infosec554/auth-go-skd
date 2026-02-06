package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"auth-go-skd/auth"
	"auth-go-skd/avatar"
	"auth-go-skd/config"
	"auth-go-skd/provider/google"
)

func main() {
	// 1. Load Config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// 2. Setup Auth Options
	opts := auth.Opts{
		SecretReader: func(id string) (string, error) {
			return "secret-key-change-me", nil
		},
		TokenDuration:  time.Minute * 15,
		CookieDuration: time.Hour * 24,
		Issuer:         cfg.App.Name,
		URL:            "http://localhost:" + cfg.HTTP.Port,
		AvatarStore:    avatar.NewLocalFS("/tmp/auth-avatars"),
	}

	// 3. Create Auth Service
	service := auth.NewService(opts)

	// 4. Add Providers
	// Google
	gProv := google.New(cfg.OAuth.Google)
	service.AddCustomProvider(gProv)

	// 5. Setup Router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// 6. Mount Auth Handlers
	authHandlers, avatarHandlers := service.Handlers()
	r.Mount("/auth", authHandlers)
	r.Mount("/avatar", avatarHandlers)

	// 7. Protected Route Example
	m := service.Middleware()
	r.Group(func(r chi.Router) {
		r.Use(m.Auth)
		r.Get("/private", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("You are authenticated!"))
		})
	})

	log.Printf("Server listening on port %s", cfg.HTTP.Port)
	if err := http.ListenAndServe(":"+cfg.HTTP.Port, r); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
