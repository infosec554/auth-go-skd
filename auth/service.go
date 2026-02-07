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

func User(r *http.Request) token.User {
	return token.MustGetUserInfo(r)
}

type Service struct {
	opts      Opts
	providers map[string]provider.Provider
	logger    *log.Logger
}

func New(opts Opts) *Service {

	if opts.TokenDuration == 0 {
		opts.TokenDuration = time.Minute * 15
	}
	if opts.CookieDuration == 0 {
		opts.CookieDuration = time.Hour * 24 * 7
	}
	if opts.AvatarStore == nil {
		opts.AvatarStore = avatar.NewLocalFS("/tmp/avatars")
	}

	return &Service{
		opts:      opts,
		providers: make(map[string]provider.Provider),
		logger:    log.Default(),
	}
}

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

func (s *Service) Add(p provider.Provider) {
	s.providers[p.Name()] = p
}

func (s *Service) ParseToken(tokenStr string) (*token.Claims, error) {
	t, err := jwt.ParseWithClaims(tokenStr, &token.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		if s.opts.SecretReader != nil {
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

func (s *Service) AddProvider(name, cid, csecret string) {

	log.Printf("AddProvider(%s) called - please use service.Add(provider) instead for now", name)
}

func (s *Service) Middleware() *Middleware {
	return &Middleware{service: s}
}
