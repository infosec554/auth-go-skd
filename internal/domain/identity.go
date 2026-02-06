package domain

import (
	"time"
)

type Identity struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	Provider   string    `json:"provider"` // google, github, linkedin
	ProviderID string    `json:"provider_id"`
	CreatedAt  time.Time `json:"created_at"`
	LastLogin  time.Time `json:"last_login"`
}
