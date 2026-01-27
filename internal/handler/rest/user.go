package rest

import (
	"net/http"

	"github.com/BruceCompiler/bank/internal/dto"
	"github.com/BruceCompiler/bank/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
)

// UserController handles the HTTP request related to user operations.
// It depends on UserService for business logic executions.
type UserController struct {
	userService *service.UserService
}

// NewUserController creates a new UserController with the provided service.
func NewUserController(u *service.UserService) *UserController {
	return &UserController{userService: u}
}

// CreateUser handles POST request to create a new user
// It expects a JSON body with username, password, full_name, email
//
// @Summary				Create a new user
// @Description			Handles a POST request to create a new user
// @Tags				user
// @Accept				json
// @Produce				json
// @Param				user  body	dto.CreateUserRequest 	true 	"user registration info"
// @Success				200	  {object}	dto.CreateUserResponse
// @Failure				400	  {object}  gin.H	"Bad Request"
// @Failure				500	  {object}	gin.H	"Internal Server Error"
// @Router				/api/v1/user	[post]
func (uc *UserController) CreateUser(ctx *gin.Context) {
	var req dto.CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := uc.userService.CreateUser(ctx.Request.Context(), req)
	if err != nil {
		if pqErr, ok := err.(*pgconn.PgError); ok {
			switch pqErr.Code {
			case "23505":
				ctx.JSON(http.StatusForbidden, gin.H{"error": "username or email already exists"})
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rsp := dto.CreateUserResponse{
		ID:       result.PublicID.String(),
		Username: result.Username,
		FullName: result.FullName,
		Email:    result.Email,
	}
	ctx.JSON(http.StatusOK, rsp)
}

func (uc *UserController) Login(ctx *gin.Context) {
	var req dto.LoginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := uc.userService.Login(ctx.Request.Context(), req)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			switch pgErr.Code {
			case "23505":
				ctx.JSON(http.StatusForbidden, gin.H{"error": "username or password is not correct"})
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}
