package auth

import (
	"net/http"

	"github.com/gghcode/go-gin-starterkit/app/api/common"
	"github.com/gghcode/go-gin-starterkit/app/api/user"
	"github.com/gghcode/go-gin-starterkit/config"
	"github.com/gghcode/go-gin-starterkit/services"
	"github.com/gin-gonic/gin"
)

// Controller is auth controller
type Controller struct {
	conf    config.JwtConfig
	service Service
}

// NewController return new auth controller instance.
func NewController(conf config.Configuration, userRepo user.Repository, passport services.Passport) *Controller {
	return &Controller{
		conf:    conf.Jwt,
		service: NewService(conf, userRepo, passport),
	}
}

// RegisterRoutes register handler routes.
func (controller *Controller) RegisterRoutes(router gin.IRouter) {
	router.Handle("POST", "/token", controller.getToken)
}

// @Description Get new access token
// @Accept json
// @Produce json
// @Param payload body auth.CreateAccessTokenRequest true "payload"
// @Success 200 {object} auth.TokenResponse "ok"
// @Failure 400 {object} common.ErrorResponse "Invalid payload"
// @Failure 401 {object} common.ErrorResponse "Invalid credential"
// @Tags Auth API
// @Router /api/oauth2/token [post]
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