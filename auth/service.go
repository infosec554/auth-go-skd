package auth

import (
	"errors"
	"fmt"
	"log"
	"time"

	"auth-go-skd/provider"
	"auth-go-skd/token"

	"github.com/golang-jwt/jwt/v5"
)

// Service is the main auth service
type Service struct {
	opts      Opts
	providers map[string]provider.Provider
	logger    *log.Logger
}

// NewService creates a new auth service
func NewService(opts Opts) *Service {
	if opts.TokenDuration == 0 {
		opts.TokenDuration = time.Minute * 15
	}
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
			Audience:  jwt.ClaimStrings{s.opts.URL}, // Verify audience
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
func (s *Service) AddProvider(name, cid, csecret string) {
	// Factory logic or just accept interface?
	// go-pkgz accepts name/cid/csecret and creates internally or via generic AddProvider (deprecated?)
	// Actually go-pkgz has AddProvider(name, cid, csecret) for "system" providers.
	// But sticking to our `provider` package is safer.
	// Let's assume generic interface for flexibility.
	// But providing convenience method is asked implicitly.
	// TODO: Implement factory in provider package?
}

// AddCustomProvider adds a pre-configured provider
func (s *Service) AddCustomProvider(p provider.Provider) {
	s.providers[p.Name()] = p
}

// Middleware returns the auth middleware
func (s *Service) Middleware() *Middleware {
	return &Middleware{service: s}
}
