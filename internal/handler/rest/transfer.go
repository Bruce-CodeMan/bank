// Package rest provides REST API handlers for bank application.
// It contains controllers that handles HTTP requests and delegate
// business logic to the service layer.
package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/BruceCompiler/bank/internal/dto"
	"github.com/BruceCompiler/bank/internal/service"
)

// TransferController handles HTTP requests related to transfer operations.
// It depends on TransferService for business logic execution.
type TransferController struct {
	transferService *service.TransferService
}

// NewTransferController creates a new TransferController with the given service.
func NewTransferController(transferService *service.TransferService) *TransferController {
	return &TransferController{transferService: transferService}
}

// CreateTransfer hanles POST request to create a new transfer
// It expects a JSON body with from_account_id, to_account_id, amount, currency
//
// @Summary 		Create a new transfer
// @Description		Create a new transfer with the given details
// @Tags			transfer
// @Accept			json
// @Produce			json
// @Param			transfer	body	dto.CreateTransferRequest	true	"Transfer data"
// @Success			200		{object}
// @Failure			400		{object}	gin.H						"Invalid request"
// @Failure			500		{object}	gin.H						"Internal server error"
// @Router			/api/v1/transfer 	[post]
func (tc *TransferController) CreateTransfer(ctx *gin.Context) {
	var req dto.CreateTransferRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := tc.transferService.CreateTransfer(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, result)
}
