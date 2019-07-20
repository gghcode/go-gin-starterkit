package user

import (
	"net/http"
	"strconv"

	"github.com/gghcode/go-gin-starterkit/api/common"
	"github.com/gghcode/go-gin-starterkit/middleware"
	"github.com/gghcode/go-gin-starterkit/service"
	"github.com/gin-gonic/gin"
)

// APIPath is path prefix
const APIPath = "/users/"

// Controller is user controller
type Controller struct {
	repo     Repository
	passport service.Passport
}

// NewController return new user controller instance.
func NewController(repo Repository, passport service.Passport) *Controller {
	return &Controller{
		repo:     repo,
		passport: passport,
	}
}

// RegisterRoutes register handler routes.
func (controller *Controller) RegisterRoutes(router gin.IRouter) {
	userRouter := router.Group(APIPath)
	{
		userRouter.Handle("POST", "/", controller.createUser)

		authorized := userRouter.Use(middleware.AuthRequired())
		{
			authorized.Handle("GET", "/:username", controller.getUserByUserName)
			authorized.Handle("PUT", "/:id", controller.updateUserByID)
			authorized.Handle("DELETE", "/:id", controller.removeUserByID)
		}
	}
}

// @Description Create new user
// @Accept json
// @Produce json
// @Param payload body user.CreateUserRequest true "user payload"
// @Success 201 {object} user.UserResponse "ok"
// @Failure 400 {object} common.ErrorResponse "Invalid user payload"
// @Tags User API
// @Router /users [post]
func (controller *Controller) createUser(ctx *gin.Context) {
	var dtoReq CreateUserRequest

	if err := ctx.ShouldBindJSON(&dtoReq); err != nil {
		ctx.JSON(http.StatusBadRequest, common.NewErrResp(err))
		return
	}

	passwordHash, err := controller.passport.HashPassword(dtoReq.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.NewErrResp(err))
		return
	}

	userEntity := User{
		UserName:     dtoReq.UserName,
		PasswordHash: passwordHash,
	}

	createdUser, err := controller.repo.CreateUser(userEntity)
	if err == common.ErrAlreadyExistsEntity {
		ctx.JSON(http.StatusConflict, common.NewErrResp(err))
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.NewErrResp(err))
		return
	}

	ctx.JSON(http.StatusCreated, createdUser.Response())
}

// @Description Get user by username
// @Security ApiKeyAuth
// @Produce json
// @Param username path string true "User Name"
// @Success 200 {object} user.UserResponse "ok"
// @Failure 404 {object} common.ErrorResponse "Not found entity"
// @Tags User API
// @Router /users/{username} [get]
func (controller *Controller) getUserByUserName(ctx *gin.Context) {
	userName := ctx.Param("username")

	user, err := controller.repo.GetUserByUserName(userName)
	if err == common.ErrEntityNotFound {
		ctx.JSON(http.StatusNotFound, common.NewErrResp(err))
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.NewErrResp(err))
		return
	}

	ctx.JSON(http.StatusOK, user.Response())
}

// @Description Update new user by user id
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path string true "user id"
// @Param payload body user.UpdateUserRequest true "user payload"
// @Success 200 {object} user.UserResponse "ok"
// @Failure 400 {object} common.ErrorResponse "Invalid user payload"
// @Failure 404 {object} common.ErrorResponse "Not found entity"
// @Tags User API
// @Router /users/{id} [put]
func (controller *Controller) updateUserByID(ctx *gin.Context) {
	userID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, common.NewErrResp(common.ErrParsingFailed))
		return
	}

	var reqBody UpdateUserRequest
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	entity := User{
		UserName: reqBody.UserName,
	}

	user, err := controller.repo.UpdateUserByUserID(userID, entity)
	if err == common.ErrEntityNotFound {
		ctx.JSON(http.StatusNotFound, common.NewErrResp(err))
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.NewErrResp(err))
		return
	}

	ctx.JSON(http.StatusOK, user.Response())
}

// @Description Remove user by user id
// @Security ApiKeyAuth
// @Produce json
// @Param id path string true "user id"
// @Success 200 {object} user.UserResponse "ok"
// @Failure 404 {object} common.ErrorResponse "Not found entity"
// @Tags User API
// @Router /users/{id} [delete]
func (controller *Controller) removeUserByID(ctx *gin.Context) {
	userID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, common.NewErrResp(common.ErrParsingFailed))
		return
	}

	removedUser, err := controller.repo.RemoveUserByUserID(userID)
	if err == common.ErrEntityNotFound {
		ctx.JSON(http.StatusNotFound, common.NewErrResp(err))
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.NewErrResp(err))
		return
	}

	ctx.JSON(http.StatusOK, removedUser.Response())
}
