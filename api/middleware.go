package api

import (
	"errors"
	"fmt"
	"net/http"
	"simple_bank/token"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationBearerType = "bearer" // allows only for Bearer type
	authorizationPayloadKey = "authorization_payload"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	// return an anonymous function - the authetication middleware func that we must implement
	return func(ctx *gin.Context) {
		// first extract the authorization header
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)

		// if header is empty
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")

			// abort the request and send a json response to the client
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))

			// return immeddiately
			return
		}

		// if the authorization header is provided
		// split the authorization header by space
		fields := strings.Fields(authorizationHeader)

		// if an error occurs
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")

			// abort the request and send a json response to the client
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))

			// return immeddiately
			return

		}

		authorizationType := strings.ToLower(fields[0])

		// match the bearer type to the supported bearer type
		if authorizationType != authorizationBearerType {
			err := fmt.Errorf("unsupported authorization %s type", authorizationType)

			// abort the request and send a json response to the client
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))

			// return immeddiately
			return

		}

		// if we got here, we now have our accessToken
		accessToken := fields[1]

		// verify our token
		payload, err := tokenMaker.VerifyToken(accessToken)

		if err != nil {
			// abort the request and send a json response to the client
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))

			// return immeddiately
			return

		}

		// store the authorization payload to the context by passing a key:value pair
		ctx.Set(authorizationPayloadKey, payload)

		// finally - forward the request to next handler
		ctx.Next()

	}
}
