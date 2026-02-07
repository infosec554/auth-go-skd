package auth

import (
	"time"

	"auth-go-skd/avatar"
	"auth-go-skd/token"
)

type Opts struct {
	SecretReader   token.SecretFunc
	Secret         string
	TokenDuration  time.Duration
	CookieDuration time.Duration
	Issuer         string
	URL            string
	URLIsHTTPS     bool
	AvatarStore    avatar.Store
	Validator      token.Validator
	DisableXSRF    bool
}
