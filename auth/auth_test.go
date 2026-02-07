package auth

import (
	"auth-go-skd/token"
	"testing"
	"time"
)

func TestService_Token(t *testing.T) {

	opts := Opts{
		Secret:        "test-secret-key-12345",
		TokenDuration: time.Minute * 15,
		URL:           "http://localhost",
	}
	s := New(opts)

	user := token.User{
		ID:    "user-123",
		Name:  "Test User",
		Email: "test@example.com",
	}

	tokenStr, err := s.Token(user)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}
	if tokenStr == "" {
		t.Fatal("token is empty")
	}

	claims, err := s.ParseToken(tokenStr)
	if err != nil {
		t.Fatalf("failed to parse token: %v", err)
	}

	if claims.User.ID != user.ID {
		t.Errorf("expected user ID %s, got %s", user.ID, claims.User.ID)
	}
	if claims.User.Email != user.Email {
		t.Errorf("expected email %s, got %s", user.Email, claims.User.Email)
	}

	if claims.Issuer != opts.Issuer {
	}
}

func TestGenerateState(t *testing.T) {
	state1 := generateState()
	state2 := generateState()

	if len(state1) == 0 {
		t.Error("generated state is empty")
	}

	if state1 == state2 {
		t.Error("generateState returned duplicate values (randomness failure)")
	}

	if len(state1) < 40 {
		t.Errorf("state too short: %d", len(state1))
	}
}

func TestMiddleware_Auth_Fail(t *testing.T) {

}
