package auth

// CreateAccessTokenRequest is request model for creating todo.
type CreateAccessTokenRequest struct {
	UserName string `json:"user_name" example:"<user name>" binding:"required"`
	Password string `json:"password" example:"<password>" binding:"required"`
}

// TokenResponse is token model
type TokenResponse struct {
	Type         string `json:"type"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}
