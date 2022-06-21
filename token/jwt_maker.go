package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

const minSecreteKeySize = 32

// JWMaker is JSON Web Token Maker - JWT which implements the JWT interface
type JWTMaker struct {
	secretKey string
}

// NewJWTMaker creates a new JWT token
// must implement the JWT interface
func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecreteKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecreteKeySize)

	}

	return &JWTMaker{secretKey: secretKey}, nil

}

// CreateToken creates and signs a new token for a specific username and a valid duration
func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, *Payload, error) {
	// create a new payload by calling NewPayload fuction
	payload, err := NewPayload(username, duration)

	if err != nil {
		return "", payload, err
	}

	// if we got here, we generate a new jwtToken
	// by calling the jwt.NewWithClaims that take a signing method and a payload
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := jwtToken.SignedString([]byte(maker.secretKey))
	return token, payload, err

}

// VerifyToken takes a token to verify and returns a Payload stored inside the body of the token
// VerifyToken checks if the token is valid or not
func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}
