// Package server provides HTTP server setup and configuration.
package server

import (
	"github.com/gin-gonic/gin"

	"github.com/BruceCompiler/bank/internal/controller/rest"
	"github.com/BruceCompiler/bank/internal/repository/postgres"
	"github.com/BruceCompiler/bank/internal/service"
)

// HTTPServer encapsulates the HTTP server and its dependencies.
// It manages the Gin engine and database store.
type HTTPServer struct {
	store  postgres.Store
	engine *gin.Engine
}

// NewHTTPServer creates a new HTTPServer with the given store.
// It initializes the Gin engine and sets up all routes.
func NewHTTPServer(store postgres.Store) *HTTPServer {
	server := &HTTPServer{
		store:  store,
		engine: gin.Default(),
	}
	server.setupRoutes()
	return server
}

// setupRoutes configures all API routes.
// It creates services, controllers, and registers them with the router.
func (s *HTTPServer) setupRoutes() {
	accountService := service.NewAccountService(s.store)
	accountController := rest.NewAccountController(accountService)
	rest.RegisterRoutes(s.engine, accountController)

}

// Start begins listening for HTTP requests on the specified address.
// It blocks until the server is shut down or encounters an error
func (s *HTTPServer) Start(address string) error {
	return s.engine.Run(address)
}

func (s *HTTPServer) Router() *gin.Engine {
	return s.engine
}
