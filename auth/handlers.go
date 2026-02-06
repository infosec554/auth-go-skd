package auth

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

// Handlers returns the http handlers for auth and avatar
func (s *Service) Handlers() (http.Handler, http.Handler) {
	r := chi.NewRouter()

	r.Get("/{provider}/login", s.loginHandler)
	r.Get("/{provider}/callback", s.callbackHandler)
	r.Post("/logout", s.logoutHandler)

	// Direct auth routes (simplified)
	r.Post("/login", s.directLoginHandler)

	avatarRouter := chi.NewRouter()
	// avatarRouter.Get("/{id}", s.avatarHandler)

	return r, avatarRouter
}

func (s *Service) loginHandler(w http.ResponseWriter, r *http.Request) {
	providerName := chi.URLParam(r, "provider")
	p, ok := s.providers[providerName]
	if !ok {
		http.Error(w, "provider not found", http.StatusNotFound)
		return
	}
	url := p.GetAuthURL("state") // TODO: secure state
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (s *Service) callbackHandler(w http.ResponseWriter, r *http.Request) {
	providerName := chi.URLParam(r, "provider")
	p, ok := s.providers[providerName]
	if !ok {
		http.Error(w, "provider not found", http.StatusNotFound)
		return
	}

	code := r.URL.Query().Get("code")
	// User fetched from provider
	user, err := p.FetchUser(r.Context(), code)
	if err != nil {
		http.Error(w, "failed to login", http.StatusInternalServerError)
		return
	}

	// Create JWT
	tokenStr, err := s.Token(user)
	if err != nil {
		http.Error(w, "failed to create token", http.StatusInternalServerError)
		return
	}

	// Set Cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "JWT",
		Value:    tokenStr,
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Now().Add(s.opts.CookieDuration),
	})

	// Redirect or return JSON
	// For SDK, usually redirects to frontend or returns JSON.
	// Let's just return JSON for now as generic behavior
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token": tokenStr,
		"user":  user,
	})
}

func (s *Service) logoutHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Clear cookie
	w.WriteHeader(http.StatusOK)
}

func (s *Service) directLoginHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement direct login
	w.WriteHeader(http.StatusNotImplemented)
}
