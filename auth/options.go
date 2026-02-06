package auth

import (
	"time"

	"auth-go-skd/avatar"
	"auth-go-skd/token"
)

// Opts configuration options for the auth service
type Opts struct {
	SecretReader   token.SecretFunc
	Secret         string // Simplified single secret
	TokenDuration  time.Duration
	CookieDuration time.Duration
	Issuer         string
	URL            string
	URLIsHTTPS     bool // Enforce Secure cookies manually if not detected
	AvatarStore    avatar.Store
	Validator      token.Validator
	DisableXSRF    bool
}
