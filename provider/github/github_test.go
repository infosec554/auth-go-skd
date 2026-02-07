package github

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"golang.org/x/oauth2"
)

type MockTransport struct {
	RoundTripFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.RoundTripFunc(req)
}

func TestGithubProvider(t *testing.T) {
	p := New("id", "secret", "url")
	if p.Name() != "github" {
		t.Errorf("expected name 'github', got %s", p.Name())
	}
}

func TestFetchUser(t *testing.T) {
	p := New("id", "secret", "url")

	p.Config.Endpoint.TokenURL = "https://github.com/login/oauth/access_token"

	mockTransport := &MockTransport{
		RoundTripFunc: func(req *http.Request) (*http.Response, error) {

			if req.URL.String() == p.Config.Endpoint.TokenURL {
				respBody := map[string]interface{}{
					"access_token": "mock-github-token",
					"token_type":   "Bearer",
					"scope":        "read:user,user:email",
				}
				body, _ := json.Marshal(respBody)
				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(bytes.NewReader(body)),
					Header:     make(http.Header),
				}, nil
			}

			if req.URL.String() == "https://api.github.com/user" {
				userInfo := map[string]interface{}{
					"id":         98765,
					"login":      "githubuser",
					"name":       "Github User",
					"avatar_url": "http://github.com/avatar.jpg",

					"email": "",
				}
				body, _ := json.Marshal(userInfo)
				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(bytes.NewReader(body)),
					Header:     make(http.Header),
				}, nil
			}

			if req.URL.String() == "https://api.github.com/user/emails" {
				emails := []map[string]interface{}{
					{"email": "secondary@example.com", "primary": false, "verified": true},
					{"email": "primary@example.com", "primary": true, "verified": true},
				}
				body, _ := json.Marshal(emails)
				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(bytes.NewReader(body)),
					Header:     make(http.Header),
				}, nil
			}

			return &http.Response{
				StatusCode: 404,
				Body:       io.NopCloser(bytes.NewReader([]byte("not found: " + req.URL.String()))),
			}, nil
		},
	}

	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, &http.Client{Transport: mockTransport})

	user, err := p.FetchUser(ctx, "mock-code")
	if err != nil {
		t.Fatalf("FetchUser failed: %v", err)
	}

	if user.ID != "98765" {
		t.Errorf("expected user ID '98765', got %s", user.ID)
	}
	if user.Email != "primary@example.com" {
		t.Errorf("expected email 'primary@example.com', got %s", user.Email)
	}
	if user.Attributes["username"] != "githubuser" {
		t.Errorf("expected username 'githubuser', got %v", user.Attributes["username"])
	}
}
