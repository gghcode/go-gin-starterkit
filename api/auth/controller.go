package auth

import (
	"net/http"

	"github.com/gghcode/go-gin-starterkit/api/common"
	"github.com/gghcode/go-gin-starterkit/api/user"
	"github.com/gghcode/go-gin-starterkit/config"
	"github.com/gghcode/go-gin-starterkit/db"
	"github.com/gghcode/go-gin-starterkit/service"
	"github.com/gin-gonic/gin"
)

// APIPath is path prefix
const APIPath = "/auth/"

// Controller is auth controller
type Controller struct {
	conf    config.JwtConfig
	service Service
}

// NewController return new auth controller instance.
func NewController(
	conf config.Configuration,
	userRepo user.Repository,
	passport service.Passport,
	redisConn db.RedisConn) *Controller {
	return &Controller{
		conf:    conf.Jwt,
		service: NewService(conf, userRepo, passport, redisConn),
	}
}

// RegisterRoutes register handler routes.
func (controller *Controller) RegisterRoutes(router gin.IRouter) {
	router.Handle("POST", APIPath+"/token", controller.getToken)
}

// @Description Get new access token
// @Accept json
// @Produce json
// @Param payload body auth.CreateAccessTokenRequest true "payload"
// @Success 200 {object} auth.TokenResponse "ok"
// @Failure 400 {object} common.ErrorResponse "Invalid payload"
// @Failure 401 {object} common.ErrorResponse "Invalid credential"
// @Tags Auth API
// @Router /auth/token [post]
func (controller *Controller) getToken(ctx *gin.Context) {
	var reqPayload CreateAccessTokenRequest

	if err := ctx.ShouldBindJSON(&reqPayload); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			common.NewErrResp(common.ErrInvalidRequestPayload),
		)
		return
	}

	loginUser, err := controller.service.VerifyAuthentication(
		reqPayload.UserName,
		reqPayload.Password,
	)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, common.NewErrResp(err))
		return
	}

	accessToken, _ := controller.service.GenerateAccessToken(loginUser.ID)
	refreshToken, _ := controller.service.IssueRefreshToken(loginUser.ID)

	res := TokenResponse{
		Type:         "Bearer",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    controller.conf.AccessExpiresInSec,
	}

	ctx.JSON(http.StatusOK, res)
}
