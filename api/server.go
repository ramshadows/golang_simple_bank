package api

import (
	"fmt"
	db "simple_bank/db/sqlc"
	"simple_bank/token"
	"simple_bank/utils"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server serves all http requests for our banking service
type Server struct {
	config     utils.Config
	store      db.Store    // Allows has to interact with our db
	tokenMaker token.Maker // allows us to interact with our token maker interface
	router     *gin.Engine // Helps send http requests to the correct handler for processing

}

// NewServer creates a new http server instance and setups all api routing for our service
// on that server
func NewServer(config utils.Config, store db.Store) (*Server, error) {
	// create a new token maker object using the paseto token maker
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmeticKey)

	if err != nil {
		// return a nil object and and error description
		// notice %w - used to wrap the original error
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	// Creates a new server
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	// register the currencyValidator() with gin
	// call binding.Validator.Engine to find what type of validator gin is using
	// then convert the validator to validator.Validate
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setUpRouter()

	// return the server
	return server, nil
}

func (server *Server) setUpRouter() {
	// Creates a new router
	router := gin.Default()

	// Add routes to the router
	// create a POST route
	// pass a path /accounts in our case
	// Note: last param is the route handler
	// all other middle params are middleware
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)
	router.POST("/tokens/renew_access", server.renewAccessToken)

	// below routes need to be authorized
	// therefore we add our  middleware here
	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	// now instead of router, we use the authRoutes

	authRoutes.POST("/accounts", server.createAcccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.GET("/accounts", server.listAccounts)
	authRoutes.POST("/transfers", server.createTransfer)

	// Set this router object to server.router
	server.router = router

}

// Start runs an http request on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
