package auth

// CreateAccessTokenRequest is request model for creating todo
type CreateAccessTokenRequest struct {
	UserName string `json:"username" example:"<username>" binding:"required"`
	Password string `json:"password" example:"<password>" binding:"required"`
}

// AccessTokenByRefreshRequest is request model
type AccessTokenByRefreshRequest struct {
	Token string `json:"token" example:"<refresh token>" binding:"required"`
}

// TokenResponse is token model
type TokenResponse struct {
	Type         string `json:"type"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresIn    int64  `json:"expires_in"`
}
