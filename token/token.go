package token

import (
	"github.com/golang-jwt/jwt/v5"
)

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
