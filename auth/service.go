package auth

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"auth-go-skd/avatar"
	"auth-go-skd/provider"
	"auth-go-skd/token"

	"github.com/golang-jwt/jwt/v5"
)

// User returns the user info from the request context
// This alias allows users to access user info via `auth.User(r)` instead of `token.MustGetUserInfo(r)`
func User(r *http.Request) token.User {
	return token.MustGetUserInfo(r)
}

// Service is the main auth service
type Service struct {
	opts      Opts
	providers map[string]provider.Provider
	logger    *log.Logger
}

// New creates a new auth service with simple options
func New(opts Opts) *Service {
	// Set defaults
	if opts.TokenDuration == 0 {
		opts.TokenDuration = time.Minute * 15
	}
	if opts.CookieDuration == 0 {
		opts.CookieDuration = time.Hour * 24 * 7
	}
	if opts.AvatarStore == nil {
		opts.AvatarStore = avatar.NewLocalFS("/tmp/avatars")
	}

	// Default simplified secret reader if string secret is used
	// We might want to add a `Secret string` field to Opts for simplicity

	return &Service{
		opts:      opts,
		providers: make(map[string]provider.Provider),
		logger:    log.Default(),
	}
}

// Token creates a new JWT token for the user
func (s *Service) Token(user token.User) (string, error) {
	claims := token.Claims{
		User: &user,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.opts.Issuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.opts.TokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Audience:  jwt.ClaimStrings{s.opts.URL},
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Get secret (default "secret" if not provided)
	secret := "secret"
	if s.opts.SecretReader != nil {
		var err error
		secret, err = s.opts.SecretReader(user.ID)
		if err != nil {
			return "", err
		}
	}

	return jwtToken.SignedString([]byte(secret))
}

// Add adds a provider to the service
func (s *Service) Add(p provider.Provider) {
	s.providers[p.Name()] = p
}

// ParseToken validates and parses the token
func (s *Service) ParseToken(tokenStr string) (*token.Claims, error) {
	t, err := jwt.ParseWithClaims(tokenStr, &token.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		if s.opts.SecretReader != nil {
			// For simplicity in this step, we use a generic secret (or handle empty ID)
			// Ideally we peek claims to find ID.
			secret, err := s.opts.SecretReader("")
			return []byte(secret), err
		}
		return []byte("secret"), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := t.Claims.(*token.Claims); ok && t.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// AddProvider adds a provider to the service
// AddProvider adds a provider to the service by name (requires manual factory registration or use Add for pre-configured providers)
func (s *Service) AddProvider(name, cid, csecret string) {
	// For now, we recommend using service.Add(google.New(...)) for type-safety and flexibility
	// But we can keep this stub if we plan to add a registry later.
	log.Printf("AddProvider(%s) called - please use service.Add(provider) instead for now", name)
}

// Middleware returns the auth middleware
func (s *Service) Middleware() *Middleware {
	return &Middleware{service: s}
}
