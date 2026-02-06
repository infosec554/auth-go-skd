package auth

import (
	"net/http"
)

type Middleware struct {
	service *Service
}

func (m *Middleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Logic to check token
		next.ServeHTTP(w, r)
	})
}
