package token

import (
	"time"
)

// Maker is an interface for managing tokens
type Maker interface {
	// CreateToken creates and signs a new token for a specific username and a valid duration
	CreateToken(username string, duration time.Duration) (string, *Payload, error)

	// VerifyToken takes a token to verify and returns a Payload stored inside the body of the token
	VerifyToken(token string) (*Payload, error)
}
