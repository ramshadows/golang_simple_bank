package token

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

// Pasetomaker is PASETO token maker
type Pasetomaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

// NewPasetoMaker creates a new paseto maker instance
func NewPasetoMaker(symmetricKey string) (Maker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		// return a nil object and an error
		return nil, fmt.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize)
	}
	//otherwise, create a new paseto object
	maker := &Pasetomaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey), // convert the input string to byte slice

	}

	// note: PasetoMaker receiver maker must implement the Maker interface as below
	return maker, nil

}

// CreateToken creates and signs a new token for a specific username and a valid duration
func (maker *Pasetomaker) CreateToken(username string, duration time.Duration) (string, *Payload, error) {
	// create a new payload by calling NewPayload fuction
	payload, err := NewPayload(username, duration)

	if err != nil {
		// return an empty string and an error
		return "", payload,  err
	}

	// otherwise, we return maker encrypted
	token, err := maker.paseto.Encrypt(maker.symmetricKey, payload, nil)

	return token, payload, err

}

// VerifyToken takes a token to verify and returns a Payload stored inside the body of the token
func (maker *Pasetomaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil)

	if err != nil {
		// return a nil object and errorInvalidToken object
		return nil, ErrInvalidToken
	}

	// otherwise if there is an error
	err = payload.Valid()
	if err != nil {
		// return nil and an error
		return nil, err
	}

	return payload, nil

}
