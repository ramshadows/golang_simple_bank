package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"simple_bank/token"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func addAuthorization(
	t *testing.T,
	request *http.Request,
	tokenMaker token.Maker,
	authorizationType string,
	username string,
	duration time.Duration,
) {
	token, payload, err := tokenMaker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, token)
	request.Header.Set(authorizationHeaderKey, authorizationHeader)
}

func TestAuthMiddleware(t *testing.T) {
	// define several test cases using an anonymous struct
	testCases := []struct {
		name          string
		setUpAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		// testcase 1: Happy case
		{
			name: "OK",
			setUpAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationBearerType, "user", time.Minute)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

			},
		},

		// testcase 2: No authorization header provided
		{
			name: "NoAuthorization",
			setUpAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				// empty

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)

			},
		},

		// testcase 3: Unsupported authorization type - required - Bearer token type
		{
			name: "UnSupportedAuthorization",
			setUpAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, "unsupported", "user", time.Minute)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)

			},
		},

		// testcase 4: Invalid Authorization format
		{
			name: "InvalidAuthorizationFormat",
			setUpAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, "", "user", time.Minute)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)

			},
		},

		// testcase 5: Expired token
		{
			name: "ExpiredToken",
			setUpAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationBearerType, "user", -time.Minute)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)

			},
		},
	}

	// loop thro the test cases
	for i := range testCases {
		tc := testCases[i]

		// generate the sub-test
		t.Run(
			tc.name,
			func(t *testing.T) {
				// create a new test server
				// notice for parameter db.store we use nil since we don't need to access the db store
				server := newTestServer(t, nil)

				// create a fake api route
				authPath := "/auth"
				server.router.GET(
					authPath,
					authMiddleware(server.tokenMaker),
					func(ctx *gin.Context) {
						// send a status ok with and empty body
						ctx.JSON(http.StatusOK, gin.H{})
					},
				)

				// now send a request to this server by creating a new recorder
				recorder := httptest.NewRecorder()

				// create a new request with method Get, the route and a nil body
				request, err := http.NewRequest(http.MethodGet, authPath, nil)
				require.NoError(t, err)

				// call the setUpAuth func to add the authorization header to the request
				tc.setUpAuth(t, request, server.tokenMaker)

				// call the server.route to serve our http with the response recorder and our request
				server.router.ServeHTTP(recorder, request)

				// call the checkResponse func to verify the results
				tc.checkResponse(t, recorder)

			},
		)

	}
}
