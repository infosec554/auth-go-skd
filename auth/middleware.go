package auth

import (
	"net/http"
	"strings"

	"auth-go-skd/token"
)

type Middleware struct {
	service *Service
}

func (m *Middleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Get Token from Cookie or Header
		var tokenStr string

		// Check Cookie
		cookie, err := r.Cookie("JWT")
		if err == nil {
			tokenStr = cookie.Value
		}

		// Check Header if no cookie
		if tokenStr == "" {
			reqToken := r.Header.Get("Authorization")
			splitToken := strings.Split(reqToken, "Bearer ")
			if len(splitToken) == 2 {
				tokenStr = splitToken[1]
			}
		}

		if tokenStr == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// 2. Parse & Validate Token
		claims, err := m.service.ParseToken(tokenStr)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if claims.User == nil {
			http.Error(w, "Unauthorized (No User)", http.StatusUnauthorized)
			return
		}

		// 3. Set User in Context
		r = token.SetUserInfo(r, *claims.User)

		next.ServeHTTP(w, r)
	})
}
