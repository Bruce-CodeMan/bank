// Package rest provides REST API routing and handlers
package rest

import (
	"github.com/BruceCompiler/bank/internal/token"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes sets up all API routes for the application.
func RegisterRoutes(router *gin.Engine,
	tokenMaker token.Maker,
	ac *AccountController,
	tc *TransferController,
	uc *UserController,
) {
	api := router.Group("/api/v1")
	// User
	user := api.Group("/user")
	{
		user.POST("", uc.CreateUser)
		user.POST("/login", uc.Login)
	}

	// Account
	account := api.Group("/account")
	{
		account.POST("", ac.CreateAccount)
		account.GET("/:public_id", ac.GetAccount)
		account.GET("", ac.ListAccount)
	}

	// Transfer
	transfer := api.Group("/transfer")
	{
		transfer.POST("", tc.CreateTransfer)
	}

}
