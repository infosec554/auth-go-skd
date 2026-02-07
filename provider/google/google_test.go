package google

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

func TestGoogleProvider(t *testing.T) {
	clientID := "test-client-id"
	clientSecret := "test-client-secret"
	callbackURL := "http://localhost/callback"

	p := New(clientID, clientSecret, callbackURL)

	if p.Name() != "google" {
		t.Errorf("expected name 'google', got %s", p.Name())
	}

	if p.Config.ClientID != clientID {
		t.Errorf("expected clientID %s, got %s", clientID, p.Config.ClientID)
	}
}

func TestFetchUser(t *testing.T) {
	p := New("id", "secret", "url")

	p.Config.Endpoint.TokenURL = "https://oauth2.googleapis.com/token"

	mockTransport := &MockTransport{
		RoundTripFunc: func(req *http.Request) (*http.Response, error) {

			if req.URL.String() == p.Config.Endpoint.TokenURL {
				respBody := map[string]interface{}{
					"access_token": "mock-access-token",
					"token_type":   "Bearer",
					"expires_in":   3600,
				}
				body, _ := json.Marshal(respBody)
				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(bytes.NewReader(body)),
					Header:     make(http.Header),
				}, nil
			}

			if req.URL.String() == "https://www.googleapis.com/oauth2/v2/userinfo" {
				userInfo := map[string]string{
					"id":      "12345",
					"email":   "test@example.com",
					"name":    "Test User",
					"picture": "http://example.com/avatar.jpg",
				}
				body, _ := json.Marshal(userInfo)
				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(bytes.NewReader(body)),
					Header:     make(http.Header),
				}, nil
			}

			return &http.Response{
				StatusCode: 404,
				Body:       io.NopCloser(bytes.NewReader([]byte("not found"))),
			}, nil
		},
	}

	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, &http.Client{Transport: mockTransport})

	user, err := p.FetchUser(ctx, "mock-code")
	if err != nil {
		t.Fatalf("FetchUser failed: %v", err)
	}

	if user.ID != "12345" {
		t.Errorf("expected user ID '12345', got %s", user.ID)
	}
	if user.Email != "test@example.com" {
		t.Errorf("expected email 'test@example.com', got %s", user.Email)
	}
	if user.Attributes["provider"] != "google" {
		t.Errorf("expected provider attribute 'google', got %v", user.Attributes["provider"])
	}
}
