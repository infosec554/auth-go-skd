package provider

import (
	"auth-go-skd/token"
	"context"
)

// Provider defines the interface for an authentication provider
type Provider interface {
	Name() string
	GetAuthURL(state string) string
	FetchUser(ctx context.Context, code string) (token.User, error)
}
