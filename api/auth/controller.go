package auth

import (
	"net/http"
	"strconv"

	"github.com/gghcode/go-gin-starterkit/api/common"
	"github.com/gghcode/go-gin-starterkit/config"
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
func NewController(conf config.Configuration, service Service) *Controller {
	return &Controller{
		conf:    conf.Jwt,
		service: service,
	}
}

// RegisterRoutes register handler routes.
func (controller *Controller) RegisterRoutes(router gin.IRouter) {
	router.Handle("POST", APIPath+"/token", controller.issueToken)
	router.Handle("POST", APIPath+"/refresh", controller.refreshToken)
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
func (controller *Controller) issueToken(ctx *gin.Context) {
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

// @Description Get new access token by refreshtoken
// @Accept json
// @Produce json
// @Param payload body auth.AccessTokenByRefreshRequest true "payload"
// @Success 200 {object} auth.TokenResponse "ok"
// @Failure 400 {object} common.ErrorResponse "Invalid payload"
// @Failure 401 {object} common.ErrorResponse "Invalid credential"
// @Tags Auth API
// @Router /auth/refresh [post]
func (controller *Controller) refreshToken(ctx *gin.Context) {
	var reqPayload AccessTokenByRefreshRequest

	if err := ctx.ShouldBindJSON(&reqPayload); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			common.NewErrResp(common.ErrInvalidRequestPayload),
		)
		return
	}

	tokenClaims, err := controller.service.ExtractTokenClaims(reqPayload.Token)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, common.NewErrResp(err))
		return
	}

	userID, _ := strconv.ParseInt(tokenClaims["sub"].(string), 10, 64)
	if valid := controller.service.VerifyRefreshToken(userID, reqPayload.Token); !valid {
		ctx.JSON(http.StatusUnauthorized, common.NewErrResp(ErrInvalidRefreshToken))
		return
	}

	accessToken, _ := controller.service.GenerateAccessToken(userID)

	res := TokenResponse{
		Type:        "Bearer",
		AccessToken: accessToken,
		ExpiresIn:   controller.conf.AccessExpiresInSec,
	}

	ctx.JSON(http.StatusOK, res)
}
