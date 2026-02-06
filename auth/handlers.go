package auth

import (
	"encoding/json"
	"net/http"

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
	user, err := p.FetchUser(r.Context(), code)
	if err != nil {
		http.Error(w, "failed to login", http.StatusInternalServerError)
		return
	}

	// TODO: Create JWT, set cookie
	json.NewEncoder(w).Encode(user)
}

func (s *Service) logoutHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Clear cookie
	w.WriteHeader(http.StatusOK)
}

func (s *Service) directLoginHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement direct login
	w.WriteHeader(http.StatusNotImplemented)
}
