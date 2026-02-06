package auth

import (
	"log"
	"time"

	"auth-go-skd/provider"
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
		opts.TokenDuration = time.Minute * 5
	}
	// ... defaults ...
	return &Service{
		opts:      opts,
		providers: make(map[string]provider.Provider),
		logger:    log.Default(),
	}
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
