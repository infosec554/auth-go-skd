package providers

import (
	"context"
)

type ProviderInfo struct {
	ID        string
	Email     string
	Name      string
	AvatarURL string
}

type Provider interface {
	GetAuthURL(state string) string
	FetchUser(ctx context.Context, code string) (*ProviderInfo, error)
}
