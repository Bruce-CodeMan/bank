package api

import (
	"net/http"

	db "github.com/BruceCompiler/bank/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type createAccountRequest struct {
	PublicID uuid.UUID `json:"public_id" binding:"required"`
	Owner    string    `json:"owner" binding:"required"`
	Currency string    `json:"currency" binding:"required,oneof=USD EUR"`
}

func (s *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// 把google uuid.UUID转成pgtype.UUID
	var pgUUID pgtype.UUID
	pgUUID.Bytes = req.PublicID
	pgUUID.Valid = true

	arg := db.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
		PublicID: pgUUID,
		Balance:  0,
	}

	account, err := s.store.CreateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	ctx.JSON(http.StatusOK, account)
}
