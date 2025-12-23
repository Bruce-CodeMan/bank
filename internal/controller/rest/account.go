// Package rest provides REST API handlers for the bank application.
// It contains controllers taht handle HTTP requests and delegate
// business logic to the service layer.
package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/BruceCompiler/bank/internal/dto"
	"github.com/BruceCompiler/bank/internal/service"
)

// AccountController handles HTTP requests related to account operations.
// It depends on AccountService for business logic execution.
type AccountController struct {
	accountService *service.AccountService
}

// NewAccountController creates a new AccountController with the given service.
func NewAccountController(accountService *service.AccountService) *AccountController {
	return &AccountController{accountService: accountService}
}

// CreateAccount handles POST requests to create a new bank account.
// It expects a JSON body with owner, currency, and public_id fields
//
// Request body:
//   - public_id: UUID(required)
//   - owner: string(required)
//   - currency: string(required, must be USD or EUR)
//
// Responses:
//   - 200: Account created successfully
//   - 400: Invalid request body
//   - 500: Internal server error
func (ac *AccountController) CreateAccount(ctx *gin.Context) {
	var req dto.CreateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	account, err := ac.accountService.CreateAccount(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, account)
}

// GetAccount handles a GET request to retrieve a bank account by its UUID.
//
// @Summary 		Get account by UUID
// @Description		Retrieves account details based on the provided public_id in the URL path.
// @Tags			account
// @Param			public_id	path	string true "Account UUID"
// Success			200			{object}	db.Account
// Failure			500			{object}	gin.H
// Router			/api/v1/account/{public_id} [get]
func (ac *AccountController) GetAccount(ctx *gin.Context) {
	var req dto.GetAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	account, err := ac.accountService.GetAccountByPublicID(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, account)
}
