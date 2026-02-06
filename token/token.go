package token

import (
	"context"
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const userKey contextKey = "user"

// SecretFunc returns the secret for the given ID
type SecretFunc func(id string) (string, error)

// User defines the user information provided by the auth providers
type User struct {
	Name       string                 `json:"name"`
	ID         string                 `json:"id"`
	Picture    string                 `json:"picture"`
	IP         string                 `json:"ip,omitempty"`
	Email      string                 `json:"email,omitempty"`
	Attributes map[string]interface{} `json:"attrs,omitempty"`
}

// Claims defines the JWT claims (using jwt/v5)
type Claims struct {
	User *User `json:"user,omitempty"` // Embedded user info
	jwt.RegisteredClaims
}

// ValidatorFunc validates the token claims
type ValidatorFunc func(token string, claims Claims) bool

// Validator interface
type Validator interface {
	Validate(token string, claims Claims) bool
}

// Validate implements Validator interface
func (f ValidatorFunc) Validate(token string, claims Claims) bool {
	return f(token, claims)
}

// SetUserInfo sets the user info in the request context
func SetUserInfo(r *http.Request, user User) *http.Request {
	ctx := context.WithValue(r.Context(), userKey, user)
	return r.WithContext(ctx)
}

// GetUserInfo retrieves the user info from the request context
func GetUserInfo(r *http.Request) (User, error) {
	if user, ok := r.Context().Value(userKey).(User); ok {
		return user, nil
	}
	return User{}, errors.New("user info not found in context")
}

// MustGetUserInfo retrieves user info or panics
func MustGetUserInfo(r *http.Request) User {
	user, err := GetUserInfo(r)
	if err != nil {
		panic(err)
	}
	return user
}
