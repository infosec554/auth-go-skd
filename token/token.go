package token

import (
	"context"
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const userKey contextKey = "user"

type SecretFunc func(id string) (string, error)

type User struct {
	Name       string                 `json:"name"`
	ID         string                 `json:"id"`
	Picture    string                 `json:"picture"`
	IP         string                 `json:"ip,omitempty"`
	Email      string                 `json:"email,omitempty"`
	Attributes map[string]interface{} `json:"attrs,omitempty"`
}

type Claims struct {
	User *User `json:"user,omitempty"`
	jwt.RegisteredClaims
}

type ValidatorFunc func(token string, claims Claims) bool

type Validator interface {
	Validate(token string, claims Claims) bool
}

func (f ValidatorFunc) Validate(token string, claims Claims) bool {
	return f(token, claims)
}

func SetUserInfo(r *http.Request, user User) *http.Request {
	ctx := context.WithValue(r.Context(), userKey, user)
	return r.WithContext(ctx)
}

func GetUserInfo(r *http.Request) (User, error) {
	if user, ok := r.Context().Value(userKey).(User); ok {
		return user, nil
	}
	return User{}, errors.New("user info not found in context")
}

func MustGetUserInfo(r *http.Request) User {
	user, err := GetUserInfo(r)
	if err != nil {
		panic(err)
	}
	return user
}
