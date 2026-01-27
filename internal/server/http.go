// Package server provides HTTP server setup and configuration.
package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	db "github.com/BruceCompiler/bank/db/sqlc"
	"github.com/BruceCompiler/bank/internal/handler/rest"
	"github.com/BruceCompiler/bank/internal/service"
	"github.com/BruceCompiler/bank/internal/token"
	"github.com/BruceCompiler/bank/internal/validators"
	"github.com/BruceCompiler/bank/utils"
	"github.com/BruceCompiler/bank/worker"
)

// HTTPServer encapsulates the HTTP server and its dependencies.
// It manages the Gin engine and database store.
type HTTPServer struct {
	config          utils.Config
	store           db.Store
	tokenMaker      token.Maker
	engine          *gin.Engine
	taskDistributor worker.TaskDistributor
}

// NewHTTPServer creates a new HTTPServer with the given store.
// It initializes the Gin engine and sets up all routes.
func NewHTTPServer(config utils.Config, store db.Store, taskDistributor worker.TaskDistributor) (*HTTPServer, error) {
	tokenMaker, err := token.NewJWTMaker(config.TokenSynmmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &HTTPServer{
		config:          config,
		tokenMaker:      tokenMaker,
		store:           store,
		engine:          gin.Default(),
		taskDistributor: taskDistributor,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validators.ValidCurrency)
	}

	server.setupRoutes()
	return server, nil
}

// setupRoutes configures all API routes.
// It creates services, controllers, and registers them with the router.
func (s *HTTPServer) setupRoutes() {
	accountService := service.NewAccountService(s.store)
	accountController := rest.NewAccountController(accountService)
	transferService := service.NewTransferService(s.store)
	transferController := rest.NewTransferController(transferService)
	userService := service.NewUserService(s.store, s.tokenMaker, s.config, s.taskDistributor)
	userController := rest.NewUserController(userService)
	tokenController := rest.NewTokenController(s.store, s.tokenMaker, s.config)
	rest.RegisterRoutes(s.engine, s.tokenMaker, accountController, transferController, userController, tokenController)
}

// Start begins listening for HTTP requests on the specified address.
// It blocks until the server is shut down or encounters an error
func (s *HTTPServer) Start(address string) error {
	return s.engine.Run(address)
}

func (s *HTTPServer) Router() *gin.Engine {
	return s.engine
}

func (s *HTTPServer) TokenMaker() token.Maker {
	return s.tokenMaker
}
