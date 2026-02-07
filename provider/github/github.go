package github

import (
	"auth-go-skd/token"
	"context"
	"encoding/json"
	"fmt"

	"golang.org/x/oauth2"
	githubOAuth "golang.org/x/oauth2/github"
)

type Provider struct {
	Config *oauth2.Config
}

func New(clientID, clientSecret, callbackURL string) *Provider {
	return &Provider{
		Config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  callbackURL,
			Scopes:       []string{"read:user", "user:email"},
			Endpoint:     githubOAuth.Endpoint,
		},
	}
}

func (p *Provider) Name() string {
	return "github"
}

func (p *Provider) GetAuthURL(state string) string {
	return p.Config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (p *Provider) FetchUser(ctx context.Context, code string) (token.User, error) {
	tok, err := p.Config.Exchange(ctx, code)
	if err != nil {
		return token.User{}, fmt.Errorf("failed to exchange token: %w", err)
	}

	client := p.Config.Client(ctx, tok)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return token.User{}, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	var userInfo struct {
		ID        int    `json:"id"`
		Login     string `json:"login"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		AvatarURL string `json:"avatar_url"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return token.User{}, fmt.Errorf("failed to decode user info: %w", err)
	}

	if userInfo.Email == "" {

		emailResp, err := client.Get("https://api.github.com/user/emails")
		if err == nil {
			defer emailResp.Body.Close()
			var emails []struct {
				Email    string `json:"email"`
				Primary  bool   `json:"primary"`
				Verified bool   `json:"verified"`
			}
			if json.NewDecoder(emailResp.Body).Decode(&emails) == nil {
				for _, email := range emails {
					if email.Primary && email.Verified {
						userInfo.Email = email.Email
						break
					}
				}
			}
		}
	}

	return token.User{
		ID:      fmt.Sprintf("%d", userInfo.ID),
		Name:    userInfo.Name,
		Email:   userInfo.Email,
		Picture: userInfo.AvatarURL,
		Attributes: map[string]interface{}{
			"username": userInfo.Login,
			"provider": "github",
		},
	}, nil
}
