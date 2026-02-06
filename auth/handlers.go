package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
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

// generateState creates a random state string
func generateState() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func (s *Service) loginHandler(w http.ResponseWriter, r *http.Request) {
	providerName := chi.URLParam(r, "provider")
	p, ok := s.providers[providerName]
	if !ok {
		http.Error(w, "provider not found", http.StatusNotFound)
		return
	}

	// 1. Generate Secure State
	state := generateState()

	// 2. Set State in secure, short-lived cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Path:     "/",
		Expires:  time.Now().Add(10 * time.Minute),
		HttpOnly: true,
		Secure:   r.TLS != nil || s.opts.URLIsHTTPS, // Auto-detect HTTPS or config
		SameSite: http.SameSiteLaxMode,
	})

	// 3. Redirect to Provider
	url := p.GetAuthURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (s *Service) callbackHandler(w http.ResponseWriter, r *http.Request) {
	providerName := chi.URLParam(r, "provider")
	p, ok := s.providers[providerName]
	if !ok {
		http.Error(w, "provider not found", http.StatusNotFound)
		return
	}

	// 1. Validate State (CSRF Protection)
	stateParam := r.URL.Query().Get("state")
	stateCookie, err := r.Cookie("oauth_state")
	if err != nil || stateCookie.Value != stateParam {
		http.Error(w, "invalid oauth state", http.StatusForbidden)
		return
	}

	// Clear state cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	})

	// 2. Exchange Code for User
	code := r.URL.Query().Get("code")
	user, err := p.FetchUser(r.Context(), code)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to login: %v", err), http.StatusInternalServerError)
		return
	}

	// 3. Create JWT
	tokenStr, err := s.Token(user)
	if err != nil {
		http.Error(w, "failed to create token", http.StatusInternalServerError)
		return
	}

	// 4. Set Session Cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "JWT",
		Value:    tokenStr,
		HttpOnly: true,
		Secure:   r.TLS != nil || s.opts.URLIsHTTPS,
		Path:     "/",
		Expires:  time.Now().Add(s.opts.CookieDuration),
		SameSite: http.SameSiteLaxMode,
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
