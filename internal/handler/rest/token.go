package rest

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	db "github.com/BruceCompiler/bank/db/sqlc"
	"github.com/BruceCompiler/bank/internal/token"
	"github.com/BruceCompiler/bank/utils"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

type renewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type renewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

type TokenController struct {
	store      db.Store
	tokenMaker token.Maker
	config     utils.Config
}

func NewTokenController(s db.Store, t token.Maker, c utils.Config) *TokenController {
	return &TokenController{
		store:      s,
		tokenMaker: t,
		config:     c,
	}
}

func (t *TokenController) RenewAccessToken(ctx *gin.Context) {
	var req renewAccessTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "refresh_token is invalid"})
		return
	}

	refreshPayload, err := t.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err})
		return
	}

	session, err := t.store.GetSession(ctx, pgtype.UUID{
		Bytes: refreshPayload.ID,
		Valid: true,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "cannot find refresh token"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	if session.IsBlocked {
		err := fmt.Errorf("blocked session")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err})
		return
	}

	if session.Username != refreshPayload.Username {
		err := fmt.Errorf("incorrect session user")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err})
		return
	}

	if session.RefreshToken != req.RefreshToken {
		err := fmt.Errorf("mismatched session token")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err})
		return
	}

	if time.Now().After(session.ExpiresAt.Time) {
		err := fmt.Errorf("expired session")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err})
		return
	}

	accessToken, accessPayload, err := t.tokenMaker.CreateToken(
		refreshPayload.Username,
		t.config.RefreshTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	rsp := renewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiresAt.Time,
	}

	ctx.JSON(http.StatusOK, rsp)

}
