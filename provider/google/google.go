package google

import (
	"auth-go-skd/config"
	"auth-go-skd/token"
	"context"
	"encoding/json"
	"fmt"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Provider struct {
	config *oauth2.Config
}

func New(cfg config.Google) *Provider {
	return &Provider{
		config: &oauth2.Config{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			RedirectURL:  cfg.RedirectURL,
			Scopes:       []string{"openid", "email", "profile"},
			Endpoint:     google.Endpoint,
		},
	}
}

func (p *Provider) Name() string {
	return "google"
}

func (p *Provider) GetAuthURL(state string) string {
	return p.config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (p *Provider) FetchUser(ctx context.Context, code string) (token.User, error) {
	tok, err := p.config.Exchange(ctx, code)
	if err != nil {
		return token.User{}, fmt.Errorf("failed to exchange token: %w", err)
	}

	client := p.config.Client(ctx, tok)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return token.User{}, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	var userInfo struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return token.User{}, fmt.Errorf("failed to decode user info: %w", err)
	}

	return token.User{
		ID:      userInfo.ID,
		Name:    userInfo.Name,
		Email:   userInfo.Email,
		Picture: userInfo.Picture,
		Attributes: map[string]interface{}{
			"provider": "google",
		},
	}, nil
}
